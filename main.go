package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "./example.proto"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(analyze(content))
	fmt.Println(Language(content))
}
