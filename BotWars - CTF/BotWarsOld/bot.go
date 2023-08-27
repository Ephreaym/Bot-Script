package BotWars

import "github.com/noxworld-dev/noxscript/ns/v4"

// Bot interface contains functions common to all bots.
type Bot interface {
	// Update executes bot logic.
	Update()
}

var bots []Bot // bots array; can contain any number of bots.

func init() {
	ns.NewTimer(ns.Frames(60), func() {
		// Add this many bots on map launch.
		const (
			warriors     = 1
			wizards      = 0
			conjurers    = 0
			bluewarriors = 1
		)
		for i := 0; i < warriors; i++ {
			bots = append(bots, NewWarrior())
		}
		for i := 0; i < wizards; i++ {
			bots = append(bots, NewWizard())
		}
		for i := 0; i < conjurers; i++ {
			bots = append(bots, NewConjurer())
		}
		for i := 0; i < bluewarriors; i++ {
			bots = append(bots, NewBlueWarrior())
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
