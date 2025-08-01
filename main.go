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

	tokens := analyze(content)
	fmt.Println(parse(tokens))
}
