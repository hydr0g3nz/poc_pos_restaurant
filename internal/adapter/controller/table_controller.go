package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// TableController handles HTTP requests related to table operations
type TableController struct {
	tableUseCase usecase.TableUsecase
}

// NewTableController creates a new instance of TableController
func NewTableController(tableUseCase usecase.TableUsecase) *TableController {
	return &TableController{
		tableUseCase: tableUseCase,
	}
}

// CreateTable handles table creation
func (c *TableController) CreateTable(ctx *fiber.Ctx) error {
	var req dto.CreateTableRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err)
	}

	if req.TableNumber <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number is required and must be greater than 0",
		})
	}

	response, err := c.tableUseCase.CreateTable(ctx.Context(), &usecase.CreateTableRequest{
		TableNumber: req.TableNumber,
		Seating:     req.Seating,
	})
	if err != nil {
		return HandleError(ctx, err)
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

	response, err := c.tableUseCase.GetTable(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table retrieved successfully", response)
}

// GetTableByNumber handles getting table by number
func (c *TableController) GetTableByNumber(ctx *fiber.Ctx) error {
	tableNumber := ctx.Query("number")
	if tableNumber == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number query parameter is required",
		})
	}

	number, err := strconv.Atoi(tableNumber)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid table number format",
		})
	}

	response, err := c.tableUseCase.GetTableByNumber(ctx.Context(), number)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table retrieved successfully", response)
}

// GetTableByQRCode handles getting table by QR code
func (c *TableController) GetTableByQRCode(ctx *fiber.Ctx) error {
	qrCode := ctx.Query("qr_code")
	if qrCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "QR code query parameter is required",
		})
	}

	response, err := c.tableUseCase.GetTableByQRCode(ctx.Context(), qrCode)
	if err != nil {
		return HandleError(ctx, err)
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
		return HandleError(ctx, err)
	}

	if req.TableNumber <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Table number is required and must be greater than 0",
		})
	}

	response, err := c.tableUseCase.UpdateTable(ctx.Context(), tableID, &usecase.UpdateTableRequest{
		TableNumber: req.TableNumber,
		Seating:     req.Seating,
	})
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table updated successfully", response)
}

// DeleteTable handles table deletion
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

	err = c.tableUseCase.DeleteTable(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Table deleted successfully", nil)
}

// ListTables handles getting all tables
func (c *TableController) ListTables(ctx *fiber.Ctx) error {
	response, err := c.tableUseCase.ListTables(ctx.Context())
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Tables retrieved successfully", response)
}

// GenerateQRCode handles generating QR code for table
func (c *TableController) GenerateQRCode(ctx *fiber.Ctx) error {
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

	qrCode, err := c.tableUseCase.GenerateQRCode(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err)
	}

	return SuccessResp(ctx, fiber.StatusOK, "QR code generated successfully", map[string]string{"qr_code": qrCode})
}

// RegisterRoutes registers the routes for the table controller
func (c *TableController) RegisterRoutes(router fiber.Router) {
	tableGroup := router.Group("/tables")

	// Public routes (for customers)
	tableGroup.Get("/", c.ListTables)
	tableGroup.Get("/search", c.GetTableByNumber) // GET /tables/search?number=1
	tableGroup.Get("/qr", c.GetTableByQRCode)     // GET /tables/qr?qr_code=/order?table=1
	tableGroup.Get("/:id", c.GetTable)
	tableGroup.Post("/:id/qr", c.GenerateQRCode)

	// Admin routes (require admin role in real implementation)
	tableGroup.Post("/", c.CreateTable)
	tableGroup.Put("/:id", c.UpdateTable)
	tableGroup.Delete("/:id", c.DeleteTable)
}
