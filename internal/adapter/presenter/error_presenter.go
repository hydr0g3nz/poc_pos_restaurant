package presenter

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
)

type ErrorPresenter interface {
	PresentError(err error) *ErrorResponse
}

type errorPresenter struct {
	logger infra.Logger
}

func NewErrorPresenter(logger infra.Logger) ErrorPresenter {
	return &errorPresenter{
		logger: logger,
	}
}

type ErrorResponse struct {
	Status    int                    `json:"status"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Category  string                 `json:"category,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	TraceID   string                 `json:"trace_id"`
}

func (p *errorPresenter) PresentError(err error) *ErrorResponse {
	traceID := generateTraceID()
	timestamp := time.Now()

	// Log error for debugging
	p.logger.Error("Error occurred",
		"error", err.Error(),
		"trace_id", traceID,
		"timestamp", timestamp,
		"type", fmt.Sprintf("%T", err),
	)

	response := &ErrorResponse{
		Timestamp: timestamp,
		TraceID:   traceID,
	}

	// Handle domain errors
	if domainErr, ok := err.(errs.DomainError); ok {
		response.Status = domainErr.HTTPStatus()
		response.Code = domainErr.Code()
		response.Message = domainErr.Message()
		response.Category = string(domainErr.Category())
		response.Details = domainErr.Details()

		// Log domain error details
		p.logger.Info("Domain error handled",
			"code", domainErr.Code(),
			"category", domainErr.Category(),
			"details", domainErr.Details(),
			"trace_id", traceID,
		)

		return response
	}

	// Handle other error types
	switch err.Error() {
	case "record not found", "sql: no rows in result set":
		response.Status = 404
		response.Code = "RESOURCE_NOT_FOUND"
		response.Message = "Resource not found"
		response.Category = string(errs.CategoryNotFound)

	default:
		// Unknown/Infrastructure errors (don't expose details)
		response.Status = 500
		response.Code = "INTERNAL_ERROR"
		response.Message = "Internal server error"
		response.Category = string(errs.CategoryInternal)

		// Log detailed error for debugging
		p.logger.Error("Unhandled error",
			"error", err.Error(),
			"type", fmt.Sprintf("%T", err),
			"trace_id", traceID,
		)
	}

	return response
}

func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("trace_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
