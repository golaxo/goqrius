package goqrius

import (
	"fmt"
)

type FilterOperator string

const (
	Eq                 FilterOperator = "eq"
	NotEq              FilterOperator = "ne"
	GreaterThan        FilterOperator = "gt"
	GreaterThanOrEqual FilterOperator = "ge"
	LessThan           FilterOperator = "lt"
	LessThanOrEqual    FilterOperator = "le"
)

type (
	Node interface {
		String() string
	}

	Expression interface {
		Node
		expressionNode()
	}

	Value interface {
		Expression
		valueNode()
	}

	// AndExpr and concatenates FilterExpr.
	AndExpr struct {
		Left  Expression
		Right Expression
	}

	// OrExpr or concatenates FilterExpr.
	OrExpr struct {
		Left  Expression
		Right Expression
	}

	// NotExpr negates an Expression.
	NotExpr struct {
		Right Expression
	}

	// FilterExpr represents a key and operator and a value in a filter clause.
	FilterExpr struct {
		Left     Identifier
		Operator FilterOperator
		Right    Value
	}

	// Identifier is the Expression to indicate the key of a filter clause, e.g. `name`.
	Identifier struct {
		Value string
	}

	// IntegerLiteral is the Expression to indicate an int value of a filter clause, e.g. `1`.
	IntegerLiteral struct {
		Value string
	}

	// Null is the Expression to indicate a value that is null.
	Null struct{}

	// StringLiteral is the Expression to indicate an int value of a filter clause, e.g. `'John'`.
	StringLiteral struct {
		Value string
	}
)

func (ae *AndExpr) String() string {
	return fmt.Sprintf("(%s and %s)", ae.Left.String(), ae.Right.String())
}
func (ae *AndExpr) expressionNode() {}

func (oe *OrExpr) String() string {
	return fmt.Sprintf("(%s or %s)", oe.Left.String(), oe.Right.String())
}
func (oe *OrExpr) expressionNode() {}

func (ne *NotExpr) String() string  { return fmt.Sprintf("(not %s)", ne.Right.String()) }
func (ne *NotExpr) expressionNode() {}

func (ie *FilterExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), string(ie.Operator), ie.Right.String())
}
func (ie *FilterExpr) expressionNode() {}

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) valueNode()      {}

func (il *IntegerLiteral) String() string  { return il.Value }
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) valueNode()      {}

func (n *Null) String() string  { return "null" }
func (n *Null) expressionNode() {}
func (n *Null) valueNode()      {}

func (sl *StringLiteral) String() string  { return fmt.Sprintf("'%s'", sl.Value) }
func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) valueNode()      {}
