// Package goqrius
package goqrius

import (
	"github.com/golaxo/goqrius/lexer"
)

// Parse the input filter expression to a goqrius Expression.
func Parse(input string) (Expression, error) {
	if input == "" {
		//nolint:nilnil // TODO think about returning something like EmptyExpression{}, nil.
		return nil, nil
	}

	l := lexer.New(input)
	p := newParser(l)
	e := p.parse()

	var err error
	if len(p.Errors()) > 0 {
		err = ParseError{errors: p.Errors()}
	}

	return e, err
}

func MustParse(input string) Expression {
	e, err := Parse(input)
	if err != nil {
		panic(err)
	}

	return e
}
