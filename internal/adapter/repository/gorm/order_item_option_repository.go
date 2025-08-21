// internal/adapter/repository/order_item_option_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type orderItemOptionRepository struct {
	baseRepository
}

func NewOrderItemOptionRepository(db *gorm.DB) repository.OrderItemOptionRepository {
	return &orderItemOptionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *orderItemOptionRepository) Create(ctx context.Context, itemOption *entity.OrderItemOption) (*entity.OrderItemOption, error) {
	dbItemOption := r.entityToModel(itemOption)

	if err := r.db.WithContext(ctx).Create(dbItemOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItemOption)
}

func (r *orderItemOptionRepository) GetByOrderItemID(ctx context.Context, orderItemID int) ([]*entity.OrderItemOption, error) {
	var dbItemOptions []model.OrderItemOption

	if err := r.db.WithContext(ctx).Where("order_item_id = ?", orderItemID).Find(&dbItemOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItemOptions)
}

func (r *orderItemOptionRepository) Update(ctx context.Context, itemOption *entity.OrderItemOption) (*entity.OrderItemOption, error) {
	dbItemOption := r.entityToModel(itemOption)

	if err := r.db.WithContext(ctx).Save(dbItemOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItemOption)
}

func (r *orderItemOptionRepository) Delete(ctx context.Context, orderItemID, optionID, valueID int) error {
	return r.db.WithContext(ctx).Where("order_item_id = ? AND option_id = ? AND value_id = ?", orderItemID, optionID, valueID).Delete(&model.OrderItemOption{}).Error
}

func (r *orderItemOptionRepository) DeleteByOrderItemID(ctx context.Context, orderItemID int) error {
	return r.db.WithContext(ctx).Where("order_item_id = ?", orderItemID).Delete(&model.OrderItemOption{}).Error
}

// Helper methods
func (r *orderItemOptionRepository) entityToModel(itemOption *entity.OrderItemOption) *model.OrderItemOption {
	return &model.OrderItemOption{
		OrderItemID:     itemOption.OrderItemID,
		OptionID:        itemOption.OptionID,
		ValueID:         itemOption.ValueID,
		AdditionalPrice: itemOption.AdditionalPrice.AmountSatang(),
	}
}

func (r *orderItemOptionRepository) modelToEntity(dbItemOption *model.OrderItemOption) (*entity.OrderItemOption, error) {
	additionalPrice, err := vo.NewMoneyFromSatang(dbItemOption.AdditionalPrice)
	if err != nil {
		return nil, err
	}

	return &entity.OrderItemOption{
		OrderItemID:     dbItemOption.OrderItemID,
		OptionID:        dbItemOption.OptionID,
		ValueID:         dbItemOption.ValueID,
		AdditionalPrice: additionalPrice,
	}, nil
}

func (r *orderItemOptionRepository) modelsToEntities(dbItemOptions []model.OrderItemOption) ([]*entity.OrderItemOption, error) {
	entities := make([]*entity.OrderItemOption, len(dbItemOptions))
	for i, dbItemOption := range dbItemOptions {
		entity, err := r.modelToEntity(&dbItemOption)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
