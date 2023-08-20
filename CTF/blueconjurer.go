package EndGameBW

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewConjurer creates a new BlueConjurer bot.
func NewBlueConjurer() *BlueConjurer {
	bluecon := &BlueConjurer{}
	bluecon.init()
	return bluecon
}

// BlueConjurer bot class.
type BlueConjurer struct {
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

func (bluecon *BlueConjurer) init() {
	// Reset spells ConBot.
	bluecon.spells.Ready = true
	// Debuff spells.
	bluecon.spells.SlowReady = true
	bluecon.spells.StunReady = true
	// Offensive spells.
	bluecon.spells.MeteorReady = true
	bluecon.spells.FistOfVengeanceReady = true
	bluecon.spells.PixieSwarmReady = true
	bluecon.spells.ForceOfNatureReady = true
	bluecon.spells.ToxicCloudReady = true
	// Defensive spells.
	bluecon.spells.BlinkReady = true
	// Summons.
	bluecon.spells.SummonGhostReady = true
	bluecon.spells.SummonBomber1Ready = true
	bluecon.spells.SummonBomber2Ready = true
	// Buff spells.
	bluecon.spells.InfravisionReady = true
	bluecon.spells.VampirismReady = true
	bluecon.spells.ProtFromFireReady = true
	bluecon.spells.ProtFromPoisonReady = true
	bluecon.spells.ProtFromShockReady = true

	// Create ConBot.
	bluecon.unit = ns.CreateObject("NPC", ns.Waypoint("BotSpawnPointBlue"))
	bluecon.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	bluecon.unit.SetMaxHealth(100)
	bluecon.unit.SetStrength(55)
	bluecon.unit.SetBaseSpeed(88)
	bluecon.spells.isAlive = true
	// Set Team.
	bluecon.unit.SetOwner(TeamBlue)
	// Create ConBot mouse cursor.
	bluecon.target = TeamRed
	bluecon.cursor = bluecon.target.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	bluecon.reactionTime = 15
	// Set ConBot properties.
	bluecon.unit.MonsterStatusEnable(object.MonStatusAlwaysRun)
	bluecon.unit.MonsterStatusEnable(object.MonStatusCanCastSpells)
	bluecon.unit.MonsterStatusEnable(object.MonStatusAlert)
	bluecon.unit.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(4), func() {
		bluecon.unit.AggressionLevel(0.83)
	})
	bluecon.unit.Hunt()
	bluecon.unit.ResumeLevel(0.8)
	bluecon.unit.RetreatLevel(0.4)
	// Create and equip ConBot starting equipment. TODO: Change location of item creation OR stop them from respawning automatically.
	bluecon.items.StreetSneakers = ns.CreateObject("StreetSneakers", bluecon.unit)
	bluecon.items.StreetPants = ns.CreateObject("StreetPants", bluecon.unit)
	bluecon.items.StreetShirt = ns.CreateObject("StreetShirt", bluecon.unit)
	bluecon.unit.Equip(bluecon.items.StreetPants)
	bluecon.unit.Equip(bluecon.items.StreetShirt)
	bluecon.unit.Equip(bluecon.items.StreetSneakers)
	// Buff on respawn.
	bluecon.buffInitial()
	// Enemy sighted.
	bluecon.unit.OnEvent(ns.EventEnemySighted, bluecon.onEnemySighted)
	// On Collision.
	bluecon.unit.OnEvent(ns.EventCollision, bluecon.onCollide)
	// Retreat.
	bluecon.unit.OnEvent(ns.EventRetreat, bluecon.onRetreat)
	// Enemy lost.
	bluecon.unit.OnEvent(ns.EventLostEnemy, bluecon.onLostEnemy)
	// On death.
	bluecon.unit.OnEvent(ns.EventDeath, bluecon.onDeath)
	// On heard.
	bluecon.unit.OnEvent(ns.EventEnemyHeard, bluecon.onEnemyHeard)
	// Looking for enemies.
	bluecon.unit.OnEvent(ns.EventLookingForEnemy, bluecon.onLookingForTarget)
	//bluecon.unit.OnEvent(ns.EventChangeFocus, bluecon.onChangeFocus)
	bluecon.unit.OnEvent(ns.EventEndOfWaypoint, bluecon.onEndOfWaypoint)
}

func (bluecon *BlueConjurer) onEndOfWaypoint() {
	bluecon.BlueTeamCheckAttackOrDefend()
}

func (bluecon *BlueConjurer) buffInitial() {
	bluecon.castVampirism()
}

func (bluecon *BlueConjurer) onLookingForTarget() {
	bluecon.castInfravision()
}

func (bluecon *BlueConjurer) onEnemyHeard() {
	bluecon.castForceOfNature()
}

func (bluecon *BlueConjurer) onEnemySighted() {
	bluecon.target = ns.GetCaller()
	bluecon.castForceOfNature()
}

func (bluecon *BlueConjurer) onCollide() {
	if bluecon.spells.isAlive {
		bluecon.BlueTeamPickUpRedFlag()
		bluecon.BlueTeamCaptureTheRedFlag()
		bluecon.BlueTeamRetrievedBlueFlag()
	}
}

func (bluecon *BlueConjurer) onRetreat() {
	bluecon.castBlink()
}

func (bluecon *BlueConjurer) onLostEnemy() {
	bluecon.castInfravision()
	bluecon.BlueTeamWalkToBlueFlag()
}

func (bluecon *BlueConjurer) onDeath() {
	bluecon.spells.isAlive = false
	bluecon.spells.Ready = false
	bluecon.BlueTeamDropFlag()
	bluecon.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, bluecon.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, bluecon.unit)
		bluecon.unit.Delete()
		bluecon.items.StreetPants.Delete()
		bluecon.items.StreetShirt.Delete()
		bluecon.items.StreetSneakers.Delete()
		bluecon.init()
	})
}

func (bluecon *BlueConjurer) Update() {
	bluecon.findLoot()
	if bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
		bluecon.spells.Ready = true
	}
	if bluecon.unit.HasEnchant(enchant.HELD) || bluecon.unit.HasEnchant(enchant.SLOWED) {
		bluecon.castBlink()
	}
	if bluecon.target.HasEnchant(enchant.HELD) || bluecon.target.HasEnchant(enchant.SLOWED) {
		if bluecon.unit.CanSee(bluecon.target) {
			bluecon.castFistOfVengeance()
		}
	}
	if bluecon.spells.Ready && bluecon.unit.CanSee(bluecon.target) {
		bluecon.castStun()
		bluecon.castPixieSwarm()
		bluecon.castToxicCloud()
		bluecon.castSlow()
		bluecon.castMeteor()

	}
	if !bluecon.unit.CanSee(bluecon.target) && bluecon.spells.Ready {
		bluecon.castVampirism()
		bluecon.castProtectionFromShock()
		bluecon.castProtectionFromFire()
		bluecon.castProtectionFromPoison()
		bluecon.summonBomber1()
		bluecon.summonBomber2()
	}
}

func (bluecon *BlueConjurer) findLoot() {
	const dist = 75
	// Weapons.
	weapons := ns.FindAllObjects(
		ns.InCirclef{Center: bluecon.unit, R: dist},
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
		if bluecon.unit.CanSee(item) {
			bluecon.unit.Equip(item)
		}
	}
	// Quiver.
	quiver := ns.FindAllObjects(
		ns.InCirclef{Center: bluecon.unit, R: dist},
		ns.HasTypeName{
			// Quiver.
			"Quiver",
		},
	)
	for _, item := range quiver {
		if bluecon.unit.CanSee(item) {
			bluecon.unit.Pickup(item)
		}
	}
	// Armor.
	armor := ns.FindAllObjects(
		ns.InCirclef{Center: bluecon.unit, R: dist},
		ns.HasTypeName{
			// BlueConjurer Helm.
			"ConjurerHelm",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor", "LeatherHelm", "LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if bluecon.unit.CanSee(item) {
			bluecon.unit.Equip(item)
		}
	}
}

func (bluecon *BlueConjurer) castVampirism() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.VampirismReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.VampirismReady = false
						ns.CastSpell(spell.VAMPIRISM, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Vampirism cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							bluecon.spells.VampirismReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.ProtFromFireReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.ProtFromFireReady = false
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Protection From Fire cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluecon.spells.ProtFromFireReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castProtectionFromPoison() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.ProtFromPoisonReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhLeft, PhRight, PhDownLeft, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.ProtFromPoisonReady = false
						ns.CastSpell(spell.PROTECTION_FROM_POISON, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Protection From Poison cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluecon.spells.ProtFromPoisonReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.ProtFromShockReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.ProtFromShockReady = false
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Protection From Shock cooldown.
						ns.NewTimer(ns.Seconds(60), func() {
							bluecon.spells.ProtFromShockReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castPixieSwarm() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.PixieSwarmReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhLeft, PhDown, PhRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.PixieSwarmReady = false
						ns.CastSpell(spell.PIXIE_SWARM, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Pixie Swarm cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							bluecon.spells.PixieSwarmReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castFistOfVengeance() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.unit.CanSee(bluecon.target) && bluecon.spells.FistOfVengeanceReady && bluecon.spells.Ready {
		// Select target.
		bluecon.cursor = bluecon.target.Pos()
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(bluecon.reactionTime))
						bluecon.spells.FistOfVengeanceReady = false
						ns.CastSpell(spell.FIST, bluecon.unit, bluecon.cursor)
						// Global cooldown.
						bluecon.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Fist Of Vengeance cooldown.
							bluecon.spells.FistOfVengeanceReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castForceOfNature() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.ForceOfNatureReady && bluecon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhDownRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.spells.ForceOfNatureReady = false
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(30))
						ns.CastSpell(spell.FORCE_OF_NATURE, bluecon.unit, bluecon.target)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Force of Nature cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							bluecon.spells.ForceOfNatureReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castBlink() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.BlinkReady && bluecon.unit != BlueTeamTank {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhRight, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.BlinkReady = false
						ns.NewTrap(bluecon.unit, spell.BLINK)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							bluecon.spells.BlinkReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castStun() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.unit.CanSee(bluecon.target) && bluecon.spells.StunReady && bluecon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(bluecon.reactionTime))
						bluecon.spells.StunReady = false
						ns.CastSpell(spell.STUN, bluecon.unit, bluecon.target)
						// Global cooldown.
						bluecon.spells.Ready = true
						ns.NewTimer(ns.Seconds(5), func() {
							// Stun cooldown.
							bluecon.spells.StunReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castToxicCloud() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.unit.CanSee(bluecon.target) && bluecon.spells.ToxicCloudReady && bluecon.spells.Ready {
		// Select target.
		bluecon.cursor = bluecon.target.Pos()
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhDownLeft, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(bluecon.reactionTime))
						bluecon.spells.ToxicCloudReady = false
						ns.CastSpell(spell.TOXIC_CLOUD, bluecon.unit, bluecon.cursor)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Toxic Cloud cooldown.
						ns.NewTimer(ns.Seconds(10), func() {
							bluecon.spells.ToxicCloudReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castSlow() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.unit.CanSee(bluecon.target) && bluecon.spells.SlowReady && bluecon.spells.Ready {
		// Select target.
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhDown, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(bluecon.reactionTime))
						bluecon.spells.SlowReady = false
						ns.CastSpell(spell.SLOW, bluecon.unit, bluecon.target)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							bluecon.spells.SlowReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castMeteor() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.unit.CanSee(bluecon.target) && bluecon.spells.MeteorReady && bluecon.spells.Ready {
		// Select target.
		bluecon.cursor = bluecon.target.Pos()
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhDownLeft, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						// Aim.
						bluecon.unit.LookAtObject(bluecon.target)
						bluecon.unit.Pause(ns.Frames(bluecon.reactionTime))
						bluecon.spells.MeteorReady = false
						ns.CastSpell(spell.METEOR, bluecon.unit, bluecon.cursor)
						// Global cooldown.
						bluecon.spells.Ready = true
						ns.NewTimer(ns.Seconds(10), func() {
							// Meteor cooldown.
							bluecon.spells.MeteorReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) castInfravision() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.InfravisionReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhRight, PhLeft, PhRight, PhLeft}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.InfravisionReady = false
						ns.CastSpell(spell.INFRAVISION, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Invravision cooldown.
						ns.NewTimer(ns.Seconds(30), func() {
							bluecon.spells.InfravisionReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) summonGhost() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.SummonGhostReady {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(bluecon.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
						bluecon.spells.SummonGhostReady = false
						ns.CastSpell(spell.SUMMON_GHOST, bluecon.unit, bluecon.unit)
						// Global cooldown.
						bluecon.spells.Ready = true
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							bluecon.spells.SummonGhostReady = true
						})
					}
				})
			}
		})
	}
}

func (bluecon *BlueConjurer) summonBomber1() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.SummonBomber1Ready {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(bluecon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(bluecon.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
															bluecon.spells.SummonBomber1Ready = false
															bluecon.bomber1 = ns.CreateObject("Bomber", bluecon.unit)
															ns.AudioEvent("BomberSummon", bluecon.bomber1)
															bluecon.bomber1.SetOwner(bluecon.unit)
															bluecon.bomber1.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	bluecon.spells.SummonBomber1Ready = true
																})
															})
															bluecon.bomber1.Follow(bluecon.unit)
															bluecon.bomber1.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															bluecon.bomber1.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																bluecon.bomber1.Attack(bluecon.target)
															})
															bluecon.bomber1.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																bluecon.bomber1.Follow(bluecon.unit)
															})
															// Global cooldown.
															bluecon.spells.Ready = true
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

func (bluecon *BlueConjurer) summonBomber2() {
	// Check if cooldowns are ready.
	if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) && bluecon.spells.Ready && bluecon.spells.SummonBomber2Ready {
		// Trigger cooldown.
		bluecon.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(bluecon.reactionTime), func() {
			// Check for War Cry before chant.
			if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
				// Stun chant.
				castPhonemes(bluecon.unit, []audio.Name{PhUpLeft, PhDown}, func() {
					// Pause for concentration.
					ns.NewTimer(ns.Frames(3), func() {
						// Check for War Cry before chant.
						if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
							// Poison chant.
							castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhDownLeft}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Fist Of Vengeance chant.
										castPhonemes(bluecon.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(bluecon.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if bluecon.spells.isAlive && !bluecon.unit.HasEnchant(enchant.ANTI_MAGIC) {
															bluecon.spells.SummonBomber2Ready = false
															bluecon.bomber2 = ns.CreateObject("Bomber", bluecon.unit)
															ns.AudioEvent("BomberSummon", bluecon.bomber2)
															bluecon.bomber2.SetOwner(bluecon.unit)
															bluecon.bomber2.OnEvent(ns.ObjectEvent(ns.EventDeath), func() {
																// Summon Bomber cooldown.
																ns.NewTimer(ns.Seconds(10), func() {
																	bluecon.spells.SummonBomber2Ready = true
																})
															})
															bluecon.bomber2.Follow(bluecon.unit)
															bluecon.bomber2.TrapSpells(spell.POISON, spell.FIST, spell.STUN)
															bluecon.bomber2.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																bluecon.bomber2.Attack(bluecon.target)
															})
															bluecon.bomber2.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																bluecon.bomber2.Follow(bluecon.unit)
															})
															// Global cooldown.
															bluecon.spells.Ready = true
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
func (bluecon *BlueConjurer) BlueTeamPickUpRedFlag() {
	if ns.GetCaller() == RedFlag {
		RedFlag.Enable(false)
		ns.AudioEvent(audio.FlagPickup, ns.GetHost()) // <----- replace with all players
		// Customize code below for individual unit.
		BlueTeamTank = bluecon.unit
		BlueTeamTank.AggressionLevel(0.16)
		BlueTeamTank.WalkTo(BlueFlag.Pos())
		ns.PrintStrToAll("Team Blue has the Red flag!")
	}
}

// Capture the flag.
func (bluecon *BlueConjurer) BlueTeamCaptureTheRedFlag() {
	if ns.GetCaller() == BlueFlag && BlueFlagIsAtBase && bluecon.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagCapture, BlueTeamTank) // <----- replace with all players
		BlueTeamTank = TeamBlue
		var1 := ns.Players()
		if len(var1) > 1 {
			var1[1].ChangeScore(+1)
		}
		FlagReset()
		bluecon.unit.AggressionLevel(0.83)
		bluecon.unit.WalkTo(RedFlag.Pos())
		ns.PrintStrToAll("Team Blue has captured the Red flag!")
	}
}

// Retrieve own flag.
func (bluecon *BlueConjurer) BlueTeamRetrievedBlueFlag() {
	if ns.GetCaller() == BlueFlag && !BlueFlagIsAtBase {
		BlueFlagIsAtBase = true
		ns.AudioEvent(audio.FlagRespawn, ns.GetHost())
		BlueFlag.SetPos(ns.Waypoint("BlueFlagStart").Pos())
		bluecon.unit.WalkTo(BlueBase.Pos())
		ns.PrintStrToAll("Team Blue has retrieved the flag!")
		BlueTeamTank.WalkTo(BlueFlag.Pos())
	}
}

// Drop flag.
func (bluecon *BlueConjurer) BlueTeamDropFlag() {
	if bluecon.unit == BlueTeamTank {
		ns.AudioEvent(audio.FlagDrop, ns.GetHost()) // <----- replace with all players
		RedFlag.Enable(true)
		BlueTeamTank = TeamBlue
		ns.PrintStrToAll("Team Blue has dropped the Red flag!")
	}
}

// CTF behaviour.
// Attack enemy tank without

func (bluecon *BlueConjurer) BlueTeamWalkToBlueFlag() {
	if !BlueFlagIsAtBase && BlueFlag.IsEnabled() {
		bluecon.unit.AggressionLevel(0.16)
		bluecon.unit.WalkTo(BlueFlag.Pos())
	} else {
		bluecon.BlueTeamCheckAttackOrDefend()
	}

}

func (bluecon *BlueConjurer) BlueTeamCheckAttackOrDefend() {
	if bluecon.unit == BlueTeamTank {
		bluecon.unit.AggressionLevel(0.16)
		bluecon.unit.Guard(BlueBase.Pos(), BlueBase.Pos(), 20)
	} else {
		bluecon.unit.AggressionLevel(0.83)
		bluecon.unit.WalkTo(RedFlag.Pos())
	}
}
