package tictactoe

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"golang.org/x/image/font/opentype"
)

const (
	sampleRate = 48000
)

type ResourceManager struct {
	images       map[string]*ebiten.Image
	fonts        map[string]*opentype.Font
	audios       map[string]*audio.Player
	audioContext *audio.Context
}

func (r *ResourceManager) Init() {
	r.images = make(map[string]*ebiten.Image)
	r.fonts = make(map[string]*opentype.Font)
	r.audios = make(map[string]*audio.Player)
	r.audioContext = audio.NewContext(sampleRate)
}

func (r *ResourceManager) LoadMp3Audio(name string) (*audio.Player, error) {
	var err error
	var data []byte
	var audioPlayer *audio.Player

	data, err = LoadAudio(name)
	if err != nil {
		return nil, err
	}
	d, err := mp3.DecodeWithoutResampling(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	audioPlayer, err = r.audioContext.NewPlayer(d)
	if err != nil {
		return nil, err
	}
	r.audios[name] = audioPlayer
	return audioPlayer, nil
}

func (r *ResourceManager) GetAudio(name string) *audio.Player {
	sndPlayer, ok := r.audios[name]
	if !ok {
		return nil
	}
	return sndPlayer
}

func (r *ResourceManager) LoadImage(name string) (*ebiten.Image, error) {

	var image, err = LoadImage(name)
	if err != nil {
		return nil, err
	}

	r.images[name] = image

	return image, nil

}

func (r *ResourceManager) GetImage(name string) *ebiten.Image {
	image, ok := r.images[name]
	if !ok {
		return nil
	}
	return image
}

func (r *ResourceManager) LoadFont(name string) (*opentype.Font, error) {
	var font, err = LoadFont(name)
	if err != nil {
		return nil, err
	}
	r.fonts[name] = font
	return font, nil
}

func (r *ResourceManager) GetFont(name string) *opentype.Font {

	font, ok := r.fonts[name]
	if !ok {
		return nil
	}
	return font
}

func (r *ResourceManager) UnloadResources() {

	//unload images
	for key, val := range r.images {
		val.Dispose()
		delete(r.images, key)
	}

	for key, _ := range r.fonts {
		delete(r.fonts, key)
	}

	// unload audio
	for key, val := range r.audios {
		val.Close()
		delete(r.audios, key)
	}
}

func NewResourceManager() *ResourceManager {
	var rm = new(ResourceManager)
	rm.Init()
	return rm
}
