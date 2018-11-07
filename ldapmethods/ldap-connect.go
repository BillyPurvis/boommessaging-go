package ldapmethods

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"github.com/BillyPurvis/boommessaging-go/database"
	logrus "github.com/sirupsen/logrus"
	ldap "gopkg.in/ldap.v2"
)

var validate *validator.Validate

// Fields is a slice of maps where the key is the device var ID and key is the value

// ConnectionDetails For LDAP
type ConnectionDetails struct {
	CustomerID  string                 `json:"customer_id,validate:required"`
	Host        string                 `validate:"required"`
	Port        string                 `validate:"required"`
	CN          string                 `validate:"required"`
	BaseDN      string                 `json:"base_dn,validate:required"`
	Identifier  string                 `validate:"required"`
	Password    string                 `validate:"required"`
	RequestID   string                 `json:"request_id,validate:required"`
	Fields      map[string]string      `json:"fields,omitempty"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
	BatchLimit  string                 `json:"batch_limit,omitempty"`
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

func calculateBatchLimit(connectionDetails *ConnectionDetails) (uint32, error) {
	var batchLimit uint32 = 1000
	if connectionDetails.BatchLimit != "" {
		customBatchLimit, err := convertStringToUint32(connectionDetails.BatchLimit)
		if err != nil {
			return 0, err
		}
		// If the incoming batch limit is more than 1000, don't update batch limit
		if customBatchLimit <= 1000 {
			batchLimit = customBatchLimit
		}
	}

	return batchLimit, nil
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

	batchLimit, err := calculateBatchLimit(connectionDetails)
	if err != nil {
		return nil, err
	}

	// Build search query and parameters.
	searchQuery := buildSearchTermsString(connectionDetails)
	pagingControl := ldap.NewControlPaging(batchLimit)
	attributeFields := mapToStringRepresentation(connectionDetails.Fields)

	db := database.DBCon

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

		//fmt.Printf("\n====================\nRecord Count: %v\n====================\n", len(records.Entries))

		recordTotal += len(records.Entries) // returns slice of structs [Entry]

		// Records is a batch of 1000 records, pipe them off to be inserted batch by batch
		valuePlaceholders := make([]string, 0, 150)
		values := make([]interface{}, 0, 150)

		for i, entry := range records.Entries {
			for fieldVarID, field := range connectionDetails.Fields {
				fieldValue := entry.GetAttributeValue(field)

				valuePlaceholders = append(valuePlaceholders, "(?,?,?,?)")
				values = append(values, connectionDetails.RequestID, fieldVarID, i, fieldValue)
			}
		}

		//  field_count * placeholder_count * record_counts = max batch count
		q := fmt.Sprintf("INSERT INTO device_integration_vars (device_request_id, devices_var_name_id,record_row_index, value) VALUES %s", strings.Join(valuePlaceholders, ","))
		stmt, err := db.Prepare(q)
		if err != nil {
			fmt.Print(err.Error())
		}

		res, err := stmt.Exec(values...)
		if err != nil {
			fmt.Print(err.Error())
		}

		count, _ := res.RowsAffected()

		logrus.Info(fmt.Sprintf("Rows Insert for RequestID: %v. Inserted %v rows", connectionDetails.RequestID, count))

		// // Loop through fields and map and store.
		// updatedControl := ldap.FindControl(records.Controls, ldap.ControlTypePaging)
		// if ctrl, ok := updatedControl.(*ldap.ControlPaging); ctrl != nil && ok && len(ctrl.Cookie) != 0 {
		// 	pagingControl.SetCookie(ctrl.Cookie)
		// 	continue
		// }
		break
	}

	fmt.Printf("\n====================\nRecord Count: %v\n====================\n", recordTotal)

	// We might need to return a success message
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
