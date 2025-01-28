package font_config

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var FONT_PATH string = fmt.Sprintf("%s/assets/fonts/Roboto-Bold.ttf", sdl.GetBasePath())
var FONT_COLOR sdl.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}

const FONT_SIZE int = 27
