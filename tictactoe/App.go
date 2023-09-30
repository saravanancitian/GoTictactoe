package tictactoe

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const(
	INNER_WIDTH int = 204
	INNER_HEIGHT int = 250
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
	ttt *TicTacToe

	gcallback func()

	screenWidth  int
	screenHeight int

	prevTime    int64
	curTime     int64
	state       int
	scalefactor float64

	rm *ResourceManager
}

func (app *App) Init() {
	app.rm = NewResourceManager()
	app.scalefactor = 1
	app.state = APP_STATE_INIT
}

func (app *App) Update() error {
	switch app.state {
	case APP_STATE_INIT:
		app.ttt = NewTicTacToe(app.rm, app.screenWidth, app.screenHeight, app.gcallback)
		app.state = APP_STATE_RUNNING
	case APP_STATE_RUNNING:
		app.curTime = time.Now().UnixMilli()
		delta := app.curTime - app.prevTime
		app.prevTime = app.curTime
		app.ttt.Update(delta)
	}

	return nil
}




func (app *App) Draw(screen *ebiten.Image) {
	app.ttt.Draw(screen)
}

func (app *App) Layout(ow, oh int) (int, int) {

	app.screenWidth = INNER_WIDTH
	app.screenHeight = INNER_HEIGHT
	return  INNER_WIDTH, INNER_HEIGHT
}

func (app *App) RegisterIGameCallback(callback func()){
	app.gcallback = callback
}

func NewApp() *App {
	var app *App = new(App)
	app.Init()
	return app
}
