package main

import (
    "strconv"
)

type Scanner struct {
    source string;
    tokens []Token;
    start int;
    current int;
    line int;
    input string
}

var dataType = map[string]TokenType{
	"array":ARRAY_TYPE,
	"dict":DICTIONARY_TYPE,
	"int":INT_TYPE,
	"int8":INT8_TYPE,
	"int16":INT16_TYPE,
	"int32":INT32_TYPE,
	"int64":INT64_TYPE,
	"float":FLOAT_TYPE,
	"float8":FLOAT8_TYPE,
	"float16":FLOAT16_TYPE,
	"float32":FLOAT32_TYPE,
	"float64":FLOAT64_TYPE,
	"uint":UINT_TYPE,
	"uint8":UINT8_TYPE,
	"uint16":UINT16_TYPE,
	"uint32":UINT32_TYPE,
	"uint64":UINT64_TYPE,
	"ufloat":UFLOAT_TYPE,
	"ufloat8":UFLOAT8_TYPE,
	"ufloat16":UFLOAT16_TYPE,
	"ufloat32":UFLOAT32_TYPE,
	"ufloat64":UFLOAT64_TYPE,
	"string":STRING_TYPE,
}

var keywords = map[string]TokenType{
    "and":AND,
    "break":BREAK,
    "class":CLASS,
    "const":CONST,
    "default":DEFAULT,
    "else":ELSE,
    "false":FALSE,
    "func":FUNC,
    "for":FOR,
    "get":GET,
    "goto":GOTO,
    "if":IF,
    "not":NOT,
    "null":NULL,
    "or":OR,
    "optional":OPTIONAL,
    "package":PACKAGE,
    "print":PRINT,
    "pub":PUBLIC,
    "prv":PRIVATE,
    "return":RETURN,
    "set":SET,
    "super":SUPER,
	"struct":STRUCT,
    "switch":SWITCH,
    "this":THIS,
    "true":TRUE,
    "typedef":TYPEDEF,
    "var":VAR,
    "while":WHILE,
    "whith":WITH,
}

func (s *Scanner) ScanTokens() []Token {
    for !s.isAtEnd() {
        s.start = s.current;
        s.scanToken();
    }
    s.tokens = append(s.tokens, Token{
        tokenType: EOF,
        Lexeme: "",
        Literal: nil,
        Line: s.line,
    })
    return s.tokens;
}

func (s *Scanner) isAtEnd() bool {
    return s.current >= len(s.source);
}

func (s *Scanner) scanToken() {
    var c byte = s.advance();
    
    switch c {
    	case '(': s.addToken(LEFT_PAREN); break;
    	case ')': s.addToken(RIGHT_PAREN); break;
    	case '{': s.addToken(LEFT_BRACE); break;
    	case '}': s.addToken(RIGHT_BRACE); break;
    	case ',': s.addToken(COMMA); break;
    	case '.': s.addToken(DOT); break;
    	case '-':
        	tokenType := MINUS;
            if s.match('>') {
                tokenType = RIGHT_ARROW
            } else if s.match('=') {
                tokenType = MINUS_EQUAL
            }
        	s.addToken(tokenType); break;
    	case '+':
        	tokenType := PLUS
        	if s.match('=') {
                tokenType = PLUS_EQUAL
            }
        	s.addToken(tokenType); break;
    	case ';': s.addToken(SEMICOLON); break;
    	case '*':
        	tokenType := STAR
        	if s.match('=') {
                tokenType = STAR_EQUAL
            }
            s.addToken(tokenType); break;
    	case '!':
        	tokenType := BANG;
            if s.match('=') {
                tokenType = BANG_EQUAL;
            }
            s.addToken(tokenType); break;
    	case '=':
        	tokenType := EQUAL;
            if s.match('=') {
                tokenType = EQUAL_EQUAL;
            }
        	s.addToken(tokenType); break;
		case '<':
        	tokenType := LESS
            if s.match('=') {
                tokenType = LESS_EQUAL;
            } else if s.match('<') {
                tokenType = CONCAT;
            }
        	s.addToken(tokenType); break;
        case '>':
        	tokenType := GREATER
            if s.match('=') {
                tokenType = GREATER_EQUAL;
            }
        	s.addToken(tokenType); break;
        case '/':
        	if s.match('/') {
            	for s.peek() != '\n' && !s.isAtEnd() {
                    s.advance()
                }
            } else if s.match('*') {
                for !s.isAtEnd() {
                    if s.peek() == '*' && s.peekNext() == '/' {
                        s.advance()
                        s.advance()
                        break
                    }
                    s.advance()
                }    
            } else {
                tokenType := SLASH
                if s.match('=') {tokenType = SLASH_EQUAL}
                s.addToken(tokenType); break
            }
        break
        case ' ':
        case '\r':
        case '\t':
        	break;
        case '\n':
        	s.line++;
            break;
        case '"': s.stringify(); break;
        
        default:
        	if s.isDigit(c) {
                s.number()
            } else if s.isAlpha(c) {
                s.identifier() 
            } else{
            for {
            	s.advance()
        		if s.isAtEnd() {
                	PrintError(3, "May be due to a symbol that got in a contexte where it shouldn't be")
                	return
                }
            }
        }
    }
}

func (s *Scanner) identifier() {
    for s.isAlphaNumeric(s.peek()) {
        s.advance()
    }
    var text string
    text = s.source[s.start:s.current]
    tokenType, isKeyword := keywords[text]
    if !isKeyword {
        _tokenType, isType := dataType[text]
        if !isType {_tokenType = IDENTIFIER}
        tokenType= _tokenType
    }
    
    s.addToken(tokenType)
}

func (s *Scanner) number() {
    for s.isDigit(s.peek()) {
        s.advance()
    }
    var isFloating bool = false
    if s.peek() == '.' && s.isDigit(s.peekNext()) {
        isFloating = true
        s.advance()
        for s.isDigit(s.peek()) {
            s.advance()
        }
    }
    var val any
    var err error
    var tokenType TokenType
    if !isFloating {
        val, err = strconv.ParseInt(s.source[s.start:s.current], 10, 32)
        tokenType = INT
    } else if isFloating {
        val, err = strconv.ParseFloat(s.source[s.start:s.current], 32)
        tokenType = FLOAT
    }
    if err != nil {
        panic(err)
    }
    s.addSpecToken(tokenType, val)
}

func (s *Scanner) peekNext() byte {
    if s.current + 1 >= len(s.source) {
        return '\x00'
    }
    return s.source[s.current + 1]
}

func (s *Scanner) isAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
    return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) advance() byte {
    char := s.source[s.current];
    s.current++;
    return char;
}

func (s *Scanner) addToken(tokenType TokenType) {
    s.addSpecToken(tokenType, nil)
}

func (s *Scanner) addSpecToken(tokType TokenType, literal any) {
    var text string = s.source[s.start:s.current];
    input = text
    s.tokens = append(s.tokens, Token{
        tokenType: tokType,
        Lexeme: text,
        Literal: literal,
        Line: s.line})
}
func (s *Scanner) match(expected byte) bool {
    if s.isAtEnd() {
        return false;
    }
    if s.source[s.current] != expected {
        return false;
    }
    
    s.current++
    return true;
}

func (s *Scanner) peek() byte {
    if s.isAtEnd() {
    	return 0;
    }
    return s.source[s.current];
}

func (s *Scanner) isDigit(expected byte) bool {
    return expected >= '0' && expected <= '9'
}

func (s *Scanner) stringify() {
    for s.peek() != '"' && !s.isAtEnd() {
        if s.peek() == '\n' {
            s.line++;
        }
        s.advance();
    }
    if s.isAtEnd() {
        PrintError(8, "Missing \" on a string and reached EOF");
    }
    s.advance();
    
    var value string
    value = s.source[s.start + 1:s.current - 1]
    s.addSpecToken(STRING, value)
}