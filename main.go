package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BillyPurvis/boommessaging-go/database"
	"github.com/BillyPurvis/boommessaging-go/ldaphandler"
	"github.com/BillyPurvis/boommessaging-go/middleware"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Configure Logger
	f, fileErr := os.OpenFile("./error.log", os.O_WRONLY|os.O_CREATE, 0755)
	if fileErr != nil {
		log.Fatal(fileErr) // We must have log working
	}
	logrus.SetOutput(f)

	// Get Credentials
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	databaseCredentials := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", dbUsername, dbPassword, dbHost, dbName)

	// Open connection to DB
	var err error
	database.DBCon, err = sql.Open("mysql", databaseCredentials)
	defer database.DBCon.Close()
	if err != nil {
		log.Fatal(err) // we must have this working
	}

	fmt.Printf("Starting Server on port %v:%v\n", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))

	// Create Go Server
	router := httprouter.New()

	router.POST("/ldap/attributes", middleware.AuthenticateWare(ldaphandler.GetAttributes))
	router.POST("/ldap/contacts", middleware.AuthenticateWare(ldaphandler.GetContacts))

	log.Fatal(http.ListenAndServe(":4000", middleware.SetJSONHeader(router)))
}
