package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func updateTodoDone(id int, done bool) error {
	body, _ := json.Marshal(map[string]bool{"done": done})
	url := fmt.Sprintf("%s/%d", getTodoBackendURL(), id)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("todo-backend returned status %d", res.StatusCode)
	}
	return nil
}

func toggleDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method works", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	done := r.FormValue("done") == "true"

	if err := updateTodoDone(id, done); err != nil {
		fmt.Println("Error updating todo:", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var isHealthy atomic.Bool

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment value %s is not set", key)
	}
	return v
}

func getTodoBackendURL() string {
	return getEnv("TODO_BACKEND_URL")
}

func fetchTodos() ([]Todo, error) {
	res, err := http.Get(getTodoBackendURL())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var todos []Todo
	if err := json.NewDecoder(res.Body).Decode(&todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func createTodo(text string) error {
	body, _ := json.Marshal(map[string]string{"text": text})
	res, err := http.Post(getTodoBackendURL(), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	io.ReadAll(res.Body)
	return nil
}

var (
	imgPath = getEnv("IMG_PATH")
	imgURL  = getEnv("IMG_URL")
)

var client = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, "tcp4", addr)
		},
	},
}

func getImageTime() time.Duration {
	fileInfo, err := os.Stat(imgPath)
	if err != nil {
		return time.Hour * 24
	}
	return time.Since(fileInfo.ModTime())
}

func fetchAndSaveImage() error {
	req, err := http.NewRequest(http.MethodGet, imgURL, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(imgPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if getImageTime() > 10*time.Minute {
		fetchAndSaveImage()
	}

	data, err := os.ReadFile(imgPath)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !isHealthy.Load() {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
		return
	}

	res, err := http.Get(getTodoBackendURL())
	if err != nil || res.StatusCode >= 500 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "reason": "backend unreachable"})
		return
	}
	res.Body.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func breakHandler(w http.ResponseWriter, r *http.Request) {
	isHealthy.Store(false)
	fmt.Println("todo-app marked unhealthy via /break")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var formError string

	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")

		if text == "" {
			formError = "Todo can't be empty"
		} else if len([]rune(text)) > 140 {
			formError = "Todo must be max 140 characters"
		} else {
			if err := createTodo(text); err != nil {
				fmt.Println("Error creating todo:", err)
				formError = "Failed to save todo, try again"
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
	}

	todos, err := fetchTodos()
	if err != nil {
		fmt.Println("Error fetching todos:", err)
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<html>
			<head>
				<style>
					body { font-family: sans-serif; max-width: 600px; margin: 40px auto; }
					input[type=text] {
						width: 100%%;
						padding: 12px;
						font-size: 16px;
						box-sizing: border-box;
					}
					button {
						padding: 10px 20px;
						font-size: 16px;
						margin-top: 8px;
					}
					.error { color: red; }
					.break-btn { background: #c0392b; color: white; border: none; margin-top: 30px; }
				</style>
			</head>
			<body>
				<h1>Todo App</h1>
				<img src="/image" width="600" /><br/><br/>
	`)
	if formError != "" {
		fmt.Fprintf(w, `<p class="error">%s</p>`, formError)
	}
	fmt.Fprintf(w, `
				<form method="POST" action="/">
					<input type="text" name="text" placeholder="Enter a new todo (max 140 characters)" maxlength="140" required />
					<button type="submit">Add</button>
				</form>
				<ul>
				
	`)
	for _, t := range todos {
		checked := ""
		if t.Done {
			checked = "checked"
		}
		doneValue := "true"
		if t.Done {
			doneValue = "false"
		}
		fmt.Fprintf(w, `
		<li>
			<form method="POST" action="/toggle" style="display:inline;">
				<input type="hidden" name="id" value="%d" />
				<input type="hidden" name="done" value="%s" />
				<input type="checkbox" %s onchange="this.form.submit()" />
			</form>
			<span style="%s">%s</span>
		</li>
	`, t.ID, doneValue, checked,
			func() string {
				if t.Done {
					return "text-decoration: line-through; color: gray;"
				}
				return ""
			}(),
			t.Text)
	}
	fmt.Fprint(w, `
				</ul>
				<form method="POST" action="/break">
					<button type="submit" class="break-btn">Break the app</button>
				</form>
			</body>
		</html>
	`)
}

func main() {
	isHealthy.Store(true)
	port := getEnv("PORT")

	fmt.Printf("Server started in port %s\n", port)

	if getImageTime() > 10*time.Minute {
		fmt.Println("Fetching new image...")
		err := fetchAndSaveImage()
		if err != nil {
			fmt.Println("Error fetching image:", err)
		} else {
			fmt.Println("Image saved successfully")
		}
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/break", breakHandler)
	http.HandleFunc("/toggle", toggleDoneHandler)

	http.ListenAndServe(":"+port, nil)
}
