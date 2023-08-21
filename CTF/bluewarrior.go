package EndGameBW

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
	war := &BlueWarrior{}
	war.init()
	return war
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

func (war *BlueWarrior) init() {
	// Reset Behaviour
	war.behaviour.listening = true
	war.behaviour.attacking = false
	war.behaviour.lookingForHealing = false
	war.behaviour.charging = false
	war.behaviour.lookingForTarget = true
	// Inventory
	war.inventory.redPotionInInventory = 0
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
	war.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointBlue"))
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	war.unit.SetStrength(125)
	war.unit.SetBaseSpeed(100)
	// Set Team.
	war.unit.SetOwner(TeamBlue)
	// Create WarBot mouse cursor.
	war.target = TeamRed
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
	war.items.Longsword = ns.CreateObject("Longsword", war.unit)
	war.items.WoodenShield = ns.CreateObject("WoodenShield", war.unit)
	war.items.StreetSneakers = ns.CreateObject("StreetSneakers", war.unit)
	war.items.StreetPants = ns.CreateObject("StreetPants", war.unit)
	war.unit.Equip(war.items.Longsword)
	war.unit.Equip(war.items.WoodenShield)
	war.unit.Equip(war.items.StreetSneakers)
	war.unit.Equip(war.items.StreetPants)
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
}

func (war *BlueWarrior) onChangeFocus() {
	if !war.behaviour.lookingForHealing {
		//war.unit.Chat("onChangeFocus")
	}
}

func (war *BlueWarrior) onLookingForEnemy() {

	if !war.behaviour.lookingForHealing {
		//war.unit.Chat("onLookingForEnemy")
	}
}

func (war *BlueWarrior) onEnemyHeard() {
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

func (war *BlueWarrior) onCollide() {
	if war.abilities.isAlive {
		// CTF Logic.
		war.BlueTeamPickUpRedFlag()
		war.BlueTeamCaptureTheRedFlag()
		war.BlueTeamRetrievedBlueFlag()
		if !war.behaviour.lookingForHealing {
			//war.unit.Chat("onCollide")
			// TODO: determine tactic.
		}
		//if ns.GetCaller() == BlueFlag {
		//	RedTeamHasBlueFlag = false
		//	RedFlagIsAtBase = true
		//	ns.AudioEvent(audio.FlagRespawn, ns.GetHost()) // <----- replace with all players
		//	BlueFlagFront.SetPos(BlueFlagOutOfGame)
		//	BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
		//}

	}
}

func (war *BlueWarrior) onEnemySighted() {
	// SCRIPT FOR WEAPON SWITCHING. On HOLD FOR NOW
	//if war.unit.HasItem(ns.Object("FanChakram")) {
	//	war.unit.Chat("HELLLLOOOOOO")
	//war.unit.Equip(ns.Object("FanChakram"))
	//war.unit.HitRanged(war.target.Pos())
	//}

	if !war.behaviour.lookingForHealing {
		//war.unit.Chat("onEnemySighted")
		//war.WarBotDetectEnemy() TEMP DISALBE
		//war.useWarCry()
	}
}

func (war *BlueWarrior) onRetreat() {
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

func (war *BlueWarrior) onLostEnemy() {
	if !war.behaviour.lookingForHealing {
		war.useEyeOfTheWolf()
		//war.unit.Chat("onLostEnemy")
		//war.unit.Chat("Multi:General10")
		war.behaviour.attacking = false
		war.unit.Hunt()
	}
	war.BlueTeamWalkToBlueFlag()
}

func (war *BlueWarrior) onHit() {
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

func (war *BlueWarrior) onEndOfWaypoint() {
	if war.behaviour.lookingForHealing {
		if war.unit.CurrentHealth() >= 140 {
			//war.unit.Chat("onEndOfWaypoint")
			war.unit.AggressionLevel(0.83)
			war.unit.Hunt()
			war.behaviour.lookingForHealing = false
		} else {
			if war.inventory.redPotionInInventory <= 1 {
				war.lookForRedPotion()
			}
		}
	} else {
		if !war.behaviour.lookingForTarget {
			war.unit.Hunt()
			war.unit.AggressionLevel(0.83)
			war.behaviour.lookingForTarget = true
		}
	}
	war.BlueTeamCheckAttackOrDefend()
}

func (war *BlueWarrior) lookForRedPotion() {
	//if war.inventory.redPotionInInventory >= 1 {
	//	war.onEndOfWaypoint()
	//} else {
	//	war.behaviour.lookingForHealing = true
	//	war.unit.AggressionLevel(0.16)
	//	war.unit.WalkTo(war.targetPotion.Pos())
	//}

}

func (war *BlueWarrior) onDeath() {
	war.abilities.isAlive = false
	war.unit.DestroyChat()
	war.BlueTeamDropFlag()
	ns.AudioEvent(audio.NPCDie, war.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, war.unit)
		war.unit.Delete()
		war.items.StreetPants.Delete()
		war.items.StreetSneakers.Delete()
		war.init()
	})
}

func (war *BlueWarrior) Update() {
	if InitLoadComplete {
		if war.unit.HasEnchant(enchant.HELD) {
			ns.CastSpell(spell.SLOW, war.unit, war.unit)
			war.unit.EnchantOff(enchant.HELD)
		}
		war.findLoot()
		war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		war.targetPotion = ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
	}
}

func (war *BlueWarrior) findLoot() {
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
			"RedPotion", "CurePoisonPotion",
		},
	)
	for _, item := range potions {
		if war.unit.CanSee(item) {
			war.unit.Pickup(item)
		}
		if war.inventory.redPotionInInventory < 3 {
			war.inventory.redPotionInInventory = war.inventory.redPotionInInventory + 1
		}
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
		if war.unit.CanSee(item) {
			war.unit.Equip(item)
		}
	}
}

func (war *BlueWarrior) useWarCry() {
	// Check if cooldown is ready.
	if war.abilities.WarCryReady && !war.behaviour.charging {
		// Select target.
		war.target = ns.FindClosestObject(war.unit, ns.HasClass(object.ClassPlayer))
		// Trigger global cooldown.
		war.abilities.Ready = false
		if war.target.MaxHealth() == 150 {
		} else {
			war.abilities.WarCryReady = false
			// Check reaction time based on difficulty setting.
			ns.NewTimer(ns.Frames(war.reactionTime), func() {
				war.unit.Pause(ns.Seconds(1))
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

func (war *BlueWarrior) useEyeOfTheWolf() {
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

// ---------------------------------- CTF BOT SCRIPT ------------------------------------//
// CTF game mechanics.
// Pick up the enemy flag.
func (war *BlueWarrior) BlueTeamPickUpRedFlag() {
	if ns.GetCaller() == RedFlag {
		RedFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		BlueTeamTank = war.unit
		BlueTeamTank.AggressionLevel(0.16)
		BlueTeamTank.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Blue has the Red flag!")
	}
}

// Capture the flag.
func (war *BlueWarrior) BlueTeamCaptureTheRedFlag() {
	if ns.GetCaller() == BlueFlag && BlueFlagIsAtBase && war.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagCapture, BlueTeamTank) // <----- replace with all players
		BlueTeamTank = TeamBlue
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[1].ChangeScore(+1)
		}
		FlagReset()
		war.unit.AggressionLevel(0.83)
		war.unit.WalkTo(RedFlag.Pos())
		ns.PrintStrToAll("Team Blue has captured the Red flag!")
	}
}

// Retrieve own flag.
func (war *BlueWarrior) BlueTeamRetrievedBlueFlag() {
	if ns.GetCaller() == BlueFlag && !BlueFlagIsAtBase {
		BlueFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
		war.unit.WalkTo(BlueBase.Pos())
		ns.PrintStrToAll("Team Blue has retrieved the flag!")
		BlueTeamTank.WalkTo(BlueFlag.Pos())
	}
}

// Drop flag.
func (war *BlueWarrior) BlueTeamDropFlag() {
	if war.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		RedFlag.Enable(true)
		BlueTeamTank = TeamBlue
		ns.PrintStrToAll("Team Blue has dropped the Red flag!")
	}
}

// CTF behaviour.
// Attack enemy tank without

func (war *BlueWarrior) BlueTeamWalkToBlueFlag() {
	if !BlueFlagIsAtBase && BlueFlag.IsEnabled() {
		war.unit.AggressionLevel(0.16)
		war.unit.WalkTo(BlueFlag.Pos())
	} else {
		war.BlueTeamCheckAttackOrDefend()
	}

}

func (war *BlueWarrior) BlueTeamCheckAttackOrDefend() {
	if war.unit == BlueTeamTank {
		war.unit.AggressionLevel(0.16)
		war.unit.Guard(BlueBase.Pos(), BlueBase.Pos(), 20)
	} else {
		war.unit.AggressionLevel(0.83)
		war.unit.WalkTo(RedFlag.Pos())
	}
}
