package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
)

var InitLoadComplete bool
var GameModeIsCTF bool
var GameModeIsTeamArena bool
var BotDifficulty int
var BotRespawn bool
var BotMana bool
var NoTarget ns.Obj
var AllManaObelisksOnMap []ns.Obj
var BlueTeamBase ns.Obj
var RedTeamBase ns.Obj
var BlueFlag ns.Obj
var RedFlag ns.Obj
var TeamsEnabled bool
var ItemDropEnabled bool

func init() {
	checkTeams()
	ItemDropEnabled = true
	if TeamsEnabled {
		BotRespawn = true
		BotMana = true
		InitLoadComplete = false

		ns.NewTimer(ns.Frames(10), func() {
			CheckIfGameModeIsCTF()
			AllManaObelisksOnMap = ns.FindAllObjects(
				ns.HasTypeName{"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12"},
			)
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
	} else {
		ns.PrintStrToAll("Bot script installed successfully!")
		ns.PrintStrToAll("Bot script is disabled when playing without teams.")
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

// Server Commands.
func onCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil {
		switch msg {
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
			BotMana = false
			ns.PrintStrToAll("Bots difficulty set to hardcore.")
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
		// Bot commands.
		case "help", "Help", "Follow", "follow", "escort", "Escort", "come", "Come":
			ns.FindObjects(
				func(it ns.Obj) bool {
					if it.CanSee(p.Unit()) {
						it.Follow(p.Unit())
						random := ns.Random(1, 4)
						if random == 1 {
							it.ChatStr("I'll follow you.")
						}
						if random == 2 {
							it.ChatStr("Let's go.")
						}
						if random == 3 {
							it.ChatStr("I'll help.")
						}
						if random == 4 {
							it.ChatStr("Sure thing.")
						}
					}
					return true
				},
				ns.HasTypeName{"NPC"},
				ns.HasTeam{p.Team()},
			)
		case "guard", "stay", "Guard", "Stay":
			ns.FindObjects(
				func(it ns.Obj) bool {
					if it.CanSee(p.Unit()) {
						it.Guard(it.Pos(), it.Pos(), 300)
						random1 := ns.Random(1, 4)
						if random1 == 1 {
							it.ChatStr("I'll guard this place.")
						}
						if random1 == 2 {
							it.ChatStr("No problem.")
						}
						if random1 == 3 {
							it.ChatStr("I'll stay.")
						}
						if random1 == 4 {
							it.ChatStr("I'll hold.")
						}
					}
					return true
				},
				ns.HasTypeName{"NPC"},
				ns.HasTeam{p.Team()},
				ns.InCirclef{Center: p.Unit(), R: 100},
			)
		case "Attack", "Go", "go", "attack":
			ns.FindObjects(
				func(it ns.Obj) bool {
					if it.CanSee(p.Unit()) {
						it.Hunt()
						random2 := ns.Random(1, 4)
						if random2 == 1 {
							it.ChatStr("I'll get them.")
						}
						if random2 == 2 {
							it.ChatStr("Time to shine.")
						}
						if random2 == 3 {
							it.ChatStr("On the offense.")
						}
						if random2 == 4 {
							it.ChatStr("Time to hunt.")
						}
					}
					return true
				},
				ns.HasTypeName{"NPC"},
				ns.HasTeam{p.Team()},
			)
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
		case "server disable drops":
			ItemDropEnabled = false
			ns.PrintStrToAll("Bots no longer drop items on death.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
		case "server enable drops":
			ItemDropEnabled = true
			ns.PrintStrToAll("Bots drop items on death.")
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
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

func CheckIfGameModeIsCTF() {
	Flags := ns.FindAllObjects(ns.HasTypeName{"Flag"})
	if Flags != nil {
		GameModeIsCTF = true
		GameModeIsTeamArena = false
		ns.PrintStrToAll("Gamemode: Capture The Flag!")
	} else {
		GameModeIsCTF = false
		GameModeIsTeamArena = true
		ns.PrintStrToAll("Gamemode: Arena!")
	}
	ns.PrintStrToAll("Bot script installed successfully!")
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
