package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWarrior creates a new warrior bot.
func NewWarrior() *Warrior {
	war := &Warrior{}
	war.init()
	return war
}

// Warrior bot class.
type Warrior struct {
	unit         ns.Obj
	target       ns.Obj
	taggedPlayer ns.Obj
	items        struct {
		Longsword      ns.Obj
		WoodenShield   ns.Obj
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
	}
	abilities struct {
		EyeOfTheWolfReady    bool // Cooldown is 20 seconds.
		BerserkerChargeReady bool // Cooldown is 10 seconds.
		WarCryReady          bool // Cooldown is 10 seconds.
	}
}

func (war *Warrior) init() {
	// Reset abilities
	war.abilities.BerserkerChargeReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.WarCryReady = true
	war.unit = ns.CreateObject("NPC", RandomBotSpawn)

	war.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	war.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	war.unit.MonsterStatusEnable(object.MonStatusAlert)
	// TODO: Change location of item creation OR stop them from respawning automatically.
	war.items.Longsword = ns.CreateObject("Longsword", war.unit)
	war.items.WoodenShield = ns.CreateObject("WoodenShield", war.unit)
	war.items.StreetSneakers = ns.CreateObject("StreetSneakers", war.unit)
	war.items.StreetPants = ns.CreateObject("StreetPants", war.unit)
	war.unit.Equip(war.items.Longsword)
	war.unit.Equip(war.items.WoodenShield)
	war.unit.Equip(war.items.StreetSneakers)
	war.unit.Equip(war.items.StreetPants)
	// TODO: Give different audio and chat for each set so they feel like different characters.
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	war.unit.AggressionLevel(0.83)
	war.unit.Hunt()
	war.unit.ResumeLevel(0.8)
	war.unit.RetreatLevel(0.2)
	// WarAI.Chat("War01A.scr:Bully1") // this is a robbery! Your money AND your life!
	// ns.AudioEvent("F1ROG01E", WarAI)
	// TODO: Add audio to match the chat: F1ROG01E.
	// Enemy Sighted. //
	war.unit.OnEvent(ns.EventEnemySighted, war.onEnemySighted)
	war.unit.OnEvent(ns.EventRetreat, war.onRetreat)
	// Enemy Lost.
	war.unit.OnEvent(ns.EventLostEnemy, war.onLostEnemy)
	// On Hit.
	war.unit.OnEvent(ns.EventIsHit, war.onHit)
	// On Death.
	war.unit.OnEvent(ns.EventDeath, war.onDeath)
}

func (war *Warrior) onEnemySighted() {
	// Script out a plan of action.
	if war.abilities.WarCryReady {
		war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		if war.target.MaxHealth() == 150 {
		} else {
			war.abilities.WarCryReady = false
			ns.NewTimer(ns.Frames(15), func() {
				war.unit.Enchant(enchant.HELD, ns.Seconds(1))
				ns.AudioEvent("WarcryInvoke", war.unit)
				ns.CastSpell(spell.COUNTERSPELL, war.unit, war.target)
				war.target.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
				war.unit.EnchantOff(enchant.INVULNERABLE)
				ns.NewTimer(ns.Seconds(10), func() {
					war.abilities.WarCryReady = true
				})
			})
		}

	}
}

func (war *Warrior) onRetreat() {
	// Walk to nearest RedPotion.
}

func (war *Warrior) onLostEnemy() {
	if war.abilities.EyeOfTheWolfReady {
		war.unit.Enchant(enchant.INFRAVISION, ns.Seconds(10))
		war.abilities.EyeOfTheWolfReady = false
		ns.NewTimer(ns.Seconds(20), func() {
			war.abilities.EyeOfTheWolfReady = true
		})
	}
}

func (war *Warrior) onHit() {
	// WarAITaggedPlayer = ns.GetCaller() ---> more research needed to select target
}

func (war *Warrior) onDeath() {
	war.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, war.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, war.unit)
		war.unit.Delete()
		war.items.StreetPants.Delete()
		war.items.StreetSneakers.Delete()
		war.init()
	})
}

func (war *Warrior) Update() {
	war.findLoot()
}

func (war *Warrior) findLoot() {
	const dist = 75
	// TODO: Setup different builds and tactics / voices / dialog / chat.

	// Weapons
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			// TODO: Fix crash with thrown weapons. Maybe WarAI tries to pickup the thrown weapon?
			//"RoundChakram", "FanChakram",

			"GreatSword", "WarHammer", "MorningStar", "BattleAxe", "Sword", "OgreAxe",

			//"StaffWooden",
		},
	)
	for _, item := range weapons {
		war.unit.Equip(item)
	}

	// Armor
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			// Plate armor.
			"OrnateHelm", "SteelHelm", "Breastplate", "PlateLeggings", "PlateBoots", "PlateArms", "SteelShield",

			// Chainmail armor.
			"ChainCoif", "ChainTunic", "ChainLeggings",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor", "LeatherHelm", "LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		war.unit.Equip(item)
	}
}
