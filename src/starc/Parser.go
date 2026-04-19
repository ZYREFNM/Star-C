package main

import (
    "fmt"
)

type Parser struct {
    tokens []Token
    current int
    envi Environnement
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
    if token.tokenType == VAR {return p.varAssignment()}
    if token.tokenType == PRINT {return p.printStmt()}
    if token.tokenType == RETURN {return p.returnStmt()}
    if token.tokenType == WHILE {return p.whileStmt()}
    if token.tokenType == FUNC {return p.funcInit()}
    if token.tokenType == TYPEDEF {return p.typeDef()}
    if token.tokenType == CLASS {return p.classDef()}
    
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
    return p.tokens[p.current - 1]
}

func (p *Parser) isValidType(Type Token) bool {
    return Type.tokenType.isType() || p.envi.hasType(Type.Lexeme)
}

// Next following are the nodes’ recursion

func (p *Parser) primary() NodeExpr {
    token := p.advance()
    
    if token.tokenType.isDigit() {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == IDENTIFIER {
        if p.peek(0).tokenType == LEFT_PAREN {return p.funcCall()}
        return &NodeVariable{Name: token.Lexeme}
    }
    if token.tokenType == LEFT_PAREN {
        expr := p.grouping()
        if p.peek(0).tokenType != RIGHT_PAREN && p.isAtEnd() {
            PrintError(8, "Expected ) before End Of File")
        } else {
            p.advance()
            return &NodeGroup{Expression: expr}
        }
    }
    if token.tokenType == NULL {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == STRING {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    PrintError(5, "May be due to unknown characters")
    panic("")
}

func (p *Parser) concat() NodeExpr {
    expr := p.primary()
    
    for p.peek(0).tokenType == CONCAT {
        p.advance()
        right := p.primary()
        expr = &NodeExprConcat{From: right, To: expr}
    }
    return expr
}

func (p *Parser) unary() NodeExpr {
    token := p.peek(0)
    
    if token.tokenType == MINUS {
        p.advance()
        
        right := p.unary()
        
        return &NodeUnary{Operator: token.Lexeme, Right: right}
    }
    return p.concat()
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

func (p *Parser) varAssignment() NodeStmt {
    var varVal Node = nil
    
    if !p.isAtEnd() {
        p.advance()
        if !p.peek(0).tokenType.isType() && !p.envi.hasType(p.peek(0).Lexeme) { PrintError(7, "Unknown variable type, if new class or type you may want to know if it's in the current scope"); panic("") }
        
        varType := p.advance()
        if p.peek(0).tokenType != IDENTIFIER { PrintError(3, "Expected an identifier for function name"); panic("") }
        varName := p.advance().Lexeme
        if p.envi.hasVar(varName) {PrintError(10, "Var declared twice"); panic("")}
        
        if p.peek(0).tokenType == EQUAL {
            p.advance()
            varVal = p.expression()
        }
        
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            p.envi.Variable[varName] = varVal
            return &NodeStmtVar{Name: varName, Type: varType, Value: varVal}
        }
    }
    PrintError(8, "Reached End Of File in an invalid variable declaration")
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
    PrintError(8, "Missing semi-colon ';'")
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
    PrintError(6, "Unknown litteral or EOF to print")
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
    PrintError(6, "Missing }")
    panic("")
}

func (p *Parser) assignement() NodeStmt {
    
    if !p.isAtEnd() {
        varName := p.advance().Lexeme
        p.advance()
        varVal := p.expression()
        if p.peek(0).tokenType != SEMICOLON {PrintError(8, "Missing semi-colon ;"); panic("")}
        return &NodeAssignement{Name: varName, Value: varVal}
    }
    PrintError(6, "Variable assignment reached End Of File")
    panic("")
}

func (p *Parser) ifStmt() NodeStmt {
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        condition = p.comparison()
        if p.peek(0).tokenType != RIGHT_PAREN {PrintError(5, "Expected right parenthesize"); panic("")}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtIf{Condition: condition, Result: result}
    }
    PrintError(5, "Invalid if statement")
    panic("")
}

func (p *Parser) whileStmt() NodeStmt {
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        condition = p.comparison()
        if p.peek(0).tokenType != RIGHT_PAREN {PrintError(5, "Expected right parenthesize )"); panic("")}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtWhile{Condition: condition, Result: result}
    }
    PrintError(5, "Invalid while statement")
    panic("")
}

func (p *Parser) parseParam() NodeStmt {
    var paramVal Node = nil
    
    if !p.isAtEnd() {
        p.advance()
        if !p.peek(0).tokenType.isType() { PrintError(7, "Unknown type, you may check if you typedefed that type or created that new class"); panic("") }
        paramType := p.advance()
        if p.peek(0).tokenType != IDENTIFIER { PrintError(3, "Expected identifier"); panic("") }
        paramName := p.advance().Lexeme
        
        if p.peek(0).tokenType == COMMA {
            return &NodeStmtVar{Name: paramName, Type: paramType, Value: paramVal}
        }
        
        if p.peek(0).tokenType == RIGHT_PAREN {
            return &NodeStmtVar{Name: paramName, Type: paramType, Value: paramVal}
        }
    }
    PrintError(8, "Missing right parenthesize )")
    panic("")
}

func (p *Parser) funcInit() NodeStmt {
    var paramList []NodeStmt = nil
    if !p.isAtEnd() {
        p.advance()
        fmt.Println(p.peek(0))
        returnType := p.advance()
        fmt.Println(returnType.Lexeme, " is returned")
        if !p.isValidType(returnType) && p.peek(-1).tokenType != VOID {PrintError(7, fmt.Sprintf("Unknown return type <%s>", returnType.Lexeme)); panic("")}
        fmt.Println(p.peek(0))
        if p.peek(0).tokenType != IDENTIFIER {PrintError(5, "Expected identifier"); panic("")}
        funcName := p.advance().Lexeme
        fmt.Println(funcName)
        if p.peek(1).tokenType != RIGHT_PAREN {
            for p.peek(0).tokenType != RIGHT_PAREN {
                paramList = append(paramList, p.parseParam())
            }
            p.advance()
        } else {p.advance(); p.advance()}
        if p.peek(0).tokenType != LEFT_BRACE {PrintError(6, "Missing left brace"); panic("")}
        code := p.blockStmt()
        return &NodeStmtFuncInit{Return: returnType.Lexeme, Name: funcName, Param: paramList, Code: code}
    }
    PrintError(8, "Invalid function def statement")
    panic("")
}

func (p *Parser) funcCall() NodeExpr {
    var argsList []NodeExpr
    if !p.isAtEnd() {
        funcName := p.peek(-1).Lexeme
        p.advance()
        for p.peek(0).tokenType != RIGHT_PAREN {
            arg := p.expression()
            argsList = append(argsList, arg)
        }
        p.advance()
        return &NodeExprFuncCall{Name: funcName, Args: argsList}
    }
    PrintError(5, "Invalid function call")
    panic("")
}

func (p *Parser) typeDef() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        var typeData Token = p.advance()
        var typeName string = p.advance().Lexeme
        if p.envi.hasType(typeName) {PrintError(10, "Object or type already exist in current context"); panic("")}
        p.envi.Type[typeName] = ""
        var typeVars []NodeStmt = nil
        if typeData.tokenType == STRUCT {
            p.advance()
            for p.peek(0).tokenType != RIGHT_BRACE {
                typeVars = append(typeVars, p.varAssignment())
            }
            p.advance()
        } else {
            if p.peek(0).tokenType != SEMICOLON {PrintError(5, "Missing semi-colon ;"); panic("")}
        }
        return &NodeStmtTypeDef{Type: typeData, Name: typeName, Vars: typeVars}
    }
    PrintError(5, "Invalid typedef, not ended correctly before EOF")
    panic("")
}

/*func (p *Parser) classCode() NodeStmt {
    
}*/

func (p *Parser) classDef() NodeStmt {
    if !p.isAtEnd() {
        var classVars []NodeStmt = nil
        var classFunc []NodeStmt = nil
        var classTypes []NodeStmt = nil
        p.advance()
        className := p.advance().Lexeme
        if p.envi.hasType(className) {PrintError(10, "Class declared twice in the same context"); panic("")}
        if p.advance().tokenType != LEFT_BRACE {PrintError(8, "Expected left brace {"); panic("")}
        for p.peek(0).tokenType != RIGHT_BRACE {
            if p.peek(0).tokenType == VAR {classVars = append(classVars, p.varAssignment())}
            if p.peek(0).tokenType == FUNC {classFunc = append(classFunc, p.funcInit())}
            if p.peek(0).tokenType == TYPEDEF {classTypes = append(classTypes, p.typeDef())}
        }
        p.advance()
        fmt.Println("Fim ", p.peek(0))
        return &NodeStmtClass{Name: className, Vars: classVars, Func: classFunc, TypeDef: classTypes}
    }
    PrintError(8, "Reached precocious End Of File in class body")
    panic("")
}