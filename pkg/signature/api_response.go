package signature

import (
	"bytes"
	"encoding/json"

	"github.com/eleme/lindb/pkg/signature/apierrors"

	"net/http"
	"net/url"
)

// writeSuccessResponseXML writes success headers and response if any,
// with content-type set to `application/xml`.
func writeSuccessResponse(w http.ResponseWriter, response []byte) {
	writeResponse(w, http.StatusOK, response)
}

// Encodes the response headers into XML format.
func encodeResponse(response interface{}) []byte {
	var bytesBuffer bytes.Buffer

	e := json.NewEncoder(&bytesBuffer)
	_ = e.Encode(response)
	return bytesBuffer.Bytes()
}

func writeResponse(w http.ResponseWriter, statusCode int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response != nil {
		_, _ = w.Write(response)
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// writeErrorRespone writes error headers
func writeErrorResponse(w http.ResponseWriter, errorCode apierrors.APIErrorCode, reqURL *url.URL) {
	apiError := apierrors.GetAPIError(errorCode)
	errorResponse := apierrors.GetAPIErrorResponse(apiError, reqURL.Path)
	encodedErrorResponse := encodeResponse(errorResponse)
	writeResponse(w, apiError.HTTPStatusCode, encodedErrorResponse)
}
