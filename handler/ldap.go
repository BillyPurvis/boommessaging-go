package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ConnectionDetails For LDAP
// Uppercase struct fields denot public properties to be accessed
type ConnectionDetails struct {
	CustomerID int `json:"customer_id"`
	Host       string
	Port       int
	Identifier string
	Password   string
}

// LDAPIndex POST Endpoint to retrieve LDAP connection details from Boom API
func LDAPIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var credentials ConnectionDetails

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&credentials)

	// Check for errors in decoding
	if err != nil {
		panic(err)
	}

	// Make LDAP Connection
	// customerID := credentials.CustomerID
	// host := credentials.Host
	// port := credentials.Port
	// username := credentials.Identifier
	// password := credentials.Password

	fmt.Println(credentials)
	//TODO: Decode Response

	//TODO:  Verify API KEY
	//TODO: If API KEY true, make LDAP connection
	//TODO: Return LDAP details back in JSON format.
	//TODO: use a go routine
	fmt.Fprintf(w, "Hello Go")
}
