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

// LDAPIndex POST Endpoint to retrieve LDAP connection details from Boom API
func LDAPIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Decode request body into struct
	var credentials ConnectionDetails
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	// Check for errors in decoding
	if err != nil {
		panic(err)
	}
	// Make LDAP Connection
	LDAPSearch(&credentials)
}

// LDAPSearch Return results from LDAP
func LDAPSearch(credentials *ConnectionDetails) {

	// Create Connection to LDAP Server
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", credentials.Host, credentials.Port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create LDAP Binding
	err = conn.Bind(credentials.Identifier, credentials.Password)
	if err != nil {
		panic(err)
	}

	//TODO: Make request to return just field names from DN search

	// Make Search Request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", credentials.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=user))",
		[]string{"displayName", "mail"}, //TODO: create map of field names required to pass to string slice of required data from LDAP
		nil,
	)

	// Make Search request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		panic(err)
	}
	// Iterate through search results slice and print
	//TODO: Return them to PHP
	for _, entry := range sr.Entries {
		fmt.Printf("%v : %v\n", entry.GetAttributeValue("displayName"), entry.GetAttributeValue("mail"))
	}
}
