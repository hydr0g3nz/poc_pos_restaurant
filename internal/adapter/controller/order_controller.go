package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// OrderController handles HTTP requests related to order operations
type OrderController struct {
	orderUseCase   usecase.OrderUsecase
	errorPresenter presenter.ErrorPresenter
}

// NewOrderController creates a new instance of OrderController
func NewOrderController(orderUseCase usecase.OrderUsecase, errorPresenter presenter.ErrorPresenter) *OrderController {
	return &OrderController{
		orderUseCase:   orderUseCase,
		errorPresenter: errorPresenter,
	}
}

// CreateOrder handles order creation
func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.TableID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table ID is required and must be greater than 0",
		})
	}

	response, err := c.orderUseCase.CreateOrder(ctx.Context(), &usecase.CreateOrderRequest{
		TableID: req.TableID,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Order created successfully", response)
}

// GetOrder handles getting order by ID
func (c *OrderController) GetOrder(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("id")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.orderUseCase.GetOrder(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order retrieved successfully", response)
}

// GetOrderWithItems handles getting order with items by ID
func (c *OrderController) GetOrderWithItems(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("id")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.orderUseCase.GetOrderWithItems(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order with items retrieved successfully", response)
}

// UpdateOrder handles updating order
func (c *OrderController) UpdateOrder(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("id")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	var req dto.UpdateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.Status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order status is required",
		})
	}

	response, err := c.orderUseCase.UpdateOrder(ctx.Context(), orderID, &usecase.UpdateOrderRequest{
		Status: req.Status,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order updated successfully", response)
}

// CloseOrder handles closing an order
func (c *OrderController) CloseOrder(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("id")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.orderUseCase.CloseOrder(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order closed successfully", response)
}

// ListOrders handles getting all orders
func (c *OrderController) ListOrders(ctx *fiber.Ctx) error {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.orderUseCase.ListOrders(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Orders retrieved successfully", response)
}

// ListOrdersByTable handles getting orders by table ID
func (c *OrderController) ListOrdersByTable(ctx *fiber.Ctx) error {
	tableIDParam := ctx.Params("tableId")
	if tableIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table ID is required",
		})
	}

	tableID, err := strconv.Atoi(tableIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Table ID format",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.orderUseCase.ListOrdersByTable(ctx.Context(), tableID, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Orders by table retrieved successfully", response)
}

// GetOpenOrderByTable handles getting open order by table ID
func (c *OrderController) GetOpenOrderByTable(ctx *fiber.Ctx) error {
	tableIDParam := ctx.Params("tableId")
	if tableIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table ID is required",
		})
	}

	tableID, err := strconv.Atoi(tableIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Table ID format",
		})
	}

	response, err := c.orderUseCase.GetOpenOrderByTable(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Open order by table retrieved successfully", response)
}

// GetOrdersByStatus handles getting orders by status
func (c *OrderController) GetOrdersByStatus(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	if status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Status query parameter is required",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.orderUseCase.GetOrdersByStatus(ctx.Context(), status, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Orders by status retrieved successfully", response)
}

// GetOrdersByDateRange handles getting orders by date range
func (c *OrderController) GetOrdersByDateRange(ctx *fiber.Ctx) error {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.orderUseCase.GetOrdersByDateRange(ctx.Context(), startDate, endDate, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Orders by date range retrieved successfully", response)
}

// AddOrderItem handles adding an item to an order
func (c *OrderController) AddOrderItem(ctx *fiber.Ctx) error {
	var req dto.AddOrderItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.OrderID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required and must be greater than 0",
		})
	}

	if req.ItemID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Item ID is required and must be greater than 0",
		})
	}

	if req.Quantity <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Quantity is required and must be greater than 0",
		})
	}

	response, err := c.orderUseCase.AddOrderItem(ctx.Context(), &usecase.AddOrderItemRequest{
		OrderID:  req.OrderID,
		ItemID:   req.ItemID,
		Quantity: req.Quantity,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Order item added successfully", response)
}

// UpdateOrderItem handles updating an order item
func (c *OrderController) UpdateOrderItem(ctx *fiber.Ctx) error {
	orderItemIDParam := ctx.Params("id")
	if orderItemIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order item ID is required",
		})
	}

	orderItemID, err := strconv.Atoi(orderItemIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order item ID format",
		})
	}

	var req dto.UpdateOrderItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.Quantity <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Quantity is required and must be greater than 0",
		})
	}

	response, err := c.orderUseCase.UpdateOrderItem(ctx.Context(), orderItemID, &usecase.UpdateOrderItemRequest{
		Quantity: req.Quantity,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order item updated successfully", response)
}

// RemoveOrderItem handles removing an order item
func (c *OrderController) RemoveOrderItem(ctx *fiber.Ctx) error {
	orderItemIDParam := ctx.Params("id")
	if orderItemIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order item ID is required",
		})
	}

	orderItemID, err := strconv.Atoi(orderItemIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order item ID format",
		})
	}

	err = c.orderUseCase.RemoveOrderItem(ctx.Context(), orderItemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order item removed successfully", nil)
}

// ListOrderItems handles getting all items for an order
func (c *OrderController) ListOrderItems(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("orderId")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.orderUseCase.ListOrderItems(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order items retrieved successfully", response)
}

// CalculateOrderTotal handles calculating order total
func (c *OrderController) CalculateOrderTotal(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("orderId")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.orderUseCase.CalculateOrderTotal(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order total calculated successfully", response)
}

// RegisterRoutes registers the routes for the order controller
func (c *OrderController) RegisterRoutes(router fiber.Router) {
	orderGroup := router.Group("/orders")

	// Order routes
	orderGroup.Post("/", c.CreateOrder)
	orderGroup.Get("/", c.ListOrders)
	orderGroup.Get("/search", c.GetOrdersByStatus)        // GET /orders/search?status=open
	orderGroup.Get("/date-range", c.GetOrdersByDateRange) // GET /orders/date-range?start_date=2024-01-01&end_date=2024-01-31
	orderGroup.Get("/:id", c.GetOrder)
	orderGroup.Get("/:id/items", c.GetOrderWithItems)
	orderGroup.Put("/:id", c.UpdateOrder)
	orderGroup.Put("/:id/close", c.CloseOrder)

	// Order by table routes
	orderGroup.Get("/table/:tableId", c.ListOrdersByTable)
	orderGroup.Get("/table/:tableId/open", c.GetOpenOrderByTable)

	// Order items routes
	orderGroup.Post("/items", c.AddOrderItem)
	orderGroup.Put("/items/:id", c.UpdateOrderItem)
	orderGroup.Delete("/items/:id", c.RemoveOrderItem)
	orderGroup.Get("/:orderId/items", c.ListOrderItems)
	orderGroup.Get("/:orderId/total", c.CalculateOrderTotal)
}
