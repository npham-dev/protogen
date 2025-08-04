package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
)

func main() {
	filePath := "./example.proto"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	tokens := analyze(content)
	data, err := parse(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%# v\n", pretty.Formatter(data))
	output := generate(data)
	err = os.WriteFile("example.ts", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
