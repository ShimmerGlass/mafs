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

	ui, err := NewUI()
	if err != nil {
		log.Fatal(err)
	}

	err = ui.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func evalOne(in string) (float64, error) {
	expr := &Command{}
	err := participle.MustBuild(&Command{}, participle.Lexer(lex)).ParseString("", in, expr)
	if err != nil {
		return 0, err
	}

	return expr.Eval(NewContext())
}
