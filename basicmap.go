package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

func init() {
	ns.Music(15, 20)
}

// OnFrame is called by the server.
func OnFrame() {
	UpdateBots()
}

func DialogOptions() {
	// Usable dialog bits
	// F1GD401E "What a Wizard spy?"

	// C2NC203E "Get away from me filthy peasants"
	// C2NC202E "URHGHH"

	// C5OGK02E "too bad you must die now"
	// C5OGK01E "youre very bold for such a little man"
	// C5OGK05E "Ill crush your bones"
}
