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
var chatmapbotupdater Updater

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
	BotRespawn = true
	BotMana = true
	InitLoadComplete = false
	ns.NewTimer(ns.Frames(10), func() {
		checkTeams()
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
		// Server settings
		// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
		case "-easy":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				data.botscript.BotDifficultySetting = 45
				BotDifficulty = data.botscript.BotDifficultySetting
				data.botscript.BotManaSetting = true
				BotMana = true
			})
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
			ns.PrintStrToAll("Server settings changed to easy bot difficulty.")
		case "-normal":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				data.botscript.BotDifficultySetting = 30
				BotDifficulty = data.botscript.BotDifficultySetting
				data.botscript.BotManaSetting = true
				BotMana = true
			})
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
			ns.PrintStrToAll("Server settings changed to normal bot difficulty.")
		case "-hard":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				data.botscript.BotDifficultySetting = 15
				BotDifficulty = data.botscript.BotDifficultySetting
				data.botscript.BotManaSetting = true
				BotMana = true
			})
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
			ns.PrintStrToAll("Server settings changed to hard bot difficulty.")
		case "-insane":
			updateMyBotScriptData(p, func(data *MyAccountData) {
				data.botscript.BotDifficultySetting = 0
				BotDifficulty = data.botscript.BotDifficultySetting
				data.botscript.BotManaSetting = false
				BotMana = false
			})
			serverSettingSoundToAllPlayers := ns.Players()
			for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
				ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
			}
			ns.PrintStrToAll("Server settings changed to insane bot difficulty.")
		// FFA commands
		case "-1 war":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveWarBots > 0 {
						ns.PrintStrToAll("Warrior Bot left the game.")
						data.botscript.ActiveWarBots--
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() {
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Warrior Bot active.")
					}
				})
			}
		case "+1 war":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
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
						bots = append(bots, NewWarriorNoTeam())
					}
				})
			}
		case "-1 con":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveConBots > 0 {
						ns.PrintStrToAll("Conjurer Bot left the game.")
						data.botscript.ActiveConBots--
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() {
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))

								return
							}
						}
					} else {
						ns.PrintStrToAll("No Conjurer Bot active.")
					}
				})
			}
		case "+1 con":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
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
						bots = append(bots, NewConjurerNoTeam())
					}
				})
			}
		case "-1 wiz":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveWizBots > 0 {
						ns.PrintStrToAll("Wizard Bot left the game.")
						data.botscript.ActiveWizBots--
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() {
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Wizard Bot active.")
					}
				})
			}
		case "+1 wiz":
			if TeamsEnabled {
				ns.PrintStrToAll("Error: no team mentioned.")
			} else {
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
						bots = append(bots, NewWizardNoTeam())
					}
				})
			}
			// Team commands
		case "-1 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Warrior Bot left the game.")
								data.botscript.ActiveRedWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Red Warrior Bot active.")
					}
				})
			}
		case "-2 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Warrior Bot left the game.")
								data.botscript.ActiveRedWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Warrior Bot active.")
					}
				})
			}
		case "-3 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Warrior Bot left the game.")
								data.botscript.ActiveRedWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Warrior Bot active.")
					}
				})
			}
		case "+1 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewWarrior(Red))
					}
				})
			}
		case "+2 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWarBots == 3 {
						ns.PrintStrToAll("Red Warrior Bot limit reached.")
					} else if data.botscript.ActiveRedWarBots == 2 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot limit reached.")
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
					} else if data.botscript.ActiveRedWarBots == 1 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						data.botscript.ActiveRedWarBots++
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					} else if data.botscript.ActiveRedWarBots == 0 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						data.botscript.ActiveRedWarBots++
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					}
				})
			}
		case "+3 red war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWarBots == 3 {
						ns.PrintStrToAll("Red Warrior Bot limit reached.")
					} else if data.botscript.ActiveRedWarBots == 2 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot limit reached.")
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
					} else if data.botscript.ActiveRedWarBots == 1 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot limit reached.")
						data.botscript.ActiveRedWarBots++
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					} else if data.botscript.ActiveRedWarBots == 0 {
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						ns.PrintStrToAll("Red Warrior Bot joined the game.")
						data.botscript.ActiveRedWarBots++
						data.botscript.ActiveRedWarBots++
						data.botscript.ActiveRedWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					}
				})
			}
		case "-1 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Conjurer Bot left the game.")
								data.botscript.ActiveRedConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Red Conjurer Bot active.")
					}
				})
			}
		case "-2 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Conjurer Bot left the game.")
								data.botscript.ActiveRedConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Conjurer Bot active.")
					}
				})
			}
		case "-3 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Conjurer Bot left the game.")
								data.botscript.ActiveRedConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Conjurer Bot active.")
					}
				})
			}
		case "+1 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewConjurer(Red))
					}
				})
			}
		case "+2 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedConBots == 3 {
						ns.PrintStrToAll("Red Conjurer Bot limit reached.")
					} else if data.botscript.ActiveRedConBots == 2 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot limit reached.")
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
					} else if data.botscript.ActiveRedConBots == 1 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						data.botscript.ActiveRedConBots++
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					} else if data.botscript.ActiveRedConBots == 0 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						data.botscript.ActiveRedConBots++
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					}
				})
			}
		case "+3 red con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedConBots == 3 {
						ns.PrintStrToAll("Red Conjurer Bot limit reached.")
					} else if data.botscript.ActiveRedConBots == 2 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot limit reached.")
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
					} else if data.botscript.ActiveRedConBots == 1 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot limit reached.")
						data.botscript.ActiveRedConBots++
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					} else if data.botscript.ActiveRedConBots == 0 {
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						ns.PrintStrToAll("Red Conjurer Bot joined the game.")
						data.botscript.ActiveRedConBots++
						data.botscript.ActiveRedConBots++
						data.botscript.ActiveRedConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					}
				})
			}
		case "-1 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Wizard Bot left the game.")
								data.botscript.ActiveRedWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Red Wizard Bot active.")
					}
				})
			}
		case "-2 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Wizard Bot left the game.")
								data.botscript.ActiveRedWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Wizard Bot active.")
					}
				})
			}
		case "-3 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[0]) {
								ns.PrintStrToAll("Red Wizard Bot left the game.")
								data.botscript.ActiveRedWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Red Wizard Bot active.")
					}
				})
			}
		case "+1 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewWizard(Red))
					}
				})
			}
		case "+2 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWizBots == 3 {
						ns.PrintStrToAll("Red Wizard Bot limit reached.")
					} else if data.botscript.ActiveRedWizBots == 2 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot limit reached.")
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
					} else if data.botscript.ActiveRedWizBots == 1 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						data.botscript.ActiveRedWizBots++
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					} else if data.botscript.ActiveRedWizBots == 0 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						data.botscript.ActiveRedWizBots++
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					}
				})
			}
		case "+3 red wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveRedWizBots == 3 {
						ns.PrintStrToAll("Red Wizard Bot limit reached.")
					} else if data.botscript.ActiveRedWizBots == 2 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot limit reached.")
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
					} else if data.botscript.ActiveRedWizBots == 1 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot limit reached.")
						data.botscript.ActiveRedWizBots++
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					} else if data.botscript.ActiveRedWizBots == 0 {
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						ns.PrintStrToAll("Red Wizard Bot joined the game.")
						data.botscript.ActiveRedWizBots++
						data.botscript.ActiveRedWizBots++
						data.botscript.ActiveRedWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					}
				})
			}
		case "-1 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Warrior Bot left the game.")
								data.botscript.ActiveBlueWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Warrior Bot active.")
					}
				})
			}
		case "-2 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Warrior Bot left the game.")
								data.botscript.ActiveBlueWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Warrior Bot active.")
					}
				})
			}
		case "-3 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWarBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 150 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Warrior Bot left the game.")
								data.botscript.ActiveBlueWarBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Warrior Bot active.")
					}
				})
			}
		case "+1 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewWarrior(Blue))
					}
				})
			}
		case "+2 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWarBots == 3 {
						ns.PrintStrToAll("Blue Warrior Bot limit reached.")
					} else if data.botscript.ActiveBlueWarBots == 2 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot limit reached.")
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
					} else if data.botscript.ActiveBlueWarBots == 1 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						data.botscript.ActiveBlueWarBots++
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					} else if data.botscript.ActiveBlueWarBots == 0 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						data.botscript.ActiveBlueWarBots++
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					}
				})
			}
		case "+3 blue war":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWarBots == 3 {
						ns.PrintStrToAll("Blue Warrior Bot limit reached.")
					} else if data.botscript.ActiveBlueWarBots == 2 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot limit reached.")
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
					} else if data.botscript.ActiveBlueWarBots == 1 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot limit reached.")
						data.botscript.ActiveBlueWarBots++
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					} else if data.botscript.ActiveBlueWarBots == 0 {
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						ns.PrintStrToAll("Blue Warrior Bot joined the game.")
						data.botscript.ActiveBlueWarBots++
						data.botscript.ActiveBlueWarBots++
						data.botscript.ActiveBlueWarBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					}
				})
			}
		case "-1 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Conjurer Bot left the game.")
								data.botscript.ActiveBlueConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Conjurer Bot active.")
					}
				})
			}
		case "-2 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Conjurer Bot left the game.")
								data.botscript.ActiveBlueConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Conjurer Bot active.")
					}
				})
			}
		case "-3 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueConBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 100 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Conjurer Bot left the game.")
								data.botscript.ActiveBlueConBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Conjurer Bot active.")
					}
				})
			}
		case "+1 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewConjurer(Blue))
					}
				})
			}
		case "+2 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueConBots == 3 {
						ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
					} else if data.botscript.ActiveBlueConBots == 2 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
					} else if data.botscript.ActiveBlueConBots == 1 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						data.botscript.ActiveBlueConBots++
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					} else if data.botscript.ActiveBlueConBots == 0 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						data.botscript.ActiveBlueConBots++
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					}
				})
			}
		case "+3 blue con":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueConBots == 3 {
						ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
					} else if data.botscript.ActiveBlueConBots == 2 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
					} else if data.botscript.ActiveBlueConBots == 1 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot limit reached.")
						data.botscript.ActiveBlueConBots++
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					} else if data.botscript.ActiveBlueConBots == 0 {
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						ns.PrintStrToAll("Blue Conjurer Bot joined the game.")
						data.botscript.ActiveBlueConBots++
						data.botscript.ActiveBlueConBots++
						data.botscript.ActiveBlueConBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					}
				})
			}
		case "-1 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Wizard Bot left the game.")
								data.botscript.ActiveBlueWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								return
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Wizard Bot active.")
					}
				})
			}
		case "-2 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Wizard Bot left the game.")
								data.botscript.ActiveBlueWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 2 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Wizard Bot active.")
					}
				})
			}
		case "-3 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWizBots > 0 {
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						x := 0
						delete := ns.FindAllObjects(ns.HasTypeName{"NPC"})
						for i := 0; i < len(delete); i++ {
							if delete[i].MaxHealth() == 75 && delete[i].IsEnabled() && delete[i].HasTeam(ns.Teams()[1]) {
								ns.PrintStrToAll("Blue Wizard Bot left the game.")
								data.botscript.ActiveBlueWizBots--
								delete[i].Enable(false)
								delete[i].SetPos(ns.Ptf(150, 150))
								x++
								if x == 3 {
									return
								}
							}
						}
					} else {
						ns.PrintStrToAll("No Blue Wizard Bot active.")
					}
				})
			}
		case "+1 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
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
						bots = append(bots, NewWizard(Blue))
					}
				})
			}
		case "+2 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWizBots == 3 {
						ns.PrintStrToAll("Blue Wizard Bot limit reached.")
					} else if data.botscript.ActiveBlueWizBots == 2 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot limit reached.")
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
					} else if data.botscript.ActiveBlueWizBots == 1 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						data.botscript.ActiveBlueWizBots++
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					} else if data.botscript.ActiveBlueWizBots == 0 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						data.botscript.ActiveBlueWizBots++
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					}
				})
			}
		case "+3 blue wiz":
			if !TeamsEnabled {
				ns.PrintStrToAll("Error: teams are disabled.")
			} else {
				updateMyBotScriptData(p, func(data *MyAccountData) {
					if data.botscript.ActiveBlueWizBots == 3 {
						ns.PrintStrToAll("Blue Wizard Bot limit reached.")
					} else if data.botscript.ActiveBlueWizBots == 2 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot limit reached.")
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
					} else if data.botscript.ActiveBlueWizBots == 1 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot limit reached.")
						data.botscript.ActiveBlueWizBots++
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					} else if data.botscript.ActiveBlueWizBots == 0 {
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						ns.PrintStrToAll("Blue Wizard Bot joined the game.")
						data.botscript.ActiveBlueWizBots++
						data.botscript.ActiveBlueWizBots++
						data.botscript.ActiveBlueWizBots++
						serverSettingSoundToAllPlayers := ns.Players()
						for i := 0; i < len(serverSettingSoundToAllPlayers); i++ {
							ns.AudioEvent(audio.ServerOptionsChange, serverSettingSoundToAllPlayers[i].Unit())
						}
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					}
				})
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
	checkSolo()
	checkTeams()
	//chatmapbotupdater.EachFrame(30, onSocialBots)
}

func onSocialBots() {
	if !GameModeIsSocial {
		return
	} else {
		Flags := ns.FindAllObjects(ns.HasTypeName{"Flag"})
		RedFlag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[0]})
		BlueFlag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[1]})
		if !Flags[0].IsEnabled() {

		} else {
			if Flags[0].HasTeam(ns.Teams()[1]) {
				ns.PrintStrToAll("Team 1")
			} else if Flags[0].HasTeam(ns.Teams()[0]) {
				ns.PrintStrToAll("Team 0")
			} else {
				ns.PrintStrToAll("no team")
			}
			social := ns.FindAllObjects(ns.HasTypeName{"NPC"})
			if social != nil {
				if TeamsEnabled {
					for i := 0; i < len(social); i++ {
						arr := ns.FindClosestObject(social[i], ns.HasTypeName{"NewPlayer"})
						social[i].WalkTo(arr.Pos())
					}
				}
			}
		}

	}
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
