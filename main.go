package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/alecthomas/participle/v2"
	"github.com/shimmerglass/mafs/num"
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

	var r interface{}
	done := make(chan int, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- 0
	}()

	go func() {
		defer func() {
			if r = recover(); r != nil {
				r = fmt.Sprintf("%s\n%s", r, debug.Stack())
				done <- 1
			}
		}()
		ui.Run()
		done <- 0
	}()

	code := <-done
	ui.Stop()
	if r != nil {
		fmt.Fprint(os.Stderr, r)
	}
	os.Exit(code)
}

func evalOne(in string) (num.Number, error) {
	expr := &Program{}
	err := participle.MustBuild(
		&Program{},
		participle.Lexer(lex),
		participle.UseLookahead(30),
	).ParseString("", in, expr)
	if err != nil {
		return nil, err
	}

	return expr.Eval(NewContext())
}
