package BotWars

import (
	"fmt"

	"github.com/noxworld-dev/noxscript/ns/v4"
	"github.com/noxworld-dev/opennox-lib/player"
)

type MyAccountData struct {
	Character struct {
		// General
		Registered bool
		Name       string
		Health     int
		Mana       int
		Level      int
		Class      player.Class
		// Order
	}
	botscript struct {
		// FFA Active Bots
		ActiveWarBots int
		ActiveConBots int
		ActiveWizBots int
		// Team Active Bots
		ActiveRedWarBots  int
		ActiveRedConBots  int
		ActiveRedWizBots  int
		ActiveBlueWarBots int
		ActiveBlueConBots int
		ActiveBlueWizBots int
		// Server settings
		BotDifficultySetting int
		BotManaSetting       bool
	}
}

func loadMyBotScriptData(pl ns.Player) MyAccountData {
	var data MyAccountData
	err := pl.Store(ns.Persistent{Name: "botscript"}).Get("my-quest-name", &data)
	if err != nil {
		fmt.Println("cannot read botscript data:", err)
	}
	return data
}

func saveMyBotScriptData(pl ns.Player, data MyAccountData) {
	err := pl.Store(ns.Persistent{Name: "botscript"}).Set("my-quest-name", &data)
	if err != nil {
		fmt.Println("cannot save botscript data:", err)
	}
}

func updateMyBotScriptData(pl ns.Player, fnc func(data *MyAccountData)) {
	data := loadMyBotScriptData(pl)
	fnc(&data)
	saveMyBotScriptData(pl, data)
}
