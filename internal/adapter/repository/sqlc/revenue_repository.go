// internal/adapter/repository/sqlc/revenue_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlc "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/sqlc/generated"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"github.com/hydr0g3nz/poc_pos_restuarant/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type revenueRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRevenueRepository(db *pgxpool.Pool) repository.RevenueRepository {
	return &revenueRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *revenueRepository) GetDailyRevenue(ctx context.Context, date time.Time) (*entity.DailyRevenue, error) {
	// Calculate start and end of day
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endDate := startDate.Add(24 * time.Hour)

	// Get revenue from payments table
	totalRevenue, err := r.queries.GetDailyRevenue(ctx, sqlc.GetDailyRevenueParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert to entity
	revenueAmount, err := vo.NewMoney(utils.FromInterfaceToFloat(totalRevenue))
	if err != nil {
		return nil, err
	}

	return &entity.DailyRevenue{
		Date:         date,
		TotalRevenue: revenueAmount,
	}, nil
}

func (r *revenueRepository) GetMonthlyRevenue(ctx context.Context, year int, month int) (*entity.MonthlyRevenue, error) {
	// Calculate start and end of month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Next month

	// Get revenue from payments table
	totalRevenue, err := r.queries.GetMonthlyRevenue(ctx, sqlc.GetMonthlyRevenueParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert to entity
	revenueAmount, err := vo.NewMoney(utils.FromInterfaceToFloat(totalRevenue))
	if err != nil {
		return nil, err
	}

	return &entity.MonthlyRevenue{
		Month:        startDate,
		TotalRevenue: revenueAmount,
	}, nil
}

func (r *revenueRepository) GetDailyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyRevenue, error) {
	// Get daily revenue breakdown
	results, err := r.queries.GetDailyRevenueRange(ctx, sqlc.GetDailyRevenueRangeParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
	})
	if err != nil {
		return nil, err
	}

	// Convert to entities
	revenues := make([]*entity.DailyRevenue, len(results))
	for i, result := range results {
		totalRevenue, err := vo.NewMoney(utils.FromInterfaceToFloat(result.TotalRevenue))
		if err != nil {
			return nil, err
		}

		revenues[i] = &entity.DailyRevenue{
			Date:         result.RevenueDate.Time,
			TotalRevenue: totalRevenue,
		}
	}

	return revenues, nil
}

func (r *revenueRepository) GetMonthlyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.MonthlyRevenue, error) {
	// Get monthly revenue breakdown
	results, err := r.queries.GetMonthlyRevenueRange(ctx, sqlc.GetMonthlyRevenueRangeParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
	})
	if err != nil {
		return nil, err
	}

	// Convert to entities
	revenues := make([]*entity.MonthlyRevenue, len(results))
	for i, result := range results {
		totalRevenue, err := vo.NewMoney(utils.FromInterfaceToFloat(result.TotalRevenue))
		if err != nil {
			return nil, err
		}

		revenues[i] = &entity.MonthlyRevenue{
			Month:        result.RevenueMonth.Time,
			TotalRevenue: totalRevenue,
		}
	}

	return revenues, nil
}

func (r *revenueRepository) GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	// Get total revenue for period
	result, err := r.queries.GetTotalRevenue(ctx, sqlc.GetTotalRevenueParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return utils.FromPgNumericToFloat(result), nil
}
