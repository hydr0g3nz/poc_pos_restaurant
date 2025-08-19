package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// CategoryController handles HTTP requests related to category operations
type CategoryController struct {
	categoryUseCase usecase.CategoryUsecase
	errorPresenter  presenter.ErrorPresenter
}

// NewCategoryController creates a new instance of CategoryController
func NewCategoryController(categoryUseCase usecase.CategoryUsecase, errorPresenter presenter.ErrorPresenter) *CategoryController {
	return &CategoryController{
		categoryUseCase: categoryUseCase,
		errorPresenter:  errorPresenter,
	}
}

// CreateCategory handles category creation
func (c *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category name is required",
		})
	}

	response, err := c.categoryUseCase.CreateCategory(ctx.Context(), &usecase.CreateCategoryRequest{
		Name: req.Name,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Category created successfully", response)
}

// GetCategory handles getting category by ID
func (c *CategoryController) GetCategory(ctx *fiber.Ctx) error {
	categoryIDParam := ctx.Params("id")
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
			Message: "Invalid Category ID format",
		})
	}

	response, err := c.categoryUseCase.GetCategory(ctx.Context(), categoryID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Category retrieved successfully", response)
}

// GetCategoryByName handles getting category by name
func (c *CategoryController) GetCategoryByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category name query parameter is required",
		})
	}

	response, err := c.categoryUseCase.GetCategoryByName(ctx.Context(), name)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Category retrieved successfully", response)
}

// UpdateCategory handles updating category
func (c *CategoryController) UpdateCategory(ctx *fiber.Ctx) error {
	categoryIDParam := ctx.Params("id")
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
			Message: "Invalid Category ID format",
		})
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Category name is required",
		})
	}

	response, err := c.categoryUseCase.UpdateCategory(ctx.Context(), categoryID, &usecase.UpdateCategoryRequest{
		Name: req.Name,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Category updated successfully", response)
}

// DeleteCategory handles category deletion
func (c *CategoryController) DeleteCategory(ctx *fiber.Ctx) error {
	categoryIDParam := ctx.Params("id")
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
			Message: "Invalid Category ID format",
		})
	}

	err = c.categoryUseCase.DeleteCategory(ctx.Context(), categoryID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Category deleted successfully", nil)
}

// ListCategories handles getting all categories
func (c *CategoryController) ListCategories(ctx *fiber.Ctx) error {
	response, err := c.categoryUseCase.ListCategories(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Categories retrieved successfully", response)
}
