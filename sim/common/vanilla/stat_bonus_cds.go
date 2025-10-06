package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	helpers.NewHasteActive(9449, 500, time.Second*30, time.Hour)

	helpers.NewSpiritActive(43664, 50, time.Minute, time.Minute*5)

	helpers.NewDodgeActive(43667, 200, time.Second*15, time.Minute*5)
	core.AddEffectsToTest = false
	core.NewSimpleStatItemEffect(43656, stats.Stats{stats.Strength: 10, stats.Agility: 10, stats.Stamina: 10, stats.Intellect: 10, stats.Spirit: 10}, time.Minute, time.Minute*10)
	core.AddEffectsToTest = true
}
