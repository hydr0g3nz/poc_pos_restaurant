// internal/adapter/controller/menu_with_options_controller.go
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

type MenuWithOptionsController struct {
	menuWithOptionsUc usecase.MenuWithOptionsUsecase
	menuOptionMgmtUc  usecase.MenuOptionManagementUsecase
	errorPresenter    presenter.ErrorPresenter
}

func NewMenuWithOptionsController(
	menuWithOptionsUc usecase.MenuWithOptionsUsecase,
	menuOptionMgmtUc usecase.MenuOptionManagementUsecase,
	errorPresenter presenter.ErrorPresenter,
) *MenuWithOptionsController {
	return &MenuWithOptionsController{
		menuWithOptionsUc: menuWithOptionsUc,
		menuOptionMgmtUc:  menuOptionMgmtUc,
		errorPresenter:    errorPresenter,
	}
}

// ==================== Menu Item with Options ====================

func (c *MenuWithOptionsController) CreateMenuItemWithOptions(ctx *fiber.Ctx) error {
	var req usecase.CreateMenuItemWithOptionsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuWithOptionsUc.CreateMenuItemWithOptions(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Menu item with options created successfully", result)
}

func (c *MenuWithOptionsController) UpdateMenuItemWithOptions(ctx *fiber.Ctx) error {
	itemID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	var req usecase.UpdateMenuItemWithOptionsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuWithOptionsUc.UpdateMenuItemWithOptions(ctx.Context(), itemID, &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item with options updated successfully", result)
}

func (c *MenuWithOptionsController) GetMenuItemWithOptions(ctx *fiber.Ctx) error {
	itemID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuWithOptionsUc.GetMenuItemWithOptions(ctx.Context(), itemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu item with options retrieved successfully", result)
}

func (c *MenuWithOptionsController) ListMenuItemsWithOptions(ctx *fiber.Ctx) error {
	req := usecase.ListMenuItemsRequest{
		Limit:  10, // Default values
		Offset: 0,
	}

	// Parse query parameters
	if categoryID := ctx.QueryInt("category_id", 0); categoryID > 0 {
		req.CategoryID = &categoryID
	}
	if isActive := ctx.QueryBool("is_active"); ctx.Query("is_active") != "" {
		req.IsActive = &isActive
	}
	if isRecommended := ctx.QueryBool("is_recommended"); ctx.Query("is_recommended") != "" {
		req.IsRecommended = &isRecommended
	}
	req.Search = ctx.Query("search", "")
	req.Limit = ctx.QueryInt("limit", 10)
	req.Offset = ctx.QueryInt("offset", 0)

	result, err := c.menuWithOptionsUc.ListMenuItemsWithOptions(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu items with options retrieved successfully", result)
}

func (c *MenuWithOptionsController) BulkAssignOptionsToMenuItems(ctx *fiber.Ctx) error {
	var req usecase.BulkAssignOptionsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	err := c.menuWithOptionsUc.BulkAssignOptionsToMenuItems(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Options assigned to menu items successfully", nil)
}

// ==================== Option Management ====================

func (c *MenuWithOptionsController) CreateOptionWithValues(ctx *fiber.Ctx) error {
	var req usecase.CreateOptionWithValuesRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuOptionMgmtUc.CreateOptionWithValues(ctx.Context(), &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Option with values created successfully", result)
}

func (c *MenuWithOptionsController) UpdateOptionWithValues(ctx *fiber.Ctx) error {
	optionID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	var req usecase.UpdateOptionWithValuesRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuOptionMgmtUc.UpdateOptionWithValues(ctx.Context(), optionID, &req)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option with values updated successfully", result)
}

func (c *MenuWithOptionsController) GetOptionWithValues(ctx *fiber.Ctx) error {
	optionID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	result, err := c.menuOptionMgmtUc.GetOptionWithValues(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option with values retrieved successfully", result)
}

func (c *MenuWithOptionsController) ListOptionsWithValues(ctx *fiber.Ctx) error {
	result, err := c.menuOptionMgmtUc.ListOptionsWithValues(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Options with values retrieved successfully", result)
}

func (c *MenuWithOptionsController) DeleteOptionWithValues(ctx *fiber.Ctx) error {
	optionID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	err = c.menuOptionMgmtUc.DeleteOptionWithValues(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option with values deleted successfully", nil)
}
