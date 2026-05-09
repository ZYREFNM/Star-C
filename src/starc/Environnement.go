package main

import (
    "fmt"
)

type Environnement struct {
    Package string
    Import map[string]bool
    Type map[string]string
    Variable map[string]string
    Pointer map[string]bool
    Func map[string]string
    Static map[string]bool
    Classes map[string]*Environnement
    Parent *Environnement
}

func InitEnvi() *Environnement {
    return &Environnement{
        Package: "main",
        Import: make(map[string]bool),
        Type: make(map[string]string),
        Variable: make(map[string]string),
        Pointer: make(map[string]bool),
        Func: make(map[string]string),
        Static: make(map[string]bool),
        Classes: make(map[string]*Environnement),
    }
}

func (e *Environnement) NewScope() *Environnement {
    return &Environnement{
        Package: e.Package,
        Import: make(map[string]bool),
        Type: make(map[string]string),
        Variable: make(map[string]string),
        Pointer: make(map[string]bool),
        Func: make(map[string]string),
        Static: make(map[string]bool),
        Classes: make(map[string]*Environnement),
        Parent: e,
    }
}

func (e *Environnement) hasType(Type string) bool {
    _, exist := e.Type[Type]
    if !exist && e.Parent != nil {
        return e.Parent.hasType(Type)
    }
    return exist
}

func (e *Environnement) hasVar(Var string) bool {
    _, exist := e.Variable[Var]
    if !exist && e.Parent != nil {
        return e.Parent.hasVar(Var)
    }
    return exist
}

func (e *Environnement) hasFunc(Func string) bool {
    _, exist := e.Func[Func]
    if !exist && e.Parent != nil {
        return e.Parent.hasFunc(Func)
    }
    return exist
}

func (e *Environnement) getStatic(Name string, ClassName string) bool {
    fmt.Println("Searchin for static", Name, "in class", ClassName)
    _, exist := e.Static[Name]
    if exist {
        fmt.Println("Found")
        return true
    }
    if envi, ok := e.Classes[ClassName]; ok {
        fmt.Println(fmt.Sprintf("Class %s found", ClassName))
        return envi.getStatic(Name, ClassName)
    }
    if e.Parent != nil {
        fmt.Println("Searching in Parent")
        return e.Parent.getStatic(Name, ClassName)
    }
    return false
}