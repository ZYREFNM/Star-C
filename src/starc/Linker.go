package main

import (
    "fmt"
    "slices"
    "maps"
)

type Linker struct {
    Files [][]Node
    FilesEnvi map[string]*Environnement
    Folder map[string][]Node
    KnownPackages map[string]bool
    Statics map[string]bool
}

func NewLinker() *Linker {
    return &Linker{
        FilesEnvi: make(map[string]*Environnement),
        KnownPackages: make(map[string]bool),
        Statics: make(map[string]bool),
    }
}

func (l *Linker) MethodLook(node Node) {
    //fmt.Println(fmt.Sprintf("Checking %s of type %T\n", node, node))
    if node == nil {
        return
    }
    
    if method, ok := node.(*NodeExprMethodCall); ok {
        if l.Statics[method.Name] == true {
            method.Static = true
        }
    }
    
    for _, child := range node.Children() {
        l.MethodLook(child)
    }
}

func (l *Linker) checkStatic() {
    for _, envi := range l.FilesEnvi {
        if envi != nil {
            statics := slices.Collect(maps.Keys(envi.Static))
            for _, static := range statics {
                //fmt.Println("Static enregistré est", static)
                l.Statics[static] = true
            }
        }
    }
    for file := range l.Files {
        for stmt := range l.Files[file] {
            l.MethodLook(l.Files[file][stmt])
        }
    }
}

func (l *Linker) checkPkg() {
    for _, envi := range l.FilesEnvi {
        if envi != nil && envi.Package != "" {
            l.KnownPackages[envi.Package] = true
        }
    }
    for _, envi := range l.FilesEnvi {
        if envi == nil {fmt.Println("Skip empty envi"); continue}
        for _, pack := range envi.SearchImport() {
            fmt.Println("Import pack", pack)
            if l.KnownPackages[pack] != true {
                PrintError(5, "Invalid package import " + pack)
            }
        }
    }
    for _, envi := range l.FilesEnvi {
        packList := envi.SearchUnknownPack()
        if envi == nil {continue}
        for _, pack := range packList {
            fmt.Println("Import pack", pack)
            if l.KnownPackages[pack] != true {
                PrintError(5, "Invalid Package " + pack)
            }
        }
    }
}

func (l *Linker) addComplementModules() {
    return
}

func (l *Linker) Link() {
    l.checkPkg()
    l.checkStatic()
    l.addComplementModules()
}

func (l *Linker) GetLink() [][]Node {
    var LinkedFiles [][]Node
    for _, file := range l.Files {
        LinkedFiles = append(LinkedFiles, file)
    }
    return LinkedFiles
}