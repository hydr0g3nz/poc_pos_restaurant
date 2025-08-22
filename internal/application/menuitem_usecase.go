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

// menuItemUsecase implements MenuItemUsecase interface
type menuItemUsecase struct {
	menuItemRepo repository.MenuItemRepository
	categoryRepo repository.CategoryRepository
	logger       infra.Logger
	config       *config.Config
}

// NewMenuItemUsecase creates a new menu item usecase
func NewMenuItemUsecase(
	menuItemRepo repository.MenuItemRepository,
	categoryRepo repository.CategoryRepository,
	logger infra.Logger,
	config *config.Config,
) MenuItemUsecase {
	return &menuItemUsecase{
		menuItemRepo: menuItemRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
		config:       config,
	}
}

// CreateMenuItem creates a new menu item
func (u *menuItemUsecase) CreateMenuItem(ctx context.Context, req *CreateMenuItemRequest) (*MenuItemResponse, error) {
	u.logger.Info("Creating menu item", "name", req.Name, "categoryID", req.CategoryID)

	// Validate category exists
	category, err := u.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		u.logger.Error("Error getting category", "error", err, "categoryID", req.CategoryID)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		u.logger.Warn("Category not found", "categoryID", req.CategoryID)
		return nil, errs.ErrCategoryNotFound
	}

	// Create menu item entity
	menuItem, err := entity.NewMenuItem(req.CategoryID, req.Name, req.Description, req.Price)
	if err != nil {
		u.logger.Error("Error creating menu item entity", "error", err, "name", req.Name)
		return nil, err
	}

	// Save to database
	createdMenuItem, err := u.menuItemRepo.Create(ctx, menuItem)
	if err != nil {
		u.logger.Error("Error creating menu item", "error", err, "name", req.Name)
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	u.logger.Info("Menu item created successfully", "menuItemID", createdMenuItem.ID, "name", createdMenuItem.Name)

	return u.toMenuItemResponse(createdMenuItem, category), nil
}

// GetMenuItem retrieves menu item by ID
func (u *menuItemUsecase) GetMenuItem(ctx context.Context, id int) (*MenuItemResponse, error) {
	u.logger.Debug("Getting menu item", "menuItemID", id)

	menuItem, err := u.menuItemRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting menu item", "error", err, "menuItemID", id)
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		u.logger.Warn("Menu item not found", "menuItemID", id)
		return nil, errs.ErrMenuItemNotFound
	}

	// Get category information
	category, err := u.categoryRepo.GetByID(ctx, menuItem.CategoryID)
	if err != nil {
		u.logger.Error("Error getting category for menu item", "error", err, "categoryID", menuItem.CategoryID)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return u.toMenuItemResponse(menuItem, category), nil
}

// UpdateMenuItem updates menu item information
func (u *menuItemUsecase) UpdateMenuItem(ctx context.Context, id int, req *UpdateMenuItemRequest) (*MenuItemResponse, error) {
	u.logger.Info("Updating menu item", "menuItemID", id, "name", req.Name)

	// Get current menu item
	currentMenuItem, err := u.menuItemRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current menu item", "error", err, "menuItemID", id)
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if currentMenuItem == nil {
		return nil, errs.ErrMenuItemNotFound
	}

	// Validate category exists
	category, err := u.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		u.logger.Error("Error getting category", "error", err, "categoryID", req.CategoryID)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	// Update menu item fields
	currentMenuItem.CategoryID = req.CategoryID
	currentMenuItem.Name = req.Name
	currentMenuItem.Description = req.Description

	// Update price if different
	if err := currentMenuItem.UpdatePrice(req.Price); err != nil {
		u.logger.Error("Error updating price", "error", err, "price", req.Price)
		return nil, err
	}

	// Update menu item
	updatedMenuItem, err := u.menuItemRepo.Update(ctx, currentMenuItem)
	if err != nil {
		u.logger.Error("Error updating menu item", "error", err, "menuItemID", id)
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	u.logger.Info("Menu item updated successfully", "menuItemID", id)

	return u.toMenuItemResponse(updatedMenuItem, category), nil
}

// DeleteMenuItem deletes a menu item
func (u *menuItemUsecase) DeleteMenuItem(ctx context.Context, id int) error {
	u.logger.Info("Deleting menu item", "menuItemID", id)

	// Check if menu item exists
	menuItem, err := u.menuItemRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return errs.ErrMenuItemNotFound
	}

	// TODO: Check if menu item is used in any orders (if order system is implemented)
	// This would prevent deletion of menu items that have been ordered

	// Delete menu item
	if err := u.menuItemRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting menu item", "error", err, "menuItemID", id)
		return fmt.Errorf("failed to delete menu item: %w", err)
	}

	u.logger.Info("Menu item deleted successfully", "menuItemID", id)
	return nil
}

// ListMenuItems retrieves menu items with pagination
func (u *menuItemUsecase) ListMenuItems(ctx context.Context, limit, offset int) (*MenuItemListResponse, error) {
	u.logger.Debug("Listing menu items", "limit", limit, "offset", offset)

	menuItems, err := u.menuItemRepo.List(ctx, limit, offset)
	if err != nil {
		u.logger.Error("Error listing menu items", "error", err)
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}

	responses, err := u.toMenuItemResponses(ctx, menuItems)
	if err != nil {
		return nil, err
	}

	return &MenuItemListResponse{
		Items:  responses,
		Total:  len(responses), // In a real implementation, you'd get the total count separately
		Limit:  limit,
		Offset: offset,
	}, nil
}

// ListMenuItemsByCategory retrieves menu items by category with pagination
func (u *menuItemUsecase) ListMenuItemsByCategory(ctx context.Context, categoryID int, limit, offset int) (*MenuItemListResponse, error) {
	u.logger.Debug("Listing menu items by category", "categoryID", categoryID, "limit", limit, "offset", offset)

	// Validate category exists
	category, err := u.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		u.logger.Error("Error getting category", "error", err, "categoryID", categoryID)
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, errs.ErrCategoryNotFound
	}

	menuItems, err := u.menuItemRepo.ListByCategory(ctx, categoryID, limit, offset)
	if err != nil {
		u.logger.Error("Error listing menu items by category", "error", err, "categoryID", categoryID)
		return nil, fmt.Errorf("failed to list menu items by category: %w", err)
	}

	responses, err := u.toMenuItemResponses(ctx, menuItems)
	if err != nil {
		return nil, err
	}

	return &MenuItemListResponse{
		Items:  responses,
		Total:  len(responses), // In a real implementation, you'd get the total count separately
		Limit:  limit,
		Offset: offset,
	}, nil
}

// SearchMenuItems searches menu items by name or description
func (u *menuItemUsecase) SearchMenuItems(ctx context.Context, query string, limit, offset int) (*MenuItemListResponse, error) {
	u.logger.Debug("Searching menu items", "query", query, "limit", limit, "offset", offset)

	if query == "" {
		u.logger.Warn("Empty search query provided")
		return &MenuItemListResponse{
			Items:  []*MenuItemResponse{},
			Total:  0,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	menuItems, err := u.menuItemRepo.Search(ctx, query, limit, offset)
	if err != nil {
		u.logger.Error("Error searching menu items", "error", err, "query", query)
		return nil, fmt.Errorf("failed to search menu items: %w", err)
	}

	responses, err := u.toMenuItemResponses(ctx, menuItems)
	if err != nil {
		return nil, err
	}

	return &MenuItemListResponse{
		Items:  responses,
		Total:  len(responses), // In a real implementation, you'd get the total count separately
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Helper methods

// toMenuItemResponse converts entity to response
func (u *menuItemUsecase) toMenuItemResponse(menuItem *entity.MenuItem, category *entity.Category) *MenuItemResponse {
	response := &MenuItemResponse{
		ID:          menuItem.ID,
		CategoryID:  menuItem.CategoryID,
		Name:        menuItem.Name,
		Description: menuItem.Description,
		Price:       menuItem.Price.AmountBaht(),
		CreatedAt:   menuItem.CreatedAt,
	}

	if category != nil {
		response.Category = &CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
		}
	}

	return response
}

// toMenuItemResponses converts slice of entities to responses with category information
func (u *menuItemUsecase) toMenuItemResponses(ctx context.Context, menuItems []*entity.MenuItem) ([]*MenuItemResponse, error) {
	responses := make([]*MenuItemResponse, len(menuItems))

	// Create a map to cache categories and avoid repeated database calls
	categoryCache := make(map[int]*entity.Category)

	for i, menuItem := range menuItems {
		var category *entity.Category
		var err error

		// Check if category is already cached
		if cachedCategory, exists := categoryCache[menuItem.CategoryID]; exists {
			category = cachedCategory
		} else {
			// Get category from database and cache it
			category, err = u.categoryRepo.GetByID(ctx, menuItem.CategoryID)
			if err != nil {
				u.logger.Error("Error getting category for menu item", "error", err, "categoryID", menuItem.CategoryID)
				return nil, fmt.Errorf("failed to get category: %w", err)
			}
			if category != nil {
				categoryCache[menuItem.CategoryID] = category
			}
		}

		responses[i] = u.toMenuItemResponse(menuItem, category)
	}

	return responses, nil
}
