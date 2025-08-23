// internal/adapter/controller/kitchen_controller.go
package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// KitchenController handles HTTP requests related to kitchen operations
type KitchenController struct {
	kitchenUseCase        usecase.KitchenUsecase
	kitchenStationUsecase usecase.KitchenStationUsecase
	errorPresenter        presenter.ErrorPresenter
}

// NewKitchenController creates a new instance of KitchenController
func NewKitchenController(kitchenUseCase usecase.KitchenUsecase, kitchenStationUsecase usecase.KitchenStationUsecase, errorPresenter presenter.ErrorPresenter) *KitchenController {
	return &KitchenController{
		kitchenUseCase:        kitchenUseCase,
		errorPresenter:        errorPresenter,
		kitchenStationUsecase: kitchenStationUsecase,
	}
}

// GetKitchenQueue handles getting kitchen queue
func (c *KitchenController) GetKitchenQueue(ctx *fiber.Ctx) error {
	response, err := c.kitchenUseCase.GetKitchenQueue(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Kitchen queue retrieved successfully", response)
}

// UpdateOrderItemStatus handles updating order item status
func (c *KitchenController) UpdateOrderItemStatus(ctx *fiber.Ctx) error {
	orderItemIDParam := ctx.Params("orderItemId")
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
			Message: "Invalid order item ID format",
		})
	}

	var req dto.UpdateOrderItemStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.kitchenUseCase.UpdateOrderItemStatus(ctx.Context(), orderItemID, req.Status)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order item status updated successfully", response)
}

// GetOrderItemsByStatus handles getting order items by status
func (c *KitchenController) GetOrderItemsByStatus(ctx *fiber.Ctx) error {
	status := ctx.Query("status")
	if status == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Status query parameter is required",
		})
	}

	response, err := c.kitchenUseCase.GetOrderItemsByStatus(ctx.Context(), status)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order items retrieved successfully", response)
}

// MarkOrderItemAsReady handles marking order item as ready
func (c *KitchenController) MarkOrderItemAsReady(ctx *fiber.Ctx) error {
	orderItemIDParam := ctx.Params("orderItemId")
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
			Message: "Invalid order item ID format",
		})
	}

	response, err := c.kitchenUseCase.MarkOrderItemAsReady(ctx.Context(), orderItemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order item marked as ready successfully", response)
}

// MarkOrderItemAsServed handles marking order item as served
func (c *KitchenController) MarkOrderItemAsServed(ctx *fiber.Ctx) error {
	orderItemIDParam := ctx.Params("orderItemId")
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
			Message: "Invalid order item ID format",
		})
	}

	response, err := c.kitchenUseCase.MarkOrderItemAsServed(ctx.Context(), orderItemID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Order item marked as served successfully", response)
}

// GetKitchenOrdersByStation handles getting orders by kitchen station
func (c *KitchenController) GetKitchenOrdersByStation(ctx *fiber.Ctx) error {
	station := ctx.Query("station")
	if station == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Station query parameter is required",
		})
	}

	response, err := c.kitchenUseCase.GetKitchenOrdersByStation(ctx.Context(), station)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Kitchen orders by station retrieved successfully", response)
}

// type KitchenStationUsecase interface {
// 	CreateKitchenStation(ctx context.Context, req *CreateKitchenStationRequest) (*KitchenStationOnlyResponse, error)
// 	GetKitchenStation(ctx context.Context, id int) (*KitchenStationOnlyResponse, error)
// 	UpdateKitchenStation(ctx context.Context, id int, req *UpdateKitchenStationRequest) (*KitchenStationOnlyResponse, error)
// 	DeleteKitchenStation(ctx context.Context, id int) error
// 	ListKitchenStations(ctx context.Context) ([]*KitchenStationOnlyResponse, error)
// }

// CreateKitchenStatation handles menu option creation
func (c *KitchenController) CreateKitchenStatation(ctx *fiber.Ctx) error {
	var req dto.CreateKitchenStatationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.kitchenStationUsecase.CreateKitchenStation(ctx.Context(), &usecase.CreateKitchenStationRequest{
		Name: req.Name,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "kitchen station created successfully", response)
}

// GetKitchenStatation handles getting menu option by ID
func (c *KitchenController) GetKitchenStatation(ctx *fiber.Ctx) error {
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

	response, err := c.kitchenStationUsecase.GetKitchenStation(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "kitchen station retrieved successfully", response)
}

func (c *KitchenController) UpdateKitchenStatation(ctx *fiber.Ctx) error {
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

	var req dto.UpdateKitchenStatationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response, err := c.kitchenStationUsecase.UpdateKitchenStation(ctx.Context(), optionID, &usecase.UpdateKitchenStationRequest{
		Name: req.Name,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Menu option updated successfully", response)
}

func (c *KitchenController) DeleteKitchenStatation(ctx *fiber.Ctx) error {
	optionIDParam := ctx.Params("id")
	if optionIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "kitchen is required",
		})
	}

	optionID, err := strconv.Atoi(optionIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid kitchen format",
		})
	}

	err = c.kitchenStationUsecase.DeleteKitchenStation(ctx.Context(), optionID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "kitchen deleted successfully", nil)
}

func (c *KitchenController) ListKitchenStatations(ctx *fiber.Ctx) error {
	response, err := c.kitchenStationUsecase.ListKitchenStations(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "kitchen retrieved successfully", response)
}
