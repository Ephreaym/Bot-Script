package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/opennox-lib/object"
)

var RandomBotSpawn ns.Obj

// BlueBot Conditions
var TeamBlue ns.Obj
var BlueFlag ns.Obj
var BlueFlagFront ns.Obj
var BlueFlagOutOfGame ns.Pointf
var BlueBot ns.Obj
var BlueBotIsAlive bool
var BlueFlagOnBot bool
var BlueFlagIsAtBase bool
var BlueBase ns.Obj
var BlueBotIsAggressive bool

// BlueBot Items
var BlueHelm ns.Obj
var BlueChest ns.Obj
var BlueBoots ns.Obj
var BlueLegs ns.Obj
var BlueWrists ns.Obj
var BlueSword ns.Obj
var BlueShield ns.Obj
var BlueCloak ns.Obj

// RedBot Items
var RedHelm ns.Obj
var RedChest ns.Obj
var RedBoots ns.Obj
var RedLegs ns.Obj
var RedWrists ns.Obj
var RedSword ns.Obj
var RedShield ns.Obj
var RedCloak ns.Obj

// RedBot Conditions
var RedFlagIsAtBase bool
var RedBotIsAggressive bool
var TeamRed ns.Obj
var RedFlag ns.Obj
var RedFlagStart ns.Obj
var RedFlagFront ns.Obj
var RedFlagOutOfGame ns.Pointf
var RedBot ns.Obj
var RedBotIsAlive bool
var RedBase ns.Obj
var RedFlagOnBot bool

var InitLoadComplete bool

func init() {
	InitLoadComplete = false
	RandomBotSpawn = ns.CreateObject("InvisibleExitArea", ns.GetHost())
	ns.NewWaypoint("BotSpawnPoint", ns.GetHost().Pos())
	ns.NewTimer(ns.Frames(60), func() {
		BlueHelm = ns.Object("BlueHelm")
		BlueChest = ns.Object("BlueChest")
		BlueBoots = ns.Object("BlueBoots")
		BlueLegs = ns.Object("BlueLegs")
		BlueWrists = ns.Object("BlueWrists")
		BlueSword = ns.Object("BlueSword")
		BlueShield = ns.Object("BlueShield")
		BlueCloak = ns.Object("BlueCloak")
		RedHelm = ns.Object("RedHelm")
		RedChest = ns.Object("RedChest")
		RedBoots = ns.Object("RedBoots")
		RedLegs = ns.Object("RedLegs")
		RedWrists = ns.Object("RedWrists")
		RedSword = ns.Object("RedSword")
		RedShield = ns.Object("RedShield")
		RedCloak = ns.Object("RedCloak")
		TeamRed = ns.Object("TeamRed")
		RedFlag = ns.Object("RedFlag")
		RedBase = ns.Object("RedBase")
		BlueBase = ns.Object("BlueBase")
		ns.NewWaypoint("RedFlagStart", RedFlag.Pos())
		RedFlagFront = ns.Object("RedFlagFront")
		RedFlagOutOfGame = RedFlagFront.Pos()
		RedBotIsAlive = false
		RedFlagOnBot = false
		TeamBlue = ns.Object("TeamBlue")
		BlueFlag = ns.Object("BlueFlag")
		ns.NewWaypoint("BlueFlagStart", BlueFlag.Pos())
		BlueFlagFront = ns.Object("BlueFlagFront")
		BlueFlagOutOfGame = BlueFlagFront.Pos()
		BlueBotIsAlive = false
		BlueFlagOnBot = false
		BlueFlagIsAtBase = true
		RedFlagIsAtBase = true
		RedBotSpawn()
		BlueBotSpawn()
		ns.Players()
		InitLoadComplete = true
	})
}

func RedBotSpawn() {
	RedBot = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPoint"))
	RedBot.Equip(RedHelm)
	RedBot.Equip(RedChest)
	RedBot.Equip(RedBoots)
	RedBot.Equip(RedLegs)
	RedBot.Equip(RedWrists)
	RedBot.Equip(RedSword)
	RedBot.Equip(RedShield)
	RedBot.Equip(RedCloak)
	RedBotIsAlive = true
	RedBotIsAggressive = false
	RedBot.SetOwner(TeamRed)
	RedBot.MonsterStatusEnable(object.MonStatusAlwaysRun)
	RedBot.SetMaxHealth(150)
	RedBot.AggressionLevel(0.16)
	RedBot.SetStrength(125)
	RedBot.SetBaseSpeed(100)
	RedBot.RetreatLevel(0)
	RedBot.OnEvent(ns.EventLookingForEnemy, func() {
		//RedBot.Chat("EventLookForEnemy")
		RedBotReturnFlag()
	})
	RedBot.OnEvent(ns.EventEndOfWaypoint, func() {
		//RedBot.Chat("EventEndOfWaypoint")
		RedBotReturnFlag()
	})
	RedBot.OnEvent(ns.EventChangeFocus, func() {
		//RedBot.Chat("EventChangeFocus")
	})
	RedBot.OnEvent(ns.EventEnemyHeard, func() {
		//RedBot.Chat("EventEnemyHeard")
	})
	RedBot.OnEvent(ns.EventLostEnemy, func() {
		//RedBot.Chat("EventLostEnemy")
	})
	RedBot.OnEvent(ns.EventEnemySighted, func() {
		//RedBot.Chat("EventEnemySighted")
	})
	RedBot.OnEvent(ns.EventIsHit, func() {
		//RedBot.Chat("EventIsHit")
	})
	RedBot.OnEvent(ns.EventRetreat, func() {
		//RedBot.Chat("EventRetreat")
	})

	RedBot.OnEvent(ns.EventCollision, func() {
		// Pickup the flag
		if ns.GetCaller() == BlueFlag && !BlueFlagOnBot && RedBotIsAlive {
			BlueFlag.SetPos(BlueFlagOutOfGame)
			ns.AudioEvent(audio.FlagPickup, RedBot)
			BlueFlagOnBot = true
			BlueFlagIsAtBase = false
			RedBotReturnFlag()
			//RedBot.Chat("EventCollision: Pickup the flag")
		}
		// Capture the flag
		if ns.GetCaller() == RedFlag && BlueFlagOnBot && RedFlagIsAtBase && RedBotIsAlive {
			ns.AudioEvent(audio.FlagCapture, RedBot)
			BlueFlagOnBot = false
			FlagReset()
			//RedBot.Chat("EventCollision: Caputure the Flag")
		}
		// Retrieve the flag
		if ns.GetCaller() == RedFlag && !RedFlagIsAtBase && RedBotIsAlive {
			RedFlagOnBot = false
			RedFlagIsAtBase = true
			ns.AudioEvent(audio.FlagRespawn, RedBot)
			RedFlagFront.SetPos(RedFlagOutOfGame)
			RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
			RedBotReturnFlag()
			//RedBot.Chat("EventCollision: Retrieve the flag")
		}
	})
	RedBot.OnEvent(ns.EventDeath, func() {
		if BlueFlagOnBot {
			ns.AudioEvent(audio.FlagDrop, RedBot)
			RedBotIsAlive = false
			BlueFlagOnBot = false
			BlueFlagFront.SetPos(TeamBlue.Pos())
			BlueFlag.SetPos(RedBot.Pos())
			ns.NewTimer(ns.Frames(60), func() {
				RedBotSpawn()
			})
			BlueBotReturnFlag()
		} else {
			RedBotIsAlive = false
			BlueFlagOnBot = false
			ns.NewTimer(ns.Frames(60), func() {
				RedBotSpawn()
			})
		}
		RedBot.Delete()
	})
	RedBot.WalkTo(BlueFlag.Pos())
}

func RedBotReturnFlag() {
	RedBot.AggressionLevel(0.16)
	if !RedFlagOnBot {
		RedBot.WalkTo(RedFlag.Pos())
	} else {
		RedBot.WalkTo(RedFlagFront.Pos())
	}
}

func BlueBotReturnFlag() {
	BlueBot.AggressionLevel(0.16)
	if !BlueFlagOnBot {
		BlueBot.WalkTo(BlueFlag.Pos())
	} else {
		BlueBot.WalkTo(BlueFlagFront.Pos())
	}
}

func BlueBotSpawn() {
	BlueBot = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPoint"))
	BlueBotIsAlive = true
	BlueBotIsAggressive = false
	BlueBot.Equip(BlueHelm)
	BlueBot.Equip(BlueChest)
	BlueBot.Equip(BlueBoots)
	BlueBot.Equip(BlueLegs)
	BlueBot.Equip(BlueWrists)
	BlueBot.Equip(BlueSword)
	BlueBot.Equip(BlueShield)
	BlueBot.Equip(BlueCloak)
	BlueBot.SetOwner(TeamBlue)
	BlueBot.MonsterStatusEnable(object.MonStatusAlwaysRun)
	BlueBot.SetMaxHealth(150)
	BlueBot.AggressionLevel(0.16)
	BlueBot.SetStrength(125)
	BlueBot.SetBaseSpeed(100)
	BlueBot.RetreatLevel(0)
	BlueBot.OnEvent(ns.EventLookingForEnemy, func() {
		BlueBotReturnFlag()
	})
	BlueBot.OnEvent(ns.EventEndOfWaypoint, func() {
		//BlueBot.Chat("EventEndOfWaypoint")
		BlueBotReturnFlag()
	})
	BlueBot.OnEvent(ns.EventChangeFocus, func() {
		//BlueBot.Chat("EventChangeFocus")
	})
	BlueBot.OnEvent(ns.EventEnemyHeard, func() {
		//BlueBot.Chat("EventEnemyHeard")
	})
	BlueBot.OnEvent(ns.EventLostEnemy, func() {
		//BlueBot.Chat("EventLostEnemy")
	})
	BlueBot.OnEvent(ns.EventEnemySighted, func() {
		//BlueBot.Chat("EventEnemySighted")
	})
	BlueBot.OnEvent(ns.EventIsHit, func() {
		//BlueBot.Chat("EventIsHit")
	})
	BlueBot.OnEvent(ns.EventRetreat, func() {
		//BlueBot.Chat("EventRetreat")
	})
	BlueBot.OnEvent(ns.EventCollision, func() {
		// Pickup the flag
		if ns.GetCaller() == RedFlag && !RedFlagOnBot && BlueBotIsAlive {
			RedFlag.SetPos(RedFlagOutOfGame)
			ns.AudioEvent(audio.FlagPickup, BlueBot)
			BlueBotReturnFlag()
			RedFlagOnBot = true
			RedFlagIsAtBase = false
			//BlueBot.Chat("EventCollision: Pickup the flag")
		}
		// Capture the flag
		if ns.GetCaller() == BlueFlag && RedFlagOnBot && BlueFlagIsAtBase && BlueBotIsAlive {
			ns.AudioEvent(audio.FlagCapture, BlueBot)
			RedFlagOnBot = false
			FlagReset()
			//BlueBot.Chat("EventCollision: Caputure the Flag")
		}
		// Retrieve the flag
		if ns.GetCaller() == BlueFlag && !BlueFlagIsAtBase && BlueBotIsAlive {
			BlueFlagOnBot = false
			BlueFlagIsAtBase = true
			ns.AudioEvent(audio.FlagRespawn, BlueBot)
			BlueFlagFront.SetPos(BlueFlagOutOfGame)
			BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
			BlueBotReturnFlag()
			//BlueBot.Chat("EventCollision: Retrieve the flag")
		}
	})
	BlueBot.OnEvent(ns.EventDeath, func() {
		if RedFlagOnBot {
			ns.AudioEvent(audio.FlagDrop, BlueBot)
			BlueBotIsAlive = false
			RedFlagOnBot = false
			RedFlagFront.SetPos(TeamRed.Pos())
			RedFlag.SetPos(BlueBot.Pos())
			ns.NewTimer(ns.Frames(60), func() {
				BlueBotSpawn()
			})
			RedBotReturnFlag()
		} else {
			BlueBotIsAlive = false
			RedFlagOnBot = false
			ns.NewTimer(ns.Frames(60), func() {
				BlueBotSpawn()
			})
		}
		BlueBot.Delete()
	})
	BlueBot.WalkTo(RedFlag.Pos())
}

func FlagReset() {
	RedFlagFront.SetPos(RedFlagOutOfGame)
	BlueFlagFront.SetPos(BlueFlagOutOfGame)
	RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
	BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
	RedBot.WalkTo(BlueFlag.Pos())
	BlueBot.WalkTo(RedFlag.Pos())
}

func OnFrame() {
	if InitLoadComplete {
		RedBase.SetPos(RedFlagFront.Pos())
		BlueBase.SetPos(BlueFlagFront.Pos())
		spawns := ns.FindAllObjects(ns.HasTypeName{"PlayerStart"})
		randomIndex := ns.Random(0, len(spawns)-1)
		pick := spawns[randomIndex]
		ns.Waypoint("BotSpawnPoint").SetPos(pick.Pos())
		UpdateBots()
	}
	if BlueFlagOnBot {
		BlueFlag.SetPos(BlueFlagOutOfGame)
		BlueFlagFront.SetPos(RedBot.Pos())
	}
	if RedFlagOnBot {
		RedFlag.SetPos(RedFlagOutOfGame)
		RedFlagFront.SetPos(BlueBot.Pos())
	}
}
