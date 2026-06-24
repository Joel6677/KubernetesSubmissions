package main

import (
	"fmt"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Moi!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", getRoot)

	http.ListenAndServe(":"+port, nil)
}
