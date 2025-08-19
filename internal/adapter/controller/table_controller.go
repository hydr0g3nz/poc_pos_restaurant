package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// TableController handles HTTP requests related to table operations
type TableController struct {
	tableUsecase   usecase.TableUsecase
	qrCodeUsecase  usecase.QRCodeUsecase
	errorPresenter presenter.ErrorPresenter
}

// NewTableController creates a new instance of TableController
func NewTableController(tableUsecase usecase.TableUsecase, errorPresenter presenter.ErrorPresenter) *TableController {
	return &TableController{
		tableUsecase:   tableUsecase,
		errorPresenter: errorPresenter,
	}
}

// CreateTable handles table creation
func (c *TableController) CreateTable(ctx *fiber.Ctx) error {
	var req dto.CreateTableRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.TableNumber <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number is required and must be greater than 0",
		})
	}

	if req.Seating < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Seating capacity cannot be negative",
		})
	}

	response, err := c.tableUsecase.CreateTable(ctx.Context(), &usecase.CreateTableRequest{
		TableNumber: req.TableNumber,
		Seating:     req.Seating,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Table created successfully", response)
}

// GetTable handles getting table by ID
func (c *TableController) GetTable(ctx *fiber.Ctx) error {
	tableIDParam := ctx.Params("id")
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

	response, err := c.tableUsecase.GetTable(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table retrieved successfully", response)
}

// GetTableByNumber handles getting table by table number
func (c *TableController) GetTableByNumber(ctx *fiber.Ctx) error {
	tableNumberParam := ctx.Params("number")
	if tableNumberParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number is required",
		})
	}

	tableNumber, err := strconv.Atoi(tableNumberParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid table number format",
		})
	}

	response, err := c.tableUsecase.GetTableByNumber(ctx.Context(), tableNumber)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table retrieved successfully", response)
}

// GetTableByQRCode handles getting table by QR code
func (c *TableController) GetTableByQRCode(ctx *fiber.Ctx) error {
	qrCode := ctx.Query("qr_code")
	if qrCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "QR code is required",
		})
	}

	response, err := c.tableUsecase.GetTableByQRCode(ctx.Context(), qrCode)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table retrieved successfully", response)
}

// UpdateTable handles updating table
func (c *TableController) UpdateTable(ctx *fiber.Ctx) error {
	tableIDParam := ctx.Params("id")
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

	var req dto.UpdateTableRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.TableNumber <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number is required and must be greater than 0",
		})
	}

	if req.Seating < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Seating capacity cannot be negative",
		})
	}

	response, err := c.tableUsecase.UpdateTable(ctx.Context(), tableID, &usecase.UpdateTableRequest{
		TableNumber: req.TableNumber,
		Seating:     req.Seating,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table updated successfully", response)
}

// DeleteTable handles deleting a table
func (c *TableController) DeleteTable(ctx *fiber.Ctx) error {
	tableIDParam := ctx.Params("id")
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

	err = c.tableUsecase.DeleteTable(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table deleted successfully", nil)
}

// ListTables handles getting all tables
func (c *TableController) ListTables(ctx *fiber.Ctx) error {
	response, err := c.tableUsecase.ListTables(ctx.Context())
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	tableListResponse := &dto.TableListResponse{
		Tables: make([]*dto.TableResponse, len(response)),
		Total:  len(response),
	}

	for i, table := range response {
		tableListResponse.Tables[i] = &dto.TableResponse{
			ID:          table.ID,
			TableNumber: table.TableNumber,
			QRCode:      table.QRCode,
			Seating:     table.Seating,
		}
	}

	return SuccessResp(ctx, fiber.StatusOK, "Tables retrieved successfully", tableListResponse)
}
