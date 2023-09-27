package tictactoe

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image       *ebiten.Image
	scalefactor float64
}

func (s *Sprite) Draw(screen *ebiten.Image, x int, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Scale(s.scalefactor, s.scalefactor)
	screen.DrawImage(s.image, op)
}

func NewSprite(img *ebiten.Image, x int, y int, w int, h int, scalefactor float64) *Sprite {
	var sprite *Sprite = new(Sprite)
	bound := image.Rect(x, y, x+w, y+h)
	sprite.image = (img.SubImage(bound)).(*ebiten.Image)
	sprite.scalefactor = scalefactor
	return sprite
}
