package main

import (
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
)

type Token struct {
	purpose    TokenPurpose
	content    string
	lineNumber int
}

type TokenPurpose = string

const (
	TokenPurposeIdentifier TokenPurpose = "identifier"
	// scalar type, https://protobuf.dev/programming-guides/proto3/#scalar
	TokenPurposeType       TokenPurpose = "type"
	TokenPurposeSymbol     TokenPurpose = "symbol"
	TokenPurposeWhitespace TokenPurpose = "whitespace"
	TokenPurposeString     TokenPurpose = "string"
	TokenPurposeInteger    TokenPurpose = "integer"
	TokenPurposeComment    TokenPurpose = "comment"
	TokenPurposeReserved   TokenPurpose = "reserved"
	TokenPurposeUnknown    TokenPurpose = "unknown"
	TokenPurposeAny        TokenPurpose = "any"
)

var TOKEN_TYPES = mapset.NewSet("map", "double", "float", "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes")
var TOKEN_SYMBOLS = mapset.NewSet('<', '>', ';', '}', '{', '=', '[', ']', ',')
var TOKEN_RESERVED = mapset.NewSet("enum", "option", "optional", "package", "syntax", "message", "repeated", "public")

// does this token match another token?
// we don't use equals here b/c line numbers can differ
func (t Token) matches(token Token) bool {
	// comments just need to match types, they shouldn't impact anything
	// I don't think protobuf has comment directives
	if t.purpose == TokenPurposeComment && token.purpose == TokenPurposeComment {
		return true
	}
	return t.purpose == token.purpose && t.content == token.content
}

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
		case '<', '>', ';', '}', '{', '=', '[', ']', ',':
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
				// skip double slash
				i += 2
				for hasNext() && currChar() != '\n' {
					word += string(currChar())
					i++
				}
				tokens = append(tokens, Token{TokenPurposeComment, word, lineNumber})
				word = ""
				lineNumber++
			} else if peekChar() == '*' {
				// handle multiline comments
				addExistingWord()
				// skip over comment starter
				i += 2
				internalLineNumber := 0
				for hasNext() && !(currChar() == '*' && peekChar() == '/') {
					if currChar() == '\n' {
						// line number should reflect start of comment
						internalLineNumber++
					}
					word += string(currChar())
					i++
				}
				i++
				tokens = append(tokens, Token{TokenPurposeComment, formatMultilineComment(word), lineNumber})
				word = ""
				lineNumber += internalLineNumber
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

// give a word, determine if it should be an identifier, number, type, or reserved word
func determineTokenPurpose(word string) TokenPurpose {
	// trimmed word is blank
	// means it's whitespace
	fmtWord := strings.TrimSpace(word)
	if len(fmtWord) == 0 {
		return TokenPurposeWhitespace
	}

	// is it a special char
	if len(word) > 0 && TOKEN_SYMBOLS.Contains(rune(word[0])) {
		return TokenPurposeSymbol
	}

	// matches a reserved word
	if TOKEN_TYPES.Contains(fmtWord) {
		return TokenPurposeType
	} else if TOKEN_RESERVED.Contains(fmtWord) {
		return TokenPurposeReserved
	}

	// is it an integer
	_, err := strconv.Atoi(fmtWord)
	if err == nil {
		return TokenPurposeInteger
	}

	return TokenPurposeIdentifier
}
