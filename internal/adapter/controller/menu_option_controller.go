// internal/adapter/controller/menu_option_controller.go
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// MenuOptionController handles HTTP requests related to menu option operations
type MenuOptionController struct {
	menuOptionUseCase usecase.MenuOptionUsecase
	errorPresenter    presenter.ErrorPresenter
}

// NewMenuOptionController creates a new instance of MenuOptionController
func NewMenuOptionController(menuOptionUseCase usecase.MenuOptionUsecase, errorPresenter presenter.ErrorPresenter) *MenuOptionController {
	return &MenuOptionController{
		menuOptionUseCase: menuOptionUseCase,
		errorPresenter:    errorPresenter,
	}
}

// CreateMenuOption handles menu option creation
func (c *MenuOptionController) CreateMenuOption(ctx *fiber.Ctx) error {
	var req dto.CreateMenuOptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.menuOptionUseCase.CreateMenuOption(ctx.Context(), &usecase.CreateMenuOptionRequest{
		Name:       req.Name,
		Type:       req.Type,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Menu option created successfully", response)
}

// GetMenuOption handles getting menu option by ID
func (c *MenuOptionController) GetMenuOption(ctx *fiber.Ctx) error {
	optionIDParam := ctx.Params("id")
	if optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu option ID is required",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu option ID format",
		})
	}

	response, err := c.menuOptionUseCase.GetMenuOption(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu option retrieved successfully", response)
}

// UpdateMenuOption handles updating menu option
func (c *MenuOptionController) UpdateMenuOption(ctx *fiber.Ctx) error {
	optionIDParam := ctx.Params("id")
	if optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu option ID is required",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu option ID format",
		})
	}

	var req dto.UpdateMenuOptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.menuOptionUseCase.UpdateMenuOption(ctx.Context(), optionID, &usecase.UpdateMenuOptionRequest{
		Name:       req.Name,
		Type:       req.Type,
		IsRequired: req.IsRequired,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu option updated successfully", response)
}

// DeleteMenuOption handles menu option deletion
func (c *MenuOptionController) DeleteMenuOption(ctx *fiber.Ctx) error {
	optionIDParam := ctx.Params("id")
	if optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Menu option ID is required",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid menu option ID format",
		})
	}

	err = c.menuOptionUseCase.DeleteMenuOption(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu option deleted successfully", nil)
}

// ListMenuOptions handles getting all menu options
func (c *MenuOptionController) ListMenuOptions(ctx *fiber.Ctx) error {
	response, err := c.menuOptionUseCase.ListMenuOptions(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu options retrieved successfully", response)
}

// GetMenuOptionsByType handles getting menu options by type
func (c *MenuOptionController) GetMenuOptionsByType(ctx *fiber.Ctx) error {
	optionType := ctx.Query("type")
	if optionType == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Option type query parameter is required",
		})
	}

	response, err := c.menuOptionUseCase.GetMenuOptionsByType(ctx.Context(), optionType)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu options retrieved successfully", response)
}
