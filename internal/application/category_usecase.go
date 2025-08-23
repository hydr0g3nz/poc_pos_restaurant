package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
)

// categoryUsecase implements CategoryUsecase interface
type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
	logger       infra.Logger
	config       *config.Config
}

// NewCategoryUsecase creates a new category usecase
func NewCategoryUsecase(
	categoryRepo repository.CategoryRepository,
	logger infra.Logger,
	config *config.Config,
) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
		logger:       logger,
		config:       config,
	}
}

// CreateCategory creates a new category
func (u *categoryUsecase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	u.logger.Info("Creating category", "name", req.Name)

	// Check if category already exists
	existingCategory, err := u.categoryRepo.GetByName(ctx, req.Name)
	if err != nil {
		u.logger.Error("Error checking existing category", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to check existing category: %w", err)
	}
	if existingCategory != nil {
		u.logger.Warn("Category already exists", "name", req.Name)
		return nil, errs.ErrDuplicateCategoryName
	}

	// Create category entity
	category, err := entity.NewCategory(req.Name, req.Description, req.DisplayOrder, req.IsActive)
	if err != nil {
		u.logger.Error("Error creating category entity", "error", err, "name", req.Name)
		return nil, err
	}

	// Save to database
	createdCategory, err := u.categoryRepo.Create(ctx, category)
	if err != nil {
		u.logger.Error("Error creating category", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	u.logger.Info("Category created successfully", "categoryID", createdCategory.ID, "name", createdCategory.Name)

	return u.toCategoryResponse(createdCategory), nil
}

// GetCategory retrieves category by ID
func (u *categoryUsecase) GetCategory(ctx context.Context, id int) (*CategoryResponse, error) {
	u.logger.Debug("Getting category", "categoryID", id)

	category, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting category", "error", err, "categoryID", id)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		u.logger.Warn("Category not found", "categoryID", id)
		return nil, errs.ErrCategoryNotFound
	}

	return u.toCategoryResponse(category), nil
}

// GetCategoryByName retrieves category by name
func (u *categoryUsecase) GetCategoryByName(ctx context.Context, name string) (*CategoryResponse, error) {
	u.logger.Debug("Getting category by name", "name", name)

	category, err := u.categoryRepo.GetByName(ctx, name)
	if err != nil {
		u.logger.Error("Error getting category by name", "error", err, "name", name)
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}
	if category == nil {
		u.logger.Warn("Category not found", "name", name)
		return nil, errs.ErrCategoryNotFound
	}

	return u.toCategoryResponse(category), nil
}

// UpdateCategory updates category information
func (u *categoryUsecase) UpdateCategory(ctx context.Context, id int, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	u.logger.Info("Updating category", "categoryID", id, "name", req.Name)

	// Get current category
	currentCategory, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current category", "error", err, "categoryID", id)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if currentCategory == nil {
		return nil, errs.ErrCategoryNotFound
	}

	// Check if new name is different and unique
	if req.Name != currentCategory.Name {
		existingCategory, err := u.categoryRepo.GetByName(ctx, req.Name)
		if err != nil {
			u.logger.Error("Error checking name uniqueness", "error", err, "name", req.Name)
			return nil, fmt.Errorf("failed to check name uniqueness: %w", err)
		}
		if existingCategory != nil {
			return nil, errs.ErrDuplicateCategoryName
		}

		// Validate new category type
		// categoryType, err := vo.NewCategoryType(req.Name)
		// if err != nil {
		// 	u.logger.Error("Invalid category type", "error", err, "name", req.Name)
		// 	return nil, err
		// }
		// currentCategory.Name = categoryType
	}

	// Update category
	updatedCategory, err := u.categoryRepo.Update(ctx, currentCategory)
	if err != nil {
		u.logger.Error("Error updating category", "error", err, "categoryID", id)
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	u.logger.Info("Category updated successfully", "categoryID", id)

	return u.toCategoryResponse(updatedCategory), nil
}

// DeleteCategory deletes a category
func (u *categoryUsecase) DeleteCategory(ctx context.Context, id int) error {
	u.logger.Info("Deleting category", "categoryID", id)

	// Check if category exists
	category, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return errs.ErrCategoryNotFound
	}

	// Check if category has menu items
	hasMenuItems, err := u.categoryRepo.HasMenuItems(ctx, id)
	if err != nil {
		u.logger.Error("Error checking menu items", "error", err, "categoryID", id)
		return fmt.Errorf("failed to check menu items: %w", err)
	}
	if hasMenuItems {
		u.logger.Warn("Cannot delete category with menu items", "categoryID", id)
		return errs.ErrCannotDeleteCategoryWithItems
	}

	// Delete category
	if err := u.categoryRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting category", "error", err, "categoryID", id)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	u.logger.Info("Category deleted successfully", "categoryID", id)
	return nil
}

// ListCategories retrieves all categories
func (u *categoryUsecase) ListCategories(ctx context.Context) ([]*CategoryResponse, error) {
	u.logger.Debug("Listing categories")

	categories, err := u.categoryRepo.List(ctx)
	if err != nil {
		u.logger.Error("Error listing categories", "error", err)
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return u.toCategoryResponses(categories), nil
}

// Helper methods

// toCategoryResponse converts entity to response
func (u *categoryUsecase) toCategoryResponse(category *entity.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Description:  category.Description,
		DisplayOrder: category.DisplayOrder,
		IsActive:     category.IsActive,
		// CreatedAt: category.CreatedAt,
	}
}

// toCategoryResponses converts slice of entities to responses
func (u *categoryUsecase) toCategoryResponses(categories []*entity.Category) []*CategoryResponse {
	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = u.toCategoryResponse(category)
	}
	return responses
}
