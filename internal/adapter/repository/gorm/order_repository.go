// internal/adapter/repository/order_repository.go
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

type orderRepository struct {
	baseRepository
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	dbOrder := r.entityToModel(order)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Create(dbOrder).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOrder)
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (*entity.Order, error) {
	var dbOrder model.Order

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).First(&dbOrder, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbOrder)
}

func (r *orderRepository) GetByIDWithItems(ctx context.Context, id int) (*entity.Order, error) {
	var dbOrder model.Order

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Preload("OrderItems").First(&dbOrder, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntityWithItems(&dbOrder)
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	dbOrder := r.entityToModel(order)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Save(dbOrder).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOrder)
}

func (r *orderRepository) Delete(ctx context.Context, id int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Delete(&model.Order{}, id).Error
}

func (r *orderRepository) List(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	var dbOrders []model.Order

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOrders)
}
func (r *orderRepository) ListWithItems(ctx context.Context, limit, offset int) ([]*entity.Order, error) {
	var dbOrders []model.Order

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("OrderItems").Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOrders)
}

func (r *orderRepository) ListByTable(ctx context.Context, tableID int, limit, offset int) ([]*entity.Order, error) {
	var dbOrders []model.Order

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx).Where("table_id = ?", tableID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOrders)
}

func (r *orderRepository) GetOpenOrderByTable(ctx context.Context, tableID int) (*entity.Order, error) {
	var dbOrder model.Order

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("table_id = ? AND order_status = ?", tableID, "open").First(&dbOrder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbOrder)
}

func (r *orderRepository) GetOrderByQRCode(ctx context.Context, qrCode string) (*entity.Order, error) {
	var dbOrder model.Order

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("qr_code = ?", qrCode).First(&dbOrder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbOrder)
}
func (r *orderRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Order, error) {
	var dbOrders []model.Order

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx).Where("order_status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOrders)
}

func (r *orderRepository) ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Order, error) {
	var dbOrders []model.Order

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx).Where("created_at BETWEEN ? AND ?", startDate, endDate)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbOrders).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOrders)
}

// Helper methods
func (r *orderRepository) entityToModel(order *entity.Order) *model.Order {
	return &model.Order{
		ID:                  order.ID,
		OrderNumber:         order.OrderNumber,
		TableID:             order.TableID,
		OrderStatus:         order.OrderStatus.String(),
		PaymentStatus:       order.PaymentStatus.String(),
		QRCode:              order.QRCode,
		Notes:               order.Notes,
		SpecialInstructions: order.SpecialInstructions,
		Subtotal:            order.Subtotal.AmountSatang(),
		Discount:            order.Discount.AmountSatang(),
		TaxAmount:           order.TaxAmount.AmountSatang(),
		ServiceCharge:       order.ServiceCharge.AmountSatang(),
		Total:               order.Total.AmountSatang(),
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
		ClosedAt:            order.ClosedAt,
	}
}

func (r *orderRepository) modelToEntity(dbOrder *model.Order) (*entity.Order, error) {
	orderStatus, err := vo.NewOrderStatus(dbOrder.OrderStatus)
	if err != nil {
		return nil, err
	}

	paymentStatus, err := vo.NewPaymentStatus(dbOrder.PaymentStatus)
	if err != nil {
		return nil, err
	}

	subtotal, err := vo.NewMoneyFromSatang(dbOrder.Subtotal)
	if err != nil {
		return nil, err
	}

	discount, err := vo.NewMoneyFromSatang(dbOrder.Discount)
	if err != nil {
		return nil, err
	}

	taxAmount, err := vo.NewMoneyFromSatang(dbOrder.TaxAmount)
	if err != nil {
		return nil, err
	}

	serviceCharge, err := vo.NewMoneyFromSatang(dbOrder.ServiceCharge)
	if err != nil {
		return nil, err
	}

	total, err := vo.NewMoneyFromSatang(dbOrder.Total)
	if err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:                  dbOrder.ID,
		OrderNumber:         dbOrder.OrderNumber,
		TableID:             dbOrder.TableID,
		OrderStatus:         orderStatus,
		PaymentStatus:       paymentStatus,
		QRCode:              dbOrder.QRCode,
		Notes:               dbOrder.Notes,
		SpecialInstructions: dbOrder.SpecialInstructions,
		Subtotal:            subtotal,
		Discount:            discount,
		TaxAmount:           taxAmount,
		ServiceCharge:       serviceCharge,
		Total:               total,
		CreatedAt:           dbOrder.CreatedAt,
		UpdatedAt:           dbOrder.UpdatedAt,
		ClosedAt:            dbOrder.ClosedAt,
	}, nil
}

func (r *orderRepository) modelToEntityWithItems(dbOrder *model.Order) (*entity.Order, error) {
	order, err := r.modelToEntity(dbOrder)
	if err != nil {
		return nil, err
	}

	// Convert order items
	orderItems := make([]*entity.OrderItem, len(dbOrder.OrderItems))
	for i, dbItem := range dbOrder.OrderItems {
		orderItem, err := r.orderItemModelToEntity(&dbItem)
		if err != nil {
			return nil, err
		}
		orderItems[i] = orderItem
	}
	order.Items = orderItems

	return order, nil
}

func (r *orderRepository) orderItemModelToEntity(dbItem *model.OrderItem) (*entity.OrderItem, error) {
	unitPrice, err := vo.NewMoneyFromSatang(dbItem.UnitPrice)
	if err != nil {
		return nil, err
	}

	discount, err := vo.NewMoneyFromSatang(dbItem.Discount)
	if err != nil {
		return nil, err
	}

	total, err := vo.NewMoneyFromSatang(dbItem.Total)
	if err != nil {
		return nil, err
	}

	itemStatus, err := vo.NewItemStatus(dbItem.ItemStatus)
	if err != nil {
		return nil, err
	}

	return &entity.OrderItem{
		ID:              dbItem.ID,
		OrderID:         dbItem.OrderID,
		ItemID:          dbItem.ItemID,
		Quantity:        dbItem.Quantity,
		UnitPrice:       unitPrice,
		Name:            dbItem.Name,
		Discount:        discount,
		Total:           total,
		SpecialReq:      dbItem.SpecialReq,
		ItemStatus:      itemStatus,
		OrderNumber:     dbItem.OrderNumber,
		KitchenTicketID: dbItem.KitchenTicketID,
		KitchenStation:  dbItem.KitchenStation,
		KitchenNotes:    dbItem.KitchenNotes,
		ServedAt:        dbItem.ServedAt,
		CreatedAt:       dbItem.CreatedAt,
		UpdatedAt:       dbItem.UpdatedAt,
	}, nil
}

func (r *orderRepository) modelsToEntities(dbOrders []model.Order) ([]*entity.Order, error) {
	entities := make([]*entity.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		var entity *entity.Order
		var err error
		if dbOrder.OrderItems != nil {
			entity, err = r.modelToEntityWithItems(&dbOrder)
			if err != nil {
				return nil, err
			}
		} else {
			entity, err = r.modelToEntity(&dbOrder)
			if err != nil {
				return nil, err
			}
		}
		entities[i] = entity
	}
	return entities, nil
}

// Count(ctx context.Context) (int, error)
//
//	CountByStatus(ctx context.Context, status string) (int, error)
//	CountByTable(ctx context.Context, tableID int) (int, error)
//	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int, error)
func (r *orderRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("order_status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
func (r *orderRepository) CountByTable(ctx context.Context, tableID int) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("table_id = ?", tableID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
func (r *orderRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
func (r *orderRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
func (r *orderRepository) GetOrderIDByQRCode(ctx context.Context, qrCode string) (int, error) {
	var dbOrder model.Order

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Select("id").Where("qr_code = ?", qrCode).First(&dbOrder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	return dbOrder.ID, nil
}
