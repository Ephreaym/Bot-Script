package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
)

// General variables and booleans.
// General functions
var TestHitBox ns.Obj

// BlueBot Conditions
var TeamBlue ns.Obj
var BlueFlag ns.Obj
var BlueBot ns.Obj
var BlueBotIsAlive bool
var BlueFlagOnBot bool
var BlueFlagIsAtBase bool
var BlueBase ns.Obj
var BlueBotIsAggressive bool

// RedBot Conditions
var RedFlagIsAtBase bool
var RedBotIsAggressive bool
var TeamRed ns.Obj
var RedFlag ns.Obj
var RedFlagStart ns.Obj
var RedBot ns.Obj
var RedBotIsAlive bool
var RedBase ns.Obj
var RedFlagOnBot bool

var InitLoadComplete bool

// Behaviour CTF profiles
var RedTank bool
var RedAttacker bool
var RedDefender bool
var RedTankObject ns.Obj

var BlueTank bool
var BlueAttacker bool
var BlueDefender bool

func init() {
	InitLoadComplete = false
	//RandomBotSpawn = ns.CreateObject("InvisibleExitArea", ns.GetHost())
	ns.NewWaypoint("BotSpawnPointRed", ns.GetHost().Pos())
	ns.NewWaypoint("BotSpawnPointBlue", ns.GetHost().Pos())
	ns.NewTimer(ns.Frames(60), func() {
		TeamRed = ns.Object("TeamRed")
		RedFlag = ns.Object("RedFlag")
		RedBase = ns.Object("RedBase")
		BlueBase = ns.Object("BlueBase")
		TestHitBox = ns.Object("TestHitBox")
		ns.NewWaypoint("RedFlagStart", RedFlag.Pos())
		TeamBlue = ns.Object("TeamBlue")
		BlueFlag = ns.Object("BlueFlag")
		ns.NewWaypoint("BlueFlagStart", BlueFlag.Pos())
		BlueFlagIsAtBase = true
		RedFlagIsAtBase = true
		ns.NewWaypoint("RedFlagWoint", RedFlag.Pos())
		ns.NewWaypoint("BlueFlagWaypoint", BlueFlag.Pos())
		RedTeamTank = TeamRed
		BlueTeamTank = TeamBlue
		//RedBotSpawn()
		//BlueBotSpawn()
		InitLoadComplete = true
	})
}

func FlagReset() {
	RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
	BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
	RedFlag.Enable(true)
	BlueFlag.Enable(true)
}

func OnFrame() {
	GetListOfPlayers()
	if InitLoadComplete {
		MoveEquipedFlagWithBot()
		RandomizeBotSpawnCTF()
		CheckIfFlagsAreAtBase()
		BotConditionsWhileCarryingTheFlag()
	}
}

// CTF objects and booleans.
var RedTeamTank ns.Obj
var BlueTeamTank ns.Obj

func MoveEquipedFlagWithBot() {
	// Script for bots that moves the flag towards them each frame.
	if !BlueFlag.IsEnabled() {
		// Move the real blue flag out of the map.
		// Move the fake blue flag on the bot.
		BlueFlag.SetPos(RedTeamTank.Pos())
	}
	if !RedFlag.IsEnabled() {
		// Move the real red flag out of the map.
		// Move the fake red flag on the bot.
		RedFlag.SetPos(BlueTeamTank.Pos())
	}
}

func BotConditionsWhileCarryingTheFlag() {
	// Remove buffs from bots that can't be active while carrying the flag.
	if RedTeamTank.HasEnchant(enchant.INVISIBLE) {
		RedTeamTank.EnchantOff(enchant.INVISIBLE)
	}
	if RedTeamTank.HasEnchant(enchant.INVULNERABLE) {
		RedTeamTank.EnchantOff(enchant.INVULNERABLE)
	}
	if !RedTeamTank.HasEnchant(enchant.VILLAIN) {
		RedTeamTank.Enchant(enchant.VILLAIN, ns.Seconds(60))
	}
	if BlueTeamTank.HasEnchant(enchant.INVISIBLE) {
		BlueTeamTank.EnchantOff(enchant.INVISIBLE)
	}
	if BlueTeamTank.HasEnchant(enchant.INVULNERABLE) {
		BlueTeamTank.EnchantOff(enchant.INVULNERABLE)
	}
	if !BlueTeamTank.HasEnchant(enchant.VILLAIN) {
		BlueTeamTank.Enchant(enchant.VILLAIN, ns.Seconds(60))
	}
}

func RandomizeBotSpawnCTF() {
	// Script to select a random PlayerStart for the bot to spawn on.
	// Filter to only select PlayStart objects that are owned by the red team.
	var spawnsRed []ns.Obj
	filterRed := ns.HasTypeName{"PlayerStart"}
	ns.ObjectGroup("TeamRed").EachObject(true, func(it ns.Obj) bool {
		if filterRed.Matches(it) {
			spawnsRed = append(spawnsRed, it)
		}
		return true // keep iterating in any case
	})
	// Filter to only select PlayerStart objects that are owned by the blue team.
	var spawnsBlue []ns.Obj
	filterBlue := ns.HasTypeName{"PlayerStart"}
	ns.ObjectGroup("TeamBlue").EachObject(true, func(it ns.Obj) bool {
		if filterBlue.Matches(it) {
			spawnsBlue = append(spawnsBlue, it)
		}
		return true // keep iterating in any case
	})
	if InitLoadComplete {
		RedBase.SetPos(RedFlag.Pos())
		BlueBase.SetPos(BlueFlag.Pos())
		//spawns := ns.FindAllObjects(ns.HasTypeName{"PlayerStart"}) // <---- Use this when no teams are used.
		randomIndexRed := ns.Random(0, len(spawnsRed)-1)
		randomIndexBlue := ns.Random(0, len(spawnsRed)-1)
		pickRed := spawnsRed[randomIndexRed]
		pickBlue := spawnsBlue[randomIndexBlue]
		ns.Waypoint("BotSpawnPointRed").SetPos(pickRed.Pos())
		ns.Waypoint("BotSpawnPointBlue").SetPos(pickBlue.Pos())
		UpdateBots()
	}
}

func CheckIfFlagsAreAtBase() {
	if (ns.InCirclef{Center: RedBase, R: 20}).Matches(RedFlag) {
		RedFlagIsAtBase = true
	} else {
		RedFlagIsAtBase = false
	}
	if (ns.InCirclef{Center: BlueBase, R: 20}).Matches(BlueFlag) {
		BlueFlagIsAtBase = true
	} else {
		BlueFlagIsAtBase = false
	}
}

func GetListOfPlayers() {
	AllPlayers := ns.Players()
}
