package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
)

// Shortcuts for phoneme audio effects.
const (
	PhUp        = audio.SpellPhonemeUp
	PhDown      = audio.SpellPhonemeDown
	PhLeft      = audio.NPCSpellPhonemeLeft
	PhRight     = audio.NPCSpellPhonemeRight
	PhUpLeft    = audio.NPCSpellPhonemeUpLeft
	PhUpRight   = audio.NPCSpellPhonemeUpRight
	PhDownLeft  = audio.NPCSpellPhonemeDownLeft
	PhDownRight = audio.NPCSpellPhonemeDownRight
)

// castPhonemes emulates the spell casting audio effect and then calls a given function.
func castPhonemes(pos ns.Positioner, phonemes []audio.Name, fnc func()) {
	if len(phonemes) == 0 {
		// no phonemes left to cast
		fnc()
		return
	}
	if ph := phonemes[0]; ph != "" {
		ns.AudioEvent(ph, pos)
	}
	ns.NewTimer(ns.Frames(3), func() {
		castPhonemes(pos, phonemes[1:], fnc)
	})
}
