package deathknight

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
)

var IcyTouchActionID = core.ActionID{SpellID: 59131}

func (dk *Deathknight) registerIcyTouchSpell() {
	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()

	dk.IcyTouch = dk.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			RunicPowerGain: 10 + 2.5*float64(dk.Talents.ChillOfTheGrave),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCritRating:  dk.rimeCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(dk.Talents.ImprovedIcyTouch),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := (sim.Roll(227, 245) + sigilBonus + 0.1*dk.getImpurityBonus(spell)) *
				dk.glacielRotBonus(target) *
				dk.RoRTSBonus(target) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.SpendRefundableCost(sim, result)

			if result.Landed() {
				dk.FrostFeverExtended[target.Index] = 0
				dk.FrostFeverSpell.Cast(sim, target)
			}

			spell.DealDamage(sim, result)
		},
	})
}
func (dk *Deathknight) registerDrwIcyTouchSpell() {
	sigilBonus := dk.sigilOfTheFrozenConscienceBonus()

	dk.RuneWeapon.IcyTouch = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    IcyTouchActionID,
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		//Flags:       core.SpellFlagIgnoreAttackerModifiers,

		BonusCritRating:  dk.rimeCritBonus() * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.05*float64(dk.Talents.ImprovedIcyTouch),
		CritMultiplier:   dk.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(227, 245) + sigilBonus + 0.1*dk.RuneWeapon.getImpurityBonus(spell)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				dk.RuneWeapon.FrostFeverSpell.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})
}
