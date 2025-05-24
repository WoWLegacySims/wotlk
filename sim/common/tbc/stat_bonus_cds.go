package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	helpers.NewSpellPowerActive(29370, 158, time.Second*20, time.Minute*2) // Icon of the Silver Crescent
	helpers.NewSpellPowerActive(29376, 158, time.Second*20, time.Minute*2) // Essence of the Martyr
	helpers.NewSpellPowerActive(33829, 211, time.Second*20, time.Minute*2) // Hex Shrunken Head
	helpers.NewSpellPowerActive(38288, 153, time.Second*20, time.Minute*2) // Direbrew Hops
	helpers.NewSpellPowerActive(38290, 155, time.Second*20, time.Minute*2) // Dark Iron Smoking Pipe

	helpers.NewSpellPowerActive(34429, 320, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	helpers.NewHasteActive(32483, 175, time.Second*20, time.Minute*2) // Skull of Gul'dan

	helpers.NewAttackPowerActive(29383, 278, time.Second*20, time.Minute*2) // Bloodlust Brooch
	helpers.NewAttackPowerActive(33831, 278, time.Second*20, time.Minute*2) // Berserkers Call
	helpers.NewAttackPowerActive(38287, 278, time.Second*20, time.Minute*2) // Empty Direbrew Mug

	helpers.NewBlockValueActive(29387, 200, time.Second*40, time.Minute*2) // Gnomeregan Autoblocker
	helpers.NewBlockValueActive(38289, 200, time.Second*40, time.Minute*2) // Coren's Lucky Coin

	helpers.NewHealthActive(32501, 1750, time.Second*20, time.Minute*3) // Shadowmoon Insignia

	core.NewSimpleStatOffensiveTrinketEffect(32658, stats.Stats{stats.Agility: 150}, time.Second*20, time.Minute*2) // Badge of Tenacity
}
