// internal/adapter/repository/category_repository.go
package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"gorm.io/gorm"
)

type categoryRepository struct {
	baseRepository
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	dbCategory := r.entityToModel(category)

	if err := r.db.WithContext(ctx).Create(dbCategory).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbCategory)
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*entity.Category, error) {
	var dbCategory model.Category

	if err := r.db.WithContext(ctx).First(&dbCategory, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbCategory)
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	var dbCategory model.Category

	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&dbCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbCategory)
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	dbCategory := r.entityToModel(category)

	if err := r.db.WithContext(ctx).Save(dbCategory).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbCategory)
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) List(ctx context.Context, onlyActive bool) ([]*entity.Category, error) {
	var dbCategories []model.Category
	query := r.db.WithContext(ctx)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	if err := query.WithContext(ctx).Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbCategories)
}

func (r *categoryRepository) HasMenuItems(ctx context.Context, categoryID int) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.MenuItem{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// Helper methods
func (r *categoryRepository) entityToModel(category *entity.Category) *model.Category {
	return &model.Category{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		DisplayOrder: category.DisplayOrder,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

func (r *categoryRepository) modelToEntity(dbCategory *model.Category) (*entity.Category, error) {

	return &entity.Category{
		ID:           dbCategory.ID,
		Name:         dbCategory.Name,
		Description:  dbCategory.Description,
		DisplayOrder: dbCategory.DisplayOrder,
		IsActive:     dbCategory.IsActive,
		CreatedAt:    dbCategory.CreatedAt,
		UpdatedAt:    dbCategory.UpdatedAt,
	}, nil
}

func (r *categoryRepository) modelsToEntities(dbCategories []model.Category) ([]*entity.Category, error) {
	entities := make([]*entity.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		entity, err := r.modelToEntity(&dbCategory)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
