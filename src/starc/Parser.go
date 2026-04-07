package main

import (
    "fmt"
    //"strconv"
)

type Parser struct {
    tokens []Token
    current int
}

func (p *Parser) Parse() []Node {
    var statements []Node
    
    for !p.isAtEnd() {
        stmt := p.ParseStmt()
        if stmt != nil {
            statements = append(statements, stmt)
        }
    }
    return statements
}

func (p *Parser) ParseStmt() NodeStmt {
    token := p.peek(0)
    
    if p.isAtEnd() {return nil}
    
    if token.tokenType == IDENTIFIER && p.peek(1).tokenType == EQUAL {return p.assignement()}
    if token.tokenType == IF {return p.ifStmt()}
    if token.tokenType == VAR {return p.varAssignement()}
    if token.tokenType == PRINT {return p.printStmt()}
    if token.tokenType == RETURN {return p.returnStmt()}
    if token.tokenType == WHILE {return p.whileStmt()}
    
    if token.tokenType == SEMICOLON {
        p.advance()
        return nil
    }
    
    expr := p.expression()
    if p.peek(0).tokenType == SEMICOLON {p.advance()}
    return &NodeStmtExpr{Expr: expr}
}

func (p *Parser) peek(offset int) Token {
    return p.tokens[p.current + offset]
}

func (p *Parser) isAtEnd() bool {
    return p.peek(0).tokenType == EOF
}

func (p *Parser) advance() Token {
    if !p.isAtEnd() {p.current++}
    //fmt.Println("[", p.peek(0), "]")
    return p.tokens[p.current - 1]
}

// Next following are the nodes’ recursion

func (p *Parser) primary() NodeExpr {
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
        expr := p.grouping()
        if p.peek(0).tokenType != RIGHT_PAREN && p.isAtEnd() {
            PrintError(5)
        } else {
            p.advance()
            return &NodeGroup{Expression: expr}
        }
    }
    if token.tokenType == NULL {
        p.advance()
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == STRING {
        p.advance()
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    PrintError(5)
    panic("")
}

func (p *Parser) unary() NodeExpr {
    token := p.peek(0)
    
    if token.tokenType == MINUS {
        p.advance()
        
        right := p.unary()
        
        return &NodeUnary{Operator: token.Lexeme, Right: right}
    }
    return p.primary()
}

func (p * Parser) grouping() NodeExpr {
    expr := p.expression()
    p.advance()
    return expr
}

func (p *Parser) factor() NodeExpr {
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

func (p *Parser) binary() NodeExpr {
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

func (p *Parser) comparison() NodeExpr {
    expr := p.binary()
    
    for p.peek(0).tokenType.isBoolOperator() {
        operator := p.advance().Lexeme
        right := p.binary()
        expr = &NodeBinary{Left: expr, Operator: operator, Right: right}
    }
    return expr
}

func (p *Parser) expression() NodeExpr {
	if p.isAtEnd() {return nil}
    return p.comparison()
}

func (p *Parser) varAssignement() NodeStmt {
    var varVal Node = nil
    
    if !p.isAtEnd() {
        p.advance()
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
        //PrintError(6)
        //panic("")
    }
    PrintError(8)
    panic("")
}

func (p *Parser) returnStmt() NodeStmt {
    var val NodeExpr = nil
    
    if !p.isAtEnd() {
		p.advance()
        if p.peek(0).tokenType != SEMICOLON {
            val = p.expression()
        }
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            return &NodeStmtReturn{Value: val}
        }
        
    }
    PrintError(8)
    panic("")
}

func (p *Parser) printStmt() NodeStmt {
    
    if !p.isAtEnd() {
        p.advance()
        var valList []NodeExpr
        for p.peek(0).tokenType != SEMICOLON {
            valList = append(valList, p.expression())
            if p.peek(0).tokenType == COMMA {
                p.advance()
            }
        }
        return &NodeStmtPrint{Expressions : valList}
    }
    PrintError(6)
    panic("")
}

func (p *Parser) blockStmt() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        var stmts []NodeStmt
        for p.peek(0).tokenType != RIGHT_BRACE {
            stmts = append(stmts, p.ParseStmt())
            if p.peek(0).tokenType == SEMICOLON {p.advance()}
        }
        p.advance()
        return &NodeBlock{Instructions: stmts}
    }
    PrintError(6)
    panic("")
}

func (p *Parser) assignement() NodeStmt {
    
    if !p.isAtEnd() {
        varName := p.advance().Lexeme
        p.advance()
        varVal := p.expression()
        if p.peek(0).tokenType != SEMICOLON {PrintError(8); panic("")}
        return &NodeAssignement{Name: varName, Value: varVal}
    }
    PrintError(6)
    panic("")
}

func (p *Parser) ifStmt() NodeStmt {
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        condition = p.comparison()
        if p.peek(0).tokenType != RIGHT_PAREN {PrintError(5); panic("")}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtIf{Condition: condition, Result: result}
    }
    PrintError(5)
    panic("")
}

func (p *Parser) whileStmt() NodeStmt {
    fmt.Println("while commence")
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        fmt.Println("Avant la compa")
        condition = p.comparison()
        fmt.Println("Après")
        if p.peek(0).tokenType != RIGHT_PAREN {PrintError(5); panic("")}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtWhile{Condition: condition, Result: result}
    }
    PrintError(5)
    panic("")
}