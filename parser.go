package goqrius

import (
	"fmt"

	"github.com/golaxo/goqrius/internal/lexer"
	"github.com/golaxo/goqrius/internal/token"
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

type parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

// New creates a new Parser based on a Lexer.
// It's recommended to use goqrius.Parse instead of this.
func newParser(l *lexer.Lexer) *parser {
	p := &parser{
		l:      l,
		errors: make([]string, 0),
	}
	// Initialize tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *parser) Errors() []string { return p.errors }

// Parse parses the whole input into an Expression AST.
func (p *parser) parse() Expression {
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

	// validate semantic constraints (e.g., null usage)
	p.validate(expr, true, nil)

	return expr
}

func (p *parser) nextToken() { p.curToken = p.peekToken; p.peekToken = p.l.NextToken() }

func (p *parser) expectPeek(t token.Type) bool {
	if p.peekToken.Type == t {
		p.nextToken()

		return true
	}

	p.peekError(t)

	return false
}

func (p *parser) peekPrecedence() int {
	if pr, ok := precedences[p.peekToken.Type]; ok {
		return pr
	}

	return lowest
}

func (p *parser) curPrecedence() int {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}

	return lowest
}

func (p *parser) peekError(t token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("expected next token to be %q, got %q instead", t, p.peekToken.Type))
}

func (p *parser) parseExpression(precedence int) Expression {
	var leftExp Expression

	//nolint:exhaustive // no need.
	switch p.curToken.Type {
	case token.Ident:
		leftExp = &Identifier{Value: p.curToken.Literal}
	case token.Int:
		leftExp = &IntegerLiteral{Value: p.curToken.Literal}
	case token.String:
		leftExp = &StringLiteral{Value: p.curToken.Literal}
	case token.Null:
		leftExp = &Null{}
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

func (p *parser) validate(expr Expression, isLeft bool, parentFilterExpression *FilterExpr) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *Null:
		if isLeft {
			p.errors = append(p.errors, "invalid use of null: only allowed as right side of 'eq' or 'ne'")

			return
		}

		if parentFilterExpression == nil ||
			parentFilterExpression.Operator != token.Eq &&
				parentFilterExpression.Operator != token.NotEq {
			p.errors = append(p.errors, "invalid use of null: only allowed as right side of 'eq' or 'ne'")
		}
	case *Identifier:
		return
	case *IntegerLiteral, *StringLiteral:
		if isLeft {
			p.errors = append(p.errors, "invalid use of literals: only allowed as right side")
		}
	case *NotExpr:
		// `not null` (directly or via grouping) is invalid
		p.validate(e.Right, false, nil)
	case *FilterExpr:
		// Left side can never be null
		p.validate(e.Left, true, e)
		p.validate(e.Right, false, e)
	}
}
