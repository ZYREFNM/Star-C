package main

type ErrorHandler struct {
    Name string
    FilePath string
    Errors []error
    ErrorLine []int
    CurrentParser *Parser
}

type StarError struct {
    ErrorLine int
    Name string
}