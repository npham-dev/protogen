package main

import (
	"log"
	"os"
)

func main() {
	filePath := "./example.proto"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	analyze(content)
}
