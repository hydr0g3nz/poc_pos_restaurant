// internal/adapter/repository/option_value_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type optionValueRepository struct {
	baseRepository
}

func NewOptionValueRepository(db *gorm.DB) repository.OptionValueRepository {
	return &optionValueRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *optionValueRepository) Create(ctx context.Context, value *entity.OptionValue) (*entity.OptionValue, error) {
	dbValue := r.entityToModel(value)

	if err := r.db.WithContext(ctx).Create(dbValue).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbValue)
}

func (r *optionValueRepository) GetByID(ctx context.Context, id int) (*entity.OptionValue, error) {
	var dbValue model.OptionValue

	if err := r.db.WithContext(ctx).First(&dbValue, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbValue)
}

func (r *optionValueRepository) GetByOptionID(ctx context.Context, optionID int) ([]*entity.OptionValue, error) {
	var dbValues []model.OptionValue

	if err := r.db.WithContext(ctx).Where("option_id = ?", optionID).Find(&dbValues).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbValues)
}

func (r *optionValueRepository) Update(ctx context.Context, value *entity.OptionValue) (*entity.OptionValue, error) {
	dbValue := r.entityToModel(value)

	if err := r.db.WithContext(ctx).Save(dbValue).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbValue)
}

func (r *optionValueRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.OptionValue{}, id).Error
}

func (r *optionValueRepository) List(ctx context.Context, limit, offset int) ([]*entity.OptionValue, error) {
	var dbValues []model.OptionValue

	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbValues).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbValues)
}

// Helper methods
func (r *optionValueRepository) entityToModel(value *entity.OptionValue) *model.OptionValue {
	return &model.OptionValue{
		ID:              value.ID,
		OptionID:        value.OptionID,
		Name:            value.Name,
		IsDefault:       value.IsDefault,
		AdditionalPrice: value.AdditionalPrice.AmountSatang(),
		DisplayOrder:    value.DisplayOrder,
	}
}

func (r *optionValueRepository) modelToEntity(dbValue *model.OptionValue) (*entity.OptionValue, error) {
	additionalPrice, err := vo.NewMoneyFromSatang(dbValue.AdditionalPrice)
	if err != nil {
		return nil, err
	}

	return &entity.OptionValue{
		ID:              dbValue.ID,
		OptionID:        dbValue.OptionID,
		Name:            dbValue.Name,
		IsDefault:       dbValue.IsDefault,
		AdditionalPrice: additionalPrice,
		DisplayOrder:    dbValue.DisplayOrder,
	}, nil
}

func (r *optionValueRepository) modelsToEntities(dbValues []model.OptionValue) ([]*entity.OptionValue, error) {
	entities := make([]*entity.OptionValue, len(dbValues))
	for i, dbValue := range dbValues {
		entity, err := r.modelToEntity(&dbValue)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
