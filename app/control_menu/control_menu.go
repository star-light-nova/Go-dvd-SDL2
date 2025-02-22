package control_menu

import (
	ac "dvd/app/app_configs"

	"dvd/app/button"
	"dvd/app/checkbox"
	"dvd/app/dvd"

	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type ControlMenu struct {
	texture *sdl.Texture

	KeyEvents         chan *sdl.KeyboardEvent
	MouseButtonEvents chan *sdl.MouseButtonEvent

	isOpen bool

	button   *button.Button
	checkbox *checkbox.Checkbox

	X, Y int32
	W, H int32
}

func NewControlMenu(r *sdl.Renderer) (*ControlMenu, error) {
	texture, err := r.CreateTexture(
		uint32(sdl.PIXELFORMAT_ABGR32),
		sdl.TEXTUREACCESS_TARGET,
		MENU_WIDTH,
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

		X: -MENU_WIDTH,
		Y: 0,
		W: MENU_WIDTH,
		H: ac.SCREEN_HEIGHT,
	}, nil
}

func (cm *ControlMenu) Update(dvd *dvd.Dvd) {
	select {
	case kevent := <-cm.KeyEvents:
		cm.slidesOnKey(kevent)
	case mbevent := <-cm.MouseButtonEvents:
		cm.menuDvdMouseButtonControl(dvd, mbevent)
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

	// Pre center height
	cm.button.Y = cm.Y + 250

	buttonRect := &sdl.Rect{
		X: cm.X,
		Y: cm.button.Y,
		W: cm.button.W,
		H: cm.button.H,
	}

	// Right above button
	cm.checkbox.Y = cm.Y + 200

	checkboxRect := &sdl.Rect{
		X: cm.X,
		Y: cm.checkbox.Y,
		W: cm.checkbox.W,
		H: cm.checkbox.H,
	}

	// Some random numbers close to white.
	if err := r.SetDrawColor(130, 126, 126, 0); err != nil {
		return fmt.Errorf("Could not set a color: %v", err)
	}

	if err := r.Copy(cm.texture, nil, rect); err != nil {
		return fmt.Errorf("Could not copy a texture: %v", err)
	}

	if err := r.FillRect(rect); err != nil {
		return fmt.Errorf("Could not fill the colour: %v", err)
	}

	if err := r.Copy(cm.button.Texture(), nil, buttonRect); err != nil {
		return fmt.Errorf("Could not copy a button: %v", err)
	}

	if err := r.Copy(cm.checkbox.Texture(), nil, checkboxRect); err != nil {
		return fmt.Errorf("Could not copy checkbox: %v", err)
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
