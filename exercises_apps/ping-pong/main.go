package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

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
}

func ensureSchema() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS counter (id INT PRIMARY KEY, count INT NOT NULL)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO counter (id, count) VALUES (1, 0) ON CONFLICT (id) DO NOTHING`)
	return err
}

func incrementAndGet() (int, error) {
	var count int
	err := db.QueryRow(`UPDATE counter SET count = count + 1 WHERE id = 1 RETURNING count`).Scan(&count)
	return count, err
}

func getCount() (int, error) {
	var count int
	err := db.QueryRow(`SELECT count FROM counter WHERE id = 1`).Scan(&count)
	return count, err
}

func pingpong(w http.ResponseWriter, r *http.Request) {
	count, err := incrementAndGet()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Error incrementing counter:", err)
		return
	}
	fmt.Fprintf(w, "pong %d", count)
}

func pingCounter(w http.ResponseWriter, r *http.Request) {
	count, err := getCount()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("Error fetching counter:", err)
		return
	}
	fmt.Fprintf(w, "%d", count)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, "db not ready", http.StatusServiceUnavailable)
		return
	}
	if err := ensureSchema(); err != nil {
		http.Error(w, "db schema not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	port := getEnv("PORT")

	initDB()

	fmt.Printf("Server started in port %s\n", port)

	http.HandleFunc("/", pingpong)
	http.HandleFunc("/pings", pingCounter)

	http.HandleFunc("/healthz", healthHandler)

	http.ListenAndServe(":"+port, nil)
}
