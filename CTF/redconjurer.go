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
	con := &RedConjurer{}
	con.init()
	return con
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

func (con *RedConjurer) init() {
	// Reset spells ConBot.
	con.spells.Ready = true
	// Debuff spells.
	con.spells.SlowReady = true
	con.spells.StunReady = true
	// Offensive spells.
	con.spells.MeteorReady = true
	con.spells.FistOfVengeanceReady = true
	con.spells.PixieSwarmReady = true
	con.spells.ForceOfNatureReady = true
	con.spells.ToxicCloudReady = true
	// Defensive spells.
	con.spells.BlinkReady = true
	// Summons.
	con.spells.SummonGhostReady = true
	con.spells.SummonBomber1Ready = true
	con.spells.SummonBomber2Ready = true
	// Buff spells.
	con.spells.InfravisionReady = true
	con.spells.VampirismReady = true
	con.spells.ProtFromFireReady = true
	con.spells.ProtFromPoisonReady = true
	con.spells.ProtFromShockReady = true

	// Create ConBot.
	con.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointRed"))
	con.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	con.unit.SetMaxHealth(100)
	con.unit.SetStrength(55)
	con.unit.SetBaseSpeed(88)
	con.spells.isAlive = true
	// Set Team.
	con.unit.SetOwner(TeamRed)
	// Create ConBot mouse cursor.
	con.target = TeamBlue
	con.cursor = con.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	con.reactionTime = 15
	// Set ConBot properties.
	con.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	con.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	con.unit.MonsterStatusEnable(object.MonStatusAlert)
	con.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(4), func() {
		con.unit.AggressionLevel(0.83)
	})
	con.unit.Hunt()
	con.unit.ResumeLevel(0.8)
	con.unit.RetreatLevel(0.4)
	// Create and equip ConBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	con.items.StreetSneakers = ns.CreateObject("StreetSneakers", con.unit)
	con.items.StreetPants = ns.CreateObject("StreetPants", con.unit)
	con.items.StreetShirt = ns.CreateObject("StreetShirt", con.unit)
	con.unit.Equip(con.items.StreetPants)
	con.unit.Equip(con.items.StreetShirt)
	con.unit.Equip(con.items.StreetSneakers)
	// Buff on respawn.
	con.buffInitial()
	// Enemy sighted.
	con.unit.OnEvent(ns.EventEnemySighted, con.onEnemySighted)
	// On Collision.
	con.unit.OnEvent(ns.EventCollision, con.onCollide)
	// Retreat.
	con.unit.OnEvent(ns.EventRetreat, con.onRetreat)
	// Enemy lost.
	con.unit.OnEvent(ns.EventLostEnemy, con.onLostEnemy)
	// On death.
	con.unit.OnEvent(ns.EventDeath, con.onDeath)
	// On heard.
	con.unit.OnEvent(ns.EventEnemyHeard, con.onEnemyHeard)
	// Looking for enemies.
	con.unit.OnEvent(ns.EventLookingForEnemy, con.onLookingForTarget)
	//con.unit.OnEvent(ns.EventChangeFocus, con.onChangeFocus)
	con.unit.OnEvent(ns.EventEndOfWaypoint, con.onEndOfWaypoint)
}

func (con *RedConjurer) onEndOfWaypoint() {
	con.RedTeamCheckAttackOrDefend()
}

func (con *RedConjurer) buffInitial() {
	con.castVampirism()
}

func (con *RedConjurer) onLookingForTarget() {
	con.castInfravision()
}

func (con *RedConjurer) onEnemyHeard() {
	con.castForceOfNature()
}

func (con *RedConjurer) onEnemySighted() {
	con.target = ns.GetCaller()
	con.castForceOfNature()
}

func (con *RedConjurer) onCollide() {
	if con.spells.isAlive {
		con.RedTeamPickUpBlueFlag()
		con.RedTeamCaptureTheBlueFlag()
		con.RedTeamRetrievedRedFlag()
	}
}

func (con *RedConjurer) onRetreat() {
	con.castBlink()
}

func (con *RedConjurer) onLostEnemy() {
	con.castInfravision()
	con.RedTeamWalkToRedFlag()
}

func (con *RedConjurer) onDeath() {
	con.spells.isAlive = false
	con.spells.Ready = false
	con.RedTeamDropFlag()
	con.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, con.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, con.unit)
		con.unit.Delete()
		con.items.StreetPants.Delete()
		con.items.StreetShirt.Delete()
		con.items.StreetSneakers.Delete()
		con.init()
	})
}

func (con *RedConjurer) Update() {
	con.findLoot()
	if con.unit.HasEnchant(enchant.ANTI_MAGIC) {
		con.spells.Ready = true
	}
	if con.unit.HasEnchant(enchant.HELD) || con.unit.HasEnchant(enchant.SLOWED) {
		con.castBlink()
	}
	if con.target.HasEnchant(enchant.HELD) || con.target.HasEnchant(enchant.SLOWED) {
		if con.unit.CanSee(con.target) {
			con.castFistOfVengeance()
		}
	}
	if con.spells.Ready && con.unit.CanSee(con.target) {
		con.castStun()
		con.castPixieSwarm()
		con.castToxicCloud()
		con.castSlow()
		con.castMeteor()

	}
	if !con.unit.CanSee(con.target) && con.spells.Ready {
		con.castVampirism()
		con.castProtectionFromShock()
		con.castProtectionFromFire()
		con.castProtectionFromPoison()
		con.summonBomber1()
		con.summonBomber2()
	}
}

func (con *RedConjurer) findLoot() {
	const dist = 75
	// Weapons.
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
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
		if con.unit.CanSee(item) {
			con.unit.Equip(item)
		}
	}
	// Quiver.
	quiver := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
		ns.HasTypeName{
			// Quiver.
			"Quiver",
		},
	)
	for _, item := range quiver {
		if con.unit.CanSee(item) {
			con.unit.Pickup(item)
		}
	}
	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
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
		if con.unit.CanSee(item) {
			con.unit.Equip(item)
		}
	}
}

func (con *RedConjurer) castVampirism() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.VampirismReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.VampirismReady = false
						ns.CastSpell(spell.VAMPIRISM, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Vampirism cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							con.spells.VampirismReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.ProtFromFireReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.ProtFromFireReady = false
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							con.spells.ProtFromFireReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castProtectionFromPoison() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.ProtFromPoisonReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownLeft, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.ProtFromPoisonReady = false
						ns.CastSpell(spell.PROTECTION_FROM_POISON, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Protection From Poison cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							con.spells.ProtFromPoisonReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.ProtFromShockReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.ProtFromShockReady = false
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							con.spells.ProtFromShockReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castPixieSwarm() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.PixieSwarmReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhDown, PhRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.PixieSwarmReady = false
						ns.CastSpell(spell.PIXIE_SWARM, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Pixie Swarm cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							con.spells.PixieSwarmReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castFistOfVengeance() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.FistOfVengeanceReady && con.spells.Ready {
		// Select target.
		con.cursor = con.target.Pos()
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.FistOfVengeanceReady = false
						ns.CastSpell(spell.FIST, con.unit, con.cursor)
						// Global cooldown.
						con.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Fist Of Vengeance cooldown.
							con.spells.FistOfVengeanceReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castForceOfNature() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.ForceOfNatureReady && con.spells.Ready {
		// Select target.
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDownRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.spells.ForceOfNatureReady = false
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(30))
						ns.CastSpell(spell.FORCE_OF_NATURE, con.unit, con.target)
						// Global cooldown.
						con.spells.Ready = true
						// Force of Nature cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							con.spells.ForceOfNatureReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castBlink() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.BlinkReady && con.unit != RedTeamTank {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.BlinkReady = false
						ns.NewTrap(con.unit, spell.BLINK)
						// Global cooldown.
						con.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							con.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castStun() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.StunReady && con.spells.Ready {
		// Select target.
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.StunReady = false
						ns.CastSpell(spell.STUN, con.unit, con.target)
						// Global cooldown.
						con.spells.Ready = true
						ns.NewTimer(ns.Seconds(5), func() {
							// Stun cooldown.
							con.spells.StunReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castToxicCloud() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.ToxicCloudReady && con.spells.Ready {
		// Select target.
		con.cursor = con.target.Pos()
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpRight, PhDownLeft, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.ToxicCloudReady = false
						ns.CastSpell(spell.TOXIC_CLOUD, con.unit, con.cursor)
						// Global cooldown.
						con.spells.Ready = true
						// Toxic Cloud cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							con.spells.ToxicCloudReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castSlow() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.SlowReady && con.spells.Ready {
		// Select target.
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.SlowReady = false
						ns.CastSpell(spell.SLOW, con.unit, con.target)
						// Global cooldown.
						con.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							con.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castMeteor() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.MeteorReady && con.spells.Ready {
		// Select target.
		con.cursor = con.target.Pos()
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDownLeft, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.MeteorReady = false
						ns.CastSpell(spell.METEOR, con.unit, con.cursor)
						// Global cooldown.
						con.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Meteor cooldown.
							con.spells.MeteorReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) castInfravision() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.InfravisionReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhRight, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.InfravisionReady = false
						ns.CastSpell(spell.INFRAVISION, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Invravision cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							con.spells.InfravisionReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) summonGhost() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.SummonGhostReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.SummonGhostReady = false
						ns.CastSpell(spell.SUMMON_GHOST, con.unit, con.unit)
						// Global cooldown.
						con.spells.Ready = true
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							con.spells.SummonGhostReady = true
						})
					}
				})
			}
		})
	}
}

func (con *RedConjurer) summonBomber1() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.SummonBomber1Ready {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(con.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(con.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(con.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
															con.spells.SummonBomber1Ready = false
															con.bomber1 = ns.CreateObject("Bomber", con.unit)
															ns.AudioEvent("BomberSummon", con.bomber1)
															con.bomber1.SetOwner(con.unit)
															con.bomber1.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	con.spells.SummonBomber1Ready = true
																})
															})
															con.bomber1.Follow(con.unit)
															con.bomber1.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															con.bomber1.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																con.bomber1.Attack(con.target)
															})
															con.bomber1.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																con.bomber1.Follow(con.unit)
															})
															// Global cooldown.
															con.spells.Ready = true
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

func (con *RedConjurer) summonBomber2() {
	// Check if cooldowns are ready.
	if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.SummonBomber2Ready {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(con.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(con.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(con.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
															con.spells.SummonBomber2Ready = false
															con.bomber2 = ns.CreateObject("Bomber", con.unit)
															ns.AudioEvent("BomberSummon", con.bomber2)
															con.bomber2.SetOwner(con.unit)
															con.bomber2.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	con.spells.SummonBomber2Ready = true
																})
															})
															con.bomber2.Follow(con.unit)
															con.bomber2.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															con.bomber2.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																con.bomber2.Attack(con.target)
															})
															con.bomber2.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																con.bomber2.Follow(con.unit)
															})
															// Global cooldown.
															con.spells.Ready = true
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
func (con *RedConjurer) RedTeamPickUpBlueFlag() {
	if ns.GetCaller() == BlueFlag {
		BlueFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		RedTeamTank = con.unit
		RedTeamTank.AggressionLevel(0.16)
		RedTeamTank.WalkTo(RedBase.Pos())
		ns.PrintStrToAll("Team Red has the Blue flag!")
	}
}

// Capture the flag.
func (con *RedConjurer) RedTeamCaptureTheBlueFlag() {
	if ns.GetCaller() == RedFlag && RedFlagIsAtBase && con.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagCapture, RedTeamTank) // <----- replace with all players

		RedTeamTank = TeamRed
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[0].ChangeScore(+1)
		}
		FlagReset()
		con.unit.AggressionLevel(0.83)
		con.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has captured the Blue flag!")
	}
}

// Retrieve own flag.
func (con *RedConjurer) RedTeamRetrievedRedFlag() {
	if ns.GetCaller() == RedFlag && !RedFlagIsAtBase {
		RedFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		RedFlag.SetPos(ns.Waypoint("RedFlagStart").Pos())
		con.unit.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Red has retrieved the flag!")
		RedTeamTank.WalkTo(RedFlag.Pos())
	}
}

// Drop flag.
func (con *RedConjurer) RedTeamDropFlag() {
	if con.unit == RedTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		BlueFlag.Enable(true)
		RedTeamTank = TeamRed
		ns.PrintStrToAll("Team Red has dropped the Blue flag!")
	}
}

// CTF behaviour.

// When enemy is killed check to see if the flag is dropped, if so get it.
func (con *RedConjurer) RedTeamWalkToRedFlag() {
	if !RedFlagIsAtBase && RedFlag.IsEnabled() {
		con.unit.AggressionLevel(0.16)
		con.unit.WalkTo(BlueFlag.Pos())
	} else {
		con.RedTeamCheckAttackOrDefend()
	}
}

// At the end of waypoint see defend if tank, attack if not.
func (con *RedConjurer) RedTeamCheckAttackOrDefend() {
	if con.unit == RedTeamTank {
		con.unit.AggressionLevel(0.16)
		con.unit.Guard(RedBase.Pos(), RedBase.Pos(), 20)
	} else {
		con.unit.AggressionLevel(0.83)
		con.unit.WalkTo(BlueFlag.Pos())
	}
}
