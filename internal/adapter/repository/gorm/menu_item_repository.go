// internal/adapter/repository/menu_item_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type menuItemRepository struct {
	baseRepository
}

func NewMenuItemRepository(db *gorm.DB) repository.MenuItemRepository {
	return &menuItemRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *menuItemRepository) Create(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error) {
	dbItem := r.entityToModel(item)

	if err := r.db.WithContext(ctx).Create(dbItem).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItem)
}

func (r *menuItemRepository) GetByID(ctx context.Context, id int) (*entity.MenuItem, error) {
	var dbItem model.MenuItem

	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("KitchenStation").
		Preload("MenuItemOptions").
		Preload("MenuItemOptions.MenuOption").
		Preload("MenuItemOptions.MenuOption.OptionValues").
		First(&dbItem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbItem)
}

func (r *menuItemRepository) Update(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error) {
	dbItem := r.entityToModel(item)

	if err := r.db.WithContext(ctx).Save(dbItem).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbItem)
}

func (r *menuItemRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.MenuItem{}, id).Error
}

func (r *menuItemRepository) List(ctx context.Context, limit, offset int) ([]*entity.MenuItem, error) {
	var dbItems []model.MenuItem

	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("Category").Preload("KitchenStation").Find(&dbItems).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItems)
}

func (r *menuItemRepository) ListByCategory(ctx context.Context, categoryID int, limit, offset int) ([]*entity.MenuItem, error) {
	var dbItems []model.MenuItem

	query := r.db.WithContext(ctx).Where("category_id = ?", categoryID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Preload("Category").Preload("KitchenStation").Find(&dbItems).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItems)
}

func (r *menuItemRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.MenuItem, error) {
	var dbItems []model.MenuItem

	dbQuery := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	if err := dbQuery.Find(&dbItems).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbItems)
}

// Helper methods
func (r *menuItemRepository) entityToModel(item *entity.MenuItem) *model.MenuItem {
	return &model.MenuItem{
		ID:              item.ID,
		CategoryID:      item.CategoryID,
		Name:            item.Name,
		Description:     item.Description,
		Price:           item.Price.AmountSatang(),
		ImageURL:        item.ImageURL,
		IsRecommended:   item.IsRecommended,
		PreparationTime: item.PreparationTime,
		DisplayOrder:    item.DisplayOrder,
		KitchenID:       item.KitchenID,
		IsActive:        item.IsActive,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func (r *menuItemRepository) modelToEntity(dbItem *model.MenuItem) (*entity.MenuItem, error) {
	price, err := vo.NewMoneyFromSatang(dbItem.Price)
	if err != nil {
		return nil, err
	}
	m := &entity.MenuItem{
		ID:              dbItem.ID,
		CategoryID:      dbItem.CategoryID,
		Name:            dbItem.Name,
		Description:     dbItem.Description,
		Price:           price,
		ImageURL:        dbItem.ImageURL,
		IsRecommended:   dbItem.IsRecommended,
		PreparationTime: dbItem.PreparationTime,
		DisplayOrder:    dbItem.DisplayOrder,
		KitchenID:       dbItem.KitchenID,
		IsActive:        dbItem.IsActive,
		CreatedAt:       dbItem.CreatedAt,
		UpdatedAt:       dbItem.UpdatedAt,
	}
	if dbItem.Category != nil {
		m.Category = &entity.Category{
			ID:   dbItem.Category.ID,
			Name: dbItem.Category.Name,
		}
	}
	if dbItem.KitchenStation != nil {
		m.KitchenStation = &entity.KitchenStation{
			ID:   dbItem.KitchenStation.ID,
			Name: dbItem.KitchenStation.Name,
		}
	}
	if len(dbItem.MenuItemOptions) > 0 {
		options := model.ModelMenuItemOptionListToMenuItemOptionEntityList(dbItem.MenuItemOptions)
		m.MenuItemOptions = options
	}
	return m, nil
}

func (r *menuItemRepository) modelsToEntities(dbItems []model.MenuItem) ([]*entity.MenuItem, error) {
	entities := make([]*entity.MenuItem, len(dbItems))
	for i, dbItem := range dbItems {
		entity, err := r.modelToEntity(&dbItem)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
