package main

type Environnement struct {
    Type map[string]*NodeType
    Variable map[string]string
    Parent *Environnement
}