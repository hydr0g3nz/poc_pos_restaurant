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

// GenerateQRCode handles generating QR code for a table
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

	qrCode, err := c.tableUsecase.GenerateQRCode(ctx.Context(), tableID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	response := map[string]interface{}{
		"table_id": tableID,
		"qr_code":  qrCode,
	}

	return SuccessResp(ctx, fiber.StatusOK, "QR code generated successfully", response)
}

// ScanQRCode handles QR code scanning with complete order information
func (c *TableController) ScanQRCode(ctx *fiber.Ctx) error {
	qrCode := ctx.Query("qr_code")
	if qrCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "QR code is required",
		})
	}

	// Use QRCodeUsecase to get complete information including order status
	response, err := c.qrCodeUsecase.ScanQRCode(ctx.Context(), qrCode)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	// Convert to DTO format
	dtoResponse := &dto.QRCodeScanResponse{
		TableID: response.TableID,
		Table: &dto.TableResponse{
			ID:          response.Table.ID,
			TableNumber: response.Table.TableNumber,
			QRCode:      response.Table.QRCode,
			Seating:     response.Table.Seating,
		},
		HasOpenOrder: response.HasOpenOrder,
	}

	// Include open order if exists
	if response.OpenOrder != nil {
		dtoResponse.OpenOrder = &dto.OrderResponse{
			ID:        response.OpenOrder.ID,
			TableID:   response.OpenOrder.TableID,
			Status:    response.OpenOrder.Status,
			CreatedAt: response.OpenOrder.CreatedAt,
			ClosedAt:  response.OpenOrder.ClosedAt,
		}
	}

	return SuccessResp(ctx, fiber.StatusOK, "QR code scanned successfully", dtoResponse)
}

// CreateOrderFromQRCode handles creating order directly from QR code scan
func (c *TableController) CreateOrderFromQRCode(ctx *fiber.Ctx) error {
	qrCode := ctx.Query("qr_code")
	if qrCode == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "QR code is required",
		})
	}

	// Create order from QR code
	response, err := c.qrCodeUsecase.CreateOrderFromQRCode(ctx.Context(), qrCode)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	// Convert to DTO format
	dtoResponse := &dto.OrderResponse{
		ID:        response.ID,
		TableID:   response.TableID,
		Status:    response.Status,
		CreatedAt: response.CreatedAt,
		ClosedAt:  response.ClosedAt,
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Order created from QR code successfully", dtoResponse)
}

// RegisterRoutes registers the routes for the table controller
func (c *TableController) RegisterRoutes(router fiber.Router) {
	tableGroup := router.Group("/tables")

	// Table CRUD operations
	tableGroup.Post("/", c.CreateTable)
	tableGroup.Get("/", c.ListTables)
	tableGroup.Get("/:id", c.GetTable)
	tableGroup.Put("/:id", c.UpdateTable)
	tableGroup.Delete("/:id", c.DeleteTable)

	// Table by number
	tableGroup.Get("/number/:number", c.GetTableByNumber)

	// QR code operations
	tableGroup.Get("/qr", c.GetTableByQRCode)               // GET /tables/qr?qr_code=/order?table=1
	tableGroup.Post("/:id/qr-code", c.GenerateQRCode)       // POST /tables/1/qr-code
	tableGroup.Get("/scan", c.ScanQRCode)                   // GET /tables/scan?qr_code=/order?table=1
	tableGroup.Post("/scan/order", c.CreateOrderFromQRCode) // POST /tables/scan/order?qr_code=/order?table=1
}
