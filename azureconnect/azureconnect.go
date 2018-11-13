package azureconnect

// ConnectionDetails - user stored details
type ConnectionDetails struct {
	ClientID    string                 `json:"client_id"`
	CustomerID  string                 `json:"customer_id,validate:required"`
	RequestID   string                 `json:"request_id,validate:required"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
}

// GraphAPIDetails - Details for the API connection to Microsft's Graph Api
type GraphAPIDetails struct {
	version            string // V1.0
	resource           string
	AuthorizationToken string // Bearer // It's a 1080 char token so we'll use a text field in MySQL
}

// NewGraphDetails - Returns graph defaults
func NewGraphDetails(authtoken string, resource string) *GraphAPIDetails {
	details := new(GraphAPIDetails)
	details.version = "V1.0"
	details.resource = resource
	details.AuthorizationToken = authtoken
	return details
}
