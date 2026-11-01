package main

type ErrorHandler struct {
    Name string
    FilePath string
    Errors []error
    ErrorLine []int
    CurrentParser *Parser
}

func (e *ErrorHandler) LinkParser(p *Parser) {
    e.CurrentParser = p
}
