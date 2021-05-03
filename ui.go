package main

import (
	"fmt"
	"log"
	"strings"

	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

func init() {
	encoding.Register()
}

type History struct {
	Prog  *Program
	Value float64
}

type UI struct {
	parser *participle.Parser
	screen tcell.Screen

	ctx  *Context
	cmds map[string]cmdFunc

	input        *uiInput
	currentProg  *InteractiveProgram
	currentValue float64
	currentError error

	idx     int
	bases   []int
	history []History
}

func NewUI() *UI {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Init(); err != nil {
		log.Fatal(err)
	}

	s.SetStyle(styleDefault)

	u := &UI{
		screen: s,
		input:  &uiInput{},
		parser: participle.MustBuild(
			&InteractiveProgram{},
			participle.Lexer(lex),
			participle.Elide("Whitespace"),
			participle.UseLookahead(30),
		),
		ctx:   NewContext(),
		cmds:  map[string]cmdFunc{},
		bases: []int{10, 16, 2},
	}
	u.cmds["base"] = u.cmdSetBase
	u.cmds["dbases"] = u.cmdSetDisplayedBases

	return u
}

func (u *UI) Run() {
	u.draw()

	for {
		ev := u.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				u.input.Add(ev.Rune())
				u.evalInput()
				u.draw()
			case tcell.KeyDEL:
				u.input.Del()
				u.evalInput()
				u.draw()
			case tcell.KeyDelete:
				u.input.DelFwd()
				u.evalInput()
				u.draw()
			case tcell.KeyEnter:
				u.exec()
			case tcell.KeyLeft:
				u.input.MoveCursor(-1)
				u.draw()
			case tcell.KeyRight:
				u.input.MoveCursor(1)
				u.draw()
			case tcell.KeyEscape:
				return
			}
		case *tcell.EventResize:
			u.draw()
		}
	}
}

func (u *UI) Stop() {
	u.screen.Fini()
}

func (u *UI) exec() {
	if u.currentProg == nil {
		return
	}

	if u.currentProg.Command != nil {
		u.evalCmd()
		u.input.Clear()
		u.currentProg = nil
		u.currentValue = 0
	} else {
		u.history = append(u.history, History{
			Prog:  u.currentProg.Program,
			Value: u.currentValue,
		})
		u.ctx.Vars[strconv.Itoa(u.idx)] = u.currentValue
		u.input.SetText(u.formatValue(u.currentValue, u.ctx.Base))
		u.idx++

		u.evalInput()
	}

	u.draw()
}

func (u *UI) draw() {
	u.screen.Clear()

	_, h := u.screen.Size()
	y := h - 1
	y = u.drawBottomBar(y)
	y = u.drawPreview(y)
	y = u.drawInput(y)
	y = u.drawHistory(y)

	u.screen.Show()
}

func (u *UI) evalInput() {
	in := u.input.Text()
	if strings.TrimSpace(in) == "" {
		u.currentError = nil
		u.currentProg = nil
		u.currentValue = 0
		return
	}

	prog, v, err := u.evalExpr(in)
	if err != nil {
		u.currentError = err
		u.currentProg = nil
		u.currentValue = 0
	} else {
		u.currentProg = prog
		u.currentValue = v
		u.currentError = nil
	}
}

func (u *UI) evalExpr(ex string) (*InteractiveProgram, float64, error) {
	e := &InteractiveProgram{}
	err := u.parser.ParseString("", ex, e)
	if err != nil {
		return nil, 0, err
	}

	if e.Program != nil {
		v, err := e.Program.Eval(u.ctx)
		if err != nil {
			return nil, 0, err
		}

		return e, v, nil
	}

	return e, 0, nil
}

func (u *UI) evalCmd() {
	fn, ok := u.cmds[u.currentProg.Command.Name]
	if !ok {
		u.currentError = fmt.Errorf("command %s does not exist", u.currentProg.Command.Name)
		return
	}

	err := fn(u.currentProg.Command.Inputs...)
	if err != nil {
		u.currentError = err
		return
	}
}

func (u *UI) emitStr(x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		u.screen.SetContent(x, y, c, comb, style)
		x += w
	}
}
