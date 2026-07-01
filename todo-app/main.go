package main

import (
	"fmt"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
			<html>
				<body>
					<h1>Todo App</h1>
					<p>Moi!</p>
				</body>
			</html>
		`)
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
