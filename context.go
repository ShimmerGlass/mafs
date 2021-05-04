package main

import (
	"math"

	"github.com/shimmerglass/mafs/num"
)

const (
	typeSignedInt = "sint"
	typeFloat     = "float"
)

type Func func([]num.Number) (num.Number, error)

type Context struct {
	Base  int
	Type  string
	Vars  map[string]num.Number
	Funcs map[string]Func
}

func NewContext() *Context {
	return &Context{
		Base: 10,
		Vars: map[string]num.Number{
			"pi": num.Float(math.Pi),
		},
		Funcs: map[string]Func{
			// "sqrt": sqrt,
			// "pow":  pow,
		},
		Type: typeSignedInt,
	}
}
