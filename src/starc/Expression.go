package main

type Visitor interface {
    VisitBinaryExpr(expr *Binary) string
    VisitGroupingExpr(expr *Grouping) string
    VisitLiteralExpr(expr *Literal) string
    VisitUnaryExpr(expr *Unary) string
}
type Expr interface {
    Accept(visitor Visitor) any
}

type Binary struct {
    Left Expr
    Operator Token
    Right Expr
}
func (b *Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(b)
}
type Grouping struct {
    Expression Expr
}

func (g *Grouping) Accept(visitor Visitor) any {
    return visitor.VisitGroupingExpr(g)
}

type Literal struct {
    Value any
}

func (l *Literal) Accept(visitor Visitor) any {
    return visitor.VisitLiteralExpr(l)
}

type Unary struct {
    Operator Token
    Right Expr
}

func (u *Unary) Accept(visitor Visitor) any {
    return visitor.VisitUnaryExpr(u)
}