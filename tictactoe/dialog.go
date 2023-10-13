package tictactoe

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Dialog struct {
	width        int
	height       int
	text         string
	bgcolor      color.RGBA
	margincolor  color.RGBA
	marginStroke int
	txtFont      font.Face
}

func (d *Dialog) Draw(screen *ebiten.Image, x int, y int) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(d.width), float32(d.height), d.bgcolor, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(d.width), float32(d.height), 4, color.White, false)
	text.Draw(screen, d.text, d.txtFont, x+d.width/2, y+d.height/2, color.Black)

}

func NewDialog(w int, h int, txt string, bgcolor color.RGBA, mcolor color.RGBA, mstroke int, fontface font.Face) *Dialog {
	var dialog *Dialog = new(Dialog)

	dialog.width = w
	dialog.height = h
	dialog.text = txt
	dialog.bgcolor = bgcolor
	dialog.margincolor = mcolor
	dialog.marginStroke = mstroke
	dialog.txtFont = fontface

	return dialog
}
