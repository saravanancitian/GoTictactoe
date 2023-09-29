package tictactoe 
import (
	"golang.org/x/image/font/opentype"
	"github.com/hajimehoshi/ebiten/v2"
)
type ResourceManager struct {
	images map[string]*ebiten.Image
	fonts map[string]*opentype.Font
}

func (r *ResourceManager) Init(){
	r.images = make(map[string]*ebiten.Image)
	r.fonts = make(map[string]*opentype.Font)
}

func (r *ResourceManager) LoadImage(name string) (*ebiten.Image, error){

	var image, err = LoadImage(name);
	if err != nil {
		return nil, err
	}

	r.images[name] = image
	
	return image, nil
	
}

func (r *ResourceManager) GetImage(name string) *ebiten.Image{
	image, ok := r.images[name]
	if !ok {
		return nil
	}
	return image
}

func (r *ResourceManager) LoadFont(name string) (*opentype.Font, error){
	var font, err = LoadFont(name)
	if err != nil {
		return nil, err
	}
	r.fonts[name] = font
	return font, nil
}

func (r *ResourceManager) GetFont(name string) *opentype.Font{

	font, ok := r.fonts[name]
	if !ok {
		return nil
	}
	return font	
}

func NewResourceManager()*ResourceManager{
	var rm = new(ResourceManager)
	rm.Init()
	return rm 
}