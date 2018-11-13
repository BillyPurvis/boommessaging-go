package azurehandler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetContacts - returns contacts from Azure
func GetContacts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Hello, Azure.")
}
