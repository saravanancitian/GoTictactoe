package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"fmt"
)

type Input struct{
	touchmode bool
	touchids []ebiten.TouchID

}

var theInput *Input



func Current() *Input {
	if theInput == nil {
		theInput = new(Input)
		theInput.Init()
	}
	return theInput
}


func (i *Input) Init(){
	i.touchmode =  isTouchPrimaryInput() 

}

func (i *Input) GetPosition() (int, int){
	var x, y int  = -1, -1
	if i.touchmode {

		
		// i.touchids = inpututil.AppendJustPressedTouchIDs(i.touchids)
		// fmt.Print(" ------- ",len(i.touchids))

		for _, id := range ebiten.AppendTouchIDs(nil){
			x, y = ebiten.TouchPosition(id)
			fmt.Print("\n ------- ",x, y)

			// i.touchids = i.touchids[:0]
			break;
		}

	} else {
		if(inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)){
			x, y = ebiten.CursorPosition()
		}
	}
	return x, y
}