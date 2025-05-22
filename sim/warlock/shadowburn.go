package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerShadowBurnSpell() {
	if !warlock.Talents.Shadowburn {
		return
	}
	dbc := warlockinfo.Shadowburn.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(1, warlock.Level)
	coef := dbc.GetCoefficient(1) * dbc.GetLevelPenalty(warlock.Level) * (1 + 0.04*float64(warlock.Talents.ShadowAndFlame))

	if warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfShadowburn) {
		warlock.RegisterResetEffect(func(sim *core.Simulation) {
			sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute int32) {
				if isExecute == 35 {
					warlock.Shadowburn.BonusCrit += 20
				}
			})
		})
	}

	warlock.Shadowburn = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  warlockinfo.Shadowburn.GetAllIDs(),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.2,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault, // backdraft procs don't change the GCD of shadowburn
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		BonusCrit: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
