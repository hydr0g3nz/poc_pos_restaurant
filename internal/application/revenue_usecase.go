// internal/application/revenue_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
)

// revenueUsecase implements RevenueUsecase interface
type revenueUsecase struct {
	revenueRepo repository.RevenueRepository
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	logger      infra.Logger
	config      *config.Config
}

// NewRevenueUsecase creates a new revenue usecase
func NewRevenueUsecase(
	revenueRepo repository.RevenueRepository,
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	logger infra.Logger,
	config *config.Config,
) RevenueUsecase {
	return &revenueUsecase{
		revenueRepo: revenueRepo,
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		logger:      logger,
		config:      config,
	}
}

// GetDailyRevenue retrieves revenue for a specific date
func (u *revenueUsecase) GetDailyRevenue(ctx context.Context, date time.Time) (*DailyRevenueResponse, error) {
	u.logger.Debug("Getting daily revenue", "date", date.Format("2006-01-02"))

	// Validate date
	if date.After(time.Now()) {
		u.logger.Error("Invalid date - future date", "date", date)
		return nil, errs.ErrInvalidRevenueDate
	}

	// Get daily revenue from repository
	dailyRevenue, err := u.revenueRepo.GetDailyRevenue(ctx, date)
	if err != nil {
		u.logger.Error("Error getting daily revenue", "error", err, "date", date)
		return nil, fmt.Errorf("failed to get daily revenue: %w", err)
	}

	// If no revenue data found, return zero revenue
	if dailyRevenue == nil {
		u.logger.Info("No revenue data found for date", "date", date)
		return &DailyRevenueResponse{
			Date:         date,
			TotalRevenue: 0,
			OrderCount:   0,
		}, nil
	}

	// Get order count for the day
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)
	orders, err := u.orderRepo.ListByDateRange(ctx, startDate, endDate, 1000, 0) // Get all orders for count
	if err != nil {
		u.logger.Error("Error getting order count", "error", err, "date", date)
		// Don't fail the request, just set order count to 0
	}

	orderCount := len(orders)

	return &DailyRevenueResponse{
		Date:         dailyRevenue.Date,
		TotalRevenue: dailyRevenue.TotalRevenue.AmountBaht(),
		OrderCount:   orderCount,
	}, nil
}

// GetMonthlyRevenue retrieves revenue for a specific month
func (u *revenueUsecase) GetMonthlyRevenue(ctx context.Context, year int, month int) (*MonthlyRevenueResponse, error) {
	u.logger.Debug("Getting monthly revenue", "year", year, "month", month)

	// Validate year and month
	if year < 2000 || year > time.Now().Year()+1 {
		u.logger.Error("Invalid year", "year", year)
		return nil, errs.ErrInvalidRevenueDate
	}
	if month < 1 || month > 12 {
		u.logger.Error("Invalid month", "month", month)
		return nil, errs.ErrInvalidRevenueDate
	}

	// Get monthly revenue from repository
	monthlyRevenue, err := u.revenueRepo.GetMonthlyRevenue(ctx, year, month)
	if err != nil {
		u.logger.Error("Error getting monthly revenue", "error", err, "year", year, "month", month)
		return nil, fmt.Errorf("failed to get monthly revenue: %w", err)
	}

	// If no revenue data found, return zero revenue
	if monthlyRevenue == nil {
		monthDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		u.logger.Info("No revenue data found for month", "year", year, "month", month)
		return &MonthlyRevenueResponse{
			Month:        monthDate,
			TotalRevenue: 0,
			OrderCount:   0,
		}, nil
	}

	// Get order count for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)                                         // Next month
	orders, err := u.orderRepo.ListByDateRange(ctx, startDate, endDate, 10000, 0) // Get all orders for count
	if err != nil {
		u.logger.Error("Error getting order count", "error", err, "year", year, "month", month)
		// Don't fail the request, just set order count to 0
	}

	orderCount := len(orders)

	return &MonthlyRevenueResponse{
		Month:        monthlyRevenue.Month,
		TotalRevenue: monthlyRevenue.TotalRevenue.AmountBaht(),
		OrderCount:   orderCount,
	}, nil
}

// GetDailyRevenueRange retrieves daily revenue for a date range
func (u *revenueUsecase) GetDailyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*DailyRevenueResponse, error) {
	u.logger.Debug("Getting daily revenue range", "startDate", startDate, "endDate", endDate)

	// Validate date range
	if startDate.After(endDate) {
		u.logger.Error("Invalid date range", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidDateRange
	}

	// Validate dates are not in the future
	now := time.Now()
	if startDate.After(now) || endDate.After(now) {
		u.logger.Error("Invalid date range - future dates", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidRevenueDate
	}

	// Get daily revenue range from repository
	dailyRevenues, err := u.revenueRepo.GetDailyRevenueRange(ctx, startDate, endDate)
	if err != nil {
		u.logger.Error("Error getting daily revenue range", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, fmt.Errorf("failed to get daily revenue range: %w", err)
	}

	// Convert to response format
	responses := make([]*DailyRevenueResponse, len(dailyRevenues))
	for i, revenue := range dailyRevenues {
		// Get order count for each day
		dayStart := time.Date(revenue.Date.Year(), revenue.Date.Month(), revenue.Date.Day(), 0, 0, 0, 0, revenue.Date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		orders, err := u.orderRepo.ListByDateRange(ctx, dayStart, dayEnd, 1000, 0)
		orderCount := 0
		if err == nil {
			orderCount = len(orders)
		}

		responses[i] = &DailyRevenueResponse{
			Date:         revenue.Date,
			TotalRevenue: revenue.TotalRevenue.AmountBaht(),
			OrderCount:   orderCount,
		}
	}

	return responses, nil
}

// GetMonthlyRevenueRange retrieves monthly revenue for a date range
func (u *revenueUsecase) GetMonthlyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*MonthlyRevenueResponse, error) {
	u.logger.Debug("Getting monthly revenue range", "startDate", startDate, "endDate", endDate)

	// Validate date range
	if startDate.After(endDate) {
		u.logger.Error("Invalid date range", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidDateRange
	}

	// Get monthly revenue range from repository
	monthlyRevenues, err := u.revenueRepo.GetMonthlyRevenueRange(ctx, startDate, endDate)
	if err != nil {
		u.logger.Error("Error getting monthly revenue range", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, fmt.Errorf("failed to get monthly revenue range: %w", err)
	}

	// Convert to response format
	responses := make([]*MonthlyRevenueResponse, len(monthlyRevenues))
	for i, revenue := range monthlyRevenues {
		// Get order count for each month
		monthStart := time.Date(revenue.Month.Year(), revenue.Month.Month(), 1, 0, 0, 0, 0, revenue.Month.Location())
		monthEnd := monthStart.AddDate(0, 1, 0)
		orders, err := u.orderRepo.ListByDateRange(ctx, monthStart, monthEnd, 10000, 0)
		orderCount := 0
		if err == nil {
			orderCount = len(orders)
		}

		responses[i] = &MonthlyRevenueResponse{
			Month:        revenue.Month,
			TotalRevenue: revenue.TotalRevenue.AmountBaht(),
			OrderCount:   orderCount,
		}
	}

	return responses, nil
}

// GetTotalRevenue retrieves total revenue for a date range
func (u *revenueUsecase) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (*TotalRevenueResponse, error) {
	u.logger.Debug("Getting total revenue", "startDate", startDate, "endDate", endDate)

	// Validate date range
	if startDate.After(endDate) {
		u.logger.Error("Invalid date range", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidDateRange
	}

	// Get total revenue from repository
	totalRevenue, err := u.revenueRepo.GetTotalRevenue(ctx, startDate, endDate)
	if err != nil {
		u.logger.Error("Error getting total revenue", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	// Get order count for the period
	orders, err := u.orderRepo.ListByDateRange(ctx, startDate, endDate, 100000, 0) // Get all orders for count
	orderCount := 0
	if err == nil {
		orderCount = len(orders)
	} else {
		u.logger.Error("Error getting order count", "error", err, "startDate", startDate, "endDate", endDate)
		// Don't fail the request, just set order count to 0
	}

	return &TotalRevenueResponse{
		StartDate:    startDate,
		EndDate:      endDate,
		TotalRevenue: totalRevenue,
		OrderCount:   orderCount,
	}, nil
}
