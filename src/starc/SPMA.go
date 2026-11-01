package main

import (
    //"fmt"
)

type SuperNova struct {
    ModulesPreCompiler map[string]*PreCompiler
}

type Interpreter struct {
    
}

type PreCompiler struct {
    ASTParser *Parser
    Arena *ModuleArena
    ExternArenas []ModuleArena
    Table *MacroTable
    MetaFuncs []*MetaFunc
}

type MetaVoid interface {
    isMeta()
}

type MetaFunc struct {
    Types []*MetaType
    Name Token
    Args map[string]*MetaType
}
func (m *MetaFunc) isMeta() {}

type MetaType struct {
    IsRawData bool
    Type string
}
func (m *MetaType) isMeta() {}

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
        ASTParser: parser,
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

func (p *PreCompiler) PreCompMetaType() *MetaType {
    parser := p.ASTParser
    
    if !parser.peek(0).tokenType.equals(HASHTAG) { return &MetaType{ IsRawData: false, Type: parser.advance().Lexeme } }
    
    if parser.peek(0).Lexeme != "data" { CatchError("Expected raw data caller as '#' was found") }
    parser.expect("Expected left bracket for raw data caller", parser.peek(0), LEFT_BRACKET)
    
    ExpectedData := parser.advance().Lexeme
    parser.expect("Expected right bracket for raw data caller", parser.peek(0), RIGHT_BRACKET)
    
    return &MetaType{ IsRawData: true, Type: ExpectedData }
}

func (p *PreCompiler) PreCompile(tok Token) {
    parser := p.ASTParser
    parser.advance()
    
    parser.expect("Expected left brace", parser.peek(0), LEFT_BRACE)
    for !parser.peek(0).tokenType.equals(RIGHT_BRACE) {
        switch tok.Lexeme {
            case "func": p.PreCompFunc()
            case "arena": p.PreCompArena()
            default:
        }
        if parser.isAtEnd() { CatchError("Expected right brace before EOF") }
    }
}

func (p *PreCompiler) PreCompArena() {
    parser := p.ASTParser
    
    parser.expect("Expected left parenthesize", parser.peek(0), LEFT_PAREN)
    
    for !parser.peek(0).tokenType.equals(RIGHT_PAREN) {
        param := parser.advance().Lexeme
        arena := p.Arena
        
        parser.expect("Expected colon", parser.peek(0), COLON)
        val := parser.peek(0).Lexeme
        
        switch param {
            case "size": arena.Size = val
            case "ref": arena.LocalName = val
            case "fallback": arena.Fallback = val
            case "state": arena.State = val
            default:
        }
        if parser.peek(0).tokenType.equals(COMMA) { parser.advance() }
        if parser.isAtEnd() { return }
    }
    parser.advance()
    
}

func (p *PreCompiler) PreCompFunc() {
    parser := p.ASTParser
    MetaTypes := []*MetaType{}
    FuncArgs := make(map[string]*MetaType)
    
    for !parser.peek(0).tokenType.equals(HASHTAG) {
        MetaTypes = append(MetaTypes, p.PreCompMetaType())
        if parser.peek(0).tokenType.equals(CHAINED_LINK) { parser.advance() }
        if parser.isAtEnd() { CatchError("Expected '#'") }
    }
    parser.advance()
    
    FuncName := parser.advance()
    
    
    parser.expect("Expected left parenthesize", parser.peek(0), LEFT_PAREN)
    
    for !parser.peek(0).tokenType.equals(RIGHT_PAREN) {
        argType := p.PreCompMetaType()
        argName := parser.peek(0).Lexeme
        
        FuncArgs[argName] = argType
        if parser.isAtEnd() { CatchError("Reached EOF before meta-function's arg parsing had come to end") }
    }
    parser.advance()
    
    
    //Ended parse reading; Creating a new MetaFunc and sending it
    NewFunc := &MetaFunc{
        Types: MetaTypes,
        Name: FuncName,
        Args: FuncArgs,
    }
    
    p.MetaFuncs = append(p.MetaFuncs, NewFunc)
    
}