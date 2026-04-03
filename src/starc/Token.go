package main

import (
)



type TokenType int

type Token struct {
    tokenType TokenType;
    Lexeme string;
    Literal any;
    Line int;
}

func (t TokenType) isDigit() bool {
    return t >= INT && t <= UFLOAT64
}

func (t TokenType) isInteger() bool {
    return t >= INT && t <= INT64 || t >= UINT && t <= UINT64
}

func (t TokenType) isType() bool {
    return t >= ARRAY && t <= STRING
}

const (
    //Single char tokens
    LEFT_PAREN TokenType = iota
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    LEFT_BRACKET
    RIGHT_BRACKET
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR
    
    //One or two char tokens
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL
    
    //Literals
    ARRAY
    DICTIONARY
    IDENTIFIER
    INT
    INT8
    INT16
    INT32
    INT64
    FLOAT
    FLOAT8
    FLOAT16
    FLOAT32
    FLOAT64
    UINT
    UINT8
    UINT16
    UINT32
    UINT64
    UFLOAT
    UFLOAT8
    UFLOAT16
    UFLOAT32
    UFLOAT64
    STRING
    
    //Keywords
    AND
    BREAK
    CLASS
    CONST
    DEFAULT
    ELSE
    FALSE
    FUNC
    FOR
    GOTO
    IF
    INNER
    NOT
    NULL
    OR
    OPTIONAL
    PACKAGE
    PRINT
    PUBLIC
    PRIVATE
    RETURN
    SUPER
    SWITCH
    THIS
    TRUE
    VAR
    WHILE
    
    EOF
)