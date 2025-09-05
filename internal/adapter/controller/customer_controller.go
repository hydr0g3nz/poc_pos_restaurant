package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// CategoryController handles HTTP requests related to category operations
type CustomerController struct {
	categoryUseCase usecase.CategoryUsecase
	menuItemUseCase usecase.MenuItemUsecase
	orderUseCase    usecase.OrderUsecase
	errorPresenter  presenter.ErrorPresenter
}

// NewCategoryController creates a new instance of CategoryController
func NewCustomerController(categoryUseCase usecase.CategoryUsecase, menuItemUseCase usecase.MenuItemUsecase, orderUseCase usecase.OrderUsecase, errorPresenter presenter.ErrorPresenter) *CustomerController {
	return &CustomerController{
		orderUseCase:    orderUseCase,
		categoryUseCase: categoryUseCase,
		menuItemUseCase: menuItemUseCase,
		errorPresenter:  errorPresenter,
	}
}

// // CreateCategory handles category creation
// func (c *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
// 	var req dto.CreateCategoryRequest
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return HandleError(ctx, err, c.errorPresenter)
// 	}

// 	if req.Name == "" {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
// 			Status:  fiber.StatusBadRequest,
// 			Message: "Category name is required",
// 		})
// 	}

// 	response, err := c.categoryUseCase.CreateCategory(ctx.Context(), &usecase.CreateCategoryRequest{
// 		Name: req.Name,
// 	})
// 	if err != nil {
// 		return HandleError(ctx, err, c.errorPresenter)
// 	}

// 	return SuccessResp(ctx, fiber.StatusCreated, "Category created successfully", response)
// }

// GetCategory handles getting category by ID
func (c *CustomerController) ListCategory(ctx *fiber.Ctx) error {

	response, err := c.categoryUseCase.ListCategories(ctx.Context(), true)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}
	return SuccessResp(ctx, fiber.StatusOK, "Category retrieved successfully", response)
}
func (c *CustomerController) ListMenuItems(ctx *fiber.Ctx) error {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	response, err := c.menuItemUseCase.ListMenuItems(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items retrieved successfully", response)
}
func (c *CustomerController) GetMenuItem(ctx *fiber.Ctx) error {
	menuItemIDParam := ctx.Params("id")
	if menuItemIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu item ID is required",
		})
	}

	menuItemID, err := strconv.Atoi(menuItemIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu item ID format",
		})
	}

	response, err := c.menuItemUseCase.GetMenuItem(ctx.Context(), menuItemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item retrieved successfully", response)
}
func (c *CustomerController) SearchMenuItems(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Search query parameter 'q' is required",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	offset := (page - 1) * limit

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	response, err := c.menuItemUseCase.SearchMenuItems(ctx.Context(), query, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items search completed successfully", response)
}

// ใน controller
func (c *CustomerController) AddOrderItemList(ctx *fiber.Ctx) error {
	var req usecase.AddOrderItemListRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	responses, err := c.orderUseCase.AddOrderItemList(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Order items added successfully",
		"data":    responses,
	})
}
func (c *CustomerController) UpdateOrderItemList(ctx *fiber.Ctx) error {
	var req usecase.UpdateOrderItemListRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	responses, err := c.orderUseCase.UpdateOrderItemList(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Order items updated successfully",
		"data":    responses,
	})
}
func (c *CustomerController) ManageOrderItemList(ctx *fiber.Ctx) error {
	var req usecase.ManageOrderItemListRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	responses, err := c.orderUseCase.ManageOrderItemList(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Order items updated successfully",
		"data":    responses,
	})
}

func (c *CustomerController) GetOrderDetailWithOptions(ctx *fiber.Ctx) error {
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
			Message: "Invalid order ID format",
		})
	}

	response, err := c.orderUseCase.GetOrderDetailWithOptions(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order detail retrieved successfully", response)
}
func (c *CustomerController) GetOrderDetailByStatus(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	if status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Status is required",
		})
	}
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	offset := (page - 1) * limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	response, err := c.orderUseCase.GetOrdersByStatusWithDetails(ctx.Context(), status, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order detail retrieved successfully", response)
}
func (c *CustomerController) ListOrderItems(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	offset := (page - 1) * limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	response, err := c.orderUseCase.ListOrdersWithItemsAndCount(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order items retrieved successfully", response)
}
