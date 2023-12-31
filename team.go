package BotWars

import (
	"fmt"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/noxscript/ns/v4/audio"
	"github.com/noxworld-dev/noxscript/ns/v4/enchant"
	"github.com/noxworld-dev/opennox-lib/object"
)

var (
	Red  = NewTeam(0, "Red")
	Blue = NewTeam(1, "Blue")
)

func init() {
	Red.Enemy = Blue
	Blue.Enemy = Red
}

func NewTeam(ind int, name string) *Team {
	return &Team{team: ns.Teams()[ind], Name: name}
}

type Team struct {
	team            ns.Team
	Name            string
	Enemy           *Team
	Flag            ns.Obj
	FlagStart       ns.WaypointObj
	FlagIsAtBase    bool
	FlagInteraction bool
	TeamBase        ns.Obj
	TeamTank        ns.Obj
	spawns          []ns.Obj
}

func (t *Team) Spawns() []ns.Obj {
	if GameModeIsCTF {
		if t.spawns == nil {
			// Filter to only select PlayStart objects that are owned by the team.
			//filter := ns.HasTypeName{"PlayerStart"}
			//ns.ObjectGroup("Team"+t.Name).EachObject(true, func(it ns.Obj) bool {
			//	if filter.Matches(it) {
			//		t.spawns = append(t.spawns, it)
			//	}
			//	return true // keep iterating in any case
			//})
			t.spawns = ns.FindAllObjects(
				ns.HasTypeName{"PlayerStart"},
				ns.HasTeam{t.Team()},
			) // <---- Use this when no teams are used.
		}
		if len(t.spawns) == 0 {
			return []ns.Obj{ns.GetHost()}
		}
	} else {
		if t.spawns == nil {
			// Filter to only select PlayStart objects that are owned by the team.
			//filter := ns.HasTypeName{"PlayerStart"}
			//ns.ObjectGroup("Team"+t.Name).EachObject(true, func(it ns.Obj) bool {
			//	if filter.Matches(it) {
			//		t.spawns = append(t.spawns, it)
			//	}
			//	return true // keep iterating in any case
			//})
			t.spawns = ns.FindAllObjects(
				ns.HasTypeName{"PlayerStart"},
				//ns.HasTeam{t.Team()},
			) // <---- Use this when no teams are used.
		}
		if len(t.spawns) == 0 {
			return []ns.Obj{ns.GetHost()}
		}
	}
	return t.spawns
}

// SpawnPoint selects a random PlayerStart for the bot to spawn on.
func (t *Team) SpawnPoint() ns.Pointf {
	spawns := t.Spawns()
	i := ns.Random(0, len(spawns)-1)
	pick := spawns[i]
	return pick.Pos()
}

func (t *Team) Team() ns.Team {
	return t.team
}

func (t *Team) init() {
}

func (t *Team) lateInit() {
	if GameModeIsCTF {
		Red.Flag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[0]})
		Blue.Flag = ns.FindObject(ns.HasTypeName{"Flag"}, ns.HasTeam{ns.Teams()[1]})
		BlueTeamBase = ns.CreateObject("ExtentBoxSmall", Blue.Flag)
		BlueTeamBase.FlagsEnable(object.FlagNoCollide)
		BlueTeamBase.FlagsEnable(object.FlagNoPushCharacters)
		RedTeamBase = ns.CreateObject("ExtentBoxSmall", Red.Flag)
		RedTeamBase.FlagsEnable(object.FlagNoCollide)
		RedTeamBase.FlagsEnable(object.FlagNoPushCharacters)
		//t.TeamBase = ns.Object(t.Name + "Base")
		Red.TeamBase = RedTeamBase
		Blue.TeamBase = BlueTeamBase
		t.TeamTank = t.Spawns()[0]
		t.FlagStart = ns.NewWaypoint(t.Name+"FlagStart", t.Flag.Pos())
		t.FlagIsAtBase = true
		t.FlagInteraction = false
		ns.NewWaypoint(t.Name+"FlagWaypoint", t.Flag.Pos())
	}
}

func (t *Team) FlagStartF() ns.WaypointObj {
	return ns.Waypoint(t.Name + "FlagStart")
}

func (t *Team) FlagReset() {
	if GameModeIsCTF {
		t.Flag.SetPos(t.FlagStart.Pos())
		t.Flag.Enable(true)
	}
}

func (t *Team) PreUpdate() {
	if GameModeIsCTF {
		// Script for bots that moves the flag towards them each frame.
		t.MoveEquipedFlagWithBot()
		t.SetBasePosition()
	}
}

func (t *Team) PostUpdate() {
	if GameModeIsCTF {
		t.CheckIfFlagsAreAtBase()
		t.BotConditionsWhileCarryingTheFlag()
	}
}

func (t *Team) MoveEquipedFlagWithBot() {
	if GameModeIsCTF {
		if !t.Flag.IsEnabled() {
			// Move the real flag out of the map.
			// Move the fake flag on the bot.
			t.Flag.SetPos(t.Enemy.TeamTank.Pos())
		}
	}
}

func (t *Team) CheckIfFlagsAreAtBase() {
	if GameModeIsCTF {
		if (ns.InCirclef{Center: t.TeamBase, R: 20}).Matches(t.Flag) {
			t.FlagIsAtBase = true
		} else {
			t.FlagIsAtBase = false
		}
	}
}

func (t *Team) SetBasePosition() {
	if GameModeIsCTF {
		if InitLoadComplete {
			t.TeamBase.SetPos(t.Flag.Pos())
		}
	}
}

func (t *Team) BotConditionsWhileCarryingTheFlag() {
	if GameModeIsCTF {
		// Remove buffs from bots that can't be active while carrying the flag.
		if t.TeamTank.HasEnchant(enchant.INVISIBLE) {
			t.TeamTank.EnchantOff(enchant.INVISIBLE)
		}
		if t.TeamTank.HasEnchant(enchant.INVULNERABLE) {
			t.TeamTank.EnchantOff(enchant.INVULNERABLE)
		}
		if !t.TeamTank.HasEnchant(enchant.VILLAIN) {
			t.TeamTank.Enchant(enchant.VILLAIN, ns.Seconds(60))
		}
	}
}

// CTF game mechanics.
// Pick up the enemy flag.
func (t *Team) CheckPickUpEnemyFlag(flag, u ns.Obj) {
	enemyFlag := t.Enemy.Flag
	if flag == enemyFlag {
		enemyFlag.Enable(false)
		soundToAllPlayers1 := ns.Players()
		t.Enemy.FlagInteraction = false
		for i := 0; i < len(soundToAllPlayers1); i++ {
			ns.AudioEvent(audio.FlagPickup, soundToAllPlayers1[i].Unit())
		}
		// Customize code below for individual unit.
		t.TeamTank = u
		t.TeamTank.AggressionLevel(0.16)
		t.TeamTank.WalkTo(t.Flag.Pos())
		ns.PrintStrToAll(fmt.Sprintf("Team %s has the %s flag!", t.Name, t.Enemy.Name))
	}
}

// Capture the flag.
func (t *Team) CheckCaptureEnemyFlag(flag, u ns.Obj) {
	if flag == t.Flag && t.FlagIsAtBase && u == t.TeamTank {
		soundToAllPlayers2 := ns.Players()
		for i := 0; i < len(soundToAllPlayers2); i++ {
			ns.AudioEvent(audio.FlagCapture, soundToAllPlayers2[i].Unit())
		}
		t.TeamTank = t.Spawns()[0]
		// new code
		t.Team().ChangeScore(+1)
		//
		t.Enemy.FlagInteraction = false
		t.FlagReset()
		t.Enemy.FlagReset()
		u.AggressionLevel(0.83)
		u.WalkTo(t.Enemy.Flag.Pos())
		ns.PrintStrToAll(fmt.Sprintf("Team %s has captured the %s flag!", t.Name, t.Enemy.Name))
	}
}

// Retrieve own flag.
func (t *Team) CheckRetrievedOwnFlag(flag, u ns.Obj) {
	if flag == t.Flag && !t.FlagIsAtBase {
		t.FlagIsAtBase = true
		t.FlagInteraction = false
		soundToAllPlayers3 := ns.Players()
		for i := 0; i < len(soundToAllPlayers3); i++ {
			ns.AudioEvent(audio.FlagRespawn, soundToAllPlayers3[i].Unit())
		}
		t.Flag.SetPos(t.FlagStart.Pos())
		u.WalkTo(t.TeamBase.Pos())
		ns.PrintStrToAll(fmt.Sprintf("Team %s has retrieved the flag!", t.Name))
		t.TeamTank.WalkTo(t.Flag.Pos())
	}
}

// Drop flag.
func (t *Team) DropEnemyFlag(u ns.Obj) {
	if u == t.TeamTank {
		t.Enemy.FlagInteraction = true
		soundToAllPlayers4 := ns.Players()
		for i := 0; i < len(soundToAllPlayers4); i++ {
			ns.AudioEvent(audio.FlagDrop, soundToAllPlayers4[i].Unit())
		}
		t.Enemy.Flag.Enable(true)
		t.TeamTank = t.Spawns()[0]
		ns.PrintStrToAll(fmt.Sprintf("Team %s has dropped the %s flag!", t.Name, t.Enemy.Name))
		ns.NewTimer(ns.Seconds(30), func() {
			if t.Enemy.Flag.IsEnabled() && t.Enemy.FlagInteraction {
				t.ReturnFlagHome(u)
			}
		})
	}
}

// Return flag home.
func (t *Team) ReturnFlagHome(u ns.Obj) {
	t.Enemy.Flag.SetPos(t.Enemy.FlagStart.Pos())
	soundToAllPlayers5 := ns.Players()
	for i := 0; i < len(soundToAllPlayers5); i++ {
		ns.AudioEvent(audio.FlagRespawn, soundToAllPlayers5[i].Unit())
	}
	ns.PrintStrToAll(fmt.Sprintf("The %s flag has returned home.", t.Enemy.Name))
}

func (t *Team) WalkToOwnFlag(u ns.Obj) {
	if !t.FlagIsAtBase && t.Flag.IsEnabled() {
		u.AggressionLevel(0.16)
		u.WalkTo(t.Flag.Pos())
	} else {
		t.CheckAttackOrDefend(u)
	}
}

func (t *Team) CheckAttackOrDefend(u ns.Obj) {
	if u == t.TeamTank {
		u.AggressionLevel(0.16)
		u.Guard(t.TeamBase.Pos(), t.TeamBase.Pos(), 20)
	} else if t.Flag.IsEnabled() {
		u.AggressionLevel(0.83)
		u.WalkTo(t.Enemy.Flag.Pos())
	} else if !t.Enemy.Flag.IsEnabled() {
		u.AggressionLevel(0.83)
		u.WalkTo(t.Flag.Pos())
	}
}

func (t *Team) DialogStart(u ns.Obj) {
	u.Chat("Con03B.scr:Worker1ChatA")
	u.SetOwner(ns.GetCaller())
	u.Follow(ns.GetCaller())
	// Con03B.scr:Worker1ChatD : I'll wait here
	// Con03B.scr:Worker1ChatA : Let's go
	// Con03B.scr:Worker1ChatB : Follow me

}

func (t *Team) DialogEnd(u ns.Obj) {

}
