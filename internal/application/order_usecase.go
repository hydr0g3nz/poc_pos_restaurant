package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/service"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// orderUsecase implements OrderUsecase interface
type orderUsecase struct {
	orderRepo      repository.OrderRepository
	orderItemRepo  repository.OrderItemRepository
	tableRepo      repository.TableRepository
	menuItemRepo   repository.MenuItemRepository
	orderService   service.OrderService
	qrCodeService  service.QRCodeService
	printerService infra.PrinterService
	logger         infra.Logger
	config         *config.Config
}

// NewOrderUsecase creates a new order usecase
func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	tableRepo repository.TableRepository,
	menuItemRepo repository.MenuItemRepository,
	orderService service.OrderService,
	qrCodeService service.QRCodeService,
	printerService infra.PrinterService,
	logger infra.Logger,
	config *config.Config,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:      orderRepo,
		orderItemRepo:  orderItemRepo,
		tableRepo:      tableRepo,
		menuItemRepo:   menuItemRepo,
		orderService:   orderService,
		printerService: printerService,
		logger:         logger,
		config:         config,
		qrCodeService:  qrCodeService,
	}
}

// CreateOrder creates a new order
func (u *orderUsecase) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, string, error) {
	u.logger.Info("Creating order", "tableID", req.TableID)

	// Validate order creation
	if err := u.orderService.ValidateOrderCreation(ctx, req.TableID); err != nil {
		u.logger.Error("Order validation failed", "error", err, "tableID", req.TableID)
		return nil, "", err
	}

	// Create order entity
	order, err := entity.NewOrder(req.TableID)
	if err != nil {
		u.logger.Error("Error creating order entity", "error", err, "tableID", req.TableID)
		return nil, "", err
	}
	qrCode := u.qrCodeService.GenerateQRCodeForOrder(ctx, order.ID)
	order.QRCode = qrCode
	qrCodeImageBytes, err := u.qrCodeService.GenerateQRCodeImage(ctx, qrCode)
	if err != nil {
		u.logger.Error("Error generating QR code image", "error", err, "qrCode", qrCode)
		return nil, "", fmt.Errorf("failed to generate QR code image: %w", err)
	}
	qrcodeImageBase64 := base64.StdEncoding.EncodeToString(qrCodeImageBytes)

	// Save to database
	createdOrder, err := u.orderRepo.Create(ctx, order)
	if err != nil {
		u.logger.Error("Error creating order", "error", err, "tableID", req.TableID)
		return nil, "", fmt.Errorf("failed to create order: %w", err)
	}

	u.logger.Info("Order created successfully", "orderID", createdOrder.ID, "tableID", createdOrder.TableID)

	return u.toOrderResponse(createdOrder), qrcodeImageBase64, nil
}

// GetOrder retrieves order by ID
func (u *orderUsecase) GetOrder(ctx context.Context, id int) (*OrderResponse, error) {
	u.logger.Debug("Getting order", "orderID", id)

	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		u.logger.Warn("Order not found", "orderID", id)
		return nil, errs.ErrOrderNotFound
	}

	return u.toOrderResponse(order), nil
}

// GetOrderWithItems retrieves order with its items
func (u *orderUsecase) GetOrderWithItems(ctx context.Context, id int) (*OrderWithItemsResponse, error) {
	u.logger.Debug("Getting order with items", "orderID", id)

	order, err := u.orderRepo.GetByIDWithItems(ctx, id)
	if err != nil {
		u.logger.Error("Error getting order with items", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to get order with items: %w", err)
	}
	if order == nil {
		u.logger.Warn("Order not found", "orderID", id)
		return nil, errs.ErrOrderNotFound
	}

	return u.toOrderWithItemsResponse(order), nil
}

// UpdateOrder updates order information
func (u *orderUsecase) UpdateOrder(ctx context.Context, id int, req *UpdateOrderRequest) (*OrderResponse, error) {
	u.logger.Info("Updating order", "orderID", id, "status", req.Status)

	// Get current order
	currentOrder, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current order", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if currentOrder == nil {
		return nil, errs.ErrOrderNotFound
	}

	// Validate and update status
	newStatus, err := vo.NewOrderStatus(req.Status)
	if err != nil {
		u.logger.Error("Invalid order status", "error", err, "status", req.Status)
		return nil, err
	}

	// Update order status
	if newStatus == vo.OrderStatusCompleted && !currentOrder.IsClosed() {
		if err := u.orderService.ProcessOrderClosure(ctx, id); err != nil {
			u.logger.Error("Order closure validation failed", "error", err, "orderID", id)
			return nil, err
		}
		currentOrder.Close()
	}

	currentOrder.OrderStatus = newStatus

	// Update order
	updatedOrder, err := u.orderRepo.Update(ctx, currentOrder)
	if err != nil {
		u.logger.Error("Error updating order", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	u.logger.Info("Order updated successfully", "orderID", id)

	return u.toOrderResponse(updatedOrder), nil
}

// CloseOrder closes an order
func (u *orderUsecase) CloseOrder(ctx context.Context, id int) (*OrderResponse, error) {
	u.logger.Info("Closing order", "orderID", id)

	// Get current order
	currentOrder, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current order", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Process order closure
	if err := u.orderService.ProcessOrderClosure(ctx, id); err != nil && !errors.Is(err, errs.ErrEmptyOrder) {
		u.logger.Error("Order closure validation failed", "error", err, "orderID", id)
		return nil, err
	}
	if err == errs.ErrEmptyOrder {
		deleteErr := u.orderRepo.Delete(ctx, id)
		if deleteErr != nil {
			u.logger.Error("Error deleting empty order", "error", deleteErr, "orderID", id)
			return nil, fmt.Errorf("failed to delete empty order: %w", deleteErr)
		}
	}
	// Close order
	currentOrder.Close()

	// Update order
	updatedOrder, err := u.orderRepo.Update(ctx, currentOrder)
	if err != nil {
		u.logger.Error("Error closing order", "error", err, "orderID", id)
		return nil, fmt.Errorf("failed to close order: %w", err)
	}

	u.logger.Info("Order closed successfully", "orderID", id)

	return u.toOrderResponse(updatedOrder), nil
}

// ListOrders retrieves all orders with pagination
func (u *orderUsecase) ListOrders(ctx context.Context, limit, offset int) (*OrderListResponse, error) {
	u.logger.Debug("Listing orders", "limit", limit, "offset", offset)

	orders, err := u.orderRepo.List(ctx, limit, offset)
	if err != nil {
		u.logger.Error("Error listing orders", "error", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return &OrderListResponse{
		Orders: u.toOrderResponses(orders),
		Total:  len(orders),
		Limit:  limit,
		Offset: offset,
	}, nil
}

// ListOrdersByTable retrieves orders for a specific table
func (u *orderUsecase) ListOrdersByTable(ctx context.Context, tableID int, limit, offset int) (*OrderListResponse, error) {
	u.logger.Debug("Listing orders by table", "tableID", tableID, "limit", limit, "offset", offset)

	orders, err := u.orderRepo.ListByTable(ctx, tableID, limit, offset)
	if err != nil {
		u.logger.Error("Error listing orders by table", "error", err, "tableID", tableID)
		return nil, fmt.Errorf("failed to list orders by table: %w", err)
	}

	return &OrderListResponse{
		Orders: u.toOrderResponses(orders),
		Total:  len(orders),
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetOpenOrderByTable retrieves open order for a table
func (u *orderUsecase) GetOpenOrderByTable(ctx context.Context, tableID int) (*OrderResponse, error) {
	u.logger.Debug("Getting open order by table", "tableID", tableID)

	order, err := u.orderRepo.GetOpenOrderByTable(ctx, tableID)
	if err != nil {
		u.logger.Error("Error getting open order by table", "error", err, "tableID", tableID)
		return nil, fmt.Errorf("failed to get open order by table: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}

	return u.toOrderResponse(order), nil
}

// GetOrdersByStatus retrieves orders by status
func (u *orderUsecase) GetOrdersByStatus(ctx context.Context, status string, limit, offset int) (*OrderListResponse, error) {
	u.logger.Debug("Getting orders by status", "status", status, "limit", limit, "offset", offset)

	// Validate status
	if _, err := vo.NewOrderStatus(status); err != nil {
		u.logger.Error("Invalid order status", "error", err, "status", status)
		return nil, err
	}

	orders, err := u.orderRepo.ListByStatus(ctx, status, limit, offset)
	if err != nil {
		u.logger.Error("Error getting orders by status", "error", err, "status", status)
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	return &OrderListResponse{
		Orders: u.toOrderResponses(orders),
		Total:  len(orders),
		Limit:  limit,
		Offset: offset,
	}, nil
}

// GetOrdersByDateRange retrieves orders within date range
func (u *orderUsecase) GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) (*OrderListResponse, error) {
	u.logger.Debug("Getting orders by date range", "startDate", startDate, "endDate", endDate, "limit", limit, "offset", offset)

	// Validate date range
	if startDate.After(endDate) {
		u.logger.Error("Invalid date range", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidDateRange
	}

	orders, err := u.orderRepo.ListByDateRange(ctx, startDate, endDate, limit, offset)
	if err != nil {
		u.logger.Error("Error getting orders by date range", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, fmt.Errorf("failed to get orders by date range: %w", err)
	}

	return &OrderListResponse{
		Orders: u.toOrderResponses(orders),
		Total:  len(orders),
		Limit:  limit,
		Offset: offset,
	}, nil
}

// AddOrderItem adds an item to an order
func (u *orderUsecase) AddOrderItem(ctx context.Context, req *AddOrderItemRequest) (*OrderItemResponse, error) {
	u.logger.Info("Adding order item", "orderID", req.OrderID, "itemID", req.ItemID, "quantity", req.Quantity)

	// Validate order item
	if err := u.orderService.ValidateOrderItem(ctx, req.OrderID, req.ItemID, req.Quantity); err != nil {
		u.logger.Error("Order item validation failed", "error", err, "orderID", req.OrderID, "itemID", req.ItemID)
		return nil, err
	}

	// Get menu item to get price
	menuItem, err := u.menuItemRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		u.logger.Error("Error getting menu item", "error", err, "itemID", req.ItemID)
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	if menuItem == nil {
		return nil, errs.ErrMenuItemNotFound
	}

	// Check if order item already exists
	existingItem, err := u.orderItemRepo.GetByOrderAndItem(ctx, req.OrderID, req.ItemID)
	if err != nil {
		u.logger.Error("Error checking existing order item", "error", err, "orderID", req.OrderID, "itemID", req.ItemID)
		return nil, fmt.Errorf("failed to check existing order item: %w", err)
	}

	if existingItem != nil {
		// Update existing item quantity
		if err := existingItem.UpdateQuantity(existingItem.Quantity + req.Quantity); err != nil {
			u.logger.Error("Error updating order item quantity", "error", err, "orderItemID", existingItem.ID)
			return nil, err
		}

		updatedItem, err := u.orderItemRepo.Update(ctx, existingItem)
		if err != nil {
			u.logger.Error("Error updating order item", "error", err, "orderItemID", existingItem.ID)
			return nil, fmt.Errorf("failed to update order item: %w", err)
		}

		return u.toOrderItemResponse(updatedItem), nil
	}

	// Create new order item
	orderItem, err := entity.NewOrderItem(req.OrderID, req.ItemID, req.Quantity, menuItem.Price.AmountBaht(), menuItem.Name)
	if err != nil {
		u.logger.Error("Error creating order item entity", "error", err, "orderID", req.OrderID, "itemID", req.ItemID)
		return nil, err
	}

	// Save to database
	createdItem, err := u.orderItemRepo.Create(ctx, orderItem)
	if err != nil {
		u.logger.Error("Error creating order item", "error", err, "orderID", req.OrderID, "itemID", req.ItemID)
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}

	u.logger.Info("Order item added successfully", "orderItemID", createdItem.ID, "orderID", req.OrderID, "itemID", req.ItemID)

	return u.toOrderItemResponse(createdItem), nil
}

// UpdateOrderItem updates an order item
func (u *orderUsecase) UpdateOrderItem(ctx context.Context, id int, req *UpdateOrderItemRequest) (*OrderItemResponse, error) {
	u.logger.Info("Updating order item", "orderItemID", id, "quantity", req.Quantity)

	// Get current order item
	currentItem, err := u.orderItemRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current order item", "error", err, "orderItemID", id)
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	if currentItem == nil {
		return nil, errs.ErrOrderItemNotFound
	}

	// Check if order is still open
	order, err := u.orderRepo.GetByID(ctx, currentItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", currentItem.OrderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return nil, errs.ErrCannotModifyClosedOrder
	}

	// Update quantity
	if err := currentItem.UpdateQuantity(req.Quantity); err != nil {
		u.logger.Error("Error updating order item quantity", "error", err, "orderItemID", id, "quantity", req.Quantity)
		return nil, err
	}

	// Update order item
	updatedItem, err := u.orderItemRepo.Update(ctx, currentItem)
	if err != nil {
		u.logger.Error("Error updating order item", "error", err, "orderItemID", id)
		return nil, fmt.Errorf("failed to update order item: %w", err)
	}

	u.logger.Info("Order item updated successfully", "orderItemID", id)

	return u.toOrderItemResponse(updatedItem), nil
}

// RemoveOrderItem removes an order item
func (u *orderUsecase) RemoveOrderItem(ctx context.Context, id int) error {
	u.logger.Info("Removing order item", "orderItemID", id)

	// Get current order item
	currentItem, err := u.orderItemRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting current order item", "error", err, "orderItemID", id)
		return fmt.Errorf("failed to get order item: %w", err)
	}
	if currentItem == nil {
		return errs.ErrOrderItemNotFound
	}

	// Check if order is still open
	order, err := u.orderRepo.GetByID(ctx, currentItem.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", currentItem.OrderID)
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return errs.ErrCannotModifyClosedOrder
	}

	// Delete order item
	if err := u.orderItemRepo.Delete(ctx, id); err != nil {
		u.logger.Error("Error deleting order item", "error", err, "orderItemID", id)
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	u.logger.Info("Order item removed successfully", "orderItemID", id)
	return nil
}

// ListOrderItems retrieves all items for an order
func (u *orderUsecase) ListOrderItems(ctx context.Context, orderID int) ([]*OrderItemResponse, error) {
	u.logger.Debug("Listing order items", "orderID", orderID)

	// Check if order exists
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}

	items, err := u.orderItemRepo.ListByOrder(ctx, orderID)
	if err != nil {
		u.logger.Error("Error listing order items", "error", err, "orderID", orderID)
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}

	return u.toOrderItemResponses(items), nil
}

// CalculateOrderTotal calculates the total amount for an order
func (u *orderUsecase) CalculateOrderTotal(ctx context.Context, orderID int) (*OrderTotalResponse, error) {
	u.logger.Debug("Calculating order total", "orderID", orderID)

	// Get order
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", orderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}

	// Get order items
	items, err := u.orderItemRepo.ListByOrder(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order items", "error", err, "orderID", orderID)
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	// Calculate total
	total := 0.0
	itemCount := 0
	for _, item := range items {
		subtotal := item.CalculateSubtotal()
		total += subtotal.AmountBaht()
		itemCount += item.Quantity
	}

	return &OrderTotalResponse{
		OrderID:   orderID,
		Items:     u.toOrderItemResponses(items),
		Total:     total,
		ItemCount: itemCount,
	}, nil
}

// Helper methods for conversion

// toOrderResponse converts entity to response
func (u *orderUsecase) toOrderResponse(order *entity.Order) *OrderResponse {
	response := &OrderResponse{
		ID:        order.ID,
		TableID:   order.TableID,
		Status:    order.OrderStatus.String(),
		CreatedAt: order.CreatedAt,
	}

	if order.ClosedAt != nil {
		response.ClosedAt = order.ClosedAt
	}

	return response
}

// toOrderWithItemsResponse converts entity to response with items
func (u *orderUsecase) toOrderWithItemsResponse(order *entity.Order) *OrderWithItemsResponse {
	response := &OrderWithItemsResponse{
		ID:        order.ID,
		TableID:   order.TableID,
		Status:    order.OrderStatus.String(),
		Items:     u.toOrderItemResponses(order.Items),
		Total:     order.CalculateTotal().AmountBaht(),
		CreatedAt: order.CreatedAt,
	}

	if order.ClosedAt != nil {
		response.ClosedAt = order.ClosedAt
	}

	return response
}

// toOrderResponses converts slice of entities to responses
func (u *orderUsecase) toOrderResponses(orders []*entity.Order) []*OrderResponse {
	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = u.toOrderResponse(order)
	}
	return responses
}

// toOrderItemResponse converts entity to response
func (u *orderUsecase) toOrderItemResponse(item *entity.OrderItem) *OrderItemResponse {
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

// toOrderItemResponses converts slice of entities to responses
func (u *orderUsecase) toOrderItemResponses(items []*entity.OrderItem) []*OrderItemResponse {
	responses := make([]*OrderItemResponse, len(items))
	for i, item := range items {
		responses[i] = u.toOrderItemResponse(item)
	}
	return responses
}

func (u *orderUsecase) PrintOrderReceipt(ctx context.Context, orderID int) error {
	u.logger.Info("Printing order receipt", "orderID", orderID)
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order for printing", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to get order for printing: %w", err)
	}
	if order == nil {
		u.logger.Warn("Order not found for printing", "orderID", orderID)
		return errs.ErrOrderNotFound
	}
	// Get order items
	items, err := u.orderItemRepo.ListByOrder(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order items for printing", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to get order items for printing: %w", err)
	}
	if len(items) == 0 {
		u.logger.Warn("No items found for order", "orderID", orderID)
		return errs.ErrEmptyOrder
	}
	order.Items = items
	// Generate receipt PDF
	receiptPDF, err := u.orderService.ReceiptPdf(ctx, order)
	if err != nil {
		u.logger.Error("Error generating receipt PDF", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to generate receipt PDF: %w", err)
	}
	// Print receipt
	if err := u.printerService.Print(ctx, receiptPDF, "PDF"); err != nil {
		u.logger.Error("Error printing receipt", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to print receipt: %w", err)
	}
	return nil
}
func (u *orderUsecase) PrintOrderQRCode(ctx context.Context, orderID int) error {
	u.logger.Info("Printing order receipt", "orderID", orderID)
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting order for printing", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to get order for printing: %w", err)
	}
	if order == nil {
		u.logger.Warn("Order not found for printing", "orderID", orderID)
		return errs.ErrOrderNotFound
	}

	// Generate receipt PDF
	qrCodePDF, err := u.orderService.QRCodePdf(ctx, order)
	if err != nil {
		u.logger.Error("Error generating receipt PDF", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to generate receipt PDF: %w", err)
	}
	// Print receipt
	if err := u.printerService.Print(ctx, qrCodePDF, "PDF"); err != nil {
		u.logger.Error("Error printing receipt", "error", err, "orderID", orderID)
		return fmt.Errorf("failed to print receipt: %w", err)
	}
	return nil
}
