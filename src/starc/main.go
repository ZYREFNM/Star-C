package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "log"
    "errors"
    "strings"
)

var hadCompileError, hadRuntimeError,  hadImplementsError bool = false, false, false
var input string
var extraInput []string
var lineTracker int
var errounousLine int
var tokType TokenType
var wordTracker int
var filePath string
var subFiles []string
var subNames []string
var mainFile string

func main() {
    log.SetFlags(0);
    log.SetPrefix(">");
    
    fmt.Println(fmt.Sprintf("\n%s %s %s: %s", "Star-C", VERSION_STATE, "version", VERSION))
    hadImplementsError = false
    if len(os.Args) >= 4 {
        filePath = os.Args[3]
        if !strings.HasSuffix(filePath, ".starc") {
            PrintError(2, "Try to run a .starc file")
        }
        if os.Args[4] != "orbit" || len(os.Args) < 5 {PrintError(0, "Expected orbit or commets")}
        for _, e := range os.Args[4:len(os.Args)] {
            if e == "orbit" {continue}
            //fmt.Println("E:", e, e[:len(e)-6])
            //fmt.Println("Est-ce vide ?", subFiles)
            subFiles = append(subFiles, e)
            //fmt.Println("Regarde ici", e, len(os.Args), os.Args)
        }
        runCommand(filePath, subFiles)
    }else if len(os.Args) == 2{
        fmt.Println("Updated")
    }else if len(os.Args) == 3 {
        if len(os.Args) < 4 {
            runCommand(filePath, subFiles)
            PrintError(0, "")
        }
        filePath = os.Args[len(os.Args)-1]
        if !strings.HasSuffix(filePath, ".starc") {
            PrintError(2, "Try to run a .starc file")
        }
        runCommand(filePath, subFiles)
    }
    os.Exit(0);
}

func getError(id uint8) (uint8, error) {
    var message string = "\n";
    
    message += fmt.Sprintf("Unexpected Error <id: %d>': ", id);
    var char string = "characters";
    
    if id == 0 {
        message = ""
    }
    switch id {
        case 0: message += "Process Runned Successfully."; break;
        case 1: message += fmt.Sprintf("%s\n%s: %v", "Empty Fields or Unrecognized Command.\n"); hadImplementsError = true; break
        case 2: message += "File type unsupported..."; hadImplementsError = true; break
        case 3: message += fmt.Sprintf("Unrecognized %s <%s>", char, input); hadRuntimeError = true; break
        case 4: message += "File path invalid... Retry with a working path."; hadImplementsError = true; break
        case 5: message += fmt.Sprintf("Not a valid expression '%s'...", input); hadCompileError = true ; break
        case 6: message += fmt.Sprintf("Unidentified %s", input); hadCompileError = true; break
        case 7: message += fmt.Sprintf("Unknown type or object <%s> of type <%v> at:%d", input, tokType, wordTracker); hadCompileError = true; break
        case 8: message += "Missing character"; hadCompileError = true; break
        case 9: message += "Expected value after statement"; hadCompileError = true; break
        case 10: message += fmt.Sprintf("Object %s already exist in current scope or context", input); break
        case 11: message += "Missing implementation"; break
        case 12: message += "PlaceHolder error for $callers"; break
    }
    message += "\n"
    err := errors.New(message);
    return id, err;
}

func PrintError(err uint8, hint string) {
    errounousLine = lineTracker
    var where string = filePath;
    id, errMsg := getError(err);
    fmt.Println(errMsg)
    if hadCompileError || hadRuntimeError && !hadImplementsError {
        fmt.Println(fmt.Sprintf("Error took place in %v at line:%v", where, errounousLine))
    }
    fmt.Println("Hint: ", hint, ".")
    os.Exit(int(id))
}

func runFile(path string) string {
    var bytes, err = os.ReadFile(path);
    if err != nil {
        PrintError(4, "Check your file's path and if it's correct")
    }
    mainFile = filepath.Base(path)[:len(filepath.Base(path))-6]
    if hadCompileError {
        os.Exit(65);
    }
    if hadRuntimeError {
        os.Exit(70);
    }
    return string(bytes)
}

func ignite(main string, extraSource []string) {
    //extra := len(extraSource) < 0;
    for _, sub := range subFiles {
        //fmt.Println("Sub is", sub)
        subNames = append(subNames, filepath.Base(sub)[:len(filepath.Base(sub))-6])
    }
    var parsedFiles [][]Node
    var filesNames []string
    filesNames = append(filesNames, mainFile)
    //fmt.Println("Os args", os.Args)
    //fmt.Println("Sub files", subNames)
    for _, name := range subNames {
        filesNames = append(filesNames, name)
        //fmt.Println("Sub: ", name)
    }
    var allSource []string
    var envis []*Environnement
    allSource = append(allSource, main)
    for _, extras := range extraSource {
        allSource = append(allSource, extras)
    }
    for _, source := range allSource {
        var scanner Scanner = Scanner{source: source, line: 1}
        var tokens []Token = scanner.ScanTokens()
        lineTracker = scanner.line
        input = scanner.input
        var parser Parser = Parser{tokens: tokens, current: 0, envi: InitEnvi()}
        for token, _ := range parser.tokens {
            wordTracker = parser.current
            tokType = tokens[token].tokenType
    	}
    	nodes := parser.Parse()
        envis = append(envis, parser.envi)
        //fmt.Println("Les funcs de mon envi\n", parser.envi.Func, "\n")
        parsedFiles = append(parsedFiles, nodes)
        //fmt.Println(i, ".passage")
    }
    linker := NewLinker()
    for i, _ := range parsedFiles {
        linker.Files = parsedFiles
        linker.FilesEnvi[filesNames[i]] = envis[i]
    }
    linker.Link()
    //fmt.Println("Taille du parsedFiles", len(allSource))
    for i, file := range linker.Files {
        name := filesNames[i]
        //fmt.Println("Name:: ", name)
        var transpiler Transpiler = Transpiler{fileName: name}
        transpiler.GenerateCCode(file)
        //fmt.Println("Aide moi", transpiler.fileName, name)
    }
    var args []string
    for _, genName := range filesNames {
    	args = append(args, genName + ".c")
    }
    args = append(args, "src/compiler/runtime.c", "-Isrc/compiler", "-o", mainFile)
    cmd := exec.Command("gcc",  args...)
    if err := cmd.Run(); err != nil {
	    fmt.Println("GCC error:", err)
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
    //fmt.Println("Dans commande, les subs sont", len(subArgs), len(subSource))
    switch command {
        case "ignite": ignite(runFile(filePath), subSource); break;
        case "launch": launch(input, subFiles); break;
        case "version": fmt.Println(fmt.Sprintf("%s %v", "Star-C version", VERSION))
        default:
        	fmt.Println(command)
        	PrintError(1, "Try to run command like this 'starc [flag] <input>'");
    }
}