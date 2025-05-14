package deathknight

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func (dk *Deathknight) registerHornOfWinterSpell() {
	dbc := core.FindMaxRank(HornofWinterInfos, dk.Level)
	if dbc == nil {
		return
	}
	actionID := core.ActionID{SpellID: dbc.SpellID}
	rpMetrics := dk.NewRunicPowerMetrics(actionID)

	dk.HornOfWinter = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: 20 * time.Second,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AddRunicPower(sim, 10, rpMetrics)
		},
	})
}
