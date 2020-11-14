package path

import (
	"unicode"
	"unicode/utf8"
)

// Lexer takes a query and breaks it down into tokens that can be consumed at
// at later date.
// The lexer in question is lazy and requires the calling of next to move it
// forward.
type Lexer struct {
	input []rune
	char  rune

	position     int
	readPosition int
	line         int
	column       int
}

// NewLexer creates a new Lexer from a given input.
func NewLexer(input string) *Lexer {
	lex := &Lexer{
		input:  []rune(input),
		char:   ' ',
		line:   1,
		column: 1,
	}

	return lex
}

// ReadNext will attempt to read the next character and correctly setup the
// positional values for the input.
func (l *Lexer) ReadNext() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
		if l.char == '\n' {
			l.column = 1
			l.line++
		} else {
			l.column++
		}
	}

	l.position = l.readPosition
	l.readPosition++
}

// Peek will attempt to read the next rune if it's available.
func (l *Lexer) Peek() rune {
	return l.PeekN(0)
}

// PeekN attempts to read the next rune by a given offset, it it's available.
func (l *Lexer) PeekN(n int) rune {
	if l.readPosition+n >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+n]
}

// NextToken attempts to grab the next token available.
func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	pos := l.getPosition()
	pos.Column--

	if t, ok := tokenMap[l.char]; ok {
		switch t {
		default:
			tok = MakeToken(t, l.char)
		}

		l.ReadNext()

		tok.Pos = pos
		return tok
	}

	newToken := l.readRunesToken()
	newToken.Pos = pos
	return newToken
}

func (l *Lexer) readRunesToken() Token {
	var tok Token
	switch {
	case l.char == 0:
		tok.Literal = ""
		tok.Type = EOF
		return tok
	case isLetter(l.char):
		tok.Literal = l.readIdentifier()
		tok.Type = IDENT
		return tok
	}
	l.ReadNext()
	return MakeToken(UNKNOWN, l.char)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.char) {
		l.ReadNext()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) || isDigit(l.char) || l.char == '-' {
		l.ReadNext()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) getPosition() Position {
	return Position{
		Offset: l.position,
		Line:   l.line,
		Column: l.column,
	}
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_' || char >= utf8.RuneSelf && unicode.IsLetter(char)
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9' || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

func isQuote(char rune) bool {
	return char == 34
}
