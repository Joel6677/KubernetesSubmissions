package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	mu     sync.Mutex
	todos  = []Todo{}
	nextID = 1
)

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment value %s is not set", key)
	}
	return v
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todos)

	case http.MethodPost:
		var body struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(body.Text) > 140 {
			http.Error(w, "Todo text is over 140 characters", http.StatusBadRequest)
			return
		}

		mu.Lock()
		newTodo := Todo{ID: nextID, Text: body.Text}
		nextID++
		todos = append(todos, newTodo)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	default:
		http.Error(w, "Only get and post methods work", http.StatusMethodNotAllowed)
	}
}

func main() {
	port := os.Getenv("PORT")

	fmt.Printf("todo-backend started on port %s\n", port)
	http.HandleFunc("/todos", todosHandler)
	http.ListenAndServe(":"+port, nil)
}
