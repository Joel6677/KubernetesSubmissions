package main

import (
	"fmt"
	"net/http"
	"os"
)

var pings = 0

func pingpong(w http.ResponseWriter, r *http.Request) {
	pings++
	fmt.Fprintf(w, "pong %d", pings)
}

func pingCounter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", pings)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/pingpong", pingpong)

	http.HandleFunc("/pings", pingCounter)

	http.ListenAndServe(":"+port, nil)
}
