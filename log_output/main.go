package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomString := uuid.New().String()
	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), randomString)
		time.Sleep(5 * time.Second)
	}
}
