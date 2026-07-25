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
    Statics map[string]bool
    Protected map[string]string
    ModuleChecker *ModuleManager
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
            statics := slices.Collect(maps.Keys(envi.Static))
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

func (m *ModuleManager) assimilateModules() {
    
    //Modules assimilation
    for _, envi := range m.Parent.FilesEnvi {
        m.KnownModules[&Module{
            Name: envi.Module,
            Envi: envi,
            Ast: m.Parent.Folder[envi.Module],
        }] = envi.ModulePath
    }
    
    //Module Checking
    /*for _, envi := range m.Parent.FilesEnvi {
        
    }*/
}

func (m *ModuleManager) checkModule() {
}

func (m *ModuleManager) getPathOfModule(module string) (path string) {
    return m.Parent.FilesEnvi[module].ModulePath
}

func (m *ModuleManager) resetKnownModules() {
    for mod := range m.KnownModules {
        delete(m.KnownModules, mod)
    }
}

func (l *Linker) checkScope() {
    var protectList []string
    for _, envi := range l.FilesEnvi {
        protectList = envi.getProtected()
        for _, protect := range protectList {
            l.Protected[protect] = envi.Module
            fmt.Println("Protected:", protect, envi.Module)
        }
    }
}

func (l *Linker) addComplementModules() {
    return
}

func (m *ModuleManager) AssembleModule(ModuleStruct *Module, SubModule *Module) {
    for name, isImport := range SubModule.Envi.Import {
        ModuleStruct.Envi.Import[name] = isImport
    }
    for file, ast := range m.Parent.Folder {
        if ast == nil {continue}
        for _, node := range ast {
            switch n := node.(type) {
                case *NodeStmtModule, *NodeImport: break
                default: m.Parent.Folder[file] = append(m.Parent.Folder[file], n)
            }
        }
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
        for mod2, _ := range m.KnownModules {
            if mod2 == nil {continue}
            if mod1.Name == mod2.Name && (mod1.Name != "") && (mod1.Name != ".") {
                m.AssembleModule(mod1, mod2)
                delete(m.KnownModules, mod2)
            }
        }
    }
}

func (l *Linker) Link() {
    l.ModuleChecker.assimilateModules()
    l.addComplementModules()
    l.ModuleChecker.fuseModule()
    l.ModuleChecker.checkModule()
    l.checkStatic()
    l.checkScope()
}

func (l *Linker) GetLink() [][]Node {
    var LinkedFiles [][]Node
    for _, file := range l.Files {
        LinkedFiles = append(LinkedFiles, file)
    }
    return LinkedFiles
}