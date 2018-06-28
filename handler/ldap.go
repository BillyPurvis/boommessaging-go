package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/ldap.v2"

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

	fmt.Println(credentials)

	LDAPSearch()
	//TODO: Decode Response

	//TODO:  Verify API KEY
	//TODO: If API KEY true, make LDAP connection
	//TODO: Return LDAP details back in JSON format.
	//TODO: use a go routine
	fmt.Fprintf(w, "Hello Go")
}

// LDAPSearch Return results from LDAP
func LDAPSearch() {

	// Pull Details from struct

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "77.75.124.181", 389))
	if err != nil {
		panic(err)
	}
	defer l.Close()

	err = l.Bind("LDAP", "Boom01$")

	if err != nil {
		panic(err)
	}

}
