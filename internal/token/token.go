// Package token contains all the token.
package token

const (
	// Illegal when an operator isn't allowed.
	Illegal Type = "Illegal"
	// EOF end of filter value.
	EOF Type = "EOF"

	/* Identifier + Literals. */

	// Ident is the type of identifiers.
	Ident Type = "Ident"
	// Int is the type of integer literals.
	Int Type = "Int"
	// String is the type of string literals.
	String Type = "String"
	// Null is the type of null literals.
	Null Type = "null"

	/* Comparison Operators. */

	// Eq is the type of equality operator.
	Eq Type = "eq"
	// NotEq is the type of inequality operator.
	NotEq Type = "ne"
	// GreaterThan is the type of greater than operator.
	GreaterThan Type = "gt"
	// GreaterThanOrEqual is the type of greater than or equal operator.
	GreaterThanOrEqual Type = "ge"
	// LessThan is the type of less than operator.
	LessThan Type = "lt"
	// LessThanOrEqual is the type of less than or equal operator.
	LessThanOrEqual Type = "le"

	/* Logical Operators. */

	// And is the type of logical and operator.
	And Type = "and"
	// Or is the type of logical or operator.
	Or Type = "or"
	// Not is the type of logical not operator.
	Not Type = "not"

	// Lparen is the type of left parenthesis.
	Lparen Type = "("
	// Rparen is the type of right parenthesis.
	Rparen Type = ")"
	// Lbrace is the type of left brace.
	Lbrace Type = "{"
	// Rbrace is the type of right brace.
	Rbrace Type = "}"
)

type (
	// Type is the type of a token.
	Type string

	// Token holds the actual type and its value.
	Token struct {
		// Type of the token.
		Type Type
		// The actual value for the token.
		Literal string
		// Position of the token.
		Position int
	}
)
