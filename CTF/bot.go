package EndGameBW

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
			redwarriors   = 1
			redwizards    = 1
			redconjurers  = 1
			blueconjurers = 1
			bluewizards   = 1
			bluewarriors  = 1
		)
		ns.NewTimer(ns.Frames(1), func() {
			for i := 0; i < redwarriors; i++ {
				bots = append(bots, NewRedWarrior())
			}
		})
		ns.NewTimer(ns.Frames(2), func() {
			for i := 0; i < redwizards; i++ {
				bots = append(bots, NewRedWizard())
			}
		})
		ns.NewTimer(ns.Frames(3), func() {
			for i := 0; i < redconjurers; i++ {
				bots = append(bots, NewRedConjurer())
			}
		})
		ns.NewTimer(ns.Frames(4), func() {
			for i := 0; i < bluewarriors; i++ {
				bots = append(bots, NewBlueWarrior())
			}
		})
		ns.NewTimer(ns.Frames(5), func() {
			for i := 0; i < bluewizards; i++ {
				bots = append(bots, NewBlueWizard())
			}
		})
		ns.NewTimer(ns.Frames(6), func() {
			for i := 0; i < blueconjurers; i++ {
				bots = append(bots, NewBlueConjurer())

			}
		})
	})
}

// UpdateBots is called each frame to execute bot logic.
// It will range over the bot array.
func UpdateBots() {
	for _, bot := range bots {
		bot.Update()
	}
}
