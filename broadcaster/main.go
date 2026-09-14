package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

type TodoEvent struct {
	Event string `json:"event"`
	Todo  struct {
		ID   int    `json:"id"`
		Text string `json:"text"`
		Done bool   `json:"done"`
	} `json:"todo"`
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment value %s is not set", key)
	}
	return v
}

func sendToWebhook(webhookURL string, message string) error {
	body, _ := json.Marshal(map[string]string{
		"user":    "bot",
		"message": message,
	})
	res, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func main() {
	natsURL := getEnv("NATS_URL")
	webhookURL := getEnv("WEBHOOK_URL")

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	_, err = nc.QueueSubscribe("todos.updates", "broadcaster-group", func(msg *nats.Msg) {
		var event TodoEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("failed to unmarshal event: %v", err)
			return
		}

		var message string
		if event.Event == "created" {
			message = "A todo was created: " + event.Todo.Text
		} else {
			status := "not done"
			if event.Todo.Done {
				status = "done"
			}
			message = "A todo was updated: " + event.Todo.Text + " (" + status + ")"
		}

		if err := sendToWebhook(webhookURL, message); err != nil {
			log.Printf("failed to send to webhook: %v", err)
			return
		}
		log.Printf("forwarded event to webhook: %s", message)
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	log.Println("broadcaster started, listening for todo updates")
	select {}
}
