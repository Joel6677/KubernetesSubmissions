package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getCounter() int {
	content, err := os.ReadFile("/shared/pings.txt")
	if err != nil {
		return 0
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0
	}
	return count
}

func pingpong(w http.ResponseWriter, r *http.Request) {
	count := getCounter() + 1
	os.WriteFile("/shared/pings.txt", []byte(strconv.Itoa(count)), 0644)
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
