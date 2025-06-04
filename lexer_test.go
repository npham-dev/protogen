package main

import (
	"reflect"
	"testing"
)

type Case struct {
	name     string
	input    string
	expected []Token
}

func testCases(t *testing.T, cases []Case) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := analyze([]byte(c.input))
			if !reflect.DeepEqual(result, c.expected) {
				t.Errorf("expected %v but got %v", c.expected, result)
			}
		})
	}
}

func TestComments(t *testing.T) {
	testCases(t, []Case{
		{
			"single line comment",
			`// hello world`,
			[]Token{
				{TokenPurposeComment, "hello world", 1},
			},
		},
		{
			"multiline comment",
			`/* hello world */`,
			[]Token{
				{TokenPurposeComment, "hello world", 1},
			},
		},
		{
			"single line comment with leading spaces",
			`   // spaced comment`,
			[]Token{
				{TokenPurposeComment, "spaced comment", 1},
			},
		},
		{
			"single line comment after code",
			`foo // trailing comment`,
			[]Token{
				{TokenPurposeIdentifier, "foo", 1},
				{TokenPurposeComment, "trailing comment", 1},
			},
		},
		{
			"multiline comment spanning lines",
			`/* hello
world */`,
			[]Token{
				{TokenPurposeComment, "hello\nworld", 1},
			},
		},
		{
			"empty single line comment",
			`//`,
			[]Token{
				{TokenPurposeComment, "", 1},
			},
		},
		{
			"empty multiline comment",
			`/**/`,
			[]Token{
				{TokenPurposeComment, "", 1},
			},
		},
		{
			"multiline comment with stars",
			`/* hello * world */`,
			[]Token{
				{TokenPurposeComment, "hello * world", 1},
			},
		},
		{
			"single line comment with symbols",
			`// !@#$%^&*()`,
			[]Token{
				{TokenPurposeComment, "!@#$%^&*()", 1},
			},
		},
	})
}
