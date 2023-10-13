package tictactoe

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	INNER_WIDTH  int = 204
	INNER_HEIGHT int = 250
)

const (
	NUM_CIRCLE_FRAMES int = 8
	NUM_CROSS_FRAMES  int = 4
)

const (
	APP_STATE_INIT      int = 1
	APP_STATE_RUNNING   int = 2
	APP_STATE_PAUSED    int = 3
	APP_STATE_DESTROYED int = 4
)

const (
	WINDOWSTATE_MAXIMIZED int = 1
	WINDOWSTATE_MINIMIZED int = 2
)

type App struct {
	ttt *TicTacToe

	gameovercallback func(int, int64)

	screenWidth  int
	screenHeight int

	prevTime            int64
	curTime             int64
	state               int
	windowState         int
	prevState           int
	scalefactor         float64
	settingSoundoff     bool
	settingShowTimeroff bool

	rm *ResourceManager
}

func (app *App) Init() {
	app.rm = NewResourceManager()
	app.scalefactor = 1
	app.prevState = -1
	app.state = -1
	app.windowState = WINDOWSTATE_MAXIMIZED
	app.SetState(APP_STATE_INIT)
}

func (app *App) Update() error {
	if ebiten.IsWindowBeingClosed() {
		app.Destroy()
	} else if ebiten.IsWindowMinimized() {
		app.Pause()
		app.windowState = WINDOWSTATE_MINIMIZED
	} else if app.windowState == WINDOWSTATE_MINIMIZED {
		app.windowState = WINDOWSTATE_MAXIMIZED
		app.Resume()
	}

	switch app.state {
	case APP_STATE_INIT:
		app.ttt = NewTicTacToe(app.rm, app.screenWidth, app.screenHeight, app.gameovercallback)
		app.ttt.SetShowTimerOff(app.settingShowTimeroff)
		app.ttt.SetSoundOff(app.settingSoundoff)
		app.SetState(APP_STATE_RUNNING)
	case APP_STATE_RUNNING:
		app.curTime = time.Now().UnixMilli()
		delta := app.curTime - app.prevTime
		app.prevTime = app.curTime
		app.ttt.Update(delta)
	case APP_STATE_PAUSED:
		app.prevTime = time.Now().UnixMilli()
		app.curTime = time.Now().UnixMilli()
	}

	return nil
}

func (app *App) SetState(state int) {
	if app.state != state {
		app.prevState = app.state
		app.state = state
	}
}

func (app *App) Pause() {
	if app.state != APP_STATE_PAUSED {

		app.SetState(APP_STATE_PAUSED)
	}
}

func (app *App) Resume() {
	if app.state == APP_STATE_PAUSED {
		app.prevTime = time.Now().UnixMilli()
		app.curTime = time.Now().UnixMilli()
		app.SetState(app.prevState)
	}
}

func (app *App) Destroy() {
	app.state = APP_STATE_DESTROYED
	app.rm.UnloadResources()
}

func (app *App) Draw(screen *ebiten.Image) {
	if app.state != APP_STATE_DESTROYED {
		app.ttt.Draw(screen)
	}
}

func (app *App) Layout(ow, oh int) (int, int) {

	app.screenWidth = INNER_WIDTH
	app.screenHeight = INNER_HEIGHT
	return INNER_WIDTH, INNER_HEIGHT
}

func (app *App) RegisterIGameCallback(callback func(int, int64)) {
	app.gameovercallback = callback
}

func (app *App) PlayAgain(ngameplayed, nwin int) {
	app.ttt.StartNewGame(ngameplayed, nwin)
}

func (app *App) SetSoundOff(off bool) {
	app.settingSoundoff = off
	if app.ttt != nil {
		app.ttt.SetSoundOff(off)
	}
}

func (app *App) SetShowTimerOff(off bool) {
	app.settingShowTimeroff = off
	if app.ttt != nil {
		app.ttt.SetShowTimerOff(off)
	}
}

func NewApp() *App {
	var app *App = new(App)
	app.Init()
	return app
}
