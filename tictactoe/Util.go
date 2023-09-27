package tictactoe

import (
	"image/png"
	"path"
	"tictactoe/tictactoe/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(filename string) (*ebiten.Image, error) {
	const dir = "images"
	f, err := assets.Assets.Open(path.Join(dir, filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// img, _, err := image.Decode(f)
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	fileimage := ebiten.NewImageFromImage(img)

	return fileimage, nil
}
