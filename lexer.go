package main

import (
	"strings"

	"github.com/samber/lo"
)

// @todo store line number in struct for better error messages?
type Token struct {
	purpose    TokenPurpose
	content    string
	lineNumber int
}

type TokenPurpose = string

const (
	TokenPurposeIdentifier TokenPurpose = "identifier"
	TokenPurposeSymbol     TokenPurpose = "symbol"
	TokenPurposeWhitespace TokenPurpose = "whitespace"
	TokenPurposeString     TokenPurpose = "string"
	TokenPurposeComment    TokenPurpose = "comment"
	TokenPurposeUnknown    TokenPurpose = "unknown"
)

func analyze(content []byte) []Token {
	// split content into tokens
	var tokens []Token
	var word string
	i := 0
	lineNumber := 1

	// utility methods for common checks
	hasNext := func() bool { return i < len(content) }
	currChar := func() byte { return content[i] }
	peekChar := func() byte { return content[i+1] }
	addExistingWord := func() {
		if len(word) > 0 {
			tokens = append(tokens, Token{determineTokenPurpose(word), word, lineNumber})
		}
		word = ""
	}

	for hasNext() {
		switch currChar() {
		case '<', '>', ';', '}', '{', '=':
			addExistingWord()
			tokens = append(tokens, Token{TokenPurposeSymbol, string(currChar()), lineNumber})
			word = ""
		case ' ', '\n':
			addExistingWord()
			tokens = append(tokens, Token{TokenPurposeWhitespace, string(currChar()), lineNumber})
			word = ""
			if currChar() == '\n' {
				lineNumber++
			}
		case '"':
			// @todo escaping comments/quotes?
			i++
			for hasNext() && currChar() != '"' {
				word += string(currChar())
				i++
			}
			tokens = append(tokens, Token{TokenPurposeString, word, lineNumber})
			word = ""
		case '/':
			if peekChar() == '/' {
				addExistingWord()
				// handle slash comments
				i += 2 // skip over //
				for hasNext() && currChar() != '\n' {
					word += string(currChar())
					i++
				}
				tokens = append(tokens, Token{TokenPurposeComment, word, lineNumber})
				word = ""
				lineNumber++
			} else if peekChar() == '*' {
				addExistingWord()
				// handle multiline comments
				i += 2 // skip over /*
				for hasNext() && !(currChar() == '*' && peekChar() == '/') {
					if currChar() == '\n' {
						lineNumber++
					}
					word += string(currChar())
					i++
				}
				i++ // skip over */
				tokens = append(tokens, Token{TokenPurposeComment, formatMultilineComment(word), lineNumber})
				word = ""
			} else {
				word += string(currChar())
			}
		default:
			word += string(currChar())
		}
		i++
	}
	addExistingWord()

	// we don't skip comments here b/c we want them to appear in the generated stuff
	// skip whitespace
	tokens = lo.Filter(
		lo.Map(tokens, func(token Token, _ int) Token {
			return Token{token.purpose, strings.TrimSpace(token.content), token.lineNumber}
		}),
		func(token Token, _ int) bool {
			return !(token.purpose == TokenPurposeWhitespace)
		})
	return tokens
}

// format a multiline comment by
// 1. trimming each line
// 2. removing any asterisks (*) at the start of each line
//
// inputs should not have an opening comment symbol (/*) or closing comment symbol (*/)
func formatMultilineComment(comment string) string {
	return strings.Join(lo.Map(strings.Split(comment, "\n"), func(line string, _ int) string {
		return strings.TrimPrefix(strings.TrimSpace(line), "*")
	}), "\n")
}

// give a word, determine if it should be an identifier or reserved word
func determineTokenPurpose(word string) TokenPurpose {
	if len(word) > 0 {
		switch word[0] {
		case '<', '>', ';', '}', '{', '=':
			return TokenPurposeSymbol
		case ' ', '\n':
			return TokenPurposeWhitespace
		}
	}
	if len(strings.TrimSpace(word)) == 0 {
		return TokenPurposeWhitespace
	}
	return TokenPurposeIdentifier
}
