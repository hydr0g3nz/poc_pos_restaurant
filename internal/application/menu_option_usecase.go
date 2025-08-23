package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// menuOptionUsecase implements MenuOptionUsecase interface
type menuOptionUsecase struct {
	menuOptionRepo repository.MenuOptionRepository
	logger         infra.Logger
	config         *config.Config
}

// NewMenuOptionUsecase creates a new menu option usecase
func NewMenuOptionUsecase(
	menuOptionRepo repository.MenuOptionRepository,
	logger infra.Logger,
	config *config.Config,
) MenuOptionUsecase {
	return &menuOptionUsecase{
		menuOptionRepo: menuOptionRepo,
		logger:         logger,
		config:         config,
	}
}

// CreateMenuOption creates a new menu option
func (u *menuOptionUsecase) CreateMenuOption(ctx context.Context, req *CreateMenuOptionRequest) (*MenuOptionResponse, error) {
	u.logger.Info("Creating menu option", "name", req.Name, "type", req.Type)

	// Create menu option entity
	menuOption, err := entity.NewMenuOption(req.Name, vo.OptionType(req.Type), req.IsRequired)
	if err != nil {
		u.logger.Error("Error creating menu option entity", "error", err, "name", req.Name, "type", req.Type)
		return nil, err
	}

	// Save to database
	createdOption, err := u.menuOptionRepo.Create(ctx, menuOption)
	if err != nil {
		u.logger.Error("Error creating menu option", "error", err, "name", req.Name, "type", req.Type)
		return nil, fmt.Errorf("failed to create menu option: %w", err)
	}

	u.logger.Info("Menu option created successfully", "optionID", createdOption.ID, "name", createdOption.Name, "type", createdOption.Type)

	return u.toMenuOptionResponse(createdOption), nil
}

// GetMenuOption retrieves menu option by ID
func (u *menuOptionUsecase) GetMenuOption(ctx context.Context, id int) (*MenuOptionResponse, error) {
	u.logger.Debug("Getting menu option", "optionID", id)

	menuOption, err := u.menuOptionRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting menu option", "error", err, "optionID", id)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if menuOption == nil {
		u.logger.Warn("Menu option not found", "optionID", id)
		return nil, errs.NewNotFoundError("menu option", id)
	}

	return u.toMenuOptionResponse(menuOption), nil
}

// UpdateMenuOption updates menu option information
func (u *menuOptionUsecase) UpdateMenuOption(ctx context.Context, id int, req *UpdateMenuOptionRequest) (*MenuOptionResponse, error) {
	u.logger.Info("Updating menu option", "optionID", id, "name", req.Name, "type", req.Type)

	// Get current menu option
	currentOption, err := u.menuOptionRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current menu option", "error", err, "optionID", id)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if currentOption == nil {
		return nil, errs.NewNotFoundError("menu option", id)
	}

	// Update fields
	currentOption.Name = req.Name
	currentOption.Type = vo.OptionType(req.Type)
	currentOption.IsRequired = req.IsRequired

	// Validate updated entity
	if !currentOption.IsValid() {
		return nil, errs.ErrInvalidMenuOption
	}

	// Update menu option
	updatedOption, err := u.menuOptionRepo.Update(ctx, currentOption)
	if err != nil {
		u.logger.Error("Error updating menu option", "error", err, "optionID", id)
		return nil, fmt.Errorf("failed to update menu option: %w", err)
	}

	u.logger.Info("Menu option updated successfully", "optionID", id)

	return u.toMenuOptionResponse(updatedOption), nil
}

// DeleteMenuOption deletes a menu option
func (u *menuOptionUsecase) DeleteMenuOption(ctx context.Context, id int) error {
	u.logger.Info("Deleting menu option", "optionID", id)

	// Check if menu option exists
	menuOption, err := u.menuOptionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get menu option: %w", err)
	}
	if menuOption == nil {
		return errs.NewNotFoundError("menu option", id)
	}

	// TODO: Check if menu option is used in any menu items
	// This would prevent deletion of options that are actively used

	// Delete menu option
	if err := u.menuOptionRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting menu option", "error", err, "optionID", id)
		return fmt.Errorf("failed to delete menu option: %w", err)
	}

	u.logger.Info("Menu option deleted successfully", "optionID", id)
	return nil
}

// ListMenuOptions retrieves all menu options
func (u *menuOptionUsecase) ListMenuOptions(ctx context.Context) ([]*MenuOptionResponse, error) {
	u.logger.Debug("Listing menu options")

	menuOptions, err := u.menuOptionRepo.List(ctx, 1000, 0) // Get all options
	if err != nil {
		u.logger.Error("Error listing menu options", "error", err)
		return nil, fmt.Errorf("failed to list menu options: %w", err)
	}

	return u.toMenuOptionResponses(menuOptions), nil
}

// GetMenuOptionsByType retrieves menu options by type
func (u *menuOptionUsecase) GetMenuOptionsByType(ctx context.Context, optionType string) ([]*MenuOptionResponse, error) {
	u.logger.Debug("Getting menu options by type", "type", optionType)

	menuOptions, err := u.menuOptionRepo.GetByType(ctx, optionType)
	if err != nil {
		u.logger.Error("Error getting menu options by type", "error", err, "type", optionType)
		return nil, fmt.Errorf("failed to get menu options by type: %w", err)
	}

	return u.toMenuOptionResponses(menuOptions), nil
}

// Helper methods

// toMenuOptionResponse converts entity to response
func (u *menuOptionUsecase) toMenuOptionResponse(option *entity.MenuOption) *MenuOptionResponse {
	return &MenuOptionResponse{
		ID:         option.ID,
		Name:       option.Name,
		Type:       option.Type.String(),
		IsRequired: option.IsRequired,
	}
}

// toMenuOptionResponses converts slice of entities to responses
func (u *menuOptionUsecase) toMenuOptionResponses(options []*entity.MenuOption) []*MenuOptionResponse {
	responses := make([]*MenuOptionResponse, len(options))
	for i, option := range options {
		responses[i] = u.toMenuOptionResponse(option)
	}
	return responses
}
