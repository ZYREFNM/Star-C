package main

//Virtual Star-C's IR

type SCIR_Manager struct {
    ParsedModules []Module
    CFG Graph
}

type Graph struct {
    Head *BasicBlock
}

type BasicBlock struct {
    ID int
    ContextId int
    ContextKind BBContext
    Seq_Code []*Instruction
}

type Data interface {
    isData()
}

type BBContext int
type InstKind int
type OPKind int

//Basic Blocks context
const (
    FUNC_CONTEXT BBContext = iota
)
//Instructions kind
const (
    BINARY InstKind = iota
    CHECK_TYPE
)
const (
    ADD OPKind = iota
)

type Instruction struct {
    Action InstKind
    Left *Data
    Operator OPKind
    Right *Data
}

func (s *SCIR_Manager) LowerAstToCFG() {
}