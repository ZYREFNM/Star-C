package main

type Node interface {
    isANode()
}

type NodeBinary struct {
    Left Node
    Operator string
    Right Node
}

func (n *NodeBinary) isANode() {}

type NodeGroup struct {
    Expression Node
}

func (n *NodeGroup) isANode() {}

type NodeLiteral struct {
    Value any
}

func (n *NodeLiteral) isANode() {}

type NodeUnary struct {
    Operator string
    Right Node
}

func (n *NodeUnary) isANode() {}

type NodeVariable struct {
    Name string
}

func (n *NodeVariable) isANode() {}

type NodeAssignement struct {
    Name string
    Value Node
}

func (n *NodeAssignement) isANode() {}

type NodeBlock struct {
    Expression Node
}

func (n *NodeBlock) isANode() {}