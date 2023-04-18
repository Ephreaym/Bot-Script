package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewConjurer creates a new conjurer bot.
func NewConjurer() *Conjurer {
	con := &Conjurer{}
	con.init()
	return con
}

// Conjurer bot class.
type Conjurer struct {
	unit         ns.Obj
	target       ns.Obj
	taggedPlayer ns.Obj
	bomber       ns.Obj
	items        struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
	}
	spells struct {
		Ready                bool // Duration unknown.
		InfravisionReady     bool // Duration is 30 seconds.
		VampirismReady       bool // Duration is 30 seconds.
		BlinkReady           bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		FistOfVengeanceReady bool // No real cooldown, mana cost 60.
		StunReady            bool // No real cooldown.
		SummonBomberReady    bool // No real cooldown.
		SummonGhostReady     bool
		ProtFromFireReady    bool // Duration is 60 seconds.
		ProtFromPoisonReady  bool
		ProtFromShockReady   bool
		PixieSwarmReady      bool
	}
}

func (con *Conjurer) init() {
	// Reset spells.
	con.spells.Ready = true
	con.spells.StunReady = true
	con.spells.InfravisionReady = true
	con.spells.VampirismReady = true
	con.spells.BlinkReady = true
	con.spells.FistOfVengeanceReady = true
	con.spells.PixieSwarmReady = true
	// Summons.
	con.spells.SummonGhostReady = true
	con.spells.SummonBomberReady = true
	// Buff spells.
	con.spells.ProtFromFireReady = true
	con.spells.ProtFromPoisonReady = true
	con.spells.ProtFromShockReady = true

	// Create Conjurer bot.
	con.unit = ns.CreateObject("NPC", RandomBotSpawn)
	con.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	con.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	con.unit.MonsterStatusEnable(object.MonStatusAlert)
	con.items.StreetSneakers = ns.CreateObject("StreetSneakers", con.unit)
	con.items.StreetPants = ns.CreateObject("StreetPants", con.unit)
	con.items.StreetShirt = ns.CreateObject("StreetShirt", con.unit)
	con.unit.Equip(con.items.StreetPants)
	con.unit.Equip(con.items.StreetShirt)
	con.unit.Equip(con.items.StreetSneakers)
	con.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	con.unit.SetMaxHealth(100)
	con.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		con.unit.AggressionLevel(0.83)
	})
	con.unit.Hunt()
	con.unit.ResumeLevel(0.8)
	con.unit.RetreatLevel(0.2)
	// Buff on respawn.
	con.buffInitial()
	// When an enemy is seen.
	con.unit.OnEvent(ns.EventEnemySighted, con.onEnemySighted)
	// When touched by another object.
	con.unit.OnEvent(ns.EventCollision, con.onCollide)
	// Escape.
	con.unit.OnEvent(ns.EventRetreat, con.onRetreat)
	// Enemy Lost.
	con.unit.OnEvent(ns.EventLostEnemy, con.onLostEnemy)
	// On Death.
	con.unit.OnEvent(ns.EventDeath, con.onDeath)
	// Looking for a target.
	con.unit.OnEvent(ns.EventLookingForEnemy, con.onLookingForTarget)
	//con.unit.OnEvent(ns.EventChangeFocus, con.onChangeFocus)
}

func (con *Conjurer) onLookingForTarget() {

}

func (con *Conjurer) onCollide() {

}

func (con *Conjurer) buffInitial() {
	if con.spells.VampirismReady && con.spells.Ready {
		con.spells.VampirismReady = false
		con.spells.Ready = false
		castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
			ns.CastSpell(spell.VAMPIRISM, con.unit, con.unit)
			con.spells.Ready = true
		})
		ns.NewTimer(ns.Seconds(30), func() {
			con.spells.VampirismReady = true
		})
	}
}

func (con *Conjurer) onEnemySighted() {
	if con.spells.StunReady && con.spells.Ready {
		con.spells.StunReady = false
		con.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			// Stun chant.
			castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDown}, func() {
				con.spells.Ready = true
				con.target = ns.FindClosestObject(con.unit, ns.HasClass(object.ClassPlayer))
				con.unit.LookAtObject(con.target)
				ns.CastSpell(spell.STUN, con.unit, con.target)
				ns.NewTimer(ns.Seconds(5), func() {
					con.spells.StunReady = true
				})
			})
		})
	}
}

func (con *Conjurer) onRetreat() {
	// Cast blink when retreating. TODO: fix trap workaround.
	if con.spells.BlinkReady && con.spells.Ready {
		con.spells.BlinkReady = false
		con.spells.Ready = false
		ns.NewTimer(ns.Frames(15), func() {
			castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
				ns.NewTrap(con.unit, spell.BLINK) // TODO: FIX IT so it doesn't have to be a trap.
				con.spells.Ready = true
			})
			ns.NewTimer(ns.Seconds(2), func() {
				con.spells.BlinkReady = true
			})
		})
	}
}

func (con *Conjurer) onLostEnemy() {
	// Cast Infravision when an enemy is lost.
	if con.spells.InfravisionReady && con.spells.Ready {
		con.spells.InfravisionReady = false
		con.spells.Ready = false
		castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhRight, PhLeft}, func() {
			con.unit.Enchant(enchant.INFRAVISION, ns.Seconds(30))
			con.spells.Ready = true
		})
		ns.NewTimer(ns.Seconds(30), func() {
			con.spells.InfravisionReady = true
		})
	}
	// Summon Ghost.
	// Set maximum of summons based on size.
	if con.spells.SummonGhostReady && con.spells.Ready {
		con.spells.SummonGhostReady = false
		con.spells.Ready = false
		castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown}, func() {
			ns.CastSpell(spell.SUMMON_GHOST, con.unit, con.unit)
			con.spells.Ready = true
			ns.NewTimer(ns.Seconds(10), func() {
				con.spells.SummonGhostReady = true
			})
		})
	}
	// Summon Bomber.
	// TODO: Set maximum of 2 active.
	if con.spells.SummonBomberReady && con.spells.Ready {
		con.spells.SummonBomberReady = false
		con.spells.Ready = false
		castPhonemes(con.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
			con.bomber = ns.CreateObject("Bomber", con.unit)
			con.bomber.SetOwner(con.unit)
			con.bomber.TrapSpells(spell.POISON, spell.METEOR, spell.STUN)
			con.spells.Ready = true
		})
		ns.NewTimer(ns.Seconds(10), func() {
			con.spells.SummonBomberReady = true
		})
	}
}

func (con *Conjurer) onDeath() {
	con.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, con.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, con.unit)
		con.unit.Delete()
		con.items.StreetPants.Delete()
		con.items.StreetShirt.Delete()
		con.items.StreetSneakers.Delete()
		con.init()
	})
}

func (con *Conjurer) Update() {
	con.findLoot()
	con.target = ns.FindClosestObject(con.unit, ns.HasClass(object.ClassPlayer))
	if con.unit.HasEnchant(enchant.ANTI_MAGIC) {
		con.spells.Ready = false
		ns.NewTimer(ns.Seconds(3), func() {
			con.spells.Ready = true
		})
	}
	if con.unit.HasEnchant(enchant.HELD) {
		if con.spells.BlinkReady && con.spells.Ready {
			con.spells.BlinkReady = false
			con.spells.Ready = false
			ns.NewTimer(ns.Frames(15), func() {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					ns.NewTrap(con.unit, spell.BLINK) // TODO: FIX IT so it doesn't have to be a trap.
					con.spells.Ready = true
				})
				ns.NewTimer(ns.Seconds(2), func() {
					con.spells.BlinkReady = true
				})
			})
		}
	}
	if con.target.HasEnchant(enchant.HELD) {
		ns.NewTimer(ns.Frames(15), con.FistOfVengeance)
	}
	// Buffing logic.
	if !con.unit.HasEnchant(enchant.VAMPIRISM) {
		if con.spells.VampirismReady && con.spells.Ready {
			con.spells.VampirismReady = false
			con.spells.Ready = false
			castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
				ns.CastSpell(spell.VAMPIRISM, con.unit, con.unit)
				con.spells.Ready = true
			})
			ns.NewTimer(ns.Seconds(30), func() {
				con.spells.VampirismReady = true
			})
		}
	}
	if !con.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) {
		if con.spells.ProtFromFireReady && con.spells.Ready {
			con.spells.Ready = false
			con.spells.ProtFromFireReady = false
			// Protection from Fire chant.
			castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
				ns.CastSpell(spell.PROTECTION_FROM_FIRE, con.unit, con.unit)
				con.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					con.spells.ProtFromFireReady = true
				})
			})
		}
	}
	if !con.unit.HasEnchant(enchant.PROTECT_FROM_POISON) {
		if con.spells.ProtFromPoisonReady && con.spells.Ready {
			con.spells.Ready = false
			con.spells.ProtFromPoisonReady = false
			// Protection from Poison chant.
			castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownLeft, PhUpRight}, func() {
				ns.CastSpell(spell.PROTECTION_FROM_POISON, con.unit, con.unit)
				con.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					con.spells.ProtFromPoisonReady = true
				})
			})
		}
	}
	if !con.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) {
		if con.spells.ProtFromShockReady && con.spells.Ready {
			con.spells.Ready = false
			con.spells.ProtFromShockReady = false
			// Protection from Shock chant.
			castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
				ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, con.unit, con.unit)
				con.spells.Ready = true
				ns.NewTimer(ns.Seconds(60), func() {
					con.spells.ProtFromShockReady = true
				})
			})
		}
	}
	if con.spells.PixieSwarmReady && con.spells.Ready {
		con.spells.Ready = false
		con.spells.PixieSwarmReady = false
		// Pixie Swarm chant.
		castPhonemes(con.unit, []audio.Name{PhLeft, PhDown, PhRight, PhDown}, func() {
			ns.CastSpell(spell.PIXIE_SWARM, con.unit, con.unit)
			con.spells.Ready = true
			ns.NewTimer(ns.Seconds(10), func() {
				con.spells.PixieSwarmReady = true
			})
		})
	}
}

func (con *Conjurer) FistOfVengeance() {
	if con.spells.FistOfVengeanceReady && con.spells.Ready {
		con.spells.Ready = false
		con.spells.FistOfVengeanceReady = false
		// Fist of Vengeance chant.
		castPhonemes(con.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
			con.unit.LookAtObject(con.target)
			ns.CastSpell(spell.FIST, con.unit, con.target)
			con.spells.Ready = true
			ns.NewTimer(ns.Seconds(10), func() {
				con.spells.FistOfVengeanceReady = true
			})
		})
	}
}

func (con *Conjurer) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			"InfinitePainWand", "LesserFireballWand",
			//"SulphorousShowerWand",
			//"SulphorousFlareWand",
			//"StaffWooden",

			// Crossbow and Bow.
			"CrossBow",
			"Bow",
			"Quiver",
		},
	)
	for _, item := range weapons {
		con.unit.Equip(item)
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
		ns.HasTypeName{
			// Conjurer Helm.
			"ConjurerHelm",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor", "LeatherHelm", "LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		con.unit.Equip(item)
	}
}
