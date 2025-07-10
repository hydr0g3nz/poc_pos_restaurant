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

type paymentRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewPaymentRepository(db *pgxpool.Pool) repository.PaymentRepository {
	return &paymentRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	dbPayment, err := r.queries.CreatePayment(ctx, sqlc.CreatePaymentParams{
		OrderID:   int32(payment.OrderID),
		Amount:    utils.ConvertToPGNumericFromFloat(payment.Amount.Amount()),
		Method:    sqlc.PaymentMethod(payment.Method.String()),
		Reference: utils.ConvertToText(payment.Reference),
	})
	if err != nil {
		return nil, err
	}

	return r.dbPaymentToEntity(dbPayment)
}

func (r *paymentRepository) GetByID(ctx context.Context, id int) (*entity.Payment, error) {
	dbPayment, err := r.queries.GetPaymentByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbPaymentToEntity(dbPayment)
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID int) (*entity.Payment, error) {
	dbPayment, err := r.queries.GetPaymentByOrderID(ctx, int32(orderID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbPaymentToEntity(dbPayment)
}

func (r *paymentRepository) Update(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	dbPayment, err := r.queries.UpdatePayment(ctx, sqlc.UpdatePaymentParams{
		ID:        int32(payment.ID),
		OrderID:   int32(payment.OrderID),
		Amount:    utils.ConvertToPGNumericFromFloat(payment.Amount.Amount()),
		Method:    sqlc.PaymentMethod(payment.Method.String()),
		Reference: utils.ConvertToText(payment.Reference),
	})
	if err != nil {
		return nil, err
	}

	return r.dbPaymentToEntity(dbPayment)
}

func (r *paymentRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeletePayment(ctx, int32(id))
}

func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Payment, error) {
	dbPayments, err := r.queries.ListPayments(ctx, sqlc.ListPaymentsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbPaymentsToEntities(dbPayments)
}

func (r *paymentRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Payment, error) {
	dbPayments, err := r.queries.ListPaymentsByDateRange(ctx, sqlc.ListPaymentsByDateRangeParams{
		PaidAt:   utils.ConvertToPGTimestamp(&startDate),
		PaidAt_2: utils.ConvertToPGTimestamp(&endDate),
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbPaymentsToEntities(dbPayments)
}

func (r *paymentRepository) ListByMethod(ctx context.Context, method string, limit, offset int) ([]*entity.Payment, error) {
	dbPayments, err := r.queries.ListPaymentsByMethod(ctx, sqlc.ListPaymentsByMethodParams{
		Method: sqlc.PaymentMethod(method),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return r.dbPaymentsToEntities(dbPayments)
}

// Helper methods for conversion
func (r *paymentRepository) dbPaymentToEntity(dbPayment *sqlc.Payment) (*entity.Payment, error) {
	money, err := vo.NewMoney(utils.FromPgNumericToFloat(dbPayment.Amount))
	if err != nil {
		return nil, err
	}

	method, err := vo.NewPaymentMethod(string(dbPayment.Method))
	if err != nil {
		return nil, err
	}

	return &entity.Payment{
		ID:        int(dbPayment.ID),
		OrderID:   int(dbPayment.OrderID),
		Amount:    money,
		Method:    method,
		Reference: utils.FromPgTextToString(dbPayment.Reference),
		PaidAt:    dbPayment.PaidAt.Time,
	}, nil
}

func (r *paymentRepository) dbPaymentsToEntities(dbPayments []*sqlc.Payment) ([]*entity.Payment, error) {
	entities := make([]*entity.Payment, len(dbPayments))
	for i, dbPayment := range dbPayments {
		entity, err := r.dbPaymentToEntity(dbPayment)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
