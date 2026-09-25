package main

import (
	"fmt"
    "os"
    "strings"
    "path/filepath"
    //"slices"
)

type Compiler struct {
    fileName string
    isMain bool
    currentClass string
    Module string
    OtherModule bool
    importField string
    static bool
    CglobalVars string
    CglobalFuncs string
    HglobalVars string
    HglobalFuncs string
    Nova *SuperNova
    Optimizer AOT_Optimizer
}

type ArenaDesc struct {
    Size string
    RefName string
}

type AOT_Optimizer struct {
    
}

func (c *Compiler) WriteInFile(code string, format string) {
    codeBytes := []byte(code)
    fmt.Println(c.fileName, "is your file")
    os.WriteFile(fmt.Sprintf("%s%s", c.fileName, format), codeBytes, 0644)
}

func (c *Compiler) matchType(Type string, Size string, SubTypes []NodeExpr) string {
    var res string
    var FinalType string = Type
    if Type == "[" {
        if sub, ok := SubTypes[0].(*NodeType); ok {
            FinalType = c.matchType(sub.Type, c.TranslateC(sub.Size), sub.SubType)
            if FinalType == "string" || FinalType == "char*" {FinalType = "char"}
        }
    }
    switch FinalType {
        //case "string": res = "char*"
        case "int8": res = "int8_t"
        case "int16": res = "int16_t"
        case "int32": res = "int32_t"
        case "int64": res = "int64_t"
        case "float32": res = "float"
        case "float64": res = "double"
        default: res = FinalType
    }
    return res
}

func (c *Compiler) matchAction(Action string, Called []NodeExpr, CallerName string) string {
    switch Action {
        case "include":
            var code string
            for _, called := range Called {
                code += "#include " + c.TranslateC(called) + "\n"
            }
            return code
        case "func":
            var args []string
            for _, arg := range Called {
                args = append(args, c.TranslateC(arg))
            }
            return fmt.Sprintf("%s(%s)", CallerName, strings.Join(args, ", "))
        default: return ""
    }
}

func (c *Compiler) basicProperty(Prop string, varName string, varType string, value string, hasParent bool) string {
    var Name string = varName
    if hasParent {
        Name = "this->" + Name
    }
    switch Prop {
        case "get": return fmt.Sprintf("{\nreturn %s;\n}", Name)
        case "set": return fmt.Sprintf("{\nif (value) {\n%s = value;\n}\n}", Name)
        default: return ""
    }
}

func (c *Compiler) matchProperty(Name string, Attributes []any, varName string, varType string, Header bool) string {
    var class string = c.currentClass
    var hasParent bool
    var pointer string
    var point string
    var otherArgs string
    var Property string
    var defaultValue string
    var code string
    if c.currentClass != "" {
        hasParent = true
        class += "_"
        pointer = c.currentClass
        point = "this"
    } else {
        hasParent = false
        pointer = varType
        point = varName
    }
    if len(Attributes) > 1 {
        if params, ok := Attributes[1].([]NodeStmt); ok {
            if len(params) > 2 {
                if paramOne, isOk := params[2].(NodeStmt); isOk {
                    if param, isParam := paramOne.(*NodeStmtVar); isParam {
                        defaultValue = param.Name
                    }
                }
            }
        }
    }
    
    if name, ok := Attributes[0].(string); ok {Property = name}
    if Header {
        code = ";"
    } else {
        if len(Attributes) > 1 {
            if attributes, ok := Attributes[2].([]NodeStmt); ok {
                code += "{\n"
                for _, att := range attributes {
                    code += "    " + c.TranslateC(att) + "\n"
                }
                code += "}\n"
            }
        } else {
            code = c.basicProperty(Property, varName, varType, defaultValue, hasParent)
        }
    }
    if otherArgs == "" {
        otherArgs += fmt.Sprintf("%s value", varType)
    }
    if pointer[len(pointer)-1:len(pointer)] != "*" {
        pointer += "*"
    }
    switch Property {
        case "get": return fmt.Sprintf("%s %s%s_get(%s %s)%s", varType, class, varName, pointer, point, code)
        case "set": return fmt.Sprintf("void %s%s_set(%s %s, %s)%s", class, varName, pointer, point, otherArgs, code)
        default: return ""
    }
}

func (c *Compiler) matchAllocate(alloc string) string {
    switch alloc {
        case "mem": return "arena_malloc"
        case "free": return "arena_free"
        case "clean": return "arena_calloc"
        case "size": return "sizeof"
    	default: return alloc
    }
}

func (c *Compiler) TranslateH(node Node) string {
    switch n := node.(type) {
        
        case *NodeLiteral:
            val := n.Value
            //if val == ""
            return fmt.Sprintf("%v", val)
        
        case *NodeType: return c.matchType(n.Type, c.TranslateH(n.Size), n.SubType)
        
        case *NodeBinary: return fmt.Sprintf("%s %s %s", c.TranslateH(n.Left), n.Operator, c.TranslateH(n.Right))
        
        case *NodeUnary: return fmt.Sprintf("%s (%s)", n.Operator, c.TranslateH(n.Right))
        
        case *NodeGroup: return fmt.Sprintf("(%s)", c.TranslateH(n.Expression))
        
        case *NodeStmtVar:
        	Type := c.TranslateH(n.Type)
            varEnd := ";"
            if n.Value != nil {varEnd = fmt.Sprintf(" = %s;", c.TranslateH(n.Value))}
            if n.Properties != nil {
                c.HglobalFuncs += "\n"
                for key, prop := range n.Properties {
                	c.HglobalFuncs += c.matchProperty(key, prop, n.Name, Type, true) + "\n"
                }
            }
        	code := fmt.Sprintf("%s %s%s\n", Type, n.Name, varEnd)
            if n.Global {
                c.HglobalVars += code
                return ""
            }
            return code
        
        case *NodeStmtMultiVars:
            var vars string
            for _, Var := range n.Vars {
                vars += c.TranslateH(Var)
            }
            return vars
        
        case *NodeStmtConst:
            var Type string = c.TranslateH(n.Type)
            c.HglobalVars += fmt.Sprintf("extern const %s %s;\n", Type, n.Name)
            return ""
        
        case *NodeStaticStmt:
            c.static = true
            stmt := c.TranslateH(n.Stmt)
            c.static = false
            return stmt
        
        case *NodeStmtFuncInit:
        	var funcName string = n.Name
            var paramList []string
            if n.Name == "main" {
                funcName = "star_main"
            }
            Type := c.TranslateH(n.Return)
            
            mod := c.Module
            if mod != "" {
                mod += "__"
            }
            
            receiverType, receiver := n.ReceiverType, n.ReceiverName
            
            if receiverType.Lexeme != "" {
                funcName = receiverType.Lexeme + "_" + funcName
                if !c.static {
                    paramList = append(paramList, fmt.Sprintf("%s* %s", receiverType.Lexeme, receiver.Lexeme))
                }
            }
            
            for _, p := range n.Param {
                param := c.TranslateH(p)
                param = strings.TrimSuffix(param, "\n")
                param = strings.TrimSuffix(param, ";")
                paramList = append(paramList, param)
            }
            c.HglobalFuncs += fmt.Sprintf("%s %s%s(%s);\n\n", Type, mod, funcName, strings.Join(paramList, ", "))
            return ""
            
        case *NodeStmtConstructor:
            var list []string
            Type := c.TranslateC(n.Return)
            for _, p := range n.Param {
                param := c.TranslateH(p)
                param = strings.TrimSuffix(param, "\n")
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            return fmt.Sprintf("%s* %s_new(%s);", Type, Type, strings.Join(list, ", "))
        
        case *NodeStmtTypeDef:
            typeName := n.Name
            code := ""
            var Type string = c.TranslateH(n.Type)
            if Type == "struct" || Type == "enum" {
                code += " {\n"
                for _, init := range n.Vars {
                    lit := c.TranslateC(init)
                    if Type == "enum" {lit = lit[:len(lit)-1] + ","}
                    code += "    " + lit + "\n"
                }
                code += "}"
            }
        	c.HglobalVars += fmt.Sprintf("typedef %s%s %s;\n\n", Type, code, typeName)
            return ""
        case *NodeImport:
            var importCode string
            for _, pkg := range n.Names {
                importCode += "#include " + pkg + "\n"
                fmt.Println("Your pack", pkg)
            }
            c.importField += importCode
            return ""
        default: return ""
    }
}

func (c *Compiler) TranslateC(node Node) string {
    switch n := node.(type) {
        
        case *NodeStmtExpr:
        	return c.TranslateC(n.Expr) + ";"
        
        case *NodeType:
            return c.matchType(n.Type, c.TranslateC(n.Size), n.SubType)
        
        case *NodeBinary:
        	return fmt.Sprintf("%s %s %s", c.TranslateC(n.Left), n.Operator, c.TranslateC(n.Right))
            
        case *NodeLiteral:
            var val any
            switch n.Type.tokenType {
                case STRING: val = fmt.Sprintf("string_new(%s).lit", n.Value); fmt.Println("DEBUG!!!")
                case NULL: val = "NULL"
                default: val = n.Value
            }
            return fmt.Sprintf("%v", val)
        
        case *NodeUnary:
            op := n.Operator
            right := c.TranslateC(n.Right)
            if op != "-" && op != "!" {
                op = c.matchType(op, "", nil)
                op = "(" + op + ")"
            } else {
                right = "(" + right + ")"
            }
            return fmt.Sprintf("%s%s", op, right)
        
        case *NodeExprConcat:
            left , right := c.TranslateC(n.To), c.TranslateC(n.From)
            left = strings.TrimSuffix(left, ".lit")
            right = strings.TrimSuffix(right, ".lit")
        	return fmt.Sprintf("star_concat(%s, %s)", left, right)
            
        case *NodeGroup: return fmt.Sprintf("(%s)", c.TranslateC(n.Expression))
        
        case *NodeGet:
            target := c.TranslateC(n.Object)
            //fmt.Println(fmt.Sprintf("Getting -> %s %s %s", target, symbol, n.Field))
            return fmt.Sprintf("%s%s%s", target, n.Symbol, n.Field)
        
        case *NodePkgResolve:
            c.OtherModule = true
            pkg := n.Pkg
            if pkg == "" {
                pkg = c.Module
            }
            res := c.TranslateC(n.Resolution)
            c.OtherModule = false
            return pkg + "__" + res
        
        case *NodeStmtVar:
            var varName string = n.Name
            var Type string
            var nodeSize string
            if node, ok := n.Type.(*NodeType); ok {
                nodeSize = c.TranslateC(node.Size)
                Type = c.matchType(node.Type, nodeSize, node.SubType)
            }
            if nodeSize != "" {
                varName += "[" + nodeSize + "]"
            }
            varEnd := ";"
            if n.Value != nil {
                val := c.TranslateC(n.Value)
                if strings.HasSuffix(val, ".lit") {val = strings.TrimSuffix(val, ".lit")}
                varEnd = fmt.Sprintf(" = %s;", val)
            }
            if n.Properties != nil {
                c.CglobalFuncs += "\n"
                for key, prop := range n.Properties {
                	c.CglobalFuncs += c.matchProperty(key, prop, n.Name, Type, false) + "\n"
                }
            }
        	code := fmt.Sprintf("%s %s%s\n", Type, varName, varEnd)
            if n.Global {
                c.CglobalVars += code
                return ""
            }
            return code
        
        case *NodeStmtMultiVars:
            var vars string
            for _, Var := range n.Vars {
                vars += c.TranslateC(Var)
            }
            return vars
        
        case *NodeStmtConst:
            var code string
        	var Type string = c.TranslateC(n.Type)
            val := c.TranslateC(n.Value)
            if !MatchStr(Type, "enum", "struct") {
                val += ";"
            }
            if MatchStr(val, "") {
                val = ";"
            } else {val = "= " + val}
            code = fmt.Sprintf("const %s %s %s", Type, n.Name, val)
            if n.Global {
                c.CglobalVars += code
                return ""
            }
            return code
            
        case *NodeAssignment:
            val := c.TranslateC(n.Value)
            if strings.HasSuffix(val, ".lit") {val = strings.TrimSuffix(val, ".lit")}
            return fmt.Sprintf("%s %s %s;", c.TranslateC(n.Target), n.Symbol.Lexeme, val)
        
        case *NodeVariable: return n.Name
        
        case *NodeExprAlloc:
            alloc := c.matchAllocate(n.Allocation)
            args := c.TranslateC(n.Size)
            if n.Allocation != "size" {args = "global_arena, " + args}
            return fmt.Sprintf("%s(%s)", alloc, args)
        
        case *NodeBlock:
        	var code string = "{\n"
            for _, stmt := range n.Instructions {
                code += "	" + c.TranslateC(stmt) + "\n"
            }
            code += "}\n"
        	return code
            
        case *NodeStmtReturn:
        	if n.Value == nil {return "return;"}
            return fmt.Sprintf("return %s;", c.TranslateC(n.Value))
        
        case *NodeStmtC:
            return c.matchAction(n.Action, n.Called, n.CallerName) + ";"
        
        case *NodeExprC:
            return c.matchAction(n.Action, n.Called, n.CallerName)
        
        case *NodeStmtIf:
        	condition := c.TranslateC(n.Condition)
            result := c.TranslateC(n.Result)
            return fmt.Sprintf("if (%s) %s", condition, result)
        
        case *NodeStmtLoop:
            loops := c.TranslateC(n.Looping)
            result := c.TranslateC(n.Result)
            return  fmt.Sprintf("for %s %s", loops, result)
        
        case *NodeExprIter:
            var Var string
            var Type string
            if node_var, ok := n.Iter.(*NodeStmtVar); ok {
                Var = node_var.Name
                Type = c.TranslateC(node_var.Type) + " "
            } else {Var = c.TranslateC(n.Iter)}
            if Var == "" {Var = "i"; Type = "int "}
            Min := c.TranslateC(n.Min)
            if Min == "" {Min = "0"}
            Max := c.TranslateC(n.Max)
            Increm := c.TranslateC(n.Increm)
            if Increm == "" {Increm = "1"}
            Operator := n.Operator
            if Operator == "" {Operator = "<"}
            return fmt.Sprintf("(%s%s = %s; %s %s %s; %s = %s + %s)", Type, Var, Min, Var, Operator, Max, Var, Var, Increm)
        
        case *NodeStmtWhile:
        	condition := c.TranslateC(n.Condition)
            result := c.TranslateC(n.Result)
            return fmt.Sprintf("while (%s) %s", condition, result)
        
        case *NodeStmtBreak:
            return "break;"
            
        case *NodeStmtFuncInit:
        	var funcName string = n.Name
            var pack string
        	var list []string
            receiverType, receiver := n.ReceiverType, n.ReceiverName
            if receiverType.Lexeme != "" {
                funcName = receiverType.Lexeme + "_" + funcName
                if !c.static {
                    list = append(list, fmt.Sprintf("%s* %s", receiverType.Lexeme, receiver.Lexeme))
                }
            }
            if c.Module != "" {
                if funcName != "main" {
                    pack = c.Module + "__"
                }
            }
        	for _, p := range n.Param {
                param := c.TranslateC(p)
                param = strings.TrimSuffix(param, "\n")
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            fmt.Println("Your param", strings.Join(list, ", "))
            code := c.TranslateC(n.Code)
            if n.Name == "main" {
                funcName = "star_main"
                code = code[:len(code)-2] + "return 0;\n}\n"
            }
        	return fmt.Sprintf("%s %s%s(%s) %s", c.TranslateC(n.Return), pack, funcName, strings.Join(list, ", "), code)
        
        case *NodeStmtConstructor:
            var list []string
            var code string
            receiverType, receiver := n.ReceiverType.Lexeme, n.ReceiverName.Lexeme
            Type := c.TranslateC(n.Return)
            for _, p := range n.Param {
                param := c.TranslateC(p)
                param = strings.TrimSuffix(param, "\n")
                param = strings.TrimSuffix(param, ";")
                list = append(list, param)
            }
            code = c.TranslateC(n.Code)
            code = code[:2] + fmt.Sprintf("	%s* %s = malloc(sizeof(%s));\n", receiverType, receiver, Type) + code[2:len(code)-2] + fmt.Sprintf("\n	return %s; \n}", receiver)
            return fmt.Sprintf("%s* %s_new(%s) %s", Type, Type, strings.Join(list, ", "), code)
        case *NodeExprFuncCall:
        	var argsList []string
            var module string
            if c.OtherModule == false {module = c.Module + "__"}
            for _, arg := range n.Args {
                trans_arg := c.TranslateC(arg)
                if strings.HasSuffix(trans_arg, ".lit") {
                    trans_arg = strings.TrimSuffix(trans_arg, ".lit")
                }
                argsList = append(argsList, trans_arg)
            }
        	return fmt.Sprintf("%s%s(%s)", module, n.Name, strings.Join(argsList, ", "))
        
        case *NodeExprMethodCall:
            var class string = n.Class
        	var argsList []string
            var parPointer string = "&" + c.TranslateC(n.Parent)
            var multiargs string
            for _, arg := range n.Args {
                multiargs = ", "
                trans_arg := c.TranslateC(arg)
                if strings.HasSuffix(trans_arg, ".lit") {
                    trans_arg = strings.TrimSuffix(trans_arg, ".lit")
                }
                argsList = append(argsList, trans_arg)
            }
            
            fmt.Println("NAME", n.Name)
            
            if n.Name == "new" {
                class = "*" + class
            } else if n.Name == "new*" {
                n.Name = "new"
            }
            
            if n.Static == true || n.Name == "new"{
                parPointer = ""
                multiargs = ""
            }
        	return fmt.Sprintf("%s_%s(%s%s%s)", class, n.Name, parPointer, multiargs, strings.Join(argsList, ", "))
        
        case *NodeExprGetter:
            var varName string
            var varParam []string
            
            parPointer := c.TranslateC(n.Expr)
            
        	for k, v := range n.Vars {
                varName = k
                if parPointer == "" {
                    parPointer = varName
                }
                for _, e := range v {
                    varParam = append(varParam, c.TranslateC(e))
                }
            }
            if len(varParam) != 0 {parPointer += ", "}
            return fmt.Sprintf("%s_get(&%s%s)", varName, parPointer, strings.Join(varParam, ", "))
            return ""
        
        case *NodeStmtSetter:
            var varName string
            var varParam []string
            
            parPointer := c.TranslateC(n.Expr)
            
        	for k, v := range n.Vars {
                varName = k
                if parPointer == "" {
                    parPointer = varName
                }
                for _, e := range v {
                    varParam = append(varParam, c.TranslateC(e))
                }
            }
            if len(varParam) != 0 {
                parPointer += ", "
            }
            return fmt.Sprintf("%s_set(&%s%s)", varName, parPointer, strings.Join(varParam, ", "))
            return ""    
        
        case *NodeExprListIndex:
            return fmt.Sprintf("%s[%s]", c.TranslateC(n.Expr), c.TranslateC(n.Index))
        
        case *NodeStaticStmt:
            c.static = true
            stmt := c.TranslateC(n.Stmt)
            c.static = false
            return stmt
        case *NodeScopeAcces:
            var code string
            for _, stmt := range n.Stmt {
                code += c.TranslateC(stmt) + "\n"
            }
            return code
        case *NodeStmtTypeDef:
            c.currentClass = n.Name
            typeName := n.Name
            code := ""
            var Type string = c.TranslateC(n.Type)
            if Type == "struct" || Type == "enum" {
                code += " {\n"
                for _, init := range n.Vars {
                    lit := c.TranslateC(init)
                    if Type == "enum" {lit = lit[:len(lit)-1] + ","}
                    code += "    " + lit + "\n"
                }
                code += "}"
            }
        	c.CglobalVars += fmt.Sprintf("typedef %s%s %s;\n\n", Type, code, typeName)
            c.currentClass = ""
            return ""
        case *NodeStmtModule:
            c.Module = filepath.Base(n.Name)
            fmt.Println(Red, "Module de salope", Reset, filepath.Base(n.Name), c.Module)
            return ""
        
        case *NodeImport:
            var Modules []string
            for _, module := range n.Names {
                fmt.Println(module, module[len(module)-5:len(module)], "MODULE", module[len(module)-5:len(module)] == "/*.h\"")
                if strings.HasSuffix(module, "/*.h\"") {
                    module = strings.TrimPrefix(module, "\".")
                    dir := module[:len(module)-5]
                    fmt.Println("Yes", dir)
                    files := DirFiles(dir)
                    for _, file := range files {
                        file = "\"" + module[:len(module)-4] + file[1:len(file)-1] + "\""
                        fmt.Println("TES FICHIERS", file)
                        if !strings.HasSuffix(file, ".starc") {continue}
                        fmt.Println("TES FICHIERS APRÈS", file)
                        Modules = append(Modules, strings.TrimSuffix(file, ".starc"))
                    }
                } else {fmt.Println("Non"); Modules = append(Modules, module)}
            }
            var importCode string
            for _, module := range Modules {
                importCode = "#include " + module + "\n"
                fmt.Println("Your pack", module)
            }
            c.importField += importCode
            return ""
        
        default: return ""
    }
}

func (c *Compiler) GenerateCCode(nodes []Node) {
    var (
        CBuilder strings.Builder
        HBuilder strings.Builder
        mainContents string
        headerContents string
        NovaArena *ModuleArena = c.Nova.ModulesPreCompiler[c.Module].Arena
    )
    for _, node := range nodes {
        line := c.TranslateC(node)
        head := c.TranslateH(node)
        mainContents += fmt.Sprintf("%s\n", line)
        headerContents += fmt.Sprintf("%s\n", head)
    }
    var BasicMacros string
    if c.isMain {
        BasicMacros = fmt.Sprintf("#define __ARENA_ALLOCATOR_SIZE__ %s\n", NovaArena.Size)
    }
    CBasicImport := "#include <stdlib.h>\n#include <stdio.h>\n#include <stdint.h>\n#include <stdbool.h>\n#include \"../src/compiler/runtime.h\"\n" + BasicMacros
    CBuilder.WriteString(CBasicImport)
    HBuilder.WriteString(CBasicImport)
    var new_global_arena string
    var mainFunc string
    if c.isMain {
        new_global_arena = fmt.Sprintf("t_Arena* %s;\n", NovaArena.LocalName)
        mainFunc = fmt.Sprintf("int main() {\n%s = new_arena(__ARENA_ALLOCATOR_SIZE__ * 1024);\nint exit = star_main();\nkill_arena(%s);\nreturn exit;\n}\n", NovaArena.LocalName, NovaArena.LocalName)
    }
    CBuilder.WriteString(new_global_arena)
    CBuilder.WriteString(fmt.Sprintf("%s\n\n%s%s%s%s", c.importField, c.CglobalVars, c.CglobalFuncs, mainContents, mainFunc))
    HBuilder.WriteString(fmt.Sprintf("\n#ifndef %s_H\n#define %s_H\n", strings.ToUpper(filepath.Base(c.fileName)), strings.ToUpper(filepath.Base(c.fileName))))
    HBuilder.WriteString(fmt.Sprintf("%s\n%s\n%s\n%s\n#endif", c.importField, c.HglobalVars, c.HglobalFuncs, headerContents))
    c.WriteInFile(HBuilder.String(), ".h")
	c.WriteInFile(CBuilder.String(), ".c")
}