package internal

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
	TokenPurposeSpecial    TokenPurpose = "special"
	TokenPurposeWhitespace TokenPurpose = "whitespace"
	TokenPurposeString     TokenPurpose = "string"
	TokenPurposeComment    TokenPurpose = "comment"
)

func (l *Lexer) analyze() []Token {
	// split content into tokens
	var tokens []Token
	var word string
	i := 0
	for i < len(l.content) {
		curr_char := l.content[i]
		next_char := l.content[i+1]

		// special chars
		switch curr_char {
		case '<':
		case '>':
		case ';':
			tokens = append(tokens, Token{purpose: TokenPurposeIdentifier, content: word})
			tokens = append(tokens, Token{purpose: TokenPurposeSpecial, content: string(curr_char)})
			word = ""
		case ' ', '\n':
			tokens = append(tokens, Token{purpose: TokenPurposeWhitespace, content: word})
			word = ""
		default:
			if curr_char == '"' {
				// strings
			} else if curr_char == '/' && next_char == '/' {
				// slash comments
			} else if curr_char == '/' && next_char == '*' {
				// star comments
			} else {
				word += string(curr_char)
			}
		}
		i++
	}
	if len(word) > 0 {
		tokens = append(tokens, Token{purpose: TokenPurposeIdentifier, content: word})
	}

	return tokens
}
