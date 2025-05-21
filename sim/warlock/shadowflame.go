package warlock

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerShadowflameSpell() {
	dbc := warlockinfo.Shadowflame.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	dbcDot := warlockinfo.ShadowflameDot.GetByID(dbc.Effects[1].TriggerSpell)
	if dbcDot == nil {
		panic(fmt.Sprintf("No Dot found for Shadowflame %d", dbc.SpellID))
	}
	bp, die := dbc.GetBPDie(0, warlock.Level)
	bpDot, _ := dbcDot.GetBPDie(0, warlock.Level)

	warlock.ShadowflameDot = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbcDot.SpellID},
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
				dot.SnapshotBaseDamage = bpDot + 0.0667*dot.Spell.SpellPower()
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
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
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

		BonusCrit: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			dmgFromSP := 0.1064 * spell.SpellPower()
			for _, target := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(bp, die) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					warlock.ShadowflameDot.Cast(sim, target)
				}
			}
		},
	})
}
