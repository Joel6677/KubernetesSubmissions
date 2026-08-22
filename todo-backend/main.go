package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var db *sql.DB

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment value %s is not set", key)
	}
	return v
}

func initDB() {
	connStr := getEnv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
}

func getTodos() ([]Todo, error) {
	rows, err := db.Query(`SELECT id, text FROM todos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func insertTodo(text string) (Todo, error) {
	var t Todo
	t.Text = text
	err := db.QueryRow(
		`INSERT INTO todos (text) VALUES ($1) RETURNING id`,
		text,
	).Scan(&t.ID)
	return t, err
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		next(wrapped, r)

		log.Printf(
			"method=%s path=%s status=%d duration=%s remote=%s",
			r.Method, r.URL.Path, wrapped.statusCode, time.Since(start), r.RemoteAddr,
		)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos, err := getTodos()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Error fetching todos:", err)
			return
		}
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
			log.Printf("todo is too long (%d chars): %q", len(body.Text), body.Text)
			http.Error(w, "Todo text is over 140 characters", http.StatusBadRequest)
			return
		}

		newTodo, err := insertTodo(body.Text)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			fmt.Println("Error inserting todo:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	default:
		http.Error(w, "Only get and post methods work", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	port := getEnv("PORT")
	initDB()
	fmt.Printf("todo-backend started on port %s\n", port)
	http.HandleFunc("/todos", loggingMiddleware(todosHandler))
	http.ListenAndServe(":"+port, nil)
}
