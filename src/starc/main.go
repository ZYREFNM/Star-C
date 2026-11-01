package main

import (
    "fmt"
    "strconv"
    "os"
    "os/exec"
    "path/filepath"
    "log"
    "errors"
    "strings"
)

var CurrentParser *Parser
var ErrorCatcher *ErrorHandler = &ErrorHandler{Name: "terminal"}
var HasCompileError bool = false
var input string
var extraInput []string
var filePath string
var subFiles []string
var subNames []string
var mainFile string

func main() {
    log.SetFlags(0);
    log.SetPrefix(">");
    
    fmt.Println(fmt.Sprintf("\n%s %s %s: %s", "Star-C", VERSION_STATE, "version", VERSION))
    if len(os.Args) >= 5 {
        filePath = os.Args[3]
        if !IsSCFile(filePath) {
            CatchError("Try to run a .starc file")
        }
        if os.Args[4] != "orbit" || len(os.Args) < 5 {CatchError("Expected orbit or commets")}
        for _, e := range os.Args[4:len(os.Args)] {
            if e == "orbit" {continue}
            subFiles = append(subFiles, e)
        }
        subFiles = append(subFiles, "std/Math.starc", "std/IO.starc", "std/Vector.starc")
        runCommand(filePath, subFiles)
    }else if len(os.Args) == 2{
    }else if len(os.Args) == 4 {
        filePath = os.Args[3]
        filePath = os.Args[len(os.Args)-1]
        if !strings.HasSuffix(filePath, ".starc") {
            CatchError("Try to run a .starc file")
        }
        PrintError()
        runCommand(filePath, subFiles)
    }
    os.Exit(0);
}

func IsSCFile(file string) bool {
    return strings.HasSuffix(file, ".starc")
}

func CatchError(output string) {
    err := errors.New(output);
    ErrorCatcher.Errors = append(ErrorCatcher.Errors, err)
    if CurrentParser == nil { return }
    ErrorCatcher.Name = ErrorCatcher.CurrentParser.filepath
    ErrorCatcher.ErrorLine = append(ErrorCatcher.ErrorLine, ErrorCatcher.CurrentParser.peek(0).Line)
}

func PrintError() {
    if len(ErrorCatcher.Errors) == 0 {return}
    fmt.Printf("%sStarC%s %sfailed!%s\n", Yellow, Reset, BrightPurple, Reset)
    var lines []int = ErrorCatcher.ErrorLine
    var absPath string = ErrorCatcher.Name
    if HasCompileError {fmt.Println("Compilation error!! :")}
    fmt.Printf("Catched %s%d%s errors!!!\n", Red, len(ErrorCatcher.Errors), Reset)
    for i, err := range ErrorCatcher.Errors {
        strLine := " at " + Red + strconv.Itoa(lines[i]) + Reset
        if !HasCompileError {strLine = ""}
        fmt.Printf("%sError%s occured%s in %s: %s\n", Red, Reset, strLine, absPath, err)
    }
    os.Exit(1)
}

func runFile(path string) string {
    var bytes, err = os.ReadFile(path);
    if err != nil {
        CatchError("Invalid filepath")
    }
    PrintError()
    return string(bytes)
}

func ignite(main string, extraSource []string) {
    for _, sub := range subFiles {
        subNames = append(subNames, filepath.Base(sub)[:len(filepath.Base(sub))-6])
    }
    var (
        parsedFiles [][]Node
        filesNames []string
        importedFiles []string
        allSource []string
        envis []*Environnement
        supernova *SuperNova = &SuperNova{
            ModulesPreCompiler: make(map[string]*PreCompiler),
        }
    )
    filesNames = append(filesNames, filePath[:len(filePath)-6])
    HasCompileError = true
    for _, name := range subFiles {
        filesNames = append(filesNames, name[:len(name)-6])
    }
    allSource = append(allSource, main)
    for _, extras := range extraSource {
        allSource = append(allSource, extras)
    }
    fmt.Println("Parsing")
    for i, source := range allSource {
        var (
            scanner Scanner = Scanner{source: source, line: 1}
            tokens []Token = scanner.ScanTokens()
            parser *Parser = &Parser{filepath: filesNames[i], tokens: tokens, current: 0, envi: NewEnvi()}
        )
    	parser.CompTime = NewPreCompiler(parser)
        CurrentParser = parser
        ErrorCatcher.LinkParser(CurrentParser)
    	nodes := parser.Parse()
        if parser.envi != nil {
            supernova.ModulesPreCompiler[parser.Module] = parser.CompTime
            for module, present := range parser.envi.ImportedModules {
                if !present { continue }
                if !strings.HasSuffix(module, "/*\"") { importedFiles = append(importedFiles, module); continue }
                file := module[:len(module)-3] + "\""
                dirs := DirFiles("./" + file[1:len(file)-1])
                for _, element := range dirs {
                    if !IsSCFile(element) { continue }
                    importedFiles = append(importedFiles, file[:len(file)-1] + "/" + element[:len(element)-6] + "\"")
                }
            }
        }
        envis = append(envis, parser.envi)
        parsedFiles = append(parsedFiles, nodes)
    }
    PrintError()
    for _, file := range importedFiles {
        fileName := file
        filesNames = append(filesNames, fileName)
        code := runFile("./" + fileName + ".starc")
        var (
            scanner Scanner = Scanner{source: code, line: 1}
            tokens []Token = scanner.ScanTokens()
            parser *Parser = &Parser{tokens: tokens, filepath: fileName + ".starc", current: 0, envi: NewEnvi()}
        )
        parser.CompTime = NewPreCompiler(parser)
        ErrorCatcher.LinkParser(parser)
        ast := parser.Parse()
        supernova.ModulesPreCompiler[parser.Module] = parser.CompTime
        envis = append(envis, parser.envi)
        parsedFiles = append(parsedFiles, ast)
    }
    PrintError()
    linker := NewLinker()
    for i, _ := range parsedFiles {
        linker.Files = parsedFiles
        linker.FilesEnvi[filesNames[i]] = envis[i]
        fmt.Println(filesNames[i], "->", parsedFiles[i], "-->", envis[i])
        linker.Folder[filesNames[i]] = parsedFiles[i]
    }
    linker.Link()
    for module, ast := range linker.ModuleChecker.ModuleTree {
        name := module.Name
        isMain := false
        if name == filePath[:len(filePath)-6] {isMain = true}
        var compiler Compiler = Compiler{fileName: name, OtherModule: false, isMain: isMain, Nova: supernova,}
        compiler.GenerateCCode(ast)
    }
    PrintError()
    var args []string
    for _, genName := range filesNames {
    	args = append(args, genName + ".c")
    }
    args = append(args, "src/compiler/runtime.c", "src/compiler/arena.c", "-Isrc/compiler", "-o", mainFile)
    cmd := exec.Command("gcc",  args...)
    if err := cmd.Run(); err != nil {
	    fmt.Println("GCC error:", err)
    }
    if !KEEP_C_FILES {
        for _, genFiles := range filesNames {
            os.Remove(genFiles + ".c")
        }
    }
    if !KEEP_H_FILES {
        for _, genFiles := range filesNames {
            os.Remove(genFiles + ".h")
        }
    }
}    

func launch(source string, subSource []string) {
    
}

func runCommand(arg string, subArgs []string) {
    command := os.Args[2]
    var subSource []string
    if len(subArgs) > 0 {
        for _, sub := range subArgs {
            subSource = append(subSource, runFile(sub))
        }
    }
    PrintError()
    switch command {
        case "ignite": mainFile = filepath.Base(filePath)[:len(filepath.Base(filePath))-6]; ignite(runFile(filePath), subSource)
        case "launch": launch(input, subFiles)
        case "version": fmt.Println(fmt.Sprintf("%s %v", "Star-C version", VERSION))
        default:
        	CatchError("Try to run command like this 'starc [flag] <input>'");
    }
    PrintError()
}