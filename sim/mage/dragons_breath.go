package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func (mage *Mage) registerDragonsBreathSpell() {
	if !mage.Talents.DragonsBreath {
		return
	}

	mage.DragonsBreath = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 42950},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 20,
			},
		},
		BonusCrit:                float64(mage.Talents.CriticalMass + mage.Talents.WorldInFlames),
		DamageMultiplierAdditive: 1 + .02*float64(mage.Talents.FirePower),
		CritMultiplier:           mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier:         1 - 0.1*float64(mage.Talents.BurningSoul),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(1101, 1279) + 0.193*spell.SpellPower()
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
