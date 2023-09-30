package tictactoe

import (
	"image/color"
	"math/rand"
	"time"
	"tictactoe/tictactoe/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"github.com/hajimehoshi/ebiten/v2/text"
	"fmt"

)

const (
	STATE_INIT_NEW_GAME        int = 1
	STATE_AI_PLAYER_TURN       int = 2
	STATE_HUMAN_PLAYER_TURN    int = 3
	STATE_GAME_OVER            int = 4
	STATE_ANIMATE_AI_PLAYER    int = 5
	STATE_ANIMATE_HUMAN_PLAYER int = 6
)


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
	NUM_COL int = 3
	NUM_ROW int = 3
	MAX_MOVES_CHECK_GAME_OVER int = 5

	TXT_YOUR_TURN string = "Your Turn"
)


type TicTacToe struct{
	GameOverCallBack func()

	youwin         *ebiten.Image
	loser          *ebiten.Image
	tide           *ebiten.Image
	gx int
	gy int

	txtTimer string
	showTimer bool
	showTurnText bool

	gameScreenWidth int
	gameScreenHeight int



	
	cell [NUM_ROW][NUM_COL]int
	
	JustAnotherHand_ttf *opentype.Font
	normalFont font.Face


	state int

	numHumanMove int
	numAIMove    int

	winner       int

	playtime int64
	playtimecalc int64

	animtime    int64
	isAnimating bool
	isGameover  bool


	animRow int
	animCol int

	board *Board
	player *Player

	overlayColor  color.RGBA
	rm *ResourceManager
}

func (t *TicTacToe) SetCallback(callback func()){
	t.GameOverCallBack = callback
}

func (t *TicTacToe) LoadSprite(rm *ResourceManager) {
	var err error
	t.youwin, err = rm.LoadImage("youwin.png")
	if err != nil {
		panic(err)
	}

	t.loser, err = rm.LoadImage("loser.png")
	if err != nil {
		panic(err)
	}

	t.tide, err = rm.LoadImage("tide.png")
	if err != nil {
		panic(err)
	}

	t.JustAnotherHand_ttf, err = rm.LoadFont("JustAnotherHand.ttf")
	if err != nil {
		panic(err)
	}

	t.normalFont, err = opentype.NewFace(t.JustAnotherHand_ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}



}


func NewTicTacToe(rm *ResourceManager, screenWidth int, screenHeight int, callback func()) *TicTacToe{
	var tictactoe = new(TicTacToe)
	tictactoe.Init(rm, screenWidth, screenHeight, callback)
	return tictactoe
}


func (t *TicTacToe) Init(rm *ResourceManager, screenWidth int, screenHeight int, callback func()){
	rand.Seed(time.Now().UnixNano())
	t.rm = rm
	t.GameOverCallBack = callback
	t.LoadSprite(rm)

	t.gx = 0
	t.gy = 0	
	t.gameScreenWidth = screenWidth
	t.gameScreenHeight = screenHeight
	
	t.board = NewBoard()
	t.player = NewPlayer()
	t.player.LoadSprite(rm)

	t.txtTimer = "00:00"

	
	x := (screenWidth - t.board.width)/2 
	y := (screenHeight - t.board.height)/2
	t.board.SetXY(x, y)


	
	t.overlayColor = color.RGBA{50, 50, 50, 150}


	t.StartNewGame()
}

func (t *TicTacToe) StartNewGame() {
	t.state = STATE_INIT_NEW_GAME
}


func (t *TicTacToe) ResetBoard() {
	for i := 0; i < NUM_ROW; i++ {
		for j := 0; j < NUM_COL; j++ {
			t.cell[i][j] = 0
		}
	}
}

func (t *TicTacToe) CalculatePlayTime(delta int64){
	if !t.isGameover { 
		if t.playtimecalc >= SEC_IN_MILLIS{
			t.playtime += t.playtimecalc
			t.playtimecalc = delta

			tsec := t.playtime  / SEC_IN_MILLIS
			min := tsec / 60
			sec := tsec % 60
			t.txtTimer = fmt.Sprintf("\n%02d:%02d", min, sec)

		} else {
			t.playtimecalc += delta
		}
	}
}


func (t *TicTacToe) CanPlayerWin(row int, col int, playerType int) bool {
	var isWin bool = false

	if row == 0 && col == 0 {
		//check col, row, dia
		isWin = (t.cell[row+1][col] == playerType && t.cell[row+2][col] == playerType || t.cell[row][col+1] == playerType && t.cell[row][col+2] == playerType || t.cell[row+1][col+1] == playerType && t.cell[row+2][col+2] == playerType)
	} else if row == 0 && col == 2 {
		isWin = (t.cell[row+1][col] == playerType && t.cell[row+2][col] == playerType || t.cell[row][col-1] == playerType && t.cell[row][col-2] == playerType || t.cell[row+1][col-1] == playerType && t.cell[row+2][col-2] == playerType)
	} else if row == 2 && col == 0 {
		isWin = (t.cell[row-1][col] == playerType && t.cell[row-2][col] == playerType || t.cell[row][col+1] == playerType && t.cell[row][col+2] == playerType || t.cell[row-1][col+1] == playerType && t.cell[row-2][col+2] == playerType)
	} else if row == 2 && col == 2 {
		isWin = (t.cell[row-1][col] == playerType && t.cell[row-2][col] == playerType || t.cell[row][col-1] == playerType && t.cell[row][col-2] == playerType || t.cell[row-1][col-1] == playerType && t.cell[row-2][col-2] == playerType)
	} else if row == 1 && col == 1 {
		isWin = (t.cell[row-1][col] == playerType && t.cell[row+1][col] == playerType || t.cell[row][col-1] == playerType && t.cell[row][col+1] == playerType || t.cell[row-1][col-1] == playerType && t.cell[row+1][col+1] == playerType)
	} else if row == 0 && col == 1 {
		isWin = (t.cell[row][col-1] == playerType && t.cell[row][col+1] == playerType || t.cell[row+1][col] == playerType && t.cell[row+2][col] == playerType)
	} else if row == 1 && col == 0 {
		isWin = (t.cell[row][col+1] == playerType && t.cell[row][col+2] == playerType || t.cell[row-1][col] == playerType && t.cell[row+1][col] == playerType)
	} else if row == 2 && col == 1 {
		isWin = (t.cell[row][col-1] == playerType && t.cell[row][col+1] == playerType || t.cell[row-1][col] == playerType && t.cell[row-2][col] == playerType)
	} else if row == 1 && col == 2 {
		isWin = (t.cell[row][col-1] == playerType && t.cell[row][col-2] == playerType || t.cell[row-1][col] == playerType && t.cell[row+1][col] == playerType)

	}

	return isWin
}

func (t *TicTacToe) CheckGameOver() (bool, int) {
	var isGameOver bool = false
	var winner int = 0
	numMove := t.numHumanMove + t.numAIMove
	if numMove >= MAX_MOVES_CHECK_GAME_OVER {
		if (t.cell[0][0] != 0) && (t.cell[0][0] == t.cell[0][1] && t.cell[0][1] == t.cell[0][2] || t.cell[0][0] == t.cell[1][0] && t.cell[1][0] == t.cell[2][0] || t.cell[0][0] == t.cell[1][1] && t.cell[1][1] == t.cell[2][2]) {
			isGameOver = true
			winner = t.cell[0][0]
		} else if (t.cell[1][0] != 0) && (t.cell[1][0] == t.cell[1][1] && t.cell[1][1] == t.cell[1][2]) {
			isGameOver = true
			winner = t.cell[1][0]
		} else if (t.cell[2][0] != 0) && (t.cell[2][0] == t.cell[2][1] && t.cell[2][1] == t.cell[2][2] || t.cell[2][0] == t.cell[1][1] && t.cell[1][1] == t.cell[0][2]) {
			isGameOver = true
			winner = t.cell[2][0]
		} else if (t.cell[0][1] != 0) && (t.cell[0][1] == t.cell[1][1] && t.cell[1][1] == t.cell[2][1]) {
			isGameOver = true
			winner = t.cell[0][1]
		} else if (t.cell[0][2] != 0) && (t.cell[0][2] == t.cell[1][2] && t.cell[1][2] == t.cell[2][2]) {
			isGameOver = true
			winner = t.cell[0][2]
		}

		if !isGameOver && numMove == ((NUM_COL*NUM_ROW)-1) {
			isGameOver = true
		}

	}
	return isGameOver, winner
}


func(t *TicTacToe) Draw(screen *ebiten.Image){

	t.board.Draw(screen)

	if t.showTurnText{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate( float64(t.board.x + t.board.width - 60), float64(t.board.y - 5))
		text.DrawWithOptions(screen, TXT_YOUR_TURN, t.normalFont,op)
	  }

 	if t.showTimer{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate( float64(t.board.x + 5), float64(t.board.y - 30))		
		 text.DrawWithOptions(screen, t.txtTimer, t.normalFont,  op)
	}

	

	var x int = 0
	var y int = 0

	for i := 0; i < NUM_ROW; i++ {
		x = t.board.x + (t.board.cellwithmargine* i) + t.board.margin

		for j := 0; j < NUM_COL; j++ {
			y = t.board.y + (t.board.cellwithmargine * j) + t.board.margin
			playerType := t.cell[i][j]

			t.player.Draw(screen, playerType, x, y)
		}
	}


	if t.isGameover {
		vector.DrawFilledRect(screen, float32(t.board.x), float32(t.board.y), float32(t.board.width), float32(t.board.height), t.overlayColor, false)
		var goimg *ebiten.Image
		if t.winner == HUMAN_PLAYER {
			goimg = t.youwin
		} else if t.winner == AI_PLAYER {
			goimg = t.loser
		} else {
			goimg = t.tide
		}

		gx := t.board.x + (t.board.width-GAMEOVER_WIDTH)/2
		gy := t.board.y + (t.board.height-GAMEOVER_HEIGHT)/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(gx), float64(gy))
		screen.DrawImage(goimg, op)

	}
}






func (t *TicTacToe) Update(delta int64){
	switch t.state {
	case STATE_INIT_NEW_GAME:
		t.ResetBoard()
		t.numAIMove = 0
		t.numHumanMove = 0
	
		t.state = STATE_AI_PLAYER_TURN
		t.player.aiType = AI_TYPE_AVERAGE
		t.playtime = 0
		t.playtimecalc = 0
		t.showTimer = true
		t.showTurnText = false

	case STATE_AI_PLAYER_TURN:
		t.CalculatePlayTime(delta)
		row, col := t.GetAIMove(t.player.aiType)
		if row > -1 && col > -1 {
			t.numAIMove++
			t.cell[row][col] = ANIMATE_AI_PLAYER
			t.animRow = row
			t.animCol = col
			t.player.circleImgIdx = 0
			t.state = STATE_ANIMATE_AI_PLAYER
		}

	case STATE_HUMAN_PLAYER_TURN:
		t.CalculatePlayTime(delta)
		mx, my := input.Current().GetPosition();
		if mx >=0 && my >= 0 {
			row, col := t.board.GetSelectedCell(mx, my)
			if row > -1 && col > -1 {
				t.numHumanMove++
				t.cell[row][col] = ANIMATE_HUMAN_PLAYER
				t.animRow = row
				t.animCol = col
				t.player.crossImgIdx = 0
				t.state = STATE_ANIMATE_HUMAN_PLAYER
				t.showTurnText = false
			}
		}

	case STATE_GAME_OVER:
		t.isGameover = true
		t.GameOverCallBack()
	case STATE_ANIMATE_AI_PLAYER:
		t.CalculatePlayTime(delta)
		if t.DelayElapsed(delta) {
			if t.player.circleImgIdx >= NUM_CIRCLE_FRAMES-1 {
				t.isAnimating = false
				t.cell[t.animRow][t.animCol] = AI_PLAYER
				isGameOver, winner := t.CheckGameOver()
				if isGameOver {
					t.winner = winner

					t.state = STATE_GAME_OVER
				} else {
					t.showTurnText = true

					t.state = STATE_HUMAN_PLAYER_TURN
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
				t.isAnimating = false
				t.cell[t.animRow][t.animCol] = HUMAN_PLAYER
				isGameOver, winner := t.CheckGameOver()
				if isGameOver {
					t.winner = winner

					t.state = STATE_GAME_OVER
				} else {
					t.state = STATE_AI_PLAYER_TURN
				}

			} else {
				t.player.crossImgIdx++
				t.animtime = ANIM_DELAY
			}
		}
	}
}


func (t *TicTacToe) SetAnimation() {
	t.animtime = ANIM_DELAY
	t.isAnimating = true
}

func (t *TicTacToe) DelayElapsed(delta int64) bool {
	if t.animtime-delta > 0 {
		t.animtime -= delta
	} else {
		return true
	}
	return false

}

func (t *TicTacToe) GetSecondPlace(currow int, curcol int, playerType int) (int, int) {
	var row int = -1
	var col int = -1

	if currow == 0 && curcol == 0 {
		if t.cell[currow][curcol+1] == 0 && t.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if t.cell[currow+1][curcol] == 0 && t.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		} else if t.cell[currow+1][curcol+1] == 0 && t.cell[currow+2][curcol+2] == 0 {
			row = currow + 1
			col = curcol + 1
		}
	} else if currow == 0 && curcol == 2 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow+1][curcol] == 0 && t.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		} else if t.cell[currow+1][curcol-1] == 0 && t.cell[currow+2][curcol-2] == 0 {
			row = currow + 1
			col = curcol - 1
		}
	} else if currow == 2 && curcol == 0 {
		if t.cell[currow][curcol+1] == 0 && t.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if t.cell[currow-1][curcol+1] == 0 && t.cell[currow-2][curcol+2] == 0 {
			row = currow - 1
			col = curcol + 1
		}
	} else if currow == 2 && curcol == 2 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if t.cell[currow-1][curcol-1] == 0 && t.cell[currow-2][curcol-2] == 0 {
			row = currow - 1
			col = curcol - 1
		}
	} else if currow == 1 && curcol == 1 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow+1][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if t.cell[currow-1][curcol-1] == 0 && t.cell[currow+1][curcol+1] == 0 {
			row = currow - 1
			col = curcol - 1
		}
	} else if currow == 0 && curcol == 1 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow+1][curcol] == 0 && t.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		}

	} else if currow == 1 && curcol == 0 {
		if t.cell[currow][curcol+1] == 0 && t.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow-1][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	} else if currow == 2 && curcol == 1 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	} else if currow == 1 && curcol == 2 {
		if t.cell[currow][curcol-1] == 0 && t.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if t.cell[currow-1][curcol] == 0 && t.cell[currow+1][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	}
	return row, col
}

func (t *TicTacToe) GetAIMove(aiType int) (int, int) {

	var row int = -1
	var col int = -1

	if aiType == AI_TYPE_AVERAGE {

		if t.numAIMove == 0 {

			for cnt := 0; cnt < 9; cnt++ {
				i := rand.Intn(NUM_ROW)
				j := rand.Intn(NUM_COL)
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
							r, c := t.GetSecondPlace(i, j, AI_PLAYER)
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
					i := rand.Intn(NUM_ROW)
					j := rand.Intn(NUM_COL)
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
							r, c := t.GetSecondPlace(i, j, AI_PLAYER)
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
					i := rand.Intn(NUM_ROW)
					j := rand.Intn(NUM_COL)
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

	return row, col

}
