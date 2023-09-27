package tictactoe

import (
	"image/color"
	"math/rand"
	"time"
	"tictactoe/tictactoe/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CELL_SIZE int = 62

	GAMEOVER_WIDTH int = 159

	GAMEOVER_HEIGHT int = 93

	NUM_COL int = 3
	NUM_ROW int = 3

	MARGIN       int = 4
	HUMAN_PLAYER int = 1
	AI_PLAYER    int = -1

	ANIMATE_HUMAN_PLAYER int = 2
	ANIMATE_AI_PLAYER    int = -2

	MAX_MOVES_CHECK_GAME_OVER int   = 5
	ANIM_DELAY                int64 = 100

	AI_TYPE_BELOW_AVERAGE int = 0
	AI_TYPE_AVERAGE       int = 1
	AI_TYPE_GOOD          int = 2
)

const (
	STATE_INIT_NEW_GAME        int = 1
	STATE_AI_PLAYER_TURN       int = 2
	STATE_HUMAN_PLAYER_TURN    int = 3
	STATE_GAME_OVER            int = 4
	STATE_ANIMATE_AI_PLAYER    int = 5
	STATE_ANIMATE_HUMAN_PLAYER int = 6
)

type Board struct {
	x        int
	y        int
	width    int
	height   int
	cellsize int
	margin   int

	cirCrossOverly *ebiten.Image
	youwin         *ebiten.Image
	loser          *ebiten.Image
	tide           *ebiten.Image

	bgColor     color.RGBA
	marginColor color.RGBA

	cellX [NUM_ROW]int
	cellY [NUM_COL]int

	cell [NUM_ROW][NUM_COL]int

	circleSprites []*Sprite
	crossSprites  []*Sprite
	overlayColor  color.RGBA

	circleImgIdx int
	crossImgIdx  int

	state int

	numHumanMove int
	numAIMove    int
	winner       int

	animtime    int64
	isAnimating bool
	isGameover  bool

	aiType int

	animRow int
	animCol int

	scalefactor float64
}

func (b *Board) ResetBoard() {
	for i := 0; i < NUM_ROW; i++ {
		for j := 0; j < NUM_COL; j++ {
			b.cell[i][j] = 0
		}
	}
}

func (b *Board) init(scalefactor float64, screenWidth int, screenHeight int) {

	if scalefactor < 1 {
		scalefactor = 1
	}
	b.scalefactor = scalefactor
	b.cellsize = CELL_SIZE
	b.margin = MARGIN

	b.width = b.cellsize*3 + b.margin*4
	b.height = b.width
	var cellwithmargine = b.margin + b.cellsize

	b.x = (screenWidth - b.width) / 2
	b.y = (screenHeight - b.height) / 2

	var x int = b.x

	for i := 0; i < NUM_COL; i++ {
		b.cellX[i] = x
		x += cellwithmargine
	}

	var y int = b.y
	for i := 0; i < NUM_ROW; i++ {
		b.cellY[i] = y
		y += cellwithmargine
	}

	b.bgColor = color.RGBA{255, 245, 184, 0xff}
	b.marginColor = color.RGBA{190, 145, 51, 0xff}
	b.overlayColor = color.RGBA{50, 50, 50, 150}

	b.loadSprite()
	b.StartNewGame()

	rand.Seed(time.Now().UnixNano())

}

func (b *Board) StartNewGame() {
	b.state = STATE_INIT_NEW_GAME
}

func (b *Board) loadSprite() {
	var err error
	b.cirCrossOverly, err = LoadImage("circlecross.png")
	if err != nil {
		panic(err)
	}

	b.youwin, err = LoadImage("youwin.png")
	if err != nil {
		panic(err)
	}

	b.loser, err = LoadImage("loser.png")
	if err != nil {
		panic(err)
	}

	b.tide, err = LoadImage("tide.png")
	if err != nil {
		panic(err)
	}

	b.circleSprites = make([]*Sprite, NUM_CIRCLE_FRAMES)
	var x, y int = 0, 0
	for i := 0; i < NUM_CIRCLE_FRAMES; i++ {
		b.circleSprites[i] = NewSprite(b.cirCrossOverly, x, y, CELL_SIZE, CELL_SIZE, b.scalefactor)
		x = x + CELL_SIZE
	}

	b.crossSprites = make([]*Sprite, NUM_CROSS_FRAMES)
	x = 0
	y = CELL_SIZE
	for i := 0; i < NUM_CROSS_FRAMES; i++ {
		b.crossSprites[i] = NewSprite(b.cirCrossOverly, x, y, CELL_SIZE, CELL_SIZE, b.scalefactor)
		x = x + CELL_SIZE
	}

}

func (b *Board) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), b.bgColor, false)
	vector.StrokeRect(screen, float32(b.x), float32(b.x), float32(b.width), float32(b.height), float32(b.margin), b.marginColor, false)

	vector.StrokeLine(screen, float32(b.x+b.cellsize+b.margin), float32(b.y), float32(b.x+b.cellsize+b.margin), float32(b.y+b.height), float32(b.margin), b.marginColor, false)
	vector.StrokeLine(screen, float32(b.x+(b.cellsize+b.margin)*2), float32(b.y), float32(b.x+(b.cellsize+b.margin)*2), float32(b.y+b.height), float32(b.margin), b.marginColor, false)

	vector.StrokeLine(screen, float32(b.x), float32(b.y+b.cellsize+b.margin), float32(b.x+b.width), float32(b.y+b.cellsize+b.margin), float32(b.margin), b.marginColor, false)

	vector.StrokeLine(screen, float32(b.x), float32(b.y+(b.cellsize+b.margin)*2), float32(b.x+b.width), float32(b.y+(b.cellsize+b.margin)*2), float32(b.margin), b.marginColor, false)

	b.DrawPlayer(screen)

	if b.isGameover {
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), float32(b.width), float32(b.height), b.overlayColor, false)
		var goimg *ebiten.Image
		if b.winner == HUMAN_PLAYER {
			goimg = b.youwin
		} else if b.winner == AI_PLAYER {
			goimg = b.loser
		} else {
			goimg = b.tide
		}

		gx := b.x + (b.width-GAMEOVER_WIDTH)/2
		gy := b.y + (b.height-GAMEOVER_HEIGHT)/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(gx), float64(gy))
		screen.DrawImage(goimg, op)

	}

}

func (b *Board) DrawPlayer(screen *ebiten.Image) {
	var x int = 0
	var y int = 0

	for i := 0; i < NUM_ROW; i++ {
		y = b.cellY[i] + b.margin

		for j := 0; j < NUM_COL; j++ {
			x = b.cellY[j] + b.margin

			player := b.cell[i][j]

			if player == HUMAN_PLAYER {
				b.crossSprites[NUM_CROSS_FRAMES-1].Draw(screen, x, y)
			} else if player == AI_PLAYER {

				b.circleSprites[NUM_CIRCLE_FRAMES-1].Draw(screen, x, y)
			} else if player == ANIMATE_HUMAN_PLAYER {
				b.crossSprites[b.crossImgIdx].Draw(screen, x, y)

			} else if player == ANIMATE_AI_PLAYER {
				b.circleSprites[b.circleImgIdx].Draw(screen, x, y)
			}
		}
	}
}

func (b *Board) Update(delta int64) {
	switch b.state {
	case STATE_INIT_NEW_GAME:
		b.ResetBoard()
		b.numAIMove = 0
		b.numHumanMove = 0
		b.crossImgIdx = 0
		b.circleImgIdx = 0
		b.state = STATE_AI_PLAYER_TURN
		b.aiType = AI_TYPE_AVERAGE

	case STATE_AI_PLAYER_TURN:
		row, col := b.GetAIMove(b.aiType)
		if row > -1 && col > -1 {
			b.numAIMove++

			b.cell[row][col] = ANIMATE_AI_PLAYER
			b.animRow = row
			b.animCol = col
			b.circleImgIdx = 0
			b.state = STATE_ANIMATE_AI_PLAYER
		}

	case STATE_HUMAN_PLAYER_TURN:
		mx, my := input.Current().GetPosition();
		if mx >=0 && my >= 0 {
			row, col := b.GetSelectedCell(mx, my)
			if row > -1 && col > -1 {
				b.numHumanMove++
				b.cell[row][col] = ANIMATE_HUMAN_PLAYER
				b.animRow = row
				b.animCol = col
				b.crossImgIdx = 0
				b.state = STATE_ANIMATE_HUMAN_PLAYER
			}
		}

	case STATE_GAME_OVER:
		b.isGameover = true
	case STATE_ANIMATE_AI_PLAYER:
		if b.DelayElapsed(delta) {
			if b.circleImgIdx >= NUM_CIRCLE_FRAMES-1 {
				b.isAnimating = false
				b.cell[b.animRow][b.animCol] = AI_PLAYER
				isGameOver, winner := b.CheckGameOver()
				if isGameOver {
					b.winner = winner

					b.state = STATE_GAME_OVER
				} else {
					b.state = STATE_HUMAN_PLAYER_TURN
				}

			} else {
				b.circleImgIdx++
				b.animtime = ANIM_DELAY
			}
		}

	case STATE_ANIMATE_HUMAN_PLAYER:

		if b.DelayElapsed(delta) {
			if b.crossImgIdx >= NUM_CROSS_FRAMES-1 {
				b.isAnimating = false
				b.cell[b.animRow][b.animCol] = HUMAN_PLAYER
				isGameOver, winner := b.CheckGameOver()
				if isGameOver {
					b.winner = winner

					b.state = STATE_GAME_OVER
				} else {
					b.state = STATE_AI_PLAYER_TURN
				}

			} else {
				b.crossImgIdx++
				b.animtime = ANIM_DELAY
			}
		}
	}
}

func (b *Board) SetAnimation() {
	b.animtime = ANIM_DELAY
	b.isAnimating = true
}

func (b *Board) DelayElapsed(delta int64) bool {
	if b.animtime-delta > 0 {
		b.animtime -= delta
	} else {
		return true
	}
	return false

}

func (b *Board) canPlayerWin(row int, col int, playerType int) bool {
	var isWin bool = false

	if row == 0 && col == 0 {
		//check col, row, dia
		isWin = (b.cell[row+1][col] == playerType && b.cell[row+2][col] == playerType || b.cell[row][col+1] == playerType && b.cell[row][col+2] == playerType || b.cell[row+1][col+1] == playerType && b.cell[row+2][col+2] == playerType)
	} else if row == 0 && col == 2 {
		isWin = (b.cell[row+1][col] == playerType && b.cell[row+2][col] == playerType || b.cell[row][col-1] == playerType && b.cell[row][col-2] == playerType || b.cell[row+1][col-1] == playerType && b.cell[row+2][col-2] == playerType)
	} else if row == 2 && col == 0 {
		isWin = (b.cell[row-1][col] == playerType && b.cell[row-2][col] == playerType || b.cell[row][col+1] == playerType && b.cell[row][col+2] == playerType || b.cell[row-1][col+1] == playerType && b.cell[row-2][col+2] == playerType)
	} else if row == 2 && col == 2 {
		isWin = (b.cell[row-1][col] == playerType && b.cell[row-2][col] == playerType || b.cell[row][col-1] == playerType && b.cell[row][col-2] == playerType || b.cell[row-1][col-1] == playerType && b.cell[row-2][col-2] == playerType)
	} else if row == 1 && col == 1 {
		isWin = (b.cell[row-1][col] == playerType && b.cell[row+1][col] == playerType || b.cell[row][col-1] == playerType && b.cell[row][col+1] == playerType || b.cell[row-1][col-1] == playerType && b.cell[row+1][col+1] == playerType)
	} else if row == 0 && col == 1 {
		isWin = (b.cell[row][col-1] == playerType && b.cell[row][col+1] == playerType || b.cell[row+1][col] == playerType && b.cell[row+2][col] == playerType)
	} else if row == 1 && col == 0 {
		isWin = (b.cell[row][col+1] == playerType && b.cell[row][col+2] == playerType || b.cell[row-1][col] == playerType && b.cell[row+1][col] == playerType)
	} else if row == 2 && col == 1 {
		isWin = (b.cell[row][col-1] == playerType && b.cell[row][col+1] == playerType || b.cell[row-1][col] == playerType && b.cell[row-2][col] == playerType)
	} else if row == 1 && col == 2 {
		isWin = (b.cell[row][col-1] == playerType && b.cell[row][col-2] == playerType || b.cell[row-1][col] == playerType && b.cell[row+1][col] == playerType)

	}

	return isWin
}

func (b *Board) CheckGameOver() (bool, int) {
	var isGameOver bool = false
	var winner int = 0
	numMove := b.numHumanMove + b.numAIMove
	if numMove >= MAX_MOVES_CHECK_GAME_OVER {
		if (b.cell[0][0] != 0) && (b.cell[0][0] == b.cell[0][1] && b.cell[0][1] == b.cell[0][2] || b.cell[0][0] == b.cell[1][0] && b.cell[1][0] == b.cell[2][0] || b.cell[0][0] == b.cell[1][1] && b.cell[1][1] == b.cell[2][2]) {
			isGameOver = true
			winner = b.cell[0][0]
		} else if (b.cell[1][0] != 0) && (b.cell[1][0] == b.cell[1][1] && b.cell[1][1] == b.cell[1][2]) {
			isGameOver = true
			winner = b.cell[1][0]
		} else if (b.cell[2][0] != 0) && (b.cell[2][0] == b.cell[2][1] && b.cell[2][1] == b.cell[2][2] || b.cell[2][0] == b.cell[1][1] && b.cell[1][1] == b.cell[0][2]) {
			isGameOver = true
			winner = b.cell[2][0]
		} else if (b.cell[0][1] != 0) && (b.cell[0][1] == b.cell[1][1] && b.cell[1][1] == b.cell[2][1]) {
			isGameOver = true
			winner = b.cell[0][1]
		} else if (b.cell[0][2] != 0) && (b.cell[0][2] == b.cell[1][2] && b.cell[1][2] == b.cell[2][2]) {
			isGameOver = true
			winner = b.cell[0][2]
		}

		if !isGameOver && numMove == ((NUM_COL*NUM_ROW)-1) {
			isGameOver = true
		}

	}
	return isGameOver, winner
}

func (b *Board) GetNextEmptyCell() (int, int) {
	for row := 0; row < NUM_ROW; row++ {
		for col := 0; col < NUM_COL; col++ {
			if b.cell[row][col] == 0 {
				return row, col
			}
		}
	}
	return -1, -1
}

func (b *Board) GetSelectedCell(x, y int) (int, int) {

	for row := 0; row < NUM_ROW; row++ {
		for col := 0; col < NUM_COL; col++ {
			if x > b.cellX[col] && x < b.cellX[col]+b.cellsize && y > b.cellY[row] && y < b.cellY[row]+b.cellsize {
				return row, col
			}
		}
	}
	return -1, -1
}

func NewBoard(scalefactor float64, screenwidth, screenheight int) *Board {
	var board *Board = new(Board)
	board.init(scalefactor, screenwidth, screenheight)
	return board
}

func (b *Board) getSecondPlace(currow int, curcol int, playerType int) (int, int) {
	var row int = -1
	var col int = -1

	if currow == 0 && curcol == 0 {
		if b.cell[currow][curcol+1] == 0 && b.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if b.cell[currow+1][curcol] == 0 && b.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		} else if b.cell[currow+1][curcol+1] == 0 && b.cell[currow+2][curcol+2] == 0 {
			row = currow + 1
			col = curcol + 1
		}
	} else if currow == 0 && curcol == 2 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow+1][curcol] == 0 && b.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		} else if b.cell[currow+1][curcol-1] == 0 && b.cell[currow+2][curcol-2] == 0 {
			row = currow + 1
			col = curcol - 1
		}
	} else if currow == 2 && curcol == 0 {
		if b.cell[currow][curcol+1] == 0 && b.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if b.cell[currow-1][curcol+1] == 0 && b.cell[currow-2][curcol+2] == 0 {
			row = currow - 1
			col = curcol + 1
		}
	} else if currow == 2 && curcol == 2 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if b.cell[currow-1][curcol-1] == 0 && b.cell[currow-2][curcol-2] == 0 {
			row = currow - 1
			col = curcol - 1
		}
	} else if currow == 1 && curcol == 1 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow+1][curcol] == 0 {
			row = currow - 1
			col = curcol
		} else if b.cell[currow-1][curcol-1] == 0 && b.cell[currow+1][curcol+1] == 0 {
			row = currow - 1
			col = curcol - 1
		}
	} else if currow == 0 && curcol == 1 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow+1][curcol] == 0 && b.cell[currow+2][curcol] == 0 {
			row = currow + 1
			col = curcol
		}

	} else if currow == 1 && curcol == 0 {
		if b.cell[currow][curcol+1] == 0 && b.cell[currow][curcol+2] == 0 {
			row = currow
			col = curcol + 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow-1][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	} else if currow == 2 && curcol == 1 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol+1] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow-2][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	} else if currow == 1 && curcol == 2 {
		if b.cell[currow][curcol-1] == 0 && b.cell[currow][curcol-2] == 0 {
			row = currow
			col = curcol - 1
		} else if b.cell[currow-1][curcol] == 0 && b.cell[currow+1][curcol] == 0 {
			row = currow - 1
			col = curcol
		}
	}
	return row, col
}

func (b *Board) GetAIMove(aiType int) (int, int) {

	var row int = -1
	var col int = -1

	if aiType == AI_TYPE_AVERAGE {

		if b.numAIMove == 0 {

			for cnt := 0; cnt < 9; cnt++ {
				i := rand.Intn(NUM_ROW)
				j := rand.Intn(NUM_COL)
				if b.cell[i][j] == 0 {
					row = i
					col = j
					break
				}
			}
		} else if b.numAIMove == 1 {
			if b.numHumanMove == 2 {
				// check if possible of human win
			out1:
				for i := 0; i < NUM_ROW; i++ {
					for j := 0; j < NUM_COL; j++ {
						if b.cell[i][j] == 0 {
							if b.canPlayerWin(i, j, HUMAN_PLAYER) {
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
						if b.cell[i][j] == AI_PLAYER {
							r, c := b.getSecondPlace(i, j, AI_PLAYER)
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
					if b.cell[i][j] == 0 {
						row = i
						col = j
						break
					}
				}
			}
		} else if b.numAIMove > 1 {
		out3:
			for i := 0; i < NUM_ROW; i++ {
				for j := 0; j < NUM_COL; j++ {
					if b.cell[i][j] == 0 {
						if b.canPlayerWin(i, j, AI_PLAYER) {
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
						if b.cell[i][j] == 0 {
							if b.canPlayerWin(i, j, HUMAN_PLAYER) {
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
						if b.cell[i][j] == AI_PLAYER {
							r, c := b.getSecondPlace(i, j, AI_PLAYER)
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
					if b.cell[i][j] == 0 {
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
					if b.cell[i][j] == 0 {
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
				if b.cell[i][j] == 0 {
					row = i
					col = j
					break bk
				}
			}
		}
	}

	return row, col

}
