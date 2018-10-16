package ldapmethods

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BillyPurvis/boommessaging-go/uuid"
	ldap "gopkg.in/ldap.v2"
)

// ConnectionDetails For LDAP
type ConnectionDetails struct {
	CustomerID  string `json:"customer_id"`
	Host        string
	Port        string
	BaseDN      string `json:"base_dn"`
	Identifier  string
	Password    string
	Fields      []string               `json:"fields,omitempty"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
	Limit       string                 `json:"limit"`
}

// LDAPConnectionBind Returns LDAP Connection Binding
func LDAPConnectionBind(connectionDetails *ConnectionDetails) *ldap.Conn {
	// Create Connection to LDAP Server
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%v", connectionDetails.Host, connectionDetails.Port))
	if err != nil {
		panic(err)
	}

	// Create LDAP Binding
	err = conn.Bind(connectionDetails.Identifier, connectionDetails.Password)
	if err != nil {
		panic(err)
	}

	// Return connection binding
	return conn
}

// GetEntries Return results from LDAP
func GetEntries(connectionDetails *ConnectionDetails) []map[string]interface{} {
	conn := LDAPConnectionBind(connectionDetails)
	defer conn.Close() // Defer until end of function

	// Build concatinated byte slice of all filter options
	// (&(attribute=value/regex)(attribute=value))
	var searchQuery strings.Builder
	for serachAttribute, searchTerm := range connectionDetails.QueryParams {
		searchQuery.WriteString(fmt.Sprintf("(%v=%v)", serachAttribute, searchTerm))
	}

	filters := fmt.Sprintf("(&%v)", searchQuery.String())
	// Fields to be retrieved, IE, name, mail etc.
	attributes := connectionDetails.Fields

	// Pagination
	// We need a 32 bit in.
	pageSize, _ := strconv.ParseUint(connectionDetails.Limit, 10, 64)
	pageSizeuint := uint32(pageSize)

	pagingControl := ldap.NewControlPaging(pageSizeuint)
	controls := []ldap.Control{pagingControl}

	// Make Search Request defining base DN, attributes and filters
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", connectionDetails.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filters,
		attributes,
		controls,
	)

	// Make Search Request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	// Create list of maps
	entryList := make([]map[string]interface{}, 0)
	// Iterate through AD Records
	for i, entry := range sr.Entries {

		// Create new map item and append to list
		entryList = append(entryList, make(map[string]interface{}))

		// Iterate through requestedFields
		for _, field := range connectionDetails.Fields {

			fieldValue := entry.GetAttributeValue(field)

			uuid := uuid.CreateUUID()
			entryList[i]["uuid"] = uuid
			// Check for empty fields and assign nil if empty
			if fieldValue != "" {
				entryList[i][field] = fieldValue //
			} else {
				entryList[i][field] = nil
			}
		}
	}

	return entryList
}

// GetEntryAttributes Returns attribute field lists for an entry
func GetEntryAttributes(connectionDetails *ConnectionDetails) []string {

	conn := LDAPConnectionBind(connectionDetails)
	defer conn.Close() // Defer until end of function

	// Make Search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", connectionDetails.BaseDN),
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
