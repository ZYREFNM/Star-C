package main

import (
    "fmt"
)

type SuperNova struct {
    ModulesPreCompiler map[string]*PreCompiler
}

type PreCompiler struct {
    ParseMachine *Parser
    Arena *ModuleArena
    ExternArenas []ModuleArena
    Table *MacroTable
}

type ModuleArena struct {
    Size string
    LocalName string
    Policy string
    Fallback string
    State string
    Module string
}

type MacroTable struct {
    Types map[string]string
}

func NewPreCompiler(parser *Parser) *PreCompiler {
    return &PreCompiler{
        ParseMachine: parser,
        Arena: &ModuleArena{
            Size: "1000",
            LocalName: "__STAR_C_ARENA_ALLOCATOR__",
            Fallback: "safe",
            State: "release",
        },
        Table: &MacroTable{
            Types: make(map[string]string),
        },
    }
}

func (p *PreCompiler) ParsePreCompile() {
    p.ParseMachine.advance()
    p.ParseMachine.expect("Pre-Compiling scope, expected left brace '{'", p.ParseMachine.peek(0), LEFT_BRACE)
    
    for !p.ParseMachine.peek(0).tokenType.equals(RIGHT_BRACE) {
        if p.ParseMachine.isAtEnd() {CatchError("PreCompiler reached EOF while reading file"); return}
        
        fmt.Println("Caller:", p.ParseMachine.peek(0))
        switch p.ParseMachine.peek(0).Lexeme {
            
            case "arena":
            fmt.Println("Is here", p.ParseMachine.peek(0))
            p.ParseMachine.advance()
            p.ParseMultiCaller("arena")
                
            default: CatchError("Pre-Compiler unknown caller " + p.ParseMachine.peek(-1).Lexeme); return
        }
    }
    p.ParseMachine.advance()
}

func (p *PreCompiler) ParseCaller(name string) {
    fmt.Println("Call:", p.ParseMachine.peek(0))
    left := p.ParseMachine.expect("Expected identifier", p.ParseMachine.peek(0), IDENTIFIER)
    p.ParseMachine.expect(fmt.Sprintf("Expected colon in %s caller, found %s", name, p.ParseMachine.peek(0).Lexeme), p.ParseMachine.peek(0), COLON)
    fmt.Println("Field:", p.ParseMachine.peek(0))
    right := p.ParseMachine.expect("Expected value for property of " + name + " arena caller", p.ParseMachine.peek(0), STRING, INT, FLOAT, IDENTIFIER)
    
    switch name {
        
        case "arena":
            switch left.Lexeme {
                case "size": p.Arena.Size = right.Lexeme
                case "local_name": p.Arena.LocalName = right.Lexeme
                case "policy": p.Arena.Policy = right.Lexeme
                case "fallback": p.Arena.Fallback = right.Lexeme
                case "state": p.Arena.State = right.Lexeme
                default:
            }
        
        default:
    }
    
}

func (p *PreCompiler) ParseMultiCaller(name string) {
    fmt.Println("Left", p.ParseMachine.peek(0))
    p.ParseMachine.expect("Expected left parenthesize", p.ParseMachine.peek(0), LEFT_PAREN)
    for !p.ParseMachine.peek(0).tokenType.equals(RIGHT_PAREN) {
        fmt.Println("Calling:", p.ParseMachine.peek(0))
        if p.ParseMachine.isAtEnd() { CatchError("Reached EOF before right parenthesize ')'"); return}
        p.ParseCaller(name)
        if p.ParseMachine.peek(0).tokenType.equals(COMMA) {p.ParseMachine.advance()}
    }
    p.ParseMachine.advance()
}