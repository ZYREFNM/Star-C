package main

type Node interface {
    Children() []Node
    isANode()
}

type NodeStmt interface {
    Node
    isStmt()
}

type NodeExpr interface {
    Node
    isExpr()
}

type NodeStmtExpr struct {
    Expr NodeExpr
}
func (n *NodeStmtExpr) isANode() {}
func (n *NodeStmtExpr) isStmt() {}
func (n *NodeStmtExpr) Children() []Node {
    var expr []Node
    expr = append(expr, n.Expr)
    return expr
}


type NodeProperty struct {
}
func (n *NodeProperty) isANode() {}
func (n *NodeProperty) isStmt() {}
func (n *NodeProperty) Children() []Node {return nil}


type NodeGet struct {
    Object Node
    Symbol string
    Field string
}
func (n *NodeGet) isANode() {}
func (n *NodeGet) isExpr() {}
func (n *NodeGet) Children() []Node {
    var child []Node
    child = append(child, n.Object)
    return child
}


type NodeSet struct {
    Target NodeExpr
    Value NodeExpr
}
func (n *NodeSet) isANode() {}
func (n *NodeSet) isExpr() {}
func (n *NodeSet) Children() []Node {
    var child []Node
    child = append(child, n.Target, n.Value)
    return child
}

type NodePkgResolve struct {
    Pkg string
    Resolution NodeExpr
}
func (n *NodePkgResolve) isANode() {}
func (n *NodePkgResolve) isExpr() {}
func (n *NodePkgResolve) Children() []Node {
    var child []Node
    child = append(child, n.Resolution)
    return child
}


type NodeBinary struct {
    Left Node
    Operator string
    Right Node
}
func (n *NodeBinary) isANode() {}
func (n *NodeBinary) isExpr() {}
func (n *NodeBinary) Children() []Node {
    var child []Node
    child = append(child, n.Left, n.Right)
    return child
}


type NodeGroup struct {
    Expression NodeExpr
}
func (n *NodeGroup) isANode() {}
func (n *NodeGroup) isExpr() {}
func (n *NodeGroup) Children() []Node {
    var child []Node
    child = append(child, n.Expression)
    return child
}


type NodeLiteral struct {
    Value any
}
func (n *NodeLiteral) isANode() {}
func (n *NodeLiteral) isExpr() {}
func (n *NodeLiteral) Children() []Node {
    return nil
}


type NodeUnary struct {
    Operator string
    Right Node
}
func (n *NodeUnary) isANode() {}
func (n *NodeUnary) isExpr() {}
func (n *NodeUnary) Children() []Node {
    var child []Node
    child = append(child, n.Right)
    return child
}


type NodeExprConcat struct {
    From NodeExpr
    To NodeExpr
}
func (n *NodeExprConcat) isANode() {}
func (n *NodeExprConcat) isExpr() {}
func (n *NodeExprConcat) Children() []Node {
    var child []Node
    child = append(child, n.From, n.To)
    return child
}

type NodeType struct {
    Type string
}
func (n *NodeType) isANode() {}
func (n *NodeType) isExpr() {}
func (n *NodeType) Children() []Node {return nil}


type NodeVariable struct {
    Name string
}
func (n *NodeVariable) isANode() {}
func (n *NodeVariable) isExpr() {}
func (n *NodeVariable) Children() []Node {
    return nil
}


type NodeStmtVar struct {
    Name string
    Properties map[string][]any
    Type Token
    Value Node
    Global bool
}
func (n *NodeStmtVar) isANode() {}
func (n *NodeStmtVar) isStmt() {}
func (n *NodeStmtVar) Children() []Node {
    var child []Node
    child = append(child, n.Value)
    return child
}



type NodeStmtConst struct {
    Name string
    Type Token
    Value Node
    Global bool
}
func (n *NodeStmtConst) isANode() {}
func (n *NodeStmtConst) isStmt() {}
func (n *NodeStmtConst) Children() []Node {
    var child []Node
    child = append(child, n.Value)
    return child
}


type NodeAssignment struct {
    Target NodeExpr
    Value NodeExpr
}
func (n *NodeAssignment) isANode() {}
func (n *NodeAssignment) isStmt() {}
func (n *NodeAssignment) Children() []Node {
    var child []Node
    child = append(child, n.Target, n.Value)
    return child
}


type NodeExprAlloc struct {
    Allocation string
    Size NodeExpr
}
func (n *NodeExprAlloc) isANode() {}
func (n *NodeExprAlloc) isExpr() {}
func (n *NodeExprAlloc) Children() []Node {
    var child []Node
    child = append(child, n.Size)
    return child
}


type NodeStmtC struct {
    Action string
    Called []NodeExpr
    CallerName string
}
func (n *NodeStmtC) isANode() {}
func (n *NodeStmtC) isStmt() {}
func (n *NodeStmtC) Children() []Node {
    var child []Node
    for _, call := range n.Called {
        child = append(child, call)
    }
    return child
}


type NodeBlock struct {
	Instructions []NodeStmt
}
func (n *NodeBlock) isANode() {}
func (n *NodeBlock) isStmt() {}
func (n *NodeBlock) Children() []Node {
    var child []Node
    for _, inst := range n.Instructions {
        child = append(child, inst)
    }
    return child
}


type NodeStmtReturn struct {
    Value NodeExpr
}
func (n *NodeStmtReturn) isANode() {}
func (n *NodeStmtReturn) isStmt() {}
func (n *NodeStmtReturn) Children() []Node {
    var child []Node
    child = append(child, n.Value)
    return child
}


type NodeStmtIf struct {
    Condition NodeExpr
    Result NodeStmt
}
func (n *NodeStmtIf) isANode() {}
func (n *NodeStmtIf) isStmt() {}
func (n *NodeStmtIf) Children() []Node {
    var child []Node
    child = append(child, n.Condition, n.Result)
    return child
}

type NodeStmtLoop struct {
    Looping NodeExpr
    Result NodeStmt
}
func (n *NodeStmtLoop) isANode() {}
func (n *NodeStmtLoop) isStmt() {}
func (n *NodeStmtLoop) Children() []Node {
    var child []Node
    child = append(child, n.Looping, n.Result)
    return child
}

type NodeStmtWhile struct {
    Condition NodeExpr
    Result NodeStmt
}
func (n *NodeStmtWhile) isANode() {}
func (n *NodeStmtWhile) isStmt() {}
func (n *NodeStmtWhile) Children() []Node {
    var child []Node
    child = append(child, n.Condition, n.Result)
    return child
}


type NodeStmtFuncInit struct {
    Return string
    Name string
    Param []NodeStmt
    Code NodeStmt
}
func (n *NodeStmtFuncInit) isANode() {}
func (n *NodeStmtFuncInit) isStmt() {}
func (n *NodeStmtFuncInit) Children() []Node {
    var child []Node
    for _, param := range n.Param {
        child = append(child, param)
    }
    child = append(child, n.Code)
    return child
}


type NodeExprFuncCall struct {
    Name string
    Args []NodeExpr
}
func (n *NodeExprFuncCall) isANode() {}
func (n *NodeExprFuncCall) isExpr() {}
func (n *NodeExprFuncCall) Children() []Node {
    var child []Node
    for _, arg := range n.Args {
        child = append(child, arg)
    }
    return child
}


type NodeExprMethodCall struct {
    Class string
    Object string
    Parent NodeExpr
    Name string
    Args []NodeExpr
    Static bool
}
func (n *NodeExprMethodCall) isANode() {}
func (n *NodeExprMethodCall) isExpr() {}
func (n *NodeExprMethodCall) Children() []Node {
    var child []Node
    child = append(child, n.Parent)
    for _, arg := range n.Args {
        child = append(child, arg)
    }
    return child
}


type NodeStmtConstructor struct {
    Return string
    Param []NodeStmt
    Code NodeStmt
}
func (n *NodeStmtConstructor) isANode() {}
func (n *NodeStmtConstructor) isStmt() {}
func (n *NodeStmtConstructor) Children() []Node {
    var child []Node
    for _, param := range n.Param {
        child = append(child, param)
    }
    child = append(child, n.Code)
    return child
}


type NodeStaticStmt struct {
    Stmt NodeStmt
}
func (n *NodeStaticStmt) isANode() {}
func (n *NodeStaticStmt) isStmt() {}
func (n *NodeStaticStmt) Children() []Node {
    var child []Node
    child = append(child, n.Stmt)
    return child
}


type NodeScopeAcces struct {
    Modifier Token
    Stmt NodeStmt
}
func (n *NodeScopeAcces) isANode() {}
func (n *NodeScopeAcces) isStmt() {}
func (n *NodeScopeAcces) Children() []Node {
    var child []Node
    child = append(child, n.Stmt)
    return child
}

type NodeStmtTypeDef struct {
    Type Token
    Name string
    Vars []NodeStmt
}
func (n *NodeStmtTypeDef) isANode() {}
func (n *NodeStmtTypeDef) isStmt() {}
func (n *NodeStmtTypeDef) Children() []Node {
    var child []Node
    for _, Var := range n.Vars {
        child = append(child, Var)
    }
    return child
}


type NodeStmtClass struct {
    Name string
    Code []NodeStmt
}
func (n *NodeStmtClass) isANode() {}
func (n *NodeStmtClass) isStmt() {}
func (n *NodeStmtClass) Children() []Node {
    var child []Node
    for _, code := range n.Code {
        child = append(child, code)
    }
    return child
}


type NodeStmtPkg struct {
    Name string
}
func (n *NodeStmtPkg) isANode() {}
func (n *NodeStmtPkg) isStmt() {}
func (n *NodeStmtPkg) Children() []Node {
    return nil
}


type NodeImport struct {
    Names []string
}
func (n *NodeImport) isANode() {}
func (n *NodeImport) isStmt() {}
func (n *NodeImport) Children() []Node {
    return nil
}