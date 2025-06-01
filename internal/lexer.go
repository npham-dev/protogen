package internal

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Lexer struct {
	content []byte
	i       int
}

func newLexer(content []byte) Lexer {
	return Lexer{content: content, i: 0}
}

// @todo store line number in struct for better error messages?
type Token struct {
	purpose TokenPurpose
	content string
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

func (l *Lexer) analyze() []Token {
	// split content into tokens
	var tokens []Token
	var word string
	i := 0

	// utility methods for common checks
	hasNext := func() bool {
		return i < len(l.content)
	}

	currChar := func() byte {
		return l.content[i]
	}

	nextChar := func() byte {
		return l.content[i+1]
	}

	for hasNext() {
		switch currChar() {
		case '<':
		case '>':
		case ';':
		case '}':
		case '{':
		case '=':
			tokens = append(tokens, Token{TokenPurposeIdentifier, word})
			tokens = append(tokens, Token{TokenPurposeSymbol, string(currChar())})
			word = ""
		case ' ', '\n':
			tokens = append(tokens, Token{TokenPurposeIdentifier, word})
			tokens = append(tokens, Token{TokenPurposeWhitespace, string(currChar())})
			word = ""
		default:
			if currChar() == '"' {
				// handle strings
				// skip over current quote
				i++
				for hasNext() && currChar() != '"' {
					word += string(currChar())
					i++
				}
				tokens = append(tokens, Token{TokenPurposeString, word})
				word = ""
			} else if currChar() == '/' && nextChar() == '/' {
				// handle slash comments
				i += 2 // skip over //
				for hasNext() && currChar() != '\n' {
					word += string(currChar())
					i++
				}
				tokens = append(tokens, Token{TokenPurposeComment, word})
				word = ""
			} else if currChar() == '/' && nextChar() == '*' {
				// handle multiline comments
				i += 2 // skip over /*
				for hasNext() && !(currChar() == '*' && nextChar() == '/') {
					word += string(currChar())
					i++
				}
				i += 2 // skip over */
				tokens = append(tokens, Token{TokenPurposeComment, formatMultilineComment(word)})
				word = ""
			} else {
				word += string(currChar())
			}
		}
		i++
	}
	tokens = append(tokens, Token{TokenPurposeIdentifier, word})

	fmt.Println(tokens)

	// we don't skip comments here b/c we want them to appear in the generated stuff
	// skip whitespace
	tokens = lo.Filter(
		lo.Map(tokens, func(token Token, _ int) Token {
			return Token{token.purpose, strings.TrimSpace(token.content)}
		}),
		func(token Token, _ int) bool {
			return !(len(token.content) == 0 || token.purpose == TokenPurposeWhitespace)
		})

	fmt.Println("\n", tokens)

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
	return TokenPurposeIdentifier
}
