package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	imgPath = "/tmp/img.jpg"
	imgURL  = "https://picsum.photos/1200"
)

func getImageTime() time.Duration {
	fileInfo, err := os.Stat(imgPath)
	if err != nil {
		return time.Hour * 24
	}
	return time.Since(fileInfo.ModTime())
}

func fetchAndSaveImage() error {
	res, err := http.Get(imgURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(imgPath)
	if err != nil {
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
				    <meta http-equiv="refresh" content="600">
				</head>
				<body>
					<h1>Todo App</h1>
					<img src="/image" width="600" />
					<script>
						console.log("image loaded")
					</script>
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
		fetchAndSaveImage()
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/image", imageHandler)

	http.ListenAndServe(":"+port, nil)
}
