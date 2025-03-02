package main

import (
	"ASSIGNMENT_1/HANDLERS"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	r := http.NewServeMux()

	r.HandleFunc("/countryinfo/v1/info/", HANDLERS.Test)
	r.HandleFunc("/countryinfo/v1/population/", HANDLERS.Population)
	r.HandleFunc("/status", HANDLERS.Status)

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
