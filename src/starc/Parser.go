package main

import (
	"fmt"
    "strings"
)

type Parser struct {
    Module string
    filepath string
    tokens []Token
    current int
    envi *Environnement
    CompTime *PreCompiler
}

func (p *Parser) Parse() []Node {
    var statements []Node
    fmt.Println("Module:", p.Module, "line:", p.tokens[p.current].Line)
    for !p.isAtEnd() {
        stmt := p.ParseStmt()
        if stmt != nil {
            statements = append(statements, stmt)
            //fmt.Println(fmt.Sprintf("Stmt: %T consummed", stmt))
        }
    }
    //fmt.Println("Parsed envi", p.envi.Static)
    return statements
}

func (p *Parser) ParseStmt() NodeStmt {
    token := p.peek(0)
    
    if p.isAtEnd() {return nil}
    switch token.tokenType {
        case IDENTIFIER, THIS, GET, STAR:
            expr := p.expression()
            if p.peek(0).tokenType.equals(EQUAL, PLUS_EQUAL, MINUS_EQUAL, STAR_EQUAL, SLASH_EQUAL) {
                symbol := p.advance()
                value := p.expression()
                return &NodeAssignment{Target: expr, Value: value, Symbol: symbol}
            }
            return &NodeStmtExpr{Expr: expr}
        case IF: return p.ifStmt()
        case VAR:
            return p.parseVar()
        case CONST: return p.constAssignment()
        case RETURN: return p.returnStmt()
        case LOOP: return p.loop()
        case WHILE: return p.whileStmt()
        case BREAK: return p.breakStmt()
        case FUNC: return p.funcInit()
        case TYPEDEF: return p.typeDef()
        case CCALL:
            if !p.peek(1).tokenType.isStmtCall() {break}
            fmt.Println(Red, "Calling a stmt", Reset)
            return p.ccall()
        case REFLEX:
            p.CompTime.ParsePreCompile()
            return nil
        case SET: return p.PropertyAccessorStmt()
        case PRIVATE, PUBLIC: return p.scopeAccess()
        case STATIC: return p.staticStmt()
        case MODULE: return p.ModuleDef()
        case IMPORT: return p.importPkg()
    }
    
    
    if token.tokenType == SEMICOLON {
        p.advance()
        return nil
    }
    
    //fmt.Println("Expr quelconque", p.peek(0))
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

func (p *Parser) expect(err string, tok Token, expected ...TokenType) Token {
    for _, e := range expected {
        if e == tok.tokenType {
            p.advance()
            return tok
        }
    }
    CatchError(err)
    return p.peek(0)
}

func (p *Parser) parseType() NodeExpr {
    if !p.isAtEnd() {
        
        var parsedType Token = p.advance()
        if p.peek(0).tokenType.equals(STAR) {parsedType.Lexeme += p.advance().Lexeme}
        if parsedType.tokenType.equals(LEFT_BRACKET) || p.peek(0).tokenType.equals(MAP_TYPE) {
            var subType []NodeExpr
            var size NodeExpr
            switch parsedType.tokenType {
                case LEFT_BRACKET:
                    size = p.expression()
                    p.advance()
                    subType = append(subType, p.parseType())
                    //p.advance()
                case MAP_TYPE:
                    p.advance()
                    subType = append(subType, p.parseType())
                    p.advance()
                    subType = append(subType, p.parseType())
                default:
            }
            return &NodeType{Type: parsedType.Lexeme, Size: size, SubType: subType}
        }
        
        return &NodeType{Type: parsedType.Lexeme}
    }
    CatchError("Incomplete or Unknown Type")
    return nil
}

func (p *Parser) parseStruct() []NodeStmt {
    var structVars []NodeStmt
    p.advance()
    parScope := p.envi 
    p.envi = p.envi.NewScope()
    for p.peek(0).tokenType != RIGHT_BRACE {
        if !p.peek(0).tokenType.equals(PRIVATE, PUBLIC) {
            structVars = append(structVars, p.parseVar())
        } else {
            structVars = append(structVars, p.scopeAccess())
        }
        if p.peek(0).tokenType.equals(SEMICOLON) {p.advance()}
    }
    p.envi = parScope
    p.advance()
    return structVars
}

func (p *Parser) parseEnum() []NodeStmt {
    //fmt.Println("This case")
    var enumVars []NodeStmt
    p.advance()
    parScope := p.envi
    p.envi = p.envi.NewScope()
    for p.peek(0).tokenType != RIGHT_BRACE {
        enumVars = append(enumVars, &NodeStmtExpr{Expr: &NodeVariable{Name: p.advance().Lexeme}})
        p.envi.Table[p.peek(-1).Lexeme] = &Symbol{
            Name: p.peek(-1).Lexeme,
            Kind: S_CONST,
            Module: p.Module,
            Type: &NodeType{Type: "int"},
        }
        if p.peek(0).tokenType.equals(COMMA) {p.advance()}
    }
    p.envi = parScope
    p.advance()
    return enumVars
}

func (p *Parser) primary() NodeExpr {
    token := p.advance()
    
    if token.tokenType.isDigit() {
        return &NodeLiteral{Value: token.Lexeme, Type: token}
    }
    if token.tokenType.equals(IDENTIFIER, THIS, GET, SET, STAR) {
        var expr NodeExpr
        
        if token.tokenType.equals(STAR) {
            token = p.advance()
            token.Lexeme = "*" + token.Lexeme
        }
        if token.tokenType == THIS {
            
        	expr = &NodeVariable{Name: "this"}
            
        } else if p.peek(0).tokenType == LEFT_PAREN {
            
            expr = p.funcCall()
            
        } else if p.peek(0).tokenType == SCOPE_RESOLVE {
            
            p.advance()
            expr = &NodePkgResolve{Pkg: token.Lexeme, Resolution: p.expression()}
            
        } else if token.tokenType.equals(GET) {
            
            var propIndices []NodeExpr
            //fmt.Println(Red, "Mattes ça", Reset, p.peek(0))
            if p.peek(0).tokenType.equals(LEFT_BRACKET) {
                //fmt.Println(Red, "Mattes ça", Reset, p.peek(0))
                p.advance()
                for !p.peek(0).tokenType.equals(RIGHT_BRACKET) {
                    if p.isAtEnd() {CatchError("Expected right bracket but found End Of File"); return nil}
                    //fmt.Println(Red, "Mattes ça", Reset, p.peek(0))
                    propIndices = append(propIndices, p.expression())
                    //fmt.Println(Red, "Mattes ça", Reset, p.peek(0))
                    if p.peek(0).tokenType.equals(COMMA) {p.advance()}
                    
                }
                
                //fmt.Println(Red, "Mattes ça", Reset, p.peek(0))
                p.advance()
                
            }
            
            if !p.peek(0).tokenType.equals(LEFT_BRACE) {CatchError("Missing left bracket"); return nil}
            var varParam []NodeExpr
            var varList map[string][]NodeExpr = make(map[string][]NodeExpr)
            p.advance()
            
            for !p.peek(0).tokenType.equals(RIGHT_BRACE) {
                
                currentVar := p.advance().Lexeme
                p.advance()
                for p.peek(0).tokenType != RIGHT_PAREN {
                    varParam = append(varParam, p.primary()) 
                    if p.peek(0).tokenType == COMMA {p.advance()}
                    
                }
                
                p.advance()
                if p.peek(0).tokenType == COMMA {p.advance()}
                varList[currentVar] = varParam
                
            }
            
            p.advance()
            
            switch token.tokenType {
                case GET: return &NodeExprGetter{Expr: nil, Vars: varList, Indices: propIndices}
                default: return nil
                
            }
            
        } else {
            
            expr = &NodeVariable{Name: token.Lexeme}
            
        }
        
        if p.peek(0).tokenType.equals(LEFT_BRACKET) {
            
                var index NodeExpr
                p.advance()
                
                if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {
                    
                    index = p.expression()
                }
                
                p.advance()
                return &NodeExprListIndex{Expr: expr, Index: index}
                
        }
        
        var class string
        
        //if field {class = p.envi.Variable[p.peek(-1).Lexeme]}
        
        for p.peek(0).tokenType == DOT {
            
            //object := token.Lexeme
            
            p.advance()
            field := p.advance().Lexeme
            var objName Token = p.peek(-3)
            
            //if class != "" {class += "_"}
            if p.peek(0).tokenType.equals(STAR) {
                field += p.advance().Lexeme
            }
            
            if p.peek(0).tokenType == LEFT_PAREN {
                
                p.advance()
                var argsList []NodeExpr
                
                for p.peek(0).tokenType != RIGHT_PAREN {
                    argsList = append(argsList, p.expression())
                    if p.peek(0).tokenType == COMMA {p.advance()}
                }
                p.advance()
                //isStatic := p.envi.getStatic(field, object)
                //if objName.Lexeme == "this" {class = p.envi.Name}
                if class == "" { class += objName.Lexeme }
                
                fmt.Println("CLASS1", class, p.envi.Module)
                expr = &NodeExprMethodCall{Class: class, Parent: expr, Name: field, Args: argsList, Static: false}
                
            } else if p.peek(0).tokenType.equals(LEFT_BRACE) || p.peek(0).tokenType.equals(LEFT_BRACKET) && (field == "get" || field == "set") {
                var propIndices []NodeExpr
                if p.peek(0).tokenType.equals(LEFT_BRACKET) {
                    p.advance()
                    for !p.peek(0).tokenType.equals(RIGHT_BRACKET) {
                        if p.isAtEnd() {CatchError("Expected right bracket but found End Of File")}
                        propIndices = append(propIndices, p.expression())
                        if p.peek(0).tokenType.equals(COMMA) {p.advance()}
                    }
                    p.advance()
                }
                p.advance()
                if field == "get" {
                    var varParam []NodeExpr
                    var varList map[string][]NodeExpr = make(map[string][]NodeExpr)
                    if p.peek(0).tokenType != LEFT_BRACE {CatchError("Missing left brace")}
                    p.advance()
                    for p.peek(0).tokenType != RIGHT_BRACE {
                        currentVar := p.advance().Lexeme
                        p.advance()
                        for p.peek(0).tokenType != RIGHT_PAREN {
                            varParam = append(varParam, p.expression()) 
                            if p.peek(0).tokenType == COMMA {p.advance()}
                            if p.isAtEnd() {CatchError("EOF while property access never took end"); return nil}
                        }
                        p.advance()
                        if p.peek(0).tokenType == COMMA {p.advance()}
                        varList[currentVar] = varParam
                        if p.isAtEnd() {CatchError("EOF while property access never took end"); return nil}
                    }
                    p.advance()
                    switch field {
                        case "get": return &NodeExprGetter{Expr: expr, Vars: varList, Indices: propIndices}
                        default: break
                    }
                }
            } else if p.peek(0).tokenType.equals(LEFT_BRACKET) && !p.peek(0).tokenType.isProperty(){
                var index NodeExpr
                p.advance()
                if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {
                    index = p.expression()
                }
                p.advance()
                return &NodeExprListIndex{Expr: expr, Index: index}
            } else {
                expr = &NodeGet{Object: expr, Symbol: ".", Field: field}
            }
            if field != "" {class += "_"}
            //class = p.envi.Table[objName.Lexeme + field].Name
        }
        return expr
    }
    if token.tokenType == LEFT_PAREN {
        expr := p.expression()
        if p.peek(0).tokenType != RIGHT_PAREN && p.isAtEnd() {
            CatchError("Expected ) before End Of File")
        } else {
            p.advance()
            return &NodeGroup{Expression: expr}
        }
    }
    
    if token.tokenType.equals(CCALL) {
        fmt.Println(Red, "Calling an expr", Reset)
        var callerName string
        var called []NodeExpr
        if !p.peek(0).tokenType.isExprCall() {fmt.Println(p.peek(0)); CatchError("Expected action for the C-Expr-Caller found " + p.peek(0).Lexeme); return nil}
        actionName := p.advance().Lexeme
        fmt.Println("Call is", actionName)
        switch actionName {
            case "call":
            case "func":
                if p.peek(0).tokenType != IDENTIFIER {CatchError("Expected Identifier"); return nil}
                callerName = p.advance().Lexeme
                //fmt.Println("Caller", callerName)
                if p.peek(0).tokenType != LEFT_PAREN {CatchError("Expected left parenthesize for C function call"); return nil}
                p.advance()
                for p.peek(0).tokenType != RIGHT_PAREN {
                    called = append(called, p.expression())
                    if p.peek(0).tokenType == COMMA {p.advance()}
                }
                p.advance()
            default:
            break
        }
        //fmt.Println("Token actuel", p.peek(0))
        return &NodeExprC{Action: actionName, Called: called, CallerName: callerName}
    }
    
    if token.tokenType == NULL {
        return &NodeLiteral{Value: token.Lexeme, Type: token}
    }
    
    if token.tokenType == TRUE {
        return &NodeLiteral{Value: token.Lexeme, Type: token}
    }
    
    if token.tokenType == FALSE {
        return &NodeLiteral{Value: token.Lexeme, Type: token}
    }
    
    if token.tokenType == STRING {
        fmt.Println(token, "Is a String")
        if p.peek(0).tokenType.equals(LEFT_BRACKET) {
            expr := &NodeLiteral{Value: token.Lexeme, Type: token}
            var index NodeExpr
            p.advance()
            if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {
                index = p.expression()
            }
            p.advance()
            return &NodeExprListIndex{Expr: expr, Index: index}
        }
        return &NodeLiteral{Value: token.Lexeme, Type: token}
    }
    //fmt.Println("Token: ", p.peek(0))
    CatchError("Unknown " + token.Lexeme)
    return nil
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
    var right NodeExpr
    token := p.peek(0)
    
    if token.tokenType.equals(MINUS, BANG) || token.tokenType.isType() {
        if token.tokenType.isType() {
            p.advance()
        }
        p.advance()
        
        if token.tokenType.isType() {
            right = p.expression()
            p.advance()
        } else {
            right = p.unary()
        }
        
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
        if token.tokenType.equals(STAR, SLASH) {
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

func (p *Parser) properties() map[string][]any {
    if !p.isAtEnd() {
        var propertyList map[string][]any = make(map[string][]any)
        p.advance()
        for p.peek(0).tokenType != RIGHT_PAREN {
            if !p.peek(0).tokenType.isProperty() {CatchError("Expected property call")}
            //fmt.Println("Property call", p.peek(0))
            var currentProp string = p.peek(0).Lexeme
            var basicSlice []any
            if p.peek(1).tokenType.equals(LEFT_BRACKET) {
                //fmt.Println("index-access case", p.peek(0))
                p.advance()
                p.advance()
                //fmt.Println("advance", p.peek(0))
                if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {CatchError("Expected left bracket for property index-access")}
                currentProp += "[]"
            }
            basicSlice = append(basicSlice, currentProp)
            propertyList[currentProp] = basicSlice
            p.advance()
            if p.peek(0).tokenType == COMMA {p.advance()}
        }
        return propertyList
    }
    CatchError("Property call declared but with no end")
    return nil
}

func (p *Parser) propertiesDesc() map[string][]any {
    var propList map[string][]any = make(map[string][]any)
    
    p.advance()
    for p.peek(0).tokenType != RIGHT_BRACE {
        
        var paramList []NodeStmt
        var propCode []NodeStmt
    	var propAttributes []any
        var arrayIndex NodeExpr
        
        if !p.peek(0).tokenType.isProperty() {CatchError("Expected a property from this variable")}
        //First we get the name as the key
        propName := p.advance().Lexeme
        if p.peek(0).tokenType.equals(LEFT_BRACKET) {
            arrayIndex = p.primary()
            if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {CatchError("Expected right bracket for array-property")}
            p.advance()
        }
        //Then we consume all the params and make a slice of the as the first attribute
        for p.peek(0).tokenType != RIGHT_PAREN {
            paramList = append(paramList, p.parseParam())
        }
        p.advance()
        p.advance()
        //Then we consume the code between braces as the second attribute
        for p.peek(0).tokenType != RIGHT_BRACE {
            propCode = append(propCode, p.ParseStmt())
        }
        p.advance()
        propAttributes = append(propAttributes, propName, paramList, propCode, arrayIndex)
        propList[propName] = propAttributes
    }
    p.advance()
    return propList
}

func (p *Parser) PropertyAccessorStmt() NodeStmt {
    if !p.isAtEnd() {
        var exprs []NodeExpr
        accessor := p.advance()
        p.expect("Expected left brace '{'.", p.peek(0), LEFT_BRACE)
        for !p.peek(0).tokenType.equals(RIGHT_BRACE) {
            exprs = append(exprs, p.primary())
            if p.peek(0).tokenType.equals(COMMA) { p.advance() }
        }
        p.expect("Expected right brace '{'.", p.peek(0), RIGHT_BRACE)
        
        switch accessor.Lexeme {
            case "set": return &NodeStmtSetter{Expr: nil}
            default:
        }
    }
    CatchError("Unexpected error while parsing property accessor " + p.peek(0).Lexeme)
    return nil
}

func (p *Parser) parseVar() NodeStmt {
    p.advance()
    if p.peek(0).tokenType.equals(LEFT_BRACE) {
        return p.multiVars()
    }
    return p.varAssignment(0, nil)
}

func (p *Parser) varAssignment(pass int, outProp map[string][]any) NodeStmt {
    var varVal NodeExpr = nil
    var varType NodeExpr
    var TypeTok Token
    var Vars []NodeStmt
    //var registeredType string
    var propertyList map[string][]any
    var isPointer bool = false
    
    if !p.isAtEnd() {
        
        if p.peek(0).tokenType.equals(LEFT_PAREN) && pass < 1 {
            propertyList = make(map[string][]any)
            propertyList = p.properties()
            p.advance()
        }
        if pass >= 1 {propertyList = outProp}
        if p.peek(0).tokenType.equals(STAR) {
            p.advance()
            isPointer = true
        }
        TypeTok = p.peek(0)
        varType = p.parseType()
        //if Type, ok := varType.(*NodeType); ok {
        //    registeredType = Type.Type
        //    if registeredType == "[" {registeredType = "array"}
        //}
        
        if p.peek(0).tokenType != IDENTIFIER { CatchError("Expected an identifier for function name"); return nil}
        varName := p.advance().Lexeme
        
        if p.peek(0).tokenType == EQUAL {
            p.advance()
            varVal = p.expression()
        }
        
        global := false
        if p.envi.isGlobal() {global = true}
        
        if p.peek(0).tokenType != SEMICOLON && p.peek(0).tokenType != LEFT_BRACE {
            CatchError("Errounous var end, forgot semicolon or var properties description")
            return nil
        }
        
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
        } else {
            if propertyList != nil {
                propertyList = p.propertiesDesc()
                for key, val := range propertyList {
                propertyList[varName + "_" + key] = val
                delete(propertyList, key)
                }
            } else if TypeTok.tokenType.equals(STRUCT_TYPE, ENUM_TYPE) {
                switch TypeTok.tokenType {
                    case STRUCT_TYPE: Vars = p.parseStruct()
                    case ENUM_TYPE: Vars = p.parseEnum()
                    default:
                }
            } else {CatchError("Unexpected" + p.peek(0).Lexeme + "in variable declaration statement"); return nil}
        }
        p.envi.Table[varName] = &Symbol{
            Name: varName,
            Kind: S_VARIABLE,
            Module: p.Module,
            //Type: varType,
            Pointer: isPointer,
            Scope: S_PROTECTED,
        }
        return &NodeStmtVar{Name: varName, Properties: propertyList, Type: varType, Value: varVal, Global: global, Vars: Vars}
    }
    CatchError("Reached End Of File in an invalid variable declaration")
    return nil
}

func (p *Parser) multiVars() NodeStmt {
    if !p.isAtEnd() {
        var vars []NodeStmt
        var pass int = 0;
        var prop map[string][]any
        p.advance()
        if p.peek(0).tokenType.equals(LEFT_PAREN) {
            prop = make(map[string][]any)
            prop = p.properties()
            pass = 1
        }
        for !p.peek(0).tokenType.equals(RIGHT_BRACE) {
            vars = append(vars, p.varAssignment(pass, prop))
        }
        p.advance()
        return &NodeStmtMultiVars{Vars: vars}
    }
    CatchError("Invalid multivars")
    return nil
}

func (p *Parser) constAssignment() NodeStmt {
    var constVal NodeExpr = nil
    
    if !p.isAtEnd() {
        p.advance()
        
        //Type := p.peek(0).Lexeme
        constType := p.parseType()
        
        if p.peek(0).tokenType != IDENTIFIER { CatchError("Expected an identifier for const name"); return nil}
        constName := p.advance().Lexeme
        
        if p.peek(0).tokenType != EQUAL {
            CatchError("Invalid Const Stmt")
        }
        p.advance()
        constVal = p.expression()
        
        global := false
        if p.envi.isGlobal() {global = true}
        
        if p.peek(0).tokenType == SEMICOLON {
            p.advance()
            p.envi.Table[constName] = &Symbol{
                Name: constName,
                Kind: S_CONST,
                Module: p.Module,
            }
            return &NodeStmtConst{Name: constName, Type: constType, Value: constVal, Global: global}
        }
    }
    CatchError("Reached End Of File in an invalid const declaration")
    return nil
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
    CatchError("Missing semi-colon ';'")
    return nil
}

func (p *Parser) ccall() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        var callerName string
        var called []NodeExpr
        if !p.peek(0).tokenType.isAction() {fmt.Println(p.peek(0)); CatchError("Expected action for the C-Caller"); return nil}
        actionName := p.advance().Lexeme
        switch actionName {
        case "include":
            for p.peek(0).tokenType != SEMICOLON {
                if !p.peek(0).tokenType.equals(LESS) {
                    called = append(called, p.expression())
                } else {
                    p.advance()
                    included := &NodeLiteral{Value: "<" + p.advance().Lexeme + ".h>"}
                    p.advance()
                    called = append(called, included)
                }
                if p.peek(0).tokenType == COMMA {p.advance()}
            }
            break
        case "write":
            code := &NodeLiteral{Value: p.primary()}
            called = append(called, code)
        }    
        if p.peek(0).tokenType != SEMICOLON {fmt.Println(p.peek(0)); CatchError("Expected semicolon"); return nil}
        p.advance()
        return &NodeStmtC{Action: actionName, Called: called, CallerName: callerName}
    }
    CatchError("End Of File reached in C call")
    return nil
}

func (p *Parser) blockStmt(sym ...*Symbol) NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        parScope := p.envi
        p.envi = p.envi.NewScope()
        for _, s := range sym {
            if s != nil {
                p.envi.Table[s.Name] = s
            }
        }
        var stmts []NodeStmt
        for p.peek(0).tokenType != RIGHT_BRACE {
            stmts = append(stmts, p.ParseStmt())
            if p.peek(0).tokenType == SEMICOLON {p.advance()}
        }
        p.advance()
        p.envi = parScope
        return &NodeBlock{Instructions: stmts}
    }
    CatchError("Missing }")
    return nil
}

func (p *Parser) ifStmt() NodeStmt {
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        condition = p.comparison()
        if p.peek(0).tokenType != RIGHT_PAREN {CatchError("Expected right parenthesize"); return nil}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtIf{Condition: condition, Result: result}
    }
    CatchError("Invalid if statement")
    return nil
}

func (p *Parser) loop() NodeStmt {
    var looping NodeExpr
    var result NodeStmt
    
    if !p.isAtEnd() {
        p.advance()
        looping = p.parseLoop()
        //fmt.Println("LOOOOOP", p.peek(0))
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtLoop{Looping: looping, Result: result}
    }
    CatchError("Invalid loop stmt")
    return nil
}

func (p *Parser) parseLoop() NodeExpr {
    if !p.isAtEnd() {
        var iterated NodeStmt
        var operator string        
        var min NodeExpr
        var max NodeExpr
        var increm NodeExpr
        var each NodeExpr
        if p.peek(0).tokenType.isType() || p.peek(0).tokenType.equals(IDENTIFIER) {
            //fmt.Println("TYPE", p.peek(0))
            _type := p.parseType()
            //fmt.Println("NAME", p.peek(0))
            name := p.advance().Lexeme
            //fmt.Println("OP", p.peek(0))
            iterated = &NodeStmtVar{Name: name, Type: _type}
            operator = p.advance().Lexeme
        }
        //fmt.Println("Op:", operator, operator == ":")
        p.advance()
        if operator != ":" {
            //fmt.Println("HERE 1", p.peek(0))
            min = p.expression()
            //fmt.Println("HERE 2", p.peek(0))
            if p.peek(0).tokenType.equals(COMMA) {
                p.advance()
                max = p.expression()
                //fmt.Println("HERE 3", p.peek(0))
                if p.peek(0).tokenType.equals(COMMA) {
                    p.advance()
                    increm = p.expression()
                    //fmt.Println("HERE 4", p.peek(0))
                }
            } else {max = min; min = nil}
        } else {
            each = p.expression()
            fmt.Println(p.peek(0))
            if !p.peek(0).tokenType.equals(RIGHT_BRACKET) {CatchError("Expected right bracket for loop-each statement")}
        }
        p.advance()
        return &NodeExprIter{Min: min, Max: max, Increm: increm, Each: each, Iter: iterated, Operator: operator}
    }
    CatchError("Loop args reached EOF")
    return nil
}

func (p *Parser) whileStmt() NodeStmt {
    var condition NodeExpr
    var result NodeStmt
    if !p.isAtEnd() {
        p.advance()
        if p.peek(0).tokenType == LEFT_PAREN {p.advance();}
        condition = p.comparison()
        if p.peek(0).tokenType != RIGHT_PAREN {CatchError("Expected right parenthesize )"); return nil}
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            result = p.blockStmt()
        } else {
            result = p.ParseStmt()
        }
        return &NodeStmtWhile{Condition: condition, Result: result}
    }
    CatchError("Invalid while statement")
    return nil
}

func (p *Parser) breakStmt() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        p.expect("Expected semi-colon after break statement", p.peek(0), SEMICOLON)
        return &NodeStmtBreak{}
    }
    CatchError("Invalid Break Stmt")
    return nil
}

func (p *Parser) parseParam() NodeStmt {
    var paramVal Node = nil
    
    if !p.isAtEnd() {
        if p.peek(0).tokenType == RIGHT_PAREN { return nil}
        paramType := p.parseType()
        if p.peek(0).tokenType != IDENTIFIER {fmt.Println(p.peek(0)); CatchError("Expected identifier"); return nil}
        paramName := p.advance().Lexeme
        
        if p.peek(0).tokenType == COMMA {
            p.advance()
        }
        return &NodeStmtVar{Name: paramName, Type: paramType, Value: paramVal}
        
    }
    CatchError("Missing right parenthesize )")
    return nil
}

func (p *Parser) funcInit() NodeStmt {
    var (
        paramList []NodeStmt = nil
        rType Token
        rName Token
        specialCases bool = false
        argsSymbol []Symbol
        typeIden *Symbol
    )
    if !p.isAtEnd() {
        p.advance()
        //Type := p.peek(0).Lexeme
        returnType := p.parseType()
        if p.peek(0).tokenType.equals(DOT) {
            p.expect("Expected an identifier", p.peek(1), IDENTIFIER)
            specialCases = true
            rType = p.peek(-2)
        }
        if p.peek(0).tokenType != IDENTIFIER {CatchError("Expected identifier found " + p.peek(0).Lexeme); return nil}
        funcName := p.advance().Lexeme
        if !specialCases {
            if p.peek(0).tokenType.equals(DOT) {p.advance(); rType = p.peek(-2); funcName = p.advance().Lexeme;}
        }

        if p.peek(1).tokenType != RIGHT_PAREN {
            p.advance()
            for p.peek(0).tokenType != RIGHT_PAREN {
                paramList = append(paramList, p.parseParam())
            }
            p.advance()
        } else {
            p.advance()
            p.advance()
        }
        
        p.envi.Table[funcName] = &Symbol{
            Name: funcName,
            Kind: S_FUNC,
            Module: p.Module,
            Params: argsSymbol,
        }
        if p.peek(0).tokenType.equals(COLON) {
            p.advance()
            fmt.Println("Iden", p.peek(-1), p.peek(0), p.peek(1), p.peek(0).tokenType.equals(IDENTIFIER))
            rName = p.expect("Expected identifier", p.peek(0), IDENTIFIER)
            typeIden = &Symbol{
                Name: rName.Lexeme,
                Kind: S_VARIABLE,
                Module: p.Module,
                Pointer: true,
            }
        }
        if p.peek(0).tokenType != LEFT_BRACE {CatchError("Missing left brace"); return nil}
        
        code := p.blockStmt(typeIden)
        switch funcName {
            case "new": return &NodeStmtConstructor{Return: returnType, Param: paramList, Code: code, ReceiverType: rType, ReceiverName: rName}
            default: return &NodeStmtFuncInit{Return: returnType, Name: funcName, Param: paramList, Code: code, ReceiverType: rType, ReceiverName: rName}
        }
    }
    CatchError("Invalid function def statement")
    return nil
}

func (p *Parser) funcCall() NodeExpr {
    var argsList []NodeExpr
    if !p.isAtEnd() {
        funcName := p.peek(-1).Lexeme
        p.advance()
        for p.peek(0).tokenType != RIGHT_PAREN {
            arg := p.expression()
            argsList = append(argsList, arg)
            if p.peek(0).tokenType == COMMA {p.advance()}
        }
        p.advance()
        return &NodeExprFuncCall{Name: funcName, Args: argsList}
    }
    CatchError("Invalid function call")
    return nil
}

func (p *Parser) staticStmt() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        stmt := p.ParseStmt()
        return &NodeStaticStmt{Stmt: stmt}
    }
    CatchError("Errounous static call")
    return nil
}

func (p *Parser) scopeAccess() NodeStmt {
    if !p.isAtEnd() {
        modifier := p.advance()
        var stmtList []NodeStmt
        //fmt.Println("Current point", p.peek(0))
        if !p.peek(0).tokenType.equals(LEFT_BRACE) {
            stmtList = append(stmtList, p.parseScope(modifier))
        } else {
            p.advance()
            for !p.peek(0).tokenType.equals(RIGHT_BRACE) {
                stmtList = append(stmtList, p.parseScope(modifier))
            }
            p.advance()
        }
        //fmt.Println(stmt)
        return &NodeScopeAcces{Modifier: modifier, Stmt: stmtList}
    }
    CatchError("Scope access in End Of File")
    return nil
}

func (p *Parser) parseScope(modifier Token) NodeStmt {
    return nil
}

func (p *Parser) typeDef() NodeStmt {
    if !p.isAtEnd() {
        p.advance()
        Type := p.peek(0)
        typeData := p.parseType()
        var typeName string = p.advance().Lexeme
        var typeVars []NodeStmt = nil
        fmt.Println(Type, "type's type", Type.tokenType == STRUCT_TYPE)
        if Type.tokenType.equals(STRUCT_TYPE) {
            typeVars = p.parseStruct()
        } else if Type.tokenType.equals(ENUM_TYPE) {
            typeVars = p.parseEnum()
        } else {
            if !Type.tokenType.equals(SEMICOLON) { CatchError("Missing Semicolon") }
        }
        return &NodeStmtTypeDef{Type: typeData, Name: typeName, Vars: typeVars}
    }
    CatchError("Invalid typedef, not ended correctly before EOF")
    return nil
}

func (p *Parser) ModuleDef() NodeStmt {
    if !p.isAtEnd() {
        var name string
        p.advance()
        for !p.peek(0).tokenType.equals(SEMICOLON) {
            name += p.advance().Lexeme
            if p.peek(0).tokenType.equals(SLASH) {
                name += p.advance().Lexeme
                name += p.advance().Lexeme
            }
            p.envi.Module = name
            p.envi.FilePath = filePath
        }
        if p.advance().tokenType != SEMICOLON {CatchError("Missing Semicolon"); return nil}
        return &NodeStmtModule{Name: name}
    }
    CatchError("Invalid Package Definition")
    return nil
}

func (p *Parser) importPkg() NodeStmt {
    if !p.isAtEnd() {
        var pathes []string
        p.advance()
        if p.peek(0).tokenType == LEFT_BRACE {
            p.advance()
            for p.peek(0).tokenType != RIGHT_BRACE {
                path := "\"" + p.peek(0).Lexeme[1:len(p.peek(0).Lexeme)-1] + ".h" + "\""
                if strings.Contains(path, "/") {path = "\"../" + path[1:]}
                p.envi.ImportedModules[path[1:len(path)-3]] = true
                pathes = append(pathes, path)
                p.advance()
                if p.peek(0).tokenType == COMMA {p.advance()}
            }
        } else if p.peek(0).tokenType == STRING {
            path := "\"" + p.peek(0).Lexeme[1:len(p.peek(0).Lexeme)-1] + ".h" + "\""
            if strings.Contains(path, "/") {path = "\"../" + path[1:]}
            p.envi.ImportedModules[path[2:len(path)-3]] = true
            pathes = append(pathes, path)
            p.advance()
            if p.peek(0).tokenType != SEMICOLON {CatchError("Something")}
        }
        p.advance()
        return &NodeImport{Names: pathes}
    }
    CatchError("Invalid Package Definition")
    return nil
}