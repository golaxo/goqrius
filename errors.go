package goqrius

import "strings"

var _ error = new(ParseError)

type ParseError struct {
	errors []string
}

func (p ParseError) Error() string {
	return strings.Join(p.errors, ",")
}
