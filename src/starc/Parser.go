package main

import (
	"fmt"
)

type Parser struct {
    Package string
    tokens []Token
    current int
    envi *Environnement
}

func (p *Parser) Parse() []Node {
    var statements []Node
    for !p.isAtEnd() {
        stmt := p.ParseStmt()
        if stmt != nil {
            fmt.Println(fmt.Sprintf("New stmt: %T", stmt))
            statements = append(statements, stmt)
            fmt.Println(fmt.Sprintf("Stmt: %T consummed", stmt))
        }
    }
    //fmt.Println("Parsed envi", p.envi.Static)
    return statements
}

func (p *Parser) ParseStmt() NodeStmt {
    token := p.peek(0)
    
    if p.isAtEnd() {return nil}
    fmt.Println(token.tokenType == IDENTIFIER)
    switch token.tokenType {
        case THIS:
        	//fmt.Println("Ident", p.peek(0).Lexeme)
            expr := p.expression()
            //fmt.Println("Ident", expr)
            if p.peek(0).tokenType == EQUAL {
                //fmt.Println("Ident", p.peek(0))
                //fmt.Println("Ident", expr)
                p.advance()
                value := p.expression()
                return &NodeAssignment{Target: expr, Value: value}
            }
            //fmt.Println("L'expression: ", expr)
            return &NodeStmtExpr{Expr: expr}
        case IDENTIFIER:
            //fmt.Println("Ident", p.peek(0).Lexeme)
            expr := p.expression()
            //fmt.Println("Ident", expr)
            if p.peek(0).tokenType == EQUAL {
                //fmt.Println("Ident", p.peek(0))
                //fmt.Println("Ident", expr)
                p.advance()
                value := p.expression()
                return &NodeAssignment{Target: expr, Value: value}
            }    
            //fmt.Println("L'expression: ", expr)
            return &NodeStmtExpr{Expr: expr}
        
        case IF: return p.ifStmt()
        case VAR: return p.varAssignment()
        case CONST: return p.constAssignment()
        case PRINT: return p.printStmt()
        case RETURN: return p.returnStmt()
        case LOOP: return p.loop()
        case WHILE: return p.whileStmt()
        case FUNC: return p.funcInit()
        case TYPEDEF: return p.typeDef()
        case CLASS: return p.classDef()
        case CCALL: return p.ccall()
        case PRIVATE: return p.scopeAccess()
        case STATIC: return p.staticStmt()
        case PACKAGE: return p.packageDef()
        case IMPORT: return p.importPkg()
    }
    
    
    if token.tokenType == SEMICOLON {
        p.advance()
        return nil
    }
    
    fmt.Println("Expr quelconque")
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
    //fmt.Println("Advanced, ", p.peek(0), ", consumed: ", p.peek(-1))
    return p.tokens[p.current - 1]
}

func (p *Parser) isValidType(Type Token) bool {
    //fmt.Println("Checkin' type")
    return Type.tokenType.isType() || p.envi.hasType(Type.Lexeme)
}

// Next following are the nodes’ recursion

func (p *Parser) primary() NodeExpr {
    //fmt.Println("Running through primary")
    token := p.advance()
    
    if token.tokenType.isDigit() {
        return &NodeLiteral{Value: token.Lexeme}
    }
    if token.tokenType == IDENTIFIER || token.tokenType == THIS {
        var expr NodeExpr
        
        if token.tokenType == THIS {
        	expr = &NodeVariable{Name: "this"}
        } else if p.peek(0).tokenType == LEFT_PAREN {
            expr = p.funcCall()
        } else if p.peek(0).tokenType == SCOPE_RESOLVE {
            p.advance()
            expr = &NodePkgResolve{Pkg: token.Lexeme, Resolution: p.expression()}
        } else {
            expr = &NodeVariable{Name: token.Lexeme}
        }
        
        var class string
        if p.envi.hasType(p.envi.Variable[p.peek(-1).Lexeme]) {class = p.envi.Variable[p.peek(-1).Lexeme]}
        for p.peek(0).tokenType == DOT {
            
            object := token.Lexeme
            if object == "this" {p.envi.Pointer[object] = true}
            p.advance()
            field := p.advance().Lexeme
            var objName Token = p.peek(-3)
            
            if p.peek(0).tokenType == LEFT_PAREN {
                p.advance()
                var argsList []NodeExpr
                
                for p.peek(0).tokenType != RIGHT_PAREN {
                    argsList = append(argsList, p.expression())
                    if p.peek(0).tokenType == COMMA {p.advance()}
                }
                p.advance()
                if field == "new" {
                    if !p.isValidType(objName) {p.envi.Unknown.Type[p.peek(0).Lexeme] = p.Package}
                //    class = objName.Lexeme
                }
                if class != "" {class += "_"}
                isStatic := p.envi.getStatic(field, object)
                //fmt.Println("Static ?", field, isStatic, p.envi.Static[field])
                expr = &NodeExprMethodCall{Class: class + objName.Lexeme, Parent: expr, Name: field, Args: argsList, Static: isStatic}
                class = p.envi.Func[objName.Lexeme + "_" + field]
            } else {
                symbol := "."
                if p.envi.Pointer[object] == true {symbol = "->"}
                expr = &NodeGet{Object: expr, Symbol: symbol, Field: field}
            }    
        }
        fmt.Println(fmt.Sprintf("Expr %v %T", expr, expr))
        return expr
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
    
    if token.tokenType == DOLLAR {
		if p.peek(0).tokenType == ALLOCATE {
            p.advance()
            if p.advance().tokenType != COLON {PrintError(12, "Missing ':' after $call")}
            if !p.peek(0).tokenType.isMemManage() {PrintError(12, "Unknown alloc caller")}
            call := p.advance().Lexeme
            if p.peek(0).tokenType != RIGHT_ARROW {PrintError(12, "Missing arrow '->' that point towards the allocation size")}
            p.advance()
            size := p.expression()
			return &NodeExprAlloc{Allocation: call, Size: size}
		}
    	PrintError(12, "Unknown or empty $call")
	}
    
    if token.tokenType == NULL {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == TRUE {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == FALSE {
        return &NodeLiteral{Value: token.Lexeme}
    }
    
    if token.tokenType == STRING {
        return &NodeLiteral{Value: token.Lexeme}
    }
    //fmt.Println("Token: ", p.peek(0))
    PrintError(5, "May be due to unknown character " + token.Lexeme)
    panic("")
}

func (p *Parser) concat() NodeExpr {
    //fmt.Println("Running through concat")
    expr := p.primary()
    
    for p.peek(0).tokenType == CONCAT {
        p.advance()
        right := p.primary()
        expr = &NodeExprConcat{From: right, To: expr}
    }
    return expr
}

func (p *Parser) unary() NodeExpr {
    //fmt.Println("Running through unary")
    token := p.peek(0)
    
    if token.tokenType == MINUS || token.tokenType == BANG {
        p.advance()
        
        right := p.unary()
        
        return &NodeUnary{Operator: token.Lexeme, Right: right}
    }
    return p.concat()
}

func (p * Parser) grouping() NodeExpr {
    //fmt.Println("Running through group")
    expr := p.expression()
    p.advance()
    return expr
}

func (p *Parser) factor() NodeExpr {
    //fmt.Println("Running through factors")
    expr := p.unary()
    token := p.peek(0)
    
    if !p.isAtEnd() {
        if token.tokenType == STAR || token.tokenType == SLASH {
            p.advance()
            operator := token.Lexeme
            right := p.unary()
            expr = &NodeBinary{Left: expr, Operator: operator, Right: right}
        }
    }
    return expr
}

func (p *Parser) binary() NodeExpr {
    //fmt.Println("Running through binary")
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
    //fmt.Println("Running through comp")
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
    //fmt.Println("Running through expression")
    return p.comparison()
}

func (p *Parser) properties() []string {
    if !p.isAtEnd() {
        var propertyList []string
        p.advance()
        for p.peek(0).tokenType != GREATER {
            if !p.peek(0).tokenType.isProperty() {PrintError(6, "Expected property call")}
            propertyList = append(propertyList, p.peek(0).Lexeme)
            p.advance()
            if p.peek(0).tokenType == COMMA {p.advance()}
        }
        return propertyList
    }
    PrintError(8, "Property call declared but with no end")
    panic("")
}

func (p *Parser) varAssignment() NodeStmt {
    var varVal NodeExpr = nil
    var propertyList []string = nil
    
    if !p.isAtEnd() {
        p.advance()
        
        if p.peek(0).tokenType == LESS {propertyList = p.properties(); p.advance()}
        if !p.isValidType(p.peek(0)) {PrintError(5, "Unsuported form of type for now... to fix"); panic("")}
        varType := p.advance()
        
        if p.peek(0).tokenType == STAR {p.envi.Pointer[p.peek(1).Lexeme] = true; p.advance()}
        if p.peek(0).tokenType != IDENTIFIER { PrintError(3, "Expected an identifier for function name"); panic("") }
        varName := p.advance().Lexeme
        //fmt.Println("Var name", varName)
        if p.envi.hasVar(varName) {PrintError(10, "Var declared twice"); panic("")}
        //fmt.Println("Ur var in envi", p.envi.Variable[varName])
        
		for _, prop := range propertyList {
            propFunc := varName + "_" + prop
			p.envi.Func[propFunc] = varType.Lexeme
			//fmt.Println("Key:", propFunc, "Value:", p.envi.Func[propFunc])
		}
        
        if p.peek(0).tokenType == EQUAL {
            p.advance()
            varVal = p.expression()
        }
        
        global := false
        if p.envi.isGlobal() {global = true}
        
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            p.envi.Variable[varName] = varType.Lexeme
            return &NodeStmtVar{Name: varName, Properties: propertyList, Type: varType, Value: varVal, Global: global}
        }
    }
    PrintError(8, "Reached End Of File in an invalid variable declaration")
    panic("")
}

func (p *Parser) constAssignment() NodeStmt {
    var constVal NodeExpr = nil
    var propertyList []string = nil
    
    if !p.isAtEnd() {
        p.advance()
        
        if p.peek(0).tokenType == LESS {propertyList = p.properties(); p.advance()}
        if !p.isValidType(p.peek(0)) { PrintError(7, "Unknown const type, if new class or type you may want to know if it's in the current scope" + p.peek(0).Lexeme); panic("") }
        constType := p.advance()
        
        if p.peek(0).tokenType == STAR {p.envi.Pointer[p.peek(1).Lexeme] = true; p.advance()}
        if p.peek(0).tokenType != IDENTIFIER { PrintError(3, "Expected an identifier for function name"); panic("") }
        constName := p.advance().Lexeme
        if p.envi.hasConst(constName) {PrintError(10, "Const declared twice"); panic("")}
        
		for _, prop := range propertyList {
            propFunc := constName + "_" + prop
			p.envi.Func[propFunc] = constType.Lexeme
		}
        
        if p.peek(0).tokenType != EQUAL {
            PrintError(6, "Invalid Const Stmt")
        }
        p.advance()
        constVal = p.expression()
        
        global := false
        if p.envi.isGlobal() {global = true}
        
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            p.envi.Const[constName] = constType.Lexeme
            return &NodeStmtConst{Name: constName, Properties: propertyList, Type: constType, Value: constVal, Global: global}
        }
    }
    PrintError(8, "Reached End Of File in an invalid const declaration")
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

func (p *Parser) ccall() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        var callerName string
        var called []NodeExpr
        if !p.peek(0).tokenType.isAction() {fmt.Println(p.peek(0)); PrintError(6, "Expected action for the C-Caller"); panic("")}
        actionName := p.advance().Lexeme
        switch actionName {
        case "include":
            for p.peek(0).tokenType != SEMICOLON {
                called = append(called, p.expression())
                if p.peek(0).tokenType == COMMA {p.advance()}
            }
            break
        case "function":
            if p.peek(0).tokenType != IDENTIFIER {PrintError(3, "Expected Identifier"); panic("")}
            callerName = p.advance().Lexeme
            //fmt.Println("Caller", callerName)
            if p.peek(0).tokenType != LEFT_PAREN {PrintError(3, "Expected left parenthesize for C function call"); panic("")}
            p.advance()
            for p.peek(0).tokenType != RIGHT_PAREN {
                called = append(called, p.expression())
                if p.peek(0).tokenType == COMMA {p.advance()}
            }
            p.advance()
            break
        }
        //fmt.Println("Token actuel", p.peek(0))
        if p.peek(0).tokenType != SEMICOLON {fmt.Println(p.peek(0)); PrintError(3, "Expected semicolon"); panic("")}
        p.advance()
        return &NodeStmtC{Action: actionName, Called: called, CallerName: callerName}
    }
    PrintError(8, "End Of File reached in C call")
    panic("")
}

func (p *Parser) blockStmt() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        parScope := p.envi
        p.envi = p.envi.NewScope()
        var stmts []NodeStmt
        for p.peek(0).tokenType != RIGHT_BRACE {
            stmts = append(stmts, p.ParseStmt())
            if p.peek(0).tokenType == SEMICOLON {p.advance()}
        }
        p.advance()
        p.envi = parScope
        return &NodeBlock{Instructions: stmts}
    }
    PrintError(6, "Missing }")
    panic("")
}

func (p *Parser) assignement() NodeStmt {
    
    if !p.isAtEnd() {
        target := p.expression()
        p.advance()
        value := p.expression()
        if p.peek(0).tokenType != SEMICOLON {PrintError(8, "Missing semi-colon ;"); panic("")}
        return &NodeAssignment{Target: target, Value: value}
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

func (p *Parser) loop() NodeStmt {
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        p.advance()
        looping := p.expression()
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtLoop{Looping: looping, Result: result}
    }
    PrintError(5, "Invalid loop stmt")
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
        if p.peek(0).tokenType == VAR_ARGS && p.peek(2).tokenType == RIGHT_PAREN {paramType.Lexeme = paramType.Lexeme + "..."; p.advance()}
        if p.peek(0).tokenType != IDENTIFIER {fmt.Println(p.peek(0)); PrintError(3, "Expected identifier"); panic("") }
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
        returnType := p.advance()
        if !p.isValidType(returnType) && p.peek(-1).tokenType != VOID {PrintError(7, fmt.Sprintf("Unknown return type <%s>", returnType.Lexeme)); panic("")}
        if p.peek(0).tokenType != IDENTIFIER {PrintError(5, "Expected identifier"); panic("")}
        funcName := p.advance().Lexeme
        p.envi.Func[funcName] = returnType.Lexeme
        if p.peek(1).tokenType != RIGHT_PAREN {
            for p.peek(0).tokenType != RIGHT_PAREN {
                paramList = append(paramList, p.parseParam())
            }
            p.advance()
        } else {p.advance(); p.advance()}
        if p.peek(0).tokenType != LEFT_BRACE {PrintError(6, "Missing left brace"); panic("")}
        code := p.blockStmt()
        if funcName == "new" {
            return &NodeStmtConstructor{Return: returnType.Lexeme, Param: paramList, Code: code}
            }
        return &NodeStmtFuncInit{Return: returnType.Lexeme, Name: funcName, Param: paramList, Code: code}
    }
    PrintError(8, "Invalid function def statement")
    panic("")
}

func (p *Parser) funcCall() NodeExpr {
    var argsList []NodeExpr
    if !p.isAtEnd() {
        funcName := p.peek(-1).Lexeme
        if !p.envi.hasFunc(funcName) {p.envi.Unknown.Func[funcName] = p.peek(-3).Lexeme}
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

func (p *Parser) staticStmt() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        stmt := p.ParseStmt()
        if stmt, ok := stmt.(*NodeStmtFuncInit); ok {
            p.envi.addStatic(stmt.Name)
            fmt.Println("Static func: ", stmt.Name, p.envi.Static[stmt.Name])
        }
        return &NodeStaticStmt{Stmt: stmt}
    }
    PrintError(6, "Errounous static call")
    panic("")
}

func (p *Parser) scopeAccess() NodeStmt {
    if !p.isAtEnd() {
        modifier := p.advance()
        //fmt.Println("Current point", p.peek(0))
        stmt := p.ParseStmt()
        //fmt.Println(stmt)
        return &NodeScopeAcces{Modifier: modifier, Stmt: stmt}
    }
    PrintError(5, "Scope access in End Of File")
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
        p.advance()
        className := p.advance().Lexeme
        if p.envi.hasType(className) {PrintError(10, "Class declared twice in the same context"); panic("")}
        p.envi.Type[className] = ""
        if p.advance().tokenType != LEFT_BRACE {PrintError(8, "Expected left brace {"); panic("")}
        parScope := p.envi
        p.envi = p.envi.NewScope()
        envi := p.envi
        var code []NodeStmt
        for p.peek(0).tokenType != RIGHT_BRACE {
            code = append(code, p.ParseStmt())
        }
        p.advance()
        p.envi = parScope
        p.envi.Classes[className] = envi
        return &NodeStmtClass{Name: className, Code: code}
    }
    PrintError(8, "Reached precocious End Of File in class body")
    panic("")
}

func (p *Parser) packageDef() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        name := p.advance().Lexeme
        p.Package = name
        if p.advance().tokenType != SEMICOLON {PrintError(8, "Missing Semicolon"); panic("")}
        return &NodeStmtPkg{Name: name}
    }
    PrintError(5, "Invalid Package Definition")
    panic("")
}

func (p *Parser) importPkg() NodeStmt {
    if !p.isAtEnd() {
        var pathes []string
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            p.advance()
            for p.peek(0).tokenType != RIGHT_BRACE {
                pathes = append(pathes, "\"" + p.peek(0).Lexeme[1:len(p.peek(0).Lexeme)-1] + ".h" + "\"")
                p.advance()
                if p.peek(0).tokenType == COMMA {p.advance()}
            }
        } else if p.peek(0).tokenType == STRING {
            pathes = append(pathes, "\"" + p.peek(0).Lexeme[1:len(p.peek(0).Lexeme)-1] + ".h" + "\"")
            p.advance()
            if p.peek(0).tokenType != SEMICOLON {PrintError(5, "Something")}
        }
        p.advance()
        return &NodeImport{Names: pathes}
    }
    PrintError(5, "Invalid Package Definition")
    panic("")
}