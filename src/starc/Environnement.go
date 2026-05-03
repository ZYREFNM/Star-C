package main


type Environnement struct {
    Type map[string]string
    Variable map[string]string
    Pointer map[string]bool
    Func map[string]string
    Parent *Environnement
}

func (e *Environnement) NewScope() *Environnement {
    return &Environnement{
        Type: make(map[string]string),
        Variable: make(map[string]string),
        Pointer: make(map[string]bool),
        Func: make(map[string]string),
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