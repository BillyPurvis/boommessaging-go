package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BillyPurvis/boommessaging-go/Middleware"
	"github.com/BillyPurvis/boommessaging-go/handler"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Printf("Starting Server on port %v:%v\n", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))
	// Create Go Server
	router := httprouter.New()

	//router.POST("/", handler.LDAPIndex)
	router.POST("/ldap", middleware.AuthenticateWare(handler.LDAPAttributes))

	log.Fatal(http.ListenAndServe(":4000", middleware.SetJSONHeader(router)))
}
