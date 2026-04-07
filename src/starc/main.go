package main

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
var tokType TokenType
var wordTracker int
var filePath string

func main() {
    log.SetFlags(0);
    log.SetPrefix(">");
    
    fmt.Println(fmt.Sprintf("%s %s %s: %s", "Star-C", VERSION_STATE, "version", VERSION))
    hadImplementsError = false
    if len(os.Args) > 4 {
        PrintError(1)
    }else if len(os.Args) == 2{
        runPrompt()
    }else if len(os.Args) >= 3 {
        if len(os.Args) < 4 {
            runCommand(filePath)
            PrintError(0)
        }
        filePath = os.Args[len(os.Args)-1]
        if !strings.HasSuffix(filePath, ".starc") {
            PrintError(2)
        }
        runCommand(filePath)
    }
    os.Exit(0);
}

func getError(id uint8) (uint8, error) {
    var message string = "";
    var where string = filePath;
    
    message += fmt.Sprintf("%s <id: %d>': ", "Unexpected 'Error", id);
    var char string = "characters";
    
    if id == 0 {
        message = ""
    }
    switch id {
        case 0: message += "Process Runned Successfully"; hadImplementsError = true; break;
        case 1: message += fmt.Sprintf("%s\n%s: %v", "Empty Fields or Unrecognized Command\nThe correct command's format is...\n'starc <action> <input>'", "Current args are", os.Args); hadImplementsError = true; break
        case 2: message += "File type unsupported...\nPlease try again with a '.starc' file format ...\nexemple: file.starc"; hadImplementsError = true; break
        case 3: message += fmt.Sprintf("%s %s <%s>", "Unrecognized", char, input); hadRuntimeError = true; break
        case 4: message += "File path invalid... Retry with a working path"; hadImplementsError = true; break
        case 5: message += fmt.Sprintf("Not a valid expression '%s'...", input); hadCompileError = true ; break
        case 6: message += fmt.Sprintf("Unidentified <%s>", input); hadCompileError = true; break
        case 7: message += fmt.Sprintf("Unknown type or object <%s> of type <%v> at:%d", input, tokType, wordTracker); hadCompileError = true; break
        case 8: message += "Missing semi-colon"; hadCompileError = true; break
        case 9: message += "Expected value after statement"; hadCompileError = true; break
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
    id, errMsg := getError(err);
    fmt.Println(errMsg)
    os.Exit(int(id))
}

func runFile(path string) string {
    var bytes, err = os.ReadFile(path);
    if err != nil {
        PrintError(4)
    }
    
    if hadCompileError {
        os.Exit(65);
    }
    if hadRuntimeError {
        os.Exit(70);
    }
    return string(bytes)
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
        run(line)
        fmt.Println()
        hadCompileError = false;
    }
}

func compile(source string) {
    var scanner Scanner = Scanner{source: source, line: 1}
    var tokens []Token = scanner.ScanTokens()
    var parser Parser = Parser{tokens: tokens, current: 0}
    for token, _ := range parser.tokens {
        wordTracker = parser.current
        tokType = tokens[token].tokenType
    }
    lineTracker = scanner.line
    nodes := parser.Parse()
    var transpiler Transpiler = Transpiler{fileName: "simple"}
    transpiler.GenerateCCode(nodes)
    PrintError(0)
    
}
func run(source string) {
}

func launch(source string) {
    
}

func runCommand(arg string) {
    command := os.Args[2]
    switch command {
        case "ignite": compile(runFile(filePath)); break;
        case "launch": launch(input); break;
        case "version": fmt.Println(fmt.Sprintf("%s %v", "Star-C version", VERSION))
        default:
        	fmt.Println(command)
        	PrintError(1);
    }
}