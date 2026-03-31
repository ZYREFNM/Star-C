package main

import (
    //"fmt"
)

type Parser struct {
    tokens []Token
    current int
}

func (p *Parser) Parse() Expr {
    return p.expression()
}

func (p *Parser) match(types ...TokenType) bool {
    for _, Type := range types {
        if p.check(Type) {
            p.advance()
            return true
        }
    }
    return false
}

func (p *Parser) check(Type TokenType) bool {
    if p.isAtEnd() {
        return false
    }
    return p.peek().tokenType == Type
}

func (p *Parser) advance() Token {
    if !p.isAtEnd() {
    	p.current++
    }
    
    return p.previous()
}

func (p *Parser) isAtEnd() bool {
    return p.peek().tokenType == EOF
}

func (p *Parser) peek() Token {
    return p.tokens[p.current]
}
func (p *Parser) previous() Token {
    return p.tokens[p.current - 1]
}

func (p *Parser) expression() Expr {
    return p.equality()
}

func (p *Parser) equality() Expr {
    var expr Expr = p.comparison()
    
    for p.match(BANG_EQUAL,  EQUAL_EQUAL){
        var operator Token = p.previous()
        var right Expr = p.comparison()
        expr = &Binary{
            expr,
            operator,
            right,
        }
    }
    return expr
}

func (p *Parser) comparison() Expr {
    var expr Expr = p.term()
    
    for p.match(GREATER, GREATER_EQUAL,  LESS, LESS_EQUAL) {
        var operator Token = p.previous()
        var right Expr = p.term()
        expr = &Binary{
            expr,
            operator,
            right,
        }
    }
    return expr
}

func (p *Parser) term() Expr {
    var expr Expr = p.factor()
    
    for p.match(MINUS, PLUS) {
        var operator Token = p.previous()
        var right Expr = p.factor()
        expr = &Binary{
            expr,
            operator,
            right,
        }
    }
    return expr
}

func (p *Parser) factor() Expr {
    var expr Expr = p.unary()
    
    for p.match(SLASH, STAR) {
        var operator Token = p.previous()
        var right Expr = p.unary()
        expr = &Binary{
            expr,
            operator,
            right,
        }
    }
    return expr
}

func (p *Parser) unary() Expr {
    if p.match(BANG, MINUS) {
        var operator Token = p.previous()
        var right Expr = p.unary()
        return &Unary{
            operator,
            right,
        }
    }
    return p.primary()
}

func (p *Parser) primary() Expr {
    if p.match(FALSE) {
        return &Literal{false}
    }
    if p.match(TRUE) {
        return &Literal{true}
    }
    if p.match(NULL) {
        return &Literal{nil}
    }
    
    if p.match(INT, INT8, INT16, INT32, INT64, UINT, UINT8, UINT16, UINT32, UINT64, FLOAT, FLOAT8, FLOAT16, FLOAT32, FLOAT64, UFLOAT, UFLOAT8, UFLOAT16, UFLOAT32, UFLOAT64, STRING) {
        return &Literal{p.previous().Literal}
    }
    if p.match(LEFT_PAREN) {
        var expr Expr = p.expression()
        p.consume(RIGHT_PAREN, "Expect ')' after expression")
        return &Grouping{expr}
    }
   PrintError(5)
   panic("")
}

func (p *Parser) consume(Type TokenType,  message string) Token {
    if p.check(Type) {
        return p.advance()
    }
    PrintError(5)
    panic("")
}

func (p *Parser) synchronize() {
    p.advance()
    
    for !p.isAtEnd() {
        if p.previous().tokenType == SEMICOLON {
            return
        }
        switch p.peek().tokenType {
            case PACKAGE:
            case CLASS:
            case INNER:
            case FUNC:
            case VAR:
            case CONST:
            case FOR:
            case GOTO:
            case IF:
            case WHILE:
            case PRINT:
            case SWITCH:
            case RETURN:
            	return
        }
        p.advance()
    }
}