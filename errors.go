package main

import "fmt"

func syntaxError(lineNumber int, message string) error {
	return fmt.Errorf(
		"[protogen] error at line %d:\n%s",
		lineNumber,
		message,
	)
}

func generateError(message string) error {
	return fmt.Errorf("[protogen] failed to generate:\n%s", message)
}
