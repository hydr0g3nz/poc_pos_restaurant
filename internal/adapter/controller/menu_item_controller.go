package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// MenuItemController handles HTTP requests related to menu item operations
type MenuItemController struct {
	menuItemUseCase usecase.MenuItemUsecase
	errorPresenter  presenter.ErrorPresenter
}

// NewMenuItemController creates a new instance of MenuItemController
func NewMenuItemController(menuItemUseCase usecase.MenuItemUsecase, errorPresenter presenter.ErrorPresenter) *MenuItemController {
	return &MenuItemController{
		menuItemUseCase: menuItemUseCase,
		errorPresenter:  errorPresenter,
	}
}

// RegisterRoutes registers the routes for the menu item controller
func (c *MenuItemController) RegisterRoutes(router fiber.Router) {
	menuItemGroup := router.Group("/menu-items")

	// Public routes
	menuItemGroup.Get("/", c.ListMenuItems)
	menuItemGroup.Get("/search", c.SearchMenuItems)                       // GET /menu-items/search?q=ข้าวผัด
	menuItemGroup.Get("/category/:categoryId", c.ListMenuItemsByCategory) // GET /menu-items/category/1
	menuItemGroup.Get("/:id", c.GetMenuItem)

	// Admin routes (require admin role in real implementation)
	menuItemGroup.Post("/", c.CreateMenuItem)
	menuItemGroup.Put("/:id", c.UpdateMenuItem)
	menuItemGroup.Delete("/:id", c.DeleteMenuItem)
}

// CreateMenuItem handles menu item creation
func (c *MenuItemController) CreateMenuItem(ctx *fiber.Ctx) error {
	var req dto.CreateMenuItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	// Validate required fields
	if req.CategoryID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category ID is required and must be greater than 0",
		})
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu item name is required",
		})
	}

	if req.Price < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Price must be non-negative",
		})
	}

	response, err := c.menuItemUseCase.CreateMenuItem(ctx.Context(), &usecase.CreateMenuItemRequest{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Menu item created successfully", response)
}

// GetMenuItem handles getting menu item by ID
func (c *MenuItemController) GetMenuItem(ctx *fiber.Ctx) error {
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

// UpdateMenuItem handles updating menu item
func (c *MenuItemController) UpdateMenuItem(ctx *fiber.Ctx) error {
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

	var req dto.UpdateMenuItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	// Validate required fields
	if req.CategoryID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category ID is required and must be greater than 0",
		})
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu item name is required",
		})
	}

	if req.Price < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Price must be non-negative",
		})
	}

	response, err := c.menuItemUseCase.UpdateMenuItem(ctx.Context(), menuItemID, &usecase.UpdateMenuItemRequest{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item updated successfully", response)
}

// DeleteMenuItem handles menu item deletion
func (c *MenuItemController) DeleteMenuItem(ctx *fiber.Ctx) error {
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

	err = c.menuItemUseCase.DeleteMenuItem(ctx.Context(), menuItemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item deleted successfully", nil)
}

// ListMenuItems handles getting all menu items
func (c *MenuItemController) ListMenuItems(ctx *fiber.Ctx) error {
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

	response, err := c.menuItemUseCase.ListMenuItems(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items retrieved successfully", response)
}

// ListMenuItemsByCategory handles getting menu items by category
func (c *MenuItemController) ListMenuItemsByCategory(ctx *fiber.Ctx) error {
	categoryIDParam := ctx.Params("categoryId")
	if categoryIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category ID is required",
		})
	}

	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid category ID format",
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

	response, err := c.menuItemUseCase.ListMenuItemsByCategory(ctx.Context(), categoryID, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items by category retrieved successfully", response)
}

// SearchMenuItems handles searching menu items
func (c *MenuItemController) SearchMenuItems(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Search query parameter 'q' is required",
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

	response, err := c.menuItemUseCase.SearchMenuItems(ctx.Context(), query, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items search completed successfully", response)
}
