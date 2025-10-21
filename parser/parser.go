// Package parser parse the filter expression.
package parser

import (
	"fmt"

	"github.com/golaxo/goqrius/lexer"
	"github.com/golaxo/goqrius/token"
)

// Precedences.
const (
	_ int = iota
	LOWEST
	OR      // or
	AND     // and
	PREFIX  // not
	COMPARE // eq, ne, gt, ge, lt, le
)

//nolint:exhaustive,gochecknoglobals // no need to put all the tokens.
var precedences = map[token.Type]int{
	token.Or:                 OR,
	token.And:                AND,
	token.Eq:                 COMPARE,
	token.NotEq:              COMPARE,
	token.GreaterThan:        COMPARE,
	token.GreaterThanOrEqual: COMPARE,
	token.LessThan:           COMPARE,
	token.LessThanOrEqual:    COMPARE,
}

// AST nodes

type Node interface{ String() string }

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct{ Value string }

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}

type IntegerLiteral struct{ Value string }

func (il *IntegerLiteral) String() string  { return il.Value }
func (il *IntegerLiteral) expressionNode() {}

type StringLiteral struct{ Value string }

func (sl *StringLiteral) String() string  { return fmt.Sprintf("'%s'", sl.Value) }
func (sl *StringLiteral) expressionNode() {}

type NotExpr struct{ Right Expression }

func (ne *NotExpr) String() string  { return fmt.Sprintf("(not %s)", ne.Right.String()) }
func (ne *NotExpr) expressionNode() {}

type InfixExpr struct {
	Left     Expression
	Operator token.Type
	Right    Expression
}

func (ie *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), string(ie.Operator), ie.Right.String())
}
func (ie *InfixExpr) expressionNode() {}

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

	expr := p.parseExpression(LOWEST)

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

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}

	return LOWEST
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
		right := p.parseExpression(PREFIX)
		leftExp = &NotExpr{Right: right}
	case token.Lparen:
		p.nextToken()

		leftExp = p.parseExpression(LOWEST)
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
			leftExp = &InfixExpr{Left: leftExp, Operator: op, Right: right}
		default:
			return leftExp
		}
	}

	return leftExp
}
