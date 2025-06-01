package internal

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Scanner struct {
	tokens []Token
	i      int
}

func newScanner(tokens []Token) Scanner {
	return Scanner{tokens: tokens, i: 0}
}

func (s *Scanner) curr() Token {
	return s.tokens[s.i]
}

func (s *Scanner) next() Token {
	s.i += 1
	return s.tokens[s.i]
}

func (s *Scanner) hasNext() bool {
	return s.i < len(s.tokens)
}

func (s *Scanner) expect(token Token) error {
	if s.curr() != token {
		return fmt.Errorf("expected %s but got %s", token, s.curr())
	}
	return nil
}

func (s *Scanner) match(token Token) bool {
	return s.curr() == token
}

func (s *Scanner) skipUntil(token Token) {
	for s.hasNext() && s.curr() != token {
		s.next()
	}
	s.next()
}

func (s *Scanner) extract(pattern []Token) (map[string]Token, error) {
	data := make(map[string]Token)
	for _, token := range pattern {
		if strings.HasPrefix(token.content, "<") && strings.HasSuffix(token.content, ">") {
			key := strings.Trim(token.content, "<>")
			data[key] = s.curr()
		} else if s.curr() != token {
			return data, fmt.Errorf("failed to match pattern: [%s]", strings.Join(lo.Map(pattern, func(token Token, _ int) string { return token.content }), " "))
		}
		s.next()
	}

	// note that the current token is right after the matched pattern
	// ex) scanner.extract([]string{"pattern", ";"})
	//     the current token is NOT ";", it's whatever comes after it in the content
	return data, nil
}
