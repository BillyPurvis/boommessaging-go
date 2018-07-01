package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	ldap "gopkg.in/ldap.v2"
)

// ConnectionDetails For LDAP
// Uppercase struct fields denot public properties to be accessed
type ConnectionDetails struct {
	CustomerID int `json:"customer_id"`
	Host       string
	Port       int
	BaseDN     string
	Identifier string
	Password   string
}

// DataFields Field list from LDAP
type DataFields struct {
	Fields []string `json:"entry_attributes"`
}

// LDAPAttributes Returns Attributes of an entry
func LDAPAttributes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// Decode request body into struct
	var credentials ConnectionDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	if err != nil {
		panic(err)
	}

	// Get attributes and encode to struct
	data := GetEntryAttributeNames(&credentials)
	result := DataFields{data}
	json.NewEncoder(w).Encode(result)

}

// LDAPIndex POST Endpoint to retrieve LDAP connection details from Boom API
func LDAPIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	//TODO: Move to middlewear
	w.Header().Set("Content-Type", "application/json")

	// Decode request body into struct
	var credentials ConnectionDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	// Check for errors in decoding
	if err != nil {
		panic(err)
	}

	// Make LDAP Connection
	data := GetEntries(&credentials)

	// Create new struct for JSON response body of attributes
	result := DataFields{data}
	json.NewEncoder(w).Encode(result)
}

// LDAPConnectionBind Returns LDAP Connection Binding
func LDAPConnectionBind(credentials *ConnectionDetails) *ldap.Conn {
	// Create Connection to LDAP Server
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", credentials.Host, credentials.Port))
	if err != nil {
		panic(err)
	}

	// Create LDAP Binding
	err = conn.Bind(credentials.Identifier, credentials.Password)
	if err != nil {
		panic(err)
	}

	// Return connection binding
	return conn
}

// GetEntries Return results from LDAP
func GetEntries(credentials *ConnectionDetails) []string {

	conn := LDAPConnectionBind(credentials)
	defer conn.Close() // Defer until end of function

	// Make Search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", credentials.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=user))",
		[]string{},
		nil,
	)

	// Make Search Request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	// Assign Attributes slice to var
	attributesSlice := sr.Entries[0].Attributes

	// Create New Slice of attribute names and return
	var attributeNames []string
	for _, attribute := range attributesSlice {
		attributeNames = append(attributeNames, attribute.Name)
	}
	return attributeNames
}

// GetEntryAttributeNames Returns attribute field lists for an entry
func GetEntryAttributeNames(credentials *ConnectionDetails) []string {

	conn := LDAPConnectionBind(credentials)
	defer conn.Close() // Defer until end of function

	// Make Search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", credentials.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=user))",
		[]string{},
		nil,
	)

	// Make Search Request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	// Assign Attributes slice to var
	attributesSlice := sr.Entries[0].Attributes

	// Create New Slice of attribute names and return
	var attributeNames []string
	for _, attribute := range attributesSlice {
		attributeNames = append(attributeNames, attribute.Name)
	}
	return attributeNames
}
