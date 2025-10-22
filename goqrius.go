// Package goqrius
package goqrius

import (
	"strings"

	"github.com/golaxo/goqrius/lexer"
)

// Parse the input filter expression to a goqrius Expression.
func Parse(input string) (Expression, error) {
	if input == "" {
		//nolint:nilnil // TODO think about returning something like EmptyExpression{}, nil.
		return nil, nil
	}

	l := lexer.New(input)
	p := New(l)
	e := p.Parse()

	return e, ParseError{errors: p.Errors()}
}

var _ error = new(ParseError)

type ParseError struct {
	errors []string
}

func (p ParseError) Error() string {
	return strings.Join(p.errors, ",")
}
