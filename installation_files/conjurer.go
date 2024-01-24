package BotWars

import (
	"image/color"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/noxscript/ns/v4/spell"
	"github.com/noxworld-dev/noxscript/ns/v4/subclass"
	"github.com/noxworld-dev/opennox-lib/object"
	"github.com/noxworld-dev/opennox-lib/script"
)

// NewConjurer creates a new Conjurer bot.
func NewConjurer(t *Team) *Conjurer {
	con := &Conjurer{team: t}
	con.init()
	return con
}

// Conjurer bot class.
type Conjurer struct {
	team              *Team
	unit              ns.Obj
	cursor            ns.Pointf
	target            ns.Obj
	mana              int
	startingEquipment struct {
		StreetSneakers ns.Obj
		StreetPants    ns.Obj
		StreetShirt    ns.Obj
	}
	spells struct {
		isAlive              bool
		Ready                bool // Duration unknown.
		CounterspellReady    bool
		BlinkReady           bool // No real cooldown, "cooldown" implemented for balance reasons. TODO: Make random.
		FistOfVengeanceReady bool // No real cooldown, mana cost 60.
		StunReady            bool // No real cooldown.
		PixieCount           int
		ForceOfNatureReady   bool
		InversionReady       bool
		ToxicCloudReady      bool // 60 mana.
		SlowReady            bool
		MeteorReady          bool
		LesserHealReady      bool
		BurnReady            bool
	}
	summons struct {
		SummonCreatureReady bool
		CreatureCage        int
		BomberCount         int
	}
	audio struct {
		ManaRestoreSound bool
	}
	behaviour struct {
		AntiStuck       bool
		Busy            bool
		ManaOfInterest  ns.Obj
		Escorting       bool
		EscortingTarget ns.Obj
		Guarding        bool
		GuardingPos     ns.Pointf
	}
	reactionTime int
}

func (con *Conjurer) init() {
	// Reset spells ConBot.
	con.spells.Ready = true
	// Debuff spells.
	con.spells.SlowReady = true
	con.spells.StunReady = true
	// Offensive spells.
	con.spells.MeteorReady = true
	con.spells.BurnReady = true
	con.spells.FistOfVengeanceReady = true
	con.spells.PixieCount = 0
	con.spells.ForceOfNatureReady = true
	con.spells.ToxicCloudReady = true
	// Defensive spells.
	con.spells.BlinkReady = true
	con.spells.CounterspellReady = true
	con.spells.InversionReady = true
	con.spells.LesserHealReady = true
	// Summons.
	con.summons.SummonCreatureReady = true
	con.summons.CreatureCage = 0
	con.summons.BomberCount = 0
	// Behaviour.
	con.behaviour.AntiStuck = true
	con.behaviour.Busy = false
	con.behaviour.Escorting = false
	con.behaviour.Guarding = false
	// Create ConBot.
	con.unit = ns.CreateObject("NPC", con.team.SpawnPoint())
	con.unit.Enchant(enchant.INVULNERABLE, script.Frames(150))
	con.unit.SetMaxHealth(100)
	con.unit.SetStrength(55)
	con.unit.SetBaseSpeed(88)
	con.spells.isAlive = true
	con.mana = 125
	// Set Team.
	if GameModeIsCTF {
		con.unit.SetOwner(con.team.Spawns()[0])
	}
	con.unit.SetTeam(con.team.Team())
	if con.unit.HasTeam(ns.Teams()[0]) {
		con.unit.SetColor(0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		con.unit.SetColor(1, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		con.unit.SetColor(2, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		con.unit.SetColor(3, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		con.unit.SetColor(4, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
		con.unit.SetColor(5, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	} else {
		con.unit.SetColor(0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		con.unit.SetColor(1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		con.unit.SetColor(2, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		con.unit.SetColor(3, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		con.unit.SetColor(4, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
		con.unit.SetColor(5, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	}
	// Create ConBot mouse cursor.
	con.target = NoTarget
	con.cursor = NoTarget.Pos()
	// Set difficulty (0 = Botlike, 15 = hard, 30 = normal, 45 = easy, 60 = beginner)
	con.reactionTime = BotDifficulty
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
	con.startingEquipment.StreetSneakers = ns.CreateObject("StreetSneakers", ns.Ptf(150, 150))
	con.startingEquipment.StreetPants = ns.CreateObject("StreetPants", ns.Ptf(150, 150))
	con.startingEquipment.StreetShirt = ns.CreateObject("StreetShirt", ns.Ptf(150, 150))
	con.unit.Equip(con.startingEquipment.StreetPants)
	con.unit.Equip(con.startingEquipment.StreetShirt)
	con.unit.Equip(con.startingEquipment.StreetSneakers)
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
	con.unit.OnEvent(ns.EventIsHit, con.onHit)
	// Looking for enemies.
	con.unit.OnEvent(ns.EventLookingForEnemy, con.onLookingForTarget)
	//con.unit.OnEvent(ns.EventChangeFocus, con.onChangeFocus)
	con.unit.OnEvent(ns.EventEndOfWaypoint, con.onEndOfWaypoint)
	con.PassiveManaRegen()
	con.LookForWeapon()
	con.WeaponPreference()
	ns.OnChat(con.onConCommand)
	con.findLoot()
	con.checkPixieCount()
}

func (con *Conjurer) onHit() {
	if con.mana <= 20 && !con.behaviour.Busy {
		con.GoToManaObelisk()
	}
}

func (con *Conjurer) onEndOfWaypoint() {
	con.behaviour.Busy = false
	con.unit.AggressionLevel(0.83)
	if con.mana <= 49 {
		con.GoToManaObelisk()
	} else {
		if GameModeIsCTF {
			con.team.CheckAttackOrDefend(con.unit)
		} else {
			con.unit.WalkTo(con.target.Pos())
			ns.NewTimer(ns.Seconds(2), func() {
				con.unit.Hunt()
			})
		}
	}
	con.LookForNearbyItems()
}

func (con *Conjurer) buffInitial() {
	con.castVampirism()
}

func (con *Conjurer) onLookingForTarget() {
	con.castInfravision()
}

func (con *Conjurer) onEnemyHeard() {
	if !con.unit.CanSee(con.target) {
		con.castForceOfNature()
		con.castInfravision()
	}
}

func (con *Conjurer) onEnemySighted() {
	con.target = ns.GetCaller()
	con.castForceOfNature()
}

func (con *Conjurer) onCollide() {
	if con.spells.isAlive {
		caller := ns.GetCaller()
		if GameModeIsCTF {
			con.team.CheckPickUpEnemyFlag(caller, con.unit)
			con.team.CheckCaptureEnemyFlag(caller, con.unit)
			con.team.CheckRetrievedOwnFlag(caller, con.unit)
		}
		if caller == con.behaviour.ManaOfInterest {
			ns.NewTimer(ns.Seconds(1), func() {
				if con.mana > 110 {
					con.onEndOfWaypoint()
				} else {
					con.GoToManaObelisk()
				}

			})
		}
	}
}

func (con *Conjurer) onRetreat() {
	con.castBlink()
}

func (con *Conjurer) onLostEnemy() {
	con.castInfravision()
	if GameModeIsCTF {
		con.team.WalkToOwnFlag(con.unit)
	}
}

func (con *Conjurer) onDeath() {
	con.spells.isAlive = false
	con.spells.Ready = false
	con.unit.FlagsEnable(object.FlagNoCollide)
	con.team.DropEnemyFlag(con.unit)
	con.unit.DestroyChat()
	ns.AudioEvent(audio.NPCDie, con.unit)
	// TODO: Change ns.GetHost() to correct caller. Is there no Gvar1 replacement?
	// ns.GetHost().ChangeScore(+1)
	if !GameModeIsCTF {
		if con.unit.HasTeam(ns.Teams()[0]) {
			ns.Teams()[1].ChangeScore(+1)
		} else {
			ns.Teams()[0].ChangeScore(+1)
		}
	}
	if !ItemDropEnabled {
		con.startingEquipment.StreetPants.Delete()
		con.startingEquipment.StreetShirt.Delete()
		con.startingEquipment.StreetSneakers.Delete()
	}
	ns.NewTimer(ns.Frames(60), func() {
		ns.AudioEvent(audio.BlinkCast, con.unit)
		con.unit.Delete()
		if ItemDropEnabled {
			con.startingEquipment.StreetPants.Delete()
			con.startingEquipment.StreetShirt.Delete()
			con.startingEquipment.StreetSneakers.Delete()
		}
		if BotRespawn {
			con.init()
		}
	})
}

func (con *Conjurer) PassiveManaRegen() {
	if con.spells.isAlive {
		ns.NewTimer(ns.Seconds(2), func() {
			if con.mana < 125 {
				if !BotMana {
					con.mana = con.mana + 300
				}
				con.mana = con.mana + 1
			}
			con.PassiveManaRegen()
			//ns.PrintStrToAll("con mana: " + strconv.Itoa(con.mana))
		})
	}
}

func (con *Conjurer) UsePotions() {
	if con.unit.CanSee(con.target) {
		if con.unit.CurrentHealth() <= 25 && con.unit.InItems().FindObjects(nil, ns.HasTypeName{"RedPotion"}) != 0 {
			ns.AudioEvent(audio.LesserHealEffect, con.unit)
			RedPotion := con.unit.Items(ns.HasTypeName{"RedPotion"})
			con.unit.SetHealth(con.unit.CurrentHealth() + 50)
			RedPotion[0].Delete()
		}
		if con.mana <= 100 && con.unit.InItems().FindObjects(nil, ns.HasTypeName{"BluePotion"}) != 0 {
			con.mana = con.mana + 50
			ns.AudioEvent(audio.RestoreMana, con.unit)
			BluePotion := con.unit.Items(ns.HasTypeName{"BluePotion"})
			BluePotion[0].Delete()
		}
	}
}

func (con *Conjurer) GoToManaObelisk() {
	if !con.behaviour.Busy {
		con.behaviour.Busy = true
		con.unit.AggressionLevel(0.16)
		NearestObeliskWithMana := ns.FindClosestObjectIn(con.unit, ns.Objects(AllManaObelisksOnMap),
			ns.ObjCondFunc(func(it ns.Obj) bool {
				return it.CurrentMana() >= 10
			}),
		)
		if NearestObeliskWithMana != nil {
			con.behaviour.ManaOfInterest = NearestObeliskWithMana
			if con.unit == con.team.TeamTank {
				if con.unit.CanSee(NearestObeliskWithMana) {
					con.unit.WalkTo(NearestObeliskWithMana.Pos())
				}
			} else {
				con.unit.WalkTo(NearestObeliskWithMana.Pos())
			}
		}
	}
}

func (con *Conjurer) RestoreMana() {
	if con.mana < 125 {
		for i := 0; i < len(AllManaObelisksOnMap); i++ {
			if AllManaObelisksOnMap[i].CurrentMana() > 0 && con.unit.CanSee(AllManaObelisksOnMap[i]) && (ns.InCirclef{Center: con.unit, R: 50}).Matches(AllManaObelisksOnMap[i]) {
				con.mana = con.mana + 1
				AllManaObelisksOnMap[i].SetMana(AllManaObelisksOnMap[i].CurrentMana() - 1)
				con.RestoreManaSound()
			}
		}
	}
}

func (con *Conjurer) RestoreManaSound() {
	if !con.audio.ManaRestoreSound {
		con.audio.ManaRestoreSound = true
		ns.AudioEvent(audio.RestoreMana, con.unit)
		ns.NewTimer(ns.Frames(15), func() {
			con.audio.ManaRestoreSound = false
		})
	}
}

func (con *Conjurer) checkForMissiles() {
	// Maybe need to add a ns.hasteam condition. Not sure yet.
	if sp2 := ns.FindClosestObject(con.unit, ns.HasTypeName{"DeathBall"}, ns.InCirclef{Center: con.unit, R: 500}); sp2 != nil {
		{
			arr2 := ns.FindAllObjects(
				ns.HasTypeName{"NewPlayer", "NPC"},
				ns.HasTeam{con.team.Enemy.Team()},
			)
			for i := 0; i < len(arr2); i++ {
				if sp2.HasOwner(arr2[i]) {
					con.castCounterspellAtForceOfNature()
				}
			}
		}
	} else {
		if sp := ns.FindClosestObject(con.unit, ns.HasClass(object.ClassMissile), ns.InCirclef{Center: con.unit, R: 500}); sp != nil {
			if sp.HasOwner(con.target) {
				con.castInversion()
			}
		}
	}
}

func (con *Conjurer) Update() {
	con.checkForMissiles()
	con.UsePotions()
	con.RestoreMana()
	if con.mana > 125 {
		con.mana = 125
	}
	if con.unit.HasEnchant(enchant.ANTI_MAGIC) {
		con.spells.Ready = true
		// Add in response to warcry.
		// Use weapon.
	}

	if con.spells.Ready {
		con.castPixieSwarm()
		if con.mana >= 100 {
			con.castLesserHeal()
		}
	}
	if con.unit.CanSee(con.target) && con.unit.HasEnchant(enchant.HELD) || con.unit.HasEnchant(enchant.SLOWED) {
		con.castBlink()
	}
	if con.target.HasEnchant(enchant.HELD) || con.target.HasEnchant(enchant.SLOWED) {
		if con.unit.CanSee(con.target) {
			//con.castFistOfVengeance()
			con.castMeteor()
			con.castToxicCloud()
			con.castBurn()
			con.castCounterspell()
		}
	}
	if con.spells.Ready && con.unit.CanSee(con.target) {
		if !GameModeIsCTF {
			con.castStun()
		}
		con.castSlow()

	}
	if !con.unit.CanSee(con.target) && con.spells.Ready {
		con.castVampirism()
		if con.mana >= 85 {
			con.summonRandomCreature()
			con.castProtectionFromShock()
			con.castProtectionFromFire()
			con.castProtectionFromPoison()
		}
	}
	if !con.unit.HasEnchant(enchant.VAMPIRISM) || !con.unit.HasEnchant(enchant.PROTECT_FROM_POISON) || !con.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) || !con.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) || con.summons.CreatureCage <= 3 {
		con.GoToManaObelisk()
	}
}

func (con *Conjurer) summonRandomCreature() {
	random := ns.Random(1, 3)
	if random == 1 {
		if con.summons.CreatureCage == 3 {
			con.castBomber()
		}
		con.summonRandomSmallCreature()
	}
	if random == 2 {
		con.summonRandomMediumCreature()
	}
	if random == 3 {
		con.summonRandomLargeCreature()
	}
}

func (con *Conjurer) summonRandomLargeCreature() {
	random := ns.Random(1, 6)
	if random == 1 {
		con.castSummonMechanicalGolem()
	}
	if random == 2 {
		con.castSummonStoneGolem()
	}
	if random == 3 {
		con.castSummonCarnivorousPlant()
	}
	if random == 4 {
		con.castSummonWillOWisp()
	}
	if random == 5 {
		con.castSummonMimic()
	}
	if random == 6 {
		con.castSummonBeholder()
	}
}

func (con *Conjurer) summonRandomMediumCreature() {
	random := ns.Random(1, 20)
	if random == 1 {
		con.castSummonBlackBear()
	}
	if random == 2 {
		con.castSummonOgre()
	}
	if random == 3 {
		con.castSummonBlackWolf()
	}
	if random == 4 {
		con.castSummonWhiteWolf()
	}
	if random == 5 {
		con.castSummonWolf()
	}
	if random == 6 {
		con.castSummonGargoyle()
	}
	if random == 7 {
		con.castSummonOgress()
	}
	if random == 8 {
		con.castSummonOgreLord()
	}
	if random == 9 {
		con.castSummonZombie()
	}
	if random == 10 {
		con.castSummonVileZombie()
	}
	if random == 11 {
		con.castSummonEmberDemon()
	}
	if random == 12 {
		con.castSummonShade()
	}
	if random == 13 {
		con.castSummonLargeCaveSpider()
	}
	if random == 14 {
		con.castSummonGrizzlyBear()
	}
	if random == 15 {
		con.castSummonSkeleton()
	}
	if random == 16 {
		con.castSummonSkeletonLord()
	}
	if random == 17 {
		con.castSummonScorpion()
	}
	if random == 18 {
		con.castSummonSpider()
	}
	if random == 19 {
		con.castSummonSpittingSpider()
	}
	if random == 20 {
		con.castSummonTroll()
	}
}

func (con *Conjurer) summonRandomSmallCreature() {
	random := ns.Random(1, 10)
	if random == 1 {
		con.castSummonWasp()
	}
	if random == 2 {
		con.castSummonUrchin()
	}
	if random == 3 {
		con.castSummonSmallSpider()
	}
	if random == 4 {
		con.castSummonSmallCaveSpider()
	}
	if random == 5 {
		con.castSummonMechanicalFlyer()
	}
	if random == 6 {
		con.castSummonImp()
	}
	if random == 7 {
		con.castSummonGiantLeech()
	}
	if random == 8 {
		con.castSummonBat()
	}
	if random == 9 {
		con.castSummonGhost()
	}
	if random == 10 {
		con.castBomber()
	}
}

func (con *Conjurer) LookForWeapon() {
	if !con.behaviour.Busy {
		con.behaviour.Busy = true
		ItemLocation := ns.FindClosestObject(con.unit, ns.HasTypeName{"CrossBow", "InfinitePainWand"})
		if ItemLocation != nil {
			con.unit.WalkTo(ItemLocation.Pos())
		}
	}
}

func (con *Conjurer) LookForNearbyItems() {
	if !con.behaviour.Busy {
		if ns.FindAllObjects(ns.HasTypeName{"CrossBow", "InfinitePainWand", "InfinitePainWand", "LesserFireballWand", "Quiver",
			"LeatherArmoredBoots", "LeatherArmor",
			"LeatherHelm",
			"LeatherLeggings", "LeatherArmbands",
			"RedPotion",
			"ConjurerHelm",
			"CurePoisonPotion",
			"BluePotion",
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
			ns.InCirclef{Center: con.unit, R: 200}) != nil {
			if con.unit.InItems().FindObjects(nil, ns.HasTypeName{"CrossBow", "InfinitePainWand", "InfinitePainWand", "LesserFireballWand", "Quiver",
				"LeatherArmoredBoots", "LeatherArmor",
				"LeatherHelm",
				"LeatherLeggings", "LeatherArmbands",
				"RedPotion",
				"ConjurerHelm",
				"CurePoisonPotion",
				"BluePotion",
				"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"}) == 0 {
				ItemLocation := ns.FindAllObjects(ns.HasTypeName{"CrossBow", "InfinitePainWand", "InfinitePainWand", "LesserFireballWand", "Quiver",
					"LeatherArmoredBoots", "LeatherArmor",
					"LeatherHelm",
					"LeatherLeggings", "LeatherArmbands",
					"RedPotion",
					"ConjurerHelm",
					"CurePoisonPotion",
					"BluePotion",
					"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants"},
					ns.InCirclef{Center: con.unit, R: 200},
				)
				if con.unit.CanSee(ItemLocation[0]) {
					con.unit.WalkTo(ItemLocation[0].Pos())
				}
			}
		}
	}
	ns.NewTimer(ns.Seconds(5), func() {
		// prevent bots getting stuck to stay in loop.
		if con.behaviour.AntiStuck {
			con.behaviour.AntiStuck = false
			if GameModeIsCTF {
				con.team.CheckAttackOrDefend(con.unit)
			} else {
				con.behaviour.Busy = false
				con.unit.Hunt()
				con.unit.AggressionLevel(0.83)
			}
			ns.NewTimer(ns.Seconds(6), func() {
				con.behaviour.AntiStuck = true
			})
		}
	})
}

func (con *Conjurer) WeaponPreference() {
	// Priority list to get the prefered weapon.
	// TODO: Add stun and range conditions.
	if con.unit.InItems().FindObjects(nil, ns.HasTypeName{"CrossBow"}) != 0 && con.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"CrossBow"}) == 0 {
		con.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				con.unit.Equip(it)
				//war.unit.Chat("I swapped to my GreatSword!")
				return true
			},
			ns.HasTypeName{"FireStormWand"},
		)
	} else if con.unit.InItems().FindObjects(nil, ns.HasTypeName{"InfinitePainWand"}) != 0 && con.unit.InEquipment().FindObjects(nil, ns.HasTypeName{"InfinitePainWand"}) == 0 {
		con.unit.InItems().FindObjects(
			func(it ns.Obj) bool {
				con.unit.Equip(it)
				//war.unit.Chat("I swapped to my WarHammer!")
				return true
			},
			ns.HasTypeName{"ForceWand"},
		)
	}
	ns.NewTimer(ns.Seconds(10), func() {
		con.WeaponPreference()
	})
}

func (con *Conjurer) findLoot() {
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
			con.unit.Pickup(item)
			con.unit.Equip(con.unit.GetLastItem())
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
			// BlueConjurer Helm.
			//"ConjurerHelm",

			// Leather armor.
			"LeatherArmoredBoots", "LeatherArmor",
			//"LeatherHelm",
			"LeatherLeggings", "LeatherArmbands",

			// Cloth armor.
			"LeatherBoots", "MedievalCloak", "MedievalShirt", "MedievalPants",
		},
	)
	for _, item := range armor {
		if con.unit.CanSee(item) {
			con.unit.Pickup(item)
			con.unit.Equip(con.unit.GetLastItem())
		}
	}
	// Potions.
	potions := ns.FindAllObjects(
		ns.InCirclef{Center: con.unit, R: dist},
		ns.HasTypeName{
			"RedPotion",
			"CurePoisonPotion",
			"BluePotion",
		},
	)
	for _, item := range potions {
		if con.unit.CanSee(item) {
			con.unit.Pickup(item)
		}
	}
	ns.NewTimer(ns.Frames(15), func() {
		con.findLoot()
	})
}

// Checks the ammount of summons active for the Conjurer bot.
func (con *Conjurer) checkCreatureCage() {
	// Get all active sommons that belong to the Conjuer bot.
	if con.summons.CreatureCage == 4 {
	} else {
		allActiveSummons := ns.FindAllObjects(ns.HasClass(object.ClassMonster), ns.ObjCondFunc(func(it ns.Obj) bool {
			if it != con.unit || !it.HasTeam(con.unit.Team()) {
				return it.HasOwner(con.unit)
			} else {
				return false
			}
		}))
		con.summons.CreatureCage = 0
		for _, summon := range allActiveSummons {
			if summon.HasSubclass(subclass.SMALL_MONSTER) {
				// Add the summon to the Creature Cage.
				con.summons.CreatureCage = con.summons.CreatureCage + 1
				// Track if the summon is a bomber.
				if summon.HasSubclass(subclass.BOMBER) {
					con.summons.BomberCount = con.summons.BomberCount + 1
				}
				summon.OnEvent(ns.EventDeath, func() {
					// Remove summon from the Creature Cage on death.
					con.summons.CreatureCage = con.summons.CreatureCage - 1
					// Track if the summon is a bomber.
					if summon.HasSubclass(subclass.BOMBER) {
						con.summons.BomberCount = con.summons.BomberCount + -1
					}
				})
			}
			if summon.HasSubclass(subclass.MEDIUM_MONSTER) {
				con.summons.CreatureCage = con.summons.CreatureCage + 2
				summon.OnEvent(ns.EventDeath, func() {
					con.summons.CreatureCage = con.summons.CreatureCage - 2
				})
			}
			if summon.HasSubclass(subclass.LARGE_MONSTER) {
				con.summons.CreatureCage = con.summons.CreatureCage + 4
				summon.OnEvent(ns.EventDeath, func() {
					con.summons.CreatureCage = con.summons.CreatureCage - 4
				})
			}
		}
	}
	ns.NewTimer(ns.Seconds(1), func() {
		con.checkCreatureCage()
	})
}

// Checks the ammount of Pixies active for the Conjurer bot.
func (con *Conjurer) checkPixieCount() {
	allPixies := ns.FindAllObjects(ns.HasTypeName{"Pixie"})
	if allPixies == nil {
		con.spells.PixieCount = 0
	} else if allPixies[0].HasOwner(con.unit) {
		con.spells.PixieCount = 1
	}
	ns.NewTimer(ns.Seconds(1), func() {
		con.checkPixieCount()
	})
}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //

func (con *Conjurer) castLesserHeal() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CurrentHealth() <= 60 && con.spells.Ready && con.spells.LesserHealReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDownRight, PhUp, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 30
						ns.CastSpell(spell.LESSER_HEAL, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castVampirism() {
	// Check if cooldowns are ready.
	if con.mana >= 20 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && !con.unit.HasEnchant(enchant.VAMPIRISM) {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 20
						ns.CastSpell(spell.VAMPIRISM, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castProtectionFromPoison() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && !con.unit.HasEnchant(enchant.PROTECT_FROM_POISON) {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownLeft, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 30
						ns.CastSpell(spell.PROTECTION_FROM_POISON, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castPixieSwarm() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.PixieCount == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhDown, PhRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 30
						ns.CastSpell(spell.PIXIE_SWARM, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

//func (con *Conjurer) castFistOfVengeance() {
//	// Check if cooldowns are ready.
//	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.FistOfVengeanceReady && con.spells.Ready {
//		// Select target.
//		con.cursor = con.target.Pos()
//		// Trigger cooldown.
//		con.spells.Ready = false
//		// Check reaction time based on difficulty setting.
//		ns.NewTimer(ns.Frames(con.reactionTime), func() {
//			// Check for War Cry before chant.
//			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
//				castPhonemes(con.unit, []audio.Name{PhUpRight, PhUp, PhDown}, func() {
//					// Check for War Cry before spell release.
//					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
//						// Aim.
//						con.unit.LookAtObject(con.target)
//						con.unit.Pause(ns.Frames(con.reactionTime))
//						con.spells.FistOfVengeanceReady = false
//						ns.CastSpell(spell.FIST, con.unit, con.cursor)
//						con.mana = con.mana - 60
//						// Global cooldown.
//						ns.NewTimer(ns.Frames(3), func() {
//							con.spells.Ready = true
//						})
//						ns.NewTimer(ns.Seconds(5), func() {
//							// Fist Of Vengeance cooldown.
//							con.spells.FistOfVengeanceReady = true
//						})
//					}
//				})
//			} else {
//				ns.NewTimer(ns.Frames(con.reactionTime), func() {
//					con.spells.Ready = true
//				})
//			}
//		})
//	}
//}

func (con *Conjurer) castForceOfNature() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.ForceOfNatureReady && con.spells.Ready {
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
						con.unit.Pause(ns.Frames(36))
						con.mana = con.mana - 60
						ns.CastSpell(spell.FORCE_OF_NATURE, con.unit, con.target.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Force of Nature cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							con.spells.ForceOfNatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castProtectionFromFire() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && !con.unit.HasEnchant(enchant.PROTECT_FROM_FIRE) && con.spells.Ready {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhLeft, PhRight, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 30
						ns.CastSpell(spell.PROTECTION_FROM_FIRE, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castProtectionFromShock() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && !con.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) && con.spells.Ready && !con.unit.HasEnchant(enchant.PROTECT_FROM_ELECTRICITY) {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhDownRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.mana = con.mana - 30
						ns.CastSpell(spell.PROTECTION_FROM_ELECTRICITY, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castInversion() {
	// Check if cooldowns are ready.
	if con.mana >= 10 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.InversionReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, FPhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.InversionReady = false
						con.mana = con.mana - 10
						ns.CastSpell(spell.INVERSION, con.unit, con.unit)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Inversion cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							con.spells.InversionReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castBlink() {
	// Check if cooldowns are ready.
	if con.mana >= 10 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.BlinkReady && con.unit != con.team.TeamTank {
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
						con.mana = con.mana - 10
						ns.NewTrap(con.unit, spell.BLINK)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Blink cooldown.
						ns.NewTimer(ns.Seconds(1), func() {
							con.spells.BlinkReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castBurn() {
	// Check if cooldowns are ready.
	if con.mana >= 10 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && !con.target.HasEnchant(enchant.INVULNERABLE) && con.spells.BurnReady && con.spells.Ready && con.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && !con.target.HasEnchant(enchant.INVULNERABLE) {
		// Select target.
		con.cursor = con.target.Pos()
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDown, PhDown, PhUp, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) {
						// Aim.
						con.unit.LookAtObject(con.target)
						con.unit.Pause(ns.Frames(con.reactionTime))
						con.spells.BurnReady = false
						ns.CastSpell(spell.BURN, con.unit, con.cursor)
						con.mana = con.mana - 10
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Burn cooldown.
						ns.NewTimer(ns.Frames(1), func() {
							con.spells.BurnReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castStun() {
	// Check if cooldowns are ready.
	if con.mana >= 10 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.StunReady && con.spells.Ready && !con.target.HasEnchant(enchant.HELD) && !con.target.HasEnchant(enchant.SLOWED) && con.target.MaxHealth() != 150 {
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
						con.mana = con.mana - 10
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						ns.NewTimer(ns.Seconds(5), func() {
							// Stun cooldown.
							con.spells.StunReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castToxicCloud() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.ToxicCloudReady && con.spells.Ready {
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
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Toxic Cloud cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							con.spells.ToxicCloudReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSlow() {
	// Check if cooldowns are ready.
	if con.mana >= 10 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.SlowReady && con.spells.Ready && !con.target.HasEnchant(enchant.SLOWED) && !con.target.HasEnchant(enchant.REFLECTIVE_SHIELD) && !con.target.HasEnchant(enchant.HELD) {
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
						con.mana = con.mana - 10
						ns.CastSpell(spell.SLOW, con.unit, con.target)
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Slow cooldown.
						ns.NewTimer(ns.Seconds(3), func() {
							con.spells.SlowReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castCounterspell() {
	// Check if cooldowns are ready.
	if con.mana >= 20 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && !con.unit.HasEnchant(enchant.INVISIBLE) && con.target.HasEnchant(enchant.SHOCK) && con.spells.Ready && con.spells.CounterspellReady && con.unit.CanSee(con.target) {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDown, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) {
						con.spells.CounterspellReady = false
						con.mana = con.mana - 20
						ns.CastSpell(spell.COUNTERSPELL, con.unit, con.unit.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(5), func() {
							con.spells.CounterspellReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castMeteor() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.unit.CanSee(con.target) && con.spells.MeteorReady && con.spells.Ready {
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
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						ns.NewTimer(ns.Seconds(5), func() {
							// Meteor cooldown.
							con.spells.MeteorReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castInfravision() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && !con.unit.HasEnchant(enchant.INFRAVISION) {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhRight, PhLeft, PhRight, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						ns.CastSpell(spell.INFRAVISION, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castBomber() {
	// Check if cooldowns are ready.
	if con.mana >= 80 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.BomberCount <= 1 && con.summons.CreatureCage <= 3 && con.summons.SummonCreatureReady {
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
							// Burn chant.
							castPhonemes(con.unit, []audio.Name{PhDown, PhDown, PhUp, PhUp}, func() {
								// Pause for concentration.
								ns.NewTimer(ns.Frames(3), func() {
									// Check for War Cry before chant.
									if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
										// Toxic Cloud chant.
										castPhonemes(con.unit, []audio.Name{PhUpRight, PhDownLeft, PhUpLeft}, func() {
											// Pause for concentration.
											ns.NewTimer(ns.Frames(3), func() {
												// Check for War Cry before chant.
												if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
													// Glyph chant.
													castPhonemes(con.unit, []audio.Name{PhUp, PhRight, PhLeft, PhDown}, func() {
														// Check for War Cry before spell release.
														if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
															con.summons.SummonCreatureReady = false
															bomber := ns.CreateObject("Bomber", con.unit)
															ns.AudioEvent("BomberSummon", bomber)
															bomber.SetOwner(con.unit)
															bomber.SetTeam(con.team.Team())
															bomber.Follow(con.unit)
															bomber.MonsterStatusEnable(object.MonStatusAlert)
															bomber.TrapSpells(spell.BURN, spell.TOXIC_CLOUD, spell.STUN)
															bomber.OnEvent(ns.ObjectEvent(ns.EventEnemySighted), func() {
																bomber.Attack(con.target)
															})
															bomber.OnEvent(ns.ObjectEvent(ns.EventEnemyHeard), func() {
																bomber.Attack(con.target)
															})
															bomber.OnEvent(ns.ObjectEvent(ns.EventLostEnemy), func() {
																bomber.Follow(con.unit)
															})
															// Global cooldown.
															ns.NewTimer(ns.Frames(3), func() {
																con.checkCreatureCage()
																con.spells.Ready = true
																con.summons.SummonCreatureReady = true
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
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castCounterspellAtForceOfNature() {
	// Check if cooldowns are ready.
	if con.mana >= 20 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.spells.CounterspellReady {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhDown, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.spells.CounterspellReady = false
						con.mana = con.mana - 20
						ns.CastSpell(spell.COUNTERSPELL, con.unit, con.unit.Pos())
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Haste cooldown.
						ns.NewTimer(ns.Seconds(20), func() {
							con.spells.CounterspellReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

// ------------------------------------------------------------------------------------------------------------------------------------ //
// ---------------------------------------------------------------- SPELL BOOK -------------------------------------------------------- //
// ------------------------------------------------------------------------------------------------------------------------------------ //

func (con *Conjurer) onConCommand(t ns.Team, p ns.Player, obj ns.Obj, msg string) string {
	if p != nil {
		switch msg {
		// Bot commands.
		case "help", "Help", "Follow", "follow", "escort", "Escort", "come", "Come":
			if con.unit.CanSee(p.Unit()) && con.unit.Team() == p.Team() {
				con.behaviour.Escorting = true
				con.behaviour.EscortingTarget = p.Unit()
				con.behaviour.Guarding = false
				con.unit.Follow(p.Unit())
				random := ns.Random(1, 4)
				if random == 1 {
					con.unit.ChatStr("I'll follow you.")
				}
				if random == 2 {
					con.unit.ChatStr("Let's go.")
				}
				if random == 3 {
					con.unit.ChatStr("I'll help.")
				}
				if random == 4 {
					con.unit.ChatStr("Sure thing.")
				}
			}
		case "Attack", "Go", "go", "attack":
			if con.unit.CanSee(p.Unit()) && con.unit.Team() == p.Team() {
				con.behaviour.Escorting = false
				con.behaviour.Guarding = false
				con.unit.Hunt()
				random2 := ns.Random(1, 4)
				if random2 == 1 {
					con.unit.ChatStr("I'll get them.")
				}
				if random2 == 2 {
					con.unit.ChatStr("Time to shine.")
				}
				if random2 == 3 {
					con.unit.ChatStr("On the offense.")
				}
				if random2 == 4 {
					con.unit.ChatStr("Time to hunt.")
				}
			}
		case "guard", "stay", "Guard", "Stay":
			if con.unit.CanSee(p.Unit()) && con.unit.Team() == p.Team() {
				con.unit.Guard(con.unit.Pos(), con.unit.Pos(), 300)
				con.behaviour.Escorting = false
				con.behaviour.Guarding = true
				con.behaviour.GuardingPos = con.unit.Pos()
				random1 := ns.Random(1, 4)
				if random1 == 1 {
					con.unit.ChatStr("I'll guard this place.")
				}
				if random1 == 2 {
					con.unit.ChatStr("No problem.")
				}
				if random1 == 3 {
					con.unit.ChatStr("I'll stay.")
				}
				if random1 == 4 {
					con.unit.ChatStr("I'll hold.")
				}
			}
		case "vamp", "Vamp", "Vampirism", "vampirism":
			if con.unit.CanSee(p.Unit()) && con.unit.HasTeam(p.Unit().Team()) {
				if con.mana >= 20 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready {
					// Trigger cooldown.
					con.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(con.reactionTime), func() {
						// Check for War Cry before chant.
						if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
								// Check for War Cry before spell release.
								if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
									con.mana = con.mana - 20
									ns.CastSpell(spell.VAMPIRISM, con.unit, p.Unit())
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										con.spells.Ready = true
									})
								}
							})
						} else {
							ns.NewTimer(ns.Frames(con.reactionTime), func() {
								con.spells.Ready = true
							})
						}
					})
				}
				if con.mana < 20 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready {
					// Trigger cooldown.
					con.spells.Ready = false
					// Check reaction time based on difficulty setting.
					ns.NewTimer(ns.Frames(con.reactionTime), func() {
						// Check for War Cry before chant.
						if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
							castPhonemes(con.unit, []audio.Name{PhUp, PhDown, PhLeft, PhRight}, func() {
								// Check for War Cry before spell release.
								if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
									ns.AudioEvent(audio.ManaEmpty, con.unit)
									// Global cooldown.
									ns.NewTimer(ns.Frames(3), func() {
										con.spells.Ready = true
										con.unit.ChatStr("Not enough mana.")
									})
								}
							})
						} else {
							ns.NewTimer(ns.Frames(con.reactionTime), func() {
								con.spells.Ready = true
							})
						}
					})
				}
			}
		}
	}
	return msg
}

// ---------------- summon creatures ----------------- //

func (con *Conjurer) castSummonGhost() {
	// Check if cooldowns are ready.
	if con.mana >= 15 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_GHOST, con.unit, con.unit)
						con.mana = con.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonBat() {
	// Check if cooldowns are ready.
	if con.mana >= 15 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_BAT, con.unit, con.unit)
						con.mana = con.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonUrchin() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_URCHIN, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonGiantLeech() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_GIANT_LEECH, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonImp() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_IMP, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonMechanicalFlyer() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_MECHANICAL_FLYER, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonWasp() {
	// Check if cooldowns are ready.
	if con.mana >= 15 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_WASP, con.unit, con.unit)
						con.mana = con.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSmallSpider() {
	// Check if cooldowns are ready.
	if con.mana >= 15 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SMALL_SPIDER, con.unit, con.unit)
						con.mana = con.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSmallCaveSpider() {
	// Check if cooldowns are ready.
	if con.mana >= 15 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 3 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SMALL_ALBINO_SPIDER, con.unit, con.unit)
						con.mana = con.mana - 15
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(2), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonVileZombie() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_VILE_ZOMBIE, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonZombie() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_ZOMBIE, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonWhiteWolf() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_WHITE_WOLF, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonWolf() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_WOLF, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonTroll() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhUpLeft, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_TROLL, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSpittingSpider() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUpLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SPITTING_SPIDER, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSpider() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SPIDER, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSkeletonLord() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SKELETON_LORD, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonSkeleton() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SKELETON, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonShade() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SHADE, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonScorpion() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_SCORPION, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonOgress() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhUpLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_OGRE, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonOgreLord() {
	// Check if cooldowns are ready.
	if con.mana >= 85 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhUpLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_OGRE_WARLORD, con.unit, con.unit)
						con.mana = con.mana - 85
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonOgre() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_OGRE_BRUTE, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonGrizzlyBear() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_BEAR, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonGargoyle() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownLeft, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_EVIL_CHERUB, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonLargeCaveSpider() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhUpRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_ALBINO_SPIDER, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonBlackWolf() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhLeft, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_BLACK_WOLF, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonBlackBear() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownRight, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_BLACK_BEAR, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonEmberDemon() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage <= 2 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhUp}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_EMBER_DEMON, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(7), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonWillOWisp() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_WILLOWISP, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonMimic() {
	// Check if cooldowns are ready.
	if con.mana >= 85 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_MIMIC, con.unit, con.unit)
						con.mana = con.mana - 85
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonMechanicalGolem() {
	// Check if cooldowns are ready.
	if con.mana >= 85 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_MECHANICAL_GOLEM, con.unit, con.unit)
						con.mana = con.mana - 85
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonCarnivorousPlant() {
	// Check if cooldowns are ready.
	if con.mana >= 30 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhUp, PhRight}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_CARNIVOROUS_PLANT, con.unit, con.unit)
						con.mana = con.mana - 30
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonStoneGolem() {
	// Check if cooldowns are ready.
	if con.mana >= 85 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhDown, PhDownLeft}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_STONE_GOLEM, con.unit, con.unit)
						con.mana = con.mana - 85
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}

func (con *Conjurer) castSummonBeholder() {
	// Check if cooldowns are ready.
	if con.mana >= 60 && con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) && con.spells.Ready && con.summons.SummonCreatureReady && con.summons.CreatureCage == 0 {
		// Trigger cooldown.
		con.spells.Ready = false
		// Check reaction time based on difficulty setting.
		ns.NewTimer(ns.Frames(con.reactionTime), func() {
			// Check for War Cry before chant.
			if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
				castPhonemes(con.unit, []audio.Name{PhUpLeft, PhDownRight, PhUpRight, PhDownLeft, PhRight, PhDownRight, PhDown}, func() {
					// Check for War Cry before spell release.
					if con.spells.isAlive && !con.unit.HasEnchant(enchant.ANTI_MAGIC) {
						con.summons.SummonCreatureReady = false
						ns.CastSpell(spell.SUMMON_BEHOLDER, con.unit, con.unit)
						con.mana = con.mana - 60
						// Global cooldown.
						ns.NewTimer(ns.Frames(3), func() {
							con.spells.Ready = true
						})
						// Summon Ghost cooldown.
						ns.NewTimer(ns.Seconds(13), func() {
							con.checkCreatureCage()
							con.summons.SummonCreatureReady = true
						})
					}
				})
			} else {
				ns.NewTimer(ns.Frames(con.reactionTime), func() {
					con.spells.Ready = true
				})
			}
		})
	}
}
