package dvd

import (
	"fmt"
	"sync"

	ac "dvd/app/app_configs"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	FILEPATH = "./assets/images/dvd_video_logo.png"
)

type Dvd struct {
	mu      sync.RWMutex
	texture *sdl.Texture

	ControlEvents chan *sdl.KeyboardEvent

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

	kevents := make(chan *sdl.KeyboardEvent)

	dvd := &Dvd{
		texture:       t,
		ControlEvents: kevents,
		W:             156,
		H:             128,
		X:             (ac.SCREEN_WIDTH - 156) / 2,
		Y:             (ac.SCREEN_HEIGHT - 128) / 2,
	}

	return dvd, nil
}

func (d *Dvd) Paint(r *sdl.Renderer) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	rect := &sdl.Rect{X: d.X, Y: d.Y, W: d.W, H: d.H}

	if err := r.Copy(d.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}

func (d *Dvd) Destroy() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.texture.Destroy()

	// Close the channel of events.
	close(d.ControlEvents)
}

var directionX int32 = 1
var directionY int32 = 1

func (d *Dvd) Update() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case kevent := <-d.ControlEvents:
		d.controlUpdate(kevent)
	default:
		d.boundaries()

		d.X += directionX
		d.Y += directionY
	}
}

func (d *Dvd) controlUpdate(kevent *sdl.KeyboardEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	switch kevent.Keysym.Sym {
	case sdl.K_UP:
		directionY = -abs(directionY)
	case sdl.K_DOWN:
		directionY = abs(directionY)
	case sdl.K_RIGHT:
		directionX = abs(directionX)
	case sdl.K_LEFT:
		directionX = -abs(directionX)
	}

	d.boundaries()

	d.X += directionX
	d.Y += directionY
}

func (d *Dvd) boundaries() {
	if d.X < 0 || d.X+d.W >= ac.SCREEN_WIDTH {
		directionX = -directionX
	}

	if d.Y < 0 || d.Y+d.H >= ac.SCREEN_HEIGHT {
		directionY = -directionY
	}

}

func abs(number int32) int32 {
	if number < 0 {
		return -number
	}

	return number
}
