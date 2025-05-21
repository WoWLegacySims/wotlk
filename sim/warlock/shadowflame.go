package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func (warlock *Warlock) registerShadowflameSpell() {
	dotDamage := 161.0
	minDamage := 615.0
	maxDamage := 671.0

	warlock.ShadowflameDot = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61291},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,

		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Shadowflame Dot",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = dotDamage + 0.0667*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
			spell.SpellMetrics[target.UnitIndex].Hits++
		},
	})

	warlock.Shadowflame = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 61290},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.25,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * time.Duration(15),
			},
		},

		BonusCritRating: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5*core.CritRatingPerCritChance, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dmgFromSP := 0.1064 * spell.SpellPower()
			for _, target := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(minDamage, maxDamage) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					warlock.ShadowflameDot.Cast(sim, target)
				}
			}
		},
	})
}
