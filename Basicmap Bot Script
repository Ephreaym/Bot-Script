package basicmap

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// Declare WarAI variables. //
var WarAI ns.Obj
var WarAITarget ns.Obj
var WarAITaggedPlayer ns.Obj
var WarAILongsword ns.Obj
var WarAIWoodenShield ns.Obj
var WarAIStreetSneakers ns.Obj
var WarAIStreetPants ns.Obj
var WarAIEyeOfTheWolfCooldown = 1    // Cooldown is 20 seconds. //
var WarAIBerserkerChargeCooldown = 1 // Cooldown is 10 seconds. TODO: reset on kill. //

// Declare WizAI variables. //
var WizAI ns.Obj
var WizAITarget ns.Obj
var WizAITaggedPlayer ns.Obj
var WizAITrap ns.Obj
var WizAIStreetSneakers ns.Obj
var WizAIStreetPants ns.Obj
var WizAIStreetShirt ns.Obj
var WizAIWizardRobe ns.Obj
var WizAIMagicMissilesCooldown = 1 // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random. //
var WizAIGlobalCooldown = 1        // Duration unknown. //
var WizAIForceFieldCooldown = 1    // Duration unknown. //
var WizAIShockCooldown = 1         // No real cooldown,
var WizAISlowCooldown = 1          // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random. //
var WizAITrapCooldown = 1          // Only one trap is placed per life. //

// Declare ConAI variables. //
var ConAI ns.Obj
var ConAIStreetSneakers ns.Obj
var ConAIStreetPants ns.Obj
var ConAIStreetShirt ns.Obj
var ConGlobalCooldown = 1        // Duration unknown. //
var ConAIInfravisionCooldown = 1 // Duration is 30 seconds. //
var ConAIVampirismCooldown = 1   // Duration is 30 seconds. //
var ConAIBlinkCooldown = 1       // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random. //

func init() {
	// Load in AI scripts on map launch. //
	LoadWarAI()
	LoadWizAI()
	LoadConAI()
	ns.Music(15, 20)
}

func OnFrame() {
	// AI checks for items to loot. //
	FindLootWarAI()
	FindLootWizAI()
	FindLootConAI()
}

func FindLootConAI() {
	// Wands. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"InfinitePainWand"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LesserFireballWand"}))
	// ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"SulphorousShowerWand"}))
	// ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"SulphorousFlareWand"}))
	// ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"StaffWooden"}))
	// Crossbow. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"CrossBow"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"Quiver"}))
	// Bow. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"Bow"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"Quiver"}))
	// Conjurer Helm. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"ConjurerHelm"}))
	// Leather armor. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherArmoredBoots"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherArmor"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherHelm"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherLeggings"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherArmbands"}))
	// Cloth armor. //
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"LeatherBoots"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"MedievalCloak"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"MedievalShirt"}))
	ConAI.Equip(ns.FindObject(ns.InCirclef{Center: ConAI, R: 75}, ns.HasTypeName{"MedievalPants"}))
}

func FindLootWizAI() {
	// Wands. //
	// WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"DeathRayWand"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"FireStormWand"}))
	// WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"LesserFireballWand"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"ForceWand"}))
	// WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"SulphorousShowerWand"}))
	// WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"SulphorousFlareWand"}))
	// WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"StaffWooden"}))
	// Wizard armor. //
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"WizardHelm"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"WizardRobe"}))
	// Cloth armor. //
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"LeatherBoots"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"MedievalCloak"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"MedievalShirt"}))
	WizAI.Equip(ns.FindObject(ns.InCirclef{Center: WizAI, R: 75}, ns.HasTypeName{"MedievalPants"}))
}

func FindLootWarAI() {
	// Weapon. //
	// TODO: Fix crash with thrown weapons. Maybe WarAI tries to pickup the thrown weapon?
	// TODO: Setup different builds and tactics / voices / dialog / chat. //
	// WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"RoundChakram"}))
	// WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"FanChakram"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"GreatSword"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"WarHammer"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"MorningStar"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"BattleAxe"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"Sword"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"OgreAxe"}))
	// WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"StaffWooden"}))
	// Plate armor. //
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"OrnateHelm"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"SteelHelm"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"Breastplate"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"PlateLeggings"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"PlateBoots"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"PlateArms"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"SteelShield"}))
	// Chainmail armor. //
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"ChainCoif"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"ChainTunic"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"ChainLeggings"}))
	// Leather armor. //
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherArmoredBoots"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherArmor"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherHelm"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherLeggings"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherArmbands"}))
	// Cloth armor. //
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"LeatherBoots"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"MedievalCloak"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"MedievalShirt"}))
	WarAI.Equip(ns.FindObject(ns.InCirclef{Center: WarAI, R: 75}, ns.HasTypeName{"MedievalPants"}))
}

// Conjurer AI script. //
func LoadConAI() {
	ConAI = ns.CreateObject("NPC", ns.GetHost())
	ConAIStreetSneakers = ns.CreateObject("StreetSneakers", ConAI)
	ConAIStreetPants = ns.CreateObject("StreetPants", ConAI)
	ConAIStreetShirt = ns.CreateObject("StreetShirt", ConAI)
	ConAI.Equip(ConAIStreetPants)
	ConAI.Equip(ConAIStreetShirt)
	ConAI.Equip(ConAIStreetSneakers)
	ConAI.Enchant("ENCHANT_INVULNERABLE", script.Frames(150))
	ConAI.SetMaxHealth(100)
	ConAI.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		ConAI.AggressionLevel(0.83)
	})
	ConAI.Hunt()
	ConAI.ResumeLevel(0.8)
	ConAI.RetreatLevel(0.2)
	// Buff on respawn. //
	if ConAIVampirismCooldown == 0 {
	} else {
		ConAIVampirismCooldown = 0
		ConGlobalCooldown = 0
		ns.AudioEvent("NPCSpellPhonemeUp", ConAI)
		ns.NewTimer(ns.Frames(3), func() {
			ns.AudioEvent("NPCSpellPhonemeSouth", ConAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeLeft", ConAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeRight", ConAI)
					ns.NewTimer(ns.Frames(3), func() {
						ns.CastSpell(spell.VAMPIRISM, ConAI, ConAI)
						ConGlobalCooldown = 1
					})
				})
			})
		})
		ns.NewTimer(ns.Seconds(30), func() {
			ConAIVampirismCooldown = 1
		})
	}
	// Escape. //
	ConAI.OnEvent(ns.EventRetreat, func() {
		if ConAIBlinkCooldown == 0 {
		} else {
			ConAIBlinkCooldown = 0
			ConGlobalCooldown = 0
			ns.AudioEvent("NPCSpellPhonemeRight", ConAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeLeft", ConAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.NewTimer(ns.Frames(3), func() {
						ns.AudioEvent("NPCSpellPhonemeUp", ConAI)
						ns.NewTrap(ConAI, spell.BLINK) // TODO: FIX IT so it doesn't have to be a trap. //
						ConGlobalCooldown = 1
					})
				})
			})
			ns.NewTimer(ns.Seconds(2), func() {
				ConAIBlinkCooldown = 1
			})
		}
	})
	// Enemy Lost. //
	ConAI.OnEvent(ns.EventLostEnemy, func() {
		if ConAIInfravisionCooldown == 0 {
		} else {
			ConAIInfravisionCooldown = 0
			ConGlobalCooldown = 0
			ns.AudioEvent("NPCSpellPhonemeRight", ConAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeLeft", ConAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeRight", ConAI)
					ns.NewTimer(ns.Frames(3), func() {
						ns.AudioEvent("NPCSpellPhonemeLeft", ConAI)
						ConAI.Enchant("ENCHANT_INFRAVISION", ns.Seconds(30))
						ConGlobalCooldown = 1
					})
				})
			})
			ns.NewTimer(ns.Seconds(30), func() {
				ConAIInfravisionCooldown = 1
			})
		}
	})
	// On Death. //
	ConAI.OnEvent(ns.EventDeath, func() {
		ConAI.DestroyChat()
		ns.AudioEvent("NPCDie", ConAI)
		// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement? //
		ns.GetHost().ChangeScore(+1)
		ns.NewTimer(ns.Frames(60), func() {
			ns.AudioEvent("BlinkCast", ConAI)
			ConAI.Delete()
			ConAIStreetPants.Delete()
			ConAIStreetShirt.Delete()
			ConAIStreetSneakers.Delete()
			LoadConAI()
		})
	})
}

// Wizard AI script. //
func LoadWizAI() {
	WizAI = ns.CreateObject("NPC", ns.GetHost())
	WizAIStreetSneakers = ns.CreateObject("StreetSneakers", WizAI)
	WizAIStreetPants = ns.CreateObject("StreetPants", WizAI)
	WizAIStreetShirt = ns.CreateObject("StreetShirt", WizAI)
	WizAIWizardRobe = ns.CreateObject("WizardRobe", WizAI)
	WizAI.Equip(WizAIStreetSneakers)
	WizAI.Equip(WizAIStreetPants)
	WizAI.Equip(WizAIStreetShirt)
	WizAI.Equip(WizAIWizardRobe)
	WizAI.Enchant("ENCHANT_INVULNERABLE", script.Frames(150))
	WizAI.SetMaxHealth(75)
	WizAI.AggressionLevel(0.16)
	ns.NewTimer(ns.Seconds(3), func() {
		WizAI.AggressionLevel(0.83)
	})
	WizAI.Hunt()
	WizAI.ResumeLevel(0.8)
	WizAI.RetreatLevel(0.2)
	// Buff on respawn. //
	if WizAIForceFieldCooldown == 0 {
	} else {
		WizAIForceFieldCooldown = 0
		WizAIGlobalCooldown = 0
		// Force Field chant. //
		ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
		ns.NewTimer(ns.Frames(3), func() {
			ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeSouth", WizAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
					ns.NewTimer(ns.Frames(3), func() {
						ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
						ns.NewTimer(ns.Frames(3), func() {
							ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
							ns.NewTimer(ns.Frames(3), func() {
								ns.AudioEvent("NPCSpellPhonemeSouth", WizAI)
								ns.NewTimer(ns.Frames(3), func() {
									ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
									ns.NewTimer(ns.Frames(3), func() {
										ns.CastSpell(spell.SHIELD, WizAI, WizAI)
										WizAIGlobalCooldown = 1
										// Pause for concentration. //
										ns.NewTimer(ns.Frames(3), func() {
											// Haste chant. //
											ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
											ns.NewTimer(ns.Frames(3), func() {
												ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
												ns.NewTimer(ns.Frames(3), func() {
													ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
													ns.CastSpell(spell.HASTE, WizAI, WizAI)
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
		ns.NewTimer(ns.Seconds(30), func() {
			WizAIForceFieldCooldown = 1
		})
	}
	// When an enemy is seen. //
	WizAI.OnEvent(ns.EventEnemySighted, func() {
		if WizAISlowCooldown == 0 {
		} else {
			WizAISlowCooldown = 0
			WizAIGlobalCooldown = 0
			// Slow chant. //
			ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
					WizAIGlobalCooldown = 1
					WizAITarget = ns.FindClosestObject(WizAI, ns.HasClass(object.ClassPlayer))
					ns.CastSpell(spell.SLOW, WizAI, WizAITarget)
					ns.NewTimer(ns.Seconds(5), func() {
						WizAISlowCooldown = 1
					})
				})
			})
		}
	})
	// On collision. //
	WizAI.OnEvent(ns.EventCollision, func() {
		if WizAIShockCooldown == 0 {
		} else {
			WizAIGlobalCooldown = 0
			WizAIShockCooldown = 0
			// Shock chant. //
			ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
					ns.NewTimer(ns.Frames(3), func() {
						ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
						ns.NewTimer(ns.Frames(3), func() {
							ns.CastSpell(spell.SHOCK, WizAI, WizAI)
							WizAIGlobalCooldown = 1
							ns.NewTimer(ns.Seconds(10), func() {
								WizAIShockCooldown = 1
							})
						})
					})
				})
			})
		}
	})
	// Trap. TODO: define when to, ns.EventLosEnemy is placeholder. IDEA: When no enemy is in sight. //
	WizAI.OnEvent(ns.EventLostEnemy, func() {
		if WizAITrapCooldown == 0 {
		} else {
			WizAIGlobalCooldown = 0
			// WizAITrapCooldown = 0
			// Ring of Fire chant. //
			ns.AudioEvent("NPCSpellPhonemeDownRight", WizAI)
			ns.NewTimer(ns.Frames(3), func() {
				ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
				ns.NewTimer(ns.Frames(3), func() {
					ns.AudioEvent("NPCSpellPhonemeDownLeft", WizAI)
					ns.NewTimer(ns.Frames(3), func() {
						ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
						ns.NewTimer(ns.Frames(3), func() {
							// Pause of Glyph concentration. //
							ns.NewTimer(ns.Frames(3), func() {
								// Magic Missiles chant. //
								ns.NewTimer(ns.Frames(3), func() {
									ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
									ns.NewTimer(ns.Frames(3), func() {
										ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
										ns.NewTimer(ns.Frames(3), func() {
											ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
											ns.NewTimer(ns.Frames(3), func() {
												ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
												ns.NewTimer(ns.Frames(3), func() {
													// Pause of Glyph concentration. //
													ns.NewTimer(ns.Frames(3), func() {
														// Shock chant. //
														ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
														ns.NewTimer(ns.Frames(3), func() {
															ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
															ns.NewTimer(ns.Frames(3), func() {
																ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
																ns.NewTimer(ns.Frames(3), func() {
																	ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
																	ns.NewTimer(ns.Frames(3), func() {
																		// Pause of Glyph concentration. //
																		ns.NewTimer(ns.Frames(3), func() {
																			// Glyph chant. //
																			ns.AudioEvent("NPCSpellPhonemeUp", WizAI)
																			ns.NewTimer(ns.Frames(3), func() {
																				ns.AudioEvent("NPCSpellPhonemeRight", WizAI)
																				ns.NewTimer(ns.Frames(3), func() {
																					ns.AudioEvent("NPCSpellPhonemeLeft", WizAI)
																					ns.NewTimer(ns.Frames(3), func() {
																						ns.AudioEvent("NPCSpellPhonemeDown", WizAI)
																						ns.NewTimer(ns.Frames(3), func() {
																							ns.AudioEvent("TrapDrop", WizAI)
																							WizAITrap = ns.NewTrap(WizAI, spell.CLEANSING_FLAME, spell.MAGIC_MISSILE, spell.SHOCK)
																							WizAITrap.SetOwner(WizAI)
																							WizAIGlobalCooldown = 1
																						})
																					})
																				})
																			})
																		})
																	})
																})
															})
														})
													})
												})
											})
										})
									})
								})
							})
						})
					})
				})
			})
		}
	})
	// On Death. //
	WizAI.OnEvent(ns.EventDeath, func() {
		WizAI.DestroyChat()
		ns.AudioEvent("NPCDie", WizAI)
		// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement? //
		ns.GetHost().ChangeScore(+1)
		ns.NewTimer(ns.Frames(60), func() {
			ns.AudioEvent("BlinkCast", WizAI)
			WizAI.Delete()
			WizAIStreetPants.Delete()
			WizAIStreetSneakers.Delete()
			WizAIStreetShirt.Delete()
			WizAITrapCooldown = 1
			LoadWizAI()
		})
	})
}

// Warrior AI script. //
func LoadWarAI() {
	WarAI = ns.CreateObject("NPC", ns.GetHost())
	// TODO: Change location of item creation OR stop them from respawning automatically. //
	WarAILongsword = ns.CreateObject("Longsword", WarAI)
	WarAIWoodenShield = ns.CreateObject("WoodenShield", WarAI)
	WarAIStreetSneakers = ns.CreateObject("StreetSneakers", WarAI)
	WarAIStreetPants = ns.CreateObject("StreetPants", WarAI)
	WarAI.Equip(WarAILongsword)
	WarAI.Equip(WarAIWoodenShield)
	WarAI.Equip(WarAIStreetSneakers)
	WarAI.Equip(WarAIStreetPants)
	// TODO: Give different audio and chat for each set so they feel like different characters. //
	WarAI.Enchant("ENCHANT_INVULNERABLE", script.Frames(150))
	WarAI.SetMaxHealth(150)
	WarAI.AggressionLevel(0.83)
	WarAI.Hunt()
	WarAI.ResumeLevel(0.8)
	WarAI.RetreatLevel(0.2)
	// WarAI.Chat("War01A.scr:Bully1") // this is a robbery! Your money AND your life! //
	// ns.AudioEvent("F1ROG01E", WarAI)
	// TODO: Add audio to match the chat: F1ROG01E. //
	// Enemy Sighted. //
	WarAI.OnEvent(ns.EventEnemySighted, func() {
		// Script out a plan of action. //
	})
	WarAI.OnEvent(ns.EventRetreat, func() {
		// Walk to nearest RedPotion. //
	})
	// Enemy Lost. //
	WarAI.OnEvent(ns.EventLostEnemy, func() {
		if WarAIEyeOfTheWolfCooldown == 0 {
		} else {
			WarAI.Enchant("ENCHANT_INFRAVISION", ns.Seconds(10))
			WarAIEyeOfTheWolfCooldown = 0
			ns.NewTimer(ns.Seconds(20), func() {
				WarAIEyeOfTheWolfCooldown = 1
			})
		}
	})
	// On Hit. //
	WarAI.OnEvent(ns.EventIsHit, func() {
		// WarAITaggedPlayer = ns.GetCaller() ---> more research needed to select target
	})
	// On Death. //
	WarAI.OnEvent(ns.EventDeath, func() {
		WarAI.DestroyChat()
		ns.AudioEvent("NPCDie", WarAI)
		// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement? //
		ns.GetHost().ChangeScore(+1)
		ns.NewTimer(ns.Frames(60), func() {
			ns.AudioEvent("BlinkCast", WarAI)
			WarAI.Delete()
			WarAIStreetPants.Delete()
			WarAIStreetSneakers.Delete()
			LoadWarAI()
		})
	})
}

func DialogOptions() {
	// Usable dialog bits //
	// F1GD401E "What a Wizard spy?" //

	// C2NC203E "Get away from me filthy peasants" //
	// C2NC202E "URHGHH" //

	// C5OGK02E "too bad you must die now"//
	// C5OGK01E "youre very bold for such a little man" //
	// C5OGK05E "Ill crush your bones" //
}
