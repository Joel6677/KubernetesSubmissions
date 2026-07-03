package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var randomString = uuid.New().String()

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), randomString)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		for {
			fmt.Printf("%s: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), randomString)
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/status", getStatus)

	http.ListenAndServe(":"+port, nil)
}
