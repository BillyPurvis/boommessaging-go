package ldaphandler

import (
	"encoding/json"
	"net/http"

	"github.com/BillyPurvis/boommessaging-go/ldapmethods"
	"github.com/BillyPurvis/boommessaging-go/response"
	"github.com/julienschmidt/httprouter"
)

// DataFields Field list from LDAP
type DataFields struct {
	Fields []string `json:"entry_attributes"`
}

// GetAttributes Returns Attributes of an entry
func GetAttributes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// TODO move to function
	var credentials ldapmethods.ConnectionDetails
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		response.HTTPResponse(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	data, err := ldapmethods.GetEntryAttributes(&credentials)
	if err != nil {
		response.HTTPResponse(w, err.Error(), http.StatusBadRequest)
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
		response.HTTPResponse(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	res, err := ldapmethods.GetEntries(&credentials)
	if err != nil {
		response.HTTPResponse(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(res)
}
