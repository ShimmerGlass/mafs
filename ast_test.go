package main

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/stretchr/testify/require"
)

var parser = participle.MustBuild(
	&Command{},
	participle.Lexer(lex),
	participle.Elide("Whitespace"),
)

var evalsBase10 = map[string]float64{
	// values
	"10":     10,
	"10.5":   10.5,
	"-42.10": -42.10,
	".8":     0.8,
	"5e2":    500,

	// operations
	"2 + 5":   7,
	"3+85":    88,
	"32/4":    8,
	"14/2":    7,
	"2 + -6":  -4,
	"2+-6":    -4,
	"3 - +8":  -5,
	"10 % 3":  1,
	".1 * .5": 0.05,

	// precedence
	"2 + 2 * 2": 6,

	// functions
	"sqrt(4)":          2,
	"sqrt(sqrt(16))":   2,
	"-sqrt(4)":         -2,
	"pow(2, 3)":        8,
	"pow(-(2 + 2), 2)": 16,
}

var evalsBase16 = map[string]float64{
	"af":       0xAF,
	"$af = af": 0xAF,
}

func TestEvalBase10(t *testing.T) {
	for in, expected := range evalsBase10 {
		t.Run(in, func(t *testing.T) {
			expr := &Command{}
			err := parser.ParseString("", in, expr)
			require.NoError(t, err)

			v, err := expr.Eval(NewContext())
			require.NoError(t, err)
			require.Equal(t, expected, v)
		})
	}
}

func TestEvalBase16(t *testing.T) {
	for in, expected := range evalsBase16 {
		t.Run(in, func(t *testing.T) {
			expr := &Command{}
			err := parser.ParseString("", in, expr)
			require.NoError(t, err)

			ctx := NewContext()
			ctx.Base = 16
			v, err := expr.Eval(ctx)
			require.NoError(t, err)
			require.Equal(t, expected, v)
		})
	}
}

func TestAssign(t *testing.T) {
	expr := &Command{}

	err := parser.ParseString("", "$foo = 10", expr)
	require.NoError(t, err)

	ctx := NewContext()
	v, err := expr.Eval(ctx)
	require.NoError(t, err)
	require.Equal(t, float64(10), v)

	expr = &Command{}
	err = parser.ParseString("", "$foo", expr)
	require.NoError(t, err)
	v, err = expr.Eval(ctx)
	require.NoError(t, err)
	require.Equal(t, float64(10), v)
}

func TestVarToVarAssignFail(t *testing.T) {
	expr := &Command{}
	err := parser.ParseString("", "$a = b", expr)
	require.NoError(t, err)

	v, err := expr.Eval(NewContext())
	require.Zero(t, v)
	require.Error(t, err)
}
