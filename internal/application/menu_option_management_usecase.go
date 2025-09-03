// internal/application/usecase/menu_option_management_usecase.go
package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

type menuOptionManagementUsecase struct {
	repo               repository.Repository
	menuItemOptionRepo repository.MenuItemOptionRepository
	menuOptionRepo     repository.MenuOptionRepository
	optionValueRepo    repository.OptionValueRepository
}

func NewMenuOptionManagementUsecase(repo repository.Repository) MenuOptionManagementUsecase {
	return &menuOptionManagementUsecase{
		repo:               repo,
		menuOptionRepo:     repo.MenuOptionRepository(),
		optionValueRepo:    repo.OptionValueRepository(),
		menuItemOptionRepo: repo.MenuItemOptionRepository(),
	}
}

func (u *menuOptionManagementUsecase) CreateOptionWithValues(ctx context.Context, req *CreateOptionWithValuesRequest) (*OptionWithValuesResponse, error) {
	return u.doInTransaction(ctx, func(ctx context.Context) (*OptionWithValuesResponse, error) {
		// 1. Create menu option
		optionType := vo.OptionType(req.Type)
		option, err := entity.NewMenuOption(req.Name, optionType, req.IsRequired)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu option entity: %w", err)
		}

		createdOption, err := u.menuOptionRepo.Create(ctx, option)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu option: %w", err)
		}

		// 2. Create option values
		var valueResponses []*OptionValueResponse
		for i, valueReq := range req.Values {
			value, err := entity.NewOptionValue(createdOption.ID, valueReq.Name, valueReq.IsDefault)
			if err != nil {
				return nil, fmt.Errorf("failed to create option value entity: %w", err)
			}

			// Set additional price if provided
			if valueReq.AdditionalPrice > 0 {
				additionalPrice, err := vo.NewMoneyFromBaht(valueReq.AdditionalPrice)
				if err != nil {
					return nil, fmt.Errorf("failed to create additional price: %w", err)
				}
				value.AdditionalPrice = additionalPrice
			}

			// Set display order
			if valueReq.DisplayOrder > 0 {
				value.DisplayOrder = valueReq.DisplayOrder
			} else {
				value.DisplayOrder = i + 1 // Auto increment if not provided
			}

			createdValue, err := u.optionValueRepo.Create(ctx, value)
			if err != nil {
				return nil, fmt.Errorf("failed to create option value: %w", err)
			}

			valueResponses = append(valueResponses, &OptionValueResponse{
				ID:              createdValue.ID,
				OptionID:        createdValue.OptionID,
				Name:            createdValue.Name,
				IsDefault:       createdValue.IsDefault,
				AdditionalPrice: createdValue.AdditionalPrice.AmountBaht(),
				DisplayOrder:    createdValue.DisplayOrder,
			})
		}

		return &OptionWithValuesResponse{
			ID:         createdOption.ID,
			Name:       createdOption.Name,
			Type:       createdOption.Type.String(),
			IsRequired: createdOption.IsRequired,
			Values:     valueResponses,
		}, nil
	})
}

func (u *menuOptionManagementUsecase) UpdateOptionWithValues(ctx context.Context, optionID int, req *UpdateOptionWithValuesRequest) (*OptionWithValuesResponse, error) {
	return u.doInTransaction(ctx, func(ctx context.Context) (*OptionWithValuesResponse, error) {
		// 1. Get existing option
		existingOption, err := u.menuOptionRepo.GetByID(ctx, optionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get menu option: %w", err)
		}
		if existingOption == nil {
			return nil, errs.NewNotFoundError("menu option", optionID)
		}

		// 2. Update option
		existingOption.Name = req.Name
		existingOption.Type = vo.OptionType(req.Type)
		existingOption.IsRequired = req.IsRequired

		_, err = u.menuOptionRepo.Update(ctx, existingOption)
		if err != nil {
			return nil, fmt.Errorf("failed to update menu option: %w", err)
		}

		// 3. Handle values - process actions
		var valueResponses []*OptionValueResponse
		for _, valueReq := range req.Values {
			switch valueReq.Action {
			case "delete":
				if valueReq.ID != nil {
					err := u.optionValueRepo.Delete(ctx, *valueReq.ID)
					if err != nil {
						return nil, fmt.Errorf("failed to delete option value %d: %w", *valueReq.ID, err)
					}
				}
			case "update":
				if valueReq.ID != nil {
					existingValue, err := u.optionValueRepo.GetByID(ctx, *valueReq.ID)
					if err != nil {
						continue // Skip if not found
					}

					existingValue.Name = valueReq.Name
					existingValue.IsDefault = valueReq.IsDefault
					existingValue.DisplayOrder = valueReq.DisplayOrder

					if valueReq.AdditionalPrice > 0 {
						additionalPrice, err := vo.NewMoneyFromBaht(valueReq.AdditionalPrice)
						if err == nil {
							existingValue.AdditionalPrice = additionalPrice
						}
					}

					updatedValue, err := u.optionValueRepo.Update(ctx, existingValue)
					if err != nil {
						return nil, fmt.Errorf("failed to update option value: %w", err)
					}

					valueResponses = append(valueResponses, &OptionValueResponse{
						ID:              updatedValue.ID,
						OptionID:        updatedValue.OptionID,
						Name:            updatedValue.Name,
						IsDefault:       updatedValue.IsDefault,
						AdditionalPrice: updatedValue.AdditionalPrice.AmountBaht(),
						DisplayOrder:    updatedValue.DisplayOrder,
					})
				}
			case "add", "":
				// Create new value
				value, err := entity.NewOptionValue(optionID, valueReq.Name, valueReq.IsDefault)
				if err != nil {
					return nil, fmt.Errorf("failed to create option value entity: %w", err)
				}

				if valueReq.AdditionalPrice > 0 {
					additionalPrice, err := vo.NewMoneyFromBaht(valueReq.AdditionalPrice)
					if err == nil {
						value.AdditionalPrice = additionalPrice
					}
				}
				value.DisplayOrder = valueReq.DisplayOrder

				createdValue, err := u.optionValueRepo.Create(ctx, value)
				if err != nil {
					return nil, fmt.Errorf("failed to create option value: %w", err)
				}

				valueResponses = append(valueResponses, &OptionValueResponse{
					ID:              createdValue.ID,
					OptionID:        createdValue.OptionID,
					Name:            createdValue.Name,
					IsDefault:       createdValue.IsDefault,
					AdditionalPrice: createdValue.AdditionalPrice.AmountBaht(),
					DisplayOrder:    createdValue.DisplayOrder,
				})
			}
		}

		return &OptionWithValuesResponse{
			ID:         existingOption.ID,
			Name:       existingOption.Name,
			Type:       existingOption.Type.String(),
			IsRequired: existingOption.IsRequired,
			Values:     valueResponses,
		}, nil
	})
}

func (u *menuOptionManagementUsecase) GetOptionWithValues(ctx context.Context, optionID int) (*OptionWithValuesResponse, error) {
	// Get option
	option, err := u.menuOptionRepo.GetByID(ctx, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu option: %w", err)
	}
	if option == nil {
		return nil, errs.NewNotFoundError("menu option", optionID)
	}

	// Get values
	values, err := u.optionValueRepo.GetByOptionID(ctx, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get option values: %w", err)
	}

	var valueResponses []*OptionValueResponse
	for _, value := range values {
		valueResponses = append(valueResponses, &OptionValueResponse{
			ID:              value.ID,
			OptionID:        value.OptionID,
			Name:            value.Name,
			IsDefault:       value.IsDefault,
			AdditionalPrice: value.AdditionalPrice.AmountBaht(),
			DisplayOrder:    value.DisplayOrder,
		})
	}

	return &OptionWithValuesResponse{
		ID:         option.ID,
		Name:       option.Name,
		Type:       option.Type.String(),
		IsRequired: option.IsRequired,
		Values:     valueResponses,
	}, nil
}

func (u *menuOptionManagementUsecase) ListOptionsWithValues(ctx context.Context) ([]*OptionWithValuesResponse, error) {
	options, err := u.menuOptionRepo.List(ctx, 1000, 0) // Get all options
	if err != nil {
		return nil, fmt.Errorf("failed to list menu options: %w", err)
	}

	var responses []*OptionWithValuesResponse
	for _, option := range options {
		optionWithValues, err := u.GetOptionWithValues(ctx, option.ID)
		if err != nil {
			continue // Skip options that can't be loaded
		}
		responses = append(responses, optionWithValues)
	}

	return responses, nil
}

func (u *menuOptionManagementUsecase) DeleteOptionWithValues(ctx context.Context, optionID int) error {
	return u.doInTransactionVoid(ctx, func(ctx context.Context) error {
		// 1. Delete all values for this option
		values, err := u.optionValueRepo.GetByOptionID(ctx, optionID)
		if err != nil {
			return fmt.Errorf("failed to get option values: %w", err)
		}

		for _, value := range values {
			err := u.optionValueRepo.Delete(ctx, value.ID)
			if err != nil {
				return fmt.Errorf("failed to delete option value %d: %w", value.ID, err)
			}
		}
		// 2. Delete the option
		err = u.menuOptionRepo.Delete(ctx, optionID)
		if err != nil {
			return fmt.Errorf("failed to delete menu option: %w", err)
		}

		return nil
	})
}

// Transaction helper methods
func (u *menuOptionManagementUsecase) doInTransaction(ctx context.Context, fn func(ctx context.Context) (*OptionWithValuesResponse, error)) (*OptionWithValuesResponse, error) {
	txCtx, err := u.repo.TxManager().BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	result, err := fn(txCtx)
	if err != nil {
		if rollbackErr := u.repo.TxManager().RollbackTx(txCtx); rollbackErr != nil {
			return nil, fmt.Errorf("transaction failed: %w, rollback failed: %w", err, rollbackErr)
		}
		return nil, err
	}

	if err := u.repo.TxManager().CommitTx(txCtx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

func (u *menuOptionManagementUsecase) doInTransactionVoid(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, err := u.repo.TxManager().BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(txCtx); err != nil {
		if rollbackErr := u.repo.TxManager().RollbackTx(txCtx); rollbackErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %w", err, rollbackErr)
		}
		return err
	}

	if err := u.repo.TxManager().CommitTx(txCtx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
