package internal

import (
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
type Token = string

func (l *Lexer) analyze() []Token {
	// split content into tokens
	var tokens []Token
	var word string
	for i := range l.content {
		curr_char := l.content[i]
		switch curr_char {
		case ';':
			tokens = append(tokens, word)
			tokens = append(tokens, ";")
			word = ""
		case ' ', '\n':
			tokens = append(tokens, word)
			word = ""
		default:
			word += string(curr_char)
		}
	}
	tokens = append(tokens, word)

	// skip whitespace
	return lo.Filter(lo.Map(tokens, func(token Token, _ int) string { return strings.TrimSpace(token) }), func(token Token, _ int) bool { return len(token) != 0 })
}
