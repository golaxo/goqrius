package goqrius

import (
	"fmt"

	"github.com/golaxo/goqrius/token"
)

type (
	Node interface {
		String() string
	}

	Expression interface {
		Node
		expressionNode()
	}
)

type (
	// Identifier is the Expression to indicate the key of a filter clause, e.g. `name`.
	Identifier struct {
		Value string
	}

	// IntegerLiteral is the Expression to indicate an int value of a filter clause, e.g. `1`.
	IntegerLiteral struct {
		Value string
	}

	// StringLiteral is the Expression to indicate an int value of a filter clause, e.g. `'John'`.
	StringLiteral struct {
		Value string
	}
)

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}

func (il *IntegerLiteral) String() string  { return il.Value }
func (il *IntegerLiteral) expressionNode() {}

func (sl *StringLiteral) String() string  { return fmt.Sprintf("'%s'", sl.Value) }
func (sl *StringLiteral) expressionNode() {}

type (
	// NotExpr negates an Expression.
	NotExpr struct {
		Right Expression
	}

	// FilterExpr represents a key and operator and a value in a filter clause.
	FilterExpr struct {
		Left     Expression
		Operator token.Type
		Right    Expression
	}
)

func (ne *NotExpr) String() string  { return fmt.Sprintf("(not %s)", ne.Right.String()) }
func (ne *NotExpr) expressionNode() {}

func (ie *FilterExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), string(ie.Operator), ie.Right.String())
}
func (ie *FilterExpr) expressionNode() {}
