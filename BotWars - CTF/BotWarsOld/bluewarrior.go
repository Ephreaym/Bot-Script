package BotWars

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWarrior creates a new BlueWarrior bot.
func NewBlueWarrior() *BlueWarrior {
	bluewar := &BlueWarrior{}
	bluewar.init()
	return bluewar
}

// BlueWarrior bot class.
type BlueWarrior struct {
	unit         ns.Obj
	target       ns.Obj
	cursor       ns.Pointf
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
	behaviour struct {
		listening         bool
		lookingForHealing bool
		charging          bool
		attacking         bool
		lookingForTarget  bool
	}
	inventory struct {
		redPotionInInventory int
	}
	reactionTime int
}

func (bluewar *BlueWarrior) init() {
	// Reset Behaviour
	bluewar.behaviour.listening = true
	bluewar.behaviour.attacking = false
	bluewar.behaviour.lookingForHealing = false
	bluewar.behaviour.charging = false
	bluewar.behaviour.lookingForTarget = true
	// Inventory
	bluewar.inventory.redPotionInInventory = 0
	// Reset abilities WarBot.
	bluewar.abilities.Ready = true
	bluewar.abilities.BerserkerChargeReady = true
	bluewar.abilities.WarCryReady = true
	bluewar.abilities.HarpoonReady = true
	bluewar.abilities.EyeOfTheWolfReady = true
	bluewar.abilities.TreadLightlyReady = true
	// Select spawnpoint.
	// Create WarBot.
	bluewar.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPoint"))
	bluewar.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	bluewar.unit.SetMaxHealth(150)
	bluewar.unit.SetStrength(125)
	bluewar.unit.SetBaseSpeed(100)
	// Set Team.
	bluewar.unit.SetOwner(TeamBlue)
	// Create WarBot mouse cursor.
	bluewar.target = ns.FindClosestObject(bluewar.unit, ns.HasClass(object.ClassPlayer))
	bluewar.cursor = bluewar.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	bluewar.reactionTime = 15
	// Set WarBot properties.
	bluewar.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	bluewar.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	//bluewar.unit.MonsterStatusEnable(object.MonStatusAlert)
	bluewar.unit.AggressionLevel(0.83)
	bluewar.unit.Hunt()
	bluewar.unit.ResumeLevel(1)
	bluewar.unit.RetreatLevel(0.0)
	// Create and equip WarBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	bluewar.items.Longsword = ns.CreateObject("Longsword", bluewar.unit)
	bluewar.items.WoodenShield = ns.CreateObject("WoodenShield", bluewar.unit)
	bluewar.items.StreetSneakers = ns.CreateObject("StreetSneakers", bluewar.unit)
	bluewar.items.StreetPants = ns.CreateObject("StreetPants", bluewar.unit)
	bluewar.unit.Equip(bluewar.items.Longsword)
	bluewar.unit.Equip(bluewar.items.WoodenShield)
	bluewar.unit.Equip(bluewar.items.StreetSneakers)
	bluewar.unit.Equip(bluewar.items.StreetPants)
	// Select a WarBot loadout (tactical preference, dialog). TODO: Give different audio and chat for each set so they feel like different characters.
	// On looking for enemy.
	bluewar.unit.OnEvent(ns.EventLookingForEnemy, bluewar.onLookingForEnemy)
	// On heard.
	bluewar.unit.OnEvent(ns.EventEnemyHeard, bluewar.onEnemyHeard)
	// Enemy sighted.
	bluewar.unit.OnEvent(ns.EventEnemySighted, bluewar.onEnemySighted)
	// Enemy lost.
	bluewar.unit.OnEvent(ns.EventLostEnemy, bluewar.onLostEnemy)
	// On end of waypoint.
	bluewar.unit.OnEvent(ns.EventEndOfWaypoint, bluewar.onEndOfWaypoint)
	// On change focus.
	bluewar.unit.OnEvent(ns.EventChangeFocus, bluewar.onChangeFocus)
	// On collision.
	bluewar.unit.OnEvent(ns.EventCollision, bluewar.onCollide)
	// On hit.
	bluewar.unit.OnEvent(ns.EventIsHit, bluewar.onHit)
	// Retreat.
	bluewar.unit.OnEvent(ns.EventRetreat, bluewar.onRetreat)
	// On death.
	bluewar.unit.OnEvent(ns.EventDeath, bluewar.onDeath)
}

func (bluewar *BlueWarrior) onChangeFocus() {
	if !bluewar.behaviour.lookingForHealing {
		//bluewar.unit.Chat("onChangeFocus")
	}
}

func (bluewar *BlueWarrior) onLookingForEnemy() {

	if !bluewar.behaviour.lookingForHealing {
		//bluewar.unit.Chat("onLookingForEnemy")
	}
}

func (bluewar *BlueWarrior) onEnemyHeard() {
	if !bluewar.behaviour.lookingForHealing && !bluewar.behaviour.attacking {
		//bluewar.unit.Chat("onEnemyHeard")
		bluewar.behaviour.attacking = true
		//bluewar.WarBotDetectEnemy() TEMP DISABLE
		//if bluewar.behaviour.listening {
		//	bluewar.behaviour.listening = false
		//	bluewar.unit.Chat("Wiz06a:Guard2Listen")
		//	bluewar.unit.Guard(bluewar.target.Pos(), bluewar.target.Pos(), 300)
		//	ns.NewTimer(ns.Seconds(10), func() {
		//		bluewar.behaviour.listening = true
		//	})
		//}
	}
}

func (bluewar *BlueWarrior) onCollide() {
	if !bluewar.behaviour.lookingForHealing {
		//bluewar.unit.Chat("onCollide")
		// TODO: determine tactic.
	}
}

func (bluewar *BlueWarrior) onEnemySighted() {
	// SCRIPT FOR WEAPON SWITCHING. On HOLD FOR NOW
	//if bluewar.unit.HasItem(ns.Object("FanChakram")) {
	//	bluewar.unit.Chat("HELLLLOOOOOO")
	//bluewar.unit.Equip(ns.Object("FanChakram"))
	//bluewar.unit.HitRanged(bluewar.target.Pos())
	//}

	if !bluewar.behaviour.lookingForHealing {
		//bluewar.unit.Chat("onEnemySighted")
		//bluewar.WarBotDetectEnemy() TEMP DISALBE
		//bluewar.useWarCry()
	}
}

func (bluewar *BlueWarrior) onRetreat() {
	//bluewar.unit.Chat("onRetreat")
	// TODO: FIX IT!
	//if bluewar.behaviour.lookForHealth {
	//	bluewar.behaviour.listening = false
	//	bluewar.behaviour.lookForHealth = false
	//	bluewar.unit.Chat("Con02A:NecroTalk02")
	//	// Walk to nearest RedPotion.
	//	bluewar.targetPotion = ns.FindClosestObject(bluewar.unit, ns.HasTypeName{"RedPotion"})
	//	bluewar.unit.AggressionLevel(0.16)
	//	bluewar.unit.Guard(bluewar.targetPotion.Pos().Pos(), bluewar.targetPotion.Pos(), 50)
	//	ns.NewTimer(ns.Seconds(10), func() {
	//		bluewar.behaviour.lookForHealth = true
	//		bluewar.behaviour.listening = true
	//	})
	//}
}

func (bluewar *BlueWarrior) onLostEnemy() {
	if !bluewar.behaviour.lookingForHealing {
		bluewar.useEyeOfTheWolf()
		//bluewar.unit.Chat("onLostEnemy")
		//bluewar.unit.Chat("Multi:General10")
		bluewar.behaviour.attacking = false
		bluewar.unit.Hunt()
	}
}

func (bluewar *BlueWarrior) onHit() {
	//if bluewar.unit.CurrentHealth() <= 100 && bluewar.target.CurrentHealth() >= 50 && !bluewar.behaviour.lookingForHealing && bluewar.inventory.redPotionInInventory <= 0 {
	//	//bluewar.unit.Chat("onHit")
	//	bluewar.lookForRedPotion()
	//	//bluewar.unit.Guard(bluewar.targetPotion.Pos().Pos(), bluewar.targetPotion.Pos(), 50)
	//}
	//if bluewar.unit.CurrentHealth() <= 100 && bluewar.inventory.redPotionInInventory >= 1 {
	//		for _, it := range bluewar.unit.Items() {
	//			if it.Type().Name() == "RedPotion" {
	//				bluewar.unit.Drop(it)
	//				bluewar.inventory.redPotionInInventory = bluewar.inventory.redPotionInInventory - 2
	//			}
	//		}
	//	}
}

func (bluewar *BlueWarrior) onEndOfWaypoint() {
	if bluewar.behaviour.lookingForHealing {
		if bluewar.unit.CurrentHealth() >= 140 {
			//bluewar.unit.Chat("onEndOfWaypoint")
			bluewar.unit.AggressionLevel(0.83)
			bluewar.unit.Hunt()
			bluewar.behaviour.lookingForHealing = false
		} else {
			if bluewar.inventory.redPotionInInventory <= 1 {
				bluewar.lookForRedPotion()
			}
		}
	} else {
		if !bluewar.behaviour.lookingForTarget {
			bluewar.unit.Hunt()
			bluewar.unit.AggressionLevel(0.83)
			bluewar.behaviour.lookingForTarget = true
		}
	}
}

func (bluewar *BlueWarrior) lookForRedPotion() {
	//if bluewar.inventory.redPotionInInventory >= 1 {
	//	bluewar.onEndOfWaypoint()
	//} else {
	//	bluewar.behaviour.lookingForHealing = true
	//	bluewar.unit.AggressionLevel(0.16)
	//	bluewar.unit.WalkTo(bluewar.targetPotion.Pos())
	//}

}

func (bluewar *BlueWarrior) onDeath() {
	bluewar.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, bluewar.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, bluewar.unit)
		bluewar.unit.Delete()
		bluewar.items.StreetPants.Delete()
		bluewar.items.StreetSneakers.Delete()
		bluewar.init()
	})
}

func (bluewar *BlueWarrior) Update() {
	if InitLoadComplete {
		if bluewar.unit.HasEnchant(enchant.HELD) {
			ns.CastSpell(spell.SLOW, bluewar.unit, bluewar.unit)
			bluewar.unit.EnchantOff(enchant.HELD)
		}
		bluewar.findLoot()
		bluewar.target = ns.FindClosestObject(bluewar.unit, ns.HasClass(object.ClassPlayer))
		bluewar.targetPotion = ns.FindClosestObject(bluewar.unit, ns.HasTypeName{"RedPotion"})
	}
}

func (bluewar *BlueWarrior) findLoot() {
	const dist = 75
	// Melee weapons.
	meleeweapons := ns.FindAllObjects(
		ns.InCirclef{Center: bluewar.unit, R: dist},
		ns.HasTypeName{

			"GreatSword", "WarHammer", "MorningStar", "BattleAxe", "Sword", "OgreAxe",

			//"StaffWooden",
		},
	)
	for _, item := range meleeweapons {
		bluewar.unit.Equip(item)
	}

	// Throwing weapons.
	throwingweapons := ns.FindAllObjects(
		ns.InCirclef{Center: bluewar.unit, R: dist},
		ns.HasTypeName{
			"RoundChakram", "FanChakram",
		},
	)
	for _, item := range throwingweapons {
		bluewar.unit.Pickup(item)
	}

	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: bluewar.unit, R: dist},
		ns.HasTypeName{
			"RedPotion", "CurePoisonPotion",
		},
	)
	for _, item := range potions {
		bluewar.unit.Pickup(item)
		if bluewar.inventory.redPotionInInventory < 3 {
			bluewar.inventory.redPotionInInventory = bluewar.inventory.redPotionInInventory + 1
		}
	}

	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: bluewar.unit, R: dist},
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
		bluewar.unit.Equip(item)
	}
}

func (bluewar *BlueWarrior) useWarCry() {
	// Check if cooldown is ready.
	if bluewar.abilities.WarCryReady && !bluewar.behaviour.charging {
		// Select target.
		bluewar.target = ns.FindClosestObject(bluewar.unit, ns.HasClass(object.ClassPlayer))
		// Trigger global cooldown.
		bluewar.abilities.Ready = false
		if bluewar.target.MaxHealth() == 150 {
		} else {
			bluewar.abilities.WarCryReady = false
			// Check reaction time based on difficulty setting.
			ns.NewTimer(ns.Frames(bluewar.reactionTime), func() {
				bluewar.unit.Pause(ns.Seconds(1))
				ns.AudioEvent("WarcryInvoke", bluewar.unit)
				ns.CastSpell(spell.COUNTERSPELL, bluewar.unit, bluewar.target)
				bluewar.target.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
				bluewar.unit.EnchantOff(enchant.INVULNERABLE)
				ns.NewTimer(ns.Seconds(10), func() {
					bluewar.abilities.WarCryReady = true
				})
				ns.NewTimer(ns.Seconds(1), func() {
					bluewar.abilities.Ready = true
				})
			})
		}
	}
}

func (bluewar *BlueWarrior) useEyeOfTheWolf() {
	// Check if cooldown is ready.
	if bluewar.abilities.EyeOfTheWolfReady {
		// Trigger cooldown.
		bluewar.abilities.EyeOfTheWolfReady = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewar.reactionTime), func() {
			// Use ability.
			bluewar.unit.Enchant(enchant.INFRAVISION, ns.Seconds(10))
		})
		// Eye Of The Wolf cooldown.
		ns.NewTimer(ns.Seconds(20), func() {
			bluewar.abilities.EyeOfTheWolfReady = true
		})
	}
}
