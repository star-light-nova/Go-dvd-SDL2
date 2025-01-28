package app

import (
	ac "dvd/app/app_configs"
	"dvd/app/scene"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

func Run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialise SDL: %v", err)
	}

	defer sdl.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(ac.SCREEN_WIDTH, ac.SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not create a window: %v", err)
	}

	// MacOS hack.
	sdl.PumpEvents()

	defer w.Destroy()

	scene, err := scene.NewScene(r)
	if err != nil {
		return fmt.Errorf("Couldn not create a scene: %v", err)
	}

	defer scene.Destroy()

	events := make(chan sdl.Event)
	errc := scene.Run(events, r)

	defer close(events)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}
