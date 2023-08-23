package EndGameBW

import (
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
	cursorObject      ns.Obj
	target            ns.Obj
	trap              ns.Obj
	startingEquipment struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
		WizardRobe     ns.Obj
	}
	spells struct {
		isAlive             bool
		Ready               bool // Duration unknown.
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
	reactionTime int
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
	// Create WizBot3.
	wiz.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointBlue"))
	wiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	wiz.unit.SetMaxHealth(75)
	wiz.unit.SetStrength(35)
	wiz.unit.SetBaseSpeed(83)
	wiz.spells.isAlive = true
	// Set Team.
	wiz.unit.SetOwner(wiz.team.TeamObj)
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
	wiz.startingEquipment.StreetSneakers = ns.CreateObject("StreetSneakers", wiz.unit)
	wiz.startingEquipment.StreetPants = ns.CreateObject("StreetPants", wiz.unit)
	wiz.startingEquipment.StreetShirt = ns.CreateObject("StreetShirt", wiz.unit)
	wiz.startingEquipment.WizardRobe = ns.CreateObject("WizardRobe", wiz.unit)
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
}

func (wiz *Wizard) onEndOfWaypoint() {
	wiz.team.CheckAttackOrDefend(wiz.unit)
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
	wiz.castShock()
	wiz.castMissilesOfMagic()
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

func (wiz *Wizard) Update() {
	wiz.findLoot()
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
		wiz.castMissilesOfMagic()
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
			// BlueWizard armor.
			"WizardHelm", "WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if wiz.unit.CanSee(item) {
			wiz.unit.Equip(item)
		}
	}
}

func (wiz *Wizard) castTrap() {
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.TrapReady {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.SHOCK) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.spells.Ready && wiz.spells.ShockReady {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.INVISIBLE) && wiz.spells.Ready && wiz.spells.InvisibilityReady && wiz.unit != wiz.team.TeamTank {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.EnergyBoltReady && wiz.spells.Ready {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.DeathRayReady && wiz.spells.Ready {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.FireballReady && wiz.spells.Ready {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.spells.Ready && wiz.spells.BlinkReady && wiz.unit != wiz.team.TeamTank {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.MagicMissilesReady && wiz.spells.Ready {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && wiz.unit.CanSee(wiz.target) && wiz.spells.SlowReady && wiz.spells.Ready {
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

func (wiz *Wizard) castHaste() {
	// Check if cooldowns are ready.
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.HASTED) && wiz.spells.Ready && wiz.spells.HasteReady {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.SHIELD) && wiz.spells.Ready && wiz.spells.ForceFieldReady {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) && wiz.spells.Ready && wiz.spells.ProtFromFireReady {
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
	if wiz.spells.isAlive && !wiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !wiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) && wiz.spells.Ready && wiz.spells.ProtFromShockReady {
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
