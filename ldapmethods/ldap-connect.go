package ldapmethods

import (
	"fmt"
	"log"
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
	CN          string
	BaseDN      string `json:"base_dn"`
	Identifier  string
	Password    string
	Fields      []string               `json:"fields,omitempty"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
	Limit       string                 `json:"limit,omitempty"`
}

// LDAPConnectionBind Returns LDAP Connection Binding
func LDAPConnectionBind(connectionDetails *ConnectionDetails) *ldap.Conn {
	// By default, Port 389 should support SSL. If not, the user can define port 636 which is SSL.
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

// Convert String to a unsigned 32 bit integer
func convertStringToUint32(stringInt string) (uint32, error) {
	uInt64, err := strconv.ParseUint(stringInt, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(uInt64), nil
}

// Build concatinated byte slice of all filter options
// (&(attribute=value/regex)(attribute=value))
func buildSearchTermsString(connectionDetails *ConnectionDetails) string {
	var searchQuery strings.Builder
	for serachAttribute, searchTerm := range connectionDetails.QueryParams {
		searchQuery.WriteString(fmt.Sprintf("(%v=%v)", serachAttribute, searchTerm))
	}
	return searchQuery.String()
}

// GetEntries Return results from LDAP
func GetEntries(connectionDetails *ConnectionDetails) ([]map[string]interface{}, error) {
	conn := LDAPConnectionBind(connectionDetails)
	defer conn.Close()

	searchQuery := buildSearchTermsString(connectionDetails)
	pageSizeuint, err := convertStringToUint32(connectionDetails.Limit)

	if err != nil {
		log.Fatal(err)
	}
	pagingControl := ldap.NewControlPaging(pageSizeuint)

	// Make Search Request defining base DN, attributes and filters
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("cn=%v,dc=%v,dc=com,dc=local", connectionDetails.CN, connectionDetails.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&%v)", searchQuery),
		connectionDetails.Fields,
		[]ldap.Control{pagingControl},
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
				entryList[i][field] = fieldValue
			} else {
				entryList[i][field] = nil
			}
		}
	}
	return entryList, nil
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
