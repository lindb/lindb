package apierrors

import (
	"net/http"
)

type APIError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

func (a *APIError) SetDescription(desc string) {
	a.Description = desc
}

// APIErrorResponse - error response format
type APIErrorResponse struct {
	Code           string
	HTTPStatusCode int
	Message        string
	Resource       string
}

type APIErrorCode int

const (
	APIErrAccessDenied APIErrorCode = iota

	APIErrBadRequest

	APIErrExpiredToken

	APIErrInvalidToken

	APIErrInvalidUnauthorized

	APIErrUnsignedHeaders

	APIErrInvalidAccessKeyID

	APIErrPasswordIsEmpty

	APIErrUserNameIsEmpty

	APIErrUserName

	APIErrPassword
)

var ErrorCodeResponse = map[APIErrorCode]APIError{
	APIErrInvalidAccessKeyID: {
		Code:           "InvalidAccessKeyId",
		Description:    "The access key ID you provided does not exist in our records.",
		HTTPStatusCode: http.StatusForbidden,
	},
	APIErrUnsignedHeaders: {
		Code:           "AccessDenied",
		Description:    "There were headers present in the request which were not signed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Access Denied.",
		HTTPStatusCode: http.StatusForbidden,
	},
	APIErrUserNameIsEmpty: {
		Code:           "BadRequest",
		Description:    "User Name is Empty",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrPasswordIsEmpty: {
		Code:           "BadRequest",
		Description:    "Password is Empty",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrUserName: {
		Code:           "BadRequest",
		Description:    "User Name is Error",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrPassword: {
		Code:           "BadRequest",
		Description:    "Password is Error",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrBadRequest: {
		Code:           "BadRequest",
		Description:    "Bad request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrExpiredToken: {
		Code:           "ExpiredToken",
		Description:    "The provided token has expired",
		HTTPStatusCode: http.StatusBadRequest,
	},
	APIErrInvalidToken: {
		Code:           "InvalidToken",
		Description:    "The provided token is malformed or otherwise invalid",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	APIErrInvalidUnauthorized: {
		Code:           "InvalidUnauthorized",
		Description:    "Unauthorized access to this resource",
		HTTPStatusCode: http.StatusUnauthorized,
	},
}

// GetAPIError provides API Error for input API error code.
func GetAPIError(code APIErrorCode) APIError {
	return ErrorCodeResponse[code]
}

// GetAPIErrorResponse gets in standard error and resource value and
// provides a encodable populated response values
func GetAPIErrorResponse(err APIError, resource string) APIErrorResponse {
	return APIErrorResponse{
		Code:           err.Code,
		Message:        err.Description,
		HTTPStatusCode: err.HTTPStatusCode,
		Resource:       resource,
	}
}
