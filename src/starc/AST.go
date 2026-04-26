package main

type Node interface {
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


type NodeProperty struct {
    
}
func (n *NodeProperty) isANode() {}
func (n *NodeProperty) isStmt() {}


type NodeGet struct {
    Object Node
    Symbol string
    Field string
}
func (n *NodeGet) isANode() {}
func (n *NodeGet) isExpr() {}


type NodeSet struct {
    Target NodeExpr
    Value NodeExpr
}
func (n *NodeSet) isANode() {}
func (n *NodeSet) isExpr() {}


type NodeBinary struct {
    Left Node
    Operator string
    Right Node
}
func (n *NodeBinary) isANode() {}
func (n *NodeBinary) isExpr() {}


type NodeGroup struct {
    Expression NodeExpr
}
func (n *NodeGroup) isANode() {}
func (n *NodeGroup) isExpr() {}


type NodeLiteral struct {
    Value any
}
func (n *NodeLiteral) isANode() {}
func (n *NodeLiteral) isExpr() {}


type NodeUnary struct {
    Operator string
    Right Node
}
func (n *NodeUnary) isANode() {}
func (n *NodeUnary) isExpr() {}


type NodeExprConcat struct {
    From NodeExpr
    To NodeExpr
}
func (n *NodeExprConcat) isANode() {}
func (n *NodeExprConcat) isExpr() {}

type NodeType struct {
    Type string
}
func (n *NodeType) isANode() {}
func (n *NodeType) isExpr() {}


type NodeVariable struct {
    Name string
}
func (n *NodeVariable) isANode() {}
func (n *NodeVariable) isExpr() {}


type NodeStmtVar struct {
    Name string
    Type Token
    Value Node
}
func (n *NodeStmtVar) isANode() {}
func (n *NodeStmtVar) isStmt() {}


type NodeAssignment struct {
    Target NodeExpr
    Value NodeExpr
}
func (n *NodeAssignment) isANode() {}
func (n *NodeAssignment) isStmt() {}


type NodeStmtC struct {
    Action string
    Called []NodeExpr
}
func (n *NodeStmtC) isANode() {}
func (n *NodeStmtC) isStmt() {}


type NodeBlock struct {
	Instructions []NodeStmt
}
func (n *NodeBlock) isANode() {}
func (n *NodeBlock) isStmt() {}


type NodeStmtReturn struct {
    Value NodeExpr
}
func (n *NodeStmtReturn) isANode() {}
func (n *NodeStmtReturn) isStmt() {}


type NodeStmtPrint struct {
    Expressions []NodeExpr
}
func (n *NodeStmtPrint) isANode() {}
func (n *NodeStmtPrint) isStmt() {}


type NodeStmtIf struct {
    Condition NodeExpr
    Result NodeStmt
}
func (n *NodeStmtIf) isANode() {}
func (n *NodeStmtIf) isStmt() {}

type NodeStmtWhile struct {
    Condition NodeExpr
    Result NodeStmt
}
func (n *NodeStmtWhile) isANode() {}
func (n *NodeStmtWhile) isStmt() {}


type NodeStmtFuncInit struct {
    Return string
    Name string
    Param []NodeStmt
    Code NodeStmt
}
func (n *NodeStmtFuncInit) isANode() {}
func (n *NodeStmtFuncInit) isStmt() {}


type NodeExprFuncCall struct {
    Name string
    Args []NodeExpr
}
func (n *NodeExprFuncCall) isANode() {}
func (n *NodeExprFuncCall) isExpr() {}


type NodeExprMethodCall struct {
    Class string
    Parent NodeExpr
    Name string
    Args []NodeExpr
}
func (n *NodeExprMethodCall) isANode() {}
func (n *NodeExprMethodCall) isExpr() {}


type NodeStmtConstructor struct {
    Return string
    Param []NodeStmt
    Code NodeStmt
}
func (n *NodeStmtConstructor) isANode() {}
func (n *NodeStmtConstructor) isStmt() {}


type NodeStmtTypeDef struct {
    Type Token
    Name string
    Vars []NodeStmt
}
func (n *NodeStmtTypeDef) isANode() {}
func (n *NodeStmtTypeDef) isStmt() {}


type NodeStmtClass struct {
    Name string
    Vars []NodeStmt
    Func []NodeStmt
    TypeDef []NodeStmt
}
func (n *NodeStmtClass) isANode() {}
func (n *NodeStmtClass) isStmt() {}