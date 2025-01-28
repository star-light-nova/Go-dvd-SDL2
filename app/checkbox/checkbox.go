package checkbox

import (
	fc "dvd/app/font_config"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Checkbox struct {
	texture *sdl.Texture

	IsSelected bool

	X, Y int32
	W, H int32
}

func NewCheckbox(r *sdl.Renderer, label string) (*Checkbox, error) {
	f, err := ttf.OpenFont(fc.FONT_PATH, fc.FONT_SIZE)
	if err != nil {
		return nil, fmt.Errorf("Could not open font for Checkbox: %v", err)
	}

	defer f.Close()

	sur, err := f.RenderUTF8Solid(label, fc.FONT_COLOR)
	if err != nil {
		return nil, fmt.Errorf("Could not render text for Checkbox: %v", err)
	}

	defer sur.Free()

	texture, err := r.CreateTextureFromSurface(sur)
	if err != nil {
		return nil, fmt.Errorf("Could not create a texture for Checkbox: %v", err)
	}

	return &Checkbox{
		texture:    texture,
		IsSelected: false,
		X:          0,
		Y:          128,
		W:          128,
		H:          64,
	}, nil
}

func (c *Checkbox) IsHover(mouseEvent *sdl.MouseButtonEvent) bool {
	if mouseEvent.X >= c.X && mouseEvent.X <= c.X+c.W {
		if mouseEvent.Y >= c.Y && mouseEvent.Y <= c.Y+c.H {
			return true
		}
	}

	return false
}

func (c *Checkbox) Click(mouseEvent *sdl.MouseButtonEvent) {
	if mouseEvent.Button == sdl.BUTTON_LEFT && mouseEvent.State == sdl.RELEASED {
		c.IsSelected = !c.IsSelected
		log.Printf("Checkbox is selected %v", c.IsSelected)
	}
}

func (c *Checkbox) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: c.X, Y: c.Y, W: c.W, H: c.H}

	if err := r.Copy(c.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not render Checkbox: %v", err)
	}

	return nil
}

func (c *Checkbox) Destroy() {
	c.texture.Destroy()
}
