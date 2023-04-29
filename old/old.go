package basicmap_old

import . "github.com/noxworld-dev/noxscript/ns/v3"

func MapInitialize() {
	WarBot = Object("WarBot")
	WizBot = Object("WizBot")
	WarBotCorpse = Object("WarBotCorpse")
	WizBotCorpse = Object("WizBotCorpse")
	OutOfGameWar = Waypoint("OutOfGameWar")
	OutOfGameWiz = Waypoint("OutOfGameWiz")
	WarSound = Waypoint("WarSound")
	WizSound = Waypoint("WizSound")
	WarCryCooldown = 0
	EyeOfTheWolfCooldown = 0
	BerserkerChargeCooldown = 0
	GlobalCooldown = 0
	RespawnCooldownDelay = 0
	SecondTimer(2, Respawn)
	SecondTimer(2, RespawnWiz)
}
func CheckCurrentHealth() {
	if CurrentHealth(WarBot) < 1001 {
		ChangeScore(GetCaller(), 1)
		DeadBotWar()
	}
	if CurrentHealth(WizBot) < 1001 {
		ChangeScore(GetCaller(), 1)
		DeadBotWiz()
	}
}
func DeadBotWar() {
	MoveWaypoint(WarSound, GetObjectX(WarBot), GetObjectY(WarBot))
	AudioEvent("NPCDie", WarSound)
	MoveObject(WarBotCorpse, GetObjectX(WarBot), GetObjectY(WarBot))
	MoveObject(WarBot, GetWaypointX(OutOfGameWar), GetWaypointY(OutOfGameWar))
	SecondTimer(2, Respawn)
}
func DeadBotWiz() {
	MoveWaypoint(WizSound, GetObjectX(WizBot), GetObjectY(WizBot))
	AudioEvent("NPCDie", WizSound)
	MoveObject(WizBotCorpse, GetObjectX(WizBot), GetObjectY(WizBot))
	MoveObject(WizBot, GetWaypointX(OutOfGameWiz), GetWaypointY(OutOfGameWiz))
	SecondTimer(2, RespawnWiz)
}
func Respawn() {
	RestoreHealth(WarBot, MaxHealth(WarBot)-CurrentHealth(WarBot))
	WarCryCooldown = 0
	EyeOfTheWolfCooldown = 0
	BerserkerChargeCooldown = 0
	MoveWaypoint(WarSound, GetObjectX(WarBotCorpse), GetObjectY(WarBotCorpse))
	AudioEvent("BlinkCast", WarSound)
	BotSpawn = Object("BotSpawn" + IntToString(Random(1, 14)))
	MoveObject(WarBot, GetObjectX(BotSpawn), GetObjectY(BotSpawn))
	MoveObject(WarBotCorpse, GetWaypointX(OutOfGameWar), GetWaypointY(OutOfGameWar))
	Enchant(WarBot, "ENCHANT_INVULNERABLE", 5.0)
	RespawnCooldownDelay = 1
	SecondTimer(10, RespawnCooldownDelayReset)
}
func RespawnWiz() {
	RestoreHealth(WizBot, MaxHealth(WizBot)-CurrentHealth(WizBot))
	MoveWaypoint(WizSound, GetObjectX(WizBotCorpse), GetObjectY(WizBotCorpse))
	AudioEvent("BlinkCast", WizSound)
	BotSpawn = Object("BotSpawn" + IntToString(Random(1, 14)))
	MoveObject(WizBot, GetObjectX(BotSpawn), GetObjectY(BotSpawn))
	MoveObject(WizBotCorpse, GetWaypointX(OutOfGameWiz), GetWaypointY(OutOfGameWiz))
	Enchant(WizBot, "ENCHANT_INVULNERABLE", 5.0)
}
func WarCry() {
	self, other := GetTrigger(), GetCaller()
	if MaxHealth(other) == 150 {
	} else {
		if (WarCryCooldown == 0) && (GlobalCooldown == 0) {
			MoveWaypoint(WarSound, GetObjectX(WarBot), GetObjectY(WarBot))
			AudioEvent("WarcryInvoke", WarSound)
			PauseObject(WarBot, 45)
			Enchant(WarBot, "ENCHANT_HELD", 1.0)
			Enchant(other, "ENCHANT_ANTI_MAGIC", 3.0)
			EnchantOff(self, "ENCHANT_SHOCK")
			EnchantOff(self, "ENCHANT_INVULNERABLE")
			WarCryCooldown = 1
			GlobalCooldown = 1
			SecondTimer(10, WarCryCooldownReset)
			SecondTimer(1, GlobalCooldownReset)
		} else {
		}
	}
}
func WarCryCooldownReset() {
	if RespawnCooldownDelay == 0 {
		WarCryCooldown = 0
	} else {
	}
}
func EyeOfTheWolf() {
	self := GetTrigger()
	Wander(WarBot)
	if EyeOfTheWolfCooldown == 0 {
		Enchant(self, "ENCHANT_INFRAVISION", 10.0)
		EyeOfTheWolfCooldown = 1
		SecondTimer(20, EyeOfTheWolfCooldownReset)
	} else {
	}
}
func EyeOfTheWolfCooldownReset() {
	if RespawnCooldownDelay == 0 {
		EyeOfTheWolfCooldown = 0
	} else {
	}
}

// ------------------------------------- BERSERKERCHARGE SCRIPT ---------------------------------------------- //

func UnitRatioX(unit, target ObjectID, size float32) float32 {
	return (GetObjectX(unit) - GetObjectX(target)) * size / Distance(GetObjectX(unit), GetObjectY(unit), GetObjectX(target), GetObjectY(target))
}

func UnitRatioY(unit, target ObjectID, size float32) float32 {
	return (GetObjectY(unit) - GetObjectY(target)) * size / Distance(GetObjectX(unit), GetObjectY(unit), GetObjectX(target), GetObjectY(target))
}
func WarBotDetectEnemy() {
	if (BerserkerChargeCooldown == 0) && (GlobalCooldown == 0) {
		rnd := Random(0, 2)

		if (rnd == 0) || (rnd == 1) {
			BerserkerChargeCooldown = 1
			GlobalCooldown = 1
			SecondTimer(1, GlobalCooldownReset)
			BerserkerInRange(GetTrigger(), GetCaller(), 10)
		}
	} else {
		if WarCryCooldown == 0 {
			WarCry()
		} else {
		}
	}
}
func CheckUnitFrontSight(unit ObjectID, dtX, dtY float32) bool {
	MoveWaypoint(1, GetObjectX(unit)+dtX, GetObjectY(unit)+dtY)
	temp := CreateObject("InvisibleLightBlueHigh", 1)
	res := IsVisibleTo(unit, temp)

	Delete(temp)
	return res
}
func pointedByPlr() {
}
func BerserkerInRange(owner, target ObjectID, wait int) {
	if CurrentHealth(owner) != 0 && CurrentHealth(target) != 0 {
		if !HasEnchant(owner, "ENCHANT_ETHEREAL") {
			Enchant(owner, "ENCHANT_ETHEREAL", 0.0)
			MoveWaypoint(1, GetObjectX(owner), GetObjectY(owner))
			unit := CreateObject("InvisibleLightBlueHigh", 1)
			MoveWaypoint(1, GetObjectX(unit), GetObjectY(unit))
			CreateObject("InvisibleLightBlueHigh", 1)
			Raise(unit, ToFloat(owner))
			Raise(unit+1, ToFloat(target))
			LookWithAngle(unit, wait)
			FrameTimerWithArg(1, unit, BerserkerWaitStrike)
		}
	}
}

func BerserkerWaitStrike(ptr ObjectID) {
	count := GetDirection(ptr)
	owner := ToInt(GetObjectZ(ptr))
	target := ToInt(GetObjectZ(ptr + 1))

	for {
		if IsObjectOn(ptr) && CurrentHealth(owner) != 0 && CurrentHealth(target) != 0 && IsObjectOn(owner) {
			if count != 0 {
				if IsVisibleTo(owner, target) && Distance(GetObjectX(owner), GetObjectY(owner), GetObjectX(target), GetObjectY(target)) < 400.0 {
					BerserkerCharge(owner, target)
				} else {
					LookWithAngle(ptr, count-1)
					FrameTimerWithArg(6, ptr, BerserkerWaitStrike)
					break
				}
			}
		}
		if CurrentHealth(owner) != 0 {
			EnchantOff(owner, "ENCHANT_ETHEREAL")
		}
		if IsObjectOn(ptr) {
			Delete(ptr)
			Delete(ptr + 1)
		}
		break
	}
}

func BerserkerCharge(owner, target ObjectID) {
	if CurrentHealth(owner) != 0 && CurrentHealth(target) != 0 {
		EnchantOff(owner, "ENCHANT_INVULNERABLE")
		MoveWaypoint(2, GetObjectX(owner), GetObjectY(owner))
		AudioEvent("BerserkerChargeInvoke", 2)
		MoveWaypoint(1, GetObjectX(owner), GetObjectY(owner))
		unit := CreateObject("InvisibleLightBlueHigh", 1)
		MoveWaypoint(1, GetObjectX(unit), GetObjectY(unit))
		Raise(CreateObject("InvisibleLightBlueHigh", 1), UnitRatioX(target, owner, 23.0))
		Raise(CreateObject("InvisibleLightBlueHigh", 1), UnitRatioY(target, owner, 23.0))
		CreateObject("InvisibleLightBlueHigh", 1)
		Raise(unit+3, ToFloat(owner))
		LookWithAngle(GetLastItem(owner), 0)
		SetCallback(owner, 9, BerserkerTouched)
		Raise(unit, ToFloat(target))
		LookAtObject(unit+1, target)
		FrameTimerWithArg(1, unit, BerserkerLoop)
	}
}

func BerserkerLoop(ptr ObjectID) {
	owner := ToInt(GetObjectZ(ptr + 3))
	count := GetDirection(ptr)

	if CurrentHealth(owner) != 0 && count < 60 && IsObjectOn(ptr) && IsObjectOn(owner) {
		if CheckUnitFrontSight(owner, GetObjectZ(ptr+1)*1.5, GetObjectZ(ptr+2)*1.5) && GetDirection(GetLastItem(owner)) == 0 {
			MoveObject(owner, GetObjectX(owner)+GetObjectZ(ptr+1), GetObjectY(owner)+GetObjectZ(ptr+2))
			LookWithAngle(owner, GetDirection(ptr+1))
			Walk(owner, GetObjectX(owner), GetObjectY(owner))
		} else {
			LookWithAngle(ptr, 100)
		}
		FrameTimerWithArg(1, ptr, BerserkerLoop)
	} else {
		SetCallback(owner, 9, NullCollide)
		Delete(ptr)
		Delete(ptr + 1)
		Delete(ptr + 2)
		Delete(ptr + 3)
	}
}

func BerserkerTouched() {
	self, other := GetTrigger(), GetCaller()
	if IsObjectOn(self) {
		for {
			if GetCaller() == 0 || (HasClass(other, "IMMOBILE") && !HasClass(other, "DOOR") && !HasClass(other, "TRIGGER")) && !HasClass(other, "DANGEROUS") {
				MoveWaypoint(2, GetObjectX(self), GetObjectY(self))
				AudioEvent("FleshHitStone", 2)

				Enchant(self, "ENCHANT_HELD", 2.0)
			} else if CurrentHealth(other) != 0 {
				if IsAttackedBy(self, other) {
					MoveWaypoint(2, GetObjectX(self), GetObjectY(self))
					AudioEvent("FleshHitFlesh", 2)
					Damage(other, self, 100, 2)
				} else {
					break
				}
			} else {
				break
			}
			LookWithAngle(GetLastItem(self), 1)
			break
		}
	}
	Wander(WarBot)
	SecondTimer(10, BerserkerChargeCooldownReset)
}
func NullCollide() {
	return
}
func BerserkerChargeCooldownReset() {
	if RespawnCooldownDelay == 0 {
		BerserkerChargeCooldown = 0
	} else {
	}
}
func GlobalCooldownReset() {
	GlobalCooldown = 0
}
func RespawnCooldownDelayReset() {
	RespawnCooldownDelay = 0
}

var (
	BotSpawn                ObjectID
	WarBot                  ObjectID
	WizBot                  ObjectID
	WizBotCorpse            ObjectID
	WarBotCorpse            ObjectID
	OutOfGameWar            WaypointID
	OutOfGameWiz            WaypointID
	WarSound                WaypointID
	WizSound                WaypointID
	WarCryCooldown          int
	EyeOfTheWolfCooldown    int
	BerserkerChargeCooldown int
	GlobalCooldown          int
	RespawnCooldownDelay    int
)
