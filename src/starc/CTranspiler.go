package main

import (
	"fmt"
    "os"
    "strings"
    "path/filepath"
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
        case *NodeStmtVar: return fmt.Sprintf("%s %s;", n.Type, n.Name)
        case *NodeVariable: return n.Name
        default: return ""
    }
}

func (t *Transpiler) GenerateCCode(nodes []Node) {
    var CBuilder strings.Builder
    CBuilder.WriteString("#include <stdio.h>\n\n")
    var mainContents string
    for _, node := range nodes {
        line := t.Translate(node)
        mainContents += fmt.Sprintf("%s\n", line)
    }
    CBuilder.WriteString(fmt.Sprintf("int main() {\n%s\n}", mainContents))
    
	t.WriteInFile(CBuilder.String())
    fmt.Println(fmt.Sprintf("File %s.c", t.fileName))
    fmt.Println(filepath.Abs(fmt.Sprintf("%s.c", t.fileName)))
}