// nolint: govet
package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2/lexer/stateful"
)

var lex = stateful.MustSimple([]stateful.Rule{
	{"Var", `\$\w+`, nil},
	{"Ref", `#\d+`, nil},
	{"Command", `/[a-z]+`, nil},
	{"Call", `\w+\(`, nil},
	{"Number", `([\da-fA-F]+(?:\.\d+)?)|\.\d+`, nil},
	{"Int", `\d+`, nil},
	{"Punct", `[*/+\-(),=?%]`, nil},
	{"Whitespace", `[ \t\n\r]+`, nil},
})

type Operator string

const (
	OpMul Operator = "*"
	OpDiv Operator = "/"
	OpMod Operator = "%"
	OpAdd Operator = "+"
	OpSub Operator = "-"
)

type UnaryOperator string

const (
	OpPlus  UnaryOperator = "+"
	OpMinus UnaryOperator = "-"
)

type Value struct {
	Operator *UnaryOperator `@("+" | "-")?`
	Number   *string        `(   @Number`
	Call     *Call          `  | @@`
	Variable *string        `  | (@Var | @Ref)`

	Subexpression *Expression `| "(" @@ ")" )`
}

type Call struct {
	Function string        `@Call`
	Inputs   []*Expression `( @@ ( "," @@ )* )? ")"`
}

type OpFactor struct {
	Operator Operator `@("*" | "/" | "%")`
	Factor   *Value   `@@`
}

type Term struct {
	Left  *Value      `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type Assignment struct {
	Name  string      `@Var "="`
	Value *Expression `@@`
}

type Command struct {
	Assignment *Assignment `  @@ `
	Expr       *Expression `| @@`
}

type InteractiveCommand struct {
	CtxCommand *string  `  @(Command | "?")`
	Command    *Command `| @@`
}

// Display

func (o Operator) String() string {
	return string(o)
}

func (v *Value) String() string {
	if v.Number != nil {
		return *v.Number
	}
	if v.Variable != nil {
		return *v.Variable
	}
	return "(" + v.Subexpression.String() + ")"
}

func (o *OpFactor) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Factor)
}

func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (e *Expression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// Evaluation

func (o Operator) Eval(l, r float64) (float64, error) {
	switch o {
	case OpMul:
		return l * r, nil
	case OpDiv:
		return l / r, nil
	case OpMod:
		return math.Mod(l, r), nil
	case OpAdd:
		return l + r, nil
	case OpSub:
		return l - r, nil
	}
	panic("unsupported operator")
}

func (o UnaryOperator) Eval(v float64) (float64, error) {
	switch o {
	case OpMinus:
		return -v, nil
	case OpPlus:
		return v, nil
	}
	panic("unsupported operator")
}

func (v *Value) Eval(ctx *Context) (float64, error) {
	var value float64
	switch {
	case v.Number != nil:
		val, err := ParseNumber(ctx.Base, *v.Number)
		if err != nil {
			return 0, err
		}
		value = val
	case v.Variable != nil:
		val, ok := ctx.Vars[*v.Variable]
		if !ok {
			return 0, fmt.Errorf("no such variable " + *v.Variable)
		}
		value = val.V
	case v.Call != nil:
		val, err := v.Call.Eval(ctx)
		if err != nil {
			return 0, err
		}
		value = val
	case v.Subexpression != nil:
		val, err := v.Subexpression.Eval(ctx)
		if err != nil {
			return 0, err
		}
		value = val
	}

	if v.Operator != nil {
		return v.Operator.Eval(value)
	}

	return value, nil
}

func (t *Term) Eval(ctx *Context) (float64, error) {
	n, err := t.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}
	for _, r := range t.Right {
		fact, err := r.Factor.Eval(ctx)
		if err != nil {
			return 0, err
		}
		n, err = r.Operator.Eval(n, fact)
		if err != nil {
			return 0, err
		}
	}
	return n, nil
}

func (e *Expression) Eval(ctx *Context) (float64, error) {
	l, err := e.Left.Eval(ctx)
	if err != nil {
		return 0, err
	}
	for _, r := range e.Right {
		term, err := r.Term.Eval(ctx)
		if err != nil {
			return 0, err
		}
		l, err = r.Operator.Eval(l, term)
		if err != nil {
			return 0, err
		}
	}
	return l, nil
}

func (a *Assignment) Eval(ctx *Context) (float64, error) {
	v, err := a.Value.Eval(ctx)
	if err != nil {
		return 0, err
	}

	ctx.Vars[a.Name] = Var{v, ""}
	return v, nil
}

func (c *Call) Eval(ctx *Context) (float64, error) {
	name := c.Function[:len(c.Function)-1]
	f, ok := ctx.Funcs[name]
	if !ok {
		return 0, fmt.Errorf("no such function %s", name)
	}

	args := make([]float64, len(c.Inputs))
	for i := range c.Inputs {
		v, err := c.Inputs[i].Eval(ctx)
		if err != nil {
			return 0, err
		}
		args[i] = v
	}

	return f.F(args)
}

func (c *Command) Eval(ctx *Context) (float64, error) {
	switch {
	case c.Expr != nil:
		return c.Expr.Eval(ctx)

	default:
		return c.Assignment.Eval(ctx)
	}
}

func (c *InteractiveCommand) Exec(ctx *InteractiveContext) {
	switch {
	case c.CtxCommand != nil:
		cmd, ok := ctx.Commands[*c.CtxCommand]
		if !ok {
			ctx.PrintError(ctx, fmt.Sprintf("no such command %s", *c.CtxCommand))
			return
		}
		cmd(ctx)

	default:
		v, err := c.Command.Eval(ctx.Context)
		if err != nil {
			ctx.PrintError(ctx, err.Error())
			return
		}

		ctx.PrintValue(ctx, v)
		ctx.Vars["#"+strconv.Itoa(ctx.Idx)] = Var{V: v, Help: ""}
		ctx.Idx++
		ctx.LastResult = v
		ctx.HasLastResult = true
	}
}
