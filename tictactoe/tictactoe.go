package tictactoe

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"tictactoe/tictactoe/input"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	STATE_GAME_OVER_HALT       int = 0
	STATE_INIT_NEW_GAME        int = 1
	STATE_AI_PLAYER_TURN       int = 2
	STATE_HUMAN_PLAYER_TURN    int = 3
	STATE_GAME_OVER            int = 4
	STATE_ANIMATE_AI_PLAYER    int = 5
	STATE_ANIMATE_HUMAN_PLAYER int = 6
	GAME_TIDE                  int = 0
)

// const (
// 	TURN1 int = iota
// 	TURN2
// 	TURN3
// 	TURN4
// 	TURN5
// 	TURN6
// 	TURN7
// 	TURN8
// 	TURN9
// )

const (
	HUMAN_PLAYER int = 1
	AI_PLAYER    int = -1

	ANIMATE_HUMAN_PLAYER int = 2
	ANIMATE_AI_PLAYER    int = -2

	AI_TYPE_BELOW_AVERAGE int = 0
	AI_TYPE_AVERAGE       int = 1
	AI_TYPE_GOOD          int = 2
)

const (
	NUM_COL                   int = 3
	NUM_ROW                   int = 3
	MAX_MOVES_CHECK_GAME_OVER int = 5

	TXT_YOUR_TURN string = "Your Turn"

	TXT_YOU_WON string = "You Won"
	TXT_I_WON   string = "I Won"
	TXT_TIED    string = "Tied"

	ANIM_DELAY    int64 = 100
	SEC_IN_MILLIS int64 = 1000
)

type TicTacToe struct {
	GameOverCallBack func(int, int64)

	youwin   *Dialog
	gamelost *Dialog
	gametide *Dialog

	winSnd   *audio.Player
	loseSnd  *audio.Player
	crossSnd *audio.Player
	roundSnd *audio.Player

	gx int
	gy int

	gameScreenWidth  int
	gameScreenHeight int

	cell [NUM_ROW][NUM_COL]int

	JustAnotherHand_ttf *opentype.Font
	normalFont          font.Face
	board               *Board
	player              *Player

	overlayColor color.RGBA
	rm           *ResourceManager

	state int

	numHumanMove int
	numAIMove    int

	winner int

	playtime     int64
	playtimecalc int64

	animtime    int64
	isAnimating bool
	isGameover  bool

	animRow int
	animCol int

	txtTimer  string
	showTimer bool
	showMsg   bool

	settingTimerOff bool
	settingSoundOff bool

	txtMsg string

	random        *rand.Rand
	prevTurnStart int
	// toggleRand    bool

	s1 image.Point
	s2 image.Point

	totalGamePlayed int
	totalWin        int
}

func (t *TicTacToe) SetSoundOff(off bool) {
	t.settingSoundOff = off
}

func (t *TicTacToe) SetShowTimerOff(off bool) {
	t.settingTimerOff = off
}

func (t *TicTacToe) SetCallback(callback func(int, int64)) {
	t.GameOverCallBack = callback
}

func (t *TicTacToe) LoadSnd(rm *ResourceManager) {
	var err error
	t.crossSnd, err = rm.LoadMp3Audio("cross.mp3")
	if err != nil {
		log.Fatal(err)
	}

	t.roundSnd, err = rm.LoadMp3Audio("round.mp3")
	if err != nil {
		log.Fatal(err)
	}

	t.winSnd, err = rm.LoadMp3Audio("win.mp3")
	if err != nil {
		log.Fatal(err)
	}

	t.loseSnd, err = rm.LoadMp3Audio("lose.mp3")
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TicTacToe) LoadSprite(rm *ResourceManager) {
	var err error

	t.JustAnotherHand_ttf, err = rm.LoadFont("JustAnotherHand.ttf")
	if err != nil {
		log.Fatal(err)
	}

	t.normalFont, err = opentype.NewFace(t.JustAnotherHand_ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func NewTicTacToe(rm *ResourceManager, screenWidth int, screenHeight int, callback func(int, int64)) *TicTacToe {
	var tictactoe = new(TicTacToe)
	tictactoe.Init(rm, screenWidth, screenHeight, callback)
	return tictactoe
}

func (t *TicTacToe) Init(rm *ResourceManager, screenWidth int, screenHeight int, callback func(int, int64)) {

	t.settingSoundOff = false
	t.settingTimerOff = false
	t.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	t.rm = rm
	t.GameOverCallBack = callback
	t.LoadSprite(rm)
	t.LoadSnd(rm)
	t.gx = 0
	t.gy = 0
	t.gameScreenWidth = screenWidth
	t.gameScreenHeight = screenHeight

	t.board = NewBoard()
	t.player = NewPlayer()
	t.player.LoadSprite(rm)

	t.youwin = NewDialog(150, 100, TXT_YOU_WON, color.RGBA{251, 206, 105, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}, 4, t.normalFont)
	t.gamelost = NewDialog(150, 100, TXT_I_WON, color.RGBA{251, 206, 105, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}, 4, t.normalFont)
	t.gametide = NewDialog(150, 100, TXT_TIED, color.RGBA{251, 206, 105, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff}, 4, t.normalFont)

	x := (screenWidth - t.board.width) / 2
	y := (screenHeight - t.board.height) / 2
	t.board.SetXY(x, y)

	t.overlayColor = color.RGBA{50, 50, 50, 150}

	t.prevTurnStart = 0
	t.totalGamePlayed = 0
	t.totalWin = 0
	t.StartNewGame(t.totalGamePlayed, t.totalWin)
}

func (t *TicTacToe) StartNewGame(ngameplayed, nwin int) {
	t.state = STATE_INIT_NEW_GAME
	t.totalGamePlayed = ngameplayed
	t.totalWin = nwin
}

func (t *TicTacToe) ResetBoard() {
	for i := 0; i < NUM_ROW; i++ {
		for j := 0; j < NUM_COL; j++ {
			t.cell[i][j] = 0
		}
	}
}

func (t *TicTacToe) CalculatePlayTime(delta int64) {
	if !t.isGameover {
		if t.playtimecalc >= SEC_IN_MILLIS {
			t.playtime += t.playtimecalc
			t.playtimecalc = delta

			tsec := t.playtime / SEC_IN_MILLIS
			min := tsec / 60
			sec := tsec % 60
			t.txtTimer = fmt.Sprintf("\n%02d:%02d", min, sec)

		} else {
			t.playtimecalc += delta
		}
	}
}

func (t *TicTacToe) setStrikeLine(row1, col1, row2, col2 int) {

	var x1, y1 = t.board.GetXY(row1, col1)

	var x2, y2 = t.board.GetXY(row2, col2)

	if row1 == 0 && col1 == 0 && row2 == 0 && col2 == 2 {
		//(0, 0, 0, 2)
		x1 = x1 + t.board.cellsize/4
		y1 = y1 + t.board.cellsize/2

		x2 = x2 + t.board.cellsize/2 + t.board.cellsize/4
		y2 = y2 + t.board.cellsize/2
	} else if row1 == 0 && col1 == 0 && row2 == 2 && col2 == 0 {
		//(0, 0, 2, 0)
		x1 = x1 + t.board.cellsize/2
		y1 = y1 + t.board.cellsize/4

		x2 = x2 + t.board.cellsize/2
		y2 = y2 + t.board.cellsize/2 + t.board.cellsize/4
	} else if row1 == 0 && col1 == 0 && row2 == 2 && col2 == 2 {
		//(0, 0, 2, 2)
		x1 = x1 + t.board.cellsize/4
		y1 = y1 + t.board.cellsize/4

		x2 = x2 + t.board.cellsize/2 + t.board.cellsize/4
		y2 = y2 + t.board.cellsize/2 + t.board.cellsize/4
	} else if row1 == 1 && col1 == 0 && row2 == 1 && col2 == 2 {
		//(1, 0, 1, 2)
		x1 = x1 + t.board.cellsize/4
		y1 = y1 + t.board.cellsize/2

		x2 = x2 + t.board.cellsize/2 + t.board.cellsize/4
		y2 = y2 + t.board.cellsize/2
	} else if row1 == 2 && col1 == 0 && row2 == 2 && col2 == 2 {
		//(2, 0, 2, 2)
		x1 = x1 + t.board.cellsize/4
		y1 = y1 + t.board.cellsize/2

		x2 = x2 + t.board.cellsize/2 + t.board.cellsize/4
		y2 = y2 + t.board.cellsize/2
	} else if row1 == 2 && col1 == 0 && row2 == 0 && col2 == 2 {
		//(2, 0, 0, 2)
		x1 = x1 + t.board.cellsize/4
		y1 = y1 + t.board.cellsize/2 + t.board.cellsize/4

		x2 = x2 + t.board.cellsize/2 + t.board.cellsize/4
		y2 = y2 + t.board.cellsize/4
	} else if row1 == 0 && col1 == 1 && row2 == 2 && col2 == 1 {
		//(0, 1, 2, 1)
		x1 = x1 + t.board.cellsize/2
		y1 = y1 + t.board.cellsize/4

		x2 = x2 + t.board.cellsize/2
		y2 = y2 + t.board.cellsize/2 + t.board.cellsize/4
	} else if row1 == 0 && col1 == 2 && row2 == 2 && col2 == 2 {
		//(0, 2, 2, 2)
		x1 = x1 + t.board.cellsize/2
		y1 = y1 + t.board.cellsize/4

		x2 = x2 + t.board.cellsize/2
		y2 = y2 + t.board.cellsize/2 + t.board.cellsize/4
	}

	t.s1 = image.Point{
		X: x1,
		Y: y1,
	}

	t.s2 = image.Point{
		X: x2,
		Y: y2,
	}
}

func (t *TicTacToe) CheckGameOver() (bool, int) {
	var isGameOver bool = false
	var winner int = 0
	numMove := t.numHumanMove + t.numAIMove
	if numMove >= MAX_MOVES_CHECK_GAME_OVER {
		if (t.cell[0][0] != 0) && (t.cell[0][0] == t.cell[0][1] && t.cell[0][1] == t.cell[0][2] || t.cell[0][0] == t.cell[1][0] && t.cell[1][0] == t.cell[2][0] || t.cell[0][0] == t.cell[1][1] && t.cell[1][1] == t.cell[2][2]) {
			isGameOver = true
			winner = t.cell[0][0]

			if t.cell[0][0] == t.cell[0][1] && t.cell[0][1] == t.cell[0][2] {
				t.setStrikeLine(0, 0, 0, 2) //
			} else if t.cell[0][0] == t.cell[1][0] && t.cell[1][0] == t.cell[2][0] {
				t.setStrikeLine(0, 0, 2, 0) //
			} else if t.cell[0][0] == t.cell[1][1] && t.cell[1][1] == t.cell[2][2] {
				t.setStrikeLine(0, 0, 2, 2) //
			}

		} else if (t.cell[1][0] != 0) && (t.cell[1][0] == t.cell[1][1] && t.cell[1][1] == t.cell[1][2]) {
			isGameOver = true
			winner = t.cell[1][0]
			t.setStrikeLine(1, 0, 1, 2) //

		} else if (t.cell[2][0] != 0) && (t.cell[2][0] == t.cell[2][1] && t.cell[2][1] == t.cell[2][2] || t.cell[2][0] == t.cell[1][1] && t.cell[1][1] == t.cell[0][2]) {
			isGameOver = true
			winner = t.cell[2][0]
			if t.cell[2][0] == t.cell[2][1] && t.cell[2][1] == t.cell[2][2] {
				t.setStrikeLine(2, 0, 2, 2) //
			} else if t.cell[2][0] == t.cell[1][1] && t.cell[1][1] == t.cell[0][2] {
				t.setStrikeLine(2, 0, 0, 2) //
			}

		} else if (t.cell[0][1] != 0) && (t.cell[0][1] == t.cell[1][1] && t.cell[1][1] == t.cell[2][1]) {
			isGameOver = true
			winner = t.cell[0][1]
			t.setStrikeLine(0, 1, 2, 1) //

		} else if (t.cell[0][2] != 0) && (t.cell[0][2] == t.cell[1][2] && t.cell[1][2] == t.cell[2][2]) {
			isGameOver = true
			winner = t.cell[0][2]
			t.setStrikeLine(0, 2, 2, 2)
		}

		if !isGameOver && numMove == (NUM_COL*NUM_ROW) {
			isGameOver = true
			winner = GAME_TIDE
		}

	}
	return isGameOver, winner
}

func (t *TicTacToe) Draw(screen *ebiten.Image) {

	t.board.Draw(screen)

	if t.showMsg {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.board.x+t.board.width-60), float64(t.board.y-5))
		text.DrawWithOptions(screen, t.txtMsg, t.normalFont, op)
	}

	if t.showTimer {
		if !t.settingTimerOff {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.board.x+5), float64(t.board.y-30))
			text.DrawWithOptions(screen, t.txtTimer, t.normalFont, op)
		}
	}

	var x int = 0
	var y int = 0

	for i := 0; i < NUM_ROW; i++ {
		y = t.board.y + (t.board.cellwithmargine * i) + t.board.margin

		for j := 0; j < NUM_COL; j++ {
			x = t.board.x + (t.board.cellwithmargine * j) + t.board.margin

			playerType := t.cell[i][j]

			t.player.Draw(screen, playerType, x, y)
		}
	}

	if t.isGameover {
		if t.winner != GAME_TIDE {
			vector.StrokeLine(screen, float32(t.s1.X), float32(t.s1.Y), float32(t.s2.X), float32(t.s2.Y), 5, color.Black, false)
		}
		if !IsMobileBuild() {
			t.DrawGameOver(screen)
		}
	}
}

func (t *TicTacToe) DrawGameOver(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(t.board.x), float32(t.board.y), float32(t.board.width), float32(t.board.height), t.overlayColor, false)
	var dialog *Dialog
	if t.winner == HUMAN_PLAYER {
		dialog = t.youwin
	} else if t.winner == AI_PLAYER {
		dialog = t.gamelost
	} else {
		dialog = t.gametide
	}

	gx := t.board.x + (t.board.width-GAMEOVER_WIDTH)/2
	gy := t.board.y + (t.board.height-GAMEOVER_HEIGHT)/2

	dialog.Draw(screen, gx, gy)
}

func (t *TicTacToe) playSound(sndPlayer *audio.Player) {
	if !t.settingSoundOff {
		sndPlayer.Rewind()
		sndPlayer.Play()
	}
}

func (t *TicTacToe) SetStartTurn() {
	var turn int = AI_PLAYER
	if t.winner == GAME_TIDE {
		if t.prevTurnStart == 0 {
			var trn = []int{AI_PLAYER, HUMAN_PLAYER}
			turn = trn[t.random.Intn(2)]

		} else if t.prevTurnStart == AI_PLAYER {
			turn = HUMAN_PLAYER
		} else if t.prevTurnStart == HUMAN_PLAYER {
			turn = AI_PLAYER
		}
	} else if t.winner == AI_PLAYER {
		turn = AI_PLAYER
	} else if t.winner == HUMAN_PLAYER {
		turn = HUMAN_PLAYER
	}

	if turn == AI_PLAYER {
		t.state = STATE_AI_PLAYER_TURN
		t.prevTurnStart = AI_PLAYER
		t.showMsg = false

		// t.player.aiType = AI_TYPE_AVERAGE

		// if t.toggleRand {
		// 	var aitypes = []int{AI_TYPE_BELOW_AVERAGE, AI_TYPE_AVERAGE, AI_TYPE_GOOD}
		// 	t.player.aiType = aitypes[t.random.Intn(3)]
		// }
		// t.toggleRand = !t.toggleRand

	} else if turn == HUMAN_PLAYER {
		t.state = STATE_HUMAN_PLAYER_TURN
		t.prevTurnStart = HUMAN_PLAYER
		t.showMsg = true

	}

}

func (t *TicTacToe) Update(delta int64) {
	switch t.state {
	case STATE_GAME_OVER_HALT:
		if !IsMobileBuild() {
			if t.DelayElapsed(delta) {
				t.StartNewGame(t.totalGamePlayed, t.totalWin)
			}
		}

	case STATE_INIT_NEW_GAME:
		t.ResetBoard()
		t.SetStartTurn()
		t.numAIMove = 0
		t.numHumanMove = 0
		t.winner = -1000
		t.playtime = 0
		t.playtimecalc = 0
		t.animtime = 0
		t.isAnimating = false
		t.isGameover = false
		t.animRow = -1
		t.animCol = -1
		t.txtTimer = "00:00"
		t.showTimer = true
		t.txtMsg = TXT_YOUR_TURN

		t.winSnd.Rewind()
		t.crossSnd.Rewind()
		t.loseSnd.Rewind()
		t.roundSnd.Rewind()

	case STATE_AI_PLAYER_TURN:
		t.CalculatePlayTime(delta)
		var seletaitype int
		if t.totalWin > 100 {
			seletaitype = AI_TYPE_GOOD
		} else if t.totalGamePlayed > 150 || t.totalWin > 75 || t.numHumanMove > 2 {
			var aitypes = []int{AI_TYPE_AVERAGE, AI_TYPE_GOOD}
			seletaitype = aitypes[t.random.Intn(2)]
		} else {
			var aitypes = []int{AI_TYPE_BELOW_AVERAGE, AI_TYPE_AVERAGE, AI_TYPE_GOOD}
			seletaitype = aitypes[t.random.Intn(3)]
		}

		if seletaitype == t.player.aiType && (seletaitype == AI_TYPE_AVERAGE || seletaitype == AI_TYPE_BELOW_AVERAGE) {
			seletaitype = AI_TYPE_GOOD
		}
		t.player.aiType = seletaitype

		row, col := t.GetAIMove(t.player.aiType)
		if row > -1 && col > -1 {
			if t.cell[row][col] == 0 {
				t.numAIMove++
				t.cell[row][col] = ANIMATE_AI_PLAYER
				t.animRow = row
				t.animCol = col
				t.player.circleImgIdx = 0
				t.state = STATE_ANIMATE_AI_PLAYER
				go t.playSound(t.roundSnd)
			}
		}

	case STATE_HUMAN_PLAYER_TURN:
		t.CalculatePlayTime(delta)

		mx, my := input.Current().GetPosition()
		if mx >= 0 && my >= 0 {
			row, col := t.board.GetSelectedCell(mx, my)
			if row > -1 && col > -1 {
				if t.cell[row][col] == 0 {
					t.numHumanMove++
					t.cell[row][col] = ANIMATE_HUMAN_PLAYER
					t.animRow = row
					t.animCol = col
					t.player.crossImgIdx = 0
					t.state = STATE_ANIMATE_HUMAN_PLAYER
					t.showMsg = false
					go t.playSound(t.crossSnd)
				}
			}
		}

	case STATE_GAME_OVER:
		t.totalGamePlayed++

		t.isGameover = true
		t.showMsg = true
		t.txtMsg = TXT_TIED
		if t.winner == HUMAN_PLAYER {
			t.totalWin++
			t.txtMsg = TXT_YOU_WON
			go t.playSound(t.winSnd)
		} else if t.winner == AI_PLAYER {
			t.txtMsg = TXT_I_WON
			go t.playSound(t.loseSnd)

		} else {
			go t.playSound(t.winSnd)

		}

		if IsMobileBuild() && t.GameOverCallBack != nil {
			t.GameOverCallBack(t.winner, t.playtime)
		} else {
			t.SetDelay(5 * SEC_IN_MILLIS)
		}
		t.state = STATE_GAME_OVER_HALT

	case STATE_ANIMATE_AI_PLAYER:
		t.CalculatePlayTime(delta)
		if t.DelayElapsed(delta) {
			if t.player.circleImgIdx >= NUM_CIRCLE_FRAMES-1 {
				if !t.roundSnd.IsPlaying() {
					t.isAnimating = false
					t.cell[t.animRow][t.animCol] = AI_PLAYER
					isGameOver, winner := t.CheckGameOver()
					if isGameOver {
						t.winner = winner

						t.state = STATE_GAME_OVER
					} else {
						t.showMsg = true

						t.state = STATE_HUMAN_PLAYER_TURN
					}
				}

			} else {
				t.player.circleImgIdx++
				t.animtime = ANIM_DELAY
			}
		}

	case STATE_ANIMATE_HUMAN_PLAYER:
		t.CalculatePlayTime(delta)
		if t.DelayElapsed(delta) {
			if t.player.crossImgIdx >= NUM_CROSS_FRAMES-1 {
				if !t.crossSnd.IsPlaying() {
					t.isAnimating = false
					t.cell[t.animRow][t.animCol] = HUMAN_PLAYER

					isGameOver, winner := t.CheckGameOver()
					if isGameOver {
						t.winner = winner

						t.state = STATE_GAME_OVER
					} else {
						t.state = STATE_AI_PLAYER_TURN
					}
				}

			} else {
				t.player.crossImgIdx++
				t.animtime = ANIM_DELAY
			}
		}
	}
}

func (t *TicTacToe) SetAnimation() {
	t.isAnimating = true
	t.SetDelay(ANIM_DELAY)
}

func (t *TicTacToe) SetDelay(delay int64) {
	t.animtime = delay

}

func (t *TicTacToe) DelayElapsed(delta int64) bool {
	if t.animtime-delta > 0 {
		t.animtime -= delta
	} else {
		return true
	}
	return false

}

func (t *TicTacToe) CanPlayerWin(row int, col int, playerType int) bool {
	var isWin bool = false

	if row == 0 && col == 0 {
		//check col, row, dia
		isWin = (t.cell[1][0] == playerType && t.cell[2][0] == playerType || t.cell[0][1] == playerType && t.cell[0][2] == playerType || t.cell[1][1] == playerType && t.cell[2][2] == playerType)
	} else if row == 0 && col == 2 {
		isWin = (t.cell[1][2] == playerType && t.cell[2][2] == playerType || t.cell[0][1] == playerType && t.cell[0][0] == playerType || t.cell[1][1] == playerType && t.cell[2][0] == playerType)
	} else if row == 2 && col == 0 {
		isWin = (t.cell[1][0] == playerType && t.cell[0][0] == playerType || t.cell[2][1] == playerType && t.cell[2][2] == playerType || t.cell[1][1] == playerType && t.cell[0][2] == playerType)
	} else if row == 2 && col == 2 {
		isWin = (t.cell[1][2] == playerType && t.cell[0][2] == playerType || t.cell[2][1] == playerType && t.cell[2][0] == playerType || t.cell[1][1] == playerType && t.cell[0][0] == playerType)
	} else if row == 1 && col == 1 {
		isWin = (t.cell[0][1] == playerType && t.cell[2][1] == playerType || t.cell[1][0] == playerType && t.cell[1][2] == playerType || t.cell[0][0] == playerType && t.cell[2][2] == playerType || t.cell[2][0] == playerType && t.cell[0][2] == playerType)
	} else if row == 0 && col == 1 {
		isWin = (t.cell[0][0] == playerType && t.cell[0][2] == playerType || t.cell[1][1] == playerType && t.cell[2][1] == playerType)
	} else if row == 1 && col == 0 {
		isWin = (t.cell[1][1] == playerType && t.cell[1][2] == playerType || t.cell[0][0] == playerType && t.cell[2][0] == playerType)
	} else if row == 2 && col == 1 {
		isWin = (t.cell[2][0] == playerType && t.cell[2][2] == playerType || t.cell[1][1] == playerType && t.cell[0][1] == playerType)
	} else if row == 1 && col == 2 {
		isWin = (t.cell[1][1] == playerType && t.cell[1][0] == playerType || t.cell[0][2] == playerType && t.cell[2][2] == playerType)
	}

	return isWin
}

func (t *TicTacToe) GetEmptyRowCol() (int, int) {
	var row int = -1
	var col int = -1
erc:
	for currow := 0; currow < NUM_ROW; currow++ {
		for curcol := 0; curcol < NUM_COL; curcol++ {
			if t.cell[currow][curcol] == 0 {
				r, c := t.GetNearPlace(currow, curcol)
				if r != -1 && c != -1 {
					row = currow
					col = curcol
					break erc
				}
			}
		}
	}

	return row, col
}

func (t *TicTacToe) GetNearPlace(currow int, curcol int) (int, int) {
	var row int = -1
	var col int = -1

	if currow == 0 && curcol == 0 {
		if t.cell[0][1] == 0 && t.cell[0][2] == 0 {
			row = 0
			col = 1
		} else if t.cell[1][0] == 0 && t.cell[2][0] == 0 {
			row = 1
			col = 0
		} else if t.cell[1][1] == 0 && t.cell[2][2] == 0 {
			row = 1
			col = 1
		}
	} else if currow == 0 && curcol == 2 {
		if t.cell[0][1] == 0 && t.cell[0][0] == 0 {
			row = 0
			col = 1
		} else if t.cell[1][2] == 0 && t.cell[2][2] == 0 {
			row = 1
			col = 2
		} else if t.cell[1][1] == 0 && t.cell[2][0] == 0 {
			row = 1
			col = 1
		}
	} else if currow == 2 && curcol == 0 {
		if t.cell[2][1] == 0 && t.cell[2][2] == 0 {
			row = 2
			col = 1
		} else if t.cell[1][0] == 0 && t.cell[0][0] == 0 {
			row = 1
			col = 0
		} else if t.cell[1][1] == 0 && t.cell[0][2] == 0 {
			row = 1
			col = 1
		}
	} else if currow == 2 && curcol == 2 {
		if t.cell[2][1] == 0 && t.cell[2][0] == 0 {
			row = 2
			col = 1
		} else if t.cell[1][2] == 0 && t.cell[0][2] == 0 {
			row = 1
			col = 2
		} else if t.cell[1][1] == 0 && t.cell[0][0] == 0 {
			row = 1
			col = 1
		}
	} else if currow == 1 && curcol == 1 {
		if t.cell[1][0] == 0 && t.cell[1][2] == 0 {
			row = 1
			col = 0
		} else if t.cell[0][1] == 0 && t.cell[2][1] == 0 {
			row = 0
			col = 1
		} else if t.cell[0][0] == 0 && t.cell[2][2] == 0 {
			row = 0
			col = 0
		} else if t.cell[2][0] == 0 && t.cell[0][2] == 0 {
			row = 2
			col = 0
		}
	} else if currow == 0 && curcol == 1 {
		if t.cell[0][0] == 0 && t.cell[0][2] == 0 {
			row = 0
			col = 0
		} else if t.cell[1][1] == 0 && t.cell[2][1] == 0 {
			row = 1
			col = 1
		}

	} else if currow == 1 && curcol == 0 {
		if t.cell[1][1] == 0 && t.cell[1][2] == 0 {
			row = 1
			col = 1
		} else if t.cell[0][0] == 0 && t.cell[2][0] == 0 {
			row = 0
			col = 0
		}
	} else if currow == 2 && curcol == 1 {
		if t.cell[2][0] == 0 && t.cell[2][2] == 0 {
			row = 2
			col = 0
		} else if t.cell[1][1] == 0 && t.cell[0][1] == 0 {
			row = 1
			col = 1
		}
	} else if currow == 1 && curcol == 2 {
		if t.cell[1][1] == 0 && t.cell[1][0] == 0 {
			row = 1
			col = 1
		} else if t.cell[0][2] == 0 && t.cell[2][2] == 0 {
			row = 0
			col = 2
		}
	}
	return row, col
}

func (t *TicTacToe) GetNextEmptyCell() (int, int) {

	var row int = -1
	var col int = -1

nec:
	for i := 0; i < NUM_ROW; i++ {
		for j := 0; j < NUM_COL; j++ {
			if t.cell[i][j] == 0 {
				row = i
				col = j
				break nec
			}
		}
	}

	return row, col
}

func (t *TicTacToe) GetRandEmptyCell() (int, int) {
	var row int = -1
	var col int = -1
	for cnt := 0; cnt < 9; cnt++ {
		i := t.random.Intn(NUM_ROW)
		j := t.random.Intn(NUM_COL)
		if t.cell[i][j] == 0 {
			row = i
			col = j
			break
		}
	}
	return row, col
}

func (t *TicTacToe) GetHumanWinCell() (int, int) {
	var row int = -1
	var col int = -1

hwc:
	for r := 0; r < NUM_ROW; r++ {
		for c := 0; c < NUM_COL; c++ {
			if t.cell[r][c] == 0 {
				if t.CanPlayerWin(r, c, HUMAN_PLAYER) {
					row = r
					col = c

					break hwc
				}
			}
		}
	}

	return row, col
}

func (t *TicTacToe) GetAIWinCell() (int, int) {
	var row int = -1
	var col int = -1

aiwc:
	for r := 0; r < NUM_ROW; r++ {
		for c := 0; c < NUM_COL; c++ {
			if t.cell[r][c] == 0 {
				if t.CanPlayerWin(r, c, AI_PLAYER) {
					row = r
					col = c

					break aiwc
				}
			}
		}
	}

	return row, col
}

func (t *TicTacToe) GetAINearCell() (int, int) {
	var row int = -1
	var col int = -1
ainc:
	for r := 0; r < NUM_ROW; r++ {
		for c := 0; c < NUM_COL; c++ {
			if t.cell[r][c] == AI_PLAYER {
				row, col = t.GetNearPlace(r, c)
				if row != -1 && col != -1 {
					break ainc
				}
			}
		}
	}

	return row, col
}

/*
func (t *TicTacToe) GetBelowAvgAIMove() (int, int) {
	var row int = -1
	var col int = -1
	var turn int = t.numAIMove + t.numHumanMove

	switch turn {
	case TURN1:
		row, col = t.GetRandEmptyCell()

	case TURN2, TURN3, TURN4:
		row, col = t.GetEmptyRowCol()
		if row == -1 || col == -1 {
			row, col = t.GetRandEmptyCell()
		}
	case TURN5, TURN6, TURN7, TURN8:
		row, col = t.GetRandEmptyCell()
		if row == -1 || col == -1 {
			row, col = t.GetAINearCell()
		}
	}

	return row, col
}

func (t *TicTacToe) GetAvgAIMove() (int, int) {

	var row int = -1
	var col int = -1
	var turn int = t.numAIMove + t.numHumanMove

	switch turn {
	case TURN1:
		row, col = t.GetRandEmptyCell()
	case TURN2:

		row, col = t.GetEmptyRowCol()

	case TURN3, TURN4:
		row, col = t.GetAINearCell()
	case TURN5:

		row, col = t.GetAIWinCell()

		if row == -1 || col == -1 {
			row, col = t.GetEmptyRowCol()
		}

	case TURN6, TURN7, TURN8, TURN9:
		row, col = t.GetAIWinCell()

		if row == -1 || col == -1 {
			row, col = t.GetAINearCell()
		}
	}

	return row, col
}

func (t *TicTacToe) GetGoodAIMove() (int, int) {

	var row int = -1
	var col int = -1
	var turn int = t.numAIMove + t.numHumanMove

	switch turn {
	case TURN1:
		row, col = t.GetRandEmptyCell()
	case TURN2:

		row, col = t.GetEmptyRowCol()

	case TURN3:
		row, col = t.GetAINearCell()
	case TURN4:
		row, col = t.GetHumanWinCell()

		if row == -1 || col == -1 {
			row, col = t.GetAINearCell()
		}

	case TURN5:
		row, col = t.GetAIWinCell()

		if row == -1 || col == -1 {
			row, col = t.GetHumanWinCell()
		}

		if row == -1 || col == -1 {
			row, col = t.GetEmptyRowCol()
		}

		if row == -1 || col == -1 {
			row, col = t.GetAINearCell()
		}

	case TURN6, TURN7, TURN8, TURN9:
		row, col = t.GetAIWinCell()

		if row == -1 || col == -1 {
			row, col = t.GetHumanWinCell()
		}

		if row == -1 || col == -1 {
			row, col = t.GetAINearCell()
		}

	}

	return row, col
}


func (t *TicTacToe) GetAIMove(aiType int) (int, int) {

	var row int = -1
	var col int = -1

	if aiType == AI_TYPE_BELOW_AVERAGE {
		row, col = t.GetBelowAvgAIMove()

	} else if aiType == AI_TYPE_AVERAGE {
		row, col = t.GetAvgAIMove()
	} else if aiType == AI_TYPE_GOOD {
		row, col = t.GetGoodAIMove()
	}

	if row == -1 || col == -1 {
		row, col = t.GetNextEmptyCell()
	}

	t.PrintCell()

	fmt.Printf("\n row = %d, col = %d", row, col)

	return row, col
}


*/

func (t *TicTacToe) PrintCell() {

	for i := 0; i < NUM_ROW; i++ {
		fmt.Println()
		for j := 0; j < NUM_COL; j++ {
			fmt.Printf("%d ", t.cell[j][i])
		}

	}

}

func (t *TicTacToe) GetAIMove(aiType int) (int, int) {

	var row int = -1
	var col int = -1

	// log.Printf("%d", aiType)

	var totalmove int = t.numAIMove + t.numHumanMove

	if t.numAIMove == 0 {
		row, col = t.GetRandEmptyCell()
	} else if t.numAIMove == 1 {
		if aiType == AI_TYPE_BELOW_AVERAGE {
			row, col = t.GetNextEmptyCell()
		} else if aiType == AI_TYPE_AVERAGE {
			row, col = t.GetAINearCell()
		} else if aiType == AI_TYPE_GOOD {
			if t.numHumanMove == 2 {
				// check if possible of human win
				row, col = t.GetHumanWinCell()
			}
			if row == -1 {
				row, col = t.GetEmptyRowCol()
			}

			if row == -1 {

				row, col = t.GetAINearCell()
			}
		}
	} else if t.numAIMove > 1 {

		if aiType == AI_TYPE_BELOW_AVERAGE {
			if row == -1 {
				row, col = t.GetAINearCell()
			}

			if row == -1 {
				row, col = t.GetRandEmptyCell()
			}
		} else if aiType == AI_TYPE_AVERAGE {
			row, col = t.GetAIWinCell()

			if row == -1 {
				if totalmove > 3 {
					row, col = t.GetHumanWinCell()
				} else {
					row, col = t.GetEmptyRowCol()
				}
			}
			if row == -1 {
				row, col = t.GetAINearCell()
			}

		} else if aiType == AI_TYPE_GOOD {
			row, col = t.GetAIWinCell()

			if row == -1 {
				row, col = t.GetHumanWinCell()
			}

			if row == -1 {
				row, col = t.GetEmptyRowCol()
			}

			if row == -1 {
				row, col = t.GetAINearCell()
			}
		}
	}

	if row == -1 {
		row, col = t.GetNextEmptyCell()
	}

	return row, col

}

/*
func (t *TicTacToe) GetAIMove(aiType int) (int, int) {

	var row int = -1
	var col int = -1

	if aiType == AI_TYPE_AVERAGE || aiType == AI_TYPE_GOOD {

		if t.numAIMove == 0 {

			for cnt := 0; cnt < 9; cnt++ {
				i := t.random.Intn(NUM_ROW)
				j := t.random.Intn(NUM_COL)
				if t.cell[i][j] == 0 {
					row = i
					col = j
					break
				}
			}
		} else if t.numAIMove == 1 {
			if t.numHumanMove == 2 {
				// check if possible of human win
			out1:
				for i := 0; i < NUM_ROW; i++ {
					for j := 0; j < NUM_COL; j++ {
						if t.cell[i][j] == 0 {
							if t.CanPlayerWin(i, j, HUMAN_PLAYER) {
								row = i
								col = j
								break out1
							}
						}
					}
				}
			}
			if row == -1 {

			out2:
				for i := 0; i < NUM_ROW; i++ {
					for j := 0; j < NUM_COL; j++ {
						if t.cell[i][j] == AI_PLAYER {
							r, c := t.GetNearPlace(i, j)
							if r != -1 && c != -1 {
								row = r
								col = c
								break out2
							}
						}
					}
				}
			}

			if row == -1 {
				for cnt := 0; cnt < 9; cnt++ {
					i := t.random.Intn(NUM_ROW)
					j := t.random.Intn(NUM_COL)
					if t.cell[i][j] == 0 {
						row = i
						col = j
						break
					}
				}
			}
		} else if t.numAIMove > 1 {
		out3:
			for i := 0; i < NUM_ROW; i++ {
				for j := 0; j < NUM_COL; j++ {
					if t.cell[i][j] == 0 {
						if t.CanPlayerWin(i, j, AI_PLAYER) {
							row = i
							col = j
							break out3
						}
					}
				}
			}

			if row == -1 {
			out4:
				for i := 0; i < NUM_ROW; i++ {
					for j := 0; j < NUM_COL; j++ {
						if t.cell[i][j] == 0 {
							if t.CanPlayerWin(i, j, HUMAN_PLAYER) {
								row = i
								col = j
								break out4
							}
						}
					}
				}
			}

			if row == -1 {
			out5:
				for i := 0; i < NUM_ROW; i++ {
					for j := 0; j < NUM_COL; j++ {
						if t.cell[i][j] == AI_PLAYER {
							r, c := t.GetNearPlace(i, j)
							if r != -1 && c != -1 {
								row = r
								col = c
								break out5
							}
						}
					}
				}
			}

			if row == -1 {

				for cnt := 0; cnt < 9; cnt++ {
					i := t.random.Intn(NUM_ROW)
					j := t.random.Intn(NUM_COL)
					if t.cell[i][j] == 0 {
						row = i
						col = j
						break
					}
				}
			}
		}

		if row == -1 {
		bk1:
			for i := 0; i < NUM_ROW; i++ {
				for j := 0; j < NUM_COL; j++ {
					if t.cell[i][j] == 0 {
						row = i
						col = j
						break bk1
					}
				}
			}
		}

	} else if aiType == AI_TYPE_BELOW_AVERAGE {
	bk:
		for i := 0; i < NUM_ROW; i++ {
			for j := 0; j < NUM_COL; j++ {
				if t.cell[i][j] == 0 {
					row = i
					col = j
					break bk
				}
			}
		}
	}

	t.PrintCell()

	fmt.Printf("\n row = %d, col = %d", row, col)

	return row, col

}
*/
