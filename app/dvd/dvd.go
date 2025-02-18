package dvd

import (
	"fmt"
	"math/rand"
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

	var x, y int32 = 156, 128

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

var directionX int32 = 1
var directionY int32 = 1

var Ytargets = [2]int32{0, ac.SCREEN_HEIGHT}
var Xtargets = [2]int32{0, ac.SCREEN_WIDTH}

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

		if d.IsTargetY {
			if d.Y == d.TargetY {
				// do nothing
			} else if d.Y+d.H < d.TargetY {
				directionY = abs(directionY)
			} else {
				directionY = -abs(directionY)
			}
		}

		if d.IsTargetX {
			if d.X == d.TargetX {
				// do nothing
			} else if d.X+d.W < d.TargetX {
				directionX = abs(directionX)
			} else {
				directionX = -abs(directionX)
			}
		}

		if d.IsTargetX == d.IsTargetY && d.IsTargetX == true {
			if (d.X == d.TargetX || d.X+d.W == d.TargetX) &&
				(d.Y == d.TargetY || d.Y+d.H == d.TargetY) {
				// [0-1]
				xRand := rand.Intn(2)
				yRand := rand.Intn(2)

				d.TargetX = Xtargets[xRand]
				d.TargetY = Ytargets[yRand]
			}
		}

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
	if d.X <= 0 || d.X+d.W >= ac.SCREEN_WIDTH {
		directionX = -directionX
	}

	if d.Y <= 0 || d.Y+d.H >= ac.SCREEN_HEIGHT {
		directionY = -directionY
	}

}

func (d *Dvd) isHover(mevent *sdl.MouseMotionEvent) bool {
	if mevent.X >= d.X && mevent.X <= d.X+d.W {
		if mevent.Y >= d.Y && mevent.Y <= d.Y+d.H {
			return true
		}
	}

	return false
}

func abs(number int32) int32 {
	if number < 0 {
		return -number
	}

	return number
}
