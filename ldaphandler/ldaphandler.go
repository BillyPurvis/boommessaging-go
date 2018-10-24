package ldaphandler

import (
	"encoding/json"
	"net/http"

	"github.com/BillyPurvis/boommessaging-go/ldapmethods"
	"github.com/julienschmidt/httprouter"
)

// DataFields Field list from LDAP
type DataFields struct {
	Fields []string `json:"entry_attributes"`
}

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

// GetAttributes Returns Attributes of an entry
func GetAttributes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// TODO move to function
	var credentials ldapmethods.ConnectionDetails
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		HTTPResponse(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	data, err := ldapmethods.GetEntryAttributes(&credentials)
	if err != nil {
		HTTPResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := DataFields{data}
	json.NewEncoder(w).Encode(result)
}

// GetContacts Returns Contacts
func GetContacts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// TODO move to function
	var credentials ldapmethods.ConnectionDetails
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		HTTPResponse(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	data, err := ldapmethods.GetEntries(&credentials)
	if err != nil {
		HTTPResponse(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(data)
}
