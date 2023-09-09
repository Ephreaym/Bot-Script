package BotWars

import (
	"image/color"

	ns3 "github.com/noxworld-dev/noxscript/ns/v3"

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
		RoundChackramReady   bool // for now cooldown 10 seconds.
	}
	behaviour struct {
		listening         bool
		lookingForHealing bool
		charging          bool
		attacking         bool
		lookingForTarget  bool
		AntiStuck         bool
		SwitchMainWeapon  bool
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
	war.behaviour.AntiStuck = true
	war.behaviour.SwitchMainWeapon = false
	// Inventory
	// Reset abilities WarBot.
	war.abilities.isAlive = true
	war.abilities.Ready = true
	war.abilities.BerserkerChargeReady = true
	war.abilities.WarCryReady = true
	war.abilities.HarpoonReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.TreadLightlyReady = true
	war.abilities.RoundChackramReady = true
	// Select spawnpoint.
	// Create WarBot.
	war.unit = ns.CreateObject("NPC", war.team.SpawnPoint())
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	war.unit.SetStrength(125)
	war.unit.SetBaseSpeed(100)
	// Set Team.
	war.unit.SetOwner(war.team.Spawns()[0])
	war.unit.SetTeam(war.team.Team())
	if war.unit.HasTeam(ns.Teams()[0]) {
		war.unit.SetColor(0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		war.unit.SetColor(1, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		war.unit.SetColor(2, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		war.unit.SetColor(3, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		war.unit.SetColor(4, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		war.unit.SetColor(5, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	} else {
		war.unit.SetColor(0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		war.unit.SetColor(1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		war.unit.SetColor(2, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		war.unit.SetColor(3, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		war.unit.SetColor(4, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		war.unit.SetColor(5, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	}
	// Create WarBot mouse cursor.
	war.target = war.team.Enemy.Spawns()[0]
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
	war.WeaponPreference()
}

func (war *Warrior) onChangeFocus() {
	//if !war.behaviour.lookingForHealing {
	//war.unit.Chat("onChangeFocus")
	//}
	war.useWarCry()
}

func (war *Warrior) onLookingForEnemy() {

	//if !war.behaviour.lookingForHealing {
	//war.unit.Chat("onLookingForEnemy")
	//}
}

func (war *Warrior) onEnemyHeard() {
	//if !war.behaviour.lookingForHealing && !war.behaviour.attacking {
	//war.unit.Chat("onEnemyHeard")
	//war.behaviour.attacking = true
	//war.WarBotDetectEnemy()
	//if war.behaviour.listening {
	//	war.behaviour.listening = false
	//	war.unit.Chat("Wiz06a:Guard2Listen")
	//	war.unit.Guard(war.target.Pos(), war.target.Pos(), 300)
	//	ns.NewTimer(ns.Seconds(10), func() {
	//		war.behaviour.listening = true
	//	})
	//}
	//}

	// WORKING SCRIPT. TEMP DISABLE.
	war.ThrowChakram()
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
	//war.WarBotDetectEnemy()

	// WORKING SCRIPT TEMP DISABLE.
	//if !war.behaviour.lookingForHealing {
	war.useWarCry()
	//}
	war.ThrowChakram()
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

//func (war *Warrior) lookForRedPotion() {
//if war.inventory.redPotionInInventory >= 1 {
//	war.onEndOfWaypoint()
//} else {
//	war.behaviour.lookingForHealing = true
//	war.unit.AggressionLevel(0.16)
//	war.unit.WalkTo(war.targetPotion.Pos())
//}

//}

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
	if ItemLocation != nil {
		war.unit.WalkTo(ItemLocation.Pos())
	}
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
	ns.NewTimer(ns.Seconds(5), func() {
		// prevent bots getting stuck to stay in loop.
		if war.behaviour.AntiStuck {
			war.behaviour.AntiStuck = false
			war.team.CheckAttackOrDefend(war.unit)
			ns.NewTimer(ns.Seconds(6), func() {
				war.behaviour.AntiStuck = true
			})
		}
	})
}

func (war *Warrior) ThrowChakram() {
	if war.abilities.RoundChackramReady && war.unit.InItems().FindObjects(nil, ns.HasTypeName{"RoundChakram"}) != 0 {
		war.abilities.RoundChackramReady = false
		//war.unit.Chat("I have a Chakram")
		war.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				war.unit.Equip(it)
				war.unit.LookAtObject(war.target)
				war.unit.HitRanged(war.target.Pos())
				ns.NewTimer(ns.Frames(5), func() {
					war.WeaponPreference()
				})
				ns.NewTimer(ns.Seconds(10), func() {
					war.abilities.RoundChackramReady = true
				})
				return true
			},
			ns.HasTypeName{"RoundChakram"},
		)
	}
}

func (war *Warrior) WeaponPreference() {
	// Priority list to get the prefered weapon.
	// TODO: Add stun and range conditions.
	if war.unit.InItems().FindObjects(nil, ns.HasTypeName{"GreatSword"}) != 0 && war.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"GreatSword"}) == 0 {
		war.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				war.unit.Equip(it)
				//war.unit.Chat("I swapped to my GreatSword!")
				return true
			},
			ns.HasTypeName{"GreatSword"},
		)
	} else if war.unit.InItems().FindObjects(nil, ns.HasTypeName{"WarHammer"}) != 0 && war.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"WarHammer"}) == 0 {
		war.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				war.unit.Equip(it)
				//war.unit.Chat("I swapped to my WarHammer!")
				return true
			},
			ns.HasTypeName{"WarHammer"},
		)
	} else {
		if war.unit.InItems().FindObjects(nil, ns.HasTypeName{"WarHammer", "GreatSword"}) == 0 {
			war.unit.Equip(war.startingEquipment.Longsword)
			//war.unit.Chat("I swapped to my LongSword!")
		}
	}
	ns.NewTimer(ns.Seconds(10), func() {
		war.WeaponPreference()
	})
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
			war.unit.Pickup(item)
			war.unit.Equip(war.unit.GetLastItem())
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
			"Breastplate", "PlateLeggings", "PlateBoots", "PlateArms",
			//"SteelShield",

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
			war.unit.Pickup(item)
			war.unit.Equip(war.unit.GetLastItem())
		}
	}
}

func (war *Warrior) useWarCry() {
	// Check if cooldown is war.abilities.Ready.
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
						if war.unit.CanSee(it) && it.MaxHealth() < 150 && !it.HasEnchant(enchant.ANTI_MAGIC) {
							ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
							it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasClass(object.ClassPlayer),
					ns.HasTeam{war.team.Enemy.Team()},
				)
				// Select target.
				// Target enemy bots.
				ns.FindObjects(
					func(it ns.Obj) bool {
						if war.unit.CanSee(it) && it.MaxHealth() < 150 && !it.HasEnchant(enchant.ANTI_MAGIC) {
							ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
							it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasTypeName{"NPC"},
					ns.HasTeam{war.team.Enemy.Team()},
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
	// Check if cooldown is war.abilities.Ready.
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
		rnd := ns3.Random(1, 1)

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
