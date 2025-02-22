package control_menu

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func (cm *ControlMenu) slidesOnKey(kevent *sdl.KeyboardEvent) {
	if kevent.State == sdl.RELEASED {
		switch kevent.Keysym.Sym {
		case sdl.K_c:
			tick := time.Tick(950 * time.Microsecond)

			if cm.isOpen {
				cm.slide(-1, tick)
			} else {
				cm.slide(1, tick)
			}

			cm.isOpen = !cm.isOpen
		}
	}
}

func (cm *ControlMenu) slide(step int32, tick <-chan time.Time) {
	go func() {
		for range tick {
			cm.X += step

			if cm.X == 0 || cm.X == -MENU_WIDTH {
				return
			}
		}
	}()
}
