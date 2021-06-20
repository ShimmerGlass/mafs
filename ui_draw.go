package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mitchellh/go-wordwrap"
)

func (u *UI) drawBlock(x1, y1, x2, y2 int, style tcell.Style) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			u.screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

func (u *UI) drawPreview(y int) int {
	w, _ := u.screen.Size()
	bases := u.displayedBases()

	if u.currentError != nil {
		wrapped := strings.Split(wordwrap.WrapString(u.currentError.Error(), uint(w-2)), "\n")
		lines := len(wrapped)
		if lines < len(bases) {
			lines = len(bases)
		}
		for i := lines - 1; i >= 0; i-- {
			if i < len(wrapped) {
				u.emitStr(1, y, styleExpressionError, wrapped[i])
			}
			y--
		}
	} else {
		for i := len(bases) - 1; i >= 0; i-- {
			b := bases[i]
			u.drawBlock(0, y, w, y, stylePreview)
			u.emitStr(1, y, stylePreview, fmt.Sprintf("%2d> ", b))

			if u.currentValue != nil {
				u.printBase(w-1, y, u.currentValue, b, u.displayedBases(), stylePreview)
			}
			y--
		}
	}
	return y
}

func (u *UI) displayedBases() []int {
	if u.ctx.Type == typeFloat {
		return []int{10, 2}
	}
	return dedupInt(append([]int{u.ctx.Base}, u.bases...))
}

func (u *UI) drawInput(y int) int {
	ex := u.input.Runes()
	w, _ := u.screen.Size()

	prefix := fmt.Sprintf("$%d = ", u.idx)

	u.drawBlock(0, y, w, y, styleInput)
	u.emitStr(1, y, styleInput.Foreground(colorComment), prefix)

	for i := 0; i <= len(ex); i++ {
		style := styleInput
		if i == u.input.Cursor {
			style = styleInputCursor
		}
		r := ' '
		if i < len(ex) {
			r = ex[i]
		}

		u.screen.SetContent(len(prefix)+i+1, y, r, nil, style)
	}
	y--
	return y
}

func (u *UI) drawHistory(y int) int {
	w, _ := u.screen.Size()
	for i := len(u.history) - 1; i >= 0; i-- {
		if y < 0 {
			break
		}

		h := u.history[i]
		prefix := fmt.Sprintf("$%d = ", i)
		u.emitStr(1, y, styleExpression.Foreground(colorComment), prefix)
		u.emitStr(1+len(prefix), y, styleExpression, h.Prog.String())
		u.printBase(w-1, y, h.Value, u.ctx.Base, []int{}, styleExpression)

		y--
	}

	return y
}

func (u *UI) drawBottomBar(y int) int {
	w, _ := u.screen.Size()

	u.drawBlock(0, y, w-1, y, styleBar)
	u.emitStr(1, y, styleBar, fmt.Sprintf("Base: %d | Displayed bases: %v | Type: %s", u.ctx.Base, u.displayedBases(), u.ctx.Type))
	y--
	return y
}
