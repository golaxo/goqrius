package goqrius

import (
	"fmt"

	"github.com/golaxo/goqrius/lexer"
	"github.com/golaxo/goqrius/token"
)

// Precedences.
const (
	_ int = iota
	lowest
	or      // or
	and     // and
	prefix  // not
	compare // eq, ne, gt, ge, lt, le
)

//nolint:exhaustive,gochecknoglobals // no need to put all the tokens.
var precedences = map[token.Type]int{
	token.Or:                 or,
	token.And:                and,
	token.Eq:                 compare,
	token.NotEq:              compare,
	token.GreaterThan:        compare,
	token.GreaterThanOrEqual: compare,
	token.LessThan:           compare,
	token.LessThanOrEqual:    compare,
}

// AST nodes

type Node interface{ String() string }

type Expression interface {
	Node
	expressionNode()
}

// Identifier is the Expression to indicate the key of a filter clause, e.g. `name`.
type Identifier struct{ Value string }

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}

// IntegerLiteral is the Expression to indicate an int value of a filter clause, e.g. `1`.
type IntegerLiteral struct{ Value string }

func (il *IntegerLiteral) String() string  { return il.Value }
func (il *IntegerLiteral) expressionNode() {}

// StringLiteral is the Expression to indicate an int value of a filter clause, e.g. `'John'`.
type StringLiteral struct{ Value string }

func (sl *StringLiteral) String() string  { return fmt.Sprintf("'%s'", sl.Value) }
func (sl *StringLiteral) expressionNode() {}

// NotExpr negates an Expression.
type NotExpr struct{ Right Expression }

func (ne *NotExpr) String() string  { return fmt.Sprintf("(not %s)", ne.Right.String()) }
func (ne *NotExpr) expressionNode() {}

// FilterExpr represents a key and operator and a value in a filter clause.
type FilterExpr struct {
	Left     Expression
	Operator token.Type
	Right    Expression
}

func (ie *FilterExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), string(ie.Operator), ie.Right.String())
}
func (ie *FilterExpr) expressionNode() {}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]string, 0),
	}
	// Initialize tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string { return p.errors }

// Parse parses the whole input into an Expression AST.
func (p *Parser) Parse() Expression {
	// Start by advancing to the first actual token if cur is zero-value
	if p.curToken.Type == "" && p.peekToken.Type == "" {
		p.nextToken()
		p.nextToken()
	}

	expr := p.parseExpression(lowest)

	// consume trailing tokens until EOF
	for p.peekToken.Type != token.EOF {
		p.nextToken()
	}

	return expr
}

func (p *Parser) nextToken() { p.curToken = p.peekToken; p.peekToken = p.l.NextToken() }

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekToken.Type == t {
		p.nextToken()

		return true
	}

	p.peekError(t)

	return false
}

func (p *Parser) peekPrecedence() int {
	if pr, ok := precedences[p.peekToken.Type]; ok {
		return pr
	}

	return lowest
}

func (p *Parser) curPrecedence() int {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}

	return lowest
}

func (p *Parser) peekError(t token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("expected next token to be %q, got %q instead", t, p.peekToken.Type))
}

func (p *Parser) parseExpression(precedence int) Expression {
	var leftExp Expression

	//nolint:exhaustive // no need.
	switch p.curToken.Type {
	case token.Ident:
		leftExp = &Identifier{Value: p.curToken.Literal}
	case token.Int:
		leftExp = &IntegerLiteral{Value: p.curToken.Literal}
	case token.String:
		leftExp = &StringLiteral{Value: p.curToken.Literal}
	case token.Not:
		p.nextToken()
		right := p.parseExpression(prefix)
		leftExp = &NotExpr{Right: right}
	case token.Lparen:
		p.nextToken()

		leftExp = p.parseExpression(lowest)
		if !p.expectPeek(token.Rparen) {
			return nil
		}
	default:
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %q found", p.curToken.Type))

		return nil
	}

	for p.peekToken.Type != token.EOF && precedence < p.peekPrecedence() {
		op := p.peekToken.Type
		// Only logical and comparison operators are infix here
		//nolint:exhaustive // no need to check all the keys.
		switch op {
		case token.And, token.Or, token.Eq, token.NotEq,
			token.GreaterThan, token.GreaterThanOrEqual, token.LessThan, token.LessThanOrEqual:
			p.nextToken() // advance to operator
			prec := p.curPrecedence()
			p.nextToken() // advance to the right expression's first token
			right := p.parseExpression(prec)
			leftExp = &FilterExpr{Left: leftExp, Operator: op, Right: right}
		default:
			return leftExp
		}
	}

	return leftExp
}
