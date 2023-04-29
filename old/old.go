package basicmap_old

import (
	ns3 "github.com/noxworld-dev/noxscript/ns/v3"
)

func MapInitialize() {
	WarBot = ns3.Object("WarBot")
	WizBot = ns3.Object("WizBot")
	WarBotCorpse = ns3.Object("WarBotCorpse")
	WizBotCorpse = ns3.Object("WizBotCorpse")
	OutOfGameWar = ns3.Waypoint("OutOfGameWar")
	OutOfGameWiz = ns3.Waypoint("OutOfGameWiz")
	WarSound = ns3.Waypoint("WarSound")
	WizSound = ns3.Waypoint("WizSound")
	WarCryCooldown = 0
	EyeOfTheWolfCooldown = 0
	BerserkerChargeCooldown = 0
	GlobalCooldown = 0
	RespawnCooldownDelay = 0
	ns3.SecondTimer(2, Respawn)
	ns3.SecondTimer(2, RespawnWiz)
}
func CheckCurrentHealth() {
	if ns3.CurrentHealth(WarBot) < 1001 {
		ns3.ChangeScore(ns3.GetCaller(), 1)
		DeadBotWar()
	}
	if ns3.CurrentHealth(WizBot) < 1001 {
		ns3.ChangeScore(ns3.GetCaller(), 1)
		DeadBotWiz()
	}
}
func DeadBotWar() {
	ns3.MoveWaypoint(WarSound, ns3.GetObjectX(WarBot), ns3.GetObjectY(WarBot))
	ns3.AudioEvent("NPCDie", WarSound)
	ns3.MoveObject(WarBotCorpse, ns3.GetObjectX(WarBot), ns3.GetObjectY(WarBot))
	ns3.MoveObject(WarBot, ns3.GetWaypointX(OutOfGameWar), ns3.GetWaypointY(OutOfGameWar))
	ns3.SecondTimer(2, Respawn)
}
func DeadBotWiz() {
	ns3.MoveWaypoint(WizSound, ns3.GetObjectX(WizBot), ns3.GetObjectY(WizBot))
	ns3.AudioEvent("NPCDie", WizSound)
	ns3.MoveObject(WizBotCorpse, ns3.GetObjectX(WizBot), ns3.GetObjectY(WizBot))
	ns3.MoveObject(WizBot, ns3.GetWaypointX(OutOfGameWiz), ns3.GetWaypointY(OutOfGameWiz))
	ns3.SecondTimer(2, RespawnWiz)
}
func Respawn() {
	ns3.RestoreHealth(WarBot, ns3.MaxHealth(WarBot)-ns3.CurrentHealth(WarBot))
	WarCryCooldown = 0
	EyeOfTheWolfCooldown = 0
	BerserkerChargeCooldown = 0
	ns3.MoveWaypoint(WarSound, ns3.GetObjectX(WarBotCorpse), ns3.GetObjectY(WarBotCorpse))
	ns3.AudioEvent("BlinkCast", WarSound)
	BotSpawn = ns3.Object("BotSpawn" + ns3.IntToString(ns3.Random(1, 14)))
	ns3.MoveObject(WarBot, ns3.GetObjectX(BotSpawn), ns3.GetObjectY(BotSpawn))
	ns3.MoveObject(WarBotCorpse, ns3.GetWaypointX(OutOfGameWar), ns3.GetWaypointY(OutOfGameWar))
	ns3.Enchant(WarBot, "ENCHANT_INVULNERABLE", 5.0)
	RespawnCooldownDelay = 1
	ns3.SecondTimer(10, RespawnCooldownDelayReset)
}
func RespawnWiz() {
	ns3.RestoreHealth(WizBot, ns3.MaxHealth(WizBot)-ns3.CurrentHealth(WizBot))
	ns3.MoveWaypoint(WizSound, ns3.GetObjectX(WizBotCorpse), ns3.GetObjectY(WizBotCorpse))
	ns3.AudioEvent("BlinkCast", WizSound)
	BotSpawn = ns3.Object("BotSpawn" + ns3.IntToString(ns3.Random(1, 14)))
	ns3.MoveObject(WizBot, ns3.GetObjectX(BotSpawn), ns3.GetObjectY(BotSpawn))
	ns3.MoveObject(WizBotCorpse, ns3.GetWaypointX(OutOfGameWiz), ns3.GetWaypointY(OutOfGameWiz))
	ns3.Enchant(WizBot, "ENCHANT_INVULNERABLE", 5.0)
}
func WarCry() {
	self, other := ns3.GetTrigger(), ns3.GetCaller()
	if ns3.MaxHealth(other) != 150 {
		if (WarCryCooldown == 0) && (GlobalCooldown == 0) {
			ns3.MoveWaypoint(WarSound, ns3.GetObjectX(WarBot), ns3.GetObjectY(WarBot))
			ns3.AudioEvent("WarcryInvoke", WarSound)
			ns3.PauseObject(WarBot, 45)
			ns3.Enchant(WarBot, "ENCHANT_HELD", 1.0)
			ns3.Enchant(other, "ENCHANT_ANTI_MAGIC", 3.0)
			ns3.EnchantOff(self, "ENCHANT_SHOCK")
			ns3.EnchantOff(self, "ENCHANT_INVULNERABLE")
			WarCryCooldown = 1
			GlobalCooldown = 1
			ns3.SecondTimer(10, WarCryCooldownReset)
			ns3.SecondTimer(1, GlobalCooldownReset)
		}
	}
}
func WarCryCooldownReset() {
	if RespawnCooldownDelay == 0 {
		WarCryCooldown = 0
	}
}
func EyeOfTheWolf() {
	self := ns3.GetTrigger()
	ns3.Wander(WarBot)
	if EyeOfTheWolfCooldown == 0 {
		ns3.Enchant(self, "ENCHANT_INFRAVISION", 10.0)
		EyeOfTheWolfCooldown = 1
		ns3.SecondTimer(20, EyeOfTheWolfCooldownReset)
	}
}
func EyeOfTheWolfCooldownReset() {
	if RespawnCooldownDelay == 0 {
		EyeOfTheWolfCooldown = 0
	}
}

// ------------------------------------- BERSERKERCHARGE SCRIPT ---------------------------------------------- //

func UnitRatioX(unit, target ns3.ObjectID, size float32) float32 {
	return (ns3.GetObjectX(unit) - ns3.GetObjectX(target)) * size / ns3.Distance(ns3.GetObjectX(unit), ns3.GetObjectY(unit), ns3.GetObjectX(target), ns3.GetObjectY(target))
}

func UnitRatioY(unit, target ns3.ObjectID, size float32) float32 {
	return (ns3.GetObjectY(unit) - ns3.GetObjectY(target)) * size / ns3.Distance(ns3.GetObjectX(unit), ns3.GetObjectY(unit), ns3.GetObjectX(target), ns3.GetObjectY(target))
}
func WarBotDetectEnemy() {
	if (BerserkerChargeCooldown == 0) && (GlobalCooldown == 0) {
		rnd := ns3.Random(0, 2)

		if (rnd == 0) || (rnd == 1) {
			BerserkerChargeCooldown = 1
			GlobalCooldown = 1
			ns3.SecondTimer(1, GlobalCooldownReset)
			BerserkerInRange(ns3.GetTrigger(), ns3.GetCaller(), 10)
		}
	} else {
		if WarCryCooldown == 0 {
			WarCry()
		}
	}
}
func CheckUnitFrontSight(unit ns3.ObjectID, dtX, dtY float32) bool {
	ns3.MoveWaypoint(1, ns3.GetObjectX(unit)+dtX, ns3.GetObjectY(unit)+dtY)
	temp := ns3.CreateObject("InvisibleLightBlueHigh", 1)
	res := ns3.IsVisibleTo(unit, temp)
	ns3.Delete(temp)
	return res
}
func pointedByPlr() {
}
func BerserkerInRange(owner, target ns3.ObjectID, wait int) {
	if ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 {
		if !ns3.HasEnchant(owner, "ENCHANT_ETHEREAL") {
			ns3.Enchant(owner, "ENCHANT_ETHEREAL", 0.0)
			ns3.MoveWaypoint(1, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
			unit := ns3.CreateObject("InvisibleLightBlueHigh", 1)
			ns3.MoveWaypoint(1, ns3.GetObjectX(unit), ns3.GetObjectY(unit))
			unit1 := ns3.CreateObject("InvisibleLightBlueHigh", 1)
			ns3.LookWithAngle(unit, wait)
			ns3.FrameTimer(1, func() {
				BerserkerWaitStrike(unit, unit1, owner, target, wait)
			})
		}
	}
}

func BerserkerWaitStrike(ptr, ptr1, owner, target ns3.ObjectID, count int) {
	for {
		if ns3.IsObjectOn(ptr) && ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 && ns3.IsObjectOn(owner) {
			if count != 0 {
				if ns3.IsVisibleTo(owner, target) && ns3.Distance(ns3.GetObjectX(owner), ns3.GetObjectY(owner), ns3.GetObjectX(target), ns3.GetObjectY(target)) < 400.0 {
					BerserkerCharge(owner, target)
				} else {
					ns3.FrameTimer(6, func() {
						BerserkerWaitStrike(ptr, ptr1, owner, target, count-1)
					})
					break
				}
			}
		}
		if ns3.CurrentHealth(owner) != 0 {
			ns3.EnchantOff(owner, "ENCHANT_ETHEREAL")
		}
		if ns3.IsObjectOn(ptr) {
			ns3.Delete(ptr)
			ns3.Delete(ptr1)
		}
		break
	}
}

func BerserkerCharge(owner, target ns3.ObjectID) {
	if ns3.CurrentHealth(owner) != 0 && ns3.CurrentHealth(target) != 0 {
		ns3.EnchantOff(owner, "ENCHANT_INVULNERABLE")
		ns3.MoveWaypoint(2, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
		ns3.AudioEvent("BerserkerChargeInvoke", 2)
		ns3.MoveWaypoint(1, ns3.GetObjectX(owner), ns3.GetObjectY(owner))

		unit := ns3.CreateObject("InvisibleLightBlueHigh", 1)
		ns3.MoveWaypoint(1, ns3.GetObjectX(unit), ns3.GetObjectY(unit))

		unit1 := ns3.CreateObject("InvisibleLightBlueHigh", 1)
		ns3.LookAtObject(unit1, target)

		ns3.LookWithAngle(ns3.GetLastItem(owner), 0)
		ns3.SetCallback(owner, 9, BerserkerTouched)

		ratioX := UnitRatioX(target, owner, 23.0)
		ratioY := UnitRatioY(target, owner, 23.0)
		ns3.FrameTimer(1, func() {
			BerserkerLoop(unit, unit1, owner, target, ratioX, ratioY)
		})
	}
}

func BerserkerLoop(ptr, ptr1, owner, target ns3.ObjectID, ratioX, ratioY float32) {
	count := ns3.GetDirection(ptr)

	if ns3.CurrentHealth(owner) != 0 && count < 60 && ns3.IsObjectOn(ptr) && ns3.IsObjectOn(owner) {
		if CheckUnitFrontSight(owner, ratioX*1.5, ratioY*1.5) && ns3.GetDirection(ns3.GetLastItem(owner)) == 0 {
			ns3.MoveObject(owner, ns3.GetObjectX(owner)+ratioX, ns3.GetObjectY(owner)+ratioY)
			ns3.LookWithAngle(owner, ns3.GetDirection(ptr1))
			ns3.Walk(owner, ns3.GetObjectX(owner), ns3.GetObjectY(owner))
		} else {
			ns3.LookWithAngle(ptr, 100)
		}
		ns3.FrameTimer(1, func() {
			BerserkerLoop(ptr, ptr1, owner, target, ratioX, ratioY)
		})
	} else {
		ns3.SetCallback(owner, 9, NullCollide)
		ns3.Delete(ptr)
		ns3.Delete(ptr1)
	}
}

func BerserkerTouched() {
	self, other := ns3.GetTrigger(), ns3.GetCaller()
	if ns3.IsObjectOn(self) {
		for {
			if ns3.GetCaller() == 0 || (ns3.HasClass(other, "IMMOBILE") && !ns3.HasClass(other, "DOOR") && !ns3.HasClass(other, "TRIGGER")) && !ns3.HasClass(other, "DANGEROUS") {
				ns3.MoveWaypoint(2, ns3.GetObjectX(self), ns3.GetObjectY(self))
				ns3.AudioEvent("FleshHitStone", 2)

				ns3.Enchant(self, "ENCHANT_HELD", 2.0)
			} else if ns3.CurrentHealth(other) != 0 {
				if ns3.IsAttackedBy(self, other) {
					ns3.MoveWaypoint(2, ns3.GetObjectX(self), ns3.GetObjectY(self))
					ns3.AudioEvent("FleshHitFlesh", 2)
					ns3.Damage(other, self, 100, 2)
				} else {
					break
				}
			} else {
				break
			}
			ns3.LookWithAngle(ns3.GetLastItem(self), 1)
			break
		}
	}
	ns3.Wander(WarBot)
	ns3.SecondTimer(10, BerserkerChargeCooldownReset)
}
func NullCollide() {
}
func BerserkerChargeCooldownReset() {
	if RespawnCooldownDelay == 0 {
		BerserkerChargeCooldown = 0
	}
}
func GlobalCooldownReset() {
	GlobalCooldown = 0
}
func RespawnCooldownDelayReset() {
	RespawnCooldownDelay = 0
}

var (
	BotSpawn                ns3.ObjectID
	WarBot                  ns3.ObjectID
	WizBot                  ns3.ObjectID
	WizBotCorpse            ns3.ObjectID
	WarBotCorpse            ns3.ObjectID
	OutOfGameWar            ns3.WaypointID
	OutOfGameWiz            ns3.WaypointID
	WarSound                ns3.WaypointID
	WizSound                ns3.WaypointID
	WarCryCooldown          int
	EyeOfTheWolfCooldown    int
	BerserkerChargeCooldown int
	GlobalCooldown          int
	RespawnCooldownDelay    int
)
