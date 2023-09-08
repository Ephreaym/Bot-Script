package BotWars

import (
	ns3 "github.com/noxworld-dev/noxscript/ns/v3"
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
		isAlive              bool
		Ready                bool // Global cooldown.
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
	war.abilities.Ready = true
	war.abilities.BerserkerChargeReady = true
	war.abilities.WarCryReady = true
	war.abilities.HarpoonReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.TreadLightlyReady = true
	// Create WarBot.
	war.unit = ns.CreateObject("NPC", RandomBotSpawn)
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	// Set Team.
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
	war.WarBotDetectEnemy()
	war.useWarCry()
}

func (war *Warrior) onRetreat() {
	// TODO: FIX IT!
	// Walk to nearest RedPotion.
	war.targetPotion = ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
	// war.unit.AggressionLevel(0.16)
	// war.unit.WalkTo(war.targetPotion.Pos())

}

func (war *Warrior) onLostEnemy() {
	war.useEyeOfTheWolf()
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
	// Check if cooldown is ready.
	if war.abilities.WarCryReady {
		// Select target.
		war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		// Trigger global cooldown.
		war.abilities.Ready = false
		if war.target.MaxHealth() == 150 {
		} else {
			war.abilities.WarCryReady = false
			// Check reaction time based on difficulty setting.
			ns.NewTimer(ns.Frames(war.reactionTime), func() {
				war.unit.Enchant(enchant.HELD, ns.Seconds(1))
				ns.AudioEvent("WarcryInvoke", war.unit)
				ns.CastSpell(spell.COUNTERSPELL, war.unit, war.target)
				war.target.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
				war.unit.EnchantOff(enchant.INVULNERABLE)
				ns.NewTimer(ns.Seconds(10), func() {
					war.abilities.WarCryReady = true
				})
				ns.NewTimer(ns.Seconds(1), func() {
					war.abilities.Ready = true
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

// --------------------------------------------------------- //

var (
	WarBot                  ns3.ObjectID
	WizBot                  ns3.ObjectID
	WizBotCorpse            ns3.ObjectID
	WarBotCorpse            ns3.ObjectID
	OutOfGameWar            ns3.WaypointID
	OutOfGameWiz            ns3.WaypointID
	WarSound                ns3.WaypointID
	WizSound                ns3.WaypointID
	WarCryCooldown          bool
	EyeOfTheWolfCooldown    bool
	BerserkerChargeCooldown bool
	GlobalCooldown          bool
	RespawnCooldownDelay    bool
)

func (war *Warrior) UnitRatioX(unit, target ns3.ObjectID, size float32) float32 {
	return (ns3.GetObjectX(unit) - ns3.GetObjectX(target)) * size / ns3.Distance(ns3.GetObjectX(unit), ns3.GetObjectY(unit), ns3.GetObjectX(target), ns3.GetObjectY(target))
}

func (war *Warrior) UnitRatioY(unit, target ns3.ObjectID, size float32) float32 {
	return (ns3.GetObjectY(unit) - ns3.GetObjectY(target)) * size / ns3.Distance(ns3.GetObjectX(unit), ns3.GetObjectY(unit), ns3.GetObjectX(target), ns3.GetObjectY(target))
}
func (war *Warrior) WarBotDetectEnemy() {
	if !BerserkerChargeCooldown && !GlobalCooldown {
		rnd := ns3.Random(0, 2)

		if (rnd == 0) || (rnd == 1) {
			BerserkerChargeCooldown = true
			GlobalCooldown = true
			ns3.SecondTimer(1, war.GlobalCooldownReset)
			war.BerserkerInRange(ns3.GetTrigger(), ns3.GetCaller(), 10)
		}
	} else {
		if !war.abilities.WarCryReady {
			war.useWarCry()
		}
	}
}
func (war *Warrior) CheckUnitFrontSight(unit ns3.ObjectID, dtX, dtY float32) bool {
	ns3.MoveWaypoint(1, ns3.GetObjectX(unit)+dtX, ns3.GetObjectY(unit)+dtY)
	temp := ns3.CreateObject("InvisibleLightBlueHigh", 1)
	res := ns3.IsVisibleTo(unit, temp)
	ns3.Delete(temp)
	return res
}

func (war *Warrior) BerserkerInRange(owner, target ns3.ObjectID, wait int) {
	if ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 {
		if !ns3.HasEnchant(owner, "ENCHANT_ETHEREAL") {
			ns3.Enchant(owner, "ENCHANT_ETHEREAL", 0.0)
			ns3.MoveWaypoint(1, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
			unit := ns3.CreateObject("InvisibleLightBlueHigh", 1)
			ns3.MoveWaypoint(1, ns3.GetObjectX(unit), ns3.GetObjectY(unit))
			unit1 := ns3.CreateObject("InvisibleLightBlueHigh", 1)
			ns3.LookWithAngle(unit, wait)
			ns3.FrameTimer(1, func() {
				war.BerserkerWaitStrike(unit, unit1, owner, target, wait)
			})
		}
	}
}

func (war *Warrior) BerserkerWaitStrike(ptr, ptr1, owner, target ns3.ObjectID, count int) {
	for {
		if ns3.IsObjectOn(ptr) && ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 && ns3.IsObjectOn(owner) {
			if count != 0 {
				if ns3.IsVisibleTo(owner, target) && ns3.Distance(ns3.GetObjectX(owner), ns3.GetObjectY(owner), ns3.GetObjectX(target), ns3.GetObjectY(target)) < 400.0 {
					war.BerserkerCharge(owner, target)
				} else {
					ns3.FrameTimer(6, func() {
						war.BerserkerWaitStrike(ptr, ptr1, owner, target, count-1)
					})
					break
				}
			}
		}
		if ns3.CurrentHealth(owner) != 0 {
			ns3.EnchantOff(owner, "ENCHANT_ETHEREAL")
		}
		if ns3.IsObjectOn(ptr) {
			ns3.Delete(ptr)
			ns3.Delete(ptr1)
		}
		break
	}
}

func (war *Warrior) BerserkerCharge(owner, target ns3.ObjectID) {
	if ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 {
		ns3.EnchantOff(owner, "ENCHANT_INVULNERABLE")
		ns3.MoveWaypoint(2, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
		ns3.AudioEvent("BerserkerChargeInvoke", 2)
		ns3.MoveWaypoint(1, ns3.GetObjectX(owner), ns3.GetObjectY(owner))

		unit := ns3.CreateObject("InvisibleLightBlueHigh", 1)
		ns3.MoveWaypoint(1, ns3.GetObjectX(unit), ns3.GetObjectY(unit))

		unit1 := ns3.CreateObject("InvisibleLightBlueHigh", 1)
		ns3.LookAtObject(unit1, target)

		ns3.LookWithAngle(ns3.GetLastItem(owner), 0)
		ns3.SetCallback(owner, 9, war.BerserkerTouched)

		ratioX := war.UnitRatioX(target, owner, 23.0)
		ratioY := war.UnitRatioY(target, owner, 23.0)
		ns3.FrameTimer(1, func() {
			war.BerserkerLoop(unit, unit1, owner, target, ratioX, ratioY)
		})
	}
}

func (war *Warrior) BerserkerLoop(ptr, ptr1, owner, target ns3.ObjectID, ratioX, ratioY float32) {
	count := ns3.GetDirection(ptr)

	if ns3.CurrentHealth(owner) != 0 && count < 60 && ns3.IsObjectOn(ptr) && ns3.IsObjectOn(owner) {
		if war.CheckUnitFrontSight(owner, ratioX*1.5, ratioY*1.5) && ns3.GetDirection(ns3.GetLastItem(owner)) == 0 {
			ns3.MoveObject(owner, ns3.GetObjectX(owner)+ratioX, ns3.GetObjectY(owner)+ratioY)
			ns3.LookWithAngle(owner, ns3.GetDirection(ptr1))
			ns3.Walk(owner, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
		} else {
			ns3.LookWithAngle(ptr, 100)
		}
		ns3.FrameTimer(1, func() {
			war.BerserkerLoop(ptr, ptr1, owner, target, ratioX, ratioY)
		})
	} else {
		ns3.SetCallback(owner, 9, war.NullCollide)
		ns3.Delete(ptr)
		ns3.Delete(ptr1)
	}
}

func (war *Warrior) BerserkerTouched() {
	self, other := ns3.GetTrigger(), ns3.GetCaller()
	if ns3.IsObjectOn(self) {
		for {
			if ns3.GetCaller() == 0 || (ns3.HasClass(other, "IMMOBILE") && !ns3.HasClass(other, "DOOR") && !ns3.HasClass(other, "TRIGGER")) && !ns3.HasClass(other, "DANGEROUS") {
				ns3.MoveWaypoint(2, ns3.GetObjectX(self), ns3.GetObjectY(self))
				ns3.AudioEvent("FleshHitStone", 2)

				ns3.Enchant(self, "ENCHANT_HELD", 2.0)
			} else if ns3.CurrentHealth(other) != 0 {
				if ns3.IsAttackedBy(self, other) {
					ns3.MoveWaypoint(2, ns3.GetObjectX(self), ns3.GetObjectY(self))
					ns3.AudioEvent("FleshHitFlesh", 2)
					ns3.Damage(other, self, 100, 2)
				} else {
					break
				}
			} else {
				break
			}
			ns3.LookWithAngle(ns3.GetLastItem(self), 1)
			break
		}
	}
	war.unit.Hunt()
	ns3.SecondTimer(10, war.BerserkerChargeCooldownReset)
}
func (war *Warrior) NullCollide() {
}
func (war *Warrior) BerserkerChargeCooldownReset() {
	if !RespawnCooldownDelay {
		BerserkerChargeCooldown = false
	}
}
func (war *Warrior) GlobalCooldownReset() {
	GlobalCooldown = false
}
