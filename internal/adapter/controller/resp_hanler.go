// Updated error handling in resp_handler.go to include menu item errors

package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
)

type successResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SuccessResp builds a success response
func SuccessResp(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(successResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func HandleError(c *fiber.Ctx, err error, errorPresenter presenter.ErrorPresenter) error {
	errorResp := errorPresenter.PresentError(err)
	return c.Status(errorResp.Status).JSON(errorResp)
}
