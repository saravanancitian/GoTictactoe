package tictactoe

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CELL_SIZE       int = 62
	GAMEOVER_WIDTH  int = 159
	GAMEOVER_HEIGHT int = 93
	MARGIN          int = 4
)

type Board struct {
	x               int
	y               int
	width           int
	height          int
	cellsize        int
	margin          int
	cellwithmargine int

	bgColor     color.RGBA
	marginColor color.RGBA
}

func (b *Board) SetXY(startx, starty int) {

	b.x = startx
	b.y = starty
}

func (b *Board) Init() {

	b.cellsize = CELL_SIZE
	b.margin = MARGIN

	b.cellwithmargine = b.margin + b.cellsize

	b.width = b.cellwithmargine*NUM_COL + b.margin
	b.height = b.width

	b.bgColor = color.RGBA{255, 245, 184, 0xff}
	b.marginColor = color.RGBA{190, 145, 51, 0xff}

}

func (b *Board) GetXY(row, col int) (int, int) {

	var x, y int = 0, 0

	x = b.x + (row * b.cellwithmargine) + b.cellsize/2
	y = b.y + (col * b.cellwithmargine) + b.cellsize/2

	return x, y

}

func (b *Board) GetSelectedCell(x, y int) (int, int) {

	if x > b.x && x < (b.x+b.width-b.margin) && y > b.y && y < (b.y+b.height-b.margin) {
		var xfrombx = x - b.x
		var yfromby = y - b.y

		row := xfrombx / b.cellwithmargine
		col := yfromby / b.cellwithmargine

		return row, col
	}

	return -1, -1
}

func (b *Board) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), b.bgColor, false)
	vector.StrokeRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), float32(b.margin), b.marginColor, false)

	vector.StrokeLine(screen, float32(b.x+b.cellwithmargine), float32(b.y), float32(b.x+b.cellwithmargine), float32(b.y+b.height), float32(b.margin), b.marginColor, false)
	vector.StrokeLine(screen, float32(b.x+b.cellwithmargine*2), float32(b.y), float32(b.x+b.cellwithmargine*2), float32(b.y+b.height), float32(b.margin), b.marginColor, false)
	vector.StrokeLine(screen, float32(b.x), float32(b.y+b.cellwithmargine), float32(b.x+b.width), float32(b.y+b.cellwithmargine), float32(b.margin), b.marginColor, false)
	vector.StrokeLine(screen, float32(b.x), float32(b.y+b.cellwithmargine*2), float32(b.x+b.width), float32(b.y+b.cellwithmargine*2), float32(b.margin), b.marginColor, false)
}

func NewBoard() *Board {
	var board *Board = new(Board)
	board.Init()
	return board
}
