package main

import (
)
type SymbolKind int
type SymbolScope int

const (
    S_VARIABLE SymbolKind = iota
    S_FUNC
    S_CONST
    S_TYPE
)

const (
    S_PUB SymbolScope = iota
    S_PROTECTED
    S_PRIVATE
    S_UNACESS
)

type CallInfo struct {
    Module string
    FilePath string
    Envi *Environnement
    Line int
    Context *Node
}

type Environnement struct {
    Module string
    FilePath string
    ImportedModules map[string]bool
    Table map[string]*Symbol
    Parent *Environnement
}

func NewEnvi() *Environnement {
    return &Environnement{
        FilePath: ".",
        ImportedModules: make(map[string]bool),
        Table: make(map[string]*Symbol),
        Parent: nil,
    }
}

func (e *Environnement) NewScope() *Environnement {
    return &Environnement{
        Table: make(map[string]*Symbol),
        Parent: e,
    }
}

type Symbol struct {
    Name string
    Kind SymbolKind
    Module string
    Type *NodeType
    Params []Symbol
    Fields []Symbol
    Scope SymbolScope
    Data []Symbol
    Pointer bool
    Static bool
    CreatedAt *CallInfo
    CallAt []CallInfo
}

func (e *Environnement) isGlobal() bool {
    return e.Parent == nil
}

func (e *Environnement) hasVar(mut string) bool {
    for k, v := range e.Table {
        if k == mut {
            if v.Kind == S_VARIABLE {
                return true
            }
        }
    }
    return false
}

func (e *Environnement) getStatics() []string {
    var statics []string
    for k, v := range e.Table {
        if v.Kind == S_FUNC && v.Static {
            statics = append(statics, k)
        }
    }
    return statics
}