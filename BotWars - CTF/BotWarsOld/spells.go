package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
)

// Shortcuts for phoneme audio effects.
const (
	// Male chant.
	PhUp        = audio.SpellPhonemeUp
	PhDown      = audio.SpellPhonemeDown
	PhLeft      = audio.NPCSpellPhonemeLeft
	PhRight     = audio.NPCSpellPhonemeRight
	PhUpLeft    = audio.NPCSpellPhonemeUpLeft
	PhUpRight   = audio.NPCSpellPhonemeUpRight
	PhDownLeft  = audio.NPCSpellPhonemeDownLeft
	PhDownRight = audio.NPCSpellPhonemeDownRight
	// Female chant.
	FPhUp        = audio.FemaleSpellPhonemeUp
	FPhDown      = audio.FemaleSpellPhonemeDown
	FPhLeft      = audio.FemaleSpellPhonemeLeft
	FPhRight     = audio.FemaleSpellPhonemeRight
	FPhUpLeft    = audio.FemaleSpellPhonemeUpLeft
	FPhUpRight   = audio.FemaleSpellPhonemeUpRight
	FPhDownLeft  = audio.FemaleSpellPhonemeDownLef
	FPhDownRight = audio.FemaleSpellPhonemeDownRig
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
