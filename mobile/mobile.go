package mobile

import (
	"tictactoe/tictactoe"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

////go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile  bind -target android -javapkg com.tictactoe.tictactoe -o ./bin/tictactoe.aar .\mobile

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game

	var game= tictactoe.NewApp()

	mobile.SetGame(game)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
