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
		Ready               bool // Duration unknown.
		DeathRayReady       bool // No cooldown, 60 mana.
		MagicMissilesReady  bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		ForceFieldReady     bool // Duration unknown.
		ShockReady          bool // Duration is 20 seconds.
		SlowReady           bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		TrapReady           bool // Only one trap is placed per life.
		EnergyBoltReady     bool // No real cooldown, mana cost *.
		FireballReady       bool // No real cooldown, mana cost 30.
		ProtFromFireReady   bool // Duration is 60 seconds.
		ProtFromPoisonReady bool
		ProtFromShockReady  bool
		BlinkReady          bool
		HasteReady          bool // Duration is 20 seconds
		InvisibilityReady   bool // Duration is 60 seconds, 30 mana.
	}
	reactionTime int
}

func (wiz *Wizard) init() {
	// Reset spells WizBot.
	wiz.spells.Ready = true
	// Debuff spells.
	wiz.spells.SlowReady = true
	// Offensive spells.
	wiz.spells.MagicMissilesReady = true
	wiz.spells.TrapReady = true
	wiz.spells.DeathRayReady = true
	wiz.spells.EnergyBoltReady = true
	wiz.spells.FireballReady = true
	// Defensive spells.
	wiz.spells.BlinkReady = true
	// Buff spells
	wiz.spells.ShockReady = true
	wiz.spells.ProtFromFireReady = true
	wiz.spells.ProtFromPoisonReady = true
	wiz.spells.ProtFromShockReady = true
	wiz.spells.HasteReady = true
	wiz.spells.ForceFieldReady = true
	wiz.spells.InvisibilityReady = true
	// Create WizBot.
	wiz.unit = ns.CreateObject("NPC", RandomBotSpawn)
	wiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	wiz.unit.SetMaxHealth(75)
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	wiz.reactionTime = 15
	// Set WizBot properties.
	wiz.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	wiz.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	wiz.unit.MonsterStatusEnable(object.MonStatusAlert)
	wiz.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		wiz.unit.AggressionLevel(0.83)
	})
	wiz.unit.Hunt()
	wiz.unit.ResumeLevel(0.8)
	wiz.unit.RetreatLevel(0.2)
	// Create and equip WizBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	wiz.items.StreetSneakers = ns.CreateObject("StreetSneakers", wiz.unit)
	wiz.items.StreetPants = ns.CreateObject("StreetPants", wiz.unit)
	wiz.items.StreetShirt = ns.CreateObject("StreetShirt", wiz.unit)
	wiz.items.WizardRobe = ns.CreateObject("WizardRobe", wiz.unit)
	wiz.unit.Equip(wiz.items.StreetSneakers)
	wiz.unit.Equip(wiz.items.StreetPants)
	wiz.unit.Equip(wiz.items.StreetShirt)
	wiz.unit.Equip(wiz.items.WizardRobe)
	// Buff on respawn.
	wiz.buffInitial()
	// On retreat.
	wiz.unit.OnEvent(ns.EventRetreat, wiz.onRetreat)
	// Enemy sighted.
	wiz.unit.OnEvent(ns.EventEnemySighted, wiz.onEnemySighted)
	// On heard.
	wiz.unit.OnEvent(ns.EventEnemyHeard, wiz.onEnemyHeard)
	// On collision.
	wiz.unit.OnEvent(ns.EventCollision, wiz.onCollide)
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight.
	wiz.unit.OnEvent(ns.EventLostEnemy, wiz.onLostEnemy)
	// On Death.
	wiz.unit.OnEvent(ns.EventDeath, wiz.onDeath)
}

func (wiz *Wizard) buffInitial() {
	// Cast Force Field.
	if wiz.spells.ForceFieldReady && wiz.spells.Ready {
		wiz.spells.ForceFieldReady = false
		wiz.spells.HasteReady = false
		wiz.spells.Ready = false
		castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
			ns.CastSpell(spell.SHIELD, wiz.unit, wiz.unit)
			// Pause for concentration.
			ns.NewTimer(ns.Frames(3), func() {
				// Cast Haste.
				wiz.spells.Ready = true
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
					ns.CastSpell(spell.HASTE, wiz.unit, wiz.unit)
					wiz.spells.HasteReady = true
				})
			})
		})
		ns.NewTimer(ns.Seconds(60), func() {
			wiz.spells.ForceFieldReady = true
		})
	}
}

func (wiz *Wizard) onRetreat() {
	// Cast Blink.
	if wiz.spells.BlinkReady && wiz.spells.Ready {
		wiz.spells.BlinkReady = false
		wiz.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
				ns.NewTrap(wiz.unit, spell.BLINK)
				ns.CastSpell(spell.TRIGGER_GLYPH, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
			})
			ns.NewTimer(ns.Seconds(2), func() {
				wiz.spells.BlinkReady = true
			})
		})
	}
}

func (wiz *Wizard) onEnemySighted() {
	// Cast Slow.
	if wiz.spells.SlowReady && wiz.spells.Ready {
		wiz.spells.SlowReady = false
		wiz.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
				wiz.spells.Ready = true
				wiz.target = ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassPlayer))
				wiz.unit.LookAtObject(wiz.target)
				ns.CastSpell(spell.SLOW, wiz.unit, wiz.target)
				ns.NewTimer(ns.Seconds(5), func() {
					wiz.spells.SlowReady = true
				})
			})
		})
	}
	// Cast Death Ray.
	if wiz.spells.DeathRayReady && wiz.spells.Ready {
		wiz.spells.DeathRayReady = false
		wiz.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhDownRight}, func() {
				wiz.spells.Ready = true
				wiz.target = ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassPlayer))
				wiz.unit.LookAtObject(wiz.target)
				ns.CastSpell(spell.DEATH_RAY, wiz.unit, wiz.target)
				ns.NewTimer(ns.Seconds(10), func() {
					wiz.spells.DeathRayReady = true
				})
			})
		})
	}
}

func (wiz *Wizard) onEnemyHeard() {
	// Cast Inivsibility.
	if wiz.spells.InvisibilityReady && wiz.spells.Ready {
		wiz.spells.InvisibilityReady = false
		wiz.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
				ns.CastSpell(spell.INVISIBILITY, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					wiz.spells.InvisibilityReady = true
				})
			})
		})
	}
}

func (wiz *Wizard) onCollide() {
	// Cast Shock.
	if wiz.spells.ShockReady && wiz.spells.Ready {
		wiz.spells.Ready = false
		wiz.spells.ShockReady = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
				ns.CastSpell(spell.SHOCK, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(20), func() {
					wiz.spells.ShockReady = true
				})
			})
		})
	}
	wiz.target = ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassPlayer))
	if wiz.unit.CanSee(wiz.target) && wiz.spells.MagicMissilesReady && wiz.spells.Ready {
		// Cast Missiles of Magic.
		wiz.spells.Ready = false
		wiz.spells.MagicMissilesReady = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
				wiz.unit.LookAtObject(wiz.target)
				ns.CastSpell(spell.MAGIC_MISSILE, wiz.unit, wiz.target)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(5), func() {
					wiz.spells.MagicMissilesReady = true
				})
			})
		})
	}
}

func (wiz *Wizard) onLostEnemy() {
	// Cast trap (Ring of Fire, Missiles of Magic, Shock).
	if wiz.spells.TrapReady && wiz.spells.Ready {
		wiz.spells.Ready = false
		wiz.spells.TrapReady = false
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
	if wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
		wiz.spells.Ready = false
		ns.NewTimer(ns.Seconds(3), func() {
			wiz.spells.Ready = true
		})
	}
	// Offensive logic.
	wiz.target = ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassPlayer))
	if wiz.unit.CanSee(wiz.target) && wiz.spells.EnergyBoltReady && wiz.spells.Ready {
		wiz.spells.Ready = false
		wiz.spells.EnergyBoltReady = false
		ns.NewTimer(ns.Frames(15), func() {
			// Energy Bolt chant.
			castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhUp}, func() {
				wiz.unit.LookAtObject(wiz.target)
				ns.CastSpell(spell.LIGHTNING, wiz.unit, wiz.target)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(5), func() {
					wiz.spells.EnergyBoltReady = true
				})
			})
		})
	}
	if wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready {
		wiz.spells.Ready = false
		wiz.spells.FireballReady = false
		// Fireball chant.
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
				wiz.unit.LookAtObject(wiz.target)
				ns.CastSpell(spell.FIREBALL, wiz.unit, wiz.target)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(10), func() {
					wiz.spells.FireballReady = true
				})
			})
		})
	}
	// Buffing logic.
	if !wiz.unit.HasEnchant(enchant.HASTED) {
		if wiz.spells.HasteReady && wiz.spells.Ready {
			// Haste chant.
			wiz.spells.Ready = false
			wiz.spells.HasteReady = false
			castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
				wiz.spells.Ready = true
				ns.CastSpell(spell.HASTE, wiz.unit, wiz.unit)
				ns.NewTimer(ns.Seconds(20), func() {
					wiz.spells.HasteReady = true
				})
			})
		}
	}
	if !wiz.unit.HasEnchant(enchant.SHIELD) {
		if wiz.spells.ForceFieldReady && wiz.spells.Ready {
			wiz.spells.Ready = false
			wiz.spells.ForceFieldReady = false
			// Force Field chant.
			castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
				ns.CastSpell(spell.SHIELD, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					wiz.spells.ForceFieldReady = true
				})
			})
		}
	}
	if !wiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) {
		if wiz.spells.ProtFromFireReady && wiz.spells.Ready {
			wiz.spells.Ready = false
			wiz.spells.ProtFromFireReady = false
			// Protection from Fire chant.
			castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
				ns.CastSpell(spell.PROTECTION_FROM_FIRE, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					wiz.spells.ProtFromFireReady = true
				})
			})
		}
	}
	if !wiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) {
		if wiz.spells.ProtFromShockReady && wiz.spells.Ready {
			wiz.spells.Ready = false
			wiz.spells.ProtFromShockReady = false
			// Protection from Shock chant.
			castPhonemes(wiz.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
				ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, wiz.unit, wiz.unit)
				wiz.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					wiz.spells.ProtFromShockReady = true
				})
			})
		}
	}
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
