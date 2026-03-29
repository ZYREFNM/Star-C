package main;

import (
    "fmt"
    "os"
    "log"
    "bufio"
    "errors"
    "strings"
    //"strconv"
)

var hadCompileError, hadRuntimeError,  hadImplementsError bool = false, false, false
var input string
var lineTracker int
var filePath string

func main() {
    log.SetFlags(0);
    log.SetPrefix(">");
    
    fmt.Println(fmt.Sprintf("%s %s %s: %s", "Star-C", VERSION_STATE, "version", VERSION))
    hadImplementsError = false
    if len(os.Args) > 4 {
        PrintError(1)
    }else if len(os.Args) == 4 {
        fmt.Println("Exact number of args")
        filePath = os.Args[len(os.Args)-1]
        if !strings.HasSuffix(filePath, ".sc") {
            PrintError(2)
        }
        runCommand(filePath)
    }else if len(os.Args) >= 2 {
        runPrompt()
    }
    os.Exit(0);
}

func getError(id uint8) (uint8, error) {
    var message string = "";
    var where string = filePath;
    
    message += fmt.Sprintf("%s <id: %d>': ", "Unexpected 'Error", id);
    var char string = "character";
    
    if len(input) >= 1 {
        char = "characters"
    }
    
    switch id {
        case 0: message += "Process Runned Successfully"; hadRuntimeError = true; break;
        case 1: message += fmt.Sprintf("%s\n%s: %v", "Empty Fields or Unrecognized Command\nThe correct command's format is...\n'starc <action> <input>'", "Current args are", os.Args); hadImplementsError = true; break;
        case 2: message += "File type unsupported...\nPlease try again with a '.sc' file format..."; hadImplementsError = true; break;
        case 3: message += fmt.Sprintf("%s %s <%s>", "Unrecognized", char, input); hadRuntimeError = true; break;
        
    }
    message += ";\n"
    if hadCompileError || hadRuntimeError && !hadImplementsError {
        message += fmt.Sprintf("%s %v at:%v", "Error took place in", where, lineTracker);
    } else if hadImplementsError {
        goto end
    }
    
    end:
    err := errors.New(message);
    return id, err;
}

func PrintError(err uint8) {
    _, errMsg := getError(err);
    fmt.Printf("%s", errMsg);
    fmt.Println();
    os.Exit(int(err))
}

func runFile(path string) {
    var bytes, err = os.ReadFile(path);
    if err != nil {
        PrintError(2)
    }
    run(string(bytes));
    
    if hadCompileError {
        os.Exit(65);
    }
    if hadRuntimeError {
        os.Exit(70);
    }
}

func runPrompt() {
    var reader = bufio.NewScanner(os.Stdin);
    for {
        fmt.Print("> ");
        if !reader.Scan() {
            break;
        }
        var line string = reader.Text();
        input = line;
        run(line);
        hadCompileError = false;
    }
}

func compile(source string) {
    fmt.Println("No compiler implemented for now, use run mode");
}
func run(source string) {
    var scanner Scanner;
    scanner = Scanner{source: source, line: 1}
    var tokens []Token = scanner.ScanTokens();
    lineTracker = scanner.line
    
    for _, token := range tokens {
        fmt.Println(token);
    }
}

func ErrorReport(id uint8) {
    PrintError(id);
    hadCompileError = true;
}

func runCommand(arg string) {
    command := os.Args[2]
    switch command {
        case "ignite": compile(arg); break;
        case "version": fmt.Println(fmt.Sprintf("%s %d", "Star-C version", VERSION))
        default:
        	fmt.Println(command)
        	PrintError(1);
    }
    fmt.Println(os.Args)
}