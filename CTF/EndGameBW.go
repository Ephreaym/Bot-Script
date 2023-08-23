package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

var InitLoadComplete bool

func init() {
	InitLoadComplete = false
	Red.init()
	Blue.init()
	ns.NewTimer(ns.Frames(60), func() {
		Red.lateInit()
		Blue.lateInit()
		InitLoadComplete = true
	})
}

func OnFrame() {
	if !InitLoadComplete {
		return
	}
	Red.PreUpdate()
	Blue.PreUpdate()
	UpdateBots()
	Red.PostUpdate()
	Blue.PostUpdate()
}
