package internal

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

/*
scanner.extract(["package", "<client>", ";"])
scanner.extract(["syntax", "=", "<message>", ";"])
scanner.extract(["message", "<Message>", "{"])
*/

func Language(content []byte) (Metadata, error) {
	lexer := newLexer(content)
	scanner := newScanner(lexer.analyze())
	var metadata Metadata

	for scanner.hasNext() {
		if scanner.match("package") {
			data, err := scanner.extract([]string{"package", "<packageName>", ";"})
			if err != nil {
				return metadata, err
			}
			metadata.packageName = data["packageName"]
		} else if scanner.match("syntax") {
			data, err := scanner.extract([]string{"syntax", "=", "<syntax>", ";"})
			if err != nil {
				return metadata, err
			}
			// remove starting and ending quotes
			metadata.syntax = strings.Trim(data["syntax"], "\"")
		} else if scanner.match("message") {
			message_data, err := scanner.extract([]string{"message", "<message>", "{"})
			if err != nil {
				return metadata, err
			}

			for scanner.curr() != "}" {
				data, err := scanner.extract([]string{"<fieldType>", "<fieldName>", "=", "<fieldId>", ";"})
				if err != nil {
					return metadata, err
				}

				fmt.Println(message_data["message"], data["fieldType"], data["fieldName"], data["fieldId"])
			}
			scanner.next() // skip }
		}
	}

	return metadata, nil
}
