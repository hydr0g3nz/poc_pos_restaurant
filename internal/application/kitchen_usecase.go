package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// kitchenUsecase implements KitchenUsecase interface
type kitchenUsecase struct {
	orderItemRepo       repository.OrderItemRepository
	orderRepo           repository.OrderRepository
	menuItemRepo        repository.MenuItemRepository
	tableRepo           repository.TableRepository
	orderItemOptionRepo repository.OrderItemOptionRepository
	menuOptionRepo      repository.MenuOptionRepository
	optionValueRepo     repository.OptionValueRepository
	logger              infra.Logger
	config              *config.Config
}

// NewKitchenUsecase creates a new kitchen usecase
func NewKitchenUsecase(
	orderItemRepo repository.OrderItemRepository,
	orderRepo repository.OrderRepository,
	menuItemRepo repository.MenuItemRepository,
	tableRepo repository.TableRepository,
	orderItemOptionRepo repository.OrderItemOptionRepository,
	menuOptionRepo repository.MenuOptionRepository,
	optionValueRepo repository.OptionValueRepository,
	logger infra.Logger,
	config *config.Config,
) KitchenUsecase {
	return &kitchenUsecase{
		orderItemRepo:       orderItemRepo,
		orderRepo:           orderRepo,
		menuItemRepo:        menuItemRepo,
		tableRepo:           tableRepo,
		orderItemOptionRepo: orderItemOptionRepo,
		menuOptionRepo:      menuOptionRepo,
		optionValueRepo:     optionValueRepo,
		logger:              logger,
		config:              config,
	}
}

// GetKitchenQueue retrieves all orders in kitchen queue
func (u *kitchenUsecase) GetKitchenQueue(ctx context.Context) ([]*KitchenOrderResponse, error) {
	u.logger.Debug("Getting kitchen queue")

	// Get all open orders
	openOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOpen), 100, 0)
	if err != nil {
		u.logger.Error("Error getting open orders", "error", err)
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	// Also get orders that are ordered but not completed
	orderedOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOrdered), 100, 0)
	if err != nil {
		u.logger.Error("Error getting ordered orders", "error", err)
		return nil, fmt.Errorf("failed to get ordered orders: %w", err)
	}

	// Combine all orders
	allOrders := append(openOrders, orderedOrders...)

	return u.toKitchenOrderResponses(ctx, allOrders), nil
}

// UpdateOrderItemStatus updates the status of an order item
func (u *kitchenUsecase) UpdateOrderItemStatus(ctx context.Context, orderItemID int, status string) (*OrderItemResponse, error) {
	u.logger.Info("Updating order item status", "orderItemID", orderItemID, "status", status)

	// Validate status
	itemStatus, err := vo.NewItemStatus(status)
	if err != nil {
		u.logger.Error("Invalid item status", "error", err, "status", status)
		return nil, err
	}

	// Get current order item
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		u.logger.Error("Error getting order item", "error", err, "orderItemID", orderItemID)
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return nil, errs.ErrOrderItemNotFound
	}

	// Update status
	orderItem.ItemStatus = itemStatus

	// Update timestamps based on status
	switch itemStatus {
	case vo.ItemStatusServed:
		now := new(time.Time)
		*now = time.Now()
		orderItem.ServedAt = now
	}

	// Save updated order item
	updatedItem, err := u.orderItemRepo.Update(ctx, orderItem)
	if err != nil {
		u.logger.Error("Error updating order item", "error", err, "orderItemID", orderItemID)
		return nil, fmt.Errorf("failed to update order item: %w", err)
	}

	u.logger.Info("Order item status updated successfully", "orderItemID", orderItemID, "status", status)

	return u.toOrderItemResponse(updatedItem), nil
}

// GetOrderItemsByStatus retrieves order items by status
func (u *kitchenUsecase) GetOrderItemsByStatus(ctx context.Context, status string) ([]*OrderItemResponse, error) {
	u.logger.Debug("Getting order items by status", "status", status)

	// Validate status
	_, err := vo.NewItemStatus(status)
	if err != nil {
		u.logger.Error("Invalid item status", "error", err, "status", status)
		return nil, err
	}

	// Note: This method would need a repository method to filter by item status
	// For now, we'll get all order items and filter in memory (not ideal for production)

	// Get all open/ordered orders first
	openOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOpen), 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	orderedOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOrdered), 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered orders: %w", err)
	}

	allOrders := append(openOrders, orderedOrders...)
	var filteredItems []*OrderItemResponse

	// Get items for each order and filter by status
	for _, order := range allOrders {
		items, err := u.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			u.logger.Error("Error getting order items", "error", err, "orderID", order.ID)
			continue
		}

		for _, item := range items {
			if item.ItemStatus.String() == status {
				filteredItems = append(filteredItems, u.toOrderItemResponse(item))
			}
		}
	}

	return filteredItems, nil
}

// MarkOrderItemAsReady marks an order item as ready
func (u *kitchenUsecase) MarkOrderItemAsReady(ctx context.Context, orderItemID int) (*OrderItemResponse, error) {
	u.logger.Info("Marking order item as ready", "orderItemID", orderItemID)

	return u.UpdateOrderItemStatus(ctx, orderItemID, string(vo.ItemStatusReady))
}

// MarkOrderItemAsServed marks an order item as served
func (u *kitchenUsecase) MarkOrderItemAsServed(ctx context.Context, orderItemID int) (*OrderItemResponse, error) {
	u.logger.Info("Marking order item as served", "orderItemID", orderItemID)

	return u.UpdateOrderItemStatus(ctx, orderItemID, string(vo.ItemStatusServed))
}

// GetKitchenOrdersByStation retrieves order items by kitchen station
func (u *kitchenUsecase) GetKitchenOrdersByStation(ctx context.Context, station string) ([]*OrderItemResponse, error) {
	u.logger.Debug("Getting kitchen orders by station", "station", station)

	// Note: This would need repository support to filter by kitchen station
	// For now, we'll get all active order items and filter

	// Get all open/ordered orders
	openOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOpen), 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}

	orderedOrders, err := u.orderRepo.ListByStatus(ctx, string(vo.OrderStatusOrdered), 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered orders: %w", err)
	}

	allOrders := append(openOrders, orderedOrders...)
	var stationItems []*OrderItemResponse

	// Get items for each order and filter by station
	for _, order := range allOrders {
		items, err := u.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			u.logger.Error("Error getting order items", "error", err, "orderID", order.ID)
			continue
		}

		for _, item := range items {
			if item.KitchenStation == station {
				stationItems = append(stationItems, u.toOrderItemResponse(item))
			}
		}
	}

	return stationItems, nil
}

// Helper methods

// toOrderItemResponse converts entity to response
func (u *kitchenUsecase) toOrderItemResponse(item *entity.OrderItem) *OrderItemResponse {
	return &OrderItemResponse{
		ID:        item.ID,
		OrderID:   item.OrderID,
		ItemID:    item.ItemID,
		Quantity:  item.Quantity,
		UnitPrice: item.UnitPrice.AmountBaht(),
		Subtotal:  item.CalculateSubtotal().AmountBaht(),
		CreatedAt: item.CreatedAt,
		Name:      item.Name,
	}
}

// toKitchenOrderResponses converts orders to kitchen order responses
func (u *kitchenUsecase) toKitchenOrderResponses(ctx context.Context, orders []*entity.Order) []*KitchenOrderResponse {
	responses := make([]*KitchenOrderResponse, len(orders))

	for i, order := range orders {
		// Get table information
		var tableNumber *int
		if order.TableID > 0 {
			table, err := u.tableRepo.GetByID(ctx, order.TableID)
			if err == nil && table != nil {
				tableNumber = &table.TableNumber
			}
		}

		// Get order items
		items, err := u.orderItemRepo.ListByOrder(ctx, order.ID)
		if err != nil {
			u.logger.Error("Error getting order items for kitchen", "error", err, "orderID", order.ID)
			items = []*entity.OrderItem{} // Empty slice to avoid nil
		}

		kitchenItems := make([]*KitchenOrderItemResponse, len(items))
		for j, item := range items {
			// Get order item options
			itemOptions, err := u.orderItemOptionRepo.GetByOrderItemID(ctx, item.ID)
			if err != nil {
				u.logger.Error("Error getting order item options", "error", err, "orderItemID", item.ID)
				itemOptions = []*entity.OrderItemOption{}
			}

			// Convert options to responses
			optionResponses := make([]*OrderItemOptionResponse, len(itemOptions))
			for k, option := range itemOptions {
				// Get menu option and value details
				menuOption, _ := u.menuOptionRepo.GetByID(ctx, option.OptionID)
				optionValue, _ := u.optionValueRepo.GetByID(ctx, option.ValueID)

				optionResponses[k] = &OrderItemOptionResponse{
					OrderItemID:     option.OrderItemID,
					OptionID:        option.OptionID,
					ValueID:         option.ValueID,
					AdditionalPrice: option.AdditionalPrice.AmountBaht(),
				}

				if menuOption != nil {
					optionResponses[k].Option = &MenuOptionResponse{
						ID:         menuOption.ID,
						Name:       menuOption.Name,
						Type:       menuOption.Type,
						IsRequired: menuOption.IsRequired,
					}
				}

				if optionValue != nil {
					optionResponses[k].Value = &OptionValueResponse{
						ID:              optionValue.ID,
						OptionID:        optionValue.OptionID,
						Name:            optionValue.Name,
						IsDefault:       optionValue.IsDefault,
						AdditionalPrice: optionValue.AdditionalPrice.AmountBaht(),
						DisplayOrder:    optionValue.DisplayOrder,
					}
				}
			}

			kitchenItems[j] = &KitchenOrderItemResponse{
				ID:             item.ID,
				ItemID:         item.ItemID,
				Name:           item.Name,
				Quantity:       item.Quantity,
				Status:         item.ItemStatus.String(),
				KitchenStation: item.KitchenStation,
				KitchenNotes:   item.KitchenNotes,
				Notes:          item.SpecialReq,
				Options:        optionResponses,
				CreatedAt:      item.CreatedAt,
				ServedAt:       item.ServedAt,
			}
		}

		responses[i] = &KitchenOrderResponse{
			OrderID:     order.ID,
			OrderNumber: order.OrderNumber,
			TableNumber: tableNumber,
			OrderType:   "dine_in", // Default, could be extended
			Items:       kitchenItems,
			CreatedAt:   order.CreatedAt,
		}
	}

	return responses
}
