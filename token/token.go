// Package token contains all the token.
package token

const (
	// Illegal when an operator isn't allowed.
	Illegal Type = "Illegal"
	// EOF end of filter value.
	EOF Type = "EOF"

	/* Identifier + Literals. */

	Ident  Type = "Ident"
	Int    Type = "Int"
	String Type = "String"

	/* Comparison Operators. */

	Eq                 Type = "eq"
	NotEq              Type = "ne"
	GreaterThan        Type = "gt"
	GreaterThanOrEqual Type = "ge"
	LessThan           Type = "lt"
	LessThanOrEqual    Type = "le"

	/* Logical Operators. */

	And Type = "and"
	Or  Type = "or"
	Not Type = "not"

	Lparen Type = "("
	Rparen Type = ")"
	Lbrace Type = "{"
	Rbrace Type = "}"
)

type (
	Type string

	// Token holds the actual type and its value.
	Token struct {
		// Type of the token.
		Type Type
		// The actual value for the token.
		Literal string
	}
)
