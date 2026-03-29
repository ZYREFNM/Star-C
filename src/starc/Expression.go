package main

type Visitor interface {
    Accept(visitor Visitor)
}