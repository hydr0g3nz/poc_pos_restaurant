package errs

import (
	"fmt"
	"net/http"
	"strings"
)

// DomainError interface for self-describing errors
type DomainError interface {
	error
	Code() string
	Message() string
	HTTPStatus() int
	Details() map[string]interface{}
	Category() ErrorCategory
	WithDetails(details map[string]interface{}) DomainError
	WithField(field string, value interface{}) DomainError
}

// ErrorCategory for grouping similar errors
type ErrorCategory string

const (
	CategoryValidation   ErrorCategory = "VALIDATION"
	CategoryNotFound     ErrorCategory = "NOT_FOUND"
	CategoryConflict     ErrorCategory = "CONFLICT"
	CategoryUnauthorized ErrorCategory = "UNAUTHORIZED"
	CategoryForbidden    ErrorCategory = "FORBIDDEN"
	CategoryInternal     ErrorCategory = "INTERNAL"
	CategoryBusiness     ErrorCategory = "BUSINESS_RULE"
)

// BaseDomainError implements DomainError interface
type BaseDomainError struct {
	code       string
	message    string
	httpStatus int
	category   ErrorCategory
	details    map[string]interface{}
}

func (e *BaseDomainError) Error() string           { return e.message }
func (e *BaseDomainError) Code() string            { return e.code }
func (e *BaseDomainError) Message() string         { return e.message }
func (e *BaseDomainError) HTTPStatus() int         { return e.httpStatus }
func (e *BaseDomainError) Category() ErrorCategory { return e.category }
func (e *BaseDomainError) Details() map[string]interface{} {
	if e.details == nil {
		return make(map[string]interface{})
	}
	return e.details
}

func (e *BaseDomainError) WithDetails(details map[string]interface{}) DomainError {
	newErr := *e
	if newErr.details == nil {
		newErr.details = make(map[string]interface{})
	}
	for k, v := range details {
		newErr.details[k] = v
	}
	return &newErr
}

func (e *BaseDomainError) WithField(field string, value interface{}) DomainError {
	return e.WithDetails(map[string]interface{}{field: value})
}

// ==========================================
// 2. Error Factory Functions
// ==========================================

// NewValidationError creates validation errors
func NewValidationError(field string, constraint string, value interface{}) DomainError {
	return &BaseDomainError{
		code:       "VALIDATION_ERROR",
		message:    fmt.Sprintf("Validation failed for field '%s': %s", field, constraint),
		httpStatus: http.StatusBadRequest,
		category:   CategoryValidation,
		details: map[string]interface{}{
			"field":      field,
			"constraint": constraint,
			"value":      value,
		},
	}
}

// NewNotFoundError creates not found errors
func NewNotFoundError(resource string, identifier interface{}) DomainError {
	resourceUpper := strings.ToUpper(strings.ReplaceAll(resource, " ", "_"))
	return &BaseDomainError{
		code:       fmt.Sprintf("%s_NOT_FOUND", resourceUpper),
		message:    fmt.Sprintf("%s not found", resource),
		httpStatus: http.StatusNotFound,
		category:   CategoryNotFound,
		details: map[string]interface{}{
			"resource":   resource,
			"identifier": identifier,
		},
	}
}

// NewConflictError creates conflict errors
func NewConflictError(resource string, reason string) DomainError {
	resourceUpper := strings.ToUpper(strings.ReplaceAll(resource, " ", "_"))
	return &BaseDomainError{
		code:       fmt.Sprintf("%s_CONFLICT", resourceUpper),
		message:    fmt.Sprintf("%s conflict: %s", resource, reason),
		httpStatus: http.StatusConflict,
		category:   CategoryConflict,
		details: map[string]interface{}{
			"resource": resource,
			"reason":   reason,
		},
	}
}

// NewBusinessRuleError creates business rule violation errors
func NewBusinessRuleError(rule string, context map[string]interface{}) DomainError {
	return &BaseDomainError{
		code:       "BUSINESS_RULE_VIOLATION",
		message:    fmt.Sprintf("Business rule violation: %s", rule),
		httpStatus: http.StatusConflict,
		category:   CategoryBusiness,
		details:    context,
	}
}

// NewUnauthorizedError creates unauthorized errors
func NewUnauthorizedError(reason string) DomainError {
	return &BaseDomainError{
		code:       "UNAUTHORIZED",
		message:    fmt.Sprintf("Unauthorized: %s", reason),
		httpStatus: http.StatusUnauthorized,
		category:   CategoryUnauthorized,
		details: map[string]interface{}{
			"reason": reason,
		},
	}
}

// NewForbiddenError creates forbidden errors
func NewForbiddenError(action string, resource string) DomainError {
	return &BaseDomainError{
		code:       "FORBIDDEN",
		message:    fmt.Sprintf("Forbidden: cannot %s %s", action, resource),
		httpStatus: http.StatusForbidden,
		category:   CategoryForbidden,
		details: map[string]interface{}{
			"action":   action,
			"resource": resource,
		},
	}
}
func NewExternalServiceError(service string, reason string) DomainError {
	return &BaseDomainError{
		code:       "EXTERNAL_SERVICE_ERROR",
		message:    fmt.Sprintf("External service error: %s - %s", service, reason),
		httpStatus: http.StatusServiceUnavailable,
		category:   CategoryExternal,
		details: map[string]interface{}{
			"service": service,
			"reason":  reason,
		},
	}
}
