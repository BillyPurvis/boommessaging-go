package response

import (
	"encoding/json"
	"net/http"
)

// HTTPResponseBody Returns error
type HTTPResponseBody struct {
	Message string
	Status  int
}

// JSONResponse - Custom response handler
type JSONResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Count   int    `json:"count"`
}

// HTTPResponse Returns HTTP Response
func HTTPResponse(w http.ResponseWriter, responseMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(HTTPResponseBody{responseMsg, statusCode})
}
