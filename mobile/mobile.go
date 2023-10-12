package mobile

import (
	"tictactoe/tictactoe"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

type IGameCallback interface {
	GameOverCallBack(winner int, duration int64)
}

var game *tictactoe.App

////go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile  bind -target android -javapkg com.tictactoe.tictactoe -o ./bin/tictactoe.aar .\mobile

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game

	game = tictactoe.NewApp()

	mobile.SetGame(game)
}

func RegisterGameCallback(callback IGameCallback) {
	game.RegisterIGameCallback(func(winner int, duration int64) { callback.GameOverCallBack(winner, duration) })
}

func PlayAgain(ngameplayed, nwin int) {
	game.PlayAgain(ngameplayed, nwin)
}

func Pause() {
	game.Pause()
}

func Resume() {
	game.Resume()
}

func Destroy() {
	game.Destroy()
}

func SetSoundOff(off bool) {
	game.SetSoundOff(off)
}

func SetShowTimerOff(off bool) {
	game.SetShowTimerOff(off)
}
