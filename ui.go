package main

import (
	"fmt"
	"io"
	"math"
	"math/bits"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var colors = map[string]*color.Color{
	"chrome": color.New(color.FgHiBlack),
	"error":  color.New(color.FgHiRed),
	"-2":     color.New(color.FgHiBlack),
	"2":      color.New(color.FgHiBlue),
	"10":     color.New(color.FgHiYellow),
	"16":     color.New(color.FgHiCyan),
}

type UI struct {
	parser *participle.Parser
	ctx    *InteractiveContext
	l      *readline.Instance
}

func NewUI() (*UI, error) {

	ui := &UI{
		parser: participle.MustBuild(
			&InteractiveCommand{},
			participle.Lexer(lex),
			participle.Elide("Whitespace"),
			participle.UseLookahead(30),
		),
		ctx: NewInteractiveContext(
			func(ctx *InteractiveContext, v float64) {
				prompt := fmt.Sprintf("%s» ", strings.Repeat(" ", len(strconv.Itoa(ctx.Idx))+2))
				colors["chrome"].Print(prompt)

				colors["10"].Println(v)

				if v < 0 ||
					math.IsInf(v, 1) ||
					math.IsInf(v, -1) ||
					math.IsNaN(v) ||
					math.Trunc(v) != v {
					return
				}

				vi := uint64(v)

				colors["chrome"].Print(prompt)

				binSize := 32
				if bits.Len64(uint64(vi)) > 32 {
					binSize = 64
				}

				for i := 1; i <= binSize; i++ {
					on := vi&(1<<(uint64(binSize-1)-uint64(i-1))) > 0

					if on {
						colors["2"].Print(1)
					} else {
						colors["-2"].Print(0)
					}
					if i%4 == 0 {
						fmt.Print(" ")
					}
					if i%8 == 0 {
						fmt.Print(" ")
					}
				}

				fmt.Println()

				hex := fmt.Sprintf("%X", vi)
				padding := binSize/4 - len(hex)

				colors["chrome"].Print(prompt)
				fmt.Print(strings.Repeat(" ", padding*5+padding/2))
				for i, c := range hex {
					colors["16"].Print(string(c))
					fmt.Print("    ")
					if (padding+i+1)%2 == 0 {
						fmt.Print(" ")
					}
				}

				fmt.Println()
			},
			func(ctx *InteractiveContext, s string) {
				fmt.Println(s)
			},
			func(ctx *InteractiveContext, s string) {
				colors["error"].Println(s)
			},
		),
	}

	completer := readline.NewPrefixCompleter(
		readline.PcItemDynamic(ui.compl),
	)
	l, err := readline.NewEx(&readline.Config{
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			// block CtrlZ feature
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		return nil, err
	}

	ui.l = l

	return ui, nil
}

func (u *UI) Run() error {
	for {
		u.l.SetPrompt(colors["chrome"].Sprintf("#%d", u.ctx.Idx) + colors[fmt.Sprint(u.ctx.Base)].Sprintf(" » "))
		u.prefillLine()

		line, err := u.l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				u.ctx.HasLastResult = false
				continue
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		u.process(line)
	}

	return nil
}

func (u *UI) process(line string) {
	if strings.TrimSpace(line) == "" {
		return
	}

	expr := &InteractiveCommand{}
	err := u.parser.ParseString("", line, expr)
	if err != nil {
		u.ctx.PrintError(u.ctx, err.Error())
		return
	}

	expr.Exec(u.ctx)
}

func (u *UI) compl(in string) (res []string) {
	for name := range u.ctx.Funcs {
		if strings.HasPrefix(name, in) {
			res = append(res, name+"(")
		}
	}
	for name := range u.ctx.Vars {
		if strings.HasPrefix(name, in) {
			res = append(res, name)
		}
	}
	for name := range u.ctx.Commands {
		if strings.HasPrefix(name, in) {
			res = append(res, name)
		}
	}

	return
}

func (u *UI) prefillLine() {
	if !u.ctx.HasLastResult {
		return
	}
	if u.ctx.Base == 10 {
		u.l.WriteStdin([]byte(fmt.Sprint(u.ctx.LastResult)))
		return
	}
	u.l.WriteStdin([]byte(strconv.FormatInt(int64(u.ctx.LastResult), u.ctx.Base)))
}
