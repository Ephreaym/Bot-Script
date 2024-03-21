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

// NewWizard creates a new Wizard bot.
func NewWizard(t *Team) *Wizard {
	wiz := &Wizard{team: t}
	wiz.init()
	return wiz
}

func NewWizardNoTeam() *Wizard {
	wiz := &Wizard{}
	wiz.init()
	return wiz
}

// Wizard bot class.
type Wizard struct {
	team              *Team
	unit              ns.Obj
	cursor            ns.Pointf
	target            ns.Obj
	trap              ns.Obj
	mana              int
	startingEquipment struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
		WizardRobe     ns.Obj
	}
	spells struct {
		isAlive               bool
		AnchorReady           bool
		Ready                 bool
		BlinkReady            bool
		BurnReady             bool
		ConfuseReady          bool
		CounterspellReady     bool
		DeathRayReady         bool
		DrainManaReady        bool
		EnergyBoltReady       bool
		FireballReady         bool
		ForceFieldReady       bool
		FumbleReady           bool
		HasteReady            bool
		InversionReady        bool
		InvisibilityReady     bool
		LesserHealReady       bool
		MagicMissilesReady    bool
		ProtFromFireReady     bool
		ProtFromPoisonReady   bool
		ProtFromShockReady    bool
		RingOfFireReady       bool
		RunReady              bool
		ShockReady            bool
		SlowReady             bool
		TeleportToTargetReady bool
		TrapReady             bool
		TrapCount             int
		IsCastingDrainMana    bool
	}
	behaviour struct {
		Busy             bool
		AntiStuck        bool
		SwitchMainWeapon bool
		useWand          bool
		ManaOfInterest   ns.Obj
		Escorting        bool
		Guarding         bool
		GuardingPos      ns.Pointf
		EscortingTarget  ns.Obj
		Chatting         bool
	}
	reactionTime int
	audio        struct {
		ManaRestoreSound bool
	}
}

func (wiz *Wizard) init() {
	// Reset spells WizBot3.
	wiz.spells.Ready = true
	// Debuff spells.
	wiz.spells.SlowReady = true
	// Offensive spells.
	wiz.spells.MagicMissilesReady = true
	wiz.spells.TrapReady = true
	wiz.spells.DeathRayReady = true
	wiz.spells.EnergyBoltReady = true
	wiz.spells.FumbleReady = true
	wiz.spells.ConfuseReady = true
	wiz.spells.AnchorReady = true
	wiz.spells.FireballReady = true
	wiz.spells.DrainManaReady = true
	wiz.spells.BurnReady = true
	wiz.spells.RingOfFireReady = true
	wiz.spells.CounterspellReady = true
	wiz.spells.RunReady = true
	// Defensive spells.
	wiz.spells.InversionReady = true
	wiz.spells.LesserHealReady = true
	wiz.spells.BlinkReady = true
	// Buff spells
	wiz.spells.ShockReady = true
	wiz.spells.ProtFromFireReady = true
	wiz.spells.ProtFromPoisonReady = true
	wiz.spells.ProtFromShockReady = true
	wiz.spells.HasteReady = true
	wiz.spells.ForceFieldReady = true
	wiz.spells.InvisibilityReady = true
	wiz.spells.TeleportToTargetReady = true
	// Behaviour
	wiz.spells.IsCastingDrainMana = false
	wiz.behaviour.AntiStuck = true
	wiz.behaviour.SwitchMainWeapon = false
	wiz.behaviour.Busy = false
	wiz.behaviour.useWand = false
	wiz.behaviour.Chatting = false
	// Create WizBot.
	if TeamsEnabled {
		wiz.unit = ns.CreateObject("NPC", wiz.team.SpawnPoint())
	} else {
		randomIndex := ns.Random(0, len(botSpawnsNoTeams)-1)
		pick := botSpawnsNoTeams[randomIndex]
		wiz.unit = ns.CreateObject("NPC", pick.Pos())
	}
	wiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	wiz.unit.SetMaxHealth(75)
	wiz.unit.SetStrength(35)
	wiz.unit.SetBaseSpeed(83)
	wiz.spells.isAlive = true
	wiz.mana = 150
	wiz.PassiveManaRegen()
	// Set Team
	if GameModeIsCTF {
		wiz.unit.SetOwner(wiz.team.Spawns()[0])
	}
	if !TeamsEnabled {
		wiz.unit.SetDisplayName("Kirik", nil)
	} else {
		wiz.unit.SetDisplayName("Wizard Bot", nil)
	}
	if TeamsEnabled {
		wiz.unit.SetTeam(wiz.team.Team())
	} else {
		wiz.unit.SetOwner(ns.Object("WizardOwner"))
	}
	if TeamsEnabled {
		if wiz.unit.HasTeam(ns.Teams()[0]) {
			wiz.unit.SetColor(0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			wiz.unit.SetColor(1, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			wiz.unit.SetColor(2, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			wiz.unit.SetColor(3, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			wiz.unit.SetColor(4, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			wiz.unit.SetColor(5, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		} else {
			wiz.unit.SetColor(0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
			wiz.unit.SetColor(1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
			wiz.unit.SetColor(2, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
			wiz.unit.SetColor(3, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
			wiz.unit.SetColor(4, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
			wiz.unit.SetColor(5, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		}
	}
	// Create WizBot3 mouse cursor.
	wiz.target = NoTarget
	wiz.cursor = NoTarget.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	wiz.reactionTime = BotDifficulty
	// Set WizBot3 properties.
	wiz.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	wiz.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	wiz.unit.MonsterStatusEnable(object.MonStatusAlert)
	wiz.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		wiz.unit.AggressionLevel(0.83)
	})
	wiz.unit.Hunt()
	wiz.unit.ResumeLevel(0.8)
	wiz.unit.RetreatLevel(0.2)
	// Create and equip WizBot3 starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	wiz.startingEquipment.StreetSneakers = ns.CreateObject("StreetSneakers", ns.Ptf(150, 150))
	wiz.startingEquipment.StreetPants = ns.CreateObject("StreetPants", ns.Ptf(150, 150))
	wiz.startingEquipment.StreetShirt = ns.CreateObject("StreetShirt", ns.Ptf(150, 150))
	wiz.startingEquipment.WizardRobe = ns.CreateObject("WizardRobe", ns.Ptf(150, 150))
	wiz.unit.Equip(wiz.startingEquipment.StreetSneakers)
	wiz.unit.Equip(wiz.startingEquipment.StreetPants)
	wiz.unit.Equip(wiz.startingEquipment.StreetShirt)
	wiz.unit.Equip(wiz.startingEquipment.WizardRobe)
	// Buff on respawn.
	wiz.buffInitial()
	// On retreat.
	wiz.unit.OnEvent(ns.EventRetreat, wiz.onRetreat)
	// Enemy sighted.
	wiz.unit.OnEvent(ns.EventEnemySighted, wiz.onEnemySighted)
	// On heard.
	wiz.unit.OnEvent(ns.EventEnemyHeard, wiz.onEnemyHeard)
	// On collision.
	wiz.unit.OnEvent(ns.EventCollision, wiz.onCollide)
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight.
	wiz.unit.OnEvent(ns.EventLostEnemy, wiz.onLostEnemy)
	// On Death.
	wiz.unit.OnEvent(ns.EventDeath, wiz.onDeath)
	wiz.unit.OnEvent(ns.EventLookingForEnemy, wiz.onLookingForTarget)
	wiz.unit.OnEvent(ns.EventEndOfWaypoint, wiz.onEndOfWaypoint)
	wiz.unit.OnEvent(ns.EventIsHit, wiz.onHit)
	wiz.LookForWeapon()
	wiz.WeaponPreference()
	ns.OnChat(wiz.onWizCommand)
	wiz.findLoot()
	// CODE FOR NEW TESTING //!!!! ONWAND
	//wiz.unit.MonsterStatusEnable(object.MonsterStatus(object.MonsterImmuneFear))
	//ns.CreateObject("FireStormWand", wiz.unit.Pos())
}

func (wiz *Wizard) checkChatting() {
	if !wiz.behaviour.Chatting {
		wiz.behaviour.Chatting = true
		ns.NewTimer(ns.Seconds(2), func() {
			wiz.behaviour.Chatting = false
		})
	}
}

func (wiz *Wizard) WeaponPreference() {
	// Priority list to get the prefered weapon.
	// TODO: Add stun and range conditions.
	if wiz.unit.InItems().FindObjects(nil, ns.HasTypeName{"FireStormWand"}) != 0 && wiz.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"FireStormWand"}) == 0 {
		wiz.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				wiz.unit.Equip(it)
				return true
			},
			ns.HasTypeName{"FireStormWand"},
		)
	} else if wiz.unit.InItems().FindObjects(nil, ns.HasTypeName{"ForceWand"}) != 0 && wiz.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"ForceWand"}) == 0 {
		wiz.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				wiz.unit.Equip(it)
				return true
			},
			ns.HasTypeName{"ForceWand"},
		)
	}
	ns.NewTimer(ns.Seconds(10), func() {
		wiz.WeaponPreference()
	})
}

func (wiz *Wizard) onHit() {
	if wiz.mana <= 20 && !wiz.behaviour.Busy {
		wiz.GoToManaObelisk()
	}
}

func (wiz *Wizard) UsePotions() {
	if wiz.unit.CanSee(wiz.target) {
		if wiz.unit.CurrentHealth() <= 25 && wiz.unit.InItems().FindObjects(nil, ns.HasTypeName{"RedPotion"}) != 0 {
			ns.AudioEvent(audio.LesserHealEffect, wiz.unit)
			RedPotion := wiz.unit.Items(ns.HasTypeName{"RedPotion"})
			wiz.unit.SetHealth(wiz.unit.CurrentHealth() + 50)
			RedPotion[0].Delete()
		}
		if wiz.mana <= 100 && wiz.unit.InItems().FindObjects(nil, ns.HasTypeName{"BluePotion"}) != 0 {
			wiz.mana = wiz.mana + 50
			ns.AudioEvent(audio.RestoreMana, wiz.unit)
			BluePotion := wiz.unit.Items(ns.HasTypeName{"BluePotion"})
			BluePotion[0].Delete()
		}
	}
}

func (wiz *Wizard) onEndOfWaypoint() {
	wiz.behaviour.Busy = false
	wiz.unit.AggressionLevel(0.83)
	if wiz.mana <= 49 {
		wiz.GoToManaObelisk()
	} else {
		if GameModeIsCTF {
			wiz.team.CheckAttackOrDefend(wiz.unit)
		} else {
			wiz.unit.WalkTo(wiz.target.Pos())
			ns.NewTimer(ns.Seconds(2), func() {
				wiz.unit.Hunt()
			})
		}
	}
	wiz.LookForNearbyItems()

}

func (wiz *Wizard) buffInitial() {
	wiz.castForceField()
}

func (wiz *Wizard) onLookingForTarget() {
}

func (wiz *Wizard) onEnemyHeard() {
	if !wiz.unit.CanSee(wiz.target) {
		wiz.castFireballAtHeard()
		wiz.castInvisibility()
	}
}

func (wiz *Wizard) onEnemySighted() {
	wiz.target = ns.GetCaller()
	if !wiz.unit.HasEnchant(enchant.INVISIBLE) {
		wiz.castSlow()
	}
}

func (wiz *Wizard) onCollide() {
	if wiz.spells.isAlive {
		caller := ns.GetCaller()
		if GameModeIsCTF {
			wiz.team.CheckPickUpEnemyFlag(caller, wiz.unit)
			wiz.team.CheckCaptureEnemyFlag(caller, wiz.unit)
			wiz.team.CheckRetrievedOwnFlag(caller, wiz.unit)
		}
		if caller == wiz.behaviour.ManaOfInterest {
			ns.NewTimer(ns.Seconds(1), func() {
				if wiz.mana > 140 {
					wiz.onEndOfWaypoint()
				} else {
					wiz.GoToManaObelisk()
				}
			})
		}
	}
}

func (wiz *Wizard) onRetreat() {
	wiz.castBlink()
}

func (wiz *Wizard) onLostEnemy() {

	if GameModeIsCTF {
		wiz.team.WalkToOwnFlag(wiz.unit)
	}
}

func (wiz *Wizard) onDeath() {
	wiz.spells.isAlive = false
	wiz.spells.Ready = false
	wiz.unit.FlagsEnable(object.FlagNoCollide)
	if GameModeIsCTF {
		wiz.team.DropEnemyFlag(wiz.unit)
	}
	wiz.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, wiz.unit)
	if !GameModeIsCTF {
		if TeamsEnabled {
			if wiz.unit.HasTeam(ns.Teams()[0]) {
				ns.Teams()[1].ChangeScore(+1)
			} else {
				ns.Teams()[0].ChangeScore(+1)
			}
		}
	}
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, wiz.unit)
		wiz.unit.Delete()
		wiz.startingEquipment.StreetPants.Delete()
		wiz.startingEquipment.StreetSneakers.Delete()
		wiz.startingEquipment.StreetShirt.Delete()
		if BotRespawn {
			wiz.init()
		}
	})
}

func (wiz *Wizard) PassiveManaRegen() {
	if wiz.spells.isAlive {
		ns.NewTimer(ns.Seconds(2), func() {
			if wiz.mana < 150 {
				if !BotMana {
					wiz.mana = wiz.mana + 300
				}
				wiz.mana = wiz.mana + 1
			}
			wiz.PassiveManaRegen()
		})
	}
}

func (wiz *Wizard) GoToManaObelisk() {
	if !wiz.behaviour.Busy {
		wiz.behaviour.Busy = true
		wiz.unit.AggressionLevel(0.16)
		NearestObeliskWithMana := ns.FindClosestObjectIn(wiz.unit, ns.Objects(AllManaObelisksOnMap),
			ns.ObjCondFunc(func(it ns.Obj) bool {
				return it.CurrentMana() >= 10
			}),
		)
		if NearestObeliskWithMana != nil {
			wiz.behaviour.ManaOfInterest = NearestObeliskWithMana
			if GameModeIsCTF {
				if wiz.unit == wiz.team.TeamTank {
					if wiz.unit.CanSee(NearestObeliskWithMana) {
						wiz.unit.WalkTo(NearestObeliskWithMana.Pos())
					}
				} else {
					wiz.unit.WalkTo(NearestObeliskWithMana.Pos())
				}
			} else {
				wiz.unit.WalkTo(NearestObeliskWithMana.Pos())
			}
		}
	}
}

func (wiz *Wizard) RestoreMana() {
	if wiz.mana < 150 {
		for i := 0; i < len(AllManaObelisksOnMap); i++ {
			if AllManaObelisksOnMap[i].CurrentMana() > 0 && wiz.unit.CanSee(AllManaObelisksOnMap[i]) && (ns.InCirclef{Center: wiz.unit, R: 50}).Matches(AllManaObelisksOnMap[i]) {
				wiz.mana = wiz.mana + 1
				AllManaObelisksOnMap[i].SetMana(AllManaObelisksOnMap[i].CurrentMana() - 1)
				wiz.RestoreManaSound()
			}
		}
	}
}

func (wiz *Wizard) RestoreManaWithDrainMana() {
	if wiz.mana < 150 && wiz.spells.IsCastingDrainMana {
		for i := 0; i < len(AllManaObelisksOnMap); i++ {
			if AllManaObelisksOnMap[i].CurrentMana() > 0 && wiz.unit.CanSee(AllManaObelisksOnMap[i]) && (ns.InCirclef{Center: wiz.unit, R: 200}).Matches(AllManaObelisksOnMap[i]) {
				wiz.mana = wiz.mana + 1
				AllManaObelisksOnMap[i].SetMana(AllManaObelisksOnMap[i].CurrentMana() - 1)
				wiz.RestoreManaSound()
			}
		}
		ManaSourceEnemyPlayer := ns.FindAllObjects(
			ns.HasClass(object.ClassPlayer),
			ns.InCirclef{Center: wiz.unit, R: 200},
		)
		for i := 0; i < len(ManaSourceEnemyPlayer); i++ {
			if ManaSourceEnemyPlayer[i].CurrentMana() > 0 && wiz.unit.CanSee(ManaSourceEnemyPlayer[i]) && ManaSourceEnemyPlayer[i].MaxHealth() <= 100 && ManaSourceEnemyPlayer[i].Team() != wiz.unit.Team() {
				wiz.mana = wiz.mana + 1
				ManaSourceEnemyPlayer[i].SetMana(ManaSourceEnemyPlayer[i].CurrentMana() - 1)
				wiz.RestoreManaSound()
			}
		}
		ManaSourceEnemyNPC := ns.FindAllObjects(
			ns.HasTypeName{"NPC"},
			ns.InCirclef{Center: wiz.unit, R: 200},
		)
		for i := 0; i < len(ManaSourceEnemyNPC); i++ {
			if ManaSourceEnemyNPC[i].CurrentMana() > 0 && wiz.unit.CanSee(ManaSourceEnemyNPC[i]) && ManaSourceEnemyNPC[i].MaxHealth() <= 100 && ManaSourceEnemyPlayer[i].Team() != wiz.unit.Team() {
				wiz.mana = wiz.mana + 1
				ManaSourceEnemyNPC[i].SetMana(ManaSourceEnemyNPC[i].CurrentMana() - 1)
				wiz.RestoreManaSound()
			}
		}
	}
}

func (wiz *Wizard) RestoreManaSound() {
	if !wiz.audio.ManaRestoreSound {
		wiz.castDrainMana()
		wiz.audio.ManaRestoreSound = true
		ns.AudioEvent(audio.RestoreMana, wiz.unit)
		ns.NewTimer(ns.Frames(15), func() {
			wiz.audio.ManaRestoreSound = false
		})
	}
}

func (wiz *Wizard) checkForMissiles() {
	// Maybe need to add a ns.hasteam condition. Not sure yet.
	if sp2 := ns.FindClosestObject(wiz.unit, ns.HasTypeName{"DeathBall"}, ns.InCirclef{Center: wiz.unit, R: 500}); sp2 != nil {
		{
			arr2 := ns.FindAllObjects(
				ns.HasTypeName{"NewPlayer", "NPC"},
			)
			for i := 0; i < len(arr2); i++ {
				if sp2.HasOwner(arr2[i]) && arr2[i].Team() != wiz.unit.Team() {
					wiz.castCounterspellAtForceOfNature()
				}
			}
		}
	} else {
		if sp := ns.FindClosestObject(wiz.unit, ns.HasClass(object.ClassMissile), ns.InCirclef{Center: wiz.unit, R: 500}); sp != nil {
			if sp.HasOwner(wiz.target) {
				wiz.castInversion()
			}
		}
	}
}

func (wiz *Wizard) Update() {
	wiz.checkForMissiles()
	wiz.UsePotions()
	wiz.RestoreMana()
	wiz.RestoreManaWithDrainMana()
	if wiz.mana > 150 {
		wiz.mana = 150
	}
	if wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
		wiz.spells.Ready = true
		//wiz.useWand()
	}
	if wiz.target.HasEnchant(enchant.HELD) || wiz.target.HasEnchant(enchant.SLOWED) || wiz.unit.HasEnchant(enchant.INVISIBLE) {
		if wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
			wiz.castDeathRay()
		}
	}
	if wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
		wiz.castFireball()
		if !wiz.unit.HasEnchant(enchant.INVISIBLE) {
			wiz.castBurn()
			wiz.castRingOfFire()
			wiz.castSlow()
			wiz.castEnergyBolt()
			wiz.castMissilesOfMagic()
			wiz.castCounterspell()
			//wiz.castConfuse()
			//wiz.castAnchor()
			//wiz.castRun()
			//wiz.castFumble()
			//wiz.castTeleportToTarget()
		}
		if wiz.target.MaxHealth() == 75 || wiz.target.MaxHealth() == 100 && (ns.InCirclef{Center: wiz.unit, R: 200}).Matches(wiz.target) {
			wiz.castDrainMana()
		}
	}
	if wiz.spells.Ready {
		if !wiz.unit.HasEnchant(enchant.INVISIBLE) {
			wiz.castForceField()
			if wiz.mana >= 140 {
				wiz.castLesserHeal()
				wiz.castHaste()
				wiz.castShock()
			}
		}
	}
	if !wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
		if wiz.mana >= 140 {
			wiz.castProtectionFromShock()
			wiz.castProtectionFromFire()
			wiz.castInvisibility()
			wiz.castTrap()
		}
	}
	if wiz.unit.HasEnchant(enchant.HELD) {
		wiz.castBlink()
	}
	if !wiz.unit.HasEnchant(enchant.SHIELD) || !wiz.unit.HasEnchant(enchant.HASTED) || !wiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) || !wiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) || wiz.spells.TrapCount <= 3 {
		wiz.GoToManaObelisk()
	}
}

func (wiz *Wizard) LookForWeapon() {
	if !wiz.behaviour.Busy {
		wiz.behaviour.Busy = true
		ItemLocation := ns.FindClosestObject(wiz.unit, ns.HasTypeName{"FireStormWand", "LesserFireballWand", "ForceWand"})
		if ItemLocation != nil {
			wiz.unit.WalkTo(ItemLocation.Pos())
		}
	}
}

func (wiz *Wizard) LookForNearbyItems() {
	if !wiz.behaviour.Busy {
		wiz.behaviour.Busy = true
		if ns.FindAllObjects(ns.HasTypeName{
			"RedPotion", "FireStormWand", "LesserFireballWand", "ForceWand", "CurePoisonPotion", "WizardHelm", "WizardRobe", "BluePotion", "LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
			ns.InCirclef{Center: wiz.unit, R: 200}) != nil {
			if wiz.unit.InItems().FindObjects(nil, ns.HasTypeName{"RedPotion", "FireStormWand", "LesserFireballWand", "ForceWand", "CurePoisonPotion", "WizardHelm", "WizardRobe", "BluePotion", "LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"}) == 0 {
				ItemLocation := ns.FindAllObjects(ns.HasTypeName{
					"RedPotion", "CurePoisonPotion", "WizardHelm", "WizardRobe", "BluePotion", "FireStormWand", "LesserFireballWand", "ForceWand", "LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
					ns.InCirclef{Center: wiz.unit, R: 200},
				)
				if wiz.unit.CanSee(ItemLocation[0]) {
					wiz.unit.WalkTo(ItemLocation[0].Pos())
				}
			}
		}
	}
	ns.NewTimer(ns.Seconds(5), func() {
		// prevent bots getting stuck to stay in loop.
		if wiz.behaviour.AntiStuck {
			wiz.behaviour.AntiStuck = false
			if GameModeIsCTF {
				wiz.team.CheckAttackOrDefend(wiz.unit)
			} else {
				wiz.behaviour.Busy = false
				wiz.unit.Hunt()
				wiz.unit.AggressionLevel(0.83)
			}
			ns.NewTimer(ns.Seconds(6), func() {
				wiz.behaviour.AntiStuck = true
			})
		}
	})
}

func (wiz *Wizard) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			"DeathRayWand",
			"FireStormWand",
			"LesserFireballWand",
			"ForceWand",
			//"SulphorousShowerWand"
			//"SulphorousFlareWand"
			//"StaffWooden",
		},
	)
	for _, item := range weapons {
		if wiz.unit.CanSee(item) {
			wiz.unit.Pickup(item)
			wiz.unit.Equip(wiz.unit.GetLastItem())
		}
	}

	DeathRayWand := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			"DeathRayWand",
		},
	)
	for _, item := range DeathRayWand {
		if wiz.unit.CanSee(item) {
			wiz.unit.Pickup(item)
		}
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Armor.
			//"WizardHelm",
			"WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if wiz.unit.CanSee(item) {
			wiz.unit.Pickup(item)
			wiz.unit.Equip(wiz.unit.GetLastItem())
		}
	}

	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			"RedPotion", "CurePoisonPotion", "BluePotion",
		},
	)
	for _, item := range potions {
		if wiz.unit.CanSee(item) {
			wiz.unit.Pickup(item)
		}
	}
	ns.NewTimer(ns.Frames(15), func() {
		wiz.findLoot()
	})
}

// Checks the ammount of traps active for the WIzard bot.
func (wiz *Wizard) checkTrapCount() {
	allTraps := ns.FindAllObjects(ns.HasTypeName{"Glyph"}, ns.ObjCondFunc(func(it ns.Obj) bool {
		return it.HasOwner(wiz.unit)
	}))
	wiz.spells.TrapCount = 0
	if allTraps == nil {
	} else {
		for i := 0; i < len(allTraps); i++ {
			wiz.spells.TrapCount = wiz.spells.TrapCount + 1
			if wiz.spells.TrapCount == 4 {
				ns.NewTimer(ns.Seconds(5), func() {
					wiz.checkTrapCount()
				})
			}
		}
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //

func (wiz *Wizard) castTrap() {
	if wiz.mana >= 105 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.TrapReady && wiz.spells.TrapCount <= 3 {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Ring of Fire chant.
				castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhDown, PhDownLeft, PhUp}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Magic Missiles chant.
							castPhonemes(wiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Shock chant.
										castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(wiz.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
															wiz.spells.TrapReady = false
															ns.AudioEvent(audio.TrapDrop, wiz.unit)
															wiz.mana = wiz.mana - 105
															wiz.trap = ns.NewTrap(wiz.unit, spell.CLEANSING_FLAME, spell.MAGIC_MISSILE, spell.SHOCK)
															wiz.trap.SetOwner(wiz.unit)
															// Global cooldown.
															ns.NewTimer(ns.Frames(15), func() {
																wiz.checkTrapCount()
																wiz.spells.Ready = true
															})
															// Trap cooldown.
															ns.NewTimer(ns.Seconds(5), func() {
																wiz.spells.TrapReady = true
															})
														}
													})
												}
											})
										})
									}
								})
							})
						}
					})
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castShock() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.SHOCK) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.spells.Ready && wiz.spells.ShockReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.ShockReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.SHOCK, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Shock cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.ShockReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castRingOfFire() {
	// Check if cooldowns are ready.
	if wiz.mana >= 60 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && !wiz.target.HasEnchant(enchant.INVULNERABLE) && wiz.spells.Ready && wiz.spells.RingOfFireReady && (ns.InCirclef{Center: wiz.unit, R: 40}).Matches(wiz.target) {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhDown, PhDownLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.RingOfFireReady = false
						wiz.mana = wiz.mana - 60
						ns.CastSpell(spell.CLEANSING_FLAME, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Ring of Fire cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.ShockReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castLesserHeal() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CurrentHealth() <= 60 && wiz.spells.Ready && wiz.spells.LesserHealReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhUp, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.LesserHealReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.LESSER_HEAL, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Shock cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							wiz.spells.LesserHealReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castInversion() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.InversionReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhUpLeft, FPhUpRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.InversionReady = false
						wiz.mana = wiz.mana - 10
						ns.CastSpell(spell.INVERSION, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Inversion cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							wiz.spells.InversionReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castInvisibility() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.spells.Ready && wiz.spells.InvisibilityReady {
		if GameModeIsCTF {
			if wiz.unit == wiz.team.TeamTank {
				return
			}
		}
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.InvisibilityReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.INVISIBILITY, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Invisibility cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.InvisibilityReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castEnergyBolt() {
	// Check if cooldowns are ready.
	if wiz.mana > 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && !wiz.target.HasEnchant(enchant.INVULNERABLE) && wiz.unit.CanSee(wiz.target) && wiz.spells.EnergyBoltReady && wiz.spells.Ready && (ns.InCirclef{Center: wiz.unit, R: 200}).Matches(wiz.target) {
		// Select target.
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(30))
						wiz.spells.EnergyBoltReady = false
						ns.CastSpell(spell.LIGHTNING, wiz.unit, wiz.target)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Energy Bolt cooldown.
						ns.NewTimer(ns.Seconds(3), func() {
							wiz.spells.EnergyBoltReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castDeathRay() {
	// Check if cooldowns are ready.
	if wiz.mana >= 60 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.DeathRayReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.INVULNERABLE) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
		// Select target.
		wiz.cursor = wiz.target.Pos()
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDownRight, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.DeathRayReady = false
						ns.CastSpell(spell.DEATH_RAY, wiz.unit, wiz.cursor)
						wiz.mana = wiz.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Death Ray cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.DeathRayReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castBurn() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.target.HasEnchant(enchant.INVULNERABLE) && wiz.spells.BurnReady && wiz.spells.Ready && wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && !wiz.target.HasEnchant(enchant.INVULNERABLE) {
		// Select target.
		wiz.cursor = wiz.target.Pos()
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhUp, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.BurnReady = false
						ns.CastSpell(spell.BURN, wiz.unit, wiz.cursor)
						wiz.mana = wiz.mana - 10
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Burn cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.BurnReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castFireball() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.target.HasEnchant(enchant.INVULNERABLE) && wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
		// Select target.
		wiz.cursor = wiz.target.Pos()
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.FireballReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.FIREBALL, wiz.unit, wiz.cursor)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.FireballReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castFireballAtHeard() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.FireballReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.FIREBALL, wiz.unit, wiz.target.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.FireballReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castBlink() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.BlinkReady {
		if GameModeIsCTF {
			if wiz.unit == wiz.team.TeamTank {
				return
			}
		}
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.BlinkReady = false
						wiz.mana = wiz.mana - 10
						ns.NewTrap(wiz.unit, spell.BLINK)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							wiz.spells.BlinkReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castMissilesOfMagic() {
	// Check if cooldowns are ready.
	if wiz.mana >= 15 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.target.HasEnchant(enchant.INVULNERABLE) && wiz.unit.CanSee(wiz.target) && wiz.spells.MagicMissilesReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
		// Select target.
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.MagicMissilesReady = false
						ns.CastSpell(spell.MAGIC_MISSILE, wiz.unit, wiz.target)
						wiz.mana = wiz.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Missiles Of Magic cooldown.
						ns.NewTimer(ns.Seconds(3), func() {
							wiz.spells.MagicMissilesReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castSlow() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.SlowReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.SLOWED) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
		// Select target.
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.SlowReady = false
						wiz.mana = wiz.mana - 10
						ns.CastSpell(spell.SLOW, wiz.unit, wiz.target)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(3), func() {
							wiz.spells.SlowReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

//func (wiz *Wizard) castConfuse() {
//	// Check if cooldowns are ready.
//	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.ConfuseReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.CONFUSED) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
//		// Select target.
//		// Trigger cooldown.
//		wiz.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//			// Check for War Cry before chant.
//			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(wiz.unit, []audio.Name{PhRight, PhLeft, PhRight}, func() {
//					// Check for War Cry before spell release.
//					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//						// Aim.
//						wiz.unit.LookAtObject(wiz.target)
//						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
//						wiz.spells.ConfuseReady = false
//						wiz.mana = wiz.mana - 10
//						ns.CastSpell(spell.CONFUSE, wiz.unit, wiz.target)
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.Ready = true
//						})
//						// Confuse cooldown.
//						ns.NewTimer(ns.Seconds(3), func() {
//							wiz.spells.ConfuseReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//					wiz.spells.Ready = true
//				})
//			}
//		})
//	}
//}
//
//func (wiz *Wizard) castTeleportToTarget() {
//	// Check if cooldowns are ready.
//	if wiz.mana >= 20 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.TeleportToTargetReady && wiz.spells.Ready {
//		// Select target.
//		wiz.cursor = wiz.target.Pos()
//		// Trigger cooldown.
//		wiz.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//			// Check for War Cry before chant.
//			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(wiz.unit, []audio.Name{PhDown, PhUp, PhLeft, PhRight}, func() {
//					// Check for War Cry before spell release.
//					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) {
//						// Aim.
//						wiz.unit.LookAtObject(wiz.target)
//						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
//						wiz.spells.TeleportToTargetReady = false
//						ns.CastSpell(spell.TELEPORT_TO_TARGET, wiz.unit, wiz.cursor)
//						wiz.mana = wiz.mana - 20
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.Ready = true
//						})
//						// Teleport To Target cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.TeleportToTargetReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//					wiz.spells.Ready = true
//				})
//			}
//		})
//	}
//}
//
//func (wiz *Wizard) castAnchor() {
//	// Check if cooldowns are ready.
//	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.AnchorReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.ANCHORED) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
//		// Select target.
//		// Trigger cooldown.
//		wiz.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//			// Check for War Cry before chant.
//			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDown}, func() {
//					// Check for War Cry before spell release.
//					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//						// Aim.
//						wiz.unit.LookAtObject(wiz.target)
//						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
//						wiz.spells.AnchorReady = false
//						wiz.mana = wiz.mana - 10
//						ns.CastSpell(spell.ANCHOR, wiz.unit, wiz.target)
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.Ready = true
//						})
//						// Slow cooldown.
//						ns.NewTimer(ns.Seconds(3), func() {
//							wiz.spells.AnchorReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//					wiz.spells.Ready = true
//				})
//			}
//		})
//	}
//}
//
//func (wiz *Wizard) castRun() {
//	// Check if cooldowns are ready.
//	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.RunReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.RUN) && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) {
//		// Select target.
//		// Trigger cooldown.
//		wiz.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//			// Check for War Cry before chant.
//			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(wiz.unit, []audio.Name{PhRight, PhRight, PhRight}, func() {
//					// Check for War Cry before spell release.
//					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//						// Aim.
//						wiz.unit.LookAtObject(wiz.target)
//						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
//						wiz.spells.RunReady = false
//						wiz.mana = wiz.mana - 10
//						ns.CastSpell(spell.RUN, wiz.unit, wiz.target)
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.Ready = true
//						})
//						// Run cooldown.
//						ns.NewTimer(ns.Seconds(3), func() {
//							wiz.spells.RunReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//					wiz.spells.Ready = true
//				})
//			}
//		})
//	}
//}
//
//func (wiz *Wizard) castFumble() {
//	// Check if cooldowns are ready.
//	if wiz.mana >= 60 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.FumbleReady && wiz.spells.Ready && !wiz.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && wiz.target.HasClass(class.PLAYER) {
//		// Select target.
//		// Trigger cooldown.
//		wiz.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//			// Check for War Cry before chant.
//			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhDown, PhLeft, PhRight, PhDown}, func() {
//					// Check for War Cry before spell release.
//					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
//						// Aim.
//						wiz.unit.LookAtObject(wiz.target)
//						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
//						wiz.spells.FumbleReady = false
//						wiz.mana = wiz.mana - 60
//						ns.CastSpell(spell.FUMBLE, wiz.unit, wiz.target)
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							wiz.spells.Ready = true
//						})
//						// Fumble cooldown.
//						ns.NewTimer(ns.Seconds(3), func() {
//							wiz.spells.FumbleReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
//					wiz.spells.Ready = true
//				})
//			}
//		})
//	}
//}

func (wiz *Wizard) castDrainMana() {
	// Check if cooldowns are ready.
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.DrainManaReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhUp, PhUpLeft, PhDown, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.DrainManaReady = false
						wiz.spells.IsCastingDrainMana = true
						ManaSource := ns.FindClosestObject(wiz.unit,
							ns.HasTypeName{
								"NewPlayer", "NPC", "ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12",
							},
							ns.InCirclef{Center: wiz.unit, R: 200})
						ns.NewTimer(ns.Frames(30), func() {
							wiz.spells.IsCastingDrainMana = false
						})
						ns.CastSpell(spell.DRAIN_MANA, wiz.unit, wiz.unit.Pos())
						wiz.unit.LookAtObject(ManaSource)
						wiz.unit.Pause(ns.Frames(30))
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(3), func() {
							wiz.spells.DrainManaReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castHaste() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.HASTED) && wiz.spells.Ready && wiz.spells.HasteReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.HasteReady = false
						wiz.mana = wiz.mana - 10
						ns.CastSpell(spell.HASTE, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							wiz.spells.HasteReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castCounterspell() {
	// Check if cooldowns are ready.
	if wiz.mana >= 20 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.target.HasEnchant(enchant.SHOCK) && wiz.spells.Ready && wiz.spells.CounterspellReady && wiz.unit.CanSee(wiz.target) {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) {
						wiz.spells.CounterspellReady = false
						wiz.mana = wiz.mana - 20
						ns.CastSpell(spell.COUNTERSPELL, wiz.unit, wiz.unit.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							wiz.spells.CounterspellReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castCounterspellAtForceOfNature() {
	// Check if cooldowns are ready.
	if wiz.mana >= 20 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.CounterspellReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhDown, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.CounterspellReady = false
						wiz.mana = wiz.mana - 20
						ns.CastSpell(spell.COUNTERSPELL, wiz.unit, wiz.unit.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							wiz.spells.CounterspellReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castForceField() {
	// if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.SHIELD)
	// Check if cooldowns are ready.
	if wiz.mana >= 80 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.SHIELD) && wiz.spells.Ready && wiz.spells.ForceFieldReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.ForceFieldReady = false
						wiz.mana = wiz.mana - 80
						ns.CastSpell(spell.SHIELD, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Force Field cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							wiz.spells.ForceFieldReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) && wiz.spells.Ready && wiz.spells.ProtFromFireReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.ProtFromFireReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.ProtFromFireReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

func (wiz *Wizard) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) && wiz.spells.Ready && wiz.spells.ProtFromShockReady {
		// Trigger cooldown.
		wiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
			// Check for War Cry before chant.
			if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(wiz.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						wiz.spells.ProtFromShockReady = false
						wiz.mana = wiz.mana - 30
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, wiz.unit, wiz.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							wiz.spells.Ready = true
						})
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.ProtFromShockReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
					wiz.spells.Ready = true
				})
			}
		})
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //

func (wiz *Wizard) onWizCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil {
		switch msg {
		// Spawn commands red bots.
		// Bot commands.
		case "help", "Help", "Follow", "follow", "escort", "Escort", "come", "Come":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.Team() == p.Team() {
				wiz.behaviour.Escorting = true
				wiz.behaviour.EscortingTarget = p.Unit()
				wiz.behaviour.Guarding = false
				wiz.unit.Follow(p.Unit())
				if !wiz.behaviour.Chatting {
					wiz.checkChatting()
					random := ns.Random(1, 4)
					if random == 1 {
						wiz.unit.ChatStr("I'll follow you.")
					}
					if random == 2 {
						wiz.unit.ChatStr("Let's go.")
					}
					if random == 3 {
						wiz.unit.ChatStr("I'll help.")
					}
					if random == 4 {
						wiz.unit.ChatStr("Sure thing.")
					}
				}
			}
		case "Attack", "Go", "go", "attack":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.Team() == p.Team() {
				wiz.behaviour.Escorting = false
				wiz.behaviour.Guarding = false
				wiz.unit.Hunt()
				if !wiz.behaviour.Chatting {
					wiz.checkChatting()
					random2 := ns.Random(1, 4)
					if random2 == 1 {
						wiz.unit.ChatStr("I'll get them.")
					}
					if random2 == 2 {
						wiz.unit.ChatStr("Time to shine.")
					}
					if random2 == 3 {
						wiz.unit.ChatStr("On the offense.")
					}
					if random2 == 4 {
						wiz.unit.ChatStr("Time to hunt.")
					}
				}
			}
		case "guard", "stay", "Guard", "Stay":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.Team() == p.Team() {
				wiz.unit.Guard(wiz.unit.Pos(), wiz.unit.Pos(), 300)
				wiz.behaviour.Escorting = false
				wiz.behaviour.Guarding = true
				wiz.behaviour.GuardingPos = wiz.unit.Pos()
				random1 := ns.Random(1, 4)
				if !wiz.behaviour.Chatting {
					wiz.checkChatting()
					if random1 == 1 {
						wiz.unit.ChatStr("I'll guard this place.")
					}
					if random1 == 2 {
						wiz.unit.ChatStr("No problem.")
					}
					if random1 == 3 {
						wiz.unit.ChatStr("I'll stay.")
					}
					if random1 == 4 {
						wiz.unit.ChatStr("I'll hold.")
					}
				}
			}
		case "force field", "Force Field", "Force field", "force Field", "Shield", "shield":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.HasTeam(p.Unit().Team()) {
				if wiz.mana >= 80 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
									wiz.mana = wiz.mana - 80
									ns.CastSpell(spell.SHIELD, wiz.unit, p.Unit())
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
									})
								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
							})
						}
					})
				}
				if wiz.mana < 80 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
									ns.AudioEvent(audio.ManaEmpty, wiz.unit)
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
										if !wiz.behaviour.Chatting {
											wiz.checkChatting()
											wiz.unit.ChatStr("Not enough mana.")
										}
									})

								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
							})
						}
					})
				}
			}
		case "haste", "Haste":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.HasTeam(p.Unit().Team()) {
				if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
									wiz.mana = wiz.mana - 10
									ns.CastSpell(spell.HASTE, wiz.unit, p.Unit())
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
									})

								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
							})
						}
					})
				}
				if wiz.mana < 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {

									ns.AudioEvent(audio.ManaEmpty, wiz.unit)
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
									})

								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
								if !wiz.behaviour.Chatting {
									wiz.checkChatting()
									wiz.unit.ChatStr("Not enough mana.")
								}
							})
						}
					})
				}
			}
		case "Invis", "invis", "Invisibility", "invisibility":
			if wiz.unit.CanSee(p.Unit()) && wiz.unit.HasTeam(p.Unit().Team()) {
				if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {

									wiz.mana = wiz.mana - 30
									ns.CastSpell(spell.INVISIBILITY, wiz.unit, p.Unit())
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
									})

								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
							})
						}
					})
				}
				if wiz.mana < 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready {
					// Trigger cooldown.
					wiz.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
						// Check for War Cry before chant.
						if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(wiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
								// Check for War Cry before spell release.
								if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {

									ns.AudioEvent(audio.ManaEmpty, wiz.unit)
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										wiz.spells.Ready = true
									})

								}
							})
						} else {
							ns.NewTimer(ns.Frames(wiz.reactionTime), func() {
								wiz.spells.Ready = true
								if !wiz.behaviour.Chatting {
									wiz.checkChatting()
									wiz.unit.ChatStr("Not enough mana.")
								}
							})
						}
					})
				}
			}
		}
	}
	return msg
}
