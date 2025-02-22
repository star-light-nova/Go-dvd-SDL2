package control_menu

import (
	ac "dvd/app/app_configs"

	"dvd/app/dvd"

	"github.com/veandco/go-sdl2/sdl"
)

func (cm *ControlMenu) menuDvdMouseButtonControl(dvd *dvd.Dvd, mbevent *sdl.MouseButtonEvent) {
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
}
