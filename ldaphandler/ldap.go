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
	//w.Header().Set("Content-Type", "application/json")
	// Decode request body into struct
	var credentials ldapmethods.ConnectionDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		panic(err)
	}

	// Get attributes and encode to struct
	// data := GetEntryAttributeNames(&credentials)
	data := ldapmethods.GetEntries(&credentials)
	result := DataFields{data}
	json.NewEncoder(w).Encode(result)

}

// // LDAPContacts Returns
// func LDAPContacts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

// }
