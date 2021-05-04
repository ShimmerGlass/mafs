package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

var baseColors = map[int]tcell.Color{
	2:  colorCyan,
	-2: colorComment,
	10: colorYellow,
	16: colorBlue,
}

func (u *UI) formatValue(v float64, base int) string {
	if base == 10 {
		return strconv.FormatFloat(v, 'f', -1, 64)
	}

	return strconv.FormatInt(int64(v), base)
}

func (u *UI) printBase(x, y int, v float64, base int, style tcell.Style) {
	f := []rune(strconv.FormatInt(int64(v), base))
	l := 35

	switch base {
	case 2:
		i := len(f) - 1
		for dx := l; dx >= 0; dx-- {
			c := '0'
			if i >= 0 {
				c = f[i]
			}

			color := baseColors[-2]
			if c == '1' {
				color = baseColors[2]
			}

			u.screen.SetContent(x+dx, y, c, nil, style.Foreground(color))
			if (len(f)-i)%8 == 0 {
				dx--
			}
			i--
		}

	case 16:
		i := len(f) - 1
		for dx := l; dx >= 0; dx-- {
			c := ' '
			if i >= 0 {
				c = f[i]
			}

			u.screen.SetContent(x+dx, y, c, nil, style.Foreground(baseColors[16]))
			switch (len(f) - i) % 2 {
			case 1:
				dx -= 3
			case 0:
				dx -= 4
			}
			i--
		}

	default:
		if c, ok := baseColors[base]; ok {
			style = style.Foreground(c)
		}
		u.emitStr(x+l-len(f)+1, y, style, string(f))
	}
}
