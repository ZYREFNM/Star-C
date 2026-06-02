package main

import (
    //"fmt"
)

type Environnement struct {
    Package string
    Import map[string]bool
    Type map[string]string
    Variable map[string]string
    Const map[string]string
    Pointer map[string]bool
    Func map[string]string
    Static map[string]bool
    Classes map[string]*Environnement
    Parent *Environnement
    Unknown *UnknownLocker
}

type UnknownLocker struct {
    Type map[string]string
    Func map[string]string
    Variable map[string]string
    Const map[string]string
    Classes map[string]*Environnement
    Pointer map[string]bool
}

func InitEnvi() *Environnement {
    return &Environnement{
        Package: "main",
        Import: make(map[string]bool),
        Type: make(map[string]string),
        Variable: make(map[string]string),
        Const: make(map[string]string),
        Pointer: make(map[string]bool),
        Func: make(map[string]string),
        Static: make(map[string]bool),
        Classes: make(map[string]*Environnement),
        Unknown: InitUnknownLocker(),
    }
}

func InitUnknownLocker() *UnknownLocker {
    return &UnknownLocker{
        Type: make(map[string]string),
        Func: make(map[string]string),
        Variable: make(map[string]string),
        Const: make(map[string]string),
        Classes: make(map[string]*Environnement),
        Pointer: make(map[string]bool),
    }
}

func (e *Environnement) NewUnknownLocker() *UnknownLocker {
    return InitUnknownLocker()
}

func (e *Environnement) NewScope() *Environnement {
    return &Environnement{
        Package: e.Package,
        Import: make(map[string]bool),
        Type: make(map[string]string),
        Variable: make(map[string]string),
        Const: make(map[string]string),
        Pointer: make(map[string]bool),
        Func: make(map[string]string),
        Static: make(map[string]bool),
        Classes: make(map[string]*Environnement),
        Unknown: e.NewUnknownLocker(),
        Parent: e,
    }
}

func (e *Environnement) isGlobal() bool {
    return e.Parent == nil
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

func (e *Environnement) hasConst(Const string) bool {
    _, exist := e.Const[Const]
    if !exist && e.Parent != nil {
        return e.Parent.hasConst(Const)
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

func (e *Environnement) addStatic(Name string) {
    if !e.isGlobal() {
        e.Parent.addStatic(Name)
    } else {
        e.Static[Name] = true
    }
    return
}

func (e *Environnement) getStatic(Name string, ClassName string) bool {
    //fmt.Println("Searchin for static", Name, "in class", ClassName)
    _, exist := e.Static[Name]
    if exist {
        //fmt.Println("Found")
        return true
    }
    if envi, ok := e.Classes[ClassName]; ok {
        //fmt.Println(fmt.Sprintf("Class %s found", ClassName))
        return envi.getStatic(Name, ClassName)
    }
    if e.Parent != nil {
        //fmt.Println("Searching in Parent")
        return e.Parent.getStatic(Name, ClassName)
    }
    return false
}

func (e *Environnement) SearchImport() []string {
    var packs []string
    if e.Import != nil {
        for pack := range e.Import {
            packs = append(packs, pack)
        }
    }
    return packs
}

func (e *Environnement) SearchUnknownPack() []string {
    var packs []string
    if e.Unknown != nil {
        if e.Unknown.Type != nil {
            for unknown := range e.Unknown.Type {
                packs = append(packs, unknown)
            }
        }
        if e.Unknown.Func != nil {
            for unknown := range e.Unknown.Func {
                packs = append(packs, unknown)
            }
        }
        if e.Unknown.Variable != nil {
            for unknown := range e.Unknown.Variable {
                packs = append(packs, unknown)
            }
        }
        if e.Unknown.Const != nil {
            for unknown := range e.Unknown.Const {
                packs = append(packs, unknown)
            }
        }
    }
    return packs
}

func (e *Environnement) getUnsureFuncs() []string {
    var unsure []string
    if e.Unknown != nil {
        if e.Unknown.Func != nil {
            for Func := range e.Unknown.Func {
                unsure = append(unsure, Func)
            }
        }
    }
    return unsure
}