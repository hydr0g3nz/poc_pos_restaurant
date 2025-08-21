// internal/adapter/repository/payment_repository.go
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

type paymentRepository struct {
	baseRepository
}

func NewPaymentRepository(db *gorm.DB) repository.PaymentRepository {
	return &paymentRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	dbPayment := r.entityToModel(payment)

	if err := r.db.WithContext(ctx).Create(dbPayment).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbPayment)
}

func (r *paymentRepository) GetByID(ctx context.Context, id int) (*entity.Payment, error) {
	var dbPayment model.Payment

	if err := r.db.WithContext(ctx).First(&dbPayment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbPayment)
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID int) (*entity.Payment, error) {
	var dbPayment model.Payment

	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&dbPayment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbPayment)
}

func (r *paymentRepository) Update(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	dbPayment := r.entityToModel(payment)

	if err := r.db.WithContext(ctx).Save(dbPayment).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbPayment)
}

func (r *paymentRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Payment{}, id).Error
}

func (r *paymentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Payment, error) {
	var dbPayments []model.Payment

	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbPayments).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbPayments)
}

func (r *paymentRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Payment, error) {
	var dbPayments []model.Payment

	query := r.db.WithContext(ctx).Where("paid_at BETWEEN ? AND ?", startDate, endDate)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbPayments).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbPayments)
}

func (r *paymentRepository) ListByMethod(ctx context.Context, method string, limit, offset int) ([]*entity.Payment, error) {
	var dbPayments []model.Payment

	query := r.db.WithContext(ctx).Where("method = ?", method)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbPayments).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbPayments)
}

// Helper methods
func (r *paymentRepository) entityToModel(payment *entity.Payment) *model.Payment {
	return &model.Payment{
		ID:        payment.ID,
		OrderID:   payment.OrderID,
		Amount:    payment.Amount.AmountSatang(),
		Method:    payment.Method.String(),
		Reference: payment.Reference,
		PaidAt:    payment.PaidAt,
	}
}

func (r *paymentRepository) modelToEntity(dbPayment *model.Payment) (*entity.Payment, error) {
	amount, err := vo.NewMoneyFromSatang(dbPayment.Amount)
	if err != nil {
		return nil, err
	}

	method, err := vo.NewPaymentMethod(dbPayment.Method)
	if err != nil {
		return nil, err
	}

	return &entity.Payment{
		ID:        dbPayment.ID,
		OrderID:   dbPayment.OrderID,
		Amount:    amount,
		Method:    method,
		Reference: dbPayment.Reference,
		PaidAt:    dbPayment.PaidAt,
	}, nil
}

func (r *paymentRepository) modelsToEntities(dbPayments []model.Payment) ([]*entity.Payment, error) {
	entities := make([]*entity.Payment, len(dbPayments))
	for i, dbPayment := range dbPayments {
		entity, err := r.modelToEntity(&dbPayment)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
