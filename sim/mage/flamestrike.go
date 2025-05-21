package mage

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerFlamestrikeSpell(downrank bool) *core.Spell {
	dbc := core.Ternary(downrank, mageinfo.Flamestrike.GetDownRank(mage.Level), mageinfo.Flamestrike.GetMaxRank(mage.Level))
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	bpDot, _ := dbc.GetBPDie(1, mage.Level)
	coefDot := dbc.GetCoefficient(1) * dbc.GetLevelPenalty(mage.Level)
	rank := dbc.Rank

	actionID := core.ActionID{SpellID: dbc.SpellID}.WithTag(rank)

	label := fmt.Sprintf("Flamestrike (Rank %d)", rank)

	return mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: dbc.BaseCost / 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second * 2,
			},
		},

		BonusCrit: float64(mage.Talents.CriticalMass+mage.Talents.WorldInFlames) * 2,
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.05*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: label,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot, _ bool) {
				target := mage.CurrentTarget
				dot.SnapshotBaseDamage = bpDot + coefDot*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := coef * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(bp, die) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
			spell.AOEDot().Apply(sim)
		},
	})
}

func (mage *Mage) registerFlamestrikeSpells() {
	mage.Flamestrike = mage.registerFlamestrikeSpell(false)
	mage.FlamestrikeRank8 = mage.registerFlamestrikeSpell(true)
}
