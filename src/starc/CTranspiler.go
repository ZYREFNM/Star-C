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
}

func (t *Transpiler) WriteInFile(code string) {
    codeBytes := []byte(code)
    os.WriteFile(fmt.Sprintf("%s.c", t.fileName), codeBytes, 0644)
}

func (t *Transpiler) matchType(Type string) string {
    switch Type {
        case "string": return "char*"
        case "int8": return "int8_t"
        case "int16": return "int16_t"
        case "int32": return "int32_t"
        case "int64": return "int64_t"
        case "float32": return "float"
        case "float64": return "double"
        default: return Type
    }
}

func (t *Transpiler) matchAction(Action string) string {
    switch Action {
        case "include": return "#include "
        case "func-call": return ""
        default: return ""
    }
}

func (t *Transpiler) Translate(node Node) string {
    switch n := node.(type) {
        
        case *NodeStmtExpr:
        	return t.Translate(n.Expr) + ";"
        
        case *NodeBinary:
        	return fmt.Sprintf("%s %s %s", t.Translate(n.Left), n.Operator, t.Translate(n.Right))
            
        case *NodeLiteral: return fmt.Sprintf("%v", n.Value)
        
        case *NodeUnary: return fmt.Sprintf("%s (%s)", n.Operator, t.Translate(n.Right))
        
        case *NodeExprConcat:
        	return fmt.Sprintf("star_concat(%s, %s)", t.Translate(n.To), t.Translate(n.From))
            
        case *NodeGroup: return fmt.Sprintf("(%s)", t.Translate(n.Expression))
        
        case *NodeGet:
        	symbol := "."
            target := t.Translate(n.Object)
        	if target == "this" {symbol = "->"}
            fmt.Println(fmt.Sprintf("Getting -> %s %s %s", target, symbol, n.Field))
            return fmt.Sprintf("%s%s%s", target, symbol, n.Field)
        
        case *NodeStmtVar:
        	Type := t.matchType(n.Type.Lexeme)
            varEnd := ";"
            if n.Value != nil {varEnd = fmt.Sprintf(" = %s;", t.Translate(n.Value))}
        	return fmt.Sprintf("%s %s%s", Type, n.Name, varEnd)
            
        case *NodeAssignment: return fmt.Sprintf("%s = %s;", t.Translate(n.Target), t.Translate(n.Value))
        
        case *NodeVariable: return n.Name
        
        case *NodeBlock:
        	var code string = "{\n"
            for _, stmt := range n.Instructions {
                code += "	" + t.Translate(stmt) + "\n"
            }
            code += "}"
        	return code
            
        case *NodeStmtPrint:
        	var list string
            var format string
        	for _, expr := range n.Expressions {
                list += fmt.Sprintf("%s, ", t.Translate(expr))
                format += "%d"
            }
            list = list[0:len(list)-2]
        	return fmt.Sprintf("printf(\"%s\\n\", %s);", format, list)
            
        case *NodeStmtReturn:
        	if n.Value == nil {return "return;"}
            return fmt.Sprintf("return %s;", t.Translate(n.Value))
        case *NodeStmtC:
        	var action string = t.matchAction(n.Action)
        	var list []string
        	for _, call := range n.Called {
                list = append(list, action + t.Translate(call)[1:len(t.Translate(call)) - 1])
            }
            return strings.Join(list, ";\n")
        
        case *NodeStmtIf:
        	condition := t.Translate(n.Condition)
            result := t.Translate(n.Result)
            return fmt.Sprintf("if (%s) %s", condition, result)
            
        case *NodeStmtWhile:
        	condition := t.Translate(n.Condition)
            result := t.Translate(n.Result)
            return fmt.Sprintf("while (%s) %s", condition, result)
            
        case *NodeStmtFuncInit:
        	var funcName string = n.Name
        	var list []string
            if t.currentClass != "" {
                funcName = t.currentClass + "_" + funcName
                list = append(list, fmt.Sprintf("%s* this", t.currentClass))
            }
        	for _, p := range n.Param {
                param := t.Translate(p)
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            code := t.Translate(n.Code)
        	return fmt.Sprintf("%s %s(%s) %s", t.matchType(n.Return), funcName, strings.Join(list, ", "), code)
            
        case *NodeExprFuncCall:
        	var argsList []string
            for _, arg := range n.Args {
                argsList = append(argsList, t.Translate(arg))
            }
        	return fmt.Sprintf("%s(%s)", n.Name, strings.Join(argsList, ", "))
        
        case *NodeExprMethodCall:
        	var argsList []string
            var multiargs string
            for _, arg := range n.Args {
                multiargs = ", "
                argsList = append(argsList, t.Translate(arg))
            }
        	return fmt.Sprintf("%s_%s(&%s%s%s)", n.Class, n.Name, multiargs, t.Translate(n.Parent), strings.Join(argsList, ", "))
            
        case *NodeStmtTypeDef:
        	typeData := n.Type
            typeName := n.Name
            code := ""
            var Type string = t.matchType(typeData.Lexeme)
            if typeData.tokenType == STRUCT {
                code += " {\n"
                for _, init := range n.Vars {
                    code += "	" + t.Translate(init) + "\n"
                }
                code += "}"
                Type = "struct"
            }
        	return fmt.Sprintf("typedef %s%s %s\n", Type, code, typeName)
        case *NodeStmtClass:
        	className := n.Name
            t.currentClass = className
            var classTypes string
            var classVars string
            var classFuncs string
            for _, e := range n.Vars {
                classVars += "	" + t.Translate(e) + "\n"
            }
            for _, e := range n.TypeDef {
                classTypes += "	" + t.Translate(e) + "\n"
            }
            for _, e := range n.Func {
                classFuncs += t.Translate(e) + "\n"
            }
            t.currentClass = ""
            return fmt.Sprintf("typedef struct {\n%s} %s;\n%s\n%s", classVars, className, classTypes, classFuncs)
        
        default: return ""
    }
}

func (t *Transpiler) GenerateCCode(nodes []Node) {
    var CBuilder strings.Builder
    CBuilder.WriteString("#include <stdio.h>\n#include <stdint.h>\n#include \"src/std/runtime.h\"\n\n")
    var mainContents string
    for _, node := range nodes {
        //fmt.Println(fmt.Sprintf("Node: %s of type %T", t.Translate(node), node))
        line := t.Translate(node)
        mainContents += fmt.Sprintf("%s\n", line)
    }
    CBuilder.WriteString(fmt.Sprintf("%s", mainContents))
    
	t.WriteInFile(CBuilder.String())
}