package BotWars

import (
	"image/color"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWarrior creates a new Warrior bot.
func NewWarrior(t *Team) *Warrior {
	war := &Warrior{team: t}
	war.init()
	return war
}

// Warrior bot class.
type Warrior struct {
	team              *Team
	unit              ns.Obj
	target            ns.Obj
	cursor            ns.Pointf
	targetPotion      ns.Obj
	startingEquipment struct {
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
	behaviour struct {
		listening         bool
		lookingForHealing bool
		charging          bool
		attacking         bool
		lookingForTarget  bool
	}
	reactionTime int
}

func (war *Warrior) init() {
	// Reset Behaviour
	war.behaviour.listening = true
	war.behaviour.attacking = false
	war.behaviour.lookingForHealing = false
	war.behaviour.charging = false
	war.behaviour.lookingForTarget = true
	// Inventory
	// Reset abilities WarBot.
	war.abilities.isAlive = true
	war.abilities.Ready = true
	war.abilities.BerserkerChargeReady = true
	war.abilities.WarCryReady = true
	war.abilities.HarpoonReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.TreadLightlyReady = true
	// Select spawnpoint.
	// Create WarBot.
	war.unit = ns.CreateObject("NPC", war.team.SpawnPoint())
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	war.unit.SetStrength(125)
	war.unit.SetBaseSpeed(100)
	// Set Team.
	war.unit.SetOwner(war.team.TeamObj)
	war.unit.SetTeam(war.team.TeamObj.Team())
	if war.team.TeamObj.HasTeam(ns.Teams()[0]) {
		war.unit.SetColor(0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	} else {
		war.unit.SetColor(0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	}
	// Create WarBot mouse cursor.
	war.target = war.team.Enemy.TeamObj
	war.cursor = war.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	war.reactionTime = 15
	// Set WarBot properties.
	war.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	war.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	//war.unit.MonsterStatusEnable(object.MonStatusAlert)
	war.unit.AggressionLevel(0.83)
	war.unit.Hunt()
	war.unit.ResumeLevel(1)
	war.unit.RetreatLevel(0.0)
	// Create and equip WarBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	war.startingEquipment.Longsword = ns.CreateObject("Longsword", ns.Ptf(150, 150))
	war.startingEquipment.WoodenShield = ns.CreateObject("WoodenShield", ns.Ptf(150, 150))
	war.startingEquipment.StreetSneakers = ns.CreateObject("StreetSneakers", ns.Ptf(150, 150))
	war.startingEquipment.StreetPants = ns.CreateObject("StreetPants", ns.Ptf(150, 150))
	war.unit.Equip(war.startingEquipment.Longsword)
	war.unit.Equip(war.startingEquipment.WoodenShield)
	war.unit.Equip(war.startingEquipment.StreetSneakers)
	war.unit.Equip(war.startingEquipment.StreetPants)
	// Select a WarBot loadout (tactical preference, dialog). TODO: Give different audio and chat for each set so they feel like different characters.
	// On looking for enemy.
	war.unit.OnEvent(ns.EventLookingForEnemy, war.onLookingForEnemy)
	// On heard.
	war.unit.OnEvent(ns.EventEnemyHeard, war.onEnemyHeard)
	// Enemy sighted.
	war.unit.OnEvent(ns.EventEnemySighted, war.onEnemySighted)
	// Enemy lost.
	war.unit.OnEvent(ns.EventLostEnemy, war.onLostEnemy)
	// On end of waypoint.
	war.unit.OnEvent(ns.EventEndOfWaypoint, war.onEndOfWaypoint)
	// On change focus.
	war.unit.OnEvent(ns.EventChangeFocus, war.onChangeFocus)
	// On collision.
	war.unit.OnEvent(ns.EventCollision, war.onCollide)
	// On hit.
	war.unit.OnEvent(ns.EventIsHit, war.onHit)
	// Retreat.
	war.unit.OnEvent(ns.EventRetreat, war.onRetreat)
	// On death.
	war.unit.OnEvent(ns.EventDeath, war.onDeath)
	war.LookForWeapon()
}

func (war *Warrior) onChangeFocus() {
	//if !war.behaviour.lookingForHealing {
	//war.unit.Chat("onChangeFocus")
	//}
}

func (war *Warrior) onLookingForEnemy() {

	//if !war.behaviour.lookingForHealing {
	//war.unit.Chat("onLookingForEnemy")
	//}
}

func (war *Warrior) onEnemyHeard() {
	if !war.behaviour.lookingForHealing && !war.behaviour.attacking {
		//war.unit.Chat("onEnemyHeard")
		war.behaviour.attacking = true
		//war.WarBotDetectEnemy() TEMP DISABLE
		//if war.behaviour.listening {
		//	war.behaviour.listening = false
		//	war.unit.Chat("Wiz06a:Guard2Listen")
		//	war.unit.Guard(war.target.Pos(), war.target.Pos(), 300)
		//	ns.NewTimer(ns.Seconds(10), func() {
		//		war.behaviour.listening = true
		//	})
		//}
	}
}

func (war *Warrior) onCollide() {
	if war.abilities.isAlive {
		caller := ns.GetCaller()
		// CTF Logic.
		war.team.CheckPickUpEnemyFlag(caller, war.unit)
		war.team.CheckCaptureEnemyFlag(caller, war.unit)
		war.team.CheckRetrievedOwnFlag(caller, war.unit)
		//	if !war.behaviour.lookingForHealing {
		//war.unit.Chat("onCollide")
		// TODO: determine tactic.
		//}
		//if ns.GetCaller() == BlueFlag {
		//	RedTeamHasBlueFlag = false
		//	RedFlagIsAtBase = true
		//	ns.AudioEvent(audio.FlagRespawn, ns.GetHost()) // <----- replace with all players
		//	BlueFlagFront.SetPos(BlueFlagOutOfGame)
		//	BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
		//}

	}
}

func (war *Warrior) onEnemySighted() {
	war.target = ns.GetCaller()
	// SCRIPT FOR WEAPON SWITCHING. On HOLD FOR NOW
	//if war.unit.HasItem(ns.Object("FanChakram")) {
	//	war.unit.Chat("HELLLLOOOOOO")
	//war.unit.Equip(ns.Object("FanChakram"))
	//war.unit.HitRanged(war.target.Pos())
	//}

	if !war.behaviour.lookingForHealing {
		//war.unit.Chat("onEnemySighted")
		//war.WarBotDetectEnemy() TEMP DISALBE
		war.useWarCry()
	}
}

func (war *Warrior) onRetreat() {
	//war.unit.Chat("onRetreat")
	// TODO: FIX IT!
	//if war.behaviour.lookForHealth {
	//	war.behaviour.listening = false
	//	war.behaviour.lookForHealth = false
	//	war.unit.Chat("Con02A:NecroTalk02")
	//	// Walk to nearest RedPotion.
	//	war.targetPotion = ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
	//	war.unit.AggressionLevel(0.16)
	//	war.unit.Guard(war.targetPotion.Pos().Pos(), war.targetPotion.Pos(), 50)
	//	ns.NewTimer(ns.Seconds(10), func() {
	//		war.behaviour.lookForHealth = true
	//		war.behaviour.listening = true
	//	})
	//}
}

func (war *Warrior) onLostEnemy() {
	if !war.behaviour.lookingForHealing {
		war.useEyeOfTheWolf()
		//war.unit.Chat("onLostEnemy")
		//war.unit.Chat("Multi:General10")
		war.behaviour.attacking = false
		war.unit.Hunt()
	}
	war.team.WalkToOwnFlag(war.unit)
}

func (war *Warrior) onHit() {
	//if war.unit.CurrentHealth() <= 100 && war.target.CurrentHealth() >= 50 && !war.behaviour.lookingForHealing && war.inventory.redPotionInInventory <= 0 {
	//	//war.unit.Chat("onHit")
	//	war.lookForRedPotion()
	//	//war.unit.Guard(war.targetPotion.Pos().Pos(), war.targetPotion.Pos(), 50)
	//}
	//if war.unit.CurrentHealth() <= 100 && war.inventory.redPotionInInventory >= 1 {
	//		for _, it := range war.unit.Items() {
	//			if it.Type().Name() == "RedPotion" {
	//				war.unit.Drop(it)
	//				war.inventory.redPotionInInventory = war.inventory.redPotionInInventory - 2
	//			}
	//		}
	//	}
}

func (war *Warrior) onEndOfWaypoint() {
	//if war.behaviour.lookingForHealing {
	//	if war.unit.CurrentHealth() >= 140 {
	//		//war.unit.Chat("onEndOfWaypoint")
	//		war.unit.AggressionLevel(0.83)
	//		war.unit.Hunt()
	//		war.behaviour.lookingForHealing = false
	//	} else {
	//		//if war.inventory.redPotionInInventory <= 1 {
	//		//	war.lookForRedPotion()
	//		//}
	//	}
	//} else {
	//	if !war.behaviour.lookingForTarget {
	//		war.unit.Hunt()
	//		war.unit.AggressionLevel(0.83)
	//		war.behaviour.lookingForTarget = true
	//	}
	//}
	war.team.CheckAttackOrDefend(war.unit)
	war.LookForNearbyItems()
}

func (war *Warrior) lookForRedPotion() {
	//if war.inventory.redPotionInInventory >= 1 {
	//	war.onEndOfWaypoint()
	//} else {
	//	war.behaviour.lookingForHealing = true
	//	war.unit.AggressionLevel(0.16)
	//	war.unit.WalkTo(war.targetPotion.Pos())
	//}

}

func (war *Warrior) onDeath() {
	war.abilities.isAlive = false
	war.unit.DestroyChat()
	war.team.DropEnemyFlag(war.unit)
	ns.AudioEvent(audio.NPCDie, war.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, war.unit)
		war.unit.Delete()
		war.startingEquipment.StreetPants.Delete()
		war.startingEquipment.StreetSneakers.Delete()
		war.init()
	})
}

func (war *Warrior) UsePotions() {
	if war.unit.CurrentHealth() <= 100 && war.unit.InItems().FindObjects(nil, ns.HasTypeName{"RedPotion"}) != 0 {
		ns.AudioEvent(audio.LesserHealEffect, war.unit)
		RedPotion := war.unit.Items(ns.HasTypeName{"RedPotion"})
		war.unit.SetHealth(war.unit.CurrentHealth() + 50)
		RedPotion[0].Delete()
	}
}

func (war *Warrior) Update() {
	if InitLoadComplete {
		war.UsePotions()
		if war.unit.HasEnchant(enchant.HELD) {
			ns.CastSpell(spell.SLOW, war.unit, war.unit)
			war.unit.EnchantOff(enchant.HELD)
		}
		war.findLoot()
		//war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		war.targetPotion = ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
	}
}

func (war *Warrior) LookForWeapon() {
	ItemLocation := ns.FindClosestObject(war.unit, ns.HasTypeName{"GreatSword", "WarHammer"})
	war.unit.WalkTo(ItemLocation.Pos())
}

func (war *Warrior) LookForNearbyItems() {
	if ns.FindAllObjects(ns.HasTypeName{ //"LeatherArmoredBoots", "LeatherArmor",
		//"LeatherHelm",
		"GreatSword", "WarHammer", //"MorningStar", "BattleAxe", "Sword", "OgreAxe",
		//"LeatherLeggings", "LeatherArmbands",
		"RoundChakram", //"FanChakram",
		//"CurePoisonPotion",
		// Plate armor.
		//"OrnateHelm",
		//"SteelHelm",
		//"Breastplate", "PlateLeggings", "PlateBoots", "PlateArms", "SteelShield",

		// Chainmail armor.
		//	"ChainCoif",
		//"ChainTunic", "ChainLeggings",
		//"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		"RedPotion"},
		ns.InCirclef{Center: war.unit, R: 200}) != nil {
		if war.unit.InItems().FindObjects(nil, ns.HasTypeName{"GreatSword", "WarHammer", "RoundChakram", "RedPotion"}) == 0 {
			ItemLocation := ns.FindAllObjects(ns.HasTypeName{ //"LeatherArmoredBoots", "LeatherArmor",
				//"LeatherHelm",
				"GreatSword", "WarHammer", // "MorningStar", "BattleAxe", "Sword", "OgreAxe",
				//"LeatherLeggings", "LeatherArmbands",
				"RoundChakram", //"FanChakram",
				//"CurePoisonPotion",
				// Plate armor.
				//"OrnateHelm",
				//"SteelHelm",
				//"Breastplate", "PlateLeggings", "PlateBoots", "PlateArms", "SteelShield",

				// Chainmail armor.
				//"ChainCoif",
				//"ChainTunic", "ChainLeggings",
				//"LeatherBoots", "MedievalCloak", "MedievalShirt",
				//"MedievalPants",
				"RedPotion"},
				ns.InCirclef{Center: war.unit, R: 200},
			)
			if war.unit.CanSee(ItemLocation[0]) {
				war.unit.WalkTo(ItemLocation[0].Pos())
			}
		}
	}
}

func (war *Warrior) findLoot() {
	const dist = 75
	// Melee weapons.
	meleeweapons := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"GreatSword", "WarHammer", "MorningStar", "BattleAxe", "Sword", "OgreAxe",
			//"StaffWooden",
		},
	)
	for _, item := range meleeweapons {
		if war.unit.CanSee(item) {
			war.unit.Equip(item)
		}
	}

	// Throwing weapons.
	throwingweapons := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"RoundChakram", "FanChakram",
		},
	)
	for _, item := range throwingweapons {
		if war.unit.CanSee(item) {
			war.unit.Pickup(item)
		}
	}

	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"RedPotion",
			"CurePoisonPotion",
		},
	)
	for _, item := range potions {
		if war.unit.CanSee(item) {
			war.unit.Pickup(item)
		}
	}

	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			// Plate armor.
			"OrnateHelm",
			"SteelHelm",
			"Breastplate", "PlateLeggings", "PlateBoots", "PlateArms", "SteelShield",

			// Chainmail armor.
			"ChainCoif",
			"ChainTunic", "ChainLeggings",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor",
			//"LeatherHelm",
			"LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if war.unit.CanSee(item) {
			war.unit.Equip(item)
		}
	}
}

func (war *Warrior) useWarCry() {
	// Check if cooldown is ready.
	if war.abilities.WarCryReady && !war.behaviour.charging {
		// Trigger global cooldown.
		war.abilities.Ready = false
		if war.target.MaxHealth() == 150 {
		} else {
			war.abilities.WarCryReady = false
			// Check reaction time based on difficulty setting.
			ns.NewTimer(ns.Frames(war.reactionTime), func() {
				war.unit.Pause(ns.Seconds(1))
				ns.AudioEvent("WarcryInvoke", war.unit)
				ns.FindObjects(
					// Target enemy players.
					func(it ns.Obj) bool {
						if war.unit.CanSee(it) && it.MaxHealth() < 150 {
							ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
							it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasClass(object.ClassPlayer),
					ns.HasTeam{war.team.Enemy.TeamObj.Team()},
				)
				// Select target.
				// Target enemy bots.
				ns.FindObjects(
					func(it ns.Obj) bool {
						if war.unit.CanSee(it) && it.MaxHealth() < 150 {
							ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
							it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasTypeName{"NPC"},
					ns.HasTeam{war.team.Enemy.TeamObj.Team()},
				)
				//Target enemy monsters small.
				//	ns.FindObjects(
				//		func(it ns.Obj) bool {
				//			ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
				//			it.Enchant(enchant.HELD, ns.Seconds(3))
				//			return true
				//		},
				//		ns.InCirclef{Center: war.unit, R: 300},
				//		ns.HasTypeName{"Urchin", "Bat, Bomber", "SmallSpider", "Ghost", "Imp", "FlyingGolem"},
				//		// "HasOwner in Enemy.Team"
				//	)
				//	 Target enemy monsters casters.
				// continue script.
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
