package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
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
	unit  ns.Obj
	items struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
	}
	spells struct {
		Ready            bool // Duration unknown.
		InfravisionReady bool // Duration is 30 seconds.
		VampirismReady   bool // Duration is 30 seconds.
		BlinkReady       bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
	}
}

func (con *Conjurer) init() {
	// Reset spells
	con.spells.Ready = true
	con.spells.InfravisionReady = true
	con.spells.VampirismReady = true
	con.spells.BlinkReady = true

	con.unit = ns.CreateObject("NPC", ns.GetHost())
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
	// Escape.
	con.unit.OnEvent(ns.EventRetreat, con.onRetreat)
	// Enemy Lost.
	con.unit.OnEvent(ns.EventLostEnemy, con.onLostEnemy)
	// On Death.
	con.unit.OnEvent(ns.EventDeath, con.onDeath)
}

func (con *Conjurer) buffInitial() {
	if con.spells.VampirismReady {
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

func (con *Conjurer) onRetreat() {
	if con.spells.BlinkReady {
		con.spells.BlinkReady = false
		con.spells.Ready = false
		castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
			ns.NewTrap(con.unit, spell.BLINK) // TODO: FIX IT so it doesn't have to be a trap.
			con.spells.Ready = true
		})
		ns.NewTimer(ns.Seconds(2), func() {
			con.spells.BlinkReady = true
		})
	}
}

func (con *Conjurer) onLostEnemy() {
	if con.spells.InfravisionReady {
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
