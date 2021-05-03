package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
)

func main() {
	if len(os.Args) > 1 {
		v, err := evalOne(strings.Join(os.Args[1:], " "))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(v)
		return
	}

	ui := NewUI()

	err := ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func evalOne(in string) (float64, error) {
	expr := &Command{}
	err := participle.MustBuild(
		&Command{},
		participle.Lexer(lex),
		participle.UseLookahead(30),
	).ParseString("", in, expr)
	if err != nil {
		return 0, err
	}

	return expr.Eval(NewContext())
}
