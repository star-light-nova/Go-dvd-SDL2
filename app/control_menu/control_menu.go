package control_menu

import (
	ac "dvd/app/app_configs"
	"dvd/app/button"
	"dvd/app/checkbox"
	"dvd/app/dvd"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type ControlMenu struct {
	texture *sdl.Texture

	KeyEvents         chan *sdl.KeyboardEvent
	MouseButtonEvents chan *sdl.MouseButtonEvent
	isOpen            bool

	button   *button.Button
	checkbox *checkbox.Checkbox

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

	keyEvents := make(chan *sdl.KeyboardEvent)
	mouseButtonEvents := make(chan *sdl.MouseButtonEvent)

	b, err := button.NewButton(r, "Reset position")
	if err != nil {
		return nil, fmt.Errorf("Could not create Button: %v", err)
	}

	c, err := checkbox.NewCheckbox(r, "Always hit corners")
	if err != nil {
		return nil, fmt.Errorf("Could not create Checkbox: %v", err)
	}

	return &ControlMenu{
		texture: texture,

		KeyEvents:         keyEvents,
		MouseButtonEvents: mouseButtonEvents,

		isOpen: false,

		button:   b,
		checkbox: c,

		X: -250,
		Y: 0,
		W: 250,
		H: ac.SCREEN_HEIGHT,
	}, nil
}

func (cm *ControlMenu) Update(dvd *dvd.Dvd) {
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
	case kevent := <-cm.KeyEvents:
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
	case mbevent := <-cm.MouseButtonEvents:
		if cm.isOpen {
			if cm.checkbox.IsHover(mbevent) {
				if cm.checkbox.Click(mbevent) {
					dvd.IsTargetX = cm.checkbox.IsSelected
					dvd.IsTargetY = dvd.IsTargetX
				}
			} else if cm.button.IsHover(mbevent) {
				if cm.button.Click(mbevent) {
					dvd.Y = (ac.SCREEN_HEIGHT - dvd.H) / 2
					dvd.X = (ac.SCREEN_WIDTH - dvd.W) / 2
				}
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

	cm.button.Y = cm.Y + 250

	buttonRect := &sdl.Rect{
		X: cm.X,
		Y: cm.button.Y,
		W: cm.button.W,
		H: cm.button.H,
	}

	cm.checkbox.Y = cm.Y + 200

	checkboxRect := &sdl.Rect{
		X: cm.X,
		Y: cm.checkbox.Y,
		W: cm.checkbox.W,
		H: cm.checkbox.H,
	}

	if err := r.SetDrawColor(130, 126, 126, 0); err != nil {
		return err
	}

	if err := r.Copy(cm.texture, nil, rect); err != nil {
		return err
	}

	if err := r.FillRect(rect); err != nil {
		return err
	}

	if err := r.Copy(cm.button.Texture(), nil, buttonRect); err != nil {
		return err
	}

	if err := r.Copy(cm.checkbox.Texture(), nil, checkboxRect); err != nil {
		return err
	}

	return nil
}

func (cm *ControlMenu) Destroy() {
	cm.texture.Destroy()
	cm.button.Destroy()
	cm.checkbox.Destroy()

	close(cm.KeyEvents)
	close(cm.MouseButtonEvents)
}
