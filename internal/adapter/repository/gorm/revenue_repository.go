// internal/adapter/repository/revenue_repository.go
package repository

import (
	"context"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type revenueRepository struct {
	baseRepository
}

func NewRevenueRepository(db *gorm.DB) repository.RevenueRepository {
	return &revenueRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *revenueRepository) GetDailyRevenue(ctx context.Context, date time.Time) (*entity.DailyRevenue, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalAmount int64

	err := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("paid_at >= ? AND paid_at < ?", startOfDay, endOfDay).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount).Error

	if err != nil {
		return nil, err
	}

	revenue, err := vo.NewMoneyFromSatang(totalAmount)
	if err != nil {
		return nil, err
	}

	return &entity.DailyRevenue{
		Date:         startOfDay,
		TotalRevenue: revenue,
	}, nil
}

func (r *revenueRepository) GetMonthlyRevenue(ctx context.Context, year int, month int) (*entity.MonthlyRevenue, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var totalAmount int64

	err := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("paid_at >= ? AND paid_at < ?", startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount).Error

	if err != nil {
		return nil, err
	}

	revenue, err := vo.NewMoneyFromSatang(totalAmount)
	if err != nil {
		return nil, err
	}

	return &entity.MonthlyRevenue{
		Month:        startOfMonth,
		TotalRevenue: revenue,
	}, nil
}

func (r *revenueRepository) GetDailyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyRevenue, error) {
	type DailyRevenueResult struct {
		Date   time.Time
		Amount int64
	}

	var results []DailyRevenueResult

	err := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Select("DATE(paid_at) as date, COALESCE(SUM(amount), 0) as amount").
		Where("paid_at >= ? AND paid_at <= ?", startDate, endDate).
		Group("DATE(paid_at)").
		Order("date").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	revenues := make([]*entity.DailyRevenue, len(results))
	for i, result := range results {
		revenue, err := vo.NewMoneyFromSatang(result.Amount)
		if err != nil {
			return nil, err
		}

		revenues[i] = &entity.DailyRevenue{
			Date:         result.Date,
			TotalRevenue: revenue,
		}
	}

	return revenues, nil
}

func (r *revenueRepository) GetMonthlyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.MonthlyRevenue, error) {
	type MonthlyRevenueResult struct {
		Month  time.Time
		Amount int64
	}

	var results []MonthlyRevenueResult

	err := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Select("DATE_TRUNC('month', paid_at) as month, COALESCE(SUM(amount), 0) as amount").
		Where("paid_at >= ? AND paid_at <= ?", startDate, endDate).
		Group("DATE_TRUNC('month', paid_at)").
		Order("month").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	revenues := make([]*entity.MonthlyRevenue, len(results))
	for i, result := range results {
		revenue, err := vo.NewMoneyFromSatang(result.Amount)
		if err != nil {
			return nil, err
		}

		revenues[i] = &entity.MonthlyRevenue{
			Month:        result.Month,
			TotalRevenue: revenue,
		}
	}

	return revenues, nil
}

func (r *revenueRepository) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	var totalAmount int64

	err := r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where("paid_at >= ? AND paid_at <= ?", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}

	revenue, err := vo.NewMoneyFromSatang(totalAmount)
	if err != nil {
		return 0, err
	}

	return revenue.AmountBaht(), nil
}
