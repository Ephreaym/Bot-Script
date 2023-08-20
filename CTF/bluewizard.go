package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewWizard creates a new BlueWizard bot.
func NewBlueWizard() *BlueWizard {
	bluewiz := &BlueWizard{}
	bluewiz.init()
	return bluewiz
}

// BlueWizard bot class.
type BlueWizard struct {
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

func (bluewiz *BlueWizard) init() {
	// Reset spells WizBot3.
	bluewiz.spells.Ready = true
	// Debuff spells.
	bluewiz.spells.SlowReady = true
	// Offensive spells.
	bluewiz.spells.MagicMissilesReady = true
	bluewiz.spells.TrapReady = true
	bluewiz.spells.DeathRayReady = true
	bluewiz.spells.EnergyBoltReady = true
	bluewiz.spells.FireballReady = true
	// Defensive spells.
	bluewiz.spells.BlinkReady = true
	// Buff spells
	bluewiz.spells.ShockReady = true
	bluewiz.spells.ProtFromFireReady = true
	bluewiz.spells.ProtFromPoisonReady = true
	bluewiz.spells.ProtFromShockReady = true
	bluewiz.spells.HasteReady = true
	bluewiz.spells.ForceFieldReady = true
	bluewiz.spells.InvisibilityReady = true
	// Create WizBot3.
	bluewiz.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointBlue"))
	bluewiz.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	bluewiz.unit.SetMaxHealth(75)
	bluewiz.unit.SetStrength(35)
	bluewiz.unit.SetBaseSpeed(83)
	bluewiz.spells.isAlive = true
	// Set Team.
	bluewiz.unit.SetOwner(TeamBlue)
	// Create WizBot3 mouse cursor.
	bluewiz.target = TeamRed
	bluewiz.cursor = bluewiz.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	bluewiz.reactionTime = 15
	// Set WizBot3 properties.
	bluewiz.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	bluewiz.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	bluewiz.unit.MonsterStatusEnable(object.MonStatusAlert)
	bluewiz.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		bluewiz.unit.AggressionLevel(0.83)
	})
	bluewiz.unit.Hunt()
	bluewiz.unit.ResumeLevel(0.8)
	bluewiz.unit.RetreatLevel(0.2)
	// Create and equip WizBot3 starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	bluewiz.items.StreetSneakers = ns.CreateObject("StreetSneakers", bluewiz.unit)
	bluewiz.items.StreetPants = ns.CreateObject("StreetPants", bluewiz.unit)
	bluewiz.items.StreetShirt = ns.CreateObject("StreetShirt", bluewiz.unit)
	bluewiz.items.WizardRobe = ns.CreateObject("WizardRobe", bluewiz.unit)
	bluewiz.unit.Equip(bluewiz.items.StreetSneakers)
	bluewiz.unit.Equip(bluewiz.items.StreetPants)
	bluewiz.unit.Equip(bluewiz.items.StreetShirt)
	bluewiz.unit.Equip(bluewiz.items.WizardRobe)
	// Buff on respawn.
	bluewiz.buffInitial()
	// On retreat.
	bluewiz.unit.OnEvent(ns.EventRetreat, bluewiz.onRetreat)
	// Enemy sighted.
	bluewiz.unit.OnEvent(ns.EventEnemySighted, bluewiz.onEnemySighted)
	// On heard.
	bluewiz.unit.OnEvent(ns.EventEnemyHeard, bluewiz.onEnemyHeard)
	// On collision.
	bluewiz.unit.OnEvent(ns.EventCollision, bluewiz.onCollide)
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight.
	bluewiz.unit.OnEvent(ns.EventLostEnemy, bluewiz.onLostEnemy)
	// On Death.
	bluewiz.unit.OnEvent(ns.EventDeath, bluewiz.onDeath)
	bluewiz.unit.OnEvent(ns.EventLookingForEnemy, bluewiz.onLookingForTarget)
	bluewiz.unit.OnEvent(ns.EventEndOfWaypoint, bluewiz.onEndOfWaypoint)
}

func (bluewiz *BlueWizard) onEndOfWaypoint() {
	bluewiz.BlueTeamCheckAttackOrDefend()
}

func (bluewiz *BlueWizard) buffInitial() {
	bluewiz.castForceField()
}

func (bluewiz *BlueWizard) onLookingForTarget() {
}

func (bluewiz *BlueWizard) onEnemyHeard() {
	bluewiz.castFireballAtHeard()
	bluewiz.castInvisibility()
}

func (bluewiz *BlueWizard) onEnemySighted() {
	bluewiz.target = ns.GetCaller()
	bluewiz.castSlow()
}

func (bluewiz *BlueWizard) onCollide() {
	bluewiz.castShock()
	bluewiz.castMissilesOfMagic()
	if bluewiz.spells.isAlive {
		bluewiz.BlueTeamPickUpRedFlag()
		bluewiz.BlueTeamCaptureTheRedFlag()
		bluewiz.BlueTeamRetrievedBlueFlag()
	}
}

func (bluewiz *BlueWizard) onRetreat() {
	bluewiz.castBlink()
}

func (bluewiz *BlueWizard) onLostEnemy() {
	bluewiz.castTrap()
	bluewiz.BlueTeamWalkToBlueFlag()
}

func (bluewiz *BlueWizard) onDeath() {
	bluewiz.spells.isAlive = false
	bluewiz.spells.Ready = false
	bluewiz.BlueTeamDropFlag()
	bluewiz.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, bluewiz.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, bluewiz.unit)
		bluewiz.unit.Delete()
		bluewiz.items.StreetPants.Delete()
		bluewiz.items.StreetSneakers.Delete()
		bluewiz.items.StreetShirt.Delete()
		bluewiz.init()
	})
}

func (bluewiz *BlueWizard) Update() {
	bluewiz.findLoot()
	if bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
		bluewiz.spells.Ready = true
	}
	if bluewiz.target.HasEnchant(enchant.HELD) || bluewiz.target.HasEnchant(enchant.SLOWED) {
		if bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.Ready {
			bluewiz.castDeathRay()
		}
	}
	if bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.Ready {
		bluewiz.castFireball()
		bluewiz.castSlow()
		bluewiz.castEnergyBolt()
		bluewiz.castMissilesOfMagic()
		bluewiz.castForceField()
		bluewiz.castShock()
	}
	if !bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.Ready {
		bluewiz.castHaste()
		bluewiz.castProtectionFromShock()
		bluewiz.castProtectionFromFire()
		bluewiz.castInvisibility()
	}
}

func (bluewiz *BlueWizard) findLoot() {
	const dist = 75
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: bluewiz.unit, R: dist},
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
		if bluewiz.unit.CanSee(item) {
			bluewiz.unit.Equip(item)
		}
	}

	armor := ns.FindAllObjects(
		ns.InCirclef{Center: bluewiz.unit, R: dist},
		ns.HasTypeName{
			// BlueWizard armor.
			"WizardHelm", "WizardRobe",
			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if bluewiz.unit.CanSee(item) {
			bluewiz.unit.Equip(item)
		}
	}
}

func (bluewiz *BlueWizard) castTrap() {
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.spells.Ready && bluewiz.spells.TrapReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Ring of Fire chant.
				castPhonemes(bluewiz.unit, []audio.Name{PhDownRight, PhDown, PhDownLeft, PhUp}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Magic Missiles chant.
							castPhonemes(bluewiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Shock chant.
										castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(bluewiz.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
															bluewiz.spells.TrapReady = false
															ns.AudioEvent(audio.TrapDrop, bluewiz.unit)
															bluewiz.trap = ns.NewTrap(bluewiz.unit, spell.CLEANSING_FLAME, spell.MAGIC_MISSILE, spell.SHOCK)
															bluewiz.trap.SetOwner(bluewiz.unit)
															// Global cooldown.
															bluewiz.spells.Ready = true
															// Trap cooldown.
															ns.NewTimer(ns.Seconds(5), func() {
																bluewiz.spells.TrapReady = true
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

func (bluewiz *BlueWizard) castShock() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.SHOCK) && !bluewiz.unit.HasEnchant(enchant.INVISIBLE) && bluewiz.spells.Ready && bluewiz.spells.ShockReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhLeft}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.ShockReady = false
						ns.CastSpell(spell.SHOCK, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Shock cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							bluewiz.spells.ShockReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castInvisibility() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.INVISIBLE) && bluewiz.spells.Ready && bluewiz.spells.InvisibilityReady && bluewiz.unit != BlueTeamTank {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhLeft, PhRight, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.InvisibilityReady = false
						ns.CastSpell(spell.INVISIBILITY, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Invisibility cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluewiz.spells.InvisibilityReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castEnergyBolt() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.EnergyBoltReady && bluewiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.EnergyBoltReady = false
						ns.CastSpell(spell.LIGHTNING, bluewiz.unit, bluewiz.target)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Energy Bolt cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							bluewiz.spells.EnergyBoltReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castDeathRay() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.spells.DeathRayReady && bluewiz.spells.Ready {
		// Select target.
		bluewiz.cursor = bluewiz.target.Pos()
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDownRight, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.DeathRayReady = false
						ns.CastSpell(spell.DEATH_RAY, bluewiz.unit, bluewiz.cursor)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Death Ray cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							bluewiz.spells.DeathRayReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castFireball() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.FireballReady && bluewiz.spells.Ready {
		// Select target.
		bluewiz.cursor = bluewiz.target.Pos()
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.FireballReady = false
						ns.CastSpell(spell.FIREBALL, bluewiz.unit, bluewiz.cursor)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							bluewiz.spells.FireballReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castFireballAtHeard() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.FireballReady && bluewiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.FireballReady = false
						ns.CastSpell(spell.FIREBALL, bluewiz.unit, bluewiz.target)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Fireball cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							bluewiz.spells.FireballReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castBlink() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.spells.Ready && bluewiz.spells.BlinkReady && bluewiz.unit != BlueTeamTank {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.BlinkReady = false
						ns.NewTrap(bluewiz.unit, spell.BLINK)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							bluewiz.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castMissilesOfMagic() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.MagicMissilesReady && bluewiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhLeft, PhUp, PhRight, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.MagicMissilesReady = false
						ns.CastSpell(spell.MAGIC_MISSILE, bluewiz.unit, bluewiz.target)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Missiles Of Magic cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							bluewiz.spells.MagicMissilesReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castSlow() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && bluewiz.unit.CanSee(bluewiz.target) && bluewiz.spells.SlowReady && bluewiz.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluewiz.unit.LookAtObject(bluewiz.target)
						bluewiz.unit.Pause(ns.Frames(bluewiz.reactionTime))
						bluewiz.spells.SlowReady = false
						ns.CastSpell(spell.SLOW, bluewiz.unit, bluewiz.target)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							bluewiz.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castHaste() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.HASTED) && bluewiz.spells.Ready && bluewiz.spells.HasteReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhLeft, PhRight, PhRight}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.HasteReady = false
						ns.CastSpell(spell.HASTE, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							bluewiz.spells.HasteReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castForceField() {
	// if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.SHIELD)
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.SHIELD) && bluewiz.spells.Ready && bluewiz.spells.ForceFieldReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhUp, PhLeft, PhDown, PhRight, PhUp, PhLeft, PhDown, PhRight}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.ForceFieldReady = false
						ns.CastSpell(spell.SHIELD, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Force Field cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluewiz.spells.ForceFieldReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) && bluewiz.spells.Ready && bluewiz.spells.ProtFromFireReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.ProtFromFireReady = false
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluewiz.spells.ProtFromFireReady = true
						})
					}
				})
			}
		})
	}
}

func (bluewiz *BlueWizard) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) && !bluewiz.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) && bluewiz.spells.Ready && bluewiz.spells.ProtFromShockReady {
		// Trigger cooldown.
		bluewiz.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluewiz.reactionTime), func() {
			// Check for War Cry before chant.
			if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluewiz.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if bluewiz.spells.isAlive && !bluewiz.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluewiz.spells.ProtFromShockReady = false
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, bluewiz.unit, bluewiz.unit)
						// Global cooldown.
						bluewiz.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluewiz.spells.ProtFromShockReady = true
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
func (bluewiz *BlueWizard) BlueTeamPickUpRedFlag() {
	if ns.GetCaller() == RedFlag {
		RedFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		BlueTeamTank = bluewiz.unit
		BlueTeamTank.AggressionLevel(0.16)
		BlueTeamTank.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Blue has the Red flag!")
	}
}

// Capture the flag.
func (bluewiz *BlueWizard) BlueTeamCaptureTheRedFlag() {
	if ns.GetCaller() == BlueFlag && BlueFlagIsAtBase && bluewiz.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagCapture, BlueTeamTank) // <----- replace with all players
		BlueTeamTank = TeamBlue
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[1].ChangeScore(+1)
		}
		FlagReset()
		bluewiz.unit.AggressionLevel(0.83)
		bluewiz.unit.WalkTo(RedFlag.Pos())
		ns.PrintStrToAll("Team Blue has captured the Red flag!")
	}
}

// Retrieve own flag.
func (bluewiz *BlueWizard) BlueTeamRetrievedBlueFlag() {
	if ns.GetCaller() == BlueFlag && !BlueFlagIsAtBase {
		BlueFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
		bluewiz.unit.WalkTo(BlueBase.Pos())
		ns.PrintStrToAll("Team Blue has retrieved the flag!")
		BlueTeamTank.WalkTo(BlueFlag.Pos())
	}
}

// Drop flag.
func (bluewiz *BlueWizard) BlueTeamDropFlag() {
	if bluewiz.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		RedFlag.Enable(true)
		BlueTeamTank = TeamBlue
		ns.PrintStrToAll("Team Blue has dropped the Red flag!")
	}
}

// CTF behaviour.
// Attack enemy tank without

func (bluewiz *BlueWizard) BlueTeamWalkToBlueFlag() {
	if !BlueFlagIsAtBase && BlueFlag.IsEnabled() {
		bluewiz.unit.AggressionLevel(0.16)
		bluewiz.unit.WalkTo(BlueFlag.Pos())
	} else {
		bluewiz.BlueTeamCheckAttackOrDefend()
	}

}

func (bluewiz *BlueWizard) BlueTeamCheckAttackOrDefend() {
	if bluewiz.unit == BlueTeamTank {
		bluewiz.unit.AggressionLevel(0.16)
		bluewiz.unit.Guard(BlueBase.Pos(), BlueBase.Pos(), 20)
	} else {
		bluewiz.unit.AggressionLevel(0.83)
		bluewiz.unit.WalkTo(RedFlag.Pos())
	}
}
