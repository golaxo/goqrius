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

type ParseError struct {
	errors []error
}

func (p ParseError) Error() string {
	errorsMessage := make([]string, len(p.errors))
	for i, err := range p.errors {
		errorsMessage[i] = err.Error()
	}

	return strings.Join(errorsMessage, ",")
}

type (
	UnexpectedTokenError struct {
		Token   token.Token
		Message string
	}
)

func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("%s, at position %d", e.Message, e.Token.Position)
}
