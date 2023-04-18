package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

var RandomBotSpawn ns.Obj

func init() {
	ns.Music(15, 20)
	RandomBotSpawn = ns.CreateObject("InvisibleExitArea", ns.GetHost())
}

// OnFrame is called by the server.
func OnFrame() {
	UpdateBots()
	spawns := ns.FindAllObjects(ns.HasTypeName{"PlayerStart"})
	randomIndex := ns.Random(0, len(spawns)-1)
	pick := spawns[randomIndex]
	RandomBotSpawn.SetPos(pick.Pos())
}

func DialogOptions() {
	// Usable dialog bits
	// F1GD401E "What a Wizard spy?"

	// C2NC203E "Get away from me filthy peasants"
	// C2NC202E "URHGHH"

	// C5OGK02E "too bad you must die now"
	// C5OGK01E "youre very bold for such a little man"
	// C5OGK05E "Ill crush your bones"
}
