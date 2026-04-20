package main

import (
    "fmt"
    "os"
    "os/exec"
    "log"
    "errors"
    "strings"
)

var hadCompileError, hadRuntimeError,  hadImplementsError bool = false, false, false
var input string
var lineTracker int
var errounousLine int
var tokType TokenType
var wordTracker int
var filePath string

func main() {
    log.SetFlags(0);
    log.SetPrefix(">");
    
    fmt.Println(fmt.Sprintf("\n%s %s %s: %s", "Star-C", VERSION_STATE, "version", VERSION))
    hadImplementsError = false
    if len(os.Args) > 4 {
        PrintError(1, "Try to run command like this 'starc [flag] <input>'")
    }else if len(os.Args) == 2{
        fmt.Println("Updated")
    }else if len(os.Args) >= 3 {
        if len(os.Args) < 4 {
            runCommand(filePath)
            PrintError(0, "")
        }
        filePath = os.Args[len(os.Args)-1]
        if !strings.HasSuffix(filePath, ".starc") {
            PrintError(2, "Try to run a .starc file")
        }
        runCommand(filePath)
    }
    os.Exit(0);
}

func getError(id uint8) (uint8, error) {
    var message string = "\n";
    
    message += fmt.Sprintf("%s <id: %d>': ", "Unexpected 'Error", id);
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
        case 6: message += fmt.Sprintf("Unidentified <%s>", input); hadCompileError = true; break
        case 7: message += fmt.Sprintf("Unknown type or object <%s> of type <%v> at:%d", input, tokType, wordTracker); hadCompileError = true; break
        case 8: message += "Missing character"; hadCompileError = true; break
        case 9: message += "Expected value after statement"; hadCompileError = true; break
        case 10: message += fmt.Sprintf("Object %s already exist in current scope or context", input)
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
    
    if hadCompileError {
        os.Exit(65);
    }
    if hadRuntimeError {
        os.Exit(70);
    }
    return string(bytes)
}

func ignite(source string) {
    var scanner Scanner = Scanner{source: source, line: 1}
    var tokens []Token = scanner.ScanTokens()
    lineTracker = scanner.line
    input = scanner.input
    var parser Parser = Parser{tokens: tokens, current: 0, envi: Environnement{Type: make(map[string]string), Variable: make(map[string]any)}}
    for token, _ := range parser.tokens {
        wordTracker = parser.current
        tokType = tokens[token].tokenType
    }
    nodes := parser.Parse()
    var transpiler Transpiler = Transpiler{fileName: "simple"}
    transpiler.GenerateCCode(nodes)
    fmt.Println(transpiler.fileName)
    cmd := exec.Command("gcc", transpiler.fileName + ".c", "src/std/runtime.c", "-Isrc/std", "-o", transpiler.fileName)
    if err := cmd.Run(); err != nil {
	fmt.Println("GCC error:", err)
    }
}    

func launch(source string) {
    
}

func runCommand(arg string) {
    command := os.Args[2]
    switch command {
        case "ignite": ignite(runFile(filePath)); break;
        case "launch": launch(input); break;
        case "version": fmt.Println(fmt.Sprintf("%s %v", "Star-C version", VERSION))
        default:
        	fmt.Println(command)
        	PrintError(1, "Try to run command like this 'starc [flag] <input>'");
    }
}