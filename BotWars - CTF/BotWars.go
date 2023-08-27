package BotWars

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

// Dialog options

// Con03B.scr:Worker1ChatD = I'll wait here
// Con03B.scr:Worker1ChatA = Let's go
// Con03B.scr:Worker1ChatB =Follow me
// Con02A:NecroTalk02 = Aaaaargh!
// client.c:Ping = Ping
// Con04a:NecroSaysDie = DIE, Conjurer!
// Con05:OgreKingTalk05 = I will crush your bones!
// Con05:OgreKingTalk07 = Din' din'. Come and get it!
// Con05:OgreTalk02 = What 'dat noise?
// Con05B.scr:OgreTaunt = A taunt.
// Con08b:InversionBoyTalk05 = Very good.
