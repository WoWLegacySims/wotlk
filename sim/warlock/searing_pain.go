package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerSearingPainSpell() {
	dbc := warlockinfo.SearingPain.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, warlock.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(warlock.Level)

	warlock.SearingPain = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  warlockinfo.SearingPain.GetAllIDs(),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.08,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 1500 * time.Millisecond,
			},
		},
		BonusCrit: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5, 0) +
			[]float64{0, .04, .07, .10}[warlock.Talents.ImprovedSearingPain],
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm),
		CritMultiplier: warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5+
			core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfSearingPain), 0.2, 0)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
