package main

import (
    //"fmt"
)

type Parser struct {
    tokens []Token
    current int
}

func (p *Parser) Parse() {
    
}

func (p *Parser) peek(offset int) Token {
    return p.tokens[p.current + offset]
}

func (p *Parser) isAtEnd() bool {
    return p.peek(0).tokenType == EOF
}

func (p *Parser) advance() Token {
    if !p.isAtEnd() {p.current++}
    return p.peek(-1)
}

// Next following are the nodes’ recursion

func (p * Parser) grouping() Node {
    expr := p.expression()
    return expr
}

func (p *Parser) primary() Node {
    token := p.peek(0)
    
    if token.tokenType.isDigit() {
        p.advance()
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == IDENTIFIER {
        p.advance()
        return &NodeVariable{Name: token.Lexeme}
    }
    if token.tokenType == LEFT_PAREN {
        p.advance()
        if token.tokenType != RIGHT_PAREN && p.isAtEnd() {
            PrintError(5)
        }
    }
    PrintError(5)
    panic("")
}

func (p *Parser) factor() Node {
    
    token := p.peek(0)
    expr := p.primary()
    if !p.isAtEnd() {
        if token.tokenType == STAR || token.tokenType == SLASH {
            operator := token.Lexeme
            right := p.primary()
            expr = &NodeBinary{Left: expr, Operator: operator, Right: right}
        }
    }
    return expr
}

func (p *Parser) binary() Node {
    expr := p.factor()
    token := p.peek(0)
    
    if !p.isAtEnd() {
        if token.tokenType == PLUS || token.tokenType == MINUS {
            operator := token.Lexeme
            right := p.factor()
            expr = &NodeBinary{Left: expr, Operator: operator, Right: right}
        }
    }
    return expr
}

func (p *Parser) expression() Node {
    return p.binary()
}