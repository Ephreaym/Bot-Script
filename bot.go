package basicmap

// Bot interface contains functions common to all bots.
type Bot interface {
	// Update executes bot logic.
	Update()
}

var bots []Bot // bots array; can contain any number of bots.

func init() {

	// Add this many bots on map launch.
	const (
		warriors  = 1
		wizards   = 1
		conjurers = 1
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
}

// UpdateBots is called each frame to execute bot logic.
// It will range over the bot array.
func UpdateBots() {
	for _, bot := range bots {
		bot.Update()
	}
}
