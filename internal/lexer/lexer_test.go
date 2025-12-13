package lexer

import (
	"testing"

	"github.com/golaxo/goqrius/internal/token"
)

func TestNextToken(t *testing.T) {
	t.Parallel()

	tdt := map[string]struct {
		input    string
		expected []struct {
			expectedType    token.Type
			expectedLiteral string
		}
	}{
		"simple equal condition": {
			input: `key eq value`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.Eq, string(token.Eq)},
				{token.Ident, "value"},
				{token.EOF, ""},
			},
		},
		"simple equal null": {
			input: `key eq null`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.Eq, string(token.Eq)},
				{token.Null, "null"},
				{token.EOF, ""},
			},
		},
		"simple not equal null": {
			input: `key ne null`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.NotEq, string(token.NotEq)},
				{token.Null, "null"},
				{token.EOF, ""},
			},
		},
		"simple not equal condition": {
			input: `key ne value`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.NotEq, string(token.NotEq)},
				{token.Ident, "value"},
				{token.EOF, ""},
			},
		},
		"simple greater than condition": {
			input: `key gt 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.GreaterThan, string(token.GreaterThan)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"simple greater than or equal condition": {
			input: `key ge 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.GreaterThanOrEqual, string(token.GreaterThanOrEqual)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"simple less than condition": {
			input: `key lt 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.LessThan, string(token.LessThan)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"simple less than or equal condition": {
			input: `key le 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "key"},
				{token.LessThanOrEqual, string(token.LessThanOrEqual)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"concatenate and condition": {
			input: `name eq john and age le 50`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "name"},
				{token.Eq, string(token.Eq)},
				{token.Ident, "john"},
				{token.And, string(token.And)},
				{token.Ident, "age"},
				{token.LessThanOrEqual, string(token.LessThanOrEqual)},
				{token.Int, "50"},
				{token.EOF, ""},
			},
		},
		"concatenate or condition": {
			input: `name eq john or age le 50`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "name"},
				{token.Eq, string(token.Eq)},
				{token.Ident, "john"},
				{token.Or, string(token.Or)},
				{token.Ident, "age"},
				{token.LessThanOrEqual, string(token.LessThanOrEqual)},
				{token.Int, "50"},
				{token.EOF, ""},
			},
		},
		"not condition": {
			input: `not name eq 'john'`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Not, string(token.Not)},
				{token.Ident, "name"},
				{token.Eq, string(token.Eq)},
				{token.String, "john"},
				{token.EOF, ""},
			},
		},
		"not null condition": {
			input: `not null`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Not, string(token.Not)},
				{token.Null, string(token.Null)},
				{token.EOF, ""},
			},
		},
		"parenthesis condition": {
			input: `(age gt 0 and age le 18)`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Lparen, string(token.Lparen)},
				{token.Ident, "age"},
				{token.GreaterThan, string(token.GreaterThan)},
				{token.Int, "0"},
				{token.And, string(token.And)},
				{token.Ident, "age"},
				{token.LessThanOrEqual, string(token.LessThanOrEqual)},
				{token.Int, "18"},
				{token.Rparen, string(token.Rparen)},
			},
		},
		"parenthesis around null": {
			input: `(null)`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Lparen, string(token.Lparen)},
				{token.Null, string(token.Null)},
				{token.Rparen, string(token.Rparen)},
			},
		},
		"uppercased NULL is identifier": {
			input: `NULL eq 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "NULL"},
				{token.Eq, string(token.Eq)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"capitalized Null is identifier": {
			input: `Null eq 1`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "Null"},
				{token.Eq, string(token.Eq)},
				{token.Int, "1"},
				{token.EOF, ""},
			},
		},
		"multiple conditions with parenthesis": {
			input: `name eq 'John' or (age gt 0 and age le 18)`,
			expected: []struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.Ident, "name"},
				{token.Eq, string(token.Eq)},
				{token.String, "John"},
				{token.Or, string(token.Or)},
				{token.Lparen, string(token.Lparen)},
				{token.Ident, "age"},
				{token.GreaterThan, string(token.GreaterThan)},
				{token.Int, "0"},
				{token.And, string(token.And)},
				{token.Ident, "age"},
				{token.LessThanOrEqual, string(token.LessThanOrEqual)},
				{token.Int, "18"},
				{token.Rparen, string(token.Rparen)},
			},
		},
	}

	for name, test := range tdt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			input := test.input

			l := New(input)

			for i, tt := range test.expected {
				tok := l.NextToken()

				if tok.Type != tt.expectedType {
					t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
				}

				if tok.Literal != tt.expectedLiteral {
					t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
				}
			}
		})
	}
}
