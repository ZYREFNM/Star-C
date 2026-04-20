package main

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
    return t >= ARRAY_TYPE && t <= STRING_TYPE
}

func (t TokenType) isDigitType() bool {
    return t >= INT_TYPE && t <= INT64_TYPE || t >= UINT_TYPE && t <= UINT64_TYPE || t >= FLOAT_TYPE && t <= FLOAT64_TYPE || t >= UFLOAT_TYPE && t <= UFLOAT64
}

func (t TokenType) isBoolOperator() bool {
    return t >= BANG_EQUAL && t <= LESS_EQUAL && t != EQUAL
}

func (t TokenType) isAssignOperator() bool {
    return t == EQUAL || t >= PLUS_EQUAL && t <= SLASH_EQUAL
}

func (t TokenType) isModifier() bool {
    return t >= PUBLIC && t <= SET
}

func (t TokenType) isVarMod() bool {
    return t >= GET && t <= SET
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
    PLUS_EQUAL
    MINUS_EQUAL
    STAR_EQUAL
    SLASH_EQUAL
    RIGHT_ARROW
    CONCAT
    
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
    
    //Data-Type
    ARRAY_TYPE
    DICTIONARY_TYPE
    IDENTIFIER_TYPE
    INT_TYPE
    INT8_TYPE
    INT16_TYPE
    INT32_TYPE
    INT64_TYPE
    FLOAT_TYPE
    FLOAT8_TYPE
    FLOAT16_TYPE
    FLOAT32_TYPE
    FLOAT64_TYPE
    UINT_TYPE
    UINT8_TYPE
    UINT16_TYPE
    UINT32_TYPE
    UINT64_TYPE
    UFLOAT_TYPE
    UFLOAT8_TYPE
    UFLOAT16_TYPE
    UFLOAT32_TYPE
    UFLOAT64_TYPE
    STRING_TYPE
    
    //Keywords
    AND
    BREAK
    CONST
    CLASS
    DEFAULT
    ELSE
    FALSE
    FUNC
    FOR
    GOTO
    IF
    NOT
    NULL
    OR
    OPTIONAL
    PACKAGE
    PRINT
    RETURN
    SUPER
    STRUCT
    SWITCH
    THIS
    TRUE
    TYPEDEF
    VOID
    VAR
    WHILE
    WITH
    
    // Modifiers' key-words
    PUBLIC
    PRIVATE
    GET
    SET
    
    EOF
)