package main

import "github.com/gdamore/tcell/v2"

var (
	colorBase03  = tcell.NewRGBColor(0x00, 0x2b, 0x36)
	colorBase02  = tcell.NewRGBColor(0x07, 0x36, 0x42)
	colorBase01  = tcell.NewRGBColor(0x58, 0x6e, 0x75)
	colorBase00  = tcell.NewRGBColor(0x65, 0x7b, 0x83)
	colorBase0   = tcell.NewRGBColor(0x83, 0x94, 0x96)
	colorBase1   = tcell.NewRGBColor(0x93, 0xa1, 0xa1)
	colorBase2   = tcell.NewRGBColor(0xee, 0xe8, 0xd5)
	colorBase3   = tcell.NewRGBColor(0xfd, 0xf6, 0xe3)
	colorYellow  = tcell.NewRGBColor(0xb5, 0x89, 0x00)
	colorOrange  = tcell.NewRGBColor(0xcb, 0x4b, 0x16)
	colorRed     = tcell.NewRGBColor(0xdc, 0x32, 0x2f)
	colorMagenta = tcell.NewRGBColor(0xd3, 0x36, 0x82)
	colorViolet  = tcell.NewRGBColor(0x6c, 0x71, 0xc4)
	colorBlue    = tcell.NewRGBColor(0x26, 0x8b, 0xd2)
	colorCyan    = tcell.NewRGBColor(0x2a, 0xa1, 0x98)
	colorGreen   = tcell.NewRGBColor(0x85, 0x99, 0x00)

	styleDefault = tcell.StyleDefault.Foreground(colorBase0).Background(colorBase03)

	styleInput       = tcell.StyleDefault.Foreground(colorViolet).Background(colorBase02)
	styleInputCursor = tcell.StyleDefault.Foreground(colorBase02).Background(colorViolet)

	stylePreview = tcell.StyleDefault.Foreground(colorBase0).Background(colorBase03)
	styleBar     = tcell.StyleDefault.Foreground(colorGreen).Background(colorBase02)

	styleExpression      = tcell.StyleDefault.Foreground(colorViolet).Background(colorBase03)
	styleExpressionError = tcell.StyleDefault.Foreground(colorRed).Background(colorBase03)
)
