package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	imgPath = "/shared/img.jpg"
	imgURL  = "https://picsum.photos/1200"
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
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
			<html>
				<head>
					<script>
						function addTodo() {
						    const input = document.getElementById('todo');
						    const todo = input.value;
						    if (todo.length === 0) return;
						    if (todo.length > 140) {
							alert('Todo is over 140 characters!');
							return;
						    }
						    const li = document.createElement('li');
						    li.textContent = todo;
						    document.getElementById('todo-list').appendChild(li);
						    input.value = '';
						}
					</script>
				</head>
				<body>
					<h1>Todo App</h1>
					<img src="/image" width="600" />
					<br><br>
					<form>
						<input 
							type="text" 
							id="todo" 
							maxlength="140" 
							placeholder="Enter a new todo (max 140 characters)"
							size="50"
						/>
						<button type="button" onclick="addTodo()">Send</button>
					</form>
					<ul id="todo-list">
					<li>Learn kubernetes basics</li>
					<li>Deploy application to cluster</li>
					<li>Configure persistent volumes</li>
					</ul>
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
