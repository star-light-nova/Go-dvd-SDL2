package control_menu

import (
	ac "dvd/app/app_configs"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type ControlMenu struct {
	texture *sdl.Texture

	ControlMenuEvents chan *sdl.KeyboardEvent
	isOpen            bool

	X, Y int32
	W, H int32
}

func NewControlMenu(r *sdl.Renderer) (*ControlMenu, error) {
	texture, err := r.CreateTexture(
		uint32(sdl.PIXELFORMAT_ABGR32),
		sdl.TEXTUREACCESS_TARGET,
		250,
		ac.SCREEN_HEIGHT,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not initialise Control Menu: %v", err)
	}

	controlMenuEvents := make(chan *sdl.KeyboardEvent)

	return &ControlMenu{
		texture:           texture,
		ControlMenuEvents: controlMenuEvents,
		isOpen:            false,
		X:                 -250,
		Y:                 0,
		W:                 250,
		H:                 ac.SCREEN_HEIGHT,
	}, nil
}

func (cm *ControlMenu) Update() {
	/*
				                     ＿＿
						　　　　　🌸＞　　フ
						　　　　　| 　_　 _ l
						　 　　　／` ミ＿xノ
						　　 　 /　　　 　 |
						　　　 /　 ヽ　　 ﾉ
						　 　 │　　|　|　|
						　／￣|　　 |　|　|
						　| (￣ヽ＿_ヽ_)__)
						　＼二つ
		                What are you going to see after that
		                might shock you. Please, come back
		                to see the cat again, if you can't
		                handle it.
	*/
	select {
	case kevent := <-cm.ControlMenuEvents:
		if kevent.State == sdl.RELEASED {
			switch kevent.Keysym.Sym {
			case sdl.K_c:
				tick := time.Tick(500 * time.Microsecond)

				if cm.isOpen {
					go func() {
						for range tick {
							cm.X -= 1

							if cm.X == -cm.W {
								return
							}
						}
					}()
				} else {
					go func() {
						for range tick {
							cm.X += 1

							if cm.X == 0 {
								return
							}
						}
					}()
				}

				cm.isOpen = !cm.isOpen
			}
		}
	default:
		return
	}
}

func (cm *ControlMenu) Paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{
		X: cm.X,
		Y: cm.Y,
		W: cm.W,
		H: cm.H,
	}

	if err := r.SetDrawColor(128, 128, 64, 0); err != nil {
		return err
	}

	if err := r.Copy(cm.texture, nil, rect); err != nil {
		return err
	}

	if err := r.FillRect(rect); err != nil {
		return err
	}

	return nil
}

func (cm *ControlMenu) Destroy() {
	cm.texture.Destroy()

	close(cm.ControlMenuEvents)
}
