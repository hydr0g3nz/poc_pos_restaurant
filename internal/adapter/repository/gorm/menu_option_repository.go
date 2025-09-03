// internal/adapter/repository/menu_option_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type menuOptionRepository struct {
	baseRepository
}

func NewMenuOptionRepository(db *gorm.DB) repository.MenuOptionRepository {
	return &menuOptionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *menuOptionRepository) Create(ctx context.Context, option *entity.MenuOption) (*entity.MenuOption, error) {
	dbOption := r.entityToModel(option)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Create(dbOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOption), nil
}

func (r *menuOptionRepository) GetByID(ctx context.Context, id int) (*entity.MenuOption, error) {
	var dbOption model.MenuOption

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).First(&dbOption, id).Error; err != nil {
		// if err == gorm.ErrRecordNotFound {
		// 	return nil, nil
		// }
		return nil, err
	}

	return r.modelToEntity(&dbOption), nil
}

func (r *menuOptionRepository) Update(ctx context.Context, option *entity.MenuOption) (*entity.MenuOption, error) {
	dbOption := r.entityToModel(option)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Save(dbOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbOption), nil
}

func (r *menuOptionRepository) Delete(ctx context.Context, id int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Delete(&model.MenuOption{}, id).Error
}

func (r *menuOptionRepository) List(ctx context.Context, limit, offset int) ([]*entity.MenuOption, error) {
	var dbOptions []model.MenuOption

	db := getDB(r.db, ctx)
	query := db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOptions), nil
}

func (r *menuOptionRepository) GetByType(ctx context.Context, optionType string) ([]*entity.MenuOption, error) {
	var dbOptions []model.MenuOption

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("type = ?", optionType).Find(&dbOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbOptions), nil
}

// Helper methods
func (r *menuOptionRepository) entityToModel(option *entity.MenuOption) *model.MenuOption {
	return &model.MenuOption{
		ID:         option.ID,
		Name:       option.Name,
		Type:       option.Type.String(),
		IsRequired: option.IsRequired,
	}
}

func (r *menuOptionRepository) modelToEntity(dbOption *model.MenuOption) *entity.MenuOption {
	return &entity.MenuOption{
		ID:         dbOption.ID,
		Name:       dbOption.Name,
		Type:       vo.OptionType(dbOption.Type),
		IsRequired: dbOption.IsRequired,
	}
}

func (r *menuOptionRepository) modelsToEntities(dbOptions []model.MenuOption) []*entity.MenuOption {
	entities := make([]*entity.MenuOption, len(dbOptions))
	for i, dbOption := range dbOptions {
		entities[i] = r.modelToEntity(&dbOption)
	}
	return entities
}
