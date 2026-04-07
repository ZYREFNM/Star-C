package main

import (
	"fmt"
    "os"
    "strings"
    //"path/filepath"
)

type Transpiler struct {
    fileName string
    //parsedFile string
    //current int
}

func (t *Transpiler) WriteInFile(code string) {
    codeBytes := []byte(code)
    os.WriteFile(fmt.Sprintf("%s.c", t.fileName), codeBytes, 0644)
}

func (t *Transpiler) Translate(node Node) string {
    switch n := node.(type) {
        case *NodeBinary:
        	return fmt.Sprintf("%s %s %s", t.Translate(n.Left), n.Operator, t.Translate(n.Right))
        case *NodeLiteral: return fmt.Sprintf("%v", n.Value)
        case *NodeUnary: return fmt.Sprintf("%s (%s)", n.Operator, t.Translate(n.Right))
        case *NodeGroup: return fmt.Sprintf("(%s)", t.Translate(n.Expression))
        case *NodeStmtVar:
        	Type := n.Type.Lexeme
            varEnd := ";"
            if n.Value != nil {varEnd = fmt.Sprintf(" = %s;", t.Translate(n.Value))}
        	if n.Type.tokenType == STRING_TYPE {
                Type = "char*"
                }
        	return fmt.Sprintf("%s %s%s", Type, n.Name, varEnd)
        case *NodeAssignement: return fmt.Sprintf("%s = %s;", n.Name, t.Translate(n.Value))
        case *NodeVariable: return n.Name
        case *NodeBlock:
        	var code string = "{\n"
            for _, stmt := range n.Instructions {
                code += "	" + t.Translate(stmt) + "\n"
                code += "}"
            }
        	return code
        case *NodeStmtPrint:
        	var list string
        	for _, expr := range n.Expressions {
                list += fmt.Sprintf("%s, ", t.Translate(expr))
            }
            list = list[0:len(list)-2]
        	return fmt.Sprintf("printf(\"%%s\\n\", %s);", list)
        case *NodeStmtReturn:
        	if n.Value == nil {return "return;"}
            return fmt.Sprintf("return %s;", t.Translate(n.Value))
        case *NodeStmtIf:
        	condition := t.Translate(n.Condition)
            result := t.Translate(n.Result)
            return fmt.Sprintf("if (%s) %s", condition, result)
        case *NodeStmtWhile:
        	condition := t.Translate(n.Condition)
            result := t.Translate(n.Result)
            return fmt.Sprintf("while (%s) %s", condition, result)
        default: return ""
    }
}

func (t *Transpiler) GenerateCCode(nodes []Node) {
    var CBuilder strings.Builder
    CBuilder.WriteString("#include <stdio.h>\n\n")
    var mainContents string
    for _, node := range nodes {
        line := t.Translate(node)
        mainContents += fmt.Sprintf("	%s\n", line)
    }
    CBuilder.WriteString(fmt.Sprintf("int main() {\n%s\n}", mainContents))
    
	t.WriteInFile(CBuilder.String())
}