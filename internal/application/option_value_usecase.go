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

// optionValueUsecase implements OptionValueUsecase interface
type optionValueUsecase struct {
	optionValueRepo repository.OptionValueRepository
	menuOptionRepo  repository.MenuOptionRepository
	logger          infra.Logger
	config          *config.Config
}

// NewOptionValueUsecase creates a new option value usecase
func NewOptionValueUsecase(
	optionValueRepo repository.OptionValueRepository,
	menuOptionRepo repository.MenuOptionRepository,
	logger infra.Logger,
	config *config.Config,
) OptionValueUsecase {
	return &optionValueUsecase{
		optionValueRepo: optionValueRepo,
		menuOptionRepo:  menuOptionRepo,
		logger:          logger,
		config:          config,
	}
}

// CreateOptionValue creates a new option value
func (u *optionValueUsecase) CreateOptionValue(ctx context.Context, req *CreateOptionValueRequest) (*OptionValueResponse, error) {
	u.logger.Info("Creating option value", "optionID", req.OptionID, "name", req.Name)

	// Validate that the option exists
	menuOption, err := u.menuOptionRepo.GetByID(ctx, req.OptionID)
	if err != nil {
		u.logger.Error("Error getting menu option", "error", err, "optionID", req.OptionID)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if menuOption == nil {
		u.logger.Warn("Menu option not found", "optionID", req.OptionID)
		return nil, errs.NewNotFoundError("menu option", req.OptionID)
	}

	// Create option value entity
	optionValue, err := entity.NewOptionValue(req.OptionID, req.Name, req.IsDefault)
	if err != nil {
		u.logger.Error("Error creating option value entity", "error", err, "optionID", req.OptionID, "name", req.Name)
		return nil, err
	}

	// Set additional price if provided
	if req.AdditionalPrice > 0 {
		additionalPrice, err := vo.NewMoneyFromBaht(req.AdditionalPrice)
		if err != nil {
			u.logger.Error("Error creating additional price", "error", err, "price", req.AdditionalPrice)
			return nil, err
		}
		optionValue.AdditionalPrice = additionalPrice
	}

	// Set display order
	optionValue.DisplayOrder = req.DisplayOrder

	// Save to database
	createdValue, err := u.optionValueRepo.Create(ctx, optionValue)
	if err != nil {
		u.logger.Error("Error creating option value", "error", err, "optionID", req.OptionID, "name", req.Name)
		return nil, fmt.Errorf("failed to create option value: %w", err)
	}

	u.logger.Info("Option value created successfully", "valueID", createdValue.ID, "optionID", req.OptionID, "name", createdValue.Name)

	return u.toOptionValueResponse(createdValue, menuOption), nil
}

// GetOptionValue retrieves option value by ID
func (u *optionValueUsecase) GetOptionValue(ctx context.Context, id int) (*OptionValueResponse, error) {
	u.logger.Debug("Getting option value", "valueID", id)

	optionValue, err := u.optionValueRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting option value", "error", err, "valueID", id)
		return nil, fmt.Errorf("failed to get option value: %w", err)
	}
	if optionValue == nil {
		u.logger.Warn("Option value not found", "valueID", id)
		return nil, errs.NewNotFoundError("option value", id)
	}

	// Get the related menu option
	menuOption, err := u.menuOptionRepo.GetByID(ctx, optionValue.OptionID)
	if err != nil {
		u.logger.Error("Error getting related menu option", "error", err, "optionID", optionValue.OptionID)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}

	return u.toOptionValueResponse(optionValue, menuOption), nil
}

// GetOptionValuesByOptionID retrieves all option values for a specific option
func (u *optionValueUsecase) GetOptionValuesByOptionID(ctx context.Context, optionID int) ([]*OptionValueResponse, error) {
	u.logger.Debug("Getting option values by option ID", "optionID", optionID)

	// Validate that the option exists
	menuOption, err := u.menuOptionRepo.GetByID(ctx, optionID)
	if err != nil {
		u.logger.Error("Error getting menu option", "error", err, "optionID", optionID)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if menuOption == nil {
		return nil, errs.NewNotFoundError("menu option", optionID)
	}

	optionValues, err := u.optionValueRepo.GetByOptionID(ctx, optionID)
	if err != nil {
		u.logger.Error("Error getting option values by option ID", "error", err, "optionID", optionID)
		return nil, fmt.Errorf("failed to get option values: %w", err)
	}

	return u.toOptionValueResponses(optionValues, menuOption), nil
}

// UpdateOptionValue updates option value information
func (u *optionValueUsecase) UpdateOptionValue(ctx context.Context, id int, req *UpdateOptionValueRequest) (*OptionValueResponse, error) {
	u.logger.Info("Updating option value", "valueID", id, "name", req.Name)

	// Get current option value
	currentValue, err := u.optionValueRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current option value", "error", err, "valueID", id)
		return nil, fmt.Errorf("failed to get option value: %w", err)
	}
	if currentValue == nil {
		return nil, errs.NewNotFoundError("option value", id)
	}

	// Update fields
	currentValue.Name = req.Name
	currentValue.IsDefault = req.IsDefault
	currentValue.DisplayOrder = req.DisplayOrder

	// Update additional price
	if req.AdditionalPrice >= 0 {
		additionalPrice, err := vo.NewMoneyFromBaht(req.AdditionalPrice)
		if err != nil {
			u.logger.Error("Error creating additional price", "error", err, "price", req.AdditionalPrice)
			return nil, err
		}
		currentValue.AdditionalPrice = additionalPrice
	}

	// Validate updated entity
	if !currentValue.IsValid() {
		return nil, errs.ErrInvalidMenuOptionValue
	}

	// Update option value
	updatedValue, err := u.optionValueRepo.Update(ctx, currentValue)
	if err != nil {
		u.logger.Error("Error updating option value", "error", err, "valueID", id)
		return nil, fmt.Errorf("failed to update option value: %w", err)
	}

	// Get the related menu option
	menuOption, err := u.menuOptionRepo.GetByID(ctx, updatedValue.OptionID)
	if err != nil {
		u.logger.Error("Error getting related menu option", "error", err, "optionID", updatedValue.OptionID)
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}

	u.logger.Info("Option value updated successfully", "valueID", id)

	return u.toOptionValueResponse(updatedValue, menuOption), nil
}

// DeleteOptionValue deletes an option value
func (u *optionValueUsecase) DeleteOptionValue(ctx context.Context, id int) error {
	u.logger.Info("Deleting option value", "valueID", id)

	// Check if option value exists
	optionValue, err := u.optionValueRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get option value: %w", err)
	}
	if optionValue == nil {
		return errs.NewNotFoundError("option value", id)
	}

	// TODO: Check if option value is used in any order items
	// This would prevent deletion of values that are actively used

	// Delete option value
	if err := u.optionValueRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting option value", "error", err, "valueID", id)
		return fmt.Errorf("failed to delete option value: %w", err)
	}

	u.logger.Info("Option value deleted successfully", "valueID", id)
	return nil
}

// Helper methods

// toOptionValueResponse converts entity to response
func (u *optionValueUsecase) toOptionValueResponse(value *entity.OptionValue, option *entity.MenuOption) *OptionValueResponse {
	response := &OptionValueResponse{
		ID:              value.ID,
		OptionID:        value.OptionID,
		Name:            value.Name,
		IsDefault:       value.IsDefault,
		AdditionalPrice: value.AdditionalPrice.AmountBaht(),
		DisplayOrder:    value.DisplayOrder,
	}

	if option != nil {
		response.Option = &MenuOptionResponse{
			ID:         option.ID,
			Name:       option.Name,
			Type:       option.Type.String(),
			IsRequired: option.IsRequired,
		}
	}

	return response
}

// toOptionValueResponses converts slice of entities to responses
func (u *optionValueUsecase) toOptionValueResponses(values []*entity.OptionValue, option *entity.MenuOption) []*OptionValueResponse {
	responses := make([]*OptionValueResponse, len(values))
	for i, value := range values {
		responses[i] = u.toOptionValueResponse(value, option)
	}
	return responses
}
