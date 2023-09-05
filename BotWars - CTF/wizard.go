package BotWars

import (
	"image/color"
	"strconv"

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
		isAlive             bool
		Ready               bool // Duration unknown.
		DrainManaReady      bool
		DeathRayReady       bool // No cooldown, 60 mana.
		MagicMissilesReady  bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		ForceFieldReady     bool // Duration unknown.
		ShockReady          bool // Duration is 20 seconds.
		SlowReady           bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		TrapReady           bool // Only one trap is placed per life.
		EnergyBoltReady     bool // No real cooldown, mana cost *.
		FireballReady       bool // No real cooldown, mana cost 30.
		ProtFromFireReady   bool // Duration is 60 seconds.
		ProtFromPoisonReady bool
		ProtFromShockReady  bool
		BlinkReady          bool
		HasteReady          bool // Duration is 20 seconds
		InvisibilityReady   bool // Duration is 60 seconds, 30 mana.
	}
	behaviour struct {
		focussed          bool
		IsCastinDrainMana bool
		Busy              bool
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
	wiz.spells.FireballReady = true
	wiz.spells.DrainManaReady = true
	// Defensive spells.
	wiz.spells.BlinkReady = true
	// Buff spells
	wiz.spells.ShockReady = true
	wiz.spells.ProtFromFireReady = true
	wiz.spells.ProtFromPoisonReady = true
	wiz.spells.ProtFromShockReady = true
	wiz.spells.HasteReady = true
	wiz.spells.ForceFieldReady = true
	wiz.spells.InvisibilityReady = true
	// Behaviour
	wiz.behaviour.IsCastinDrainMana = false
	// Create WizBot3.
	wiz.unit = ns.CreateObject("NPC", wiz.team.SpawnPoint())
	wiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	wiz.unit.SetMaxHealth(75)
	wiz.unit.SetStrength(35)
	wiz.unit.SetBaseSpeed(83)
	wiz.spells.isAlive = true
	wiz.mana = 150
	wiz.PassiveManaRegen()
	// Set Team.
	wiz.unit.SetOwner(wiz.team.TeamObj)
	wiz.unit.SetTeam(wiz.team.TeamObj.Team())
	if wiz.team.TeamObj.HasTeam(ns.Teams()[0]) {
		wiz.unit.SetColor(0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	} else {
		wiz.unit.SetColor(0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	}
	// Create WizBot3 mouse cursor.
	wiz.target = wiz.team.Enemy.TeamObj
	wiz.cursor = wiz.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	wiz.reactionTime = 15
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
	//wiz.onChat()
	wiz.LookForWeapon()
}

//func (wiz *Wizard) onChat() {
//	ns.OnChat(func(t ns.Team, pl ns.Player, obj ns.Obj, msg string) string {
//		fmt.Printf("wizard %p heard message: team=%v, player=%v, obj=%v, msg=%s\n", wiz, t, pl, obj, msg)
//		return msg // pass message to the chat
//	})
//}

func (wiz *Wizard) onHit() {
	if wiz.mana <= 49 {
		wiz.GoToManaObelisk()
	}
}

func (wiz *Wizard) UsePotions() {
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

func (wiz *Wizard) onEndOfWaypoint() {
	if wiz.mana <= 49 {
		wiz.GoToManaObelisk()
	} else {
		wiz.team.CheckAttackOrDefend(wiz.unit)
	}
	wiz.LookForNearbyItems()
}

func (wiz *Wizard) buffInitial() {
	wiz.castForceField()
}

func (wiz *Wizard) onLookingForTarget() {
}

func (wiz *Wizard) onEnemyHeard() {
	wiz.castFireballAtHeard()
	wiz.castInvisibility()
}

func (wiz *Wizard) onEnemySighted() {
	wiz.target = ns.GetCaller()
	wiz.castSlow()
}

func (wiz *Wizard) onCollide() {
	//wiz.castShock()
	//wiz.castMissilesOfMagic()
	if wiz.spells.isAlive {
		caller := ns.GetCaller()
		wiz.team.CheckPickUpEnemyFlag(caller, wiz.unit)
		wiz.team.CheckCaptureEnemyFlag(caller, wiz.unit)
		wiz.team.CheckRetrievedOwnFlag(caller, wiz.unit)
	}
}

func (wiz *Wizard) onRetreat() {
	wiz.castBlink()
}

func (wiz *Wizard) onLostEnemy() {
	wiz.castTrap()
	wiz.team.WalkToOwnFlag(wiz.unit)
}

func (wiz *Wizard) onDeath() {
	wiz.spells.isAlive = false
	wiz.spells.Ready = false
	wiz.team.DropEnemyFlag(wiz.unit)
	wiz.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, wiz.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, wiz.unit)
		wiz.unit.Delete()
		wiz.startingEquipment.StreetPants.Delete()
		wiz.startingEquipment.StreetSneakers.Delete()
		wiz.startingEquipment.StreetShirt.Delete()
		wiz.init()
	})
}

func (wiz *Wizard) PassiveManaRegen() {
	if wiz.spells.isAlive {
		ns.NewTimer(ns.Seconds(2), func() {
			if wiz.mana < 150 {
				wiz.mana = wiz.mana + 1
			}
			wiz.PassiveManaRegen()
			ns.PrintStrToAll("wiz mana: " + strconv.Itoa(wiz.mana))
		})
	}
}

func (wiz *Wizard) GoToManaObelisk() {
	//wiz.unit.AggressionLevel(0.16)
	ManaSource := ns.FindClosestObject(wiz.unit, ns.HasTypeName{
		"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12",
	})
	wiz.unit.WalkTo(ManaSource.Pos())
}

func (wiz *Wizard) RestoreMana() {
	if wiz.mana < 150 {
		ManaSource := ns.FindAllObjects(
			ns.HasTypeName{
				"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12",
			},
			ns.InCirclef{Center: wiz.unit, R: 50},
		)
		for i := 0; i < len(ManaSource); i++ {
			if ManaSource[i].CurrentMana() > 0 && wiz.unit.CanSee(ManaSource[i]) {
				wiz.mana = wiz.mana + 1
				ManaSource[i].SetMana(ManaSource[i].CurrentMana() - 1)
				wiz.RestoreManaSound()
			}
		}
	}
}

func (wiz *Wizard) RestoreManaWithDrainMana() {
	if wiz.mana < 150 && wiz.behaviour.IsCastinDrainMana {
		ManaSource := ns.FindAllObjects(
			ns.HasTypeName{
				"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12",
			},
			ns.InCirclef{Center: wiz.unit, R: 200},
		)
		for i := 0; i < len(ManaSource); i++ {
			if ManaSource[i].CurrentMana() > 0 && wiz.unit.CanSee(ManaSource[i]) {
				wiz.mana = wiz.mana + 1
				ManaSource[i].SetMana(ManaSource[i].CurrentMana() - 1)
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

func (wiz *Wizard) Update() {
	wiz.findLoot()
	wiz.UsePotions()
	wiz.RestoreMana()
	wiz.RestoreManaWithDrainMana()
	if wiz.mana > 150 {
		wiz.mana = 150
	}
	if wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
		wiz.spells.Ready = true
	}
	if wiz.target.HasEnchant(enchant.HELD) || wiz.target.HasEnchant(enchant.SLOWED) {
		if wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
			wiz.castDeathRay()
		}
	}
	if wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
		wiz.castFireball()
		wiz.castSlow()
		wiz.castEnergyBolt()
		//wiz.castMissilesOfMagic()
		wiz.castForceField()
		wiz.castShock()
	}
	if !wiz.unit.CanSee(wiz.target) && wiz.spells.Ready {
		wiz.castHaste()
		wiz.castProtectionFromShock()
		wiz.castProtectionFromFire()
		wiz.castInvisibility()
	}
}

func (wiz *Wizard) LookForWeapon() {
	ItemLocation := ns.FindClosestObject(wiz.unit, ns.HasTypeName{"FireStormWand", "LesserFireballWand", "ForceWand"})
	wiz.unit.WalkTo(ItemLocation.Pos())
}

func (wiz *Wizard) LookForNearbyItems() {
	if ns.FindAllObjects(ns.HasTypeName{
		"RedPotion",
		"FireStormWand", "LesserFireballWand", "ForceWand",
		"CurePoisonPotion",
		"WizardHelm",
		"WizardRobe",
		"BluePotion",
		"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
		ns.InCirclef{Center: wiz.unit, R: 200}) != nil {
		ItemLocation := ns.FindAllObjects(ns.HasTypeName{
			"RedPotion",
			"CurePoisonPotion",
			"WizardHelm",
			"WizardRobe",
			"BluePotion",
			"FireStormWand", "LesserFireballWand", "ForceWand",
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
			ns.InCirclef{Center: wiz.unit, R: 200},
		)
		if wiz.unit.CanSee(ItemLocation[0]) {
			wiz.unit.WalkTo(ItemLocation[0].Pos())
		}
	}
}

func (wiz *Wizard) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			//"DeathRayWand",
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
			wiz.unit.Equip(item)
		}
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			// Armor.
			"WizardHelm",
			"WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if wiz.unit.CanSee(item) {
			wiz.unit.Equip(item)
		}
	}

	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: wiz.unit, R: dist},
		ns.HasTypeName{
			"RedPotion",
			"CurePoisonPotion",
			"BluePotion",
		},
	)
	for _, item := range potions {
		if wiz.unit.CanSee(item) {
			wiz.unit.Pickup(item)
		}
	}

}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //

func (wiz *Wizard) castTrap() {
	if wiz.mana >= 105 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.TrapReady {
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
															wiz.spells.Ready = true
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
						wiz.spells.Ready = true
						// Shock cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							wiz.spells.ShockReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castInvisibility() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.spells.Ready && wiz.spells.InvisibilityReady && wiz.unit != wiz.team.TeamTank {
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
						wiz.spells.Ready = true
						// Invisibility cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.InvisibilityReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castEnergyBolt() {
	// Check if cooldowns are ready.
	if wiz.mana > 20 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.EnergyBoltReady && wiz.spells.Ready {
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
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.EnergyBoltReady = false
						ns.CastSpell(spell.LIGHTNING, wiz.unit, wiz.target)
						// Global cooldown.
						wiz.spells.Ready = true
						// Energy Bolt cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.EnergyBoltReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castDeathRay() {
	// Check if cooldowns are ready.
	if wiz.mana >= 60 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.DeathRayReady && wiz.spells.Ready {
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
					if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						wiz.unit.LookAtObject(wiz.target)
						wiz.unit.Pause(ns.Frames(wiz.reactionTime))
						wiz.spells.DeathRayReady = false
						ns.CastSpell(spell.DEATH_RAY, wiz.unit, wiz.cursor)
						wiz.mana = wiz.mana - 60
						// Global cooldown.
						wiz.spells.Ready = true
						// Death Ray cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							wiz.spells.DeathRayReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castFireball() {
	// Check if cooldowns are ready.
	if wiz.mana >= 30 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready {
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
						wiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							wiz.spells.FireballReady = true
						})
					}
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
						ns.CastSpell(spell.FIREBALL, wiz.unit, wiz.target)
						// Global cooldown.
						wiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							wiz.spells.FireballReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castBlink() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.BlinkReady && wiz.unit != wiz.team.TeamTank {
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
						wiz.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							wiz.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castMissilesOfMagic() {
	// Check if cooldowns are ready.
	if wiz.mana >= 15 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.MagicMissilesReady && wiz.spells.Ready {
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
						wiz.spells.Ready = true
						// Missiles Of Magic cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.MagicMissilesReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castSlow() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.SlowReady && wiz.spells.Ready {
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
						wiz.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (wiz *Wizard) castDrainMana() {
	// Check if cooldowns are ready.
	if wiz.mana >= 10 && wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.DrainManaReady {
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
						wiz.behaviour.IsCastinDrainMana = true
						wiz.unit.Pause(ns.Frames(30))
						ManaSource := ns.FindAllObjects(
							ns.HasTypeName{
								"ObeliskPrimitive", "Obelisk", "InvisibleObelisk", "InvisibleObeliskNWSE", "MineCrystal01", "MineCrystal02", "MineCrystal03", "MineCrystal04", "MineCrystal05", "MineCrystalDown01", "MineCrystalDown02", "MineCrystalDown03", "MineCrystalDown04", "MineCrystalDown05", "MineCrystalUp01", "MineCrystalUp02", "MineCrystalUp03", "MineCrystalUp04", "MineCrystalUp05", "MineManaCart1", "MineManaCart1", "MineManaCrystal1", "MineManaCrystal2", "MineManaCrystal3", "MineManaCrystal4", "MineManaCrystal5", "MineManaCrystal6", "MineManaCrystal7", "MineManaCrystal8", "MineManaCrystal9", "MineManaCrystal10", "MineManaCrystal11", "MineManaCrystal12",
							},
							ns.InCirclef{Center: wiz.unit, R: 200})
						wiz.unit.LookAtObject(ManaSource[0])
						ns.NewTimer(ns.Frames(30), func() {
							wiz.behaviour.IsCastinDrainMana = false
						})
						ns.CastSpell(spell.DRAIN_MANA, wiz.unit, wiz.unit.Pos())
						// Global cooldown.
						wiz.spells.Ready = true
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							wiz.spells.DrainManaReady = true
						})
					}
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
						wiz.spells.Ready = true
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							wiz.spells.HasteReady = true
						})
					}
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
						wiz.spells.Ready = true
						// Force Field cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.ForceFieldReady = true
						})
					}
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
						wiz.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.ProtFromFireReady = true
						})
					}
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
						wiz.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							wiz.spells.ProtFromShockReady = true
						})
					}
				})
			}
		})
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //
