package main

import (
    "fmt"
    //"slices"
    //"maps"
)

type Linker struct {
    Files [][]Node
    FilesEnvi map[string]*Environnement
    Folder map[string][]Node
    Statics map[string]bool
    Protected map[string]string
    ModuleChecker *ModuleManager
    SemAnalyser *SemanticAnalyser
}

type ModuleManager struct {
    ModulesPath map[*Module]string
    KnownModules map[*Module]string
    CrossModules map[*Environnement][]string
    ModuleTree map[*Module][]Node
    Parent *Linker
}

type Module struct {
    Name string
    Envi *Environnement
    Ast []Node
    SubModules []Module
}

type SemanticAnalyser struct {
    
}

func NewModuleManager(l *Linker) *ModuleManager {
    return &ModuleManager{
        ModulesPath: make(map[*Module]string),
        KnownModules: make(map[*Module]string),
        CrossModules: make(map[*Environnement][]string),
        ModuleTree: make(map[*Module][]Node),
        Parent: l,
    }
}

func NewLinker() *Linker {
    var linker *Linker
    linker = &Linker{
        FilesEnvi: make(map[string]*Environnement),
        Folder: make(map[string][]Node),
        Statics: make(map[string]bool),
        Protected: make(map[string]string),
    }
    linker.ModuleChecker = NewModuleManager(linker)
    return linker
}

func (l *Linker) StaticLook(node Node) {
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
        l.StaticLook(child)
    }
}

func (l *Linker) checkStatic() {
    for _, envi := range l.FilesEnvi {
        if envi != nil {
            statics := envi.getStatics()
            for _, static := range statics {
                fmt.Println("Static enregistré est", static)
                l.Statics[static] = true
            }
        }
    }
    for file := range l.Files {
        for stmt := range l.Files[file] {
            l.StaticLook(l.Files[file][stmt])
        }
    }
}

func (l *Linker) pointSymbol(node Node, mod *Module, index int) {
    fmt.Println("Node is", fmt.Sprintf("%t", node))
    if n, ok := node.(*NodeGet); ok {
        fmt.Println(Red, "Get node: ", Reset, n)
        if VarCall, exist := n.Object.(*NodeVariable); exist {
            if mod.Envi.Table[VarCall.Name].Pointer == true {
                n.Symbol = "->"
                l.ModuleChecker.ModuleTree[mod][index] = n
            }
            fmt.Println("Here 1")
            return
            fmt.Println("Here 2")
        }
        fmt.Println(Red, "Get end: ", Reset, n.Symbol, node.Children())
    } else {
        if len(node.Children()) == 0 {return}
        fmt.Println("Here 1")
        for _, child := range node.Children() {
            if child == nil {return}
            l.pointSymbol(child, mod, index)
        }
    }
}

func (l *Linker) linkPointer() {
    var pointers []*Symbol
    for mod, _ := range l.ModuleChecker.ModuleTree {
        for _, sym := range mod.Envi.Table {
            if sym.Pointer == true {
                pointers = append(pointers, sym)
            }
        }
    }
    fmt.Println(Red, "Get 1", Reset)
    for mod, ast := range l.ModuleChecker.ModuleTree {
        fmt.Println(Red, "Get mod: ", Reset, mod)
        fmt.Println(Red, "Get ast: ", Reset, ast)
        for index, node := range ast {
            fmt.Println(Red, "Node&.co", Reset, node, index, mod)
            l.pointSymbol(node, mod, index)
        }
    }
}

func (m *ModuleManager) assimilateModules() {
    
    //Modules assimilation
    for _, envi := range m.Parent.FilesEnvi {
        fmt.Println("NKM:", envi.Module, envi)
        m.KnownModules[&Module{
            Name: envi.Module,
            Envi: envi,
            Ast: m.Parent.Folder[envi.Module],
        }] = envi.FilePath
    }
    
    //Module Checking
    /*for _, envi := range m.Parent.FilesEnvi {
        
    }*/
}

func (m *ModuleManager) checkModule() {
}

func (m *ModuleManager) getPathOfModule(module string) (path string) {
    return m.Parent.FilesEnvi[module].FilePath
}

func (m *ModuleManager) resetKnownModules() {
}

func (l *Linker) checkScope() {
}

func (l *Linker) addComplementModules() {
    return
}

func (m *ModuleManager) AssembleModule(ModuleStruct *Module, SubModule *Module) {
    var moduleSymbolesNames []string
    for module := range ModuleStruct.Envi.Table {
        moduleSymbolesNames = append(moduleSymbolesNames, module)
    }
    for module, symbol := range SubModule.Envi.Table {
        ModuleStruct.Envi.Table[module] = symbol
        delete(SubModule.Envi.Table, module)
    }
    fmt.Println("Module M1:", m.Parent.Folder)
    for file, ast := range m.Parent.Folder {
        fmt.Println("Module M2:", file, ast)
        if ast == nil {continue}
        for _, node := range ast {
            switch n := node.(type) {
                case *NodeStmtModule, *NodeImport: break
                default: m.Parent.Folder[file] = append(m.Parent.Folder[file], n)
            }
        }
        fmt.Println("Module M3:", file, ast)
        m.ModuleTree[&Module{
            Name: file,
            Envi: m.Parent.FilesEnvi[file],
            Ast: ast,
        }] = ast
        delete(m.Parent.Folder, file)
    }
}

func (m *ModuleManager) fuseModule() {
    for mod1, _ := range m.KnownModules {
        if mod1 == nil {continue}
        fmt.Println("Here, mod 1:", mod1)
        for mod2, _ := range m.KnownModules {
            if mod2 == nil {continue}
            fmt.Println("Here, mod 2:", mod2)
            if mod1.Name == mod2.Name && (mod1.Name != "") && (mod1.Name != ".") {
                fmt.Println("Here, mods:", mod1 == mod2)
                m.AssembleModule(mod1, mod2)
                fmt.Println("Here, mods assembled", m.ModuleTree)
                delete(m.KnownModules, mod2)
            } else {fmt.Println("Not the same", mod1, mod2, mod1 == mod2)}
        }
    }
}

func (l *Linker) linkNodes() {
    l.checkStatic()
    l.checkScope()
    l.linkPointer()
}

func (l *Linker) Link() {
    l.ModuleChecker.assimilateModules()
    l.addComplementModules()
    l.ModuleChecker.fuseModule()
    l.ModuleChecker.checkModule()
    l.linkNodes()
}

func (l *Linker) GetLink() [][]Node {
    var LinkedFiles [][]Node
    for _, file := range l.Files {
        LinkedFiles = append(LinkedFiles, file)
    }
    return LinkedFiles
}