package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWizard creates a new RedWizard bot.
func NewRedWizard() *RedWizard {
	redwiz := &RedWizard{}
	redwiz.init()
	return redwiz
}

// RedWizard bot class.
type RedWizard struct {
	unit         ns.Obj
	cursor       ns.Pointf
	cursorObject ns.Obj
	target       ns.Obj
	trap         ns.Obj
	items        struct {
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

func (redwiz *RedWizard) init() {
	// Reset spells WizBot3.
	redwiz.spells.Ready = true
	// Debuff spells.
	redwiz.spells.SlowReady = true
	// Offensive spells.
	redwiz.spells.MagicMissilesReady = true
	redwiz.spells.TrapReady = true
	redwiz.spells.DeathRayReady = true
	redwiz.spells.EnergyBoltReady = true
	redwiz.spells.FireballReady = true
	// Defensive spells.
	redwiz.spells.BlinkReady = true
	// Buff spells
	redwiz.spells.ShockReady = true
	redwiz.spells.ProtFromFireReady = true
	redwiz.spells.ProtFromPoisonReady = true
	redwiz.spells.ProtFromShockReady = true
	redwiz.spells.HasteReady = true
	redwiz.spells.ForceFieldReady = true
	redwiz.spells.InvisibilityReady = true
	// Create WizBot3.
	redwiz.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointRed"))
	redwiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	redwiz.unit.SetMaxHealth(75)
	redwiz.unit.SetStrength(35)
	redwiz.unit.SetBaseSpeed(83)
	redwiz.spells.isAlive = true
	// Set Team.
	redwiz.unit.SetOwner(TeamRed)
	// Create WizBot3 mouse cursor.
	redwiz.target = TeamBlue
	redwiz.cursor = redwiz.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	redwiz.reactionTime = 15
	// Set WizBot3 properties.
	redwiz.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	redwiz.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	redwiz.unit.MonsterStatusEnable(object.MonStatusAlert)
	redwiz.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		redwiz.unit.AggressionLevel(0.83)
	})
	redwiz.unit.Hunt()
	redwiz.unit.ResumeLevel(0.8)
	redwiz.unit.RetreatLevel(0.2)
	// Create and equip WizBot3 starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	redwiz.items.StreetSneakers = ns.CreateObject("StreetSneakers", redwiz.unit)
	redwiz.items.StreetPants = ns.CreateObject("StreetPants", redwiz.unit)
	redwiz.items.StreetShirt = ns.CreateObject("StreetShirt", redwiz.unit)
	redwiz.items.WizardRobe = ns.CreateObject("WizardRobe", redwiz.unit)
	redwiz.unit.Equip(redwiz.items.StreetSneakers)
	redwiz.unit.Equip(redwiz.items.StreetPants)
	redwiz.unit.Equip(redwiz.items.StreetShirt)
	redwiz.unit.Equip(redwiz.items.WizardRobe)
	// Buff on respawn.
	redwiz.buffInitial()
	// On retreat.
	redwiz.unit.OnEvent(ns.EventRetreat, redwiz.onRetreat)
	// Enemy sighted.
	redwiz.unit.OnEvent(ns.EventEnemySighted, redwiz.onEnemySighted)
	// On heard.
	redwiz.unit.OnEvent(ns.EventEnemyHeard, redwiz.onEnemyHeard)
	// On collision.
	redwiz.unit.OnEvent(ns.EventCollision, redwiz.onCollide)
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight.
	redwiz.unit.OnEvent(ns.EventLostEnemy, redwiz.onLostEnemy)
	// On Death.
	redwiz.unit.OnEvent(ns.EventDeath, redwiz.onDeath)
	redwiz.unit.OnEvent(ns.EventLookingForEnemy, redwiz.onEndOfWaypoint)
}

func (redwiz *RedWizard) onEndOfWaypoint() {
	redwiz.RedTeamCheckAttackOrDefend()
}

func (redwiz *RedWizard) buffInitial() {
	redwiz.castForceField()
}

func (redwiz *RedWizard) onLookingForTarget() {
}

func (redwiz *RedWizard) onEnemyHeard() {
	redwiz.castFireballAtHeard()
	redwiz.castInvisibility()
}

func (redwiz *RedWizard) onEnemySighted() {
	redwiz.target = ns.GetCaller()
	redwiz.castSlow()
}

func (redwiz *RedWizard) onCollide() {
	redwiz.castShock()
	redwiz.castMissilesOfMagic()
	if redwiz.spells.isAlive {
		redwiz.RedTeamPickUpBlueFlag()
		redwiz.RedTeamCaptureTheBlueFlag()
		redwiz.RedTeamRetrievedRedFlag()
	}
}

func (redwiz *RedWizard) onRetreat() {
	redwiz.castBlink()
}

func (redwiz *RedWizard) onLostEnemy() {
	redwiz.castTrap()
	redwiz.RedTeamWalkToRedFlag()
}

func (redwiz *RedWizard) onDeath() {
	redwiz.spells.isAlive = false
	redwiz.spells.Ready = false
	redwiz.RedTeamDropFlag()
	redwiz.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, redwiz.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, redwiz.unit)
		redwiz.unit.Delete()
		redwiz.items.StreetPants.Delete()
		redwiz.items.StreetSneakers.Delete()
		redwiz.items.StreetShirt.Delete()
		redwiz.init()
	})
}

func (redwiz *RedWizard) Update() {
	redwiz.findLoot()
	if redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
		redwiz.spells.Ready = true
	}
	if redwiz.target.HasEnchant(enchant.HELD) || redwiz.target.HasEnchant(enchant.SLOWED) {
		if redwiz.unit.CanSee(redwiz.target) && redwiz.spells.Ready {
			redwiz.castDeathRay()
		}
	}
	if redwiz.unit.CanSee(redwiz.target) && redwiz.spells.Ready {
		redwiz.castFireball()
		redwiz.castSlow()
		redwiz.castEnergyBolt()
		redwiz.castMissilesOfMagic()
		redwiz.castForceField()
		redwiz.castShock()
	}
	if !redwiz.unit.CanSee(redwiz.target) && redwiz.spells.Ready {
		redwiz.castHaste()
		redwiz.castProtectionFromShock()
		redwiz.castProtectionFromFire()
		redwiz.castInvisibility()
	}
}

func (redwiz *RedWizard) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: redwiz.unit, R: dist},
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
		if redwiz.unit.CanSee(item) {
			redwiz.unit.Equip(item)
		}
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: redwiz.unit, R: dist},
		ns.HasTypeName{
			// RedWizard armor.
			"WizardHelm", "WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if redwiz.unit.CanSee(item) {
			redwiz.unit.Equip(item)
		}
	}
}

func (redwiz *RedWizard) castTrap() {
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.spells.Ready && redwiz.spells.TrapReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Ring of Fire chant.
				castPhonemes(redwiz.unit, []audio.Name{PhDownRight, PhDown, PhDownLeft, PhUp}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Magic Missiles chant.
							castPhonemes(redwiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Shock chant.
										castPhonemes(redwiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(redwiz.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
															redwiz.spells.TrapReady = false
															ns.AudioEvent(audio.TrapDrop, redwiz.unit)
															redwiz.trap = ns.NewTrap(redwiz.unit, spell.CLEANSING_FLAME, spell.MAGIC_MISSILE, spell.SHOCK)
															redwiz.trap.SetOwner(redwiz.unit)
															// Global cooldown.
															redwiz.spells.Ready = true
															// Trap cooldown.
															ns.NewTimer(ns.Seconds(5), func() {
																redwiz.spells.TrapReady = true
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

func (redwiz *RedWizard) castShock() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.SHOCK) && !redwiz.unit.HasEnchant(enchant.INVISIBLE) && redwiz.spells.Ready && redwiz.spells.ShockReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.ShockReady = false
						ns.CastSpell(spell.SHOCK, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Shock cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							redwiz.spells.ShockReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castInvisibility() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.INVISIBLE) && redwiz.spells.Ready && redwiz.spells.InvisibilityReady && redwiz.unit != RedTeamTank {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.InvisibilityReady = false
						ns.CastSpell(spell.INVISIBILITY, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Invisibility cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redwiz.spells.InvisibilityReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castEnergyBolt() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.unit.CanSee(redwiz.target) && redwiz.spells.EnergyBoltReady && redwiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.EnergyBoltReady = false
						ns.CastSpell(spell.LIGHTNING, redwiz.unit, redwiz.target)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Energy Bolt cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							redwiz.spells.EnergyBoltReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castDeathRay() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.spells.DeathRayReady && redwiz.spells.Ready {
		// Select target.
		redwiz.cursor = redwiz.target.Pos()
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDownRight, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.DeathRayReady = false
						ns.CastSpell(spell.DEATH_RAY, redwiz.unit, redwiz.cursor)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Death Ray cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							redwiz.spells.DeathRayReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castFireball() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.unit.CanSee(redwiz.target) && redwiz.spells.FireballReady && redwiz.spells.Ready {
		// Select target.
		redwiz.cursor = redwiz.target.Pos()
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.FireballReady = false
						ns.CastSpell(spell.FIREBALL, redwiz.unit, redwiz.cursor)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							redwiz.spells.FireballReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castFireballAtHeard() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.unit.CanSee(redwiz.target) && redwiz.spells.FireballReady && redwiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.FireballReady = false
						ns.CastSpell(spell.FIREBALL, redwiz.unit, redwiz.target)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							redwiz.spells.FireballReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castBlink() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.spells.Ready && redwiz.spells.BlinkReady && redwiz.unit != RedTeamTank {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.BlinkReady = false
						ns.NewTrap(redwiz.unit, spell.BLINK)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							redwiz.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castMissilesOfMagic() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.unit.CanSee(redwiz.target) && redwiz.spells.MagicMissilesReady && redwiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.MagicMissilesReady = false
						ns.CastSpell(spell.MAGIC_MISSILE, redwiz.unit, redwiz.target)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Missiles Of Magic cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							redwiz.spells.MagicMissilesReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castSlow() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && redwiz.unit.CanSee(redwiz.target) && redwiz.spells.SlowReady && redwiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redwiz.unit.LookAtObject(redwiz.target)
						redwiz.unit.Pause(ns.Frames(redwiz.reactionTime))
						redwiz.spells.SlowReady = false
						ns.CastSpell(spell.SLOW, redwiz.unit, redwiz.target)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							redwiz.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castHaste() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.HASTED) && redwiz.spells.Ready && redwiz.spells.HasteReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.HasteReady = false
						ns.CastSpell(spell.HASTE, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							redwiz.spells.HasteReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castForceField() {
	// if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.SHIELD)
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.SHIELD) && redwiz.spells.Ready && redwiz.spells.ForceFieldReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.ForceFieldReady = false
						ns.CastSpell(spell.SHIELD, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Force Field cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redwiz.spells.ForceFieldReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) && redwiz.spells.Ready && redwiz.spells.ProtFromFireReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.ProtFromFireReady = false
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redwiz.spells.ProtFromFireReady = true
						})
					}
				})
			}
		})
	}
}

func (redwiz *RedWizard) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !redwiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) && redwiz.spells.Ready && redwiz.spells.ProtFromShockReady {
		// Trigger cooldown.
		redwiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redwiz.reactionTime), func() {
			// Check for War Cry before chant.
			if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redwiz.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if redwiz.spells.isAlive && !redwiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redwiz.spells.ProtFromShockReady = false
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, redwiz.unit, redwiz.unit)
						// Global cooldown.
						redwiz.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redwiz.spells.ProtFromShockReady = true
						})
					}
				})
			}
		})
	}
}

// ---------------------------------- CTF BOT SCRIPT ------------------------------------//
// CTF game mechanics.
// Pick up the enemy flag.
func (redwiz *RedWizard) RedTeamPickUpBlueFlag() {
	if ns.GetCaller() == BlueFlag {
		BlueFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		RedTeamTank = redwiz.unit
		RedTeamTank.AggressionLevel(0.16)
		RedTeamTank.WalkTo(RedBase.Pos())
		ns.PrintStrToAll("Team Red has the Blue flag!")
	}
}

// Capture the flag.
func (redwiz *RedWizard) RedTeamCaptureTheBlueFlag() {
	if ns.GetCaller() == RedFlag && RedFlagIsAtBase && redwiz.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagCapture, RedTeamTank) // <----- replace with all players

		RedTeamTank = TeamRed
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[0].ChangeScore(+1)
		}
		FlagReset()
		redwiz.unit.AggressionLevel(0.83)
		redwiz.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has captured the Blue flag!")
	}
}

// Retrieve own flag.
func (redwiz *RedWizard) RedTeamRetrievedRedFlag() {
	if ns.GetCaller() == RedFlag && !RedFlagIsAtBase {
		RedFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
		redwiz.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has retrieved the flag!")
		RedTeamTank.WalkTo(RedFlag.Pos())
	}
}

// Drop flag.
func (redwiz *RedWizard) RedTeamDropFlag() {
	if redwiz.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		BlueFlag.Enable(true)
		RedTeamTank = TeamRed
		ns.PrintStrToAll("Team Red has dropped the Blue flag!")
	}
}

// CTF behaviour.

// When enemy is killed check to see if the flag is dropped, if so get it.
func (redwiz *RedWizard) RedTeamWalkToRedFlag() {
	if !RedFlagIsAtBase && RedFlag.IsEnabled() {
		redwiz.unit.AggressionLevel(0.16)
		redwiz.unit.WalkTo(BlueFlag.Pos())
	} else {
		redwiz.RedTeamCheckAttackOrDefend()
	}
}

// At the end of waypoint see defend if tank, attack if not.
func (redwiz *RedWizard) RedTeamCheckAttackOrDefend() {
	if redwiz.unit == RedTeamTank {
		redwiz.unit.AggressionLevel(0.16)
		redwiz.unit.Guard(RedBase.Pos(), RedBase.Pos(), 20)
	} else {
		redwiz.unit.AggressionLevel(0.83)
		redwiz.unit.WalkTo(BlueFlag.Pos())
	}
}
