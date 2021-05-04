package main

import "github.com/gdamore/tcell/v2"

var (
	colorBackground  = tcell.NewRGBColor(0x1d, 0x1f, 0x21)
	colorCurrentLine = tcell.NewRGBColor(0x28, 0x2a, 0x2e)
	colorForeground  = tcell.NewRGBColor(0xc5, 0xc8, 0xc6)
	colorComment     = tcell.NewRGBColor(0x96, 0x98, 0x96)
	colorYellow      = tcell.NewRGBColor(0xf0, 0xc6, 0x74)
	colorOrange      = tcell.NewRGBColor(0xde, 0x93, 0x5f)
	colorRed         = tcell.NewRGBColor(0xcc, 0x66, 0x66)
	colorViolet      = tcell.NewRGBColor(0xb2, 0x94, 0xbb)
	colorBlue        = tcell.NewRGBColor(0x81, 0xa2, 0xbe)
	colorCyan        = tcell.NewRGBColor(0x8a, 0xbe, 0xb7)
	colorGreen       = tcell.NewRGBColor(0xb5, 0xbd, 0x68)

	styleDefault = tcell.StyleDefault.Foreground(colorForeground).Background(colorBackground)

	styleInput       = tcell.StyleDefault.Foreground(colorViolet).Background(colorCurrentLine)
	styleInputCursor = tcell.StyleDefault.Foreground(colorCurrentLine).Background(colorViolet)

	stylePreview = tcell.StyleDefault.Foreground(colorComment).Background(colorBackground)
	styleBar     = tcell.StyleDefault.Foreground(colorGreen).Background(colorCurrentLine)

	styleExpression      = tcell.StyleDefault.Foreground(colorViolet).Background(colorBackground)
	styleExpressionError = tcell.StyleDefault.Foreground(colorRed).Background(colorBackground)
)
