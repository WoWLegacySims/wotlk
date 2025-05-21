package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

func (druid *Druid) registerDemoralizingRoarSpell() {
	dbc := druidinfo.DemoralizingRoar.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}

	druid.DemoralizingRoar = druid.RegisterSpell(Bear, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  62 * 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealOutcome(sim, aoeTarget, spell.OutcomeMagicHit)
				if result.Landed() {
					druid.DemoralizingRoarAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{druid.DemoralizingRoarAuras},
	})
}

func (druid *Druid) ShouldDemoralizingRoar(sim *core.Simulation, filler bool, maintainOnly bool) bool {
	if !druid.DemoralizingRoar.CanCast(sim, druid.CurrentTarget) {
		return false
	}

	if filler {
		return true
	}

	refreshWindow := time.Second * 2

	if (druid.MangleBear != nil) && (!druid.MangleBear.IsReady(sim)) {
		refreshWindow = druid.MangleBear.ReadyAt() - sim.CurrentTime + core.GCDDefault
	}

	return maintainOnly &&
		druid.DemoralizingRoarAuras.Get(druid.CurrentTarget).ShouldRefreshExclusiveEffects(sim, refreshWindow)
}
