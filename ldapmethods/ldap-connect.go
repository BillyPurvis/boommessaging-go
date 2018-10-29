package ldapmethods

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	logrus "github.com/sirupsen/logrus"
	ldap "gopkg.in/ldap.v2"
)

var validate *validator.Validate

// Fields is a slice of maps where the key is the device var ID and key is the value

// ConnectionDetails For LDAP
type ConnectionDetails struct {
	CustomerID  string `json:"customer_id"`
	Host        string
	Port        string
	CN          string
	BaseDN      string `json:"base_dn"`
	Identifier  string
	Password    string                 `validate:"required"`
	RequestID   string                 `json:"request_id"`
	Fields      map[string]string      `json:"fields,omitempty"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
}

type resultsSet struct {
	Result []map[string]string
}

// LDAPConnectionBind Returns LDAP Connection Binding
func LDAPConnectionBind(connectionDetails *ConnectionDetails) (*ldap.Conn, error) {
	// By default, Port 389 should support SSL. If not, the user can define port 636 which is SSL.
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%v", connectionDetails.Host, connectionDetails.Port))
	if err != nil {
		return nil, err
	}

	// Create LDAP Binding
	err = conn.Bind(connectionDetails.Identifier, connectionDetails.Password)
	if err != nil {
		return nil, err
	}

	// Return connection binding
	return conn, nil
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

// Take a map and convert it to a string representation of the map
func mapToStringRepresentation(hashMap map[string]string) []string {
	// Get attribute fields from POST Body as string
	var attributeFields []string
	for _, value := range hashMap {
		attributeFields = append(attributeFields, value)
	}
	return attributeFields
}

// GetEntries Return results from LDAP
func GetEntries(connectionDetails *ConnectionDetails) ([]map[string]interface{}, error) {

	// Validate Struct first
	validate = validator.New()
	err := validate.Struct(connectionDetails)
	if err != nil {
		return nil, err
	}

	// Make LDAP Connection
	conn, err := LDAPConnectionBind(connectionDetails)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Build search query and parameters.
	searchQuery := buildSearchTermsString(connectionDetails)
	pagingControl := ldap.NewControlPaging(10)
	attributeFields := mapToStringRepresentation(connectionDetails.Fields)

	var recordTotal int
	for {
		// Make Search Request defining base DN, attributes and filters
		searchRequest := ldap.NewSearchRequest(
			fmt.Sprintf("cn=%v,dc=%v,dc=com,dc=local", connectionDetails.CN, connectionDetails.BaseDN),
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&%v)", searchQuery),
			attributeFields,
			[]ldap.Control{pagingControl},
		)

		records, err := conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		logrus.Info("LDAP Contacts batch count: ", len(records.Entries))
		fmt.Printf("\n====================\nRecord Count: %v\n====================\n", len(records.Entries))

		recordTotal += len(records.Entries)

		// Records is a batch of 1000 records, pipe them off to be inserted batch by batch

		// //
		// for _, entry := range records.Entries {
		// 	for _, field := range connectionDetails.Fields {
		// 		// TODO: Add Logic for saving to device vars
		// 		field = entry.GetAttributeValue(field)

		// 	}
		// }
		// Loop through fields and map and store.
		updatedControl := ldap.FindControl(records.Controls, ldap.ControlTypePaging)
		if ctrl, ok := updatedControl.(*ldap.ControlPaging); ctrl != nil && ok && len(ctrl.Cookie) != 0 {
			pagingControl.SetCookie(ctrl.Cookie)
			continue
		}
		break
	}

	fmt.Printf("\n====================\nRecord Count: %v\n====================\n", recordTotal)
	// // Make Search Request
	// sr, err := conn.Search(searchRequest)
	// if err != nil {
	// 	return nil, err
	// }

	//TODO: We don't need this, we just need to return a status
	// // Create list of maps
	// entryList := make([]map[string]interface{}, 0)
	// // Iterate through AD Records
	// for i, entry := range sr.Entries {

	// 	// Create new map item and append to list
	// 	entryList = append(entryList, make(map[string]interface{}))

	// 	// Iterate through requestedFields
	// 	for _, field := range connectionDetails.Fields {

	// 		fieldValue := entry.GetAttributeValue(field)

	// 		uuid := uuid.CreateUUID()
	// 		entryList[i]["uuid"] = uuid
	// 		// Check for empty fields and assign nil if empty
	// 		if fieldValue != "" {
	// 			entryList[i][field] = fieldValue
	// 		} else {
	// 			entryList[i][field] = nil
	// 		}
	// 	}
	// }
	m := make([]map[string]interface{}, 0)
	return m, nil
}

// GetEntryAttributes Returns attribute field lists for an entry
func GetEntryAttributes(connectionDetails *ConnectionDetails) ([]string, error) {

	conn, err := LDAPConnectionBind(connectionDetails)
	if err != nil {
		return nil, err
	}
	defer conn.Close() // Defer until end of function

	// Make Search request
	searchRequest := ldap.NewSearchRequest(
		fmt.Sprintf("dc=%v,dc=com,dc=local", connectionDetails.BaseDN),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&)", //TODO We don't need to filter?
		[]string{},
		nil,
	)

	// Make Search Request
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// Assign Attributes slice to var
	attributesSlice := sr.Entries[0].Attributes

	// Create New Slice of attribute names and return
	var attributeNames []string
	for _, attribute := range attributesSlice {
		attributeNames = append(attributeNames, attribute.Name)
	}
	return attributeNames, nil
}
