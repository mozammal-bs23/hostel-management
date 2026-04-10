package main

import (
	"log"
	"net/http"
	"os"

	"example.com/hostel-management/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.New()

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
