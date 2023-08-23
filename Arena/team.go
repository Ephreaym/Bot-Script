package BWEstate

import (
	"github.com/noxworld-dev/noxscript/ns/v4"
)

var (
	Red  = NewTeam("Red")
	Blue = NewTeam("Blue")
)

func init() {
	Red.Enemy = Blue
	Blue.Enemy = Red
}

func NewTeam(name string) *Team {
	return &Team{Name: name}
}

type Team struct {
	Name            string
	Enemy           *Team
	Flag            ns.Obj
	FlagStart       ns.WaypointObj
	FlagIsAtBase    bool
	FlagInteraction bool
	TeamObj         ns.Obj
	TeamBase        ns.Obj
	TeamTank        ns.Obj
	spawns          []ns.Obj
}

// SpawnPoint selects a random PlayerStart for the bot to spawn on.
func (t *Team) SpawnPoint() ns.Pointf {
	if t.spawns == nil {
		// Filter to only select PlayStart objects that are owned by the team.
		t.spawns = ns.FindAllObjects(ns.HasTypeName{"PlayerStart"}) // <---- Use this when no teams are used.
	}
	if len(t.spawns) == 0 {
		return ns.GetHost().Pos()
	}
	i := ns.Random(0, len(t.spawns)-1)
	pick := t.spawns[i]
	return pick.Pos()
}

func (t *Team) init() {
}

func (t *Team) lateInit() {
	t.TeamObj = ns.Object("Team" + t.Name)
	//t.TeamBase = ns.Object(t.Name + "Base")
	//t.TeamTank = t.TeamObj
	//t.Flag = ns.Object(t.Name + "Flag")
	//t.FlagStart = ns.NewWaypoint(t.Name+"FlagStart", t.Flag.Pos())
	//t.FlagIsAtBase = true
	//t.FlagInteraction = false
	//ns.NewWaypoint(t.Name+"FlagWaypoint", t.Flag.Pos())
}

func (t *Team) SetNPCTeamColor() {

}

func (t *Team) FlagStartF() {
	//return ns.Waypoint(t.Name + "FlagStart")
}

func (t *Team) FlagReset() {
	//t.Flag.SetPos(t.FlagStart.Pos())
	//t.Flag.Enable(true)
}

func (t *Team) PreUpdate() {
	// Script for bots that moves the flag towards them each frame.
	//t.MoveEquipedFlagWithBot()
	//t.SetBasePosition()
}

func (t *Team) PostUpdate() {
	//t.CheckIfFlagsAreAtBase()
	//t.BotConditionsWhileCarryingTheFlag()
}

func (t *Team) MoveEquipedFlagWithBot() {
	//if !t.Flag.IsEnabled() {
	// Move the real flag out of the map.
	// Move the fake flag on the bot.
	//t.Flag.SetPos(t.Enemy.TeamTank.Pos())
	//}
}

func (t *Team) CheckIfFlagsAreAtBase() {
	//if (ns.InCirclef{Center: t.TeamBase, R: 20}).Matches(t.Flag) {
	//	t.FlagIsAtBase = true
	//} else {
	//	t.FlagIsAtBase = false
	//}
}

func (t *Team) SetBasePosition() {
	//if InitLoadComplete {
	//	t.TeamBase.SetPos(t.Flag.Pos())
	//}
}

func (t *Team) BotConditionsWhileCarryingTheFlag() {
	// Remove buffs from bots that can't be active while carrying the flag.
	//if t.TeamTank.HasEnchant(enchant.INVISIBLE) {
	//	t.TeamTank.EnchantOff(enchant.INVISIBLE)
	//}
	//if t.TeamTank.HasEnchant(enchant.INVULNERABLE) {
	//	t.TeamTank.EnchantOff(enchant.INVULNERABLE)
	//}
	//if !t.TeamTank.HasEnchant(enchant.VILLAIN) {
	//	t.TeamTank.Enchant(enchant.VILLAIN, ns.Seconds(60))
	//}
}
