package goqrius

import (
	"testing"

	"github.com/golaxo/goqrius/internal/lexer"
	"github.com/golaxo/goqrius/internal/token"
)

func TestParseExpressions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedString string
	}{
		"simple ident eq null": {
			input:          "name eq null",
			expectedString: "(name eq null)",
		},
		"simple ident ne null": {
			input:          "name ne null",
			expectedString: "(name ne null)",
		},
		"ident eq string": {
			input:          "name eq 'john'",
			expectedString: "(name eq 'john')",
		},
		"gt and le with precedence": {
			input:          "age gt 0 and age le 18",
			expectedString: "((age gt 0) and (age le 18))",
		},
		"or has lower precedence than and": {
			input:          "name eq 'John' or age gt 0 and age le 18",
			expectedString: "((name eq 'John') or ((age gt 0) and (age le 18)))",
		},
		"parentheses grouping": {
			input:          "(age gt 0 and age le 18)",
			expectedString: "((age gt 0) and (age le 18))",
		},
		"not binds tighter than and": {
			input:          "not name eq 'john' and age le 50",
			expectedString: "((not (name eq 'john')) and (age le 50))",
		},
		"identifier with dot": {
			input:          "user.name eq 'john'",
			expectedString: "(user.name eq 'john')",
		},
		"identifier named NULL stays ident (case-sensitive)": {
			input:          "NULL eq 1",
			expectedString: "(NULL eq 1)",
		},
		"identifier named Null stays ident (case-sensitive)": {
			input:          "Null eq 1",
			expectedString: "(Null eq 1)",
		},
		"identifier with dash": {
			input:          "user-name eq 1",
			expectedString: "(user-name eq 1)",
		},
		"mixed precedence with null": {
			input:          "name eq null or not age ge 18",
			expectedString: "((name eq null) or (not (age ge 18)))",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("err not expected; error=%v", err)
			}

			if expr == nil {
				t.Fatalf("expected non-nil expression")
			}

			if got := expr.String(); got != tt.expectedString {
				t.Fatalf("unexpected AST string. expected=%q got=%q", tt.expectedString, got)
			}
		})
	}
}

func TestEmptyInput(t *testing.T) {
	t.Parallel()

	expr, err := Parse("")
	if err != nil {
		t.Fatalf("err not expected; error=%v", err)
	}

	if expr != nil {
		t.Fatalf("expected nil expression")
	}
}

func TestMustParseEmptyInput(t *testing.T) {
	t.Parallel()

	expr := MustParse("")
	if expr != nil {
		t.Fatalf("expected nil expression")
	}
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedErrors []error
	}{
		"identifier as value": {
			input: "name eq value",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Ident,
						Literal:  "value",
						Position: 8,
					},
					Message: "identifier can not be used as value",
				},
			},
		},
		"missing closing paren": {
			input: "(name eq 'John'",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.EOF,
						Literal:  "",
						Position: 15,
					},
					Message: "expected next token to be \")\", got \"\"",
				},
			},
		},
		"unknown prefix": {
			input: "@ eq 1",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Illegal,
						Literal:  "@",
						Position: 0,
					},
					Message: "illegal token \"@\"",
				},
			},
		},
		"unknown operator is": {
			input: "name is null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Ident,
						Literal:  "is",
						Position: 5,
					},
					Message: "expected next token to be an operator, got \"is\"",
				},
			},
		},
		"unknown operator not": {
			input: "name not null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Not,
						Literal:  "not",
						Position: 5,
					},
					Message: "expected next token to be an operator, got \"not\"",
				},
			},
		},
		"int as left side": {
			input: "1 gt 2",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Int,
						Literal:  "1",
						Position: 0,
					},
					Message: LeftSideMustBeIdentifier,
				},
			},
		},
		"string in left side": {
			input: "'name' eq 'John'",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.String,
						Literal:  "name",
						Position: 0,
					},
					Message: LeftSideMustBeIdentifier,
				},
			},
		},
		"not null is invalid": {
			input: "not null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 4,
					},
					Message: "'not' can not be applied to a value",
				},
			},
		},
		"bare null is invalid": {
			input: "null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 0,
					},
					Message: "'null' can not be used as a standalone expression",
				},
			},
		},
		"paren bare null is invalid": {
			input: "(null)",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Rparen,
						Literal:  ")",
						Position: 5,
					},
					Message: "grouped value is not a valid expression",
				},
			},
		},
		"null on left is invalid": {
			input: "null eq name",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 0,
					},
					Message: LeftSideMustBeIdentifier,
				},
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Ident,
						Literal:  "name",
						Position: 8,
					},
					Message: "identifier can not be used as value",
				},
			},
		},
		"gt null is invalid": {
			input: "name gt null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 8,
					},
					Message: NullCannotBeUsedWithComparison,
				},
			},
		},
		"ge null is invalid": {
			input: "name ge null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 8,
					},
					Message: NullCannotBeUsedWithComparison,
				},
			},
		},
		"lt null is invalid": {
			input: "name lt null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 8,
					},
					Message: NullCannotBeUsedWithComparison,
				},
			},
		},
		"le null is invalid": {
			input: "name le null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 8,
					},
					Message: NullCannotBeUsedWithComparison,
				},
			},
		},
		"eq not null is invalid": {
			input: "name eq not null",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Not,
						Literal:  "not",
						Position: 8,
					},
					Message: "invalid value token \"not\"",
				},
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Null,
						Literal:  "null",
						Position: 12,
					},
					Message: "unexpected token \"null\"",
				},
			},
		},
		"eq (not null) is invalid": {
			input: "name eq (not null)",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Not,
						Literal:  "not",
						Position: 13,
					},
					Message: "'not' can not be applied to a value",
				},
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Not,
						Literal:  "not",
						Position: 9,
					},
					Message: "right side of comparison must be a value",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			p := newParser(lexer.New(tt.input))
			_ = p.parse()

			if len(p.Errors()) == 0 {
				t.Fatalf("expected errors, got none")
			}

			if len(p.Errors()) != len(tt.expectedErrors) {
				t.Fatalf("expected %d errors, got %d", len(tt.expectedErrors), len(p.Errors()))
			}

			for i, err := range p.Errors() {
				if err.Error() != tt.expectedErrors[i].Error() {
					t.Fatalf("expected error at [%d]: %v, got: %v", i, tt.expectedErrors[i], err)
				}
			}
		})
	}
}
