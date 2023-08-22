package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

// General variables and booleans.
// General functions
var TestHitBox ns.Obj

var InitLoadComplete bool

// Behaviour CTF profiles
// TODO: move to the Team struct
var RedTank bool
var RedAttacker bool
var RedDefender bool

var BlueTank bool
var BlueAttacker bool
var BlueDefender bool

func init() {
	InitLoadComplete = false
	//RandomBotSpawn = ns.CreateObject("InvisibleExitArea", ns.GetHost())
	Red.init()
	Blue.init()
	ns.NewTimer(ns.Frames(60), func() {
		Red.lateInit()
		Blue.lateInit()
		TestHitBox = ns.Object("TestHitBox")
		//RedBotSpawn()
		//BlueBotSpawn()
		InitLoadComplete = true
	})
}

func OnFrame() {
	GetListOfPlayers()
	if !InitLoadComplete {
		return
	}
	// Script for bots that moves the flag towards them each frame.
	Red.PreUpdate()
	Blue.PreUpdate()
	UpdateBots()
	Red.PostUpdate()
	Blue.PostUpdate()
}

func GetListOfPlayers() {
	//AllPlayers := ns.Players()
}
