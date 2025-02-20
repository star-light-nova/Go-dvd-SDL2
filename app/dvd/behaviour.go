package dvd

import (
	ac "dvd/app/app_configs"

	"github.com/veandco/go-sdl2/sdl"

	"math/rand"
)

// Changes directions wheather DVD tries to escape the window.
func (d *Dvd) boundaries() {
	if d.X <= 0 || d.X+d.W >= ac.SCREEN_WIDTH {
		directionX = -directionX
	}

	if d.Y <= 0 || d.Y+d.H >= ac.SCREEN_HEIGHT {
		directionY = -directionY
	}

}

// Refreshes target values (randomly)
// X: [0, WINDOW_WIDTH]
// Y: [0, WINDOW_HEIGHT]
func (d *Dvd) RefreshTargets() {
	xRand := rand.Intn(2)
	yRand := rand.Intn(2)

	d.TargetX = xtargets[xRand]
	d.TargetY = ytargets[yRand]
}

// * When IsTargetX and/or IsTargetY are/is set, we do not redirect from
// walls in indifferent direction, but wait till we reach the target in
// other axis.
// * When our position + WIDTH/HEIGHT is less than target we have positive
// directions, otherwise negative.
// * When we reach the target, we refresh them, until the `IsTargetX` or
// `IsTargetY` are not false.
func (d *Dvd) targetBehaviour() {
	d.mu.RLock()
	defer d.mu.RUnlock()

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

	// When we reach the target, choose random X,Y targets.
	if d.IsTargetX == d.IsTargetY && d.IsTargetX == true {
		if (d.X == d.TargetX || d.X+d.W == d.TargetX) &&
			(d.Y == d.TargetY || d.Y+d.H == d.TargetY) {
			d.RefreshTargets()
		}
	}
}

// Allows arrow keyboard to change the direction of the DVD.
func (d *Dvd) controlUpdate(kevent *sdl.KeyboardEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// It was causing overflow while targeted.
	if d.IsTargetX || d.IsTargetY {
		return
	}

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
