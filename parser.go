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
	name string
	id   string
}

type SyntaxMessage struct {
	name     string
	fields   []SyntaxMessageField
	messages []SyntaxMessage
}

type SyntaxMessageField struct {
	name string
	id   string
	kind string

	repeated bool
	optional bool

	// optional values for maps
	mapKey   *string
	mapValue *string
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
			// ex) option allow_alias = true;
			//   	 EAA_STARTED = 1;
			//		 EAA_RUNNING = 1;

			// parse enum body
			for !scanner.matches(t(TokenPurposeSymbol, "}")) && scanner.hasNext() {
				// handle message field/attribute stuff
				data, err := scanner.extract([]Token{
					t(TokenPurposeIdentifier, "{{name}}"),
					t(TokenPurposeSymbol, "="),
					t(TokenPurposeInteger, "{{id}}"),
					t(TokenPurposeSymbol, ";"),
				})
				if err != nil {
					return document, err
				}

				syntaxEnum.fields = append(syntaxEnum.fields, SyntaxEnumField{
					name: data["name"].content,
					id:   data["id"].content,
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
			// handle repeated & optional flags
			repeated := false
			optional := false
			if scanner.matches(t(TokenPurposeReserved, "repeated")) {
				repeated = true
				scanner.next()
			} else if scanner.matches(t(TokenPurposeReserved, "optional")) {
				optional = true
				scanner.next()
			}

			syntaxMessageField := SyntaxMessageField{
				repeated: repeated,
				optional: optional,
			}

			switch {
			// handle maps
			case scanner.matches(t(TokenPurposeType, "map")):
				data, err := scanner.extract([]Token{
					t(TokenPurposeType, "map"),
					t(TokenPurposeSymbol, "<"),
					t(TokenPurposeType, "{{mapKey}}"),
					t(TokenPurposeSymbol, ","),
					t(TokenPurposeAny, "{{mapValue}}"), // @todo map value can be either a type or token purpose identifier
					t(TokenPurposeSymbol, ">"),
					t(TokenPurposeIdentifier, "{{name}}"),
					t(TokenPurposeSymbol, "="),
					t(TokenPurposeInteger, "{{id}}"),
					t(TokenPurposeSymbol, ";"),
				})
				if err != nil {
					return syntaxMessage, err
				}

				syntaxMessageField.kind = "map"
				syntaxMessageField.name = data["name"].content
				syntaxMessageField.id = data["id"].content

				// validate mapKey - mapKeys can be any scalar type but floats & bytes
				// https://protobuf.dev/programming-guides/proto3/#maps
				mapKey := data["mapKey"].content
				if mapKey == "float" || mapKey == "bytes" {
					return syntaxMessage, syntaxError(data["mapKey"].lineNumber, fmt.Sprintf("map keys cannot be of type '%s'", mapKey))
				}

				// validate mapValue - mapValues can be anything except other maps
				// extract would fail if we passed another map anyways
				// we just want to check if it's a type or identifier and not something weird
				mapValue := data["mapValue"].content
				if !(data["mapValue"].purpose == TokenPurposeType || data["mapValue"].purpose == TokenPurposeIdentifier) {
					return syntaxMessage, syntaxError(data["mapKey"].lineNumber, fmt.Sprintf("map values cannot be of type '%s'", mapValue))
				}

				syntaxMessageField.mapKey = &mapKey
				syntaxMessageField.mapValue = &mapValue

			// handle built-in types & declared types like enums
			// ex) bool field = 1;
			// ex) Corpus corpus = 1;
			case scanner.matchesPurpose([]TokenPurpose{TokenPurposeType, TokenPurposeIdentifier}):
				data, err := scanner.extract([]Token{
					t(TokenPurposeAny, "{{type}}"),
					t(TokenPurposeIdentifier, "{{name}}"),
					t(TokenPurposeSymbol, "="),
					t(TokenPurposeInteger, "{{id}}"),
					t(TokenPurposeSymbol, ";"),
				})
				if err != nil {
					return syntaxMessage, err
				}

				syntaxMessageField.kind = data["type"].content
				syntaxMessageField.name = data["name"].content
				syntaxMessageField.id = data["id"].content
			}

			syntaxMessage.fields = append(syntaxMessage.fields, syntaxMessageField)
		}
	}
	scanner.next() // skip }

	return syntaxMessage, nil
}
