package goqrius

import (
	"fmt"
	"strings"

	"github.com/golaxo/goqrius/internal/token"
)

var (
	_ error = new(ParseError)
	_ error = new(UnexpectedTokenError)
)

const (
	// LeftSideMustBeIdentifier message indicates that the left side of a
	// comparison must be an identifier.
	LeftSideMustBeIdentifier = "left side of comparison must be an identifier"
	// NullCannotBeUsedWithComparison message indicates that the null value cannot be
	// used with comparison operators.
	NullCannotBeUsedWithComparison = "'null' cannot be used with comparison operator"
)

type (
	// ParseError is the error type returned by the parser.
	ParseError struct {
		errors []error
	}

	// UnexpectedTokenError is the error type returned when the lexer finds an unexpected token.
	UnexpectedTokenError struct {
		Token   token.Token
		Message string
	}
)

// Error implements the error interface and returns all the parsing errors that happened.
func (p ParseError) Error() string {
	errorsMessage := make([]string, len(p.errors))
	for i, err := range p.errors {
		errorsMessage[i] = err.Error()
	}

	return strings.Join(errorsMessage, ",")
}

// Error implements the error interface.
func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("%s, at position %d", e.Message, e.Token.Position)
}
