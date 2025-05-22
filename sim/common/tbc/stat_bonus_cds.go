package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
)

func init() {

	helpers.NewSpellPowerActive(34429, 320, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	helpers.NewHasteActive(32483, 175, time.Second*20, time.Minute*2) // Skull of Gul'dan

	helpers.NewAttackPowerActive(33831, 278, time.Second*20, time.Minute*2) // Berserkers Call
}
