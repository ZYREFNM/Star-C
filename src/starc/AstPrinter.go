package main

import (
    "fmt"
    "strings"
)

type AstPrinter struct {
    
}

func (a *AstPrinter) Print(expr Expr) string {
    expression := expr.Accept(a)
    return fmt.Sprintf("%v", expression)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) string {
    return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) string {
    return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) string {
    if expr.Value == nil {
        return "null"
    }
    return fmt.Sprintf("%v", expr)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) string {
    return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
    var builder strings.Builder
    builder.WriteString("(")
    builder.WriteString(name)
    
    for _, expr := range exprs {
        builder.WriteString(" ")
        res := expr.Accept(&AstPrinter{})
        builder.WriteString(fmt.Sprintf("%v", res))
    }
    builder.WriteString(")")
    
    return builder.String()
}