// nolint: govet
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/shimmerglass/mafs/num"
)

var lex = stateful.MustSimple([]stateful.Rule{
	{
		Name:    "Word",
		Pattern: `\w+`,
	},
	{
		Name:    "Operator",
		Pattern: `[*/%+\-]`,
	},
	{
		Name:    "Punct",
		Pattern: `[()\.,#$=]`,
	},
	{
		Name:    "Whitespace",
		Pattern: `\s+`,
	},
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
	Call     *Call          `(   @@`
	Variable *Variable      `  | @@`
	Number   *Number        `  | @@`

	Subexpression *Expression `| ( "(" @@ ")" ) )`
}

type Number struct {
	Left  *string `@Word?`
	Right *string `( "." @Word )?`
}

type Variable struct {
	Name string `"$" @Word`
}

type Call struct {
	Function string        `@Word`
	Inputs   []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
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
	Variable *Variable   `@@ "="`
	Value    *Expression `@@`
}

type Program struct {
	Assignment *Assignment `  @@ `
	Expr       *Expression `| @@`
}

type Command struct {
	Name   string   `"#" @Word`
	Inputs []string `( @Word ( "," @Word )* )?`
}

type InteractiveProgram struct {
	Command *Command `  @@`
	Program *Program `| @@`
}

// Display

func (o Operator) String() string {
	return string(o)
}

func (v *Value) String() string {
	res := ""
	if v.Operator != nil {
		res += string(*v.Operator)
	}
	switch {
	case v.Number != nil:
		res += v.Number.String()

	case v.Variable != nil:
		res += v.Variable.String()

	case v.Call != nil:
		res += v.Call.String()

	case v.Subexpression != nil:
		res += "(" + v.Subexpression.String() + ")"
	}

	return res
}

func (n *Number) String() string {
	r := ""
	if n.Left != nil {
		r += *n.Left
	} else {
		r += "0"
	}
	if n.Right != nil {
		r += "." + *n.Right
	}
	return r
}

func (v *Variable) String() string {
	return "$" + v.Name
}

func (c *Call) String() string {
	args := []string{}
	for _, a := range c.Inputs {
		args = append(args, a.String())
	}

	return fmt.Sprintf("%s(%s)", c.Function, strings.Join(args, ", "))
}

func (o *OpFactor) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Factor)
}

func (t *Term) String() string {
	out := []string{}
	if t.Left != nil {
		out = append(out, t.Left.String())
	}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (e *Expression) String() string {
	out := []string{}
	if e.Left != nil {
		out = append(out, e.Left.String())
	}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (p *Program) String() string {
	switch {
	case p.Assignment != nil:
		return fmt.Sprintf("%s = %s", p.Assignment.Variable.String(), p.Assignment.Value.String())

	case p.Expr != nil:
		return p.Expr.String()
	}

	return ""
}

// Evaluation

func (n *Number) Eval(ctx *Context) (num.Number, error) {
	switch ctx.Type {
	case typeFloat:
		s := *n.Left
		if n.Right != nil {
			s += "." + *n.Right
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return num.Float(v), nil

	case typeSignedInt:
		if n.Right != nil {
			return nil, fmt.Errorf("decimal point not allowed in int mode")
		}
		v, err := strconv.ParseInt(*n.Left, ctx.Base, 64)
		if err != nil {
			return nil, err
		}
		return num.SignedInt(v), nil

	default:
		return nil, fmt.Errorf("bad type %s", ctx.Type)
	}
}

func (o Operator) Eval(l, r num.Number) (num.Number, error) {
	if !reflect.TypeOf(l).AssignableTo(reflect.TypeOf(r)) {
		return nil, fmt.Errorf("incompatible types")
	}

	switch o {
	case OpMul:
		return l.Mul(r)
	case OpDiv:
		return l.Div(r)
	case OpMod:
		return l.Mod(r)
	case OpAdd:
		return l.Add(r)
	case OpSub:
		return l.Sub(r)
	}
	panic("unsupported operator")
}

func (o UnaryOperator) Eval(v num.Number) (num.Number, error) {
	switch o {
	case OpMinus:
		return v.Inverse()
	case OpPlus:
		return v, nil
	}
	panic("unsupported operator")
}

func (v *Value) Eval(ctx *Context) (num.Number, error) {
	var value num.Number
	switch {
	case v.Number != nil:
		val, err := v.Number.Eval(ctx)
		if err != nil {
			return nil, err
		}
		value = val
	case v.Variable != nil:
		val, ok := ctx.Vars[v.Variable.Name]
		if !ok {
			return nil, fmt.Errorf("no such variable " + v.Variable.Name)
		}
		value = val
	case v.Call != nil:
		val, err := v.Call.Eval(ctx)
		if err != nil {
			return nil, err
		}
		value = val
	case v.Subexpression != nil:
		val, err := v.Subexpression.Eval(ctx)
		if err != nil {
			return nil, err
		}
		value = val
	}

	if v.Operator != nil {
		return v.Operator.Eval(value)
	}

	return value, nil
}

func (t *Term) Eval(ctx *Context) (num.Number, error) {
	n, err := t.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range t.Right {
		fact, err := r.Factor.Eval(ctx)
		if err != nil {
			return nil, err
		}
		n, err = r.Operator.Eval(n, fact)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}

func (e *Expression) Eval(ctx *Context) (num.Number, error) {
	l, err := e.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range e.Right {
		term, err := r.Term.Eval(ctx)
		if err != nil {
			return nil, err
		}
		l, err = r.Operator.Eval(l, term)
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}

func (a *Assignment) Eval(ctx *Context) (num.Number, error) {
	v, err := a.Value.Eval(ctx)
	if err != nil {
		return nil, err
	}

	ctx.Vars[a.Variable.Name] = v
	return v, nil
}

func (c *Call) Eval(ctx *Context) (num.Number, error) {
	f, ok := ctx.Funcs[c.Function]
	if !ok {
		return nil, fmt.Errorf("no such function %s", c.Function)
	}

	args := make([]num.Number, len(c.Inputs))
	for i := range c.Inputs {
		v, err := c.Inputs[i].Eval(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = v
	}

	return f(args)
}

func (c *Program) Eval(ctx *Context) (num.Number, error) {
	switch {
	case c.Expr != nil:
		return c.Expr.Eval(ctx)

	default:
		return c.Assignment.Eval(ctx)
	}
}
