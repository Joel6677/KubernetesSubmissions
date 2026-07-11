package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func getStatus(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("/shared/output.txt")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	pings, err := os.ReadFile("/shared/pings.txt")

	pCount := "0"

	if err == nil {
		pCount = strings.TrimSpace(string(pings))
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\nPing / Pongs: %s\n", string(content), pCount)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	http.HandleFunc("/status", getStatus)
	http.ListenAndServe(":"+port, nil)
}
