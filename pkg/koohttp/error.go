package koohttp

const (
	APIErrorCodeInternalServerError = "internal_server_error"
	APIErrorCodeBadRequest          = "bad_request"
	APIErrorCodeUnauthorized        = "unauthorized"
	APIErrorCodeForbidden           = "forbidden"
	APIErrorCodeNotFound            = "not_found"
	APIErrorCodeRequestTimeout      = "request_timeout"
	APIErrorCodeConflict            = "conflict"
	APIErrorCodeUnprocessableEntity = "unprocessable_entity"
	APIErrorCodeServiceUnavailable  = "service_unavailable"
)

type APIError interface {
	error
	HTTPStatus() int
}

type APIResponseError struct {
	Status    int    `json:"status"`
	ErrorCode string `json:"errorCode"`
}

// Error implements the error interface.
func (e *APIResponseError) Error() string {
	return e.ErrorCode
}

// Is enables errors.Is() to be used on the apiError.
func (e *APIResponseError) Is(target error) bool {
	return e.ErrorCode == target.Error()
}

func (e *APIResponseError) HTTPStatus() int {
	return e.Status
}

func NewAPIError(status int, errorCode string) APIError {
	return &APIResponseError{
		Status:    status,
		ErrorCode: errorCode,
	}
}
