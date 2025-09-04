// internal/application/usecase/menu_with_options_usecase.go
package usecase

import (
	"context"
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
)

type menuWithOptionsUsecase struct {
	repo               repository.Repository
	menuItemRepo       repository.MenuItemRepository
	menuOptionRepo     repository.MenuOptionRepository
	optionValueRepo    repository.OptionValueRepository
	menuItemOptionRepo repository.MenuItemOptionRepository
	categoryRepo       repository.CategoryRepository
	kitchenStationRepo repository.KitchenStationRepository
}

func NewMenuWithOptionsUsecase(repo repository.Repository) MenuWithOptionsUsecase {
	return &menuWithOptionsUsecase{
		repo:               repo,
		menuItemRepo:       repo.MenuItemRepository(),
		menuOptionRepo:     repo.MenuOptionRepository(),
		optionValueRepo:    repo.OptionValueRepository(),
		menuItemOptionRepo: repo.MenuItemOptionRepository(),
		categoryRepo:       repo.CategoryRepository(),
		kitchenStationRepo: repo.KitchenStationRepository(),
	}
}

func (u *menuWithOptionsUsecase) CreateMenuItemWithOptions(ctx context.Context, req *CreateMenuItemWithOptionsRequest) (*MenuItemWithOptionsResponse, error) {
	// ใช้ transaction
	return u.doInTransaction(ctx, func(ctx context.Context) (*MenuItemWithOptionsResponse, error) {
		// 1. สร้าง menu item
		menuItem, err := entity.NewMenuItem(req.CategoryID, req.Name, req.Description, req.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu item entity: %w", err)
		}

		menuItem.KitchenID = req.KitchenStationID
		menuItem.IsActive = req.IsActive
		menuItem.IsRecommended = req.IsRecommended
		menuItem.DisplayOrder = req.DisplayOrder

		createdItem, err := u.menuItemRepo.Create(ctx, menuItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create menu item: %w", err)
		}

		// 2. เพิ่ม options ให้ menu item
		for _, optionReq := range req.AssignedOptions {
			_, err := u.menuItemOptionRepo.Create(ctx, &entity.MenuItemOption{
				ItemID:   createdItem.ID,
				OptionID: optionReq.OptionID,
				IsActive: optionReq.IsActive,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to assign option %d to menu item: %w", optionReq.OptionID, err)
			}
		}

		// 3. Return complete response
		return u.GetMenuItemWithOptions(ctx, createdItem.ID)
	})
}

func (u *menuWithOptionsUsecase) UpdateMenuItemWithOptions(ctx context.Context, itemID int, req *UpdateMenuItemWithOptionsRequest) (*MenuItemWithOptionsResponse, error) {
	return u.doInTransaction(ctx, func(ctx context.Context) (*MenuItemWithOptionsResponse, error) {
		// 1. Get existing menu item
		existingItem, err := u.menuItemRepo.GetByID(ctx, itemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get menu item: %w", err)
		}
		if existingItem == nil {
			return nil, errs.ErrMenuItemNotFoundWithID(itemID)
		}

		// 2. Update menu item
		existingItem.CategoryID = req.CategoryID
		existingItem.KitchenID = req.KitchenStationID
		existingItem.Name = req.Name
		existingItem.Description = req.Description
		existingItem.UpdatePrice(req.Price)
		existingItem.IsActive = req.IsActive
		existingItem.IsRecommended = req.IsRecommended
		existingItem.DisplayOrder = req.DisplayOrder

		_, err = u.menuItemRepo.Update(ctx, existingItem)
		if err != nil {
			return nil, fmt.Errorf("failed to update menu item: %w", err)
		}

		// // 3. Update options - ลบเก่าแล้วเพิ่มใหม่
		// err = u.menuItemOptionRepo.DeleteByItemID(ctx, itemID)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to remove existing options: %w", err)
		// }
		// 4. เพิ่ม options ใหม่
		for _, optionReq := range req.AssignedOptions {
			if optionReq.OptionID == 0 {
				continue
			}
			if optionReq.IsActive {
				_, err := u.menuItemOptionRepo.Create(ctx, &entity.MenuItemOption{
					ItemID:   itemID,
					OptionID: optionReq.OptionID,
					IsActive: true,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to assign option %d to menu item: %w", optionReq.OptionID, err)
				}
			} else {
				update := &entity.MenuItemOption{
					ItemID:   itemID,
					OptionID: optionReq.OptionID,
					IsActive: false,
				}
				_, err := u.menuItemOptionRepo.Update(ctx, update)
				if err != nil {
					return nil, fmt.Errorf("failed to remove option %d from menu item: %w", optionReq.OptionID, err)
				}
			}
		}

		// 5. Return updated response
		return u.GetMenuItemWithOptions(ctx, itemID)
	})
}

func (u *menuWithOptionsUsecase) GetMenuItemWithOptions(ctx context.Context, itemID int) (*MenuItemWithOptionsResponse, error) {
	// 1. Get menu item
	menuItem, err := u.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return nil, errs.ErrMenuItemNotFoundWithID(itemID)
	}

	// 2. Get category name
	category, err := u.categoryRepo.GetByID(ctx, menuItem.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// 3. Get kitchen station name
	kitchenStation, err := u.kitchenStationRepo.GetByID(ctx, menuItem.KitchenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kitchen station: %w", err)
	}

	// 4. Get menu item options
	itemOptions, err := u.menuItemOptionRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item options: %w", err)
	}

	// 5. Build detailed options response
	var optionDetails []*MenuItemOptionDetailResponse
	for _, itemOption := range itemOptions {
		// Get option details
		option, err := u.menuOptionRepo.GetByID(ctx, itemOption.OptionID)
		if err != nil {
			continue // Skip if option not found
		}

		// Get option values
		values, err := u.optionValueRepo.GetByOptionID(ctx, itemOption.OptionID)
		if err != nil {
			continue // Skip if values not found
		}

		// Convert values
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

		optionDetails = append(optionDetails, &MenuItemOptionDetailResponse{
			OptionID:   option.ID,
			OptionName: option.Name,
			OptionType: option.Type.String(),
			IsRequired: option.IsRequired,
			IsActive:   itemOption.IsActive,
			Values:     valueResponses,
		})
	}

	return &MenuItemWithOptionsResponse{
		ID:               menuItem.ID,
		CategoryID:       menuItem.CategoryID,
		KitchenStationID: menuItem.KitchenID,
		Name:             menuItem.Name,
		Description:      menuItem.Description,
		Price:            menuItem.Price.AmountBaht(),
		IsActive:         menuItem.IsActive,
		IsRecommended:    menuItem.IsRecommended,
		DisplayOrder:     menuItem.DisplayOrder,
		Category:         category.Name,
		KitchenStation:   kitchenStation.Name,
		AvailableOptions: optionDetails,
		CreatedAt:        menuItem.CreatedAt,
		UpdatedAt:        menuItem.UpdatedAt,
	}, nil
}

func (u *menuWithOptionsUsecase) ListMenuItemsWithOptions(ctx context.Context, req *ListMenuItemsRequest) (*MenuItemWithOptionsListResponse, error) {
	// Get menu items based on filters
	menuItems, err := u.menuItemRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu items: %w", err)
	}

	// Build response with options for each item
	var itemResponses []*MenuItemWithOptionsResponse
	for _, item := range menuItems {
		itemWithOptions, err := u.GetMenuItemWithOptions(ctx, item.ID)
		if err != nil {
			continue // Skip items that can't be loaded
		}

		// Apply filters
		if req.CategoryID != nil && item.CategoryID != *req.CategoryID {
			continue
		}
		if req.IsActive != nil && item.IsActive != *req.IsActive {
			continue
		}
		if req.IsRecommended != nil && item.IsRecommended != *req.IsRecommended {
			continue
		}

		itemResponses = append(itemResponses, itemWithOptions)
	}

	return &MenuItemWithOptionsListResponse{
		Items:  itemResponses,
		Total:  len(itemResponses),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

func (u *menuWithOptionsUsecase) BulkAssignOptionsToMenuItems(ctx context.Context, req *BulkAssignOptionsRequest) error {
	return u.doInTransactionVoid(ctx, func(ctx context.Context) error {
		for _, menuItemID := range req.MenuItemIDs {
			// Check if menu item exists
			menuItem, err := u.menuItemRepo.GetByID(ctx, menuItemID)
			if err != nil {
				return fmt.Errorf("failed to get menu item %d: %w", menuItemID, err)
			}
			if menuItem == nil {
				return errs.ErrMenuItemNotFoundWithID(menuItemID)
			}

			// Assign each option to the menu item
			for _, optionID := range req.OptionIDs {
				// Check if option exists
				option, err := u.menuOptionRepo.GetByID(ctx, optionID)
				if err != nil {
					return fmt.Errorf("failed to get option %d: %w", optionID, err)
				}
				if option == nil {
					continue // Skip non-existent options
				}

				// Check if already assigned
				existing, err := u.menuItemOptionRepo.GetByItemAndOption(ctx, menuItemID, optionID)
				if err != nil {
					return fmt.Errorf("failed to check existing assignment: %w", err)
				}
				if existing != nil {
					continue // Skip if already assigned
				}

				// Create assignment
				_, err = u.menuItemOptionRepo.Create(ctx, &entity.MenuItemOption{
					ItemID:   menuItemID,
					OptionID: optionID,
					IsActive: req.IsActive,
				})
				if err != nil {
					return fmt.Errorf("failed to assign option %d to menu item %d: %w", optionID, menuItemID, err)
				}
			}
		}
		return nil
	})
}

func (u *menuWithOptionsUsecase) RemoveOptionFromMenuItem(ctx context.Context, itemID, optionID int) error {
	return u.menuItemOptionRepo.Delete(ctx, itemID, optionID)
}

// Helper methods for transaction handling
func (u *menuWithOptionsUsecase) doInTransaction(ctx context.Context, fn func(ctx context.Context) (*MenuItemWithOptionsResponse, error)) (*MenuItemWithOptionsResponse, error) {
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

func (u *menuWithOptionsUsecase) doInTransactionVoid(ctx context.Context, fn func(ctx context.Context) error) error {
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
