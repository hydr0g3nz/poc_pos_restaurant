package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// MenuItemController handles HTTP requests related to menu item operations
type MenuItemController struct {
	menuItemUseCase usecase.MenuItemUsecase
}

// NewMenuItemController creates a new instance of MenuItemController
func NewMenuItemController(menuItemUseCase usecase.MenuItemUsecase) *MenuItemController {
	return &MenuItemController{
		menuItemUseCase: menuItemUseCase,
	}
}

// CreateMenuItem handles menu item creation
func (c *MenuItemController) CreateMenuItem(ctx *fiber.Ctx) error {
	var req dto.CreateMenuItemRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	// Validate required fields
	if req.Name == "" || req.CategoryID <= 0 || req.Price < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Name, CategoryID, and valid Price are required",
		})
	}

	response, err := c.menuItemUseCase.CreateMenuItem(ctx.Context(), &usecase.CreateMenuItemRequest{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return HandleError(ctx, err)
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
		return HandleError(ctx, err)
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
		return HandleError(ctx, err)
	}

	// Validate required fields
	if req.Name == "" || req.CategoryID <= 0 || req.Price < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Name, CategoryID, and valid Price are required",
		})
	}

	response, err := c.menuItemUseCase.UpdateMenuItem(ctx.Context(), menuItemID, &usecase.UpdateMenuItemRequest{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return HandleError(ctx, err)
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
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item deleted successfully", nil)
}

// ListMenuItems handles getting all menu items with pagination
func (c *MenuItemController) ListMenuItems(ctx *fiber.Ctx) error {
	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid offset parameter",
		})
	}

	response, err := c.menuItemUseCase.ListMenuItems(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err)
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

	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid offset parameter",
		})
	}

	response, err := c.menuItemUseCase.ListMenuItemsByCategory(ctx.Context(), categoryID, limit, offset)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items retrieved successfully", response)
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

	limitStr := ctx.Query("limit", "10")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid limit parameter",
		})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid offset parameter",
		})
	}

	response, err := c.menuItemUseCase.SearchMenuItems(ctx.Context(), query, limit, offset)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items search completed successfully", response)
}

// RegisterRoutes registers the routes for the menu item controller
func (c *MenuItemController) RegisterRoutes(router fiber.Router) {
	menuItemGroup := router.Group("/menu-items")

	// Public routes
	menuItemGroup.Get("/", c.ListMenuItems)
	menuItemGroup.Get("/search", c.SearchMenuItems) // GET /menu-items/search?q=ข้าวผัด
	menuItemGroup.Get("/:id", c.GetMenuItem)
	menuItemGroup.Get("/category/:categoryId", c.ListMenuItemsByCategory)

	// Admin routes (require admin role in real implementation)
	menuItemGroup.Post("/", c.CreateMenuItem)
	menuItemGroup.Put("/:id", c.UpdateMenuItem)
	menuItemGroup.Delete("/:id", c.DeleteMenuItem)
}
