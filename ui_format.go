package main

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/gdamore/tcell/v2"
	"github.com/shimmerglass/mafs/num"
)

var baseColors = map[int]tcell.Color{
	2:  colorCyan,
	-2: colorComment,
	10: colorYellow,
	16: colorBlue,
}

func (u *UI) formatValue(v num.Number, base int) string {
	switch n := v.(type) {
	case num.Float:
		return strconv.FormatFloat(float64(n), 'f', -1, 64)
	case num.SignedInt:
		return strconv.FormatInt(int64(n), base)
	}

	panic(fmt.Errorf("unknow type %T", v))
}

func (u *UI) printBase(x, y int, v num.Number, base int, adjacent []int, style tcell.Style) {
	f := []rune(u.formatValue(v, base))

	switch base {
	case 2:
		u.printBase2(x, y, v, style)

	case 16:
		hasBase2 := false
		for _, b := range adjacent {
			if b == 2 {
				hasBase2 = true
				break
			}
		}

		for i := len(f) - 1; i >= 0; i-- {
			c := f[i]

			dx := len(f) - i
			u.screen.SetContent(x-dx, y, c, nil, style.Foreground(baseColors[16]))
			if hasBase2 {
				switch (len(f) - i) % 2 {
				case 1:
					x -= 3
				case 0:
					x -= 4
				}
			}
		}

	default:
		if c, ok := baseColors[base]; ok {
			style = style.Foreground(c)
		}
		u.emitStr(x-len(f), y, style, string(f))
	}
}

func (u *UI) printBase2(x, y int, v num.Number, style tcell.Style) {
	trim := false
	signPos := -1

	var bytes []byte
	switch n := v.(type) {
	case num.Float:
		bytes = (*(*[8]byte)(unsafe.Pointer(&n)))[:]
	case num.SignedInt:
		bytes = (*(*[8]byte)(unsafe.Pointer(&n)))[:]
		trim = true
		signPos = 0
	}

	var f []bool
	for i := len(bytes) - 1; i >= 0; i-- {
		for j := 7; j >= 0; j-- {
			v := (bytes[i]>>j)&1 == 1
			f = append(f, v)
		}
	}

	if trim {
	Next:
		for i := 1; i < len(f)/8; i *= 2 {
			stop := len(f) - i*8
			for j := 0; j < stop; j++ {
				if f[j] {
					continue Next
				}
			}
			f = f[stop:]
			break
		}
	}

	for i := len(f) - 1; i >= 0; i-- {
		dx := len(f) - i
		if f[i] {
			if i == signPos {
				u.screen.SetContent(x-dx, y, '1', nil, style.Foreground(colorRed))
			} else {
				u.screen.SetContent(x-dx, y, '1', nil, style.Foreground(baseColors[2]))
			}
		} else {
			u.screen.SetContent(x-dx, y, '0', nil, style.Foreground(baseColors[-2]))
		}

		if (len(f)-i)%8 == 0 {
			x--
		}
	}
}
