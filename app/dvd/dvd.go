package dvd

import (
	"fmt"

	ac "dvd/app/app_configs"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FILEPATH = "./assets/images/dvd_video_logo.png"
)

type Dvd struct {
	texture *sdl.Texture

	X, Y int32
	W, H int32
}

func NewDvd(r *sdl.Renderer) (*Dvd, error) {
	if err := img.Init(img.INIT_PNG); err != nil {
		return nil, fmt.Errorf("Could not initialise an image: %v", err)
	}

	sur, err := img.Load(FILEPATH)
	if err != nil {
		return nil, fmt.Errorf("Could not load an image: %v", err)
	}

	defer sur.Free()

	t, err := r.CreateTextureFromSurface(sur)
	if err != nil {
		return nil, fmt.Errorf("Could not create a texture from surface: %v", err)
	}

	dvd := &Dvd{
		texture: t,
		W:       156,
		H:       128,
		X:       (ac.SCREEN_WIDTH - 156) / 2,
		Y:       (ac.SCREEN_HEIGHT - 128) / 2,
	}

	return dvd, nil
}

func (d *Dvd) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: d.X, Y: d.Y, W: d.W, H: d.H}

	if err := r.Copy(d.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}

func (d *Dvd) Destroy() {
	d.texture.Destroy()
}

var directionX int32 = 1
var directionY int32 = 1

func (d *Dvd) Update() {
	if d.X < 0 || d.X+d.W >= ac.SCREEN_WIDTH {
		directionX = -directionX
	}

	d.X += directionX

	if d.Y < 0 || d.Y+d.H >= ac.SCREEN_HEIGHT {
		directionY = -directionY
	}

	d.Y += directionY
}
