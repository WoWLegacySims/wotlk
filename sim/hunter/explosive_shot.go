package hunter

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/hunterinfo"
)

func (hunter *Hunter) registerExplosiveShotSpell(timer *core.Timer) {
	if !hunter.Talents.ExplosiveShot {
		return
	}

	hunter.ExplosiveShotR4 = hunter.makeExplosiveShotSpell(timer, false)
	hunter.ExplosiveShotR3 = hunter.makeExplosiveShotSpell(timer, true)
}

func (hunter *Hunter) makeExplosiveShotSpell(timer *core.Timer, downrank bool) *core.Spell {
	dbc := core.Ternary(downrank, hunterinfo.ExplosiveShot.GetDownRank(hunter.Level), hunterinfo.ExplosiveShot.GetMaxRank(hunter.Level))
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hunter.Level)
	rank := dbc.Rank

	return hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second * 6,
			},
		},

		BonusCrit: 0 +
			2*float64(hunter.Talents.SurvivalInstincts) +
			core.TernaryFloat64(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveShot), 4, 0),
		DamageMultiplierAdditive: 1 +
			.02*float64(hunter.Talents.TNT),
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(true, false, false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("Explosive Shot (Rank %d)", rank),
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + 0.14*dot.Spell.RangedAttackPower(target)
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeRangedHitAndCritSnapshot)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)

			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.TickOnce(sim)
			}
		},
	})
}
