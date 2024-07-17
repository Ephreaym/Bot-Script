package BotWars

import "github.com/noxworld-dev/noxscript/ns/v4"

// Bot interface contains functions common to all bots.
type Bot interface {
	// Update executes bot logic.
	Update()
}

var bots []Bot // bots array; can contain any number of bots.

func init() {
	ns.NewTimer(ns.Frames(10), func() {
		if !TeamsEnabled {
			if ns.Object("WarriorOwner") != nil && ns.Object("ConjurerOwner") != nil && ns.Object("WizardOwner") != nil {
				updateMyBotScriptData(ns.GetHost().Player(), func(data *MyAccountData) {
					if data.botscript.ActiveWarBots > 0 {
						bots = append(bots, NewWarriorNoTeam())
					}
					if data.botscript.ActiveConBots > 0 {
						bots = append(bots, NewConjurerNoTeam())
					}
					if data.botscript.ActiveWizBots > 0 {
						bots = append(bots, NewWizardNoTeam())
					}
				})
			}
			return
		}
		ns.NewTimer(ns.Frames(60), func() {
			if TeamsEnabled {
				updateMyBotScriptData(ns.GetHost().Player(), func(data *MyAccountData) {
					// Add this many bots on map launch.
					// Team Red
					switch data.botscript.ActiveRedWarBots {
					case 0:
					case 1:
						bots = append(bots, NewWarrior(Red))
					case 2:
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					case 3:
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
						bots = append(bots, NewWarrior(Red))
					}
					switch data.botscript.ActiveRedConBots {
					case 0:
					case 1:
						bots = append(bots, NewConjurer(Red))
					case 2:
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					case 3:
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
						bots = append(bots, NewConjurer(Red))
					}
					switch data.botscript.ActiveRedWizBots {
					case 0:
					case 1:
						bots = append(bots, NewWizard(Red))
					case 2:
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					case 3:
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
						bots = append(bots, NewWizard(Red))
					}
					// Team Blue
					switch data.botscript.ActiveBlueWarBots {
					case 0:
					case 1:
						bots = append(bots, NewWarrior(Blue))
					case 2:
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					case 3:
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
						bots = append(bots, NewWarrior(Blue))
					}
					switch data.botscript.ActiveBlueConBots {
					case 0:
					case 1:
						bots = append(bots, NewConjurer(Blue))
					case 2:
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					case 3:
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
						bots = append(bots, NewConjurer(Blue))
					}
					switch data.botscript.ActiveBlueWizBots {
					case 0:
					case 1:
						bots = append(bots, NewWizard(Blue))
					case 2:
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					case 3:
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
						bots = append(bots, NewWizard(Blue))
					}
				})
			}
		})
	})
}

// UpdateBots is called each frame to execute bot logic.
// It will range over the bot array.
func UpdateBots() {
	for _, bot := range bots {
		if !GameModeIsSocial {
			bot.Update()
		}
	}
}
