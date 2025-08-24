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

// orderItemOptionUsecase implements OrderItemOptionUsecase interface
type orderItemOptionUsecase struct {
	orderItemOptionRepo repository.OrderItemOptionRepository
	orderItemRepo       repository.OrderItemRepository
	menuOptionRepo      repository.MenuOptionRepository
	optionValueRepo     repository.OptionValueRepository
	orderRepo           repository.OrderRepository
	logger              infra.Logger
	config              *config.Config
}

// NewOrderItemOptionUsecase creates a new order item option usecase
func NewOrderItemOptionUsecase(
	orderItemOptionRepo repository.OrderItemOptionRepository,
	orderItemRepo repository.OrderItemRepository,
	menuOptionRepo repository.MenuOptionRepository,
	optionValueRepo repository.OptionValueRepository,
	orderRepo repository.OrderRepository,
	logger infra.Logger,
	config *config.Config,
) OrderItemOptionUsecase {
	return &orderItemOptionUsecase{
		orderItemOptionRepo: orderItemOptionRepo,
		orderItemRepo:       orderItemRepo,
		menuOptionRepo:      menuOptionRepo,
		optionValueRepo:     optionValueRepo,
		orderRepo:           orderRepo,
		logger:              logger,
		config:              config,
	}
}

// AddOptionToOrderItem adds an option to an order item
func (u *orderItemOptionUsecase) AddOptionToOrderItem(ctx context.Context, req *AddOrderItemOptionRequest) (*OrderItemOptionResponse, error) {
	u.logger.Info("Adding option to order item", "orderItemID", req.OrderItemID, "optionID", req.OptionID, "valueID", req.ValueID)

	// Validate order item exists and order is still open
	orderItem, err := u.orderItemRepo.GetByID(ctx, req.OrderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", req.OrderItemID)
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return nil, errs.ErrOrderItemNotFound
	}

	// Check if the order is still open
	order, err := u.orderRepo.GetByID(ctx, orderItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderItem.OrderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return nil, errs.ErrCannotModifyClosedOrder
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

	// Validate option value exists and belongs to the option
	optionValue, err := u.optionValueRepo.GetByID(ctx, req.ValueID)
	if err != nil {
		u.logger.Error("Error getting option value", "error", err, "valueID", req.ValueID)
		return nil, fmt.Errorf("failed to get option value: %w", err)
	}
	if optionValue == nil {
		return nil, errs.NewNotFoundError("option value", req.ValueID)
	}
	if optionValue.OptionID != req.OptionID {
		return nil, errs.NewValidationError("option_value", "does not belong to the specified option", req.ValueID)
	}

	// Create additional price money object
	var additionalPrice vo.Money
	if req.AdditionalPrice > 0 {
		additionalPrice, err = vo.NewMoneyFromBaht(req.AdditionalPrice)
		if err != nil {
			u.logger.Error("Error creating additional price", "error", err, "price", req.AdditionalPrice)
			return nil, err
		}
	} else {
		// Use the additional price from the option value if not provided
		additionalPrice = optionValue.AdditionalPrice
	}

	// Create order item option entity
	orderItemOption, err := entity.NewOrderItemOption(req.OrderItemID, req.OptionID, req.ValueID, additionalPrice)
	if err != nil {
		u.logger.Error("Error creating order item option entity", "error", err, "orderItemID", req.OrderItemID)
		return nil, err
	}

	// Save to database
	createdItemOption, err := u.orderItemOptionRepo.Create(ctx, orderItemOption)
	if err != nil {
		u.logger.Error("Error creating order item option", "error", err, "orderItemID", req.OrderItemID)
		return nil, fmt.Errorf("failed to create order item option: %w", err)
	}

	u.logger.Info("Option added to order item successfully", "orderItemID", req.OrderItemID, "optionID", req.OptionID, "valueID", req.ValueID)

	return u.toOrderItemOptionResponse(createdItemOption, menuOption, optionValue), nil
}

// UpdateOrderItemOption updates an order item option
func (u *orderItemOptionUsecase) UpdateOrderItemOption(ctx context.Context, req *UpdateOrderItemOptionRequest) (*OrderItemOptionResponse, error) {
	u.logger.Info("Updating order item option", "valueID", req.ValueID, "additionalPrice", req.AdditionalPrice)

	// Note: This method needs to identify which order item option to update
	// The current DTO doesn't include orderItemID and optionID
	// This implementation assumes they should be part of the request
	return nil, fmt.Errorf("UpdateOrderItemOption needs orderItemID and optionID to identify the record to update")
}

// RemoveOptionFromOrderItem removes an option from an order item
func (u *orderItemOptionUsecase) RemoveOptionFromOrderItem(ctx context.Context, orderItemID, optionID, valueID int) error {
	u.logger.Info("Removing option from order item", "orderItemID", orderItemID, "optionID", optionID, "valueID", valueID)

	// Validate order item exists and order is still open
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", orderItemID)
		return fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return errs.ErrOrderItemNotFound
	}

	// Check if the order is still open
	order, err := u.orderRepo.GetByID(ctx, orderItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderItem.OrderID)
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return errs.ErrCannotModifyClosedOrder
	}

	// Delete the order item option
	if err := u.orderItemOptionRepo.Delete(ctx, orderItemID, optionID, valueID); err != nil {
		u.logger.Error("Error removing option from order item", "error", err, "orderItemID", orderItemID, "optionID", optionID, "valueID", valueID)
		return fmt.Errorf("failed to remove option from order item: %w", err)
	}

	u.logger.Info("Option removed from order item successfully", "orderItemID", orderItemID, "optionID", optionID, "valueID", valueID)
	return nil
}

// GetOrderItemOptions retrieves all options for an order item
func (u *orderItemOptionUsecase) GetOrderItemOptions(ctx context.Context, orderItemID int) ([]*OrderItemOptionResponse, error) {
	u.logger.Debug("Getting order item options", "orderItemID", orderItemID)

	// Validate order item exists
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", orderItemID)
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return nil, errs.ErrOrderItemNotFound
	}

	// Get order item options
	orderItemOptions, err := u.orderItemOptionRepo.GetByOrderItemID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item options", "error", err, "orderItemID", orderItemID)
		return nil, fmt.Errorf("failed to get order item options: %w", err)
	}

	return u.toOrderItemOptionResponses(ctx, orderItemOptions), nil
}

// Helper methods

// toOrderItemOptionResponse converts entity to response
func (u *orderItemOptionUsecase) toOrderItemOptionResponse(itemOption *entity.OrderItemOption, option *entity.MenuOption, value *entity.OptionValue) *OrderItemOptionResponse {
	response := &OrderItemOptionResponse{
		OrderItemID:     itemOption.OrderItemID,
		OptionID:        itemOption.OptionID,
		ValueID:         itemOption.ValueID,
		AdditionalPrice: itemOption.AdditionalPrice.AmountBaht(),
	}

	if option != nil {
		response.Option = &MenuOptionResponse{
			ID:         option.ID,
			Name:       option.Name,
			Type:       option.Type.String(),
			IsRequired: option.IsRequired,
		}
	}

	if value != nil {
		response.Value = &OptionValueResponse{
			ID:              value.ID,
			OptionID:        value.OptionID,
			Name:            value.Name,
			IsDefault:       value.IsDefault,
			AdditionalPrice: value.AdditionalPrice.AmountBaht(),
			DisplayOrder:    value.DisplayOrder,
		}
	}

	return response
}

// toOrderItemOptionResponses converts slice of entities to responses with related data
func (u *orderItemOptionUsecase) toOrderItemOptionResponses(ctx context.Context, itemOptions []*entity.OrderItemOption) []*OrderItemOptionResponse {
	responses := make([]*OrderItemOptionResponse, len(itemOptions))

	for i, itemOption := range itemOptions {
		// Get menu option details
		menuOption, err := u.menuOptionRepo.GetByID(ctx, itemOption.OptionID)
		if err != nil {
			u.logger.Error("Error getting menu option for response", "error", err, "optionID", itemOption.OptionID)
			menuOption = nil // Handle gracefully
		}

		// Get option value details
		optionValue, err := u.optionValueRepo.GetByID(ctx, itemOption.ValueID)
		if err != nil {
			u.logger.Error("Error getting option value for response", "error", err, "valueID", itemOption.ValueID)
			optionValue = nil // Handle gracefully
		}

		responses[i] = u.toOrderItemOptionResponse(itemOption, menuOption, optionValue)
	}

	return responses
}

// เพิ่มใน internal/application/order_item_option_usecase.go

// RemoveAllOptionsFromOrderItem removes all options from an order item
func (u *orderItemOptionUsecase) RemoveAllOptionsFromOrderItem(ctx context.Context, orderItemID int) error {
	u.logger.Info("Removing all options from order item", "orderItemID", orderItemID)

	// ตรวจสอบว่า order item มีอยู่จริง
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", orderItemID)
		return fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return errs.ErrOrderItemNotFound
	}

	// ตรวจสอบว่า order ยังเปิดอยู่
	order, err := u.orderRepo.GetByID(ctx, orderItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderItem.OrderID)
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return errs.ErrCannotModifyClosedOrder
	}

	// ลบ options ทั้งหมด
	if err := u.orderItemOptionRepo.DeleteByOrderItemID(ctx, orderItemID); err != nil {
		u.logger.Error("Error removing all options from order item", "error", err, "orderItemID", orderItemID)
		return fmt.Errorf("failed to remove all options from order item: %w", err)
	}

	u.logger.Info("All options removed from order item successfully", "orderItemID", orderItemID)
	return nil
}

// เพิ่มใน internal/application/order_item_option_usecase.go

// RemoveSpecificOptionFromOrderItem removes a specific option from an order item
func (u *orderItemOptionUsecase) RemoveSpecificOptionFromOrderItem(ctx context.Context, orderItemID, optionID int) error {
	u.logger.Info("Removing specific option from order item", "orderItemID", orderItemID, "optionID", optionID)

	// ตรวจสอบว่า order item มีอยู่จริงและ order ยังเปิดอยู่
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", orderItemID)
		return fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return errs.ErrOrderItemNotFound
	}

	order, err := u.orderRepo.GetByID(ctx, orderItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderItem.OrderID)
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return errs.ErrCannotModifyClosedOrder
	}

	// หา options ที่มี optionID ตรงกัน
	existingOptions, err := u.orderItemOptionRepo.GetByOrderItemID(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to get existing options: %w", err)
	}

	// ลบทุก option ที่มี optionID ตรงกัน
	deleted := false
	for _, opt := range existingOptions {
		if opt.OptionID == optionID {
			if err := u.orderItemOptionRepo.Delete(ctx, orderItemID, optionID, opt.ValueID); err != nil {
				u.logger.Error("Error deleting option", "error", err, "orderItemID", orderItemID, "optionID", optionID, "valueID", opt.ValueID)
				return fmt.Errorf("failed to delete option: %w", err)
			}
			deleted = true
		}
	}

	if !deleted {
		return errs.NewNotFoundError("order item option", fmt.Sprintf("orderItemID: %d, optionID: %d", orderItemID, optionID))
	}

	u.logger.Info("Specific option removed from order item successfully", "orderItemID", orderItemID, "optionID", optionID)
	return nil
}
