package main

import (
	"fmt"
	"math"
	"strings"
)

type Func struct {
	F    func([]float64) (float64, error)
	Help string
}

type Var struct {
	V    float64
	Help string
}

type Context struct {
	Base  int
	Vars  map[string]Var
	Funcs map[string]Func
}

type InteractiveContext struct {
	*Context

	Idx           int
	HasLastResult bool
	LastResult    float64

	PrintValue func(*InteractiveContext, float64)
	PrintText  func(*InteractiveContext, string)
	PrintError func(*InteractiveContext, string)

	Commands map[string]func(*InteractiveContext)
}

func NewContext() *Context {
	return &Context{
		Base: 10,
		Vars: map[string]Var{
			"$pi":      {math.Pi, "Pi"},
			"$max_u8":  {float64(math.MaxUint8), "Max uint8"},
			"$max_u16": {float64(math.MaxUint16), "Max uint16"},
			"$max_u32": {float64(math.MaxUint32), "Max uint32"},
			"$max_8":   {float64(math.MaxInt8), "Max int8"},
			"$max_16":  {float64(math.MaxInt16), "Max int16"},
			"$max_32":  {float64(math.MaxInt32), "Max int32"},
		},
		Funcs: map[string]Func{
			"sqrt": {sqrt, "Computes the square root of a number."},
			"pow":  {pow, "Takes two arguments (base value and power value) and, returns the power raised to the base number."},
		},
	}
}

func NewInteractiveContext(
	printValue func(*InteractiveContext, float64),
	printText func(*InteractiveContext, string),
	printError func(*InteractiveContext, string),
) *InteractiveContext {
	return &InteractiveContext{
		Context: NewContext(),

		PrintValue: printValue,
		PrintText:  printText,
		PrintError: printError,

		Commands: map[string]func(ctx *InteractiveContext){
			"?":    func(ctx *InteractiveContext) { ctx.PrintHelp() },
			"/dec": func(ctx *InteractiveContext) { ctx.SetBase(10) },
			"/bin": func(ctx *InteractiveContext) { ctx.SetBase(2) },
			"/hex": func(ctx *InteractiveContext) { ctx.SetBase(16) },
		},
	}
}

func (c *InteractiveContext) PrintHelp() {
	b := &strings.Builder{}

	b.WriteString("Functions:\n")
	for name, f := range c.Funcs {
		b.WriteString(fmt.Sprintf("  %s: %s\n", name, f.Help))
	}

	b.WriteString("Variables:\n")
	for name, f := range c.Vars {
		b.WriteString(fmt.Sprintf("  %s: %f\n", name, f.V))
	}

	c.PrintText(c, b.String())
}

func (c *InteractiveContext) SetBase(b int) {
	c.Base = b
	c.PrintText(c, fmt.Sprintf("Base is now %d", b))
}
