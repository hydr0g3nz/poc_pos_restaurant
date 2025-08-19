package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// RevenueController handles HTTP requests related to revenue operations
type RevenueController struct {
	revenueUsecase usecase.RevenueUsecase
	errorPresenter presenter.ErrorPresenter
}

// NewRevenueController creates a new instance of RevenueController
func NewRevenueController(revenueUsecase usecase.RevenueUsecase, errorPresenter presenter.ErrorPresenter) *RevenueController {
	return &RevenueController{
		revenueUsecase: revenueUsecase,
		errorPresenter: errorPresenter,
	}
}

// GetDailyRevenue handles getting daily revenue
func (c *RevenueController) GetDailyRevenue(ctx *fiber.Ctx) error {
	dateStr := ctx.Query("date")
	if dateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Date query parameter is required (format: YYYY-MM-DD)",
		})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid date format. Use YYYY-MM-DD",
		})
	}

	response, err := c.revenueUsecase.GetDailyRevenue(ctx.Context(), date)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Daily revenue retrieved successfully", response)
}

// GetMonthlyRevenue handles getting monthly revenue
func (c *RevenueController) GetMonthlyRevenue(ctx *fiber.Ctx) error {
	yearStr := ctx.Query("year")
	monthStr := ctx.Query("month")

	if yearStr == "" || monthStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Year and month query parameters are required",
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid year format",
		})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid month format",
		})
	}

	if month < 1 || month > 12 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Month must be between 1 and 12",
		})
	}

	response, err := c.revenueUsecase.GetMonthlyRevenue(ctx.Context(), year, month)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Monthly revenue retrieved successfully", response)
}

// GetDailyRevenueRange handles getting daily revenue for a date range
func (c *RevenueController) GetDailyRevenueRange(ctx *fiber.Ctx) error {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	response, err := c.revenueUsecase.GetDailyRevenueRange(ctx.Context(), startDate, endDate)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Daily revenue range retrieved successfully", response)
}

// GetMonthlyRevenueRange handles getting monthly revenue for a date range
func (c *RevenueController) GetMonthlyRevenueRange(ctx *fiber.Ctx) error {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	response, err := c.revenueUsecase.GetMonthlyRevenueRange(ctx.Context(), startDate, endDate)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Monthly revenue range retrieved successfully", response)
}

// GetTotalRevenue handles getting total revenue for a date range
func (c *RevenueController) GetTotalRevenue(ctx *fiber.Ctx) error {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	response, err := c.revenueUsecase.GetTotalRevenue(ctx.Context(), startDate, endDate)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Total revenue retrieved successfully", response)
}
