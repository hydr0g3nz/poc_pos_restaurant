// internal/adapter/repository/menu_option_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"gorm.io/gorm"
)

type KitchenStationRepository struct {
	baseRepository
}

func NewKitchenStationRepository(db *gorm.DB) repository.KitchenStationRepository {
	return &KitchenStationRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *KitchenStationRepository) Create(ctx context.Context, option *entity.KitchenStation) (*entity.KitchenStation, error) {
	dbOption := r.entityToModel(option)

	if err := r.db.WithContext(ctx).Create(dbOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOption), nil
}

func (r *KitchenStationRepository) GetByID(ctx context.Context, id int) (*entity.KitchenStation, error) {
	var dbOption model.KitchenStation

	if err := r.db.WithContext(ctx).First(&dbOption, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbOption), nil
}

func (r *KitchenStationRepository) Update(ctx context.Context, option *entity.KitchenStation) (*entity.KitchenStation, error) {
	dbOption := r.entityToModel(option)

	if err := r.db.WithContext(ctx).Save(dbOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOption), nil
}

func (r *KitchenStationRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.KitchenStation{}, id).Error
}

func (r *KitchenStationRepository) List(ctx context.Context, onlyAvailable bool, limit, offset int) ([]*entity.KitchenStation, error) {
	var dbOptions []model.KitchenStation

	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if onlyAvailable {
		query = query.Where("is_available = ?", true)
	}
	if err := query.Find(&dbOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOptions), nil
}

func (r *KitchenStationRepository) GetByType(ctx context.Context, optionType string) ([]*entity.KitchenStation, error) {
	var dbOptions []model.KitchenStation

	if err := r.db.WithContext(ctx).Where("type = ?", optionType).Find(&dbOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOptions), nil
}

// Helper methods
func (r *KitchenStationRepository) entityToModel(option *entity.KitchenStation) *model.KitchenStation {
	return &model.KitchenStation{
		ID:          option.ID,
		Name:        option.Name,
		IsAvailable: option.IsAvailable,
	}
}

func (r *KitchenStationRepository) modelToEntity(dbOption *model.KitchenStation) *entity.KitchenStation {
	return &entity.KitchenStation{
		ID:          dbOption.ID,
		Name:        dbOption.Name,
		IsAvailable: dbOption.IsAvailable,
	}
}

func (r *KitchenStationRepository) modelsToEntities(dbOptions []model.KitchenStation) []*entity.KitchenStation {
	entities := make([]*entity.KitchenStation, len(dbOptions))
	for i, dbOption := range dbOptions {
		entities[i] = r.modelToEntity(&dbOption)
	}
	return entities
}
