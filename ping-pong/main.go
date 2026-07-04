package main

import (
	"fmt"
	"net/http"
	"os"
)

var count = 0

func pingpong(w http.ResponseWriter, r *http.Request) {
	count++
	fmt.Fprintf(w, "pong %d", count)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/pingpong", pingpong)

	http.ListenAndServe(":"+port, nil)
}
