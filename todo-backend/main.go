package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:done`
}

var (
	db        *sql.DB
	isHealthy atomic.Bool
)

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
			text TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	_, err = db.Exec(`ALTER TABLE todos ADD COLUMN IF NOT EXISTS done BOOLEAN NOT NULL DEFAULT FALSE`)
	if err != nil {
		log.Fatalf("failed to migrate table: %v", err)
	}
}

func getTodos() ([]Todo, error) {
	rows, err := db.Query(`SELECT id, text, done FROM todos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Done); err != nil {
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
		`INSERT INTO todos (text) VALUES ($1) RETURNING id, done`,
		text,
	).Scan(&t.ID, &t.Done)
	return t, err
}

func updateTodoDone(id int, done bool) (Todo, error) {
	var t Todo
	err := db.QueryRow(
		`UPDATE todos SET done = $1 WHERE id = $2 RETURNING id, text, done`,
		done, id,
	).Scan(&t.ID, &t.Text, &t.Done)
	return t, err
}

func todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only put method allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo id", http.StatusBadRequest)
		return
	}

	var body struct {
		Done bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updated, err := updateTodoDone(id, body.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("error updating todo: %v", err)
		return
	}

	log.Printf("updated todo id=%d done=%v", updated.ID, updated.Done)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
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

		log.Printf("created todo id=%d text=%q", newTodo.ID, newTodo.Text)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)

	default:
		http.Error(w, "Only get and post methods work", http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isHealthy.Load() {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
		return
	}

	if err := db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "reason": "db unreachable"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func breakHandler(w http.ResponseWriter, r *http.Request) {
	isHealthy.Store(false)
	log.Println("todo-backend unhealthy via /break")
	w.Write([]byte("todo-backend unhealthy"))
}

func main() {
	log.SetOutput(os.Stdout)
	isHealthy.Store(true)
	port := getEnv("PORT")
	initDB()
	fmt.Printf("todo-backend started on port %s\n", port)
	http.HandleFunc("/todos", loggingMiddleware(todosHandler))
	http.HandleFunc("/todos/", loggingMiddleware(todoByIDHandler))
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/break", breakHandler)
	http.ListenAndServe(":"+port, nil)
}
