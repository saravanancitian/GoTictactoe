package tictactoe

import (
	"image/png"
	"path"
	"tictactoe/tictactoe/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
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

func LoadFont(filename string) (*opentype.Font, error) {
	const dir = "fonts"
	data, err := assets.Fonts.ReadFile(path.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(data)

	if err != nil {
		return nil, err
	}

	return tt, nil
}

func LoadAudio(filename string) ([]byte, error) {
	const dir = "audio"

	data, err := assets.Audios.ReadFile(path.Join(dir, filename))
	if err != nil {
		return nil, err
	}

	return data, nil
}
