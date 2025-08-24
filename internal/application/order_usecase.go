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
	orderRepo              repository.OrderRepository
	orderItemRepo          repository.OrderItemRepository
	tableRepo              repository.TableRepository
	menuItemRepo           repository.MenuItemRepository
	orderItemOptionUsecase OrderItemOptionUsecase
	orderService           service.OrderService
	qrCodeService          service.QRCodeService
	printerService         infra.PrinterService
	tx                     repository.TxManager
	logger                 infra.Logger
	config                 *config.Config
}

// NewOrderUsecase creates a new order usecase
func NewOrderUsecase(
	orderItemOptionUsecase OrderItemOptionUsecase,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	tableRepo repository.TableRepository,
	menuItemRepo repository.MenuItemRepository,
	orderService service.OrderService,
	qrCodeService service.QRCodeService,
	printerService infra.PrinterService,
	tx repository.TxManager,
	logger infra.Logger,
	config *config.Config,
) OrderUsecase {
	return &orderUsecase{
		orderItemOptionUsecase: orderItemOptionUsecase,
		orderRepo:              orderRepo,
		orderItemRepo:          orderItemRepo,
		tableRepo:              tableRepo,
		menuItemRepo:           menuItemRepo,
		orderService:           orderService,
		printerService:         printerService,
		tx:                     tx,
		logger:                 logger,
		config:                 config,
		qrCodeService:          qrCodeService,
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
func (u *orderUsecase) ListOrdersWithItems(ctx context.Context, limit, offset int) (*OrderWithItemsListResponse, error) {
	u.logger.Debug("Listing orders", "limit", limit, "offset", offset)

	orders, err := u.orderRepo.ListWithItems(ctx, limit, offset)
	if err != nil {
		u.logger.Error("Error listing orders", "error", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return &OrderWithItemsListResponse{
		Orders: u.toOrderWithItemsResponses(orders),
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
func (u *orderUsecase) toOrderWithItemsResponses(orders []*entity.Order) []*OrderWithItemsResponse {
	responses := make([]*OrderWithItemsResponse, len(orders))
	for i, order := range orders {
		responses[i] = u.toOrderWithItemsResponse(order)
	}
	return responses
}

// toOrderItemResponse converts entity to response
func (u *orderUsecase) toOrderItemResponse(item *entity.OrderItem) *OrderItemResponse {
	return &OrderItemResponse{
		ID:             item.ID,
		OrderID:        item.OrderID,
		ItemID:         item.ItemID,
		Quantity:       item.Quantity,
		UnitPrice:      item.UnitPrice.AmountBaht(),
		Subtotal:       item.CalculateSubtotal().AmountBaht(),
		CreatedAt:      item.CreatedAt,
		Name:           item.Name,
		KitchenStation: item.KitchenStation,
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
func (u *orderUsecase) AddOrderItemList(ctx context.Context, req *AddOrderItemListRequest) ([]*OrderItemResponse, error) {
	u.logger.Info("Adding order items list", "orderID", req.OrderID, "itemCount", len(req.Items))

	// เริ่ม transaction
	txCtx, err := u.tx.BeginTx(ctx)
	if err != nil {
		u.logger.Error("Error beginning transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			u.tx.RollbackTx(txCtx)
			panic(r)
		}
	}()

	var responses []*OrderItemResponse
	var addedOrderItems []*entity.OrderItem // เก็บ order items ที่เพิ่มแล้ว

	// วนลูปเพิ่มแต่ละ item
	for i, item := range req.Items {
		u.logger.Debug("Processing order item", "index", i, "menuItemID", item.MenuItemID, "quantity", item.Quantity)

		// Validate order item
		if err := u.orderService.ValidateOrderItem(txCtx, req.OrderID, item.MenuItemID, item.Quantity); err != nil {
			u.logger.Error("Order item validation failed", "error", err, "orderID", req.OrderID, "menuItemID", item.MenuItemID, "index", i)
			u.tx.RollbackTx(txCtx)
			return nil, fmt.Errorf("validation failed for item %d: %w", i, err)
		}

		// Get menu item เพื่อเอาราคา
		menuItem, err := u.menuItemRepo.GetByID(txCtx, item.MenuItemID)
		if err != nil {
			u.logger.Error("Error getting menu item", "error", err, "menuItemID", item.MenuItemID, "index", i)
			u.tx.RollbackTx(txCtx)
			return nil, fmt.Errorf("failed to get menu item %d: %w", item.MenuItemID, err)
		}
		if menuItem == nil {
			u.tx.RollbackTx(txCtx)
			return nil, errs.ErrMenuItemNotFound
		}

		// ตรวจสอบว่ามี order item นี้อยู่แล้วหรือไม่
		existingItem, err := u.orderItemRepo.GetByOrderAndItem(txCtx, req.OrderID, item.MenuItemID)
		if err != nil {
			u.logger.Error("Error checking existing order item", "error", err, "orderID", req.OrderID, "menuItemID", item.MenuItemID, "index", i)
			u.tx.RollbackTx(txCtx)
			return nil, fmt.Errorf("failed to check existing order item: %w", err)
		}

		var orderItem *entity.OrderItem

		if existingItem != nil {
			// อัปเดตจำนวนของ item ที่มีอยู่แล้ว
			if err := existingItem.UpdateQuantity(existingItem.Quantity + item.Quantity); err != nil {
				u.logger.Error("Error updating order item quantity", "error", err, "orderItemID", existingItem.ID, "index", i)
				u.tx.RollbackTx(txCtx)
				return nil, fmt.Errorf("failed to update order item quantity: %w", err)
			}

			updatedItem, err := u.orderItemRepo.Update(txCtx, existingItem)
			if err != nil {
				u.logger.Error("Error updating order item", "error", err, "orderItemID", existingItem.ID, "index", i)
				u.tx.RollbackTx(txCtx)
				return nil, fmt.Errorf("failed to update order item: %w", err)
			}
			orderItem = updatedItem
		} else {
			// สร้าง order item ใหม่
			newOrderItem, err := entity.NewOrderItem(req.OrderID, item.MenuItemID, item.Quantity, menuItem.Price.AmountBaht(), menuItem.Name)
			if err != nil {
				u.logger.Error("Error creating order item entity", "error", err, "orderID", req.OrderID, "menuItemID", item.MenuItemID, "index", i)
				u.tx.RollbackTx(txCtx)
				return nil, fmt.Errorf("failed to create order item entity: %w", err)
			}

			// Save to database
			createdItem, err := u.orderItemRepo.Create(txCtx, newOrderItem)
			if err != nil {
				u.logger.Error("Error creating order item", "error", err, "orderID", req.OrderID, "menuItemID", item.MenuItemID, "index", i)
				u.tx.RollbackTx(txCtx)
				return nil, fmt.Errorf("failed to create order item: %w", err)
			}
			orderItem = createdItem
		}

		addedOrderItems = append(addedOrderItems, orderItem)

		// เพิ่ม options สำหรับ order item นี้ (ถ้ามี)
		if len(item.Options) > 0 {
			u.logger.Debug("Adding options to order item", "orderItemID", orderItem.ID, "optionsCount", len(item.Options))

			for j, option := range item.Options {
				u.logger.Debug("Processing option", "optionIndex", j, "optionID", option.OptionID, "valueID", option.OptionValID)

				optionReq := &AddOrderItemOptionRequest{
					OrderItemID: orderItem.ID,
					OptionID:    option.OptionID,
					ValueID:     option.OptionValID,
				}

				_, err := u.orderItemOptionUsecase.AddOptionToOrderItem(txCtx, optionReq)
				if err != nil {
					u.logger.Error("Error adding option to order item", "error", err, "orderItemID", orderItem.ID, "optionID", option.OptionID, "valueID", option.OptionValID)
					u.tx.RollbackTx(txCtx)
					return nil, fmt.Errorf("failed to add option to order item %d: %w", orderItem.ID, err)
				}
			}
		}

		// เพิ่ม response
		responses = append(responses, u.toOrderItemResponse(orderItem))

		u.logger.Info("Order item added successfully", "orderItemID", orderItem.ID, "orderID", req.OrderID, "menuItemID", item.MenuItemID, "quantity", item.Quantity, "optionsCount", len(item.Options))
	}

	// Commit transaction
	if err := u.tx.CommitTx(txCtx); err != nil {
		u.logger.Error("Error committing transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.logger.Info("Order items list added successfully", "orderID", req.OrderID, "totalItems", len(responses))
	return responses, nil
}

// เพิ่มในไฟล์ internal/application/order_usecase.go

func (u *orderUsecase) UpdateOrderItemList(ctx context.Context, req *UpdateOrderItemListRequest) ([]*OrderItemResponse, error) {
	u.logger.Info("Updating order items list", "orderID", req.OrderID, "itemCount", len(req.Items))

	// ตรวจสอบว่า order มีอยู่และยังเปิดอยู่
	order, err := u.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", req.OrderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errs.ErrOrderNotFound
	}
	if order.IsClosed() {
		return nil, errs.ErrCannotModifyClosedOrder
	}

	// เริ่ม transaction
	txCtx, err := u.tx.BeginTx(ctx)
	if err != nil {
		u.logger.Error("Error beginning transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			u.tx.RollbackTx(txCtx)
			panic(r)
		}
	}()

	var responses []*OrderItemResponse

	// วนลูปอัปเดตแต่ละ item
	for i, item := range req.Items {
		u.logger.Debug("Processing order item update", "index", i, "orderItemID", item.OrderItemID, "action", item.Action)

		// หาก action เป็น delete ให้ลบ item
		if item.Action == "delete" {
			err := u.processDeleteOrderItem(txCtx, item.OrderItemID)
			if err != nil {
				u.logger.Error("Error deleting order item", "error", err, "orderItemID", item.OrderItemID)
				u.tx.RollbackTx(txCtx)
				return nil, fmt.Errorf("failed to delete order item %d: %w", item.OrderItemID, err)
			}
			continue // ไม่เพิ่มใน response เพราะถูกลบแล้ว
		}

		// Default action เป็น update
		updatedItem, err := u.processUpdateOrderItem(txCtx, req.OrderID, item)
		if err != nil {
			u.logger.Error("Error updating order item", "error", err, "orderItemID", item.OrderItemID)
			u.tx.RollbackTx(txCtx)
			return nil, fmt.Errorf("failed to update order item %d: %w", item.OrderItemID, err)
		}

		responses = append(responses, u.toOrderItemResponse(updatedItem))
	}

	// Commit transaction
	if err := u.tx.CommitTx(txCtx); err != nil {
		u.logger.Error("Error committing transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.logger.Info("Order items list updated successfully", "orderID", req.OrderID, "updatedItems", len(responses))
	return responses, nil
}

// Helper function สำหรับลบ order item
func (u *orderUsecase) processDeleteOrderItem(ctx context.Context, orderItemID int) error {
	// ตรวจสอบว่า order item มีอยู่จริง
	orderItem, err := u.orderItemRepo.GetByID(ctx, orderItemID)
	if err != nil {
		return fmt.Errorf("failed to get order item: %w", err)
	}
	if orderItem == nil {
		return errs.ErrOrderItemNotFound
	}

	// ลบ order item options ก่อน (ถ้ามี)
	err = u.orderItemOptionUsecase.RemoveAllOptionsFromOrderItem(ctx, orderItemID)
	if err != nil {
		u.logger.Warn("Error removing options from order item", "error", err, "orderItemID", orderItemID)
		// ไม่ return error เพราะอาจไม่มี options
	}

	// ลบ order item
	return u.orderItemRepo.Delete(ctx, orderItemID)
}

// อัปเดตใน internal/application/order_usecase.go

// Helper function สำหรับอัปเดต order item (แก้ไข)
func (u *orderUsecase) processUpdateOrderItem(ctx context.Context, orderID int, item *UpdateOrderItemRequest2) (*entity.OrderItem, error) {
	// ตรวจสอบว่า order item มีอยู่จริง
	currentItem, err := u.orderItemRepo.GetByID(ctx, item.OrderItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}
	if currentItem == nil {
		return nil, errs.ErrOrderItemNotFound
	}

	// ตรวจสอบว่า order item นี้เป็นของ order ที่ถูกต้อง
	if currentItem.OrderID != orderID {
		return nil, fmt.Errorf("order item %d does not belong to order %d", item.OrderItemID, orderID)
	}

	// หาก menu item เปลี่ยน ต้อง validate menu item ใหม่
	if currentItem.ItemID != item.MenuItemID {
		if err := u.orderService.ValidateOrderItem(ctx, orderID, item.MenuItemID, item.Quantity); err != nil {
			return nil, fmt.Errorf("validation failed for new menu item: %w", err)
		}

		// อัปเดต menu item และราคา
		menuItem, err := u.menuItemRepo.GetByID(ctx, item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to get menu item: %w", err)
		}
		if menuItem == nil {
			return nil, errs.ErrMenuItemNotFound
		}

		currentItem.ItemID = item.MenuItemID
		currentItem.Name = menuItem.Name
		price, err := vo.NewMoneyFromBaht(menuItem.Price.AmountBaht())
		if err != nil {
			return nil, fmt.Errorf("failed to create price: %w", err)
		}
		currentItem.UnitPrice = price
	}

	// อัปเดต quantity
	if err := currentItem.UpdateQuantity(item.Quantity); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	// บันทึกการเปลี่ยนแปลงของ order item
	updatedItem, err := u.orderItemRepo.Update(ctx, currentItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item: %w", err)
	}

	// จัดการ options แบบละเอียด
	if len(item.Options) > 0 {
		err = u.processOrderItemOptionsUpdate(ctx, item.OrderItemID, item.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to update options: %w", err)
		}
	}

	return updatedItem, nil
}

// อัปเดตใน internal/application/order_usecase.go

// Helper function สำหรับจัดการ options แบบละเอียด (แก้ไข)
func (u *orderUsecase) processOrderItemOptionsUpdate(ctx context.Context, orderItemID int, options []*OrderItemOptionUpdateRequest) error {
	for i, option := range options {
		u.logger.Debug("Processing option", "optionIndex", i, "optionID", option.OptionID, "valueID", option.OptionValID, "action", option.Action)

		switch option.Action {
		case "delete":
			if option.OptionValID != 0 {
				// ลบ option เฉพาะ value (แบบเดิม)
				err := u.orderItemOptionUsecase.RemoveOptionFromOrderItem(ctx, orderItemID, option.OptionID, option.OptionValID)
				if err != nil {
					u.logger.Error("Error deleting specific option value", "error", err, "orderItemID", orderItemID, "optionID", option.OptionID, "valueID", option.OptionValID)
					return fmt.Errorf("failed to delete option %d (value %d) from order item %d: %w", option.OptionID, option.OptionValID, orderItemID, err)
				}
			} else {
				// ลบทุก option ที่มี optionID เดียวกัน (ใช้ฟังก์ชันใหม่)
				err := u.orderItemOptionUsecase.RemoveSpecificOptionFromOrderItem(ctx, orderItemID, option.OptionID)
				if err != nil {
					u.logger.Error("Error deleting all options with optionID", "error", err, "orderItemID", orderItemID, "optionID", option.OptionID)
					return fmt.Errorf("failed to delete all options %d from order item %d: %w", option.OptionID, orderItemID, err)
				}
			}

		case "update":
			// อัปเดต option (ลบตัวเก่า แล้วเพิ่มตัวใหม่)
			if option.OptionValID == 0 {
				return fmt.Errorf("option_val_id is required for update action")
			}

			// หา option เก่าที่มี optionID เดียวกัน
			existingOptions, err := u.orderItemOptionUsecase.GetOrderItemOptions(ctx, orderItemID)
			if err != nil {
				return fmt.Errorf("failed to get existing options: %w", err)
			}

			// ลบ option เก่าที่มี optionID เดียวกันทั้งหมด
			var oldValueFound bool
			for _, existingOpt := range existingOptions {
				if existingOpt.OptionID == option.OptionID {
					err = u.orderItemOptionUsecase.RemoveOptionFromOrderItem(ctx, orderItemID, option.OptionID, existingOpt.ValueID)
					if err != nil {
						u.logger.Warn("Error removing old option for update", "error", err, "orderItemID", orderItemID, "optionID", option.OptionID, "valueID", existingOpt.ValueID)
					} else {
						oldValueFound = true
					}
				}
			}

			if !oldValueFound {
				u.logger.Warn("No existing option found to update", "orderItemID", orderItemID, "optionID", option.OptionID)
			}

			// เพิ่ม option ใหม่
			optionReq := &AddOrderItemOptionRequest{
				OrderItemID: orderItemID,
				OptionID:    option.OptionID,
				ValueID:     option.OptionValID,
			}

			_, err = u.orderItemOptionUsecase.AddOptionToOrderItem(ctx, optionReq)
			if err != nil {
				return fmt.Errorf("failed to add updated option: %w", err)
			}

		case "add", "": // default action เป็น add
			// เพิ่ม option ใหม่
			if option.OptionValID == 0 {
				return fmt.Errorf("option_val_id is required for add action")
			}

			optionReq := &AddOrderItemOptionRequest{
				OrderItemID: orderItemID,
				OptionID:    option.OptionID,
				ValueID:     option.OptionValID,
			}

			_, err := u.orderItemOptionUsecase.AddOptionToOrderItem(ctx, optionReq)
			if err != nil {
				return fmt.Errorf("failed to add option: %w", err)
			}

		default:
			return fmt.Errorf("invalid option action: %s", option.Action)
		}
	}

	return nil
}
