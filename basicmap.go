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
	spawns := ns.FindAllObjects(ns.HasTypeName{"PlayerStart"})
	randomIndex := ns.Random(0, len(spawns)-1)
	pick := spawns[randomIndex]
	RandomBotSpawn.SetPos(pick.Pos())
	UpdateBots()
}

func DialogOptions() {
	// Usable dialog bits
	// F1GD401E "What a Wizard spy?"

	// C2NC203E "Get away from me filthy peasants"
	// C2NC202E "URHGHH"

	// C5OGK02E "too bad you must die now"
	// C5OGK01E "youre very bold for such a little man"
	// C5OGK05E "Ill crush your bones"
	// WarAI.Chat("War01A.scr:Bully1") // this is a robbery! Your money AND your life!
	// ns.AudioEvent("F1ROG01E", WarAI)
	// TODO: Add audio to match the chat: F1ROG01E.

	//     TauntLaugh,
	//     TauntShakeFist,
	//     TauntPoint,

	//HumanMaleHurtLight,
	//     HumanMaleHurtMedium,
	//     HumanMaleHurtHeavy,

	//HumanFemaleHurtLight,
	//     HumanFemaleHurtMedium,
	//     HumanFemaleHurtHeavy,
}
