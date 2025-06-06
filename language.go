package main

import (
	"fmt"
	"strings"
)

type Metadata struct {
	packageName string
	syntax      string
}

type Message struct {
	name string
}

func t(purpose TokenPurpose, content string) Token {
	return Token{
		purpose:    purpose,
		content:    content,
		lineNumber: 0,
	}
}

func Language(content []byte) (Metadata, error) {
	scanner := newScanner(analyze(content))
	var metadata Metadata

	for scanner.hasNext() {
		switch {
		// skip over comments
		case scanner.matches(t(TokenPurposeComment, "//")):
			scanner.next()
		case scanner.matches(t(TokenPurposeComment, "/*")):
			scanner.next()

		// skip over option syntax
		case scanner.matches(t(TokenPurposeIdentifier, "option")):
			scanner.skipUntil(t(TokenPurposeSymbol, ";"))

		// package name
		case scanner.matches(t(TokenPurposeIdentifier, "package")):
			data, err := scanner.extract([]Token{
				t(TokenPurposeIdentifier, "package"),
				t(TokenPurposeIdentifier, "{{packageName}}"),
				t(TokenPurposeSymbol, ";"),
			})
			if err != nil {
				return metadata, err
			}
			metadata.packageName = data["packageName"].content

		// syntax
		case scanner.matches(t(TokenPurposeIdentifier, "syntax")):
			data, err := scanner.extract([]Token{
				t(TokenPurposeIdentifier, "syntax"),
				t(TokenPurposeSymbol, "="),
				t(TokenPurposeString, "{{syntax}}"),
				t(TokenPurposeSymbol, ";"),
			})
			if err != nil {
				return metadata, err
			}
			// remove starting and ending quotes
			metadata.syntax = strings.Trim(data["syntax"].content, "\"")
		case scanner.matches(t(TokenPurposeIdentifier, "message")):
			message_data, err := scanner.extract([]Token{
				t(TokenPurposeIdentifier, "message"),
				t(TokenPurposeIdentifier, "{{message}}"),
				t(TokenPurposeSymbol, "{"),
			})
			if err != nil {
				return metadata, err
			}

			for !scanner.curr().matches(t(TokenPurposeSymbol, "}")) && scanner.hasNext() {
				// reserved syntax - just skip until ;
				// if scanner.match("reserved") {
				// }

				// handle message field/attribute stuff
				data, err := scanner.extract([]Token{
					t(TokenPurposeIdentifier, "{{fieldType}}"),
					t(TokenPurposeIdentifier, "{{fieldName}}"),
					t(TokenPurposeSymbol, "="),
					t(TokenPurposeIdentifier, "{{fieldId}}"), // @todo numbers in lexer
					t(TokenPurposeSymbol, ";"),
				})
				if err != nil {
					return metadata, err
				}

				fmt.Println(message_data["message"], data["fieldType"], data["fieldName"], data["fieldId"])
			}
			scanner.i++ // skip }
		default:
			curr := scanner.curr()
			return metadata, fmt.Errorf("unsupported syntax at line %d:\n%s", curr.lineNumber, curr.content)
		}
	}

	return metadata, nil
}
