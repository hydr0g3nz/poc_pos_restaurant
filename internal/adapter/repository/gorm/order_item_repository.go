// internal/adapter/repository/order_item_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	baseRepository
}

func NewOrderItemRepository(db *gorm.DB) repository.OrderItemRepository {
	return &orderItemRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *orderItemRepository) Create(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error) {
	dbItem := r.entityToModel(item)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Create(dbItem).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItem)
}

func (r *orderItemRepository) GetByID(ctx context.Context, id int) (*entity.OrderItem, error) {
	var dbItem model.OrderItem

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).First(&dbItem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbItem)
}

func (r *orderItemRepository) Update(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error) {
	dbItem := r.entityToModel(item)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Save(dbItem).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItem)
}

func (r *orderItemRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.OrderItem{}, id).Error
}

func (r *orderItemRepository) ListByOrder(ctx context.Context, orderID int) ([]*entity.OrderItem, error) {
	var dbItems []model.OrderItem
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("order_id = ?", orderID).Find(&dbItems).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItems)
}

func (r *orderItemRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&model.OrderItem{}).Error
}

func (r *orderItemRepository) GetByOrderAndItem(ctx context.Context, orderID, itemID int) (*entity.OrderItem, error) {
	var dbItem model.OrderItem

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("order_id = ? AND item_id = ?", orderID, itemID).First(&dbItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbItem)
}

// Helper methods
func (r *orderItemRepository) entityToModel(item *entity.OrderItem) *model.OrderItem {
	return &model.OrderItem{
		ID:              item.ID,
		OrderID:         item.OrderID,
		ItemID:          item.ItemID,
		Quantity:        item.Quantity,
		UnitPrice:       item.UnitPrice.AmountSatang(),
		Name:            item.Name,
		Discount:        item.Discount.AmountSatang(),
		Total:           item.Total.AmountSatang(),
		SpecialReq:      item.SpecialReq,
		ItemStatus:      item.ItemStatus.String(),
		OrderNumber:     item.OrderNumber,
		KitchenTicketID: item.KitchenTicketID,
		KitchenStation:  item.KitchenStation,
		KitchenNotes:    item.KitchenNotes,
		ServedAt:        item.ServedAt,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func (r *orderItemRepository) modelToEntity(dbItem *model.OrderItem) (*entity.OrderItem, error) {
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

func (r *orderItemRepository) modelsToEntities(dbItems []model.OrderItem) ([]*entity.OrderItem, error) {
	entities := make([]*entity.OrderItem, len(dbItems))
	for i, dbItem := range dbItems {
		entity, err := r.modelToEntity(&dbItem)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
