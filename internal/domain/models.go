package domain

import (
	"runtime"

	"github.com/faiface/pixel"
)

var DefaultSprites map[string]*pixel.Sprite

var SceneSprites map[string][]*pixel.Sprite

var CurrentScene string

var (
	HeroPlayerRunFrames  []pixel.Rect
	HeroPlayerStayFrames []pixel.Rect
	HeroPlayerJumpFrames []pixel.Rect
)

var (
	HeroPlayerRunAssets  pixel.Picture
	HeroPlayerStayAssets pixel.Picture
	HeroPlayerJumpAssets pixel.Picture
)

// it's bag or not?
func ReturnFilePath(fp string) string {
	var path string
	if runtime.GOOS == "windows" {
		path += ".../../../"
	}
	return path + fp
}
