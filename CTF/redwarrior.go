package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWarrior creates a new RedWarrior bot.
func NewRedWarrior() *RedWarrior {
	redwar := &RedWarrior{}
	redwar.init()
	return redwar
}

// RedWarrior bot class.
type RedWarrior struct {
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
		RedPotionInInventory int
	}
	reactionTime int
}

func (redwar *RedWarrior) init() {
	// Reset Behaviour
	redwar.behaviour.listening = true
	redwar.behaviour.attacking = false
	redwar.behaviour.lookingForHealing = false
	redwar.behaviour.charging = false
	redwar.behaviour.lookingForTarget = true
	// Inventory
	redwar.inventory.RedPotionInInventory = 0
	// Reset abilities WarBot.
	redwar.abilities.isAlive = true
	redwar.abilities.Ready = true
	redwar.abilities.BerserkerChargeReady = true
	redwar.abilities.WarCryReady = true
	redwar.abilities.HarpoonReady = true
	redwar.abilities.EyeOfTheWolfReady = true
	redwar.abilities.TreadLightlyReady = true
	// Select spawnpoint.
	// Create WarBot.
	redwar.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointRed"))
	redwar.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	redwar.unit.SetMaxHealth(150)
	redwar.unit.SetStrength(125)
	redwar.unit.SetBaseSpeed(100)
	// Set Team.
	redwar.unit.SetOwner(TeamRed)
	// Create WarBot mouse cursor.
	redwar.target = TeamBlue
	redwar.cursor = redwar.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	redwar.reactionTime = 15
	// Set WarBot properties.
	redwar.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	redwar.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	//redwar.unit.MonsterStatusEnable(object.MonStatusAlert)
	redwar.unit.AggressionLevel(0.83)
	redwar.unit.Hunt()
	redwar.unit.ResumeLevel(1)
	redwar.unit.RetreatLevel(0.0)
	// Create and equip WarBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	redwar.items.Longsword = ns.CreateObject("Longsword", redwar.unit)
	redwar.items.WoodenShield = ns.CreateObject("WoodenShield", redwar.unit)
	redwar.items.StreetSneakers = ns.CreateObject("StreetSneakers", redwar.unit)
	redwar.items.StreetPants = ns.CreateObject("StreetPants", redwar.unit)
	redwar.unit.Equip(redwar.items.Longsword)
	redwar.unit.Equip(redwar.items.WoodenShield)
	redwar.unit.Equip(redwar.items.StreetSneakers)
	redwar.unit.Equip(redwar.items.StreetPants)
	// Select a WarBot loadout (tactical preference, dialog). TODO: Give different audio and chat for each set so they feel like different characters.
	// On looking for enemy.
	redwar.unit.OnEvent(ns.EventLookingForEnemy, redwar.onLookingForEnemy)
	// On heard.
	redwar.unit.OnEvent(ns.EventEnemyHeard, redwar.onEnemyHeard)
	// Enemy sighted.
	redwar.unit.OnEvent(ns.EventEnemySighted, redwar.onEnemySighted)
	// Enemy lost.
	redwar.unit.OnEvent(ns.EventLostEnemy, redwar.onLostEnemy)
	// On end of waypoint.
	redwar.unit.OnEvent(ns.EventEndOfWaypoint, redwar.onEndOfWaypoint)
	// On change focus.
	redwar.unit.OnEvent(ns.EventChangeFocus, redwar.onChangeFocus)
	// On collision.
	redwar.unit.OnEvent(ns.EventCollision, redwar.onCollide)
	// On hit.
	redwar.unit.OnEvent(ns.EventIsHit, redwar.onHit)
	// Retreat.
	redwar.unit.OnEvent(ns.EventRetreat, redwar.onRetreat)
	// On death.
	redwar.unit.OnEvent(ns.EventDeath, redwar.onDeath)
}

func (redwar *RedWarrior) onChangeFocus() {
	if !redwar.behaviour.lookingForHealing {
		//redwar.unit.Chat("onChangeFocus")
	}
}

func (redwar *RedWarrior) onLookingForEnemy() {

	if !redwar.behaviour.lookingForHealing {
		//redwar.unit.Chat("onLookingForEnemy")
	}
}

func (redwar *RedWarrior) onEnemyHeard() {
	if !redwar.behaviour.lookingForHealing && !redwar.behaviour.attacking {
		//redwar.unit.Chat("onEnemyHeard")
		redwar.behaviour.attacking = true
		//redwar.WarBotDetectEnemy() TEMP DISABLE
		//if redwar.behaviour.listening {
		//	redwar.behaviour.listening = false
		//	redwar.unit.Chat("Wiz06a:Guard2Listen")
		//	redwar.unit.Guard(redwar.target.Pos(), redwar.target.Pos(), 300)
		//	ns.NewTimer(ns.Seconds(10), func() {
		//		redwar.behaviour.listening = true
		//	})
		//}
	}
}

func (redwar *RedWarrior) onCollide() {
	if redwar.abilities.isAlive {
		// CTF Logic.
		redwar.RedTeamPickUpBlueFlag()
		redwar.RedTeamCaptureTheBlueFlag()
		redwar.RedTeamRetrievedRedFlag()
		if !redwar.behaviour.lookingForHealing {
			//redwar.unit.Chat("onCollide")
			// TODO: determine tactic.
		}
		//if ns.GetCaller() == RedFlag {
		//	BlueTeamHasRedFlag = false
		//	BlueFlagIsAtBase = true
		//	ns.AudioEvent(audio.FlagRespawn, ns.GetHost()) // <----- replace with all players
		//	RedFlagFront.SetPos(RedFlagOutOfGame)
		//	RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
		//}

	}
}

func (redwar *RedWarrior) onEnemySighted() {
	// SCRIPT FOR WEAPON SWITCHING. On HOLD FOR NOW
	//if redwar.unit.HasItem(ns.Object("FanChakram")) {
	//	redwar.unit.Chat("HELLLLOOOOOO")
	//redwar.unit.Equip(ns.Object("FanChakram"))
	//redwar.unit.HitRanged(redwar.target.Pos())
	//}

	if !redwar.behaviour.lookingForHealing {
		//redwar.unit.Chat("onEnemySighted")
		//redwar.WarBotDetectEnemy() TEMP DISALBE
		//redwar.useWarCry()
	}
}

func (redwar *RedWarrior) onRetreat() {
	//redwar.unit.Chat("onRetreat")
	// TODO: FIX IT!
	//if redwar.behaviour.lookForHealth {
	//	redwar.behaviour.listening = false
	//	redwar.behaviour.lookForHealth = false
	//	redwar.unit.Chat("Con02A:NecroTalk02")
	//	// Walk to nearest RedPotion.
	//	redwar.targetPotion = ns.FindClosestObject(redwar.unit, ns.HasTypeName{"RedPotion"})
	//	redwar.unit.AggressionLevel(0.16)
	//	redwar.unit.Guard(redwar.targetPotion.Pos().Pos(), redwar.targetPotion.Pos(), 50)
	//	ns.NewTimer(ns.Seconds(10), func() {
	//		redwar.behaviour.lookForHealth = true
	//		redwar.behaviour.listening = true
	//	})
	//}
}

func (redwar *RedWarrior) onLostEnemy() {
	if !redwar.behaviour.lookingForHealing {
		redwar.useEyeOfTheWolf()
		//redwar.unit.Chat("onLostEnemy")
		//redwar.unit.Chat("Multi:General10")
		redwar.behaviour.attacking = false
		redwar.unit.Hunt()
	}
	redwar.RedTeamWalkToRedFlag()
}

func (redwar *RedWarrior) onHit() {
	//if redwar.unit.CurrentHealth() <= 100 && redwar.target.CurrentHealth() >= 50 && !redwar.behaviour.lookingForHealing && redwar.inventory.RedPotionInInventory <= 0 {
	//	//redwar.unit.Chat("onHit")
	//	redwar.lookForRedPotion()
	//	//redwar.unit.Guard(redwar.targetPotion.Pos().Pos(), redwar.targetPotion.Pos(), 50)
	//}
	//if redwar.unit.CurrentHealth() <= 100 && redwar.inventory.RedPotionInInventory >= 1 {
	//		for _, it := range redwar.unit.Items() {
	//			if it.Type().Name() == "RedPotion" {
	//				redwar.unit.Drop(it)
	//				redwar.inventory.RedPotionInInventory = redwar.inventory.RedPotionInInventory - 2
	//			}
	//		}
	//	}
}

func (redwar *RedWarrior) onEndOfWaypoint() {
	if redwar.behaviour.lookingForHealing {
		if redwar.unit.CurrentHealth() >= 140 {
			//redwar.unit.Chat("onEndOfWaypoint")
			redwar.unit.AggressionLevel(0.83)
			redwar.unit.Hunt()
			redwar.behaviour.lookingForHealing = false
		} else {
			if redwar.inventory.RedPotionInInventory <= 1 {
				redwar.lookForRedPotion()
			}
		}
	} else {
		if !redwar.behaviour.lookingForTarget {
			redwar.unit.Hunt()
			redwar.unit.AggressionLevel(0.83)
			redwar.behaviour.lookingForTarget = true
		}
	}
	redwar.RedTeamCheckAttackOrDefend()
}

func (redwar *RedWarrior) lookForRedPotion() {
	//if redwar.inventory.RedPotionInInventory >= 1 {
	//	redwar.onEndOfWaypoint()
	//} else {
	//	redwar.behaviour.lookingForHealing = true
	//	redwar.unit.AggressionLevel(0.16)
	//	redwar.unit.WalkTo(redwar.targetPotion.Pos())
	//}

}

func (redwar *RedWarrior) onDeath() {
	redwar.abilities.isAlive = false
	redwar.unit.DestroyChat()
	redwar.RedTeamDropFlag()
	ns.AudioEvent(audio.NPCDie, redwar.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, redwar.unit)
		redwar.unit.Delete()
		redwar.items.StreetPants.Delete()
		redwar.items.StreetSneakers.Delete()
		redwar.init()
	})
}

func (redwar *RedWarrior) Update() {
	if InitLoadComplete {
		if redwar.unit.HasEnchant(enchant.HELD) {
			ns.CastSpell(spell.SLOW, redwar.unit, redwar.unit)
			redwar.unit.EnchantOff(enchant.HELD)
		}
		redwar.findLoot()
		redwar.target = ns.FindClosestObject(redwar.unit, ns.HasClass(object.ClassPlayer))
		redwar.targetPotion = ns.FindClosestObject(redwar.unit, ns.HasTypeName{"RedPotion"})
	}
}

func (redwar *RedWarrior) findLoot() {
	const dist = 75
	// Melee weapons.
	meleeweapons := ns.FindAllObjects(
		ns.InCirclef{Center: redwar.unit, R: dist},
		ns.HasTypeName{

			"GreatSword", "WarHammer", "MorningStar", "BattleAxe", "Sword", "OgreAxe",

			//"StaffWooden",
		},
	)
	for _, item := range meleeweapons {
		if redwar.unit.CanSee(item) {
			redwar.unit.Equip(item)
		}
	}

	// Throwing weapons.
	throwingweapons := ns.FindAllObjects(
		ns.InCirclef{Center: redwar.unit, R: dist},
		ns.HasTypeName{
			"RoundChakram", "FanChakram",
		},
	)
	for _, item := range throwingweapons {
		if redwar.unit.CanSee(item) {
			redwar.unit.Pickup(item)
		}
	}

	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: redwar.unit, R: dist},
		ns.HasTypeName{
			"RedPotion", "CurePoisonPotion",
		},
	)
	for _, item := range potions {
		if redwar.unit.CanSee(item) {
			redwar.unit.Pickup(item)
		}
		if redwar.inventory.RedPotionInInventory < 3 {
			redwar.inventory.RedPotionInInventory = redwar.inventory.RedPotionInInventory + 1
		}
	}

	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: redwar.unit, R: dist},
		ns.HasTypeName{
			// Plate armor.
			"OrnateHelm", "SteelHelm", "Breastplate", "PlateLeggings", "PlateBoots", "PlateArms", "SteelShield",

			// Chainmail armor.
			"ChainCoif", "ChainTunic", "ChainLeggings",

			// Leather armor.
			"LeatherArmoBlueBoots", "LeatherArmor", "LeatherHelm", "LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if redwar.unit.CanSee(item) {
			redwar.unit.Equip(item)
		}
	}
}

func (redwar *RedWarrior) useWarCry() {
	// Check if cooldown is ready.
	if redwar.abilities.WarCryReady && !redwar.behaviour.charging {
		// Select target.
		redwar.target = ns.FindClosestObject(redwar.unit, ns.HasClass(object.ClassPlayer))
		// Trigger global cooldown.
		redwar.abilities.Ready = false
		if redwar.target.MaxHealth() == 150 {
		} else {
			redwar.abilities.WarCryReady = false
			// Check reaction time based on difficulty setting.
			ns.NewTimer(ns.Frames(redwar.reactionTime), func() {
				redwar.unit.Pause(ns.Seconds(1))
				ns.AudioEvent("WarcryInvoke", redwar.unit)
				ns.CastSpell(spell.COUNTERSPELL, redwar.unit, redwar.target)
				redwar.target.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
				redwar.unit.EnchantOff(enchant.INVULNERABLE)
				ns.NewTimer(ns.Seconds(10), func() {
					redwar.abilities.WarCryReady = true
				})
				ns.NewTimer(ns.Seconds(1), func() {
					redwar.abilities.Ready = true
				})
			})
		}
	}
}

func (redwar *RedWarrior) useEyeOfTheWolf() {
	// Check if cooldown is ready.
	if redwar.abilities.EyeOfTheWolfReady {
		// Trigger cooldown.
		redwar.abilities.EyeOfTheWolfReady = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwar.reactionTime), func() {
			// Use ability.
			redwar.unit.Enchant(enchant.INFRAVISION, ns.Seconds(10))
		})
		// Eye Of The Wolf cooldown.
		ns.NewTimer(ns.Seconds(20), func() {
			redwar.abilities.EyeOfTheWolfReady = true
		})
	}
}

// ---------------------------------- CTF BOT SCRIPT ------------------------------------//
// CTF game mechanics.
// Pick up the enemy flag.
func (redwar *RedWarrior) RedTeamPickUpBlueFlag() {
	if ns.GetCaller() == BlueFlag {
		BlueFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		RedTeamTank = redwar.unit
		RedTeamTank.AggressionLevel(0.16)
		RedTeamTank.WalkTo(RedBase.Pos())
		ns.PrintStrToAll("Team Red has the Blue flag!")
	}
}

// Capture the flag.
func (redwar *RedWarrior) RedTeamCaptureTheBlueFlag() {
	if ns.GetCaller() == RedFlag && RedFlagIsAtBase && redwar.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagCapture, RedTeamTank) // <----- replace with all players

		RedTeamTank = TeamRed
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[0].ChangeScore(+1)
		}
		FlagReset()
		redwar.unit.AggressionLevel(0.83)
		redwar.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has captured the Blue flag!")
	}
}

// Retrieve own flag.
func (redwar *RedWarrior) RedTeamRetrievedRedFlag() {
	if ns.GetCaller() == RedFlag && !RedFlagIsAtBase {
		RedFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
		redwar.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has retrieved the flag!")
		RedTeamTank.WalkTo(RedFlag.Pos())
	}
}

// Drop flag.
func (redwar *RedWarrior) RedTeamDropFlag() {
	if redwar.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		BlueFlag.Enable(true)
		RedTeamTank = TeamRed
		ns.PrintStrToAll("Team Red has dropped the Blue flag!")
	}
}

// CTF behaviour.

// When enemy is killed check to see if the flag is dropped, if so get it.
func (redwar *RedWarrior) RedTeamWalkToRedFlag() {
	if !RedFlagIsAtBase && RedFlag.IsEnabled() {
		redwar.unit.AggressionLevel(0.16)
		redwar.unit.WalkTo(BlueFlag.Pos())
	} else {
		redwar.RedTeamCheckAttackOrDefend()
	}
}

// At the end of waypoint see defend if tank, attack if not.
func (redwar *RedWarrior) RedTeamCheckAttackOrDefend() {
	if redwar.unit == RedTeamTank {
		redwar.unit.AggressionLevel(0.16)
		redwar.unit.Guard(RedBase.Pos(), RedBase.Pos(), 20)
	} else {
		redwar.unit.AggressionLevel(0.83)
		redwar.unit.WalkTo(BlueFlag.Pos())
	}
}
