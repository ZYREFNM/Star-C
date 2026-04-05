package main

import (
    "fmt"
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
    fmt.Print(p.current, p.current - 1)
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
        if p.peek(0).tokenType != RIGHT_PAREN && p.isAtEnd() {
            PrintError(5)
        } else {
            return &NodeGroup{Expression: expr}
        }
    }
    PrintError(7)
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

func (p *Parser) varAssignement() Node {
    var varVal Node = nil
    
    if !p.isAtEnd() {
        if p.peek(0).tokenType != VAR { PrintError(6); panic("") }
        if !p.peek(0).tokenType.isType() { PrintError(7); panic("") }
        varType := p.advance()
        if p.peek(0).tokenType != IDENTIFIER { PrintError(3); panic("") }
        varName := p.advance().Lexeme
        
        if p.peek(0).tokenType == EQUAL {
            p.advance()
            varVal = p.expression()
        }
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            return &NodeStmtVar{Name: varName, Type: varType, Value: varVal}
        }
        PrintError(6)
        panic("")
    }
    PrintError(6)
    panic("")
}