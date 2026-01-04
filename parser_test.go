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
		"identifier mixing characters and number": {
			input:          "nam3 eq 'John",
			expectedString: "(nam3 eq 'John')",
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
		description    string
		expectedErrors []error
	}{
		"name eq value": {
			description: "identifier as value",
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
		"(name eq 'John'": {
			description: "missing closing paren",
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
		"@ eq 1": {
			description: "illegal character",
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
		"name is null": {
			description: "unknown operator 'is'",
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
		"name not null": {
			description: "unknown operator not",
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
		"1 gt 2": {
			description: "int as left side",
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
		"'name' eq 'John'": {
			description: "string in left side",
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
		"not null": {
			description: "not null is invalid",
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
		"null": {
			description: "bare null is invalid",
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
		"(null)": {
			description: "paren bare null is invalid",
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
		"null eq name": {
			description: "null on left is invalid",
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
		"name gt null": {
			description: "gt null is invalid",
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
		"name ge null": {
			description: "ge null is invalid",
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
		"name lt null": {
			description: "lt null is invalid",
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
		"name le null": {
			description: "le null is invalid",
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
		"name eq not null": {
			description: "eq not null is invalid",
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
		"name eq (not null)": {
			description: "eq (not null) is invalid",
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
		"n@me eq 'John'": {
			description: "illegal character in identifier",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Illegal,
						Literal:  "@",
						Position: 1,
					},
					Message: "illegal token \"@\"",
				},
			},
		},
		"name @q 'John'": {
			description: "illegal character in operator",
			expectedErrors: []error{
				UnexpectedTokenError{
					Token: token.Token{
						Type:     token.Illegal,
						Literal:  "@",
						Position: 5,
					},
					Message: "illegal token \"@\"",
				},
			},
		},
	}

	for input, tt := range tests {
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			p := newParser(lexer.New(input))
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
