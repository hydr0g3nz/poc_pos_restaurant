// internal/adapter/controller/menu_item_option_controller.go
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// MenuItemOptionController handles HTTP requests related to menu item option operations
type MenuItemOptionController struct {
	menuItemOptionUseCase usecase.MenuItemOptionUsecase
	errorPresenter        presenter.ErrorPresenter
}

// NewMenuItemOptionController creates a new instance of MenuItemOptionController
func NewMenuItemOptionController(menuItemOptionUseCase usecase.MenuItemOptionUsecase, errorPresenter presenter.ErrorPresenter) *MenuItemOptionController {
	return &MenuItemOptionController{
		menuItemOptionUseCase: menuItemOptionUseCase,
		errorPresenter:        errorPresenter,
	}
}

// AddOptionToMenuItem handles adding option to menu item
func (c *MenuItemOptionController) AddOptionToMenuItem(ctx *fiber.Ctx) error {
	var req dto.AddMenuItemOptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.menuItemOptionUseCase.AddOptionToMenuItem(ctx.Context(), &usecase.AddMenuItemOptionRequest{
		ItemID:   req.ItemID,
		OptionID: req.OptionID,
		IsActive: req.IsActive,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Option added to menu item successfully", response)
}

// RemoveOptionFromMenuItem handles removing option from menu item
func (c *MenuItemOptionController) RemoveOptionFromMenuItem(ctx *fiber.Ctx) error {
	itemIDParam := ctx.Params("itemId")
	optionIDParam := ctx.Params("optionId")

	if itemIDParam == "" || optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Item ID and Option ID are required",
		})
	}

	itemID, err := strconv.Atoi(itemIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid item ID format",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid option ID format",
		})
	}

	err = c.menuItemOptionUseCase.RemoveOptionFromMenuItem(ctx.Context(), itemID, optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option removed from menu item successfully", nil)
}

// GetMenuItemOptions handles getting all options for a menu item
func (c *MenuItemOptionController) GetMenuItemOptions(ctx *fiber.Ctx) error {
	itemIDParam := ctx.Params("itemId")
	if itemIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Item ID is required",
		})
	}

	itemID, err := strconv.Atoi(itemIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid item ID format",
		})
	}

	response, err := c.menuItemOptionUseCase.GetMenuItemOptions(ctx.Context(), itemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item options retrieved successfully", response)
}
