package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:       "Blackened Naaru Sliver",
		ID:         34427,
		Duration:   time.Second * 20,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
}
