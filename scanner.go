package main

import (
	"fmt"
	"strings"
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
	if s.i < len(s.tokens) {
		return s.tokens[s.i]
	}
	return Token{}
	// @todo maybe don't return zero value? probably should let user know or what
}

func (s *Scanner) hasNext() bool {
	return s.i < len(s.tokens)
}

// func (s *Scanner) expect(token Token) error {
// 	if s.curr() != token {
// 		return fmt.Errorf("expected %s but got %s", token.content, s.curr().content)
// 	}
// 	return nil
// }

func (s *Scanner) matches(token Token) bool {
	return s.curr().matches(token)
}

func (s *Scanner) matchesPurpose(purposes []TokenPurpose) bool {
	for _, purpose := range purposes {
		if purpose == s.curr().purpose {
			return true
		}
	}
	return false
}

func (s *Scanner) skipUntil(token Token) {
	for s.hasNext() && !s.curr().matches(token) {
		s.next()
	}
	s.next()
}

func (s *Scanner) extract(pattern []Token) (map[string]Token, error) {
	data := make(map[string]Token)
	for _, token := range pattern {
		// get current token (skip over comments)
		for s.hasNext() && s.curr().purpose == TokenPurposeComment {
			s.next()
		}
		curr := s.curr()

		if strings.HasPrefix(token.content, "{{") && strings.HasSuffix(token.content, "}}") {
			if curr.purpose == token.purpose || token.purpose == TokenPurposeAny {
				key := strings.Trim(token.content, "{}")
				data[key] = curr
			} else {
				return data, syntaxError(
					curr.lineNumber,
					fmt.Sprintf(
						"expected '%s' but found '%s'",
						token.content,
						curr.content))
			}
		} else if !curr.matches(token) {
			return data, syntaxError(
				curr.lineNumber,
				fmt.Sprintf("expected '%s' but found '%s'", token.content,
					curr.content,
				))
		}
		s.next()
	}

	// note that the current token is right after the matched pattern
	// ex) scanner.extract([]string{"pattern", ";"})
	//     the current token is NOT ";", it's whatever comes after it in the content
	return data, nil
}

func syntaxError(lineNumber int, message string) error {
	return fmt.Errorf(
		"[protogen] error at line %d:\n%s",
		lineNumber,
		message,
	)
}
