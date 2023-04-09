package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWizard creates a new wizard bot.
func NewWizard() *Wizard {
	wiz := &Wizard{}
	wiz.init()
	return wiz
}

// Wizard bot class.
type Wizard struct {
	unit         ns.Obj
	target       ns.Obj
	taggedPlayer ns.Obj
	trap         ns.Obj
	items        struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
		WizardRobe     ns.Obj
	}
	spells struct {
		Ready              bool // Duration unknown.
		MagicMissilesReady bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		ForceFieldReady    bool // Duration unknown.
		ShockReady         bool // No real cooldown,
		SlowReady          bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		TrapReady          bool // Only one trap is placed per life.
	}
}

func (wiz *Wizard) init() {
	// Reset spells
	wiz.spells.Ready = true
	wiz.spells.MagicMissilesReady = true
	wiz.spells.ForceFieldReady = true
	wiz.spells.ShockReady = true
	wiz.spells.SlowReady = true
	wiz.spells.TrapReady = true

	wiz.unit = ns.CreateObject("NPC", ns.GetHost())
	wiz.items.StreetSneakers = ns.CreateObject("StreetSneakers", wiz.unit)
	wiz.items.StreetPants = ns.CreateObject("StreetPants", wiz.unit)
	wiz.items.StreetShirt = ns.CreateObject("StreetShirt", wiz.unit)
	wiz.items.WizardRobe = ns.CreateObject("WizardRobe", wiz.unit)
	wiz.unit.Equip(wiz.items.StreetSneakers)
	wiz.unit.Equip(wiz.items.StreetPants)
	wiz.unit.Equip(wiz.items.StreetShirt)
	wiz.unit.Equip(wiz.items.WizardRobe)
	wiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	wiz.unit.SetMaxHealth(75)
	wiz.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		wiz.unit.AggressionLevel(0.83)
	})
	wiz.unit.Hunt()
	wiz.unit.ResumeLevel(0.8)
	wiz.unit.RetreatLevel(0.2)
	// Buff on respawn. //
	wiz.buffInitial()
	// When an enemy is seen. //
	wiz.unit.OnEvent(ns.EventEnemySighted, wiz.onEnemySighted)
	// On collision.
	wiz.unit.OnEvent(ns.EventCollision, wiz.onCollide)
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight.
	wiz.unit.OnEvent(ns.EventLostEnemy, wiz.onLostEnemy)
	// On Death.
	wiz.unit.OnEvent(ns.EventDeath, wiz.onDeath)
}

func (wiz *Wizard) buffInitial() {
	if wiz.spells.ForceFieldReady {
		wiz.spells.ForceFieldReady = true
		wiz.spells.Ready = true
		// Force Field chant.
		castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
			ns.CastSpell(spell.SHIELD, wiz.unit, wiz.unit)
			wiz.spells.Ready = true
			// Pause for concentration.
			ns.NewTimer(ns.Frames(3), func() {
				// Haste chant.
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
					ns.CastSpell(spell.HASTE, wiz.unit, wiz.unit)
				})
			})
		})
		ns.NewTimer(ns.Seconds(30), func() {
			wiz.spells.ForceFieldReady = true
		})
	}
}
func (wiz *Wizard) onEnemySighted() {
	if wiz.spells.SlowReady {
		wiz.spells.SlowReady = false
		wiz.spells.Ready = false
		// Slow chant.
		castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
			wiz.spells.Ready = true
			wiz.target = ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassPlayer))
			ns.CastSpell(spell.SLOW, wiz.unit, wiz.target)
			ns.NewTimer(ns.Seconds(5), func() {
				wiz.spells.SlowReady = true
			})
		})
	}
}
func (wiz *Wizard) onCollide() {
	if wiz.spells.ShockReady {
		wiz.spells.Ready = false
		wiz.spells.ShockReady = false
		// Shock chant.
		castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
			ns.CastSpell(spell.SHOCK, wiz.unit, wiz.unit)
			wiz.spells.Ready = true
			ns.NewTimer(ns.Seconds(10), func() {
				wiz.spells.ShockReady = true
			})
		})
	}
}

func (wiz *Wizard) onLostEnemy() {
	if wiz.spells.TrapReady {
		wiz.spells.Ready = true
		// WizAITrapCooldown = 0
		// Ring of Fire chant.
		castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhDown, PhDownLeft, PhUp}, func() {
			// Pause of Glyph concentration.
			ns.NewTimer(ns.Frames(3), func() {
				// Magic Missiles chant.
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
					// Pause of Glyph concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Shock chant.
						castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
							// Pause of Glyph concentration.
							ns.NewTimer(ns.Frames(3), func() {
								// Glyph chant.
								castPhonemes(wiz.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
									ns.AudioEvent(audio.TrapDrop, wiz.unit)
									wiz.trap = ns.NewTrap(wiz.unit, spell.CLEANSING_FLAME, spell.MAGIC_MISSILE, spell.SHOCK)
									wiz.trap.SetOwner(wiz.unit)
									wiz.spells.Ready = true
								})
							})
						})
					})
				})
			})
		})
	}
}

func (wiz *Wizard) onDeath() {
	wiz.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, wiz.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, wiz.unit)
		wiz.unit.Delete()
		wiz.items.StreetPants.Delete()
		wiz.items.StreetSneakers.Delete()
		wiz.items.StreetShirt.Delete()
		wiz.init()
	})
}

func (wiz *Wizard) Update() {
	wiz.findLoot()
}

func (wiz *Wizard) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			//"DeathRayWand",
			"FireStormWand",
			//"LesserFireballWand",
			"ForceWand",
			//"SulphorousShowerWand"
			//"SulphorousFlareWand"
			//"StaffWooden",
		},
	)
	for _, item := range weapons {
		wiz.unit.Equip(item)
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Wizard armor.
			"WizardHelm", "WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		wiz.unit.Equip(item)
	}
}
