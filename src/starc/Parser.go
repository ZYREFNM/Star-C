package main

import (
    //"fmt"
    //"strconv"
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
    return p.tokens[p.current - 1]
}

// Next following are the nodes’ recursion

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
        expr := p.expression()
        if token.tokenType != RIGHT_PAREN && p.isAtEnd() {
            PrintError(5)
        } else {
            return &NodeGroup{Expression: expr}
        }
    }
    PrintError(5)
    panic("")
}

func (p *Parser) unary() Node {
    token := p.peek(0)
    
    if token.tokenType == MINUS {
        p.advance()
        
        right := p.unary()
        
        return &NodeUnary{Operator: token.Lexeme, Right: right}
    }
    return p.primary()
}

func (p * Parser) grouping() Node {
    expr := p.expression()
    return expr
}

func (p *Parser) factor() Node {
    expr := p.unary()
    token := p.peek(0)
    
    if !p.isAtEnd() {
        if token.tokenType == STAR || token.tokenType == SLASH {
            p.advance()
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
            p.advance()
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

func (p *Parser) assignement() Node {
    token := p.peek(0)
    var varVal Node
    
    if !p.isAtEnd() {
    	if token.tokenType == VAR {
            p.advance()
            varType := token.Lexeme
            p.advance()
            varName := token.Lexeme
            if p.peek(-1).tokenType.isType() {
            	p.advance()
                if p.peek(1).tokenType == EQUAL {
                    p.advance()
                    varVal = p.expression()
                } else if p.peek(1).tokenType == SEMICOLON {
                    varVal = nil
                } else { PrintError(6) }
                return &NodeStmtVar{Name: varName, Type: varType, Value: varVal}
            } else {
                PrintError(6)
            }
        }
	}
    PrintError(6)
    panic("")
}