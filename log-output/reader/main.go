package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getPings() string {
	pingPongURL := os.Getenv("PINGPONG_URL")
	if pingPongURL == "" {
		pingPongURL = "http://ping-pong-svc:2345/pings"
	}

	res, err := http.Get(pingPongURL)
	if err != nil {
		fmt.Println("Error fetching ping count:", err)
		return "0"
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "0"
	}
	return string(body)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("/shared/output.txt")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileContent, err := os.ReadFile("/config/information.txt")
	fileText := ""
	if err != nil {
		fileText = "could not read file"
	} else {
		fileText = strings.TrimSpace(string(fileContent))
	}

	message := os.Getenv("MESSAGE")

	pCount := getPings()

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "file content: %s\nenv variable: MESSAGE=%s\n%s\nPing / Pongs: %s\n",
		fileText, message, string(content), pCount)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started in port %s\n", port)
	http.HandleFunc("/", getStatus)
	http.ListenAndServe(":"+port, nil)
}
