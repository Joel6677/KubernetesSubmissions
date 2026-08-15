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
	"time"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var formError string

	if r.Method == http.MethodPost {
		r.ParseForm()
		text := r.FormValue("text")

		if text == "" {
			formError = "Todo cant be empty"
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
					<li>Learn kubernetes basics</li>
					<li>Deploy application to cluster</li>
					<li>Configure persistent volumes</li>

	`)
	for _, t := range todos {
		fmt.Fprintf(w, "<li>%s</li>\n", t.Text)
	}
	fmt.Fprint(w, `
				</ul>
			</body>
		</html>
	`)
}

func main() {
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

	http.ListenAndServe(":"+port, nil)
}
