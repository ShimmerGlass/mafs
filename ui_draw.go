package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func (u *UI) drawBlock(x1, y1, x2, y2 int, style tcell.Style) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			u.screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

func (u *UI) drawInput() {
	ex := u.input.Runes()

	w, h := u.screen.Size()
	y := h - 2

	v, ok := u.evalExpr(u.input.Text())

	for i := len(u.ctx.DisplayedBases) - 1; i >= 0; i-- {
		b := u.ctx.DisplayedBases[i]
		u.drawBlock(0, y, w, y, stylePreview)
		if ok {
			u.emitStr(1, y, stylePreview, fmt.Sprintf("%2d> ", b))
			u.printBase(w-36-1, y, v, b, stylePreview)
		}
		y--
	}

	u.drawBlock(0, y, w, y, styleInput)
	for i := 0; i <= len(ex); i++ {
		style := styleInput
		if i == u.input.Cursor {
			style = styleInputCursor
		}
		r := ' '
		if i < len(ex) {
			r = ex[i]
		}

		u.screen.SetContent(i+1, y, r, nil, style)
	}
}

func (u *UI) drawHistory() {
	w, h := u.screen.Size()
	y := h - 2 - u.inputBarSize()

	for i := len(u.history) - 1; i >= 0; i-- {
		if y < 0 {
			break
		}

		h := u.history[i]
		u.emitStr(1, y, styleExpression, h.Expr)
		if h.Error != nil {
			e := h.Error.Error()
			u.emitStr(w-len(e)-1, y, styleExpressionError, e)
		} else {
			u.printBase(w-36-1, y, h.Value, u.ctx.Base, styleExpression)
		}

		y--
	}
}

func (u *UI) drawBottomBar() {
	w, h := u.screen.Size()
	y := h - 1
	for x := 0; x < w; x++ {
		u.screen.SetContent(x, y, ' ', nil, styleBar)
	}

	u.emitStr(1, y, styleBar, fmt.Sprintf("Base: %d", u.ctx.Base))
}
