package main

type Environnement struct {
    Type map[string]string
    Variable map[string]any
    Func map[string]any
    Parent *Environnement
}

func (e *Environnement) hasType(Type string) bool {
    _, exist := e.Type[Type]
    return exist
}

func (e *Environnement) hasVar(Var string) bool {
    _, exist := e.Variable[Var]
    return exist
}