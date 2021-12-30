package main

import (
	database "github.com/muhammadnajie/kp/internal/pkg/db/mysql"
	"github.com/muhammadnajie/kp/router"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8090"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()

	routes := router.Routes()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, routes))
}
