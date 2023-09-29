package tictactoe

import (
	"github.com/hajimehoshi/ebiten/v2"
)


type Player  struct{
	cirCrossOverly *ebiten.Image
	circleSprites []*Sprite
	crossSprites  []*Sprite

	circleImgIdx int
	crossImgIdx  int
	aiType int
}

func NewPlayer() *Player{
	var player = new(Player)
	player.Init()
	return player
} 

func (p *Player) Init(){
	p.crossImgIdx = 0
	p.circleImgIdx = 0
	p.aiType = AI_TYPE_AVERAGE

}


func (p *Player) LoadSprite(rm *ResourceManager) {
	var err error
	p.cirCrossOverly, err = rm.LoadImage("circlecross.png")
	if err != nil {
		panic(err)
	}


	p.circleSprites = make([]*Sprite, NUM_CIRCLE_FRAMES)
	var x, y int = 0, 0
	for i := 0; i < NUM_CIRCLE_FRAMES; i++ {
		p.circleSprites[i] = NewSprite(p.cirCrossOverly, x, y, CELL_SIZE, CELL_SIZE, 1)
		x = x + CELL_SIZE
	}

	p.crossSprites = make([]*Sprite, NUM_CROSS_FRAMES)
	x = 0
	y = CELL_SIZE
	for i := 0; i < NUM_CROSS_FRAMES; i++ {
		p.crossSprites[i] = NewSprite(p.cirCrossOverly, x, y, CELL_SIZE, CELL_SIZE, 1)
		x = x + CELL_SIZE
	}

}


func (p *Player) Draw(screen *ebiten.Image, playerType int, x int , y int){
	if playerType == HUMAN_PLAYER {
		p.crossSprites[NUM_CROSS_FRAMES-1].Draw(screen, x, y)
	} else if playerType == AI_PLAYER {

		p.circleSprites[NUM_CIRCLE_FRAMES-1].Draw(screen, x, y)
	} else if playerType == ANIMATE_HUMAN_PLAYER {
		p.crossSprites[p.crossImgIdx].Draw(screen, x, y)

	} else if playerType == ANIMATE_AI_PLAYER {
		p.circleSprites[p.circleImgIdx].Draw(screen, x, y)
	}

}


