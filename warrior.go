package BotWars

import (
	"image/color"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/class"
	"github.com/noxworld-dev/noxscript/ns/v4/damage"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/noxscript/ns/v4/subclass"
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
	berserkcursor     ns.Obj
	vec               ns.Pointf
	startingEquipment struct {
		Longsword      ns.Obj
		WoodenShield   ns.Obj
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
	}
	abilities struct {
		BerserkerChargeIsEnabled    bool
		isAlive                     bool
		Ready                       bool // Global cooldown.
		BerserkerChargeReady        bool // Cooldown is 10 seconds.
		BerserkerTarget             bool
		BerserkerChargeActive       bool
		BerserkerStunActive         bool
		BerserkerChargeResetOnKill  bool
		BerserkerChareCooldownTimer int
		BomberStunActive            bool
		WarCryReady                 bool // Cooldown is 10 seconds.
		WarCryActive                bool
		Harpoon                     ns.Obj
		HarpoonMask                 ns.Obj
		HarpoonTarget               ns.Obj
		HarpoonReady                bool
		HarpoonFlying               bool
		HarpoonAttached             bool
		EyeOfTheWolfReady           bool // Cooldown is 20 seconds.
		TreadLightlyReady           bool
		RoundChackramReady          bool // for now cooldown 10 seconds.
	}
	behaviour struct {
		listening          bool
		lookingForHealing  bool
		charging           bool
		attacking          bool
		lookingForTarget   bool
		AntiStuck          bool
		SwitchMainWeapon   bool
		Busy               bool
		targetTeleportWake ns.Obj
	}
	inventory struct {
		crown bool
	}
	reactionTime int
}

func (war *Warrior) init() {
	// TEMP bool to toggle berserk for testing.
	war.abilities.BerserkerChargeIsEnabled = true
	// Reset Behaviour
	war.behaviour.listening = true
	war.behaviour.attacking = false
	war.behaviour.lookingForHealing = false
	war.behaviour.charging = false
	war.behaviour.lookingForTarget = true
	war.behaviour.AntiStuck = true
	war.behaviour.SwitchMainWeapon = false
	war.abilities.BomberStunActive = false
	war.behaviour.Busy = false
	// Inventory
	war.inventory.crown = false
	// Reset abilities WarBot.
	war.abilities.isAlive = true
	war.abilities.Ready = false
	war.abilities.BerserkerChareCooldownTimer = 0
	war.abilities.BerserkerChargeReady = true
	war.abilities.BerserkerChargeActive = false
	war.abilities.BerserkerStunActive = false
	war.abilities.BerserkerChargeResetOnKill = false
	war.abilities.WarCryReady = true
	war.abilities.HarpoonReady = true
	war.abilities.EyeOfTheWolfReady = true
	war.abilities.TreadLightlyReady = true
	war.abilities.RoundChackramReady = true
	war.abilities.HarpoonAttached = false
	war.abilities.HarpoonFlying = false
	// Select spawnpoint.
	// Create WarBot.
	war.unit = ns.CreateObject("NPC", war.team.SpawnPoint())
	war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	war.unit.SetMaxHealth(150)
	war.unit.SetStrength(125)
	war.unit.SetBaseSpeed(100)
	// Set Team.
	if GameModeIsCTF {
		war.unit.SetOwner(war.team.Spawns()[0])
	}
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
	war.target = NoTarget
	war.cursor = NoTarget.Pos()
	war.berserkcursor = war.target
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	war.reactionTime = BotDifficulty
	// Set WarBot properties.
	war.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	war.unit.MonsterStatusEnable(object.MonStatusCanDodge)
	war.unit.MonsterStatusEnable(object.MonStatusAlert)
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
	war.findLoot()
	//war.onSlowUpdate()
	ns.NewTimer(ns.Frames(3+war.reactionTime), func() {
		war.abilities.Ready = true
	})
}

func (war *Warrior) onSlowUpdate() {
	ns.NewTimer(ns.Seconds(1), func() {

	})
}

func (war *Warrior) onChangeFocus() {
	war.useHarpoon()
	war.useBerserkerCharge()
	war.useWarCry()
}

func (war *Warrior) onLookingForEnemy() {
}

func (war *Warrior) onEnemyHeard() {
	war.ThrowChakram()
}

func (war *Warrior) onCollide() {
	if war.abilities.isAlive {
		// When the Warriors hits a wall with Berserker Charge.
		if ns.GetCaller() == nil {
			if war.abilities.BerserkerChargeActive && war.abilities.isAlive {
				war.abilities.BerserkerStunActive = true
				ns.AudioEvent(audio.BerserkerChargeOff, war.unit)
				ns.AudioEvent(audio.FleshHitStone, war.unit)
				war.unit.Enchant(enchant.HELD, ns.Seconds(2))
				ns.NewTimer(ns.Seconds(2), func() {
					war.abilities.BerserkerStunActive = false
				})
				war.StopBerserkLoop()
			}
		}
		// When the Warrior has drawn the target nearby with Harpoon.
		if ns.GetCaller() == war.abilities.HarpoonTarget && war.abilities.HarpoonAttached {
			war.abilities.HarpoonAttached = false
			war.abilities.HarpoonMask.Delete()
		}
		// CTF mechanics for flag collision.
		caller := ns.GetCaller()
		if GameModeIsCTF {
			war.team.CheckPickUpEnemyFlag(caller, war.unit)
			war.team.CheckCaptureEnemyFlag(caller, war.unit)
			war.team.CheckRetrievedOwnFlag(caller, war.unit)
		}
		// Fix to enable stun when a Warrior is hit by a Bomber.
		if ns.GetCaller() != nil && ns.GetCaller().HasSubclass(subclass.BOMBER) && ns.GetCaller().HasTeam(war.team.Enemy.team) {
			war.abilities.BomberStunActive = true
			ns.NewTimer(ns.Seconds(2), func() {
				war.abilities.BomberStunActive = false
			})
		}
		// Berserker Charge when nearby.
		if ns.GetCaller() == war.target && !war.target.Flags().Has(object.FlagDead) {
			targettime := ns.GetCaller()
			ns.NewTimer(ns.Frames(war.reactionTime*2), func() {
				if targettime == war.target && !war.target.Flags().Has(object.FlagDead) {
					war.useBerserkerCharge()
				}
			})
		}
		// Berserker Charge damage and cooldown reset after a kill with Berserker Charge.
		if war.abilities.BerserkerChargeActive && war.abilities.isAlive && !ns.GetCaller().Flags().Has(object.FlagDead) {
			if ns.GetCaller() != nil && !ns.GetCaller().Flags().Has(object.FlagDead) && war.abilities.isAlive && ns.GetCaller().HasClass(class.PLAYER) || ns.GetCaller().HasClass(class.MONSTER) {
				ns.AudioEvent(audio.BerserkerChargeOff, war.unit)
				ns.AudioEvent(audio.FleshHitFlesh, war.unit)
				ns.GetCaller().Damage(war.unit, 150, 2)
				if ns.GetCaller().Flags().Has(object.FlagDead) {
					war.abilities.BerserkerChargeResetOnKill = true
					ns.NewTimer(ns.Frames(war.reactionTime+3), func() {
						war.abilities.BerserkerChargeReady = true
					})
				}
				war.StopBerserkLoop()
			} else if ns.GetCaller() != nil && war.abilities.isAlive && ns.GetCaller().HasClass(class.IMMOBILE) && !ns.GetCaller().HasClass(class.DOOR) && !ns.GetCaller().HasClass(class.FIRE) && !ns.GetCaller().HasClass(class.MISSILE) && !ns.GetCaller().Flags().Has(object.FlagDead) {
				war.abilities.BerserkerStunActive = true
				ns.AudioEvent(audio.BerserkerChargeOff, war.unit)
				ns.AudioEvent(audio.FleshHitStone, war.unit)
				war.unit.Enchant(enchant.HELD, ns.Seconds(2))
				ns.NewTimer(ns.Seconds(2), func() {
					war.abilities.BerserkerStunActive = false
				})
				war.StopBerserkLoop()
			}
		}
	}
}

func (war *Warrior) onEnemySighted() {
	war.target = ns.GetCaller()
	war.useHarpoon()
	war.useBerserkerCharge()
	war.useWarCry()
	war.ThrowChakram()
}

func (war *Warrior) onRetreat() {
}

func (war *Warrior) onLostEnemy() {
	war.behaviour.targetTeleportWake = ns.FindClosestObject(war.unit, ns.HasTypeName{"TeleportWake"})
	if war.behaviour.targetTeleportWake != nil {
		war.onCheckBlinkWakeRange()
		war.unit.WalkTo(war.behaviour.targetTeleportWake.Pos())
	}
	ns.NewTimer(ns.Frames(15), func() {
		war.useEyeOfTheWolf()
		if GameModeIsCTF {
			war.team.WalkToOwnFlag(war.unit)
		}
	})
}

func (war *Warrior) onCheckBlinkWakeRange() {
	if !(ns.InCirclef{Center: war.unit, R: 100}).Matches(war.behaviour.targetTeleportWake) {
		war.unit.Attack(war.target)
		war.unit.DestroyChat()
		return
	}
	ns.NewTimer(ns.Frames(1), func() {
		war.onCheckBlinkWakeRange()
	})
}

func (war *Warrior) onHit() {
	if war.unit.CurrentHealth() < 100 && !war.behaviour.Busy {
		war.GoToRedPotion()
	}
}

func (war *Warrior) onEndOfWaypoint() {
	war.behaviour.Busy = false
	war.unit.AggressionLevel(0.83)
	if GameModeIsCTF {
		war.team.CheckAttackOrDefend(war.unit)
		war.LookForNearbyItems()
	} else {
		war.unit.Hunt()
		war.LookForNearbyItems()
	}

}

func (war *Warrior) GoToRedPotion() {
	if !war.behaviour.Busy {
		NearestRedPotion := ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
		if NearestRedPotion != nil {
			war.behaviour.Busy = true
			war.unit.AggressionLevel(0.16)
			if war.unit == war.team.TeamTank {
				if war.unit.CanSee(NearestRedPotion) {
					war.unit.WalkTo(NearestRedPotion.Pos())
				}
			} else {
				war.unit.WalkTo(NearestRedPotion.Pos())
			}
		}
	}
}

func (war *Warrior) onDeath() {
	war.abilities.isAlive = false
	war.StopBerserkLoop()
	war.unit.DestroyChat()
	war.unit.FlagsEnable(object.FlagNoCollide)
	war.team.DropEnemyFlag(war.unit)
	ns.AudioEvent(audio.NPCDie, war.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	if !GameModeIsCTF {
		if war.unit.HasTeam(ns.Teams()[0]) {
			ns.Teams()[1].ChangeScore(+1)
		} else {
			ns.Teams()[0].ChangeScore(+1)
		}
	}
	if !ItemDropEnabled {
		war.startingEquipment.Longsword.Delete()
		war.startingEquipment.WoodenShield.Delete()
		war.startingEquipment.StreetPants.Delete()
		war.startingEquipment.StreetSneakers.Delete()
	}
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, war.unit)
		war.unit.Delete()
		if ItemDropEnabled {
			war.startingEquipment.StreetPants.Delete()
			war.startingEquipment.StreetSneakers.Delete()

		}
		if BotRespawn {
			war.init()
		}
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
		if war.unit.HasEnchant(enchant.HELD) && !war.abilities.BerserkerStunActive && !war.abilities.BomberStunActive {
			ns.CastSpell(spell.SLOW, war.unit, war.unit)
			war.unit.EnchantOff(enchant.HELD)
		}
	}
}

func (war *Warrior) LookForWeapon() {
	if !war.behaviour.Busy {
		war.behaviour.Busy = true
		ItemLocation := ns.FindClosestObject(war.unit, ns.HasTypeName{"GreatSword", "WarHammer"})
		if ItemLocation != nil {
			war.unit.WalkTo(ItemLocation.Pos())
		}
	}
}

func (war *Warrior) LookForNearbyItems() {
	if !war.behaviour.Busy {
		war.behaviour.Busy = true
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
	ns.NewTimer(ns.Seconds(5), func() {
		// prevent bots getting stuck to stay in loop.
		if war.behaviour.AntiStuck {
			war.behaviour.AntiStuck = false
			if GameModeIsCTF {
				war.team.CheckAttackOrDefend(war.unit)
			} else {
				war.behaviour.Busy = false
				war.unit.Hunt()
				war.unit.AggressionLevel(0.83)
			}
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
			//"OrnateHelm",
			//"SteelHelm",
			"Breastplate", "PlateLeggings", "PlateBoots", "PlateArms",
			//"SteelShield",

			// Chainmail armor.
			//"ChainCoif",
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
	ns.NewTimer(ns.Frames(15), func() {
		war.findLoot()
	})
}

// ------------------------------------------------------------- WARRIOR ABILITIES --------------------------------------------------------------- //

func (war *Warrior) useBerserkerCharge() {
	// Check if cooldowns are ready.
	if !war.abilities.HarpoonFlying && war.abilities.BerserkerChargeIsEnabled && war.unit.CanSee(war.target) && war.abilities.Ready && war.abilities.BerserkerChargeReady && war.abilities.isAlive && war.unit != war.team.TeamTank && !war.target.HasEnchant(enchant.INVULNERABLE) && !war.target.Flags().Has(object.FlagDead) {
		// Select target.
		war.cursor = war.target.Pos()
		war.vec = war.unit.Pos().Sub(war.cursor).Normalize()
		// Trigger cooldown.
		war.abilities.Ready = false
		war.abilities.BerserkerChargeReady = false
		war.abilities.BerserkerChargeActive = true
		// Check reaction time based on difficulty setting.
		//ns.NewTimer(ns.Frames(war.reactionTime), func() {
		if war.abilities.BerserkerChargeActive && war.abilities.isAlive {
			war.unit.EnchantOff(enchant.INVULNERABLE)
			ns.AudioEvent(audio.BerserkerChargeInvoke, war.unit)
			war.unit.LookAtObject(war.target.Pos())
			war.abilities.BerserkerChargeActive = true
			war.BerserkLoop()
		}
		ns.NewTimer(ns.Seconds(3), func() {
			if war.abilities.isAlive {
				war.StopBerserkLoop()
				war.abilities.Ready = true
			}
		})
		// Stop berserk if no object is hit/max range berserk.
		ns.NewTimer(ns.Frames(96), func() {
			if war.abilities.BerserkerChargeActive {
				war.StopBerserkLoop()
			}
		})
		war.abilities.BerserkerChareCooldownTimer = 10
		war.abilities.BerserkerChargeResetOnKill = false
		war.BerserkerChargeCooldownManager()
	}
}

func (war *Warrior) BerserkerChargeCooldownManager() {
	if !war.abilities.BerserkerChargeResetOnKill {
		if war.abilities.BerserkerChareCooldownTimer == 0 {
			war.abilities.BerserkerChargeReady = true
			war.abilities.BerserkerChargeResetOnKill = false
		} else {
			ns.NewTimer(ns.Seconds(1), func() {
				war.abilities.BerserkerChareCooldownTimer = war.abilities.BerserkerChareCooldownTimer - 1
				war.BerserkerChargeCooldownManager()
			})
		}
	} else {
		return
	}
}

func (war *Warrior) StopBerserkLoop() {
	if war.abilities.isAlive {
		war.abilities.Ready = true
		war.abilities.BerserkerChargeActive = false
		war.abilities.BerserkerTarget = false
		war.berserkcursor.Delete()
	}
}

func (war *Warrior) BerserkLoop() {
	if war.abilities.BerserkerChargeActive && war.abilities.isAlive {
		war.cursor = war.berserkcursor.Pos()
		war.unit.Pause(ns.Frames(1))
		war.unit.ApplyForce(war.vec.Mul(-12))
		ns.NewTimer(ns.Frames(1), func() {
			war.BerserkLoop()
		})
	} else {
		war.StopBerserkLoop()
	}
}

func (war *Warrior) useWarCry() {
	// Check if cooldown is war.abilities.Ready.
	if war.abilities.WarCryReady && !war.abilities.BerserkerChargeActive && !war.abilities.HarpoonFlying {
		if war.target.MaxHealth() == 150 {
		} else {
			// Trigger global cooldown.
			war.abilities.Ready = false
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

// ------------------ Harpoon ---------------- //

func (war *Warrior) useHarpoon() {
	if war.abilities.HarpoonReady && !war.target.HasEnchant(enchant.INVULNERABLE) && !war.abilities.BerserkerChargeActive {
		// Create objects, set properties and shoot harpoon.
		war.abilities.HarpoonTarget = war.target
		war.abilities.HarpoonReady = false
		war.abilities.HarpoonFlying = true
		war.unit.LookAtObject(war.abilities.HarpoonTarget)
		ns.NewTimer(ns.Frames(1), func() {})
		war.abilities.HarpoonMask = ns.CreateObject("HarpoonBolt", war.unit)
		war.abilities.HarpoonMask.FlagsEnable(object.FlagNoCollide)
		ns.AudioEvent(audio.HarpoonInvoke, war.unit)
		war.onHarpoonFlyingLoop()
		// No target hit.
		ns.NewTimer(ns.Frames(15), func() {
			if war.abilities.HarpoonFlying {
				ns.AudioEvent(audio.HarpoonBroken, war.unit)
				war.abilities.HarpoonFlying = false
				war.abilities.HarpoonMask.Delete()
			}
		})
		// Reel max duration.
		ns.NewTimer(ns.Seconds(5), func() {
			if war.abilities.HarpoonAttached {
				ns.AudioEvent(audio.HarpoonBroken, war.unit)
				war.abilities.HarpoonAttached = false
				war.abilities.HarpoonFlying = false
				war.abilities.HarpoonMask.Delete()
			}
			war.abilities.HarpoonReady = true
		})
	}
}

func (war *Warrior) onHarpoonFlyingLoop() {
	if war.abilities.HarpoonFlying && war.abilities.isAlive {
		//ns.Effect(effect.SENTRY_RAY, war.unit.Pos(), war.abilities.HarpoonMask.Pos())
		war.abilities.HarpoonMask.PushTo(war.abilities.HarpoonTarget, -15)
		ns.NewTimer(ns.Frames(1), func() {
			war.onHarpoonFlyingLoop()
		})
		if (ns.InCirclef{Center: war.abilities.HarpoonMask, R: 50}.Matches(war.abilities.HarpoonTarget)) {
			war.abilities.HarpoonFlying = false
			war.abilities.HarpoonAttached = true
			war.onHarpoonHit()
		}
	}
}

func (war *Warrior) onHarpoonHit() {
	prev := war.abilities.HarpoonTarget.CurrentHealth()
	war.abilities.HarpoonTarget.Damage(war.unit, 1, damage.IMPALE)
	if prev != war.abilities.HarpoonTarget.CurrentHealth() {
		ns.AudioEvent(audio.HarpoonReel, war.unit)
		war.onHarpoonReelLoop()
		if war.abilities.BerserkerChargeReady {
			war.useBerserkerCharge()
		}
	} else {
		ns.AudioEvent(audio.HitMetalShield, war.unit)
		war.abilities.HarpoonFlying = false
		war.abilities.HarpoonAttached = false
		war.abilities.HarpoonMask.Delete()
	}

}

func (war *Warrior) onHarpoonReelLoop() {
	if war.unit.CanSee(war.abilities.HarpoonTarget) && war.abilities.HarpoonAttached && (ns.InCirclef{Center: war.unit, R: 300}.Matches(war.abilities.HarpoonTarget)) && !war.abilities.HarpoonTarget.Flags().Has(object.FlagDead) && war.abilities.isAlive {
		//ns.Effect(effect.SENTRY_RAY, war.unit.Pos(), war.abilities.HarpoonTarget.Pos())
		war.abilities.HarpoonMask.SetPos(war.abilities.HarpoonTarget.Pos())
		vec := war.abilities.HarpoonTarget.Pos().Sub(war.unit.Pos())
		war.abilities.HarpoonTarget.ApplyForce(vec.Mul(-0.03))
		ns.NewTimer(ns.Frames(1), func() {
			war.onHarpoonReelLoop()
		})
	} else {
		ns.AudioEvent(audio.HarpoonBroken, war.unit)
		war.abilities.HarpoonFlying = false
		war.abilities.HarpoonAttached = false
		war.abilities.HarpoonMask.Delete()
	}
}
