package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) registerGreaterHealSpell() {
	dbc := priestinfo.GreaterHeal.GetMaxRank(priest.Level)
	if dbc == nil {
		dbc = priestinfo.Heal.GetMaxRank(priest.Level)
		if dbc == nil {
			dbc = priestinfo.LesserHeal.GetMaxRank(priest.Level)
			if dbc == nil {
				return
			}
		}
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	coef := dbc.GetCoefficient(0)
	spellCoeff := (coef + 0.08*float64(priest.Talents.EmpoweredHealing)) * dbc.GetLevelPenalty(priest.Level)

	priest.GreaterHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48063},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: dbc.BaseCost / 100,
			Multiplier: 1 *
				(1 - .05*float64(priest.Talents.ImprovedHealing)) *
				core.TernaryFloat64(priest.HasSetBonus(ItemSetRegaliaOfFaith, 4), .95, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: dbc.CastTime - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCrit: float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetVestmentsOfAbsolution, 4), 1.05, 1),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseHealing := sim.Roll(bp, die) + spellCoeff*spell.HealingPower(target)
			spell.CalcAndDealHealing(sim, target, baseHealing, spell.OutcomeHealingCrit)
		},
	})
}
