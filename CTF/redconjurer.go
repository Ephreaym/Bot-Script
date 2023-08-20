package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewConjurer creates a new RedConjurer bot.
func NewRedConjurer() *RedConjurer {
	redcon := &RedConjurer{}
	redcon.init()
	return redcon
}

// RedConjurer bot class.
type RedConjurer struct {
	unit    ns.Obj
	cursor  ns.Pointf
	target  ns.Obj
	bomber1 ns.Obj
	bomber2 ns.Obj
	items   struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
	}
	spells struct {
		isAlive              bool
		Ready                bool // Duration unknown.
		InfravisionReady     bool // Duration is 30 seconds.
		VampirismReady       bool // Duration is 30 seconds.
		BlinkReady           bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		FistOfVengeanceReady bool // No real cooldown, mana cost 60.
		StunReady            bool // No real cooldown.
		SummonBomber1Ready   bool // No real cooldown.
		SummonBomber2Ready   bool
		SummonGhostReady     bool
		ProtFromFireReady    bool // Duration is 60 seconds.
		ProtFromPoisonReady  bool
		ProtFromShockReady   bool
		PixieSwarmReady      bool
		ForceOfNatureReady   bool
		ToxicCloudReady      bool // 60 mana.
		SlowReady            bool
		MeteorReady          bool
	}
	reactionTime int
}

func (redcon *RedConjurer) init() {
	// Reset spells ConBot.
	redcon.spells.Ready = true
	// Debuff spells.
	redcon.spells.SlowReady = true
	redcon.spells.StunReady = true
	// Offensive spells.
	redcon.spells.MeteorReady = true
	redcon.spells.FistOfVengeanceReady = true
	redcon.spells.PixieSwarmReady = true
	redcon.spells.ForceOfNatureReady = true
	redcon.spells.ToxicCloudReady = true
	// Defensive spells.
	redcon.spells.BlinkReady = true
	// Summons.
	redcon.spells.SummonGhostReady = true
	redcon.spells.SummonBomber1Ready = true
	redcon.spells.SummonBomber2Ready = true
	// Buff spells.
	redcon.spells.InfravisionReady = true
	redcon.spells.VampirismReady = true
	redcon.spells.ProtFromFireReady = true
	redcon.spells.ProtFromPoisonReady = true
	redcon.spells.ProtFromShockReady = true

	// Create ConBot.
	redcon.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointRed"))
	redcon.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	redcon.unit.SetMaxHealth(100)
	redcon.unit.SetStrength(55)
	redcon.unit.SetBaseSpeed(88)
	redcon.spells.isAlive = true
	// Set Team.
	redcon.unit.SetOwner(TeamRed)
	// Create ConBot mouse cursor.
	redcon.target = TeamBlue
	redcon.cursor = redcon.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	redcon.reactionTime = 15
	// Set ConBot properties.
	redcon.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	redcon.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	redcon.unit.MonsterStatusEnable(object.MonStatusAlert)
	redcon.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(4), func() {
		redcon.unit.AggressionLevel(0.83)
	})
	redcon.unit.Hunt()
	redcon.unit.ResumeLevel(0.8)
	redcon.unit.RetreatLevel(0.4)
	// Create and equip ConBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	redcon.items.StreetSneakers = ns.CreateObject("StreetSneakers", redcon.unit)
	redcon.items.StreetPants = ns.CreateObject("StreetPants", redcon.unit)
	redcon.items.StreetShirt = ns.CreateObject("StreetShirt", redcon.unit)
	redcon.unit.Equip(redcon.items.StreetPants)
	redcon.unit.Equip(redcon.items.StreetShirt)
	redcon.unit.Equip(redcon.items.StreetSneakers)
	// Buff on respawn.
	redcon.buffInitial()
	// Enemy sighted.
	redcon.unit.OnEvent(ns.EventEnemySighted, redcon.onEnemySighted)
	// On Collision.
	redcon.unit.OnEvent(ns.EventCollision, redcon.onCollide)
	// Retreat.
	redcon.unit.OnEvent(ns.EventRetreat, redcon.onRetreat)
	// Enemy lost.
	redcon.unit.OnEvent(ns.EventLostEnemy, redcon.onLostEnemy)
	// On death.
	redcon.unit.OnEvent(ns.EventDeath, redcon.onDeath)
	// On heard.
	redcon.unit.OnEvent(ns.EventEnemyHeard, redcon.onEnemyHeard)
	// Looking for enemies.
	redcon.unit.OnEvent(ns.EventLookingForEnemy, redcon.onLookingForTarget)
	//redcon.unit.OnEvent(ns.EventChangeFocus, redcon.onChangeFocus)
	redcon.unit.OnEvent(ns.EventEndOfWaypoint, redcon.onEndOfWaypoint)
}

func (redcon *RedConjurer) onEndOfWaypoint() {
	redcon.RedTeamCheckAttackOrDefend()
}

func (redcon *RedConjurer) buffInitial() {
	redcon.castVampirism()
}

func (redcon *RedConjurer) onLookingForTarget() {
	redcon.castInfravision()
}

func (redcon *RedConjurer) onEnemyHeard() {
	redcon.castForceOfNature()
}

func (redcon *RedConjurer) onEnemySighted() {
	redcon.target = ns.GetCaller()
	redcon.castForceOfNature()
}

func (redcon *RedConjurer) onCollide() {
	if redcon.spells.isAlive {
		redcon.RedTeamPickUpBlueFlag()
		redcon.RedTeamCaptureTheBlueFlag()
		redcon.RedTeamRetrievedRedFlag()
	}
}

func (redcon *RedConjurer) onRetreat() {
	redcon.castBlink()
}

func (redcon *RedConjurer) onLostEnemy() {
	redcon.castInfravision()
	redcon.RedTeamWalkToRedFlag()
}

func (redcon *RedConjurer) onDeath() {
	redcon.spells.isAlive = false
	redcon.spells.Ready = false
	redcon.RedTeamDropFlag()
	redcon.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, redcon.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, redcon.unit)
		redcon.unit.Delete()
		redcon.items.StreetPants.Delete()
		redcon.items.StreetShirt.Delete()
		redcon.items.StreetSneakers.Delete()
		redcon.init()
	})
}

func (redcon *RedConjurer) Update() {
	redcon.findLoot()
	if redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
		redcon.spells.Ready = true
	}
	if redcon.unit.HasEnchant(enchant.HELD) || redcon.unit.HasEnchant(enchant.SLOWED) {
		redcon.castBlink()
	}
	if redcon.target.HasEnchant(enchant.HELD) || redcon.target.HasEnchant(enchant.SLOWED) {
		if redcon.unit.CanSee(redcon.target) {
			redcon.castFistOfVengeance()
		}
	}
	if redcon.spells.Ready && redcon.unit.CanSee(redcon.target) {
		redcon.castStun()
		redcon.castPixieSwarm()
		redcon.castToxicCloud()
		redcon.castSlow()
		redcon.castMeteor()

	}
	if !redcon.unit.CanSee(redcon.target) && redcon.spells.Ready {
		redcon.castVampirism()
		redcon.castProtectionFromShock()
		redcon.castProtectionFromFire()
		redcon.castProtectionFromPoison()
		redcon.summonBomber1()
		redcon.summonBomber2()
	}
}

func (redcon *RedConjurer) findLoot() {
	const dist = 75
	// Weapons.
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: redcon.unit, R: dist},
		ns.HasTypeName{
			// Wands.
			"InfinitePainWand", "LesserFireballWand",
			//"SulphorousShowerWand",
			//"SulphorousFlareWand",
			//"StaffWooden",

			// Crossbow and Bow.
			"CrossBow",
			"Bow",
			"Quiver",
		},
	)
	for _, item := range weapons {
		if redcon.unit.CanSee(item) {
			redcon.unit.Equip(item)
		}
	}
	// Quiver.
	quiver := ns.FindAllObjects(
		ns.InCirclef{Center: redcon.unit, R: dist},
		ns.HasTypeName{
			// Quiver.
			"Quiver",
		},
	)
	for _, item := range quiver {
		if redcon.unit.CanSee(item) {
			redcon.unit.Pickup(item)
		}
	}
	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: redcon.unit, R: dist},
		ns.HasTypeName{
			// RedConjurer Helm.
			"ConjurerHelm",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor", "LeatherHelm", "LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if redcon.unit.CanSee(item) {
			redcon.unit.Equip(item)
		}
	}
}

func (redcon *RedConjurer) castVampirism() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.VampirismReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.VampirismReady = false
						ns.CastSpell(spell.VAMPIRISM, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Vampirism cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							redcon.spells.VampirismReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.ProtFromFireReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.ProtFromFireReady = false
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redcon.spells.ProtFromFireReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castProtectionFromPoison() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.ProtFromPoisonReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhLeft, PhRight, PhDownLeft, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.ProtFromPoisonReady = false
						ns.CastSpell(spell.PROTECTION_FROM_POISON, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Protection From Poison cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redcon.spells.ProtFromPoisonReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.ProtFromShockReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.ProtFromShockReady = false
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							redcon.spells.ProtFromShockReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castPixieSwarm() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.PixieSwarmReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhLeft, PhDown, PhRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.PixieSwarmReady = false
						ns.CastSpell(spell.PIXIE_SWARM, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Pixie Swarm cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							redcon.spells.PixieSwarmReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castFistOfVengeance() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.unit.CanSee(redcon.target) && redcon.spells.FistOfVengeanceReady && redcon.spells.Ready {
		// Select target.
		redcon.cursor = redcon.target.Pos()
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(redcon.reactionTime))
						redcon.spells.FistOfVengeanceReady = false
						ns.CastSpell(spell.FIST, redcon.unit, redcon.cursor)
						// Global cooldown.
						redcon.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Fist Of Vengeance cooldown.
							redcon.spells.FistOfVengeanceReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castForceOfNature() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.ForceOfNatureReady && redcon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhDownRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.spells.ForceOfNatureReady = false
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(30))
						ns.CastSpell(spell.FORCE_OF_NATURE, redcon.unit, redcon.target)
						// Global cooldown.
						redcon.spells.Ready = true
						// Force of Nature cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							redcon.spells.ForceOfNatureReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castBlink() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.BlinkReady && redcon.unit != RedTeamTank {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.BlinkReady = false
						ns.NewTrap(redcon.unit, spell.BLINK)
						// Global cooldown.
						redcon.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							redcon.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castStun() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.unit.CanSee(redcon.target) && redcon.spells.StunReady && redcon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(redcon.reactionTime))
						redcon.spells.StunReady = false
						ns.CastSpell(spell.STUN, redcon.unit, redcon.target)
						// Global cooldown.
						redcon.spells.Ready = true
						ns.NewTimer(ns.Seconds(5), func() {
							// Stun cooldown.
							redcon.spells.StunReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castToxicCloud() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.unit.CanSee(redcon.target) && redcon.spells.ToxicCloudReady && redcon.spells.Ready {
		// Select target.
		redcon.cursor = redcon.target.Pos()
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhDownLeft, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(redcon.reactionTime))
						redcon.spells.ToxicCloudReady = false
						ns.CastSpell(spell.TOXIC_CLOUD, redcon.unit, redcon.cursor)
						// Global cooldown.
						redcon.spells.Ready = true
						// Toxic Cloud cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							redcon.spells.ToxicCloudReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castSlow() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.unit.CanSee(redcon.target) && redcon.spells.SlowReady && redcon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(redcon.reactionTime))
						redcon.spells.SlowReady = false
						ns.CastSpell(spell.SLOW, redcon.unit, redcon.target)
						// Global cooldown.
						redcon.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							redcon.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castMeteor() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.unit.CanSee(redcon.target) && redcon.spells.MeteorReady && redcon.spells.Ready {
		// Select target.
		redcon.cursor = redcon.target.Pos()
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhDownLeft, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						redcon.unit.LookAtObject(redcon.target)
						redcon.unit.Pause(ns.Frames(redcon.reactionTime))
						redcon.spells.MeteorReady = false
						ns.CastSpell(spell.METEOR, redcon.unit, redcon.cursor)
						// Global cooldown.
						redcon.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Meteor cooldown.
							redcon.spells.MeteorReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) castInfravision() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.InfravisionReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhRight, PhLeft, PhRight, PhLeft}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.InfravisionReady = false
						ns.CastSpell(spell.INFRAVISION, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Invravision cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							redcon.spells.InfravisionReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) summonGhost() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.SummonGhostReady {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(redcon.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						redcon.spells.SummonGhostReady = false
						ns.CastSpell(spell.SUMMON_GHOST, redcon.unit, redcon.unit)
						// Global cooldown.
						redcon.spells.Ready = true
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							redcon.spells.SummonGhostReady = true
						})
					}
				})
			}
		})
	}
}

func (redcon *RedConjurer) summonBomber1() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.SummonBomber1Ready {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(redcon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(redcon.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
															redcon.spells.SummonBomber1Ready = false
															redcon.bomber1 = ns.CreateObject("Bomber", redcon.unit)
															ns.AudioEvent("BomberSummon", redcon.bomber1)
															redcon.bomber1.SetOwner(redcon.unit)
															redcon.bomber1.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	redcon.spells.SummonBomber1Ready = true
																})
															})
															redcon.bomber1.Follow(redcon.unit)
															redcon.bomber1.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															redcon.bomber1.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																redcon.bomber1.Attack(redcon.target)
															})
															redcon.bomber1.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																redcon.bomber1.Follow(redcon.unit)
															})
															// Global cooldown.
															redcon.spells.Ready = true
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

func (redcon *RedConjurer) summonBomber2() {
	// Check if cooldowns are ready.
	if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) && redcon.spells.Ready && redcon.spells.SummonBomber2Ready {
		// Trigger cooldown.
		redcon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(redcon.reactionTime), func() {
			// Check for War Cry before chant.
			if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(redcon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(redcon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(redcon.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if redcon.spells.isAlive && !redcon.unit.HasEnchant(enchant.ANTI_MAGIC) {
															redcon.spells.SummonBomber2Ready = false
															redcon.bomber2 = ns.CreateObject("Bomber", redcon.unit)
															ns.AudioEvent("BomberSummon", redcon.bomber2)
															redcon.bomber2.SetOwner(redcon.unit)
															redcon.bomber2.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	redcon.spells.SummonBomber2Ready = true
																})
															})
															redcon.bomber2.Follow(redcon.unit)
															redcon.bomber2.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															redcon.bomber2.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																redcon.bomber2.Attack(redcon.target)
															})
															redcon.bomber2.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																redcon.bomber2.Follow(redcon.unit)
															})
															// Global cooldown.
															redcon.spells.Ready = true
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

// ---------------------------------- CTF BOT SCRIPT ------------------------------------//
// CTF game mechanics.
// Pick up the enemy flag.
func (redcon *RedConjurer) RedTeamPickUpBlueFlag() {
	if ns.GetCaller() == BlueFlag {
		BlueFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		RedTeamTank = redcon.unit
		RedTeamTank.AggressionLevel(0.16)
		RedTeamTank.WalkTo(RedBase.Pos())
		ns.PrintStrToAll("Team Red has the Blue flag!")
	}
}

// Capture the flag.
func (redcon *RedConjurer) RedTeamCaptureTheBlueFlag() {
	if ns.GetCaller() == RedFlag && RedFlagIsAtBase && redcon.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagCapture, RedTeamTank) // <----- replace with all players

		RedTeamTank = TeamRed
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[0].ChangeScore(+1)
		}
		FlagReset()
		redcon.unit.AggressionLevel(0.83)
		redcon.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has captured the Blue flag!")
	}
}

// Retrieve own flag.
func (redcon *RedConjurer) RedTeamRetrievedRedFlag() {
	if ns.GetCaller() == RedFlag && !RedFlagIsAtBase {
		RedFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
		redcon.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has retrieved the flag!")
		RedTeamTank.WalkTo(RedFlag.Pos())
	}
}

// Drop flag.
func (redcon *RedConjurer) RedTeamDropFlag() {
	if redcon.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		BlueFlag.Enable(true)
		RedTeamTank = TeamRed
		ns.PrintStrToAll("Team Red has dropped the Blue flag!")
	}
}

// CTF behaviour.

// When enemy is killed check to see if the flag is dropped, if so get it.
func (redcon *RedConjurer) RedTeamWalkToRedFlag() {
	if !RedFlagIsAtBase && RedFlag.IsEnabled() {
		redcon.unit.AggressionLevel(0.16)
		redcon.unit.WalkTo(BlueFlag.Pos())
	} else {
		redcon.RedTeamCheckAttackOrDefend()
	}
}

// At the end of waypoint see defend if tank, attack if not.
func (redcon *RedConjurer) RedTeamCheckAttackOrDefend() {
	if redcon.unit == RedTeamTank {
		redcon.unit.AggressionLevel(0.16)
		redcon.unit.Guard(RedBase.Pos(), RedBase.Pos(), 20)
	} else {
		redcon.unit.AggressionLevel(0.83)
		redcon.unit.WalkTo(BlueFlag.Pos())
	}
}
