package BotWars

import "github.com/noxworld-dev/noxscript/ns/v4"

// Bot interface contains functions common to all bots.
type Bot interface {
	// Update executes bot logic.
	Update()
}

var bots []Bot // bots array; can contain any number of bots.

func init() {
	if !TeamsEnabled {
		return
	}
	ns.NewTimer(ns.Frames(60), func() {
		// Add this many bots on map launch.
		// 1 = active bot.
		// 0 = deactivated bot.
		const (
			// Team red.
			// Warriors.
			redwarriors01 = 0
			redwarriors02 = 0
			redwarriors03 = 0
			// Conjurers.
			redconjurers01 = 0
			redconjurers02 = 0
			redconjurers03 = 0
			// Wizards.
			redwizards01 = 0
			redwizards02 = 0
			redwizards03 = 0

			// Team Blue.
			// Warriors.
			bluewarriors01 = 0
			bluewarriors02 = 0
			bluewarriors03 = 0
			// Conjurers.
			blueconjurers01 = 0
			blueconjurers02 = 0
			blueconjurers03 = 0
			// Wizards.
			bluewizards01 = 0
			bluewizards02 = 0
			bluewizards03 = 0
		)

		for i := 0; i < redwarriors01; i++ {
			bots = append(bots, NewWarrior(Red))
		}
		for i := 0; i < redwizards01; i++ {
			bots = append(bots, NewWizard(Red))
		}
		for i := 0; i < redconjurers01; i++ {
			bots = append(bots, NewConjurer(Red))
		}
		for i := 0; i < bluewarriors01; i++ {
			bots = append(bots, NewWarrior(Blue))
		}
		for i := 0; i < bluewizards01; i++ {
			bots = append(bots, NewWizard(Blue))
		}
		for i := 0; i < blueconjurers01; i++ {
			bots = append(bots, NewConjurer(Blue))
		}
		for i := 0; i < redwarriors02; i++ {
			bots = append(bots, NewWarrior(Red))
		}
		for i := 0; i < redwizards02; i++ {
			bots = append(bots, NewWizard(Red))
		}
		for i := 0; i < redconjurers02; i++ {
			bots = append(bots, NewConjurer(Red))
		}
		for i := 0; i < bluewarriors02; i++ {
			bots = append(bots, NewWarrior(Blue))
		}
		for i := 0; i < bluewizards02; i++ {
			bots = append(bots, NewWizard(Blue))
		}
		for i := 0; i < blueconjurers02; i++ {
			bots = append(bots, NewConjurer(Blue))
		}
		for i := 0; i < redwarriors03; i++ {
			bots = append(bots, NewWarrior(Red))
		}
		for i := 0; i < redwizards03; i++ {
			bots = append(bots, NewWizard(Red))
		}
		for i := 0; i < redconjurers03; i++ {
			bots = append(bots, NewConjurer(Red))
		}
		for i := 0; i < bluewarriors03; i++ {
			bots = append(bots, NewWarrior(Blue))
		}
		for i := 0; i < bluewizards03; i++ {
			bots = append(bots, NewWizard(Blue))
		}
		for i := 0; i < blueconjurers03; i++ {
			bots = append(bots, NewConjurer(Blue))
		}
	})
}

// UpdateBots is called each frame to execute bot logic.
// It will range over the bot array.
func UpdateBots() {
	for _, bot := range bots {
		bot.Update()
	}
}
