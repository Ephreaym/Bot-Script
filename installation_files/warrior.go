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

func NewWarriorNoTeam() *Warrior {
	war := &Warrior{}
	war.init()
	return war
}

// Warrior bot class.
type Warrior struct {
	team                            *Team
	unit                            ns.Obj
	target                          ns.Obj
	cursor                          ns.Pointf
	berserkcursor                   ns.Obj
	vec                             ns.Pointf
	findLootT                       Updater
	weaponPreferenceT               Updater
	lookForWeaponT                  Updater
	BerserkerChargeCooldownManagerT Updater
	MissileBlockT                   Updater
	startingEquipment               struct {
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
		BerserkUpdateBool           bool
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
		listening                 bool
		lookingForHealing         bool
		charging                  bool
		attacking                 bool
		lookingForTarget          bool
		AntiStuck                 bool
		SwitchMainWeapon          bool
		Busy                      bool
		targetTeleportWake        ns.Obj
		Escorting                 bool
		Guarding                  bool
		GuardingPos               ns.Pointf
		EscortingTarget           ns.Obj
		Chatting                  bool
		VocalReady                bool
		ObjectOfInterest          ns.Obj
		LongswordAndShieldEquiped bool
		GreatswordEquiped         bool
		HammerEquiped             bool
		blinkWakeOutOfRange       bool
	}
	inventory struct {
		crown      bool
		Greatsword ns.Obj
		WarHammer  ns.Obj
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
	war.behaviour.Chatting = false
	war.behaviour.blinkWakeOutOfRange = false
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
	if TeamsEnabled {
		war.unit = ns.CreateObject("NPC", war.team.SpawnPoint())
		if GameModeIsSocial {
			war.unit.Guard(war.team.SpawnPoint(), war.team.SpawnPoint(), 20)
		}
	} else {
		randomIndex := ns.Random(0, len(botSpawnsNoTeams)-1)
		pick := botSpawnsNoTeams[randomIndex]
		war.unit = ns.CreateObject("NPC", pick.Pos())
		if GameModeIsSocial {
			war.unit.Guard(pick.Pos(), pick.Pos(), 20)
		}
	}
	if !GameModeIsSocial {
		war.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	}
	war.unit.SetMaxHealth(150)
	war.unit.SetStrength(125)
	war.unit.SetBaseSpeed(100)
	// Set Team.
	if GameModeIsCTF {
		war.unit.SetOwner(war.team.Spawns()[0])
	}
	if TeamsEnabled {
		war.unit.SetTeam(war.team.Team())
	} else {
		war.unit.SetOwner(ns.Object("WarriorOwner"))
	}
	war.unit.SetDisplayName("Warrior Bot", nil)
	if TeamsEnabled {
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
	}
	// Create WarBot mouse cursor.
	war.target = NoTarget
	war.cursor = NoTarget.Pos()
	war.berserkcursor = war.target
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	war.reactionTime = BotDifficulty
	// Set WarBot properties.
	if !GameModeIsSocial {
		war.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
		war.unit.MonsterStatusEnable(object.MonStatusCanDodge)
		war.unit.MonsterStatusEnable(object.MonStatusAlert)
	}
	war.unit.AggressionLevel(0.16)
	if !GameModeIsSocial {
		war.unit.AggressionLevel(0.83)
		war.unit.Hunt()
		war.unit.ResumeLevel(1)
		war.unit.RetreatLevel(0.0)
	}
	// Create and equip WarBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	war.startingEquipment.Longsword = ns.CreateObject("Longsword", ns.Ptf(150, 150))
	war.startingEquipment.WoodenShield = ns.CreateObject("WoodenShield", ns.Ptf(150, 150))
	war.startingEquipment.StreetSneakers = ns.CreateObject("StreetSneakers", ns.Ptf(150, 150))
	war.startingEquipment.StreetPants = ns.CreateObject("StreetPants", ns.Ptf(150, 150))
	war.behaviour.LongswordAndShieldEquiped = true
	war.unit.Equip(war.startingEquipment.Longsword)
	war.unit.Equip(war.startingEquipment.WoodenShield)
	war.unit.Equip(war.startingEquipment.StreetSneakers)
	war.unit.Equip(war.startingEquipment.StreetPants)
	// Select a WarBot loadout (tactical preference, dialog). TODO: Give different audio and chat for each set so they feel like different characters.
	// On looking for enemy.
	if !GameModeIsSocial {
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
	ns.OnChat(war.onWarCommand)
	ns.NewTimer(ns.Frames(3+war.reactionTime), func() {
		war.abilities.Ready = true
	})
}

func (war *Warrior) checkChatting() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !war.behaviour.Chatting {
		war.behaviour.Chatting = true
		ns.NewTimer(ns.Seconds(2), func() {
			war.behaviour.Chatting = false
		})
	}
}

func (war *Warrior) onChangeFocus() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	war.useHarpoon()
	war.useBerserkerCharge()
	war.useWarCry()
}

func (war *Warrior) onLookingForEnemy() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
}

func (war *Warrior) onEnemyHeard() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	war.ThrowChakram()
}

func (war *Warrior) onCollide() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
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
		if ns.GetCaller() != nil && ns.GetCaller().HasSubclass(subclass.BOMBER) && ns.GetCaller().Team() != war.unit.Team() {
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
				war.abilities.BerserkerChargeActive = false
				ns.AudioEvent(audio.BerserkerChargeOff, war.unit)
				ns.AudioEvent(audio.FleshHitFlesh, war.unit)
				ns.GetCaller().Damage(war.unit, 150, 2)
				war.StopBerserkLoop()
				if ns.GetCaller().Flags().Has(object.FlagDead) {
					war.abilities.BerserkerChargeResetOnKill = true
					ns.NewTimer(ns.Frames(war.reactionTime+3), func() {
						war.abilities.BerserkerChargeReady = true
					})
				}
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	war.target = ns.GetCaller()
	war.useHarpoon()
	war.useBerserkerCharge()
	war.useWarCry()
	war.ThrowChakram()
}

func (war *Warrior) onRetreat() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
}

func (war *Warrior) onLostEnemy() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	war.behaviour.targetTeleportWake = ns.FindClosestObject(war.unit, ns.HasTypeName{"TeleportWake"})
	if war.behaviour.targetTeleportWake != nil {
		war.behaviour.blinkWakeOutOfRange = true
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() || !war.behaviour.blinkWakeOutOfRange {
		return
	}
	if war.behaviour.targetTeleportWake != nil {
		if !(ns.InCirclef{Center: war.unit, R: 100}).Matches(war.behaviour.targetTeleportWake) {
			war.behaviour.blinkWakeOutOfRange = false
			war.unit.Attack(war.target)
			return
		}
	}
}

func (war *Warrior) onHit() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.unit.CurrentHealth() < 100 && !war.behaviour.Busy {
		war.GoToRedPotion()
	}
}

func (war *Warrior) onCheckIfObjectOfInterestIsPickedUp() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.behaviour.ObjectOfInterest == nil {
		return
	} else {
		if war.unit.HasItem(war.behaviour.ObjectOfInterest) {
			war.behaviour.ObjectOfInterest = nil
			war.onEndOfWaypoint()
		}
	}
}

func (war *Warrior) onEndOfWaypoint() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	war.behaviour.Busy = false
	war.unit.AggressionLevel(0.83)
	if GameModeIsCTF {
		war.team.CheckAttackOrDefend(war.unit)
	} else {
		war.unit.Hunt()
	}
}

func (war *Warrior) GoToRedPotion() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !war.behaviour.Busy {
		NearestRedPotion := ns.FindClosestObject(war.unit, ns.HasTypeName{"RedPotion"})
		if NearestRedPotion != nil {
			war.behaviour.Busy = true
			war.unit.AggressionLevel(0.16)
			if GameModeIsCTF {
				if war.unit == war.team.TeamTank {
					if war.unit.CanSee(NearestRedPotion) {
						war.unit.WalkTo(NearestRedPotion.Pos())
					}
				} else {
					war.unit.WalkTo(NearestRedPotion.Pos())
				}
			} else {
				war.unit.WalkTo(NearestRedPotion.Pos())
			}

		}
	}
}

func (war *Warrior) onDeath() {
	if !BotRespawn {
		BotRespawn = true
		return
	}
	war.abilities.isAlive = false
	war.StopBerserkLoop()
	war.unit.DestroyChat()
	war.unit.FlagsEnable(object.FlagNoCollide)
	if GameModeIsCTF {
		war.team.DropEnemyFlag(war.unit)
	}
	ns.AudioEvent(audio.NPCDie, war.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	if !GameModeIsCTF {
		if TeamsEnabled {
			if war.unit.HasTeam(ns.Teams()[0]) {
				ns.Teams()[1].ChangeScore(+1)
			} else {
				ns.Teams()[0].ChangeScore(+1)
			}
		}
	}
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, war.unit)
		war.unit.Delete()
		war.startingEquipment.StreetPants.Delete()
		war.startingEquipment.StreetSneakers.Delete()
		if war.unit.IsEnabled() {
			war.init()
		}
	})
}

func (war *Warrior) UsePotions() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.unit.CurrentHealth() <= 100 && war.unit.InItems().FindObjects(nil, ns.HasTypeName{"RedPotion"}) != 0 {
		ns.AudioEvent(audio.LesserHealEffect, war.unit)
		RedPotion := war.unit.Items(ns.HasTypeName{"RedPotion"})
		war.unit.SetHealth(war.unit.CurrentHealth() + 50)
		RedPotion[0].Delete()
	}
}

func (war *Warrior) Update() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !InitLoadComplete {
		return
	}
	//war.MissileBlockT.EachFrame(15, war.onDefensiveWeaponChoice)
	war.UsePotions()
	war.onCheckBlinkWakeRange()
	war.BerserkLoop()
	war.onHarpoonFlyingLoop()
	war.onHarpoonReelLoop()
	war.BerserkerChargeCooldownManagerT.EachFrame(30, war.BerserkerChargeCooldownManager)
	war.onCheckIfObjectOfInterestIsPickedUp()
	war.findLootT.EachFrame(15, war.findLoot)
	//war.weaponPreferenceT.EachFrame(90, war.WeaponPreference)
	if war.unit.HasEnchant(enchant.HELD) && !war.abilities.BerserkerStunActive && !war.abilities.BomberStunActive {
		ns.CastSpell(spell.SLOW, war.unit, war.unit)
		war.unit.EnchantOff(enchant.HELD)
	}
}

func (war *Warrior) LookForWeapon() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.inventory.Greatsword == nil && war.inventory.WarHammer == nil {
		if !war.behaviour.Busy {
			war.behaviour.Busy = true
			war.behaviour.ObjectOfInterest = ns.FindClosestObject(war.unit, ns.HasTypeName{"GreatSword", "WarHammer"})
			if war.behaviour.ObjectOfInterest != nil {
				war.unit.WalkTo(war.behaviour.ObjectOfInterest.Pos())
			}
		}
	}
}

func (war *Warrior) ThrowChakram() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.abilities.RoundChackramReady && war.unit.InItems().FindObjects(nil, ns.HasTypeName{"RoundChakram"}) != 0 {
		war.abilities.RoundChackramReady = false
		war.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				war.unit.Equip(it)
				war.unit.LookAtObject(war.target)
				war.unit.HitRanged(war.target.Pos())
				ns.NewTimer(ns.Seconds(10), func() {
					war.abilities.RoundChackramReady = true
				})
				return true
			},
			ns.HasTypeName{"RoundChakram"},
		)
	}
}

func (war *Warrior) equipLongswordAndShield() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.startingEquipment.Longsword == nil || war.startingEquipment.WoodenShield == nil {
		return
	} else if war.behaviour.LongswordAndShieldEquiped {
		return
	} else {
		war.behaviour.LongswordAndShieldEquiped = true
		war.behaviour.GreatswordEquiped = false
		war.behaviour.HammerEquiped = false
		war.unit.Equip(war.startingEquipment.Longsword)
		war.unit.Equip(war.startingEquipment.WoodenShield)
	}
}

func (war *Warrior) equipGreatsword() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.inventory.Greatsword == nil {
		return
	} else if war.behaviour.GreatswordEquiped {
		return
	} else {
		war.behaviour.LongswordAndShieldEquiped = false
		war.behaviour.GreatswordEquiped = true
		war.behaviour.HammerEquiped = false
		war.unit.Equip(war.inventory.Greatsword)
	}
}

func (war *Warrior) equipWarHammer() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.inventory.WarHammer == nil {
		return
	} else if war.behaviour.HammerEquiped {
		return
	} else {
		war.behaviour.LongswordAndShieldEquiped = false
		war.behaviour.GreatswordEquiped = false
		war.behaviour.HammerEquiped = true
		war.unit.Equip(war.inventory.WarHammer)
	}
}

func (war *Warrior) WeaponPreference() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !war.behaviour.HammerEquiped && war.inventory.WarHammer != nil {
		war.equipWarHammer()
	}
}

func (war *Warrior) onDefensiveWeaponChoice() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	// Sword and shield on stun
	if war.abilities.BerserkerStunActive && war.unit.HasEnchant(enchant.HELD) {
		war.equipLongswordAndShield()
		// Reflect Force of Nature
	} else if sp2 := ns.FindClosestObject(war.unit, ns.HasTypeName{"DeathBall", "Fireball"}, ns.InCirclef{Center: war.unit, R: 100}); sp2 != nil {
		{
			arr2 := ns.FindAllObjects(
				ns.HasTypeName{"NewPlayer", "NPC"},
			)
			for i := 0; i < len(arr2); i++ {
				if sp2.HasOwner(arr2[i]) && arr2[i].Team() != war.unit.Team() {
					// FIXME: Equip weapon
				}
			}
		}
		// Reflect Missile attacks
	} else if sp := ns.FindClosestObject(war.unit, ns.HasClass(object.ClassMissile), ns.InCirclef{Center: war.unit, R: 100}); sp != nil {
		if sp.HasOwner(war.target) {
			// FIXME: Equip weapon
		}
	}
}

func (war *Warrior) findLoot() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	const dist = 80
	// Melee weapons.
	meleeweapons := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"MorningStar", "BattleAxe", "Sword", "OgreAxe",
			//"StaffWooden",
		},
	)
	for _, item := range meleeweapons {
		if war.unit.CanSee(item) {
			war.unit.Pickup(item)
			war.unit.Equip(war.unit.GetLastItem())
		}
	}

	greatswords := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"GreatSword",
		},
	)
	for _, item := range greatswords {
		if war.unit.CanSee(item) {
			war.inventory.Greatsword = item
			war.equipGreatsword()
		}
	}

	hammers := ns.FindAllObjects(
		ns.InCirclef{Center: war.unit, R: dist},
		ns.HasTypeName{
			"WarHammer",
		},
	)
	for _, item := range hammers {
		if war.unit.CanSee(item) {
			war.inventory.WarHammer = item
			war.equipWarHammer()
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
}

func (war *Warrior) onWarCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil && !war.unit.Flags().Has(object.FlagDead) {
		switch msg {
		// Spawn commands red bots.
		// Bot commands.
		case "help", "Help", "Follow", "follow", "escort", "Escort", "come", "Come":
			if war.unit.CanSee(p.Unit()) && war.unit.Team() == p.Team() {
				war.behaviour.Escorting = true
				war.behaviour.EscortingTarget = p.Unit()
				war.behaviour.Guarding = false
				war.unit.Follow(p.Unit())
				if !war.behaviour.Chatting {
					war.checkChatting()
					random := ns.Random(1, 6)
					if random == 1 {
						war.unit.ChatStr("I'll follow you.")
					}
					if random == 2 {
						war.unit.ChatStr("Let's go.")
					}
					if random == 3 {
						war.unit.ChatStr("I'll help.")
					}
					if random == 4 {
						war.unit.ChatStr("Sure thing.")
					}
					if random == 5 {
						war.unit.ChatStr("Lead the way.")
					}
					if random == 6 {
						war.unit.ChatStr("I'll escort you.")
					}
				}

			}
		case "Attack", "Go", "go", "attack":
			if war.unit.CanSee(p.Unit()) && war.unit.Team() == p.Team() {
				war.behaviour.Escorting = false
				war.behaviour.Guarding = false
				war.unit.Hunt()
				if !war.behaviour.Chatting {
					war.checkChatting()
					random2 := ns.Random(1, 4)
					if random2 == 1 {
						war.unit.ChatStr("I'll get them.")
					}
					if random2 == 2 {
						war.unit.ChatStr("Time to shine.")
					}
					if random2 == 3 {
						war.unit.ChatStr("On the offense.")
					}
					if random2 == 4 {
						war.unit.ChatStr("Time to hunt.")
					}
				}
			}
		case "guard", "stay", "Guard", "Stay":
			if war.unit.CanSee(p.Unit()) && war.unit.Team() == p.Team() {
				war.unit.Guard(war.unit.Pos(), war.unit.Pos(), 300)
				war.behaviour.Escorting = false
				war.behaviour.Guarding = true
				war.behaviour.GuardingPos = war.unit.Pos()
				if !war.behaviour.Chatting {
					war.checkChatting()
					random1 := ns.Random(1, 4)
					if random1 == 1 {
						war.unit.ChatStr("I'll guard this place.")
					}
					if random1 == 2 {
						war.unit.ChatStr("No problem.")
					}
					if random1 == 3 {
						war.unit.ChatStr("I'll stay.")
					}
					if random1 == 4 {
						war.unit.ChatStr("I'll hold.")
					}
				}
			}
		}
	}
	return msg
}

// ------------------------------------------------------------- WARRIOR ABILITIES --------------------------------------------------------------- //

func (war *Warrior) useBerserkerCharge() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	// Check if cooldowns are ready.
	if !war.abilities.HarpoonFlying && !war.abilities.BerserkerChargeActive && war.abilities.BerserkerChargeIsEnabled && war.unit.CanSee(war.target) && war.abilities.Ready && war.abilities.BerserkerChargeReady && war.abilities.isAlive && !war.target.HasEnchant(enchant.INVULNERABLE) && !war.target.Flags().Has(object.FlagDead) {
		if GameModeIsCTF {
			if war.unit == war.team.TeamTank {
				return
			}
		}
		// Select target.
		war.cursor = war.target.Pos()
		war.vec = war.unit.Pos().Sub(war.cursor).Normalize()
		// Trigger cooldown.
		war.abilities.Ready = false
		war.abilities.BerserkerChargeReady = false
		war.abilities.BerserkerChargeActive = true
		war.abilities.BerserkerTarget = true
		// Check reaction time based on difficulty setting.
		//ns.NewTimer(ns.Frames(war.reactionTime), func() {
		if war.abilities.BerserkerChargeActive && war.abilities.isAlive {
			war.unit.EnchantOff(enchant.INVULNERABLE)
			ns.AudioEvent(audio.BerserkerChargeInvoke, war.unit)
			war.unit.LookAtObject(war.target.Pos())
			war.abilities.BerserkerChargeActive = true
			war.BerserkLoop()
			// Stop berserk if no object is hit/max range berserk.
			ns.NewTimer(ns.Seconds(3), func() {
				if war.abilities.isAlive && war.abilities.BerserkerChargeActive {
					war.StopBerserkLoop()
					war.abilities.Ready = true
				}
			})
			war.abilities.BerserkerChareCooldownTimer = 10
			war.abilities.BerserkerChargeResetOnKill = false
		}

	}
}

func (war *Warrior) BerserkerChargeCooldownManager() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.abilities.BerserkerChargeActive {
		return
	} else {
		if !war.abilities.BerserkerChargeResetOnKill {
			if war.abilities.BerserkerChareCooldownTimer == 0 {
				war.abilities.BerserkerChargeReady = true
				war.abilities.BerserkerChargeResetOnKill = false
			} else if !war.abilities.BerserkerChargeReady {
				war.abilities.BerserkerChareCooldownTimer = war.abilities.BerserkerChareCooldownTimer - 1
			}
		} else {
			return
		}
	}
}

func (war *Warrior) StopBerserkLoop() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.abilities.isAlive && war.abilities.BerserkerTarget {
		war.abilities.Ready = true
		war.abilities.BerserkerChargeActive = false
		war.abilities.BerserkerTarget = false
		war.berserkcursor.Delete()
	}
}

func (war *Warrior) BerserkLoop() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !war.abilities.BerserkerChargeActive {
		return
	}
	if war.abilities.isAlive && war.abilities.BerserkerTarget {
		war.cursor = war.berserkcursor.Pos()
		war.unit.Pause(ns.Frames(1))
		war.unit.ApplyForce(war.vec.Mul(-12))
	} else {
		war.StopBerserkLoop()
	}
}

func (war *Warrior) useWarCry() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	// Check if cooldown is war.abilities.Ready.
	if war.abilities.WarCryReady && !war.abilities.BerserkerChargeActive && !war.abilities.HarpoonFlying && !war.target.Flags().Has(object.FlagDead) && war.unit.CanSee(war.target) {
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
							if TeamsEnabled {
								if it.Team() != war.unit.Team() {
									ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
									it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
								}
							} else {
								ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
								it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
							}

						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasClass(object.ClassPlayer),
				)
				// Select target.
				// Target enemy bots.
				ns.FindObjects(
					func(it ns.Obj) bool {
						if war.unit.CanSee(it) && it.MaxHealth() < 150 && !it.HasEnchant(enchant.ANTI_MAGIC) {
							if TeamsEnabled {
								if it.Team() != war.unit.Team() {
									ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
									it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
								} else {
									ns.CastSpell(spell.COUNTERSPELL, war.unit, it)
									it.Enchant(enchant.ANTI_MAGIC, ns.Seconds(3))
								}
							}
						}
						return true
					},
					ns.InCirclef{Center: war.unit, R: 300},
					ns.HasTypeName{"NPC"},
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.abilities.HarpoonReady && !war.target.HasEnchant(enchant.INVULNERABLE) && !war.abilities.BerserkerChargeActive && !war.target.Flags().Has(object.FlagDead) && war.unit.CanSee(war.target) {
		// Create objects, set properties and shoot harpoon.
		war.abilities.HarpoonTarget = war.target
		war.abilities.HarpoonReady = false
		war.abilities.HarpoonFlying = true
		war.unit.LookAtObject(war.abilities.HarpoonTarget)
		war.abilities.HarpoonMask = ns.CreateObject("HarpoonBolt", war.unit)
		war.abilities.HarpoonMask.FlagsEnable(object.FlagNoCollide)
		ns.AudioEvent(audio.HarpoonInvoke, war.unit)
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if war.abilities.HarpoonFlying && war.abilities.isAlive {
		//ns.Effect(effect.SENTRY_RAY, war.unit.Pos(), war.abilities.HarpoonMask.Pos())
		war.abilities.HarpoonMask.PushTo(war.abilities.HarpoonTarget, -15)
		if (ns.InCirclef{Center: war.abilities.HarpoonMask, R: 50}.Matches(war.abilities.HarpoonTarget)) {
			war.abilities.HarpoonFlying = false
			war.abilities.HarpoonAttached = true
			war.onHarpoonHit()
		}
	}
}

func (war *Warrior) onHarpoonHit() {
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
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
	if war.unit.Flags().Has(object.FlagDead) || war.unit == nil || !war.unit.IsEnabled() {
		return
	}
	if !war.abilities.HarpoonAttached {
		return
	} else if war.unit.CanSee(war.abilities.HarpoonTarget) && war.abilities.HarpoonAttached && (ns.InCirclef{Center: war.unit, R: 300}.Matches(war.abilities.HarpoonTarget)) && !war.abilities.HarpoonTarget.Flags().Has(object.FlagDead) && war.abilities.isAlive {
		//ns.Effect(effect.SENTRY_RAY, war.unit.Pos(), war.abilities.HarpoonTarget.Pos())
		war.abilities.HarpoonMask.SetPos(war.abilities.HarpoonTarget.Pos())
		vec := war.abilities.HarpoonTarget.Pos().Sub(war.unit.Pos())
		war.abilities.HarpoonTarget.ApplyForce(vec.Mul(-0.03))
	} else {
		ns.AudioEvent(audio.HarpoonBroken, war.unit)
		war.abilities.HarpoonFlying = false
		war.abilities.HarpoonAttached = false
		war.abilities.HarpoonMask.Delete()
	}
}
