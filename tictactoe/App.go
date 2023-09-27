package tictactoe

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	NUM_CIRCLE_FRAMES int = 8
	NUM_CROSS_FRAMES  int = 4
)

const (
	APP_STATE_INIT    int = 1
	APP_STATE_RUNNING int = 2
)

type App struct {
	board *Board

	screenWidth  int
	screenHeight int

	prevTime    int64
	curTime     int64
	state       int
	scalefactor float64
}

func (app *App) Init() {

	app.scalefactor = 1
	app.screenWidth = 204
	app.screenHeight = 204
	app.state = APP_STATE_INIT
}

func (app *App) Update() error {
	switch app.state {
	case APP_STATE_INIT:
		app.board = NewBoard(app.scalefactor, app.screenWidth, app.screenHeight)
		app.state = APP_STATE_RUNNING
	case APP_STATE_RUNNING:
		app.curTime = time.Now().UnixMilli()
		delta := app.curTime - app.prevTime
		app.prevTime = app.curTime
		app.board.Update(delta)
	}

	return nil
}

func (app *App) ProcessEvent() {

}

func (app *App) Draw(screen *ebiten.Image) {
	app.board.Draw(screen)
}

func (app *App) Layout(ow, oh int) (int, int) {
	return  app.screenWidth,  app.screenHeight
}

func NewApp() *App {
	var app *App = new(App)

	app.Init()

	return app
}
