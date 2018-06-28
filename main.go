package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BillyPurvis/boom/handler"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
)

func main() {

	fmt.Println("Brining Server online...")

	// Create Go Server
	router := httprouter.New()

	router.POST("/", handler.LDAPIndex)

	log.Fatal(http.ListenAndServe(":4000", router))

}
