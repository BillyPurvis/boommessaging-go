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

// GetAttributes Returns Attributes of an entry
func GetAttributes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request body into struct
	var credentials ldapmethods.ConnectionDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	// Get attributes and encode to struct
	data := ldapmethods.GetEntryAttributes(&credentials)
	result := DataFields{data}
	json.NewEncoder(w).Encode(result)
}

// GetContacts Returns Contacts
func GetContacts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var credentials ldapmethods.ConnectionDetails
	err := json.NewDecoder(r.Body).Decode(&credentials)

	// Check err
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	data, err := ldapmethods.GetEntries(&credentials)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w).Encode(data)
}
