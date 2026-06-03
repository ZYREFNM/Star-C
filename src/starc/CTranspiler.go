package main

import (
	"fmt"
    "os"
    "strings"
    //"slices"
)

type Transpiler struct {
    fileName string
    currentClass string
    Package string
    importField string
    static bool
    CglobalVars string
    CglobalFuncs string
    HglobalVars string
    HglobalFuncs string
}

func (t *Transpiler) WriteInFile(code string, format string) {
    codeBytes := []byte(code)
    os.WriteFile(fmt.Sprintf("%s%s", t.fileName, format), codeBytes, 0644)
}

func (t *Transpiler) matchType(Type string) string {
    var res string
    var var_args bool = false
    if Type[len(Type)-3:len(Type)] == "..." {
        Type = Type[:len(Type)-3]
        var_args = true
    }
    switch Type {
        case "string": res = "char*"
        //case "bool": return "boolean"
        case "int8": res = "int8_t"
        case "int16": res = "int16_t"
        case "int32": res = "int32_t"
        case "int64": res = "int64_t"
        case "float32": res = "float"
        case "float64": res = "double"
        default: res = Type
    }
    if var_args {res = res + "..."}
    return res
}

func (t *Transpiler) matchAction(Action string, Called []NodeExpr, CallerName string) string {
    switch Action {
        case "include":
            var code string
            for _, called := range Called {
                code += "#include " + t.TranslateC(called) + "\n"
            }
            return code
        case "function":
            var args []string
            for _, arg := range Called {
                args = append(args, t.TranslateC(arg))
            }
            return fmt.Sprintf("%s(%s)", CallerName, strings.Join(args, ", "))
        default: return ""
    }
}

func (t *Transpiler) basicProperty(Prop string, varName string, varType string, value string) string {
    switch Prop {
        case "get": return fmt.Sprintf("{\nreturn this->%s;\n}", varName)
        case "set": return fmt.Sprintf("{\nif (value) {\nthis->%s = value;\n}\n}", varName)
        default: return ""
    }
}

func (t *Transpiler) matchProperty(Name string, Attributes []any, varName string, varType string, Header bool) string {
    var class string = t.currentClass
    var pointer string
    var point string
    var otherArgs string
    var Property string
    var defaultValue string
    var code string
    if t.currentClass != "" {
        class += "_"
        pointer = t.currentClass
        point = "this"
    } else {
        pointer = varType
        point = varName
    }
    if len(Attributes) > 1 {
        if params, ok := Attributes[1].([]NodeStmt); ok {
            if paramOne, isOk := params[2].(NodeStmt); isOk {
                if param, isParam := paramOne.(*NodeStmtVar); isParam {
                    defaultValue = param.Name
                }
            }
        }
    }
    
    if name, ok := Attributes[0].(string); ok {Property = name}
    if Header {
        code = ";"
    } else {
        if len(Attributes) > 1 {
            if attributes, ok := Attributes[1].([]NodeStmt); ok {
                code += "{\n"
                for _, att := range attributes {
                    code += "    " + t.TranslateC(att) + "\n"
                }
                code += "}\n"
            }
        } else {
            code = t.basicProperty(Property, varName, varType, defaultValue)
        }
    }
    if otherArgs == "" {
        otherArgs += fmt.Sprintf("%s value", varType)
    }
    switch Property {
        case "get": return fmt.Sprintf("%s %s%s_get(%s* %s)%s", varType, class, varName, pointer, point, code)
        case "set": return fmt.Sprintf("void %s%s_set(%s* %s, %s)%s", class, varName, pointer, point, otherArgs, code)
        default: return ""
    }
}

func (t *Transpiler) matchAllocate(alloc string) string {
    switch alloc {
        case "memory": return "malloc"
        case "free": return "free"
        case "clean": return "calloc"
        case "size": return "sizeof"
    	default: return alloc
    }
}

func (t *Transpiler) TranslateH(node Node) string {
    switch n := node.(type) {
        
        case *NodeLiteral: return fmt.Sprintf("%v", n.Value)
        
        case *NodeBinary: return fmt.Sprintf("%s %s %s", t.TranslateH(n.Left), n.Operator, t.TranslateH(n.Right))
        
        case *NodeUnary: return fmt.Sprintf("%s (%s)", n.Operator, t.TranslateH(n.Right))
        
        case *NodeGroup: return fmt.Sprintf("(%s)", t.TranslateH(n.Expression))
        
        case *NodeStmtVar:
        	Type := t.matchType(n.Type.Lexeme)
            varEnd := ";"
            if n.Value != nil {varEnd = fmt.Sprintf(" = %s;", t.TranslateH(n.Value))}
            if n.Properties != nil {
                t.HglobalFuncs += "\n"
                for key, prop := range n.Properties {
                	t.HglobalFuncs += t.matchProperty(key, prop, n.Name, Type, true) + "\n"
                }
            }
        	code := fmt.Sprintf("%s %s%s", Type, n.Name, varEnd)
            if n.Global {
                t.HglobalVars += code
                return ""
            }
            return code
        
        case *NodeStmtConst:
            Type := t.matchType(n.Type.Lexeme)
            t.HglobalVars += fmt.Sprintf("extern const %s %s;\n", Type, n.Name)
            return ""
        
        case *NodeStaticStmt:
            t.static = true
            stmt := t.TranslateH(n.Stmt)
            t.static = false
            return stmt
        
        case *NodeStmtFuncInit:
        	var funcName string = n.Name
            var paramList []string
            if n.Name == "main" {return ""}
            Type := t.matchType(n.Return)
            
            if t.currentClass != "" {
                funcName = t.currentClass + "_" + funcName
                if !t.static {
                    paramList = append(paramList, fmt.Sprintf("%s* this", t.currentClass))
                }
            }
            
            for _, p := range n.Param {
                param := strings.TrimSuffix(t.TranslateH(p), ";")
                paramList = append(paramList, param)
            }
            t.HglobalFuncs += fmt.Sprintf("%s %s__%s(%s);\n", Type, t.Package, funcName, strings.Join(paramList, ", "))
            return ""
        
        case *NodeStmtTypeDef:
            typeData := n.Type
            typeName := n.Name
            code := ""
            var Type string = t.matchType(typeData.Lexeme)
            if typeData.tokenType == STRUCT {
                code += " {\n"
                for _, init := range n.Vars {
                    code += "    " + t.TranslateH(init) + "\n"
                }
                code += "}"
                Type = "struct"
            }
        	t.HglobalVars += fmt.Sprintf("typedef %s%s %s;\n", Type, code, typeName)
            return ""
        
        case *NodeStmtClass:
        	className := n.Name
            var classVars string
            var classCode string
            t.currentClass = className
            for _, e := range n.Code {
                if _, ok := e.(*NodeStmtVar); ok {
                    classVars += "    " + t.TranslateH(e) + "\n"
                } else {
                    classCode += t.TranslateH(e) + "\n"
                }
            }
            t.currentClass = ""
            t.HglobalVars += fmt.Sprintf("typedef struct {\n%s} %s;\n%s", classVars, className, classCode)
            return ""
        
        default: return ""
    }
}

func (t *Transpiler) TranslateC(node Node) string {
    switch n := node.(type) {
        
        case *NodeStmtExpr:
        	return t.TranslateC(n.Expr) + ";"
        
        case *NodeBinary:
        	return fmt.Sprintf("%s %s %s", t.TranslateC(n.Left), n.Operator, t.TranslateC(n.Right))
            
        case *NodeLiteral: return fmt.Sprintf("%v", n.Value)
        
        case *NodeUnary: return fmt.Sprintf("%s (%s)", n.Operator, t.TranslateC(n.Right))
        
        case *NodeExprConcat:
        	return fmt.Sprintf("star_concat(%s, %s)", t.TranslateC(n.To), t.TranslateC(n.From))
            
        case *NodeGroup: return fmt.Sprintf("(%s)", t.TranslateC(n.Expression))
        
        case *NodeGet:
        	symbol := "."
            target := t.TranslateC(n.Object)
        	if target == "this" {symbol = "->"}
            //fmt.Println(fmt.Sprintf("Getting -> %s %s %s", target, symbol, n.Field))
            return fmt.Sprintf("%s%s%s", target, symbol, n.Field)
        
        case *NodePkgResolve:
            pkg := n.Pkg
            if pkg == "" {
                pkg = t.Package
            }
            res := t.TranslateC(n.Resolution)
            //fmt.Println("Resolution", pkg, res)
            return pkg + "__" + res
        
        case *NodeStmtVar:
        	Type := t.matchType(n.Type.Lexeme)
            varEnd := ";"
            if n.Value != nil {varEnd = fmt.Sprintf(" = %s;", t.TranslateC(n.Value))}
            if n.Properties != nil {
                t.CglobalFuncs += "\n"
                for key, prop := range n.Properties {
                	t.CglobalFuncs += t.matchProperty(key, prop, n.Name, Type, false) + "\n"
                }
            }
        	code := fmt.Sprintf("%s %s%s", Type, n.Name, varEnd)
            if n.Global {
                t.CglobalVars += code
                return ""
            }
            return code
        
        case *NodeStmtConst:
            var code string
        	Type := t.matchType(n.Type.Lexeme)
            code = fmt.Sprintf("const %s %s = %s;", Type, n.Name, t.TranslateC(n.Value))
            if n.Global {
                t.CglobalVars += code
                return ""
            }
            return code
            
        case *NodeAssignment: return fmt.Sprintf("%s = %s;", t.TranslateC(n.Target), t.TranslateC(n.Value))
        
        case *NodeVariable: return n.Name
        
        case *NodeExprAlloc:
            return fmt.Sprintf("%s(%s)", t.matchAllocate(n.Allocation), t.TranslateC(n.Size))
        
        case *NodeBlock:
        	var code string = "{\n"
            for _, stmt := range n.Instructions {
                code += "	" + t.TranslateC(stmt) + "\n"
            }
            code += "}"
        	return code
            
        case *NodeStmtReturn:
        	if n.Value == nil {return "return;"}
            return fmt.Sprintf("return %s;", t.TranslateC(n.Value))
        case *NodeStmtC:
            return t.matchAction(n.Action, n.Called, n.CallerName) + ";"
        
        case *NodeStmtIf:
        	condition := t.TranslateC(n.Condition)
            result := t.TranslateC(n.Result)
            return fmt.Sprintf("if (%s) %s", condition, result)
        case *NodeStmtLoop:
            loops := t.TranslateC(n.Looping)
            result := t.TranslateC(n.Result)
            return  fmt.Sprintf("for (int i = 0; i < %s; i++) %s", loops, result)
        case *NodeStmtWhile:
        	condition := t.TranslateC(n.Condition)
            result := t.TranslateC(n.Result)
            return fmt.Sprintf("while (%s) %s", condition, result)
            
        case *NodeStmtFuncInit:
        	var funcName string = n.Name
            var pack string
        	var list []string
            if t.currentClass != "" {
                funcName = t.currentClass + "_" + funcName
                if !t.static {
                    list = append(list, fmt.Sprintf("%s* this", t.currentClass))
                }
            }
            if t.Package != "" {
                if funcName != "main" {
                    pack = t.Package + "__"
                }
            }
        	for _, p := range n.Param {
                param := t.TranslateC(p)
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            code := t.TranslateC(n.Code)
        	return fmt.Sprintf("%s %s%s(%s) %s", t.matchType(n.Return), pack, funcName, strings.Join(list, ", "), code)
        
        case *NodeStmtConstructor:
            var list []string
            var code string
            for _, p := range n.Param {
                param := t.TranslateC(p)
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            code = t.TranslateC(n.Code)
            code = code[:2] + fmt.Sprintf("	%s* this = malloc(sizeof(%s));\n", n.Return, n.Return) + code[2:len(code)-2] + "\n	return this;"
            return fmt.Sprintf("%s* %s_new(%s) %s", n.Return, n.Return, strings.Join(list, ", "), code)
        case *NodeExprFuncCall:
        	var argsList []string
            for _, arg := range n.Args {
                argsList = append(argsList, t.TranslateC(arg))
            }
        	return fmt.Sprintf("%s(%s)", n.Name, strings.Join(argsList, ", "))
        
        case *NodeExprMethodCall:
        	var argsList []string
            var parPointer string = "&" + t.TranslateC(n.Parent)
            var multiargs string
            for _, arg := range n.Args {
                multiargs = ", "
                argsList = append(argsList, t.TranslateC(arg))
            }
            
            //fmt.Println("Static func call: ", n.Name, n.Static)
            if n.Static == true {
                parPointer = ""
                multiargs = ""
            }
        	return fmt.Sprintf("%s_%s(%s%s%s)", n.Class, n.Name, parPointer, multiargs, strings.Join(argsList, ", "))
        
        case *NodeExprGetter:
            class := n.Class
            var varName string
            var varParam []string
            
            parPointer := t.TranslateC(n.Expr)
            
        	for k, v := range n.Vars {
                varName = k
                for _, e := range v {
                    varParam = append(varParam, t.TranslateC(e))
                }
            }
            if len(varParam) != 0 {parPointer += ", "}
            return fmt.Sprintf("%s%s_get(&%s%s)", class, varName, parPointer, strings.Join(varParam, ", "))
            return ""
        
        case *NodeExprSetter:
            class := n.Class
            var varName string
            var varParam []string
            
            parPointer := t.TranslateC(n.Expr)
            
        	for k, v := range n.Vars {
                varName = k
                for _, e := range v {
                    varParam = append(varParam, t.TranslateC(e))
                }
            }
            if len(varParam) != 0 {
                parPointer += ", "
            }
            return fmt.Sprintf("%s%s_set(&%s%s)", class, varName, parPointer, strings.Join(varParam, ", "))
            return ""    
        
        case *NodeStaticStmt:
            t.static = true
            stmt := t.TranslateC(n.Stmt)
            t.static = false
            return stmt
            
        case *NodeStmtTypeDef:
        	typeData := n.Type
            typeName := n.Name
            code := ""
            var Type string = t.matchType(typeData.Lexeme)
            if typeData.tokenType == STRUCT {
                code += " {\n"
                for _, init := range n.Vars {
                    code += "    " + t.TranslateC(init) + "\n"
                }
                code += "}"
                Type = "struct"
            }
        	t.CglobalVars += fmt.Sprintf("typedef %s%s %s;\n", Type, code, typeName)
            return ""
        case *NodeStmtClass:
        	className := n.Name
            var classVars string
            var classCode string
            t.currentClass = className
            for _, e := range n.Code {
                if _, ok := e.(*NodeStmtVar); ok {
                    classVars += "    " + t.TranslateC(e) + "\n"
                } else {
                    classCode += t.TranslateC(e) + "\n"
                }
            }
            t.currentClass = ""
            t.CglobalVars += fmt.Sprintf("typedef struct {\n%s} %s;\n%s", classVars, className, classCode)
            return ""
        
        case *NodeStmtPkg:
            t.Package = n.Name
            return ""
        
        case *NodeImport:
            var importCode string
            for _, pkg := range n.Names {
                importCode += "#include " + pkg + "\n"
                fmt.Println("Your pack", pkg)
            }
            t.importField += importCode
            return ""
        
        default: return ""
    }
}

func (t *Transpiler) GenerateCCode(nodes []Node) {
    var CBuilder strings.Builder
    var HBuilder strings.Builder
    CBuilder.WriteString("#include <stdlib.h>\n#include <stdio.h>\n#include <stdint.h>\n#include <stdbool.h>\n#include \"src/compiler/runtime.h\"\n")
    var mainContents string
    var headerContents string
    for _, node := range nodes {
        //fmt.Println(fmt.Sprintf("Node: %s of type %T", t.TranslateC(node), node))
        line := t.TranslateC(node)
        head := t.TranslateH(node)
        //fmt.Println("Line: ", line)
        //fmt.Println("Head: ", head)
        mainContents += fmt.Sprintf("%s\n", line)
        headerContents += fmt.Sprintf("%s\n", head)
    }
    fmt.Println("global funcs", t.HglobalFuncs)
    CBuilder.WriteString(fmt.Sprintf("%s\n\n%s%s%s", t.importField, t.CglobalVars, t.CglobalFuncs, mainContents))
    HBuilder.WriteString(fmt.Sprintf("#ifndef %s_H\n#define %s_H\n", strings.ToUpper(t.fileName), strings.ToUpper(t.fileName)))
    HBuilder.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n#endif", t.importField, t.HglobalVars, t.HglobalFuncs, headerContents))
    t.WriteInFile(HBuilder.String(), ".h")
	t.WriteInFile(CBuilder.String(), ".c")
}