package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerSoulFireSpell() {
	dbc := warlockinfo.SoulFire.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, warlock.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(warlock.Level)

	warlock.SoulFire = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:   warlockinfo.SoulFire.GetAllIDs(),
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.09,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * time.Duration(6000-400*warlock.Talents.Bane),
			},
		},

		BonusCrit: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
