package utils

import "math"

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// PaginationResponse represents a paginated API response
type PaginationResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) Response {
	return Response{
		Status: "success",
		Data:   data,
	}
}

// SuccessResponseWithMessage creates a success response with a custom message
func SuccessResponseWithMessage(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string) Response {
	return Response{
		Status:  "error",
		Message: message,
	}
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(errors interface{}) Response {
	return Response{
		Status: "error",
		Message: "Validation failed",
		Errors: errors,
	}
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(data interface{}, page, limit int, total int64) PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}
}
