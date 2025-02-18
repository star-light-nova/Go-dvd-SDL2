package checkbox

import (
	fc "dvd/app/font_config"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Checkbox struct {
	textures [2]*sdl.Texture

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

	surF, surT, err := createSurfaces(f, label)
	if err != nil {
		return nil, fmt.Errorf("Could not create a T/F surfaces: %v", err)
	}

	defer surF.Free()
	defer surT.Free()

	textureF, textureT, err := createTextures(r, surF, surT)
	if err != nil {
		return nil, fmt.Errorf("Could not create a T/F textures: %v", err)
	}

	return &Checkbox{
		textures:   [2]*sdl.Texture{textureF, textureT},
		IsSelected: false,
		X:          0,
		Y:          0,
		W:          250,
		H:          50,
	}, nil
}

// surfaces with True and False values
func createSurfaces(f *ttf.Font, label string) (surF, surT *sdl.Surface, err error) {
	labelFalse := fmt.Sprintf("%s: %t", label, false)

	surF, err = f.RenderUTF8Solid(labelFalse, fc.FONT_COLOR)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not render text for Checkbox: %v", err)
	}

	labelTrue := fmt.Sprintf("%s: %t", label, true)
	surT, err = f.RenderUTF8Solid(labelTrue, fc.FONT_COLOR)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not render text for Checkbox: %v", err)
	}

	return
}

// Texture with True and False values
func createTextures(r *sdl.Renderer, surF, surT *sdl.Surface) (textureF, textureT *sdl.Texture, err error) {
	textureF, err = r.CreateTextureFromSurface(surF)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not create a texture for Checkbox: %v", err)
	}

	textureT, err = r.CreateTextureFromSurface(surT)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not create a texture for Checkbox: %v", err)
	}

	return
}

func (c *Checkbox) IsHover(mouseEvent *sdl.MouseButtonEvent) bool {
	if mouseEvent.X >= c.X && mouseEvent.X <= c.X+c.W {
		if mouseEvent.Y >= c.Y && mouseEvent.Y <= c.Y+c.H {
			return true
		}
	}

	return false
}

// Returns `true` if click has happened, otherwise `false`.
func (c *Checkbox) Click(mouseEvent *sdl.MouseButtonEvent) bool {
	if mouseEvent.Button == sdl.BUTTON_LEFT && mouseEvent.State == sdl.RELEASED {
		c.IsSelected = !c.IsSelected

		return true
	}

	return false
}

func (c *Checkbox) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: c.X, Y: c.Y, W: c.W, H: c.H}

	if c.IsSelected {
		if err := r.Copy(c.textures[1], nil, rect); err != nil {
			return fmt.Errorf("Could not render Checkbox: %v", err)
		}
	} else {
		if err := r.Copy(c.textures[0], nil, rect); err != nil {
			return fmt.Errorf("Could not render Checkbox: %v", err)
		}
	}

	return nil
}

func (c *Checkbox) Texture() *sdl.Texture {
	if c.IsSelected {
		return c.textures[1]
	}

	return c.textures[0]
}

func (c *Checkbox) Destroy() {
	for _, texture := range c.textures {
		texture.Destroy()
	}
}
