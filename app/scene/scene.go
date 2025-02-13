package scene

import (
	"dvd/app/control_menu"
	"dvd/app/dvd"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"fmt"
)

type Scene struct {
	dvd         *dvd.Dvd
	controlMenu *control_menu.ControlMenu
}

func NewScene(r *sdl.Renderer) (*Scene, error) {
	d, err := dvd.NewDvd(r)
	if err != nil {
		return nil, fmt.Errorf("Could not create a DVD: %v", err)
	}

	cm, err := control_menu.NewControlMenu(r)
	if err != nil {
		return nil, fmt.Errorf("Could not create Control Menu: %v", err)
	}

	scene := &Scene{dvd: d, controlMenu: cm}

	return scene, nil
}

func (scene *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)

		tick := time.Tick(10 * time.Millisecond)

		for {
			select {
			case e := <-events:
				if isExit := scene.handleEvent(e); isExit {
					return
				}
			case <-tick:
				scene.Update()

				if err := scene.Paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (scene *Scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		mouseEvent := event.(*sdl.MouseButtonEvent)

		go func() { scene.controlMenu.MouseButtonEvents <- mouseEvent }()
	case *sdl.KeyboardEvent:
		kevent := event.(*sdl.KeyboardEvent)

		go func() {
			scene.dvd.ControlEvents <- kevent
			scene.controlMenu.KeyEvents <- kevent
		}()
	}

	return false
}

func (scene *Scene) Update() {
	scene.dvd.Update()
	scene.controlMenu.Update(scene.dvd)
}

func (scene *Scene) Paint(r *sdl.Renderer) error {
	// Default colour
	if err := r.SetDrawColor(0, 0, 0, 0); err != nil {
		return err
	}

	r.Clear()

	if err := scene.controlMenu.Paint(r); err != nil {
		return fmt.Errorf("Could not scene the Control Menu: %v", err)
	}

	if err := scene.dvd.Paint(r); err != nil {
		return fmt.Errorf("Could not scene the dvd: %v", err)
	}

	r.Present()

	return nil
}

func (scene *Scene) Destroy() {
	scene.dvd.Destroy()
	scene.controlMenu.Destroy()
}
