package response

import "github.com/shoelfikar/voucher-management-system/pkg/utils"

// Response wraps the utils.Response for HTTP delivery layer
type Response = utils.Response

// PaginationResponse wraps the utils.PaginationResponse for HTTP delivery layer
type PaginationResponse = utils.PaginationResponse

// SuccessResponse creates a success response
func SuccessResponse(data interface{}) Response {
	return utils.SuccessResponse(data)
}

// SuccessResponseWithMessage creates a success response with a custom message
func SuccessResponseWithMessage(message string, data interface{}) Response {
	return utils.SuccessResponseWithMessage(message, data)
}

// ErrorResponse creates an error response
func ErrorResponse(message string) Response {
	return utils.ErrorResponse(message)
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(errors interface{}) Response {
	return utils.ValidationErrorResponse(errors)
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(data interface{}, page, limit int, total int64) PaginationResponse {
	return utils.PaginatedResponse(data, page, limit, total)
}
