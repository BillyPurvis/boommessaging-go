package ldapmethods

import (
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// ConnectionDetails For LDAP
type ConnectionDetails struct {
	CustomerID int `json:"customer_id"`
	Host       string
	Port       int
	BaseDN     string
	Identifier string
	Password   string
	Fields     []string `json:"fields,omitempty"`
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
func GetEntries(credentials *ConnectionDetails) []map[string]*string {

	conn := LDAPConnectionBind(credentials)
	defer conn.Close() // Defer until end of function

	//	l := len(credentials.Fields)

	// Make Search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", credentials.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=user))",
		credentials.Fields,
		nil,
	)

	// Make Search Request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	// Create list of maps
	entryList := make([]map[string]*string, 0)
	// Iterate through AD Records
	for i, entry := range sr.Entries {

		// Create new map item and append to list
		entryList = append(entryList, make(map[string]*string))

		// Iterate through requestedFields
		for _, field := range credentials.Fields {

			fieldValue := entry.GetAttributeValue(field)

			// Check for empty fields and assign nil if empty
			if fieldValue != "" {
				entryList[i][field] = &fieldValue //
			} else {
				entryList[i][field] = nil
			}
		}
	}

	return entryList
}

// GetEntryAttributes Returns attribute field lists for an entry
func GetEntryAttributes(credentials *ConnectionDetails) []string {

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
