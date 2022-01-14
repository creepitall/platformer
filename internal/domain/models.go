package domain

import (
	"runtime"
	"time"

	"github.com/faiface/pixel"
)

const (
	GlobalGravity = -512
)

var (
	PreviousTime time.Time
)

// Главный герой фреймы
var (
	HeroPlayerRunFrames  []pixel.Rect
	HeroPlayerStayFrames []pixel.Rect
	HeroPlayerJumpFrames []pixel.Rect
)

// Главный герой ассеты
var (
	HeroPlayerRunAssets  pixel.Picture
	HeroPlayerStayAssets pixel.Picture
	HeroPlayerJumpAssets pixel.Picture
)

// TODO: вынести в core
func ReturnFilePath(fp string) string {
	var path string
	if runtime.GOOS == "windows" {
		path += ".../../../"
	}
	return path + fp
}
