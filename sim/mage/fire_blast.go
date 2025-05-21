package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerFireBlastSpell() {
	dbc := mageinfo.FireBlast.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second*8 - time.Second*time.Duration(mage.Talents.ImprovedFireBlast),
			},
		},

		BonusCrit: 0 +
			float64(mage.Talents.CriticalMass+mage.Talents.Incineration)*2,
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
