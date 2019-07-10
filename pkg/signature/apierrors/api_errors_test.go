package apierrors

import (
	"github.com/stretchr/testify/assert"

	"net/http"

	"testing"
)

func Test_GetAPIError(t *testing.T) {
	apiError := GetAPIError(APIErrUserName)
	assert.Equal(t, http.StatusBadRequest, apiError.HTTPStatusCode)
	assert.Equal(t, "BadRequest", apiError.Code)
	assert.Contains(t, apiError.Description, "User Name is Error")

	apiError = GetAPIError(APIErrAccessDenied)
	assert.Equal(t, http.StatusForbidden, apiError.HTTPStatusCode)
	assert.Equal(t, "AccessDenied", apiError.Code)
	assert.Contains(t, apiError.Description, "Access Denied")

	apiError = GetAPIError(APIErrBadRequest)
	assert.Equal(t, http.StatusBadRequest, apiError.HTTPStatusCode)
	assert.Equal(t, "BadRequest", apiError.Code)
	assert.Contains(t, apiError.Description, "Bad request")

	apiError = GetAPIError(APIErrExpiredToken)
	assert.Equal(t, http.StatusBadRequest, apiError.HTTPStatusCode)
	assert.Equal(t, "ExpiredToken", apiError.Code)
	assert.Contains(t, apiError.Description, "provided token has expi")

	apiError = GetAPIError(APIErrInvalidToken)
	assert.Equal(t, http.StatusUnauthorized, apiError.HTTPStatusCode)
	assert.Equal(t, "InvalidToken", apiError.Code)
	assert.Contains(t, apiError.Description, "malformed or otherwise")
}

func Test_GetAPIErrorResponse(t *testing.T) {

	apiError := GetAPIError(APIErrUserName)
	apiErrorResponse := GetAPIErrorResponse(apiError, "")
	assert.Equal(t, "BadRequest", apiError.Code)
	assert.Contains(t, apiErrorResponse.Message, "User Name is Error")

	apiError = GetAPIError(APIErrAccessDenied)
	apiErrorResponse = GetAPIErrorResponse(apiError, "")
	assert.Equal(t, "AccessDenied", apiError.Code)
	assert.Contains(t, apiErrorResponse.Message, "Access Denied")

	apiError = GetAPIError(APIErrBadRequest)
	apiErrorResponse = GetAPIErrorResponse(apiError, "")
	assert.Equal(t, "BadRequest", apiError.Code)
	assert.Contains(t, apiErrorResponse.Message, "Bad request")

	apiError = GetAPIError(APIErrExpiredToken)
	apiErrorResponse = GetAPIErrorResponse(apiError, "")
	assert.Equal(t, "ExpiredToken", apiError.Code)
	assert.Contains(t, apiErrorResponse.Message, "provided token has expi")

	apiError = GetAPIError(APIErrInvalidToken)
	apiErrorResponse = GetAPIErrorResponse(apiError, "")
	assert.Equal(t, "InvalidToken", apiError.Code)
	assert.Contains(t, apiErrorResponse.Message, "malformed or otherwise")

	assert.Equal(t, 11, len(ErrorCodeResponse))

	apiError.SetDescription("test me")
	assert.Equal(t, "test me", apiError.Description)
}
