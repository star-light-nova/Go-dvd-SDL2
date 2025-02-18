package button

import (
	"fmt"

	fc "dvd/app/font_config"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Button struct {
	texture *sdl.Texture

	X, Y int32
	H, W int32
}

func NewButton(r *sdl.Renderer, label string) (*Button, error) {
	f, err := ttf.OpenFont(fc.FONT_PATH, fc.FONT_SIZE)
	if err != nil {
		return nil, fmt.Errorf("Could not init Button Font: %v", err)
	}

	defer f.Close()

	surface, err := f.RenderUTF8Solid(label, fc.FONT_COLOR)
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
		X:       0,
		Y:       0,
		H:       50,
		W:       250,
	}, nil

}

func (b *Button) IsHover(mouseEvent *sdl.MouseButtonEvent) bool {
	if mouseEvent.X >= b.X && mouseEvent.X <= b.X+b.W {
		if mouseEvent.Y >= b.Y && mouseEvent.Y <= b.Y+b.H {
			return true
		}
	}

	return false
}

// Returns `true` if click has happened, otherwise `false`.
func (b *Button) Click(mouseEvent *sdl.MouseButtonEvent) bool {
	if mouseEvent.Button == sdl.BUTTON_LEFT && mouseEvent.State == sdl.RELEASED {
		return true
	}

	return false
}

func (b *Button) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H}

	if err := r.Copy(b.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not pain Button: %v", err)
	}

	return nil
}

func (b *Button) Texture() *sdl.Texture {
	return b.texture
}

func (b *Button) Destroy() {
	b.texture.Destroy()
}
