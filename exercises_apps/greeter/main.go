package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	greeting := os.Getenv("GREETING")
	if greeting == "" {
		greeting = "hello"
	}
	version := os.Getenv("VERSION")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s (from %s)", greeting, version)
	})

	log.Println("greeter started")
	http.ListenAndServe(":8080", nil)
}
