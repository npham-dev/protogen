package main

import (
	"fmt"
	"log"
	"os"

	"github.com/natmfat/protogen/internal"
)

func main() {
	filePath := "./example.proto"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	metadata, err := internal.Language(content)
	fmt.Println(metadata, err)
	// language := newLanguage(tokens)
	// i := 0
	// for i := 0; i < len(content); i++ {
	// 	char :=

	// 	// fmt.Printf("%c\n", content[i])
	// }
}
