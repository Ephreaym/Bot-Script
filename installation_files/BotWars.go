package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
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

// Social
var GameModeIsSocial bool

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
	if p != nil && p == ns.GetHost().Player() {
		switch msg {
		// FFA commands
		case "-1 war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveWarBots > 0 {
					ns.PrintStrToAll("Warrior Bot left the game.")
					data.botscript.ActiveWarBots--
					println(data.botscript.ActiveWarBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}

				} else {
					ns.PrintStrToAll("No Warrior Bot active.")
				}
			})
		case "+1 war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveWarBots == 1 {
					ns.PrintStrToAll("Warrior Bot limit reached for FFA.")
				} else {
					ns.PrintStrToAll("Warrior Bot joined the game.")
					data.botscript.ActiveWarBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveConBots > 0 {
					ns.PrintStrToAll("Conjurer Bot left the game.")
					data.botscript.ActiveConBots--
					println(data.botscript.ActiveConBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Conjurer Bot active.")
				}
			})
		case "+1 con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveConBots == 1 {
					ns.PrintStrToAll("Conjurer Bot limit reached for FFA.")
				} else {
					ns.PrintStrToAll("Conjurer Bot joined the game.")
					data.botscript.ActiveConBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveWizBots > 0 {
					ns.PrintStrToAll("Wizard Bot left the game.")
					data.botscript.ActiveWizBots--
					println(data.botscript.ActiveWizBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Wizard Bot active.")
				}
			})
		case "+1 wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveWizBots == 1 {
					ns.PrintStrToAll("Wizard Bot limit reached for FFA.")
				} else {
					ns.PrintStrToAll("Wizard Bot joined the game.")
					data.botscript.ActiveWizBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
			// Team commands
		case "-1 red war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedWarBots > 0 {
					ns.PrintStrToAll("Red Warrior Bot left the game.")
					data.botscript.ActiveRedWarBots--
					println(data.botscript.ActiveRedWarBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Red Warrior Bot active.")
				}
			})
		case "+1 red war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedWarBots == 3 {
					ns.PrintStrToAll("Red Warrior Bot limit reached.")
				} else {
					ns.PrintStrToAll("Red Warrior Bot joined the game.")
					data.botscript.ActiveRedWarBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 red con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedConBots > 0 {
					ns.PrintStrToAll("Red Conjurer Bot left the game.")
					data.botscript.ActiveRedConBots--
					println(data.botscript.ActiveRedConBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Red Conjurer Bot active.")
				}
			})
		case "+1 red con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedConBots == 3 {
					ns.PrintStrToAll("Red Conjurer Bot limit reached.")
				} else {
					ns.PrintStrToAll("Red Conjurer Bot joined the game.")
					data.botscript.ActiveRedConBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 red wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedWizBots > 0 {
					ns.PrintStrToAll("Red Wizard Bot left the game.")
					data.botscript.ActiveRedWizBots--
					println(data.botscript.ActiveRedWizBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Red Wizard Bot active.")
				}
			})
		case "+1 red wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveRedWizBots == 3 {
					ns.PrintStrToAll("Red Wizard Bot limit reached.")
				} else {
					ns.PrintStrToAll("Red Wizard Bot joined the game.")
					data.botscript.ActiveRedWizBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 blue war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueWarBots > 0 {
					ns.PrintStrToAll("Blue Warrior Bot left the game.")
					data.botscript.ActiveBlueWarBots--
					println(data.botscript.ActiveBlueWarBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Blue Warrior Bot active.")
				}
			})
		case "+1 blue war":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueWarBots == 3 {
					ns.PrintStrToAll("Blue Warrior Bot limit reached.")
				} else {
					ns.PrintStrToAll("Blue Warrior Bot joined the game.")
					data.botscript.ActiveBlueWarBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 blue con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueConBots > 0 {
					ns.PrintStrToAll("Blue Conjurer Bot left the game.")
					data.botscript.ActiveBlueConBots--
					println(data.botscript.ActiveBlueConBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Blue Conjurer Bot active.")
				}
			})
		case "+1 blue con":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueConBots == 3 {
					ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
				} else {
					ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
					data.botscript.ActiveBlueConBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
			})
		case "-1 blue wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueWizBots > 0 {
					ns.PrintStrToAll("Blue Wizard Bot left the game.")
					data.botscript.ActiveBlueWizBots--
					println(data.botscript.ActiveBlueWizBots)
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				} else {
					ns.PrintStrToAll("No Blue Wizard Bot active.")
				}
			})
		case "+1 blue wiz":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				if data.botscript.ActiveBlueWizBots == 3 {
					ns.PrintStrToAll("Blue Wizard Bot limit reached.")
				} else {
					ns.PrintStrToAll("Blue Wizard Bot joined the game.")
					data.botscript.ActiveBlueWizBots++
					serverSettingSoundToAllPlayers := ns.Players()
					for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
						ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
					}
				}
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
	arenaobjects := ns.FindAllObjects(ns.HasTypeName{"Obelisk", "Quiver"})
	if arenaobjects == nil {
		GameModeIsSocial = true
	} else if Flags != nil {
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
