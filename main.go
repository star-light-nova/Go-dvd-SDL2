package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TODO: minimise
func main() {
	if err := AppRun(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}

// TODO: Separate
func AppRun() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialise SDL: %v", err)
	}

	defer sdl.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not create a window: %v", err)
	}

    sdl.PumpEvents()

	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("Could not draw a title: %v", err)
	}

	time.Sleep(3 * time.Second)

	return nil
}

// TODO: get rid of this one
func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialise ttf: %v", err)
	}

	f, err := ttf.OpenFont("./assets/static/Roboto-Bold.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}

	defer f.Close()

	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	s, err := f.RenderUTF8Blended("Salem Alem", c)
	if err != nil {
		return fmt.Errorf("Could not create a surface: %v", err)
	}

	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("Could not create a texture from surface: %v", err)
	}

	defer t.Destroy()

	err = r.Copy(t, nil, nil)
	if err != nil {
		return fmt.Errorf("Could not create a copy the texture: %v", err)
	}

	r.Present()

	return nil
}
