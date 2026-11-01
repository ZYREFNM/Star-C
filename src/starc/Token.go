package main

type TokenType int

type Token struct {
    tokenType TokenType;
    Lexeme string;
    Literal any;
    Line int;
}

func (t TokenType) isDigit() bool {
    return t == INT
}

func (t TokenType) isInteger() bool {
    return t == INT || t == UINT
}

func (t TokenType) isType() bool {
    return t >= ARRAY_TYPE && t <= STRUCT_TYPE
}

func (t TokenType) isDigitType() bool {
    return t >= INT_TYPE && t <= INT64_TYPE || t >= UINT_TYPE && t <= UINT64_TYPE || t >= FLOAT_TYPE && t <= FLOAT64_TYPE || t >= UFLOAT_TYPE && t <= UFLOAT64_TYPE
}

func (t TokenType) isBoolOperator() bool {
    return t >= BANG_EQUAL && t <= LESS_EQUAL
}

func (t TokenType) isAssignOperator() bool {
    return t == EQUAL || t >= PLUS_EQUAL && t <= SLASH_EQUAL
}

func (t TokenType) isModifier() bool {
    return t >= PUBLIC && t <= SET
}

func (t TokenType) isAction() bool {
    return t >= CALL && t <= INCLUDE
}

func (t TokenType) isStmtCall() bool {
    return t >= WRITE && t <= INCLUDE
}

func (t TokenType) isExprCall() bool {
    return (t >= CALL && t <= FUNCCALL) || t == FUN
}

func (t TokenType) isProperty() bool {
    return t >= GET && t <= SET
}

func (t TokenType) equals(Types ...TokenType) bool {
    for _, token := range Types {
        if token == t {return true} else {continue}
    }
    return false
}

const (
    //Single char tokens
    LEFT_PAREN TokenType = iota
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    LEFT_BRACKET
    RIGHT_BRACKET
    COLON
    COMMA
    DOT
    DOLLAR
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR
    HASHTAG
    AT_SIGN
    
    //One or more char tokens
    BANG
    EQUAL
    BANG_EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL
    PLUS_EQUAL
    MINUS_EQUAL
    STAR_EQUAL
    SLASH_EQUAL
    DASH_R_ARROW
    FAT_R_ARROW
    CONCAT
    VAR_ARGS
    TYPE_CONV
    CHAINED_LINK
    SCOPE_RESOLVE
    
    //Literals
    ARRAY
    BOOL
    BYTE
    ENUM
    IDENTIFIER
    INT
    FLOAT
    MAP
    UINT
    UFLOAT
    STRING
    STRUCT
    
    //Data-Type
    ARRAY_TYPE
    BYTE_TYPE
    BOOL_TYPE
    ENUM_TYPE
    MAP_TYPE
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
    STRUCT_TYPE
    
    //Keywords
    ALLOCATE
    AND
    ARENA
    BREAK
    CCALL
    CONST
    CLASS
    DEFAULT
    DEFINE
    ELSE
    EXTENDS
    FALSE
    FUN
    FOR
    GOTO
    IF
    IMPORT
    IOTA
    LOOP
    MODULE
    NOT
    NULL
    OR
    OPTIONAL
    PRECOMP
    RETURN
    SUPER
    STATIC
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
    
    //C-Caller's actions
    CALL
    FUNCCALL
    //------
    WRITE
    INCLUDE
    
    //Memory management keywords
    HEAP
    FREE
    SIZEOF
    LENOF
    EMPTY
    
    //Precomp reserved words
    
    EOF
)