package scene

import (
	"dvd/app/button"
	"dvd/app/checkbox"
	"dvd/app/dvd"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"fmt"
)

type Scene struct {
	dvd      *dvd.Dvd
	button   *button.Button
	checkbox *checkbox.Checkbox
}

func NewScene(r *sdl.Renderer) (*Scene, error) {
	d, err := dvd.NewDvd(r)
	if err != nil {
		return nil, fmt.Errorf("Could not create a DVD: %v", err)
	}

	b, err := button.NewButton(r, "Control")
	if err != nil {
		return nil, fmt.Errorf("Could not create Button: %v", err)
	}

	c, err := checkbox.NewCheckbox(r, "Always hit corners")
	if err != nil {
		return nil, fmt.Errorf("Could not create Checkbox: %v", err)
	}

	scene := &Scene{dvd: d, button: b, checkbox: c}

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

		if scene.checkbox.IsHover(mouseEvent) {
			scene.checkbox.Click(mouseEvent)
		} else if scene.button.IsHover(mouseEvent) {
			scene.button.Click(mouseEvent)
		}
	}

	return false
}

func (scene *Scene) Update() {
	scene.dvd.Update()
}

func (scene *Scene) Paint(r *sdl.Renderer) error {
	r.Clear()

	if err := scene.dvd.Paint(r); err != nil {
		return fmt.Errorf("Could not scene the dvd: %v", err)
	}

	if err := scene.button.Paint(r); err != nil {
		return fmt.Errorf("Could not scene the button: %v", err)
	}

	if err := scene.checkbox.Paint(r); err != nil {
		return fmt.Errorf("Could not scene the checkbox: %v", err)
	}

	r.Present()

	return nil
}

func (scene *Scene) Destroy() {
	scene.dvd.Destroy()
	scene.button.Destroy()
	scene.checkbox.Destroy()
}
