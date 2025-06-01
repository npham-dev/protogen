package internal

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

	scanner.hasNext()

	// for scanner.hasNext() {
	// 	switch {
	// 	// skip over comments
	// 	// case scanner.match(Token{purpose: TokenPurposeComment, content: "//"}):
	// 	// scanner.skipUntil(Token{purpose: TokenPurposeWhitespace, content: "\n"})
	// 	// case scanner.match(Token{purpose: TokenPurposeComment, content: "/*"}):
	// 	// scanner.skipUntil("*/")
	// 	// skip over option syntax
	// 	case scanner.match(Token{TokenPurposeIdentifier, "option"}):
	// 		scanner.skipUntil(Token{TokenPurposeSymbol, ";"})
	// 	case scanner.match(Token{TokenPurposeIdentifier, "package"}):
	// 		data, err := scanner.extract([]string{"package", "<packageName>", ";"})
	// 		if err != nil {
	// 			return metadata, err
	// 		}
	// 		metadata.packageName = data["packageName"].content
	// 	case scanner.match(Token{TokenPurposeIdentifier, "syntax"}):
	// 		data, err := scanner.extract([]string{"syntax", "=", "<syntax>", ";"})
	// 		if err != nil {
	// 			return metadata, err
	// 		}
	// 		// remove starting and ending quotes
	// 		metadata.syntax = strings.Trim(data["syntax"].content, "\"")
	// 	case scanner.match(Token{TokenPurposeIdentifier, "message"}):
	// 		message_data, err := scanner.extract([]string{"message", "<message>", "{"})
	// 		if err != nil {
	// 			return metadata, err
	// 		}

	// 		for scanner.curr() != "}" {
	// 			// reserved syntax - just skip
	// 			// if scanner.match("reserved") {
	// 			// }

	// 			// handle message field/attribute stuff
	// 			data, err := scanner.extract([]string{"<fieldType>", "<fieldName>", "=", "<fieldId>", ";"})
	// 			if err != nil {
	// 				return metadata, err
	// 			}

	// 			fmt.Println(message_data["message"], data["fieldType"], data["fieldName"], data["fieldId"])
	// 		}
	// 		scanner.i++ // skip }
	// 	default:
	// 		return metadata, errors.New("unsupported syntax")
	// 	}
	// }

	return metadata, nil
}
