package response

import (
	"encoding/json"
	"net/http"
)

// HTTPError Returns error
type HTTPError struct {
	Message string
	Status  int
}

// HTTPResponse Returns HTTP Response
func HTTPResponse(w http.ResponseWriter, responseMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(HTTPError{responseMsg, statusCode})
}
