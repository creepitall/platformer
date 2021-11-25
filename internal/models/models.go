package models

import "github.com/faiface/pixel"

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
