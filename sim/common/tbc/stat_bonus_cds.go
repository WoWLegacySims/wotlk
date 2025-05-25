package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	helpers.NewSpellPowerActive(27828, 150, time.Second*20, time.Minute*2)
	helpers.NewSpellPowerActive(28040, 120, time.Second*15, time.Second*90)
	helpers.NewSpellPowerActive(28223, 167, time.Second*20, time.Minute*2)
	helpers.NewSpellPowerActive(29370, 158, time.Second*20, time.Minute*2) // Icon of the Silver Crescent
	helpers.NewSpellPowerActive(29376, 158, time.Second*20, time.Minute*2) // Essence of the Martyr
	helpers.NewSpellPowerActive(33829, 211, time.Second*20, time.Minute*2) // Hex Shrunken Head
	helpers.NewSpellPowerActive(38288, 153, time.Second*20, time.Minute*2) // Direbrew Hops
	helpers.NewSpellPowerActive(38290, 155, time.Second*20, time.Minute*2) // Dark Iron Smoking Pipe

	helpers.NewSpellPowerActive(34429, 320, time.Second*15, time.Second*90) // Shifting Naaru Sliver

	helpers.NewHasteActive(32483, 175, time.Second*20, time.Minute*2) // Skull of Gul'dan
	helpers.NewHasteActive(28288, 260, time.Second*10, time.Minute*2)

	helpers.NewAttackPowerActive(29383, 278, time.Second*20, time.Minute*2) // Bloodlust Brooch
	helpers.NewAttackPowerActive(33831, 278, time.Second*20, time.Minute*2) // Berserkers Call
	helpers.NewAttackPowerActive(38287, 278, time.Second*20, time.Minute*2) // Empty Direbrew Mug
	helpers.NewAttackPowerActive(28041, 200, time.Second*15, time.Second*90)

	helpers.NewBlockValueActive(29387, 200, time.Second*40, time.Minute*2) // Gnomeregan Autoblocker
	helpers.NewBlockValueActive(38289, 200, time.Second*40, time.Minute*2) // Coren's Lucky Coin

	helpers.NewHealthActive(32501, 1750, time.Second*20, time.Minute*3) // Shadowmoon Insignia
	helpers.NewHealthActive(28042, 900, time.Second*15, time.Minute*5)

	helpers.NewArmorActive(27891, 1280, time.Second*20, time.Minute*2)

	helpers.NewArmorPenActive(28121, 85, time.Second*20, time.Minute*2)

	helpers.NewDodgeActive(28528, 300, time.Second*10, time.Minute*2)

	core.AddEffectsToTest = false
	core.NewSimpleStatOffensiveTrinketEffect(32658, stats.Stats{stats.Agility: 150}, time.Second*20, time.Minute*2) // Badge of Tenacity

	core.NewSimpleStatDefensiveTrinketEffect(28484, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30)
	core.NewSimpleStatDefensiveTrinketEffect(28485, stats.Stats{stats.Health: 1500, stats.Strength: 150}, time.Second*15, time.Minute*30)
	core.AddEffectsToTest = true
}
