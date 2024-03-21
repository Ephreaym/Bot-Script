package BotWars

import (
	"fmt"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
)

type MyBotScriptData struct {
}

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
var AllBotsLinked bool
var HorstLinked bool
var LanceLinked bool
var KirikLinked bool

// Server settings
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
	ns.OnPlayerJoin(playerJoin)
	ns.OnPlayerLeave(playerLeave)
	ns.OnPlayerDeath(playerDeath)
}

func playerJoin(p ns.Player) bool {
	return true
}

func playerLeave(p ns.Player) {
}

func playerDeath(p ns.Player, k ns.Obj) {
	if k != nil {
	} else {
		// currently returns nil
	}

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
	ns.NewTimer(ns.Seconds(2), func() {
		checkSolo()
	})
}

//func ScoringSystem() {
//	if len(ns.Players()) >= 3 {
//		for i := 0; i < len(ns.Players()); i++ {
//			if ns.Players()[i].Name() == "Lance" {
//				arr := ns.FindAllObjects(ns.HasTypeName{"NPC"})
//				for i := 0; i < len(arr); i++ {
//					if arr[i].MaxHealth() == 150 {
//						ns.PrintStrToAll("I'm LANCE!")
//						arr[i].SetOwner(ns.Players()[i].Unit())
//						ns.Players()[i].Unit().SetPos(ns.Ptf(150, 150))
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoCollide)
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoUpdate)
//						LanceLinked = true
//					}
//				}
//			}
//			if ns.Players()[i].Name() == "Horst" {
//				arr := ns.FindAllObjects(ns.HasTypeName{"NPC"})
//				for i := 0; i < len(arr); i++ {
//					if arr[i].MaxHealth() == 100 {
//						ns.PrintStrToAll("I'm Horst!")
//						arr[i].SetOwner(ns.Players()[i].Unit())
//						ns.Players()[i].Unit().SetPos(ns.Ptf(150, 150))
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoCollide)
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoUpdate)
//						HorstLinked = true
//					}
//				}
//			}
//			if ns.Players()[i].Name() == "Kirik" {
//				arr := ns.FindAllObjects(ns.HasTypeName{"NPC"})
//				for i := 0; i < len(arr); i++ {
//					if arr[i].MaxHealth() == 75 {
//						ns.PrintStrToAll("I'm Kirik!")
//						arr[i].SetOwner(ns.Players()[i].Unit())
//						ns.Players()[i].Unit().SetPos(ns.Ptf(150, 150))
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoCollide)
//						ns.Players()[i].Unit().FlagsEnable(object.FlagNoUpdate)
//						KirikLinked = true
//					}
//				}
//			}
//		}
//		if KirikLinked && HorstLinked && LanceLinked {
//			AllBotsLinked = true
//			ns.PrintStrToAll("All bots linked")
//		}
//	}
//}

// Server Commands.
func onCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil {
		switch msg {
		case "test":
		// Spawn commands red bots.
		case "server spawn red war":
			bots = append(bots, NewWarrior(Red))
			ns.PrintStrToAll("A Warrior bot has joined team Red!")
		case "server spawn red con":
			bots = append(bots, NewConjurer(Red))
			ns.PrintStrToAll("A Conjurer bot has joined team Red!")
		case "server spawn red wiz":
			bots = append(bots, NewWizard(Red))
			ns.PrintStrToAll("A Wizard bot has joined team Red!")
			// Spawn commands blue bots.
		case "server spawn blue war":
			bots = append(bots, NewWarrior(Blue))
			ns.PrintStrToAll("A Warrior bot has joined team Blue!")
		case "server spawn blue con":
			bots = append(bots, NewConjurer(Blue))
			ns.PrintStrToAll("A Conjurer bot has joined team Blue!")
		case "server spawn blue wiz":
			bots = append(bots, NewWizard(Blue))
			ns.PrintStrToAll("A Wizard bot has joined team Blue!")
		case "server spawn bots 3v3":
			bots = append(bots, NewWarrior(Red))
			bots = append(bots, NewConjurer(Red))
			bots = append(bots, NewWizard(Red))
			bots = append(bots, NewWarrior(Blue))
			bots = append(bots, NewConjurer(Blue))
			bots = append(bots, NewWizard(Blue))
			ns.PrintStrToAll("Both the Red and Blue team now have 3 bots active!")
			// Remove all bots.
		case "server remove all bots":
			ns.PrintStrToAll("All bots have been removed from the game.")
			ns.FindObjects(
				func(it ns.Obj) bool {
					//BotRespawn = false
					ns.PrintStrToAll("Bot removal not yet implemented!")
					return true
				},
				ns.HasTypeName{"NPC"},
			)
			// Set bot difficulty.
		case "server hardcore bots":
			BotDifficulty = 0
			//BotMana = false
			//ns.PrintStrToAll("Bots difficulty set to hardcore.")
			ns.PrintStrToAll("Hardcore mode is disabled this build due to instability.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
		case "server hard bots":
			BotDifficulty = 15
			BotMana = true
			ns.PrintStrToAll("Bots difficulty set to hard.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
		case "server normal bots":
			BotDifficulty = 30
			BotMana = true
			ns.PrintStrToAll("Bots difficulty set to normal.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
		case "server easy bots":
			BotDifficulty = 45
			BotMana = true
			ns.PrintStrToAll("Bots difficulty set to easy.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
		case "server beginner bots":
			BotDifficulty = 60
			BotMana = true
			ns.PrintStrToAll("Bots difficulty set to beginner.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
			// Bot chat responses.
		case "hello", "Hello", "Yo", "yo", "what's up?", "What's up?", "hi", "Hi", "Hey", "hey", "sup", "Sup":
			it := ns.FindClosestObject(p.Unit(), ns.HasTypeName{"NPC"})
			random := ns.Random(1, 4)
			ns.NewTimer(ns.Seconds(1), func() {
				if random == 1 {
					it.ChatStr("Hey!")
				}
				if random == 2 {
					it.ChatStr("Hello!")
				}
				if random == 3 {
					it.ChatStr("Sup!")
				}
				if random == 4 {
					it.ChatStr("Greetings!")
				}
			})
		case "gg", "GG", "Gg", "GG!", "gg!", "Good game!", "good game", "Good game":
			it := ns.FindClosestObject(p.Unit(), ns.HasTypeName{"NPC"})
			random := ns.Random(1, 2)
			ns.NewTimer(ns.Seconds(1), func() {
				if random == 1 {
					it.ChatStr("GG!")
				}
				if random == 2 {
					it.ChatStr("Good game!")
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

func loadMyQuestData(pl ns.Player) MyBotScriptData {
	var data MyBotScriptData
	err := pl.Store(ns.Persistent{Name: "noxworld"}).Get("my-quest-name", &data)
	if err != nil {
		fmt.Println("cannot read quest data:", err)
	}
	return data
}

func saveMyQuestData(pl ns.Player, data MyBotScriptData) {
	err := pl.Store(ns.Persistent{Name: "noxworld"}).Set("my-quest-name", &data)
	if err != nil {
		fmt.Println("cannot save quest data:", err)
	}
}

func updateMyQuestData(pl ns.Player, fnc func(data *MyBotScriptData)) {
	data := loadMyQuestData(pl)
	fnc(&data)
	saveMyQuestData(pl, data)
}
