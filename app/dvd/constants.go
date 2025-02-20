package dvd

import ac "dvd/app/app_configs"

const (
	FILEPATH = "./assets/images/dvd_video_logo.png"
	x        = 156
	y        = 128
)

var directionX int32 = 1
var directionY int32 = 1

var ytargets = [2]int32{0, ac.SCREEN_HEIGHT}
var xtargets = [2]int32{0, ac.SCREEN_WIDTH}
