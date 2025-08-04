package main

import (
	"fmt"
	"strings"
)

// I feel like this isn't really an ast? the structure is pretty flat
// idk but I'm prefixing all outputs with "Syntax" like "Syntax Tree"
type SyntaxDocument struct {
	packageName string
	syntax      string

	enums    []SyntaxEnum
	messages []SyntaxMessage
}

type SyntaxEnum struct {
	name   string
	fields []SyntaxEnumField
}

type SyntaxEnumField struct {
	name  string
	value string
}

type SyntaxMessage struct {
	name     string
	fields   []SyntaxMessageField
	messages []SyntaxMessage
}

type SyntaxMessageField struct {
	name      string
	value     string
	fieldType string
}

func t(purpose TokenPurpose, content string) Token {
	return Token{
		purpose:    purpose,
		content:    content,
		lineNumber: 0,
	}
}

func parse(tokens []Token) (SyntaxDocument, error) {
	scanner := newScanner(tokens)
	var document SyntaxDocument

	for scanner.hasNext() {
		switch {
		// skip over comments
		// @todo add comments to output to build documentation
		case scanner.matches(t(TokenPurposeComment, "//")):
			scanner.next()
		case scanner.matches(t(TokenPurposeComment, "/*")):
			scanner.next()

		// skip over option syntax (file level directives)
		// https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/descriptor.proto
		case scanner.matches(t(TokenPurposeReserved, "option")):
			scanner.skipUntil(t(TokenPurposeSymbol, ";"))

		// package name
		// ex) package client;
		case scanner.matches(t(TokenPurposeReserved, "package")):
			data, err := scanner.extract([]Token{
				t(TokenPurposeReserved, "package"),
				t(TokenPurposeIdentifier, "{{packageName}}"),
				t(TokenPurposeSymbol, ";"),
			})
			if err != nil {
				return document, err
			}
			document.packageName = data["packageName"].content

		// syntax
		// ex) syntax = "proto3";
		case scanner.matches(t(TokenPurposeReserved, "syntax")):
			data, err := scanner.extract([]Token{
				t(TokenPurposeReserved, "syntax"),
				t(TokenPurposeSymbol, "="),
				t(TokenPurposeString, "{{syntax}}"),
				t(TokenPurposeSymbol, ";"),
			})
			if err != nil {
				return document, err
			}
			// remove starting and ending quotes
			document.syntax = strings.Trim(data["syntax"].content, "\"")

		// enum
		case scanner.matches(t(TokenPurposeReserved, "enum")):
			enumData, err := scanner.extract([]Token{
				t(TokenPurposeReserved, "enum"),
				t(TokenPurposeIdentifier, "{{name}}"),
				t(TokenPurposeSymbol, "{"),
			})
			if err != nil {
				return document, err
			}

			syntaxEnum := SyntaxEnum{name: enumData["name"].content}

			// @todo "Enum Value Aliases"
			// parse enum body
			for !scanner.matches(t(TokenPurposeSymbol, "}")) && scanner.hasNext() {
				// handle message field/attribute stuff
				data, err := scanner.extract([]Token{
					t(TokenPurposeIdentifier, "{{fieldName}}"),
					t(TokenPurposeSymbol, "="),
					t(TokenPurposeInteger, "{{fieldValue}}"),
					t(TokenPurposeSymbol, ";"),
				})
				if err != nil {
					return document, err
				}

				syntaxEnum.fields = append(syntaxEnum.fields, SyntaxEnumField{
					name:  data["fieldName"].content,
					value: data["fieldValue"].content,
				})
			}

			scanner.next() // skip }

			document.enums = append(document.enums, syntaxEnum)

		// message
		case scanner.matches(t(TokenPurposeReserved, "message")):
			syntaxMessage, err := parseMessage(&scanner)
			if err != nil {
				return document, err
			}
			document.messages = append(document.messages, syntaxMessage)

		default:
			curr := scanner.curr()
			return document, fmt.Errorf("unsupported syntax at line %d:\n%s", curr.lineNumber, curr.content)
		}
	}

	return document, nil
}

// dedicated method because we need to recursively parse nested messages
func parseMessage(scanner *Scanner) (SyntaxMessage, error) {
	syntaxMessage := SyntaxMessage{}
	messageData, err := scanner.extract([]Token{
		t(TokenPurposeReserved, "message"),
		t(TokenPurposeIdentifier, "{{name}}"),
		t(TokenPurposeSymbol, "{"),
	})
	if err != nil {
		return syntaxMessage, err
	}

	syntaxMessage.name = messageData["name"].content

	for !scanner.matches(t(TokenPurposeSymbol, "}")) && scanner.hasNext() {
		switch {
		// reserved syntax - just skip until ;
		// has no use for generation
		case scanner.matches(t(TokenPurposeIdentifier, "reserved")):
			scanner.skipUntil(t(TokenPurposeSymbol, ";"))

		// handle nested messages
		case scanner.matches(t(TokenPurposeReserved, "message")):
			childSyntaxMessage, parseMessageErr := parseMessage(scanner)
			if parseMessageErr != nil {
				return syntaxMessage, err
			}
			syntaxMessage.messages = append(syntaxMessage.messages, childSyntaxMessage)

		// handle field/attribute stuff
		default:
			data, err := scanner.extract([]Token{
				// @todo first might be token purpose identifier (enum or another message)
				// @todo repeated or option
				t(TokenPurposeType, "{{fieldType}}"),
				t(TokenPurposeIdentifier, "{{fieldName}}"),
				t(TokenPurposeSymbol, "="),
				t(TokenPurposeInteger, "{{fieldId}}"),
				t(TokenPurposeSymbol, ";"),
			})
			if err != nil {
				return syntaxMessage, err
			}

			syntaxMessage.fields = append(syntaxMessage.fields, SyntaxMessageField{
				name:      data["fieldName"].content,
				value:     data["fieldId"].content,
				fieldType: data["fieldType"].content,
			})
		}
	}
	scanner.next() // skip }

	return syntaxMessage, nil
}
