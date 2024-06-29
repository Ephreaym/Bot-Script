package BotWars

import (
	"time"

	"github.com/noxworld-dev/noxscript/ns/v4"
)

type Updater struct {
	last time.Duration
}

func (u *Updater) EachSec(sec float64, fnc func()) {
	dt := time.Duration(sec * float64(time.Second))
	now := ns.Now()
	if u.last != 0 && u.last+dt < now {
		return
	}
	u.last = now
	fnc()
}

func (u *Updater) EachFrame(frame int, fnc func()) {
	if ns.Frame()%frame == 0 {
		fnc()
	}
}
