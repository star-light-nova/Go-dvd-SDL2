package button

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var FONT_PATH = fmt.Sprintf("%s/assets/fonts/Roboto-Bold.ttf", sdl.GetBasePath())

const FONT_SIZE int = 27

type Button struct {
	texture *sdl.Texture

	H, W int32
	X, Y int32
}

func NewButton(r *sdl.Renderer, label string) (*Button, error) {
	if err := ttf.Init(); err != nil {
		return nil, fmt.Errorf("Could not init font (button): %v", err)
	}

	f, err := ttf.OpenFont(FONT_PATH, FONT_SIZE)
	if err != nil {
		return nil, fmt.Errorf("Could not init Button Font: %v", err)
	}

	defer f.Close()

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}

	surface, err := f.RenderUTF8Blended(label, c)
	if err != nil {
		return nil, fmt.Errorf("Could not init Button Surface: %v", err)
	}

	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("Could not init Button Texture: %v", err)
	}

	return &Button{
		texture: texture,
		H:       50,
		W:       50,
		X:       0,
		Y:       0,
	}, nil

}

func (b *Button) Click(mouseEvent *sdl.MouseButtonEvent) {
	// TODO: do something
	if mouseEvent.X >= b.X && mouseEvent.X <= b.X+b.W {
		if mouseEvent.Y >= b.Y && mouseEvent.Y <= b.Y+b.H {
			if mouseEvent.Button == sdl.BUTTON_LEFT && mouseEvent.State == sdl.RELEASED {
				f, err := sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT, Timestamp: sdl.GetTicks()})

				if err != nil {
					log.Printf("Could not send an quit event: %v %v", err, f)
				}

				log.Printf("Button has clicked")
			}
		}
	}
}

func (b *Button) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H}

	if err := r.Copy(b.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not pain Button: %v", err)
	}

	return nil
}

func (b *Button) Destroy() {
	b.texture.Destroy()
}
