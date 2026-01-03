package lexer

import (
	"unicode"

	"github.com/golaxo/goqrius/internal/token"
)

// Lexer to parse the input.
type Lexer struct {
	input string
	// Current position in input (points to current char)
	position int
	// Current reading position in input (after current char)
	readPosition int
}

// New creates a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	// Initialize positions so that getChar works correctly
	l.position = 0
	l.readPosition = 0

	return l
}

// NextToken returns the next token parsed, or token.EOF if finished.
//
//nolint:funlen // refactor later
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	startPos := l.readPosition

	ch, ok := l.peekChar()
	if !ok {
		return token.Token{Type: token.EOF, Literal: "", Position: startPos}
	}

	// Single-character tokens (delimiters)
	switch ch {
	case '(':
		l.readChar()

		return token.Token{Type: token.Lparen, Literal: string(token.Lparen), Position: startPos}
	case ')':
		l.readChar()

		return token.Token{Type: token.Rparen, Literal: string(token.Rparen), Position: startPos}
	case '{':
		l.readChar()

		return token.Token{Type: token.Lbrace, Literal: string(token.Lbrace), Position: startPos}
	case '}':
		l.readChar()

		return token.Token{Type: token.Rbrace, Literal: string(token.Rbrace), Position: startPos}
	case '\'':
		// String literal
		str := l.readSingleQuoted()

		return token.Token{Type: token.String, Literal: str, Position: startPos}
	}

	// Numbers
	if isDigit(ch) {
		num := l.readWhile(isDigit)

		return token.Token{Type: token.Int, Literal: num, Position: startPos}
	}

	// Identifiers and keywords (and, or, not, eq, ne, gt, ge, lt, le)
	ident := l.readWhile(isIdentChar)
	// If we didn't read any identifier characters, this is an unknown/illegal character.
	if ident == "" {
		// consume the offending character and return Illegal so the parser can handle it
		if offendingCh, isOk := l.peekChar(); isOk {
			l.readChar()

			return token.Token{Type: token.Illegal, Literal: string(offendingCh), Position: startPos}
		}

		return token.Token{Type: token.EOF, Literal: "", Position: startPos}
	}

	switch ident {
	case string(token.Null):
		return newTokenFromType(token.Null, startPos)
	case string(token.And):
		return newTokenFromType(token.And, startPos)
	case string(token.Or):
		return newTokenFromType(token.Or, startPos)
	case string(token.Not):
		return newTokenFromType(token.Not, startPos)
	case string(token.Eq):
		return newTokenFromType(token.Eq, startPos)
	case string(token.NotEq):
		return newTokenFromType(token.NotEq, startPos)
	case string(token.GreaterThan):
		return newTokenFromType(token.GreaterThan, startPos)
	case string(token.GreaterThanOrEqual):
		return newTokenFromType(token.GreaterThanOrEqual, startPos)
	case string(token.LessThan):
		return newTokenFromType(token.LessThan, startPos)
	case string(token.LessThanOrEqual):
		return newTokenFromType(token.LessThanOrEqual, startPos)
	default:
		return token.Token{Type: token.Ident, Literal: ident, Position: startPos}
	}
}

// readChar advances the cursor by one rune (byte, since input is ASCII expected).
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.position = l.readPosition
	}

	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns current unread char without consuming.
func (l *Lexer) peekChar() (byte, bool) {
	if l.readPosition >= len(l.input) {
		return 0, false
	}

	return l.input[l.readPosition], true
}

func (l *Lexer) skipWhitespace() {
	for {
		ch, ok := l.peekChar()
		if !ok || !unicode.IsSpace(rune(ch)) {
			return
		}

		l.readChar()
	}
}

func (l *Lexer) readWhile(pred func(byte) bool) string {
	start := l.readPosition
	for ch, ok := l.peekChar(); ok && pred(ch); ch, ok = l.peekChar() {
		l.readChar()
	}

	return l.input[start:l.readPosition]
}

// readSingleQuoted reads content inside single quotes, consuming both quotes.
// If no closing quote is found, it reads until end and returns what was found (without the opening quote).
func (l *Lexer) readSingleQuoted() string {
	// consume opening quote
	l.readChar()

	start := l.readPosition
	for ch, ok := l.peekChar(); ok; ch, ok = l.peekChar() {
		if ch == '\'' {
			// end of string
			literal := l.input[start:l.readPosition]
			l.readChar() // consume closing quote

			return literal
		}

		l.readChar()
	}
	// EOF reached without closing quote
	return l.input[start:l.readPosition]
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

func isIdentStart(ch byte) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isIdentChar(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '-' || ch == '.'
}

func newTokenFromType(t token.Type, position int) token.Token {
	return token.Token{Type: t, Literal: string(t), Position: position}
}
