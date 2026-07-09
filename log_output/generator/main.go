package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomString := uuid.New().String()
	filePath := "/shared/output.txt"

	for {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		line := fmt.Sprintf("%s: %s\n", timestamp, randomString)

		os.WriteFile(filePath, []byte(line), 0644)
		fmt.Print(line)
		time.Sleep(5 * time.Second)
	}
}
