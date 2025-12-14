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

// New creates a new parser based on a lexer.Lexer.
func newParser(l *lexer.Lexer) *parser {
	p := &parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *parser) Errors() []string { return p.errors }

func (p *parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *parser) parse() Expression {
	// handle empty input
	if p.curToken.Type == token.EOF && p.peekToken.Type == token.EOF {
		return nil
	}

	expr := p.parseExpression(lowest)

	// if expression is a bare value or identifier, it's invalid
	switch expr.(type) {
	case *IntegerLiteral, *StringLiteral, *Null, *Identifier:
		p.errors = append(p.errors, "invalid bare expression")
	}

	// consume remaining tokens and mark errors if any leftover meaningful tokens
	for p.peekToken.Type != token.EOF {
		p.nextToken()

		if p.curToken.Type != token.EOF {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token %q", p.curToken.Literal))
		}
	}

	return expr
}

//nolint:exhaustive,funlen,gocognit // refactor later
func (p *parser) parseExpression(precedence int) Expression {
	var leftExp Expression

	switch p.curToken.Type {
	case token.Ident:
		leftExp = &Identifier{Value: p.curToken.Literal}
	case token.Int:
		// bare int is invalid as an expression, record error but continue
		leftExp = &IntegerLiteral{Value: p.curToken.Literal}
	case token.String:
		leftExp = &StringLiteral{Value: p.curToken.Literal}
	case token.Null:
		// bare null is invalid
		leftExp = &Null{}

		p.errors = append(p.errors, "null cannot be used as a standalone expression")
	case token.Not:
		// prefix not
		p.nextToken()
		right := p.parseExpression(prefix)
		// Disallow 'not' applied to a bare value
		switch right.(type) {
		case *IntegerLiteral, *StringLiteral, *Null:
			p.errors = append(p.errors, "not cannot be applied to a value")
		case nil:
			p.errors = append(p.errors, "missing expression after not")
		}

		leftExp = &NotExpr{Right: right}
	case token.Lparen:
		// consume '(' and parse subexpression
		p.nextToken()

		inner := p.parseExpression(lowest)
		p.expectPeek(token.Rparen)
		// Disallow grouping a bare value as a full expression like (null)
		switch inner.(type) {
		case *IntegerLiteral, *StringLiteral, *Null:
			p.errors = append(p.errors, "grouped value is not a valid expression")
		}

		leftExp = inner
	case token.Illegal:
		p.errors = append(p.errors, fmt.Sprintf("illegal token %q", p.curToken.Literal))

		return nil
	default:
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %q", p.curToken.Literal))

		return nil
	}

	// Infix parsing loop
	for p.peekToken.Type != token.EOF && precedence < p.peekPrecedence() {
		switch p.peekToken.Type {
		case token.And:
			p.nextToken() // move to 'and'
			opPrec := p.curPrecedence()
			p.nextToken() // move to right prefix
			right := p.parseExpression(opPrec)
			leftExp = &AndExpr{Left: leftExp, Right: right}
		case token.Or:
			p.nextToken() // move to 'or'
			opPrec := p.curPrecedence()
			p.nextToken()
			right := p.parseExpression(opPrec)
			leftExp = &OrExpr{Left: leftExp, Right: right}
		case token.Eq, token.NotEq, token.GreaterThan, token.GreaterThanOrEqual, token.LessThan, token.LessThanOrEqual:
			// comparisons bind tighter than and/or
			p.nextToken() // move to operator
			operator := p.curToken.Type

			// left must be identifier
			ident, ok := leftExp.(*Identifier)
			if !ok {
				p.errors = append(p.errors, "left side of comparison must be an identifier")
			}

			// parse right value
			p.nextToken()
			val := p.parseValue()

			// validate null with comparison
			if _, isNull := val.(*Null); isNull {
				switch operator {
				case token.GreaterThan, token.GreaterThanOrEqual, token.LessThan, token.LessThanOrEqual:
					p.errors = append(p.errors, "null cannot be used with comparison operator")
				}
			}

			if ident == nil {
				// fabricate to proceed
				ident = &Identifier{Value: ""}
			}

			leftExp = &FilterExpr{Left: *ident, Operator: FilterOperator(operator), Right: val}
		default:
			return leftExp
		}
	}

	return leftExp
}

//nolint:exhaustive // no need to check all the tokens.
func (p *parser) parseValue() Value {
	switch p.curToken.Type {
	case token.Int:
		return &IntegerLiteral{Value: p.curToken.Literal}
	case token.String:
		return &StringLiteral{Value: p.curToken.Literal}
	case token.Null:
		return &Null{}
	case token.Ident:
		p.errors = append(p.errors, "identifier cannot be used as value")

		return nil
	case token.Lparen:
		// value cannot be a grouped expression (e.g., (not null)) per tests
		p.nextToken()

		inner := p.parseExpression(lowest)
		p.expectPeek(token.Rparen)

		p.errors = append(p.errors, "invalid value expression")
		// try to convert to some value to continue
		switch v := inner.(type) {
		case *Identifier:
			return v
		case *IntegerLiteral:
			return v
		case *StringLiteral:
			return v
		case *Null:
			return v
		default:
			return &Identifier{Value: ""}
		}
	default:
		p.errors = append(p.errors, fmt.Sprintf("invalid value token %q", p.curToken.Literal))

		return &Identifier{Value: ""}
	}
}

func (p *parser) expectPeek(t token.Type) {
	if p.peekToken.Type == t {
		p.nextToken()

		return
	}

	p.errors = append(
		p.errors,
		fmt.Sprintf("expected next token to be %q, got %q instead", string(t), p.peekToken.Literal))
}

func (p *parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return lowest
}

func (p *parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return lowest
}
