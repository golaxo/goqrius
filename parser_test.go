package goqrius

import (
	"testing"

	"github.com/golaxo/goqrius/lexer"
)

func TestParseExpressions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedString string
	}{
		"simple ident eq ident": {
			input:          "name eq value",
			expectedString: "(name eq value)",
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
		"identifier with dash": {
			input:          "user-name eq 1",
			expectedString: "(user-name eq 1)",
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

func TestParseErrors(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"missing closing paren": "(name eq 'John'",
		"unknown prefix":        "@ eq 1",
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			l := lexer.New(input)
			p := New(l)
			_ = p.Parse()

			if len(p.Errors()) == 0 {
				t.Fatalf("expected errors, got none for input: %q", input)
			}
		})
	}
}
