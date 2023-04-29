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
	cursor       ns.Pointf
	taggedPlayer ns.Obj
	targetPotion ns.Obj
	items        struct {
		Longsword      ns.Obj
		WoodenShield   ns.Obj
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
	}
	abilities struct {
		BerserkerChargeReady bool // Cooldown is 10 seconds.
		WarCryReady          bool // Cooldown is 10 seconds.
		HarpoonReady         bool
		EyeOfTheWolfReady    bool // Cooldown is 20 seconds.
		TreadLightlyReady    bool
	}
	reactionTime int
}

func (war *Warrior) init() {
	// Reset abilities WarBot.
	war.abilities.BerserkerChargeReady = true
	war.abilities.WarCryReady = true
	war.abilities.HarpoonReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.TreadLightlyReady = true
	// Create WarBot.
	war.unit = ns.CreateObject("NPC", RandomBotSpawn)
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	// Create WarBot mouse cursor.
	war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
	war.cursor = war.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	war.reactionTime = 15
	// Set WarBot properties.
	war.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	war.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	war.unit.MonsterStatusEnable(object.MonStatusAlert)
	war.unit.AggressionLevel(0.83)
	war.unit.Hunt()
	war.unit.ResumeLevel(0.8)
	war.unit.RetreatLevel(0.4)
	// Create and equip WarBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	war.items.Longsword = ns.CreateObject("Longsword", war.unit)
	war.items.WoodenShield = ns.CreateObject("WoodenShield", war.unit)
	war.items.StreetSneakers = ns.CreateObject("StreetSneakers", war.unit)
	war.items.StreetPants = ns.CreateObject("StreetPants", war.unit)
	war.unit.Equip(war.items.Longsword)
	war.unit.Equip(war.items.WoodenShield)
	war.unit.Equip(war.items.StreetSneakers)
	war.unit.Equip(war.items.StreetPants)
	// Select a WarBot loadout (tactical preference, dialog). TODO: Give different audio and chat for each set so they feel like different characters.
	// Enemy sighted.
	war.unit.OnEvent(ns.EventEnemySighted, war.onEnemySighted)
	// Retreat.
	war.unit.OnEvent(ns.EventRetreat, war.onRetreat)
	// Enemy lost.
	war.unit.OnEvent(ns.EventLostEnemy, war.onLostEnemy)
	// On hit.
	war.unit.OnEvent(ns.EventIsHit, war.onHit)
	// On collision.
	war.unit.OnEvent(ns.EventCollision, war.onCollide)
	// On death.
	war.unit.OnEvent(ns.EventDeath, war.onDeath)
}

func (war *Warrior) onCollide() {
	// TODO: determine tactic.
}

func (war *Warrior) onEnemySighted() {
	// war.useWarCry()
	// war.useBerserkerCharge()
}

func (war *Warrior) onRetreat() {
	// TODO: FIX IT!
	// Walk to nearest RedPotion.
	war.targetPotion = ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
	// war.unit.AggressionLevel(0.16)
	// war.unit.WalkTo(war.targetPotion.Pos())

}

func (war *Warrior) onLostEnemy() {
	// war.useEyeOfTheWolf()
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
		war.items.Longsword.Delete()
		war.items.WoodenShield.Delete()
		war.init()
	})
}

func (war *Warrior) Update() {
	war.findLoot()
	war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
}

func (war *Warrior) findLoot() {
	const dist = 75
	// Weapons.
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

	// Armor.
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

func (war *Warrior) useWarCry() {
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

func (war *Warrior) useEyeOfTheWolf() {
	// Check if cooldown is ready.
	if war.abilities.EyeOfTheWolfReady {
		// Trigger cooldown.
		war.abilities.EyeOfTheWolfReady = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(war.reactionTime), func() {
			// Use ability.
			war.unit.Enchant(enchant.INFRAVISION, ns.Seconds(10))
		})
		// Eye Of The Wolf cooldown.
		ns.NewTimer(ns.Seconds(20), func() {
			war.abilities.EyeOfTheWolfReady = true
		})
	}
}

func (war *Warrior) useBerserkerCharge() {
	// Check if cooldown is ready.
	if war.abilities.BerserkerChargeReady && war.unit.CanSee(war.target) {
		ns.PrintStrToAll("useBeserkerCharge")
		// Select target.
		war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		war.cursor = war.target.Pos()
		// Trigger cooldown.
		war.abilities.BerserkerChargeReady = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(war.reactionTime), func() {
			// Use ability.
			ns.AudioEvent("BerserkerChargeInvoke", war.unit)
			war.unit.LookAtObject(war.cursor)
			war.loopBerserker()
		})
		// Berserker Charge cooldown.
		ns.NewTimer(ns.Seconds(10), func() {
			war.abilities.BerserkerChargeReady = true
		})
	}
}

func (war *Warrior) loopBerserker() {
	war.unit.SetZ(3)
	//war.unit.Move(war.cursor)
	war.unit.WalkTo(war.cursor.Pos())
	ns.NewTimer(ns.Frames(1), war.loopBerserker)
}
