package main

import (
	"log"

	"github.com/alecthomas/participle/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

func init() {
	encoding.Register()
}

type History struct {
	Expr  string
	Value float64
	Error error
}

type UI struct {
	parser *participle.Parser
	ctx    *InteractiveContext
	screen tcell.Screen

	input   *uiInput
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

	return &UI{
		screen: s,
		input:  &uiInput{},
		parser: participle.MustBuild(
			&InteractiveCommand{},
			participle.Lexer(lex),
			participle.Elide("Whitespace"),
			participle.UseLookahead(30),
		),
		ctx: NewInteractiveContext(),
	}
}

func (u *UI) Run() error {
	u.draw()

	for {
		ev := u.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				u.input.Add(ev.Rune())
				u.draw()
			case tcell.KeyDEL:
				u.input.Del()
				u.draw()
			case tcell.KeyDelete:
				u.input.DelFwd()
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
				return nil
			}
		case *tcell.EventResize:
			u.draw()
		}
	}

}

func (u *UI) exec() {
	ex := u.input.Text()

	cmd := &InteractiveCommand{}
	err := u.parser.ParseString("", ex, cmd)

	h := History{
		Expr:  ex,
		Error: err,
	}
	if err == nil {
		v, err := cmd.Exec(u.ctx)
		h.Error = err
		h.Value = v
	}

	u.history = append(u.history, h)
	u.input.Clear()

	u.draw()
}

func (u *UI) draw() {
	u.screen.Clear()
	u.drawInput()
	u.drawHistory()
	u.drawBottomBar()
	u.screen.Sync()
}

func (u *UI) evalExpr(ex string) (float64, bool) {
	e := &InteractiveCommand{}
	err := u.parser.ParseString("", ex, e)
	if err != nil {
		return 0, false
	}

	v, err := e.Exec(u.ctx)
	if err != nil {
		return 0, false
	}

	return v, true
}

func (u *UI) inputBarSize() int {
	return len(u.ctx.DisplayedBases) + 1
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
