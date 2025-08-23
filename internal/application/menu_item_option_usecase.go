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

// menuItemOptionUsecase implements MenuItemOptionUsecase interface
type menuItemOptionUsecase struct {
	menuItemOptionRepo repository.MenuItemOptionRepository
	menuItemRepo       repository.MenuItemRepository
	menuOptionRepo     repository.MenuOptionRepository
	optionValueRepo    repository.OptionValueRepository
	logger             infra.Logger
	config             *config.Config
}

// NewMenuItemOptionUsecase creates a new menu item option usecase
func NewMenuItemOptionUsecase(
	menuItemOptionRepo repository.MenuItemOptionRepository,
	menuItemRepo repository.MenuItemRepository,
	menuOptionRepo repository.MenuOptionRepository,
	optionValueRepo repository.OptionValueRepository,
	logger infra.Logger,
	config *config.Config,
) MenuItemOptionUsecase {
	return &menuItemOptionUsecase{
		menuItemOptionRepo: menuItemOptionRepo,
		menuItemRepo:       menuItemRepo,
		menuOptionRepo:     menuOptionRepo,
		optionValueRepo:    optionValueRepo,
		logger:             logger,
		config:             config,
	}
}

// AddOptionToMenuItem adds an option to a menu item
func (u *menuItemOptionUsecase) AddOptionToMenuItem(ctx context.Context, req *AddMenuItemOptionRequest) (*MenuItemOptionResponse, error) {
	u.logger.Info("Adding option to menu item", "itemID", req.ItemID, "optionID", req.OptionID)

	// Validate menu item exists
	menuItem, err := u.menuItemRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		u.logger.Error("Error getting menu item", "error", err, "itemID", req.ItemID)
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return nil, errs.ErrMenuItemNotFound
	}

	// Validate menu option exists
	menuOption, err := u.menuOptionRepo.GetByID(ctx, req.OptionID)
	if err != nil {
		u.logger.Error("Error getting menu option", "error", err, "optionID", req.OptionID)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if menuOption == nil {
		return nil, errs.NewNotFoundError("menu option", req.OptionID)
	}

	// Check if the relationship already exists
	existing, err := u.menuItemOptionRepo.GetByItemAndOption(ctx, req.ItemID, req.OptionID)
	if err != nil {
		u.logger.Error("Error checking existing menu item option", "error", err, "itemID", req.ItemID, "optionID", req.OptionID)
		return nil, fmt.Errorf("failed to check existing menu item option: %w", err)
	}
	if existing != nil {
		return nil, errs.NewConflictError("menu item option", "option already added to menu item")
	}

	// Create menu item option entity
	menuItemOption, err := entity.NewMenuItemOption(req.ItemID, req.OptionID, req.IsActive)
	if err != nil {
		u.logger.Error("Error creating menu item option entity", "error", err, "itemID", req.ItemID, "optionID", req.OptionID)
		return nil, err
	}

	// Save to database
	createdItemOption, err := u.menuItemOptionRepo.Create(ctx, menuItemOption)
	if err != nil {
		u.logger.Error("Error creating menu item option", "error", err, "itemID", req.ItemID, "optionID", req.OptionID)
		return nil, fmt.Errorf("failed to create menu item option: %w", err)
	}

	u.logger.Info("Option added to menu item successfully", "itemID", req.ItemID, "optionID", req.OptionID)

	return u.toMenuItemOptionResponse(createdItemOption, menuOption, nil), nil
}

// RemoveOptionFromMenuItem removes an option from a menu item
func (u *menuItemOptionUsecase) RemoveOptionFromMenuItem(ctx context.Context, itemID, optionID int) error {
	u.logger.Info("Removing option from menu item", "itemID", itemID, "optionID", optionID)

	// Check if the relationship exists
	existing, err := u.menuItemOptionRepo.GetByItemAndOption(ctx, itemID, optionID)
	if err != nil {
		u.logger.Error("Error getting menu item option", "error", err, "itemID", itemID, "optionID", optionID)
		return fmt.Errorf("failed to get menu item option: %w", err)
	}
	if existing == nil {
		return errs.NewNotFoundError("menu item option", fmt.Sprintf("itemID: %d, optionID: %d", itemID, optionID))
	}

	// Delete the relationship
	if err := u.menuItemOptionRepo.Delete(ctx, itemID, optionID); err != nil {
		u.logger.Error("Error removing option from menu item", "error", err, "itemID", itemID, "optionID", optionID)
		return fmt.Errorf("failed to remove option from menu item: %w", err)
	}

	u.logger.Info("Option removed from menu item successfully", "itemID", itemID, "optionID", optionID)
	return nil
}

// GetMenuItemOptions retrieves all options for a menu item
func (u *menuItemOptionUsecase) GetMenuItemOptions(ctx context.Context, itemID int) ([]*MenuItemOptionResponse, error) {
	u.logger.Debug("Getting menu item options", "itemID", itemID)

	// Validate menu item exists
	menuItem, err := u.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		u.logger.Error("Error getting menu item", "error", err, "itemID", itemID)
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return nil, errs.ErrMenuItemNotFound
	}

	// Get menu item options
	menuItemOptions, err := u.menuItemOptionRepo.GetByItemID(ctx, itemID)
	if err != nil {
		u.logger.Error("Error getting menu item options", "error", err, "itemID", itemID)
		return nil, fmt.Errorf("failed to get menu item options: %w", err)
	}

	return u.toMenuItemOptionResponses(ctx, menuItemOptions), nil
}

// UpdateMenuItemOption updates a menu item option
func (u *menuItemOptionUsecase) UpdateMenuItemOption(ctx context.Context, req *UpdateMenuItemOptionRequest) (*MenuItemOptionResponse, error) {
	u.logger.Info("Updating menu item option", "isActive", req.IsActive)

	// Note: The request doesn't include itemID and optionID, so we need to add them
	// This might need to be adjusted based on how you want to identify which menu item option to update
	// For now, I'll assume they should be part of the request or passed separately

	return nil, fmt.Errorf("UpdateMenuItemOption needs itemID and optionID to identify the record to update")
}

// Helper methods

// toMenuItemOptionResponse converts entity to response
func (u *menuItemOptionUsecase) toMenuItemOptionResponse(itemOption *entity.MenuItemOption, option *entity.MenuOption, values []*entity.OptionValue) *MenuItemOptionResponse {
	response := &MenuItemOptionResponse{
		ItemID:   itemOption.ItemID,
		OptionID: itemOption.OptionID,
		IsActive: itemOption.IsActive,
	}

	if option != nil {
		response.Option = &MenuOptionResponse{
			ID:         option.ID,
			Name:       option.Name,
			Type:       option.Type.String(),
			IsRequired: option.IsRequired,
		}
	}

	if values != nil {
		response.Values = make([]*OptionValueResponse, len(values))
		for i, value := range values {
			response.Values[i] = &OptionValueResponse{
				ID:              value.ID,
				OptionID:        value.OptionID,
				Name:            value.Name,
				IsDefault:       value.IsDefault,
				AdditionalPrice: value.AdditionalPrice.AmountBaht(),
				DisplayOrder:    value.DisplayOrder,
			}
		}
	}

	return response
}

// toMenuItemOptionResponses converts slice of entities to responses with related data
func (u *menuItemOptionUsecase) toMenuItemOptionResponses(ctx context.Context, itemOptions []*entity.MenuItemOption) []*MenuItemOptionResponse {
	responses := make([]*MenuItemOptionResponse, len(itemOptions))

	for i, itemOption := range itemOptions {
		// Get menu option details
		menuOption, err := u.menuOptionRepo.GetByID(ctx, itemOption.OptionID)
		if err != nil {
			u.logger.Error("Error getting menu option for response", "error", err, "optionID", itemOption.OptionID)
			menuOption = nil // Handle gracefully
		}

		// Get option values
		var optionValues []*entity.OptionValue
		if menuOption != nil {
			values, err := u.optionValueRepo.GetByOptionID(ctx, menuOption.ID)
			if err != nil {
				u.logger.Error("Error getting option values for response", "error", err, "optionID", menuOption.ID)
			} else {
				optionValues = values
			}
		}

		responses[i] = u.toMenuItemOptionResponse(itemOption, menuOption, optionValues)
	}

	return responses
}
