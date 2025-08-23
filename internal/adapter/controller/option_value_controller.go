// internal/adapter/controller/option_value_controller.go
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// OptionValueController handles HTTP requests related to option value operations
type OptionValueController struct {
	optionValueUseCase usecase.OptionValueUsecase
	errorPresenter     presenter.ErrorPresenter
}

// NewOptionValueController creates a new instance of OptionValueController
func NewOptionValueController(optionValueUseCase usecase.OptionValueUsecase, errorPresenter presenter.ErrorPresenter) *OptionValueController {
	return &OptionValueController{
		optionValueUseCase: optionValueUseCase,
		errorPresenter:     errorPresenter,
	}
}

// CreateOptionValue handles option value creation
func (c *OptionValueController) CreateOptionValue(ctx *fiber.Ctx) error {
	var req dto.CreateOptionValueRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.optionValueUseCase.CreateOptionValue(ctx.Context(), &usecase.CreateOptionValueRequest{
		OptionID:        req.OptionID,
		Name:            req.Name,
		IsDefault:       req.IsDefault,
		AdditionalPrice: req.AdditionalPrice,
		DisplayOrder:    req.DisplayOrder,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Option value created successfully", response)
}

// GetOptionValue handles getting option value by ID
func (c *OptionValueController) GetOptionValue(ctx *fiber.Ctx) error {
	valueIDParam := ctx.Params("id")
	if valueIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Option value ID is required",
		})
	}

	valueID, err := strconv.Atoi(valueIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid option value ID format",
		})
	}

	response, err := c.optionValueUseCase.GetOptionValue(ctx.Context(), valueID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option value retrieved successfully", response)
}

// GetOptionValuesByOptionID handles getting option values by option ID
func (c *OptionValueController) GetOptionValuesByOptionID(ctx *fiber.Ctx) error {
	optionIDParam := ctx.Params("optionId")
	if optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Option ID is required",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid option ID format",
		})
	}

	response, err := c.optionValueUseCase.GetOptionValuesByOptionID(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option values retrieved successfully", response)
}

// UpdateOptionValue handles updating option value
func (c *OptionValueController) UpdateOptionValue(ctx *fiber.Ctx) error {
	valueIDParam := ctx.Params("id")
	if valueIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Option value ID is required",
		})
	}

	valueID, err := strconv.Atoi(valueIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid option value ID format",
		})
	}

	var req dto.UpdateOptionValueRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.optionValueUseCase.UpdateOptionValue(ctx.Context(), valueID, &usecase.UpdateOptionValueRequest{
		Name:            req.Name,
		IsDefault:       req.IsDefault,
		AdditionalPrice: req.AdditionalPrice,
		DisplayOrder:    req.DisplayOrder,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option value updated successfully", response)
}

// DeleteOptionValue handles option value deletion
func (c *OptionValueController) DeleteOptionValue(ctx *fiber.Ctx) error {
	valueIDParam := ctx.Params("id")
	if valueIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Option value ID is required",
		})
	}

	valueID, err := strconv.Atoi(valueIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid option value ID format",
		})
	}

	err = c.optionValueUseCase.DeleteOptionValue(ctx.Context(), valueID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Option value deleted successfully", nil)
}
