// internal/adapter/repository/menu_item_option_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type menuItemOptionRepository struct {
	baseRepository
}

func NewMenuItemOptionRepository(db *gorm.DB) repository.MenuItemOptionRepository {
	return &menuItemOptionRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *menuItemOptionRepository) Create(ctx context.Context, itemOption *entity.MenuItemOption) (*entity.MenuItemOption, error) {
	dbItemOption := r.entityToModel(itemOption)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "item_id"}, {Name: "option_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_active"}),
	}).Create(dbItemOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItemOption), nil
}

func (r *menuItemOptionRepository) GetByItemID(ctx context.Context, itemID int) ([]*entity.MenuItemOption, error) {
	var dbItemOptions []model.MenuItemOption

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("item_id = ?", itemID).Find(&dbItemOptions).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItemOptions), nil
}

func (r *menuItemOptionRepository) GetByItemAndOption(ctx context.Context, itemID, optionID int) (*entity.MenuItemOption, error) {
	var dbItemOption model.MenuItemOption

	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Where("item_id = ? AND option_id = ?", itemID, optionID).First(&dbItemOption).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbItemOption), nil
}

func (r *menuItemOptionRepository) Update(ctx context.Context, itemOption *entity.MenuItemOption) (*entity.MenuItemOption, error) {
	dbItemOption := r.entityToModel(itemOption)
	db := getDB(r.db, ctx)
	if err := db.WithContext(ctx).Save(dbItemOption).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItemOption), nil
}

func (r *menuItemOptionRepository) Delete(ctx context.Context, itemID, optionID int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Where("item_id = ? AND option_id = ?", itemID, optionID).Delete(&model.MenuItemOption{}).Error
}

func (r *menuItemOptionRepository) DeleteByItemID(ctx context.Context, itemID int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Where("item_id = ?", itemID).Delete(&model.MenuItemOption{}).Error
}

// Helper methods
func (r *menuItemOptionRepository) entityToModel(itemOption *entity.MenuItemOption) *model.MenuItemOption {
	return &model.MenuItemOption{
		ItemID:   itemOption.ItemID,
		OptionID: itemOption.OptionID,
		IsActive: itemOption.IsActive,
	}
}

func (r *menuItemOptionRepository) modelToEntity(dbItemOption *model.MenuItemOption) *entity.MenuItemOption {
	return &entity.MenuItemOption{
		ItemID:   dbItemOption.ItemID,
		OptionID: dbItemOption.OptionID,
		IsActive: dbItemOption.IsActive,
	}
}

func (r *menuItemOptionRepository) modelsToEntities(dbItemOptions []model.MenuItemOption) []*entity.MenuItemOption {
	entities := make([]*entity.MenuItemOption, len(dbItemOptions))
	for i, dbItemOption := range dbItemOptions {
		entities[i] = r.modelToEntity(&dbItemOption)
	}
	return entities
}
func (r *menuItemOptionRepository) DeleteByOptionID(ctx context.Context, optionID int) error {
	db := getDB(r.db, ctx)
	return db.WithContext(ctx).Where("option_id = ?", optionID).Delete(&model.MenuItemOption{}).Error
}
