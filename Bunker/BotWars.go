package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

// Global Bot Script variables
var InitLoadComplete bool
var BotRespawn bool
var AllManaObelisksOnMap []ns.Obj
var NoTarget ns.Obj
var BotMana bool
var soloPlay bool
var soloPlayer ns.Obj
var soloPlayerHasFlag bool
var botSpawnsNoTeams []ns.Obj

// Server settings

// -- General Bot Settings
var BotDifficulty int
var TeamsEnabled bool

// Capture the Flag
var GameModeIsCTF bool
var BlueTeamBase ns.Obj
var RedTeamBase ns.Obj
var BlueFlag ns.Obj
var RedFlag ns.Obj

// King of the Realm
var GameModeIsTeamKOTR bool
var Crowns []ns.Obj
var CrownRed ns.Obj
var CrownBlue ns.Obj

// Arena
var GameModeIsTeamArena bool

var TestManaObelisk ns.Obj

func init() {
	checkTeams()
	BotRespawn = true
	BotMana = true
	InitLoadComplete = false
	ns.NewTimer(ns.Frames(10), func() {
		getGameMode()
		loadMapObjects()
		checkSolo()
	})
	ns.NewTimer(ns.Frames(20), func() {
		Red.init()
		Blue.init()
	})
	ns.NewTimer(ns.Frames(60), func() {
		Red.lateInit()
		Blue.lateInit()
		InitLoadComplete = true
	})
	ns.OnChat(onCommand)
	NoTarget = ns.CreateObject("InvisibleExitArea", ns.Ptf(150, 150))
	ns.PrintStrToAll("Bot script installed successfully.")
}

func checkTeams() {
	AllTeams := ns.Teams()
	TeamsCheck := len(AllTeams)
	if TeamsCheck == 0 {
		TeamsEnabled = false
	} else {
		TeamsEnabled = true
	}
}

func loadMapObjects() {
	AllManaObelisksOnMap = ns.FindAllObjects(
		ns.HasTypeName{"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12", "LOTDManaObelisk"},
	)
	botSpawnsNoTeams = ns.FindAllObjects(ns.HasTypeName{"PlayerStart"})
}

func checkSolo() {
	if len(ns.Players()) > 1 {
		if soloPlay {
			// Solo play disabled.
			soloPlay = false
			//soloPlayerHasFlag = false
		}
	} else {
		if !soloPlay {
			// Solo play active.
			soloPlay = true
			soloPlayer = ns.Players()[0].Unit()
		}
	}
}

// Server Commands.
func onCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil {
		switch msg {
		case "test":
			ns.PrintStrToAll("Test")
			updateMyBotScriptData(p, func(data *MyAccountData) {
			})
		}
	}
	return msg
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
	checkSolo()
}

func getGameMode() {
	Flags := ns.FindAllObjects(ns.HasTypeName{"Flag"})
	Crowns = ns.FindAllObjects(ns.HasTypeName{"Crown"})
	if Flags != nil {
		GameModeIsTeamKOTR = false
		GameModeIsCTF = true
		GameModeIsTeamArena = false
		RedFlag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[0]})
		BlueFlag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[1]})

	} else if Crowns != nil {
		GameModeIsTeamKOTR = true
		GameModeIsCTF = false
		GameModeIsTeamArena = false
	} else {
		GameModeIsTeamKOTR = false
		GameModeIsCTF = false
		GameModeIsTeamArena = true
	}
}
