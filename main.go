package main

import (
	"tictactoe/tictactoe"

	"github.com/hajimehoshi/ebiten/v2"
)

const SCREEN_WIDTH int = 400
const SCREEN_HEIGHT int = 400

func main() {
	ebiten.SetWindowTitle("Tic Tac Toe")
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.RunGame(tictactoe.NewApp())
}
