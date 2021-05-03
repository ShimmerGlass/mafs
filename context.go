package main

import (
	"math"
)

type Func func([]float64) (float64, error)

type Context struct {
	Base  int
	Vars  map[string]float64
	Funcs map[string]Func
}

func NewContext() *Context {
	return &Context{
		Base: 10,
		Vars: map[string]float64{
			"$pi":      math.Pi,
			"$max_u8":  float64(math.MaxUint8),
			"$max_u16": float64(math.MaxUint16),
			"$max_u32": float64(math.MaxUint32),
			"$max_8":   float64(math.MaxInt8),
			"$max_16":  float64(math.MaxInt16),
			"$max_32":  float64(math.MaxInt32),
		},
		Funcs: map[string]Func{
			"sqrt": sqrt,
			"pow":  pow,
		},
	}
}
