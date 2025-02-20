package dvd

import (
	"fmt"
	"sync"

	ac "dvd/app/app_configs"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Dvd struct {
	mu      sync.RWMutex
	texture *sdl.Texture

	ControlEvents     chan *sdl.KeyboardEvent
	MouseMotionEvents chan *sdl.MouseMotionEvent

	IsTargetX bool
	TargetX   int32

	IsTargetY bool
	TargetY   int32

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
	mevents := make(chan *sdl.MouseMotionEvent)

	dvd := &Dvd{
		texture:           t,
		ControlEvents:     kevents,
		MouseMotionEvents: mevents,
		X:                 (ac.SCREEN_WIDTH - x) / 2,
		Y:                 (ac.SCREEN_HEIGHT - y) / 2,
		W:                 x,
		H:                 y,
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
	close(d.MouseMotionEvents)
}

func (d *Dvd) Update() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case kevent := <-d.ControlEvents:
		d.controlUpdate(kevent)
	case mevent := <-d.MouseMotionEvents:
		d.X = mevent.X
		d.Y = mevent.Y
	default:
		d.boundaries()

		d.targetBehaviour()

		d.X += directionX
		d.Y += directionY
	}
}
