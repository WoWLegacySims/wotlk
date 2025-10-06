package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) registerBindingHealSpell() {
	dbc := priestinfo.BindingHeal.GetMaxRank(priest.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	spellCoeff := (dbc.GetCoefficient(0) + 0.04*float64(priest.Talents.EmpoweredHealing)) * dbc.GetLevelPenalty(priest.Level)

	priest.BindingHeal = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  priestinfo.BindingHeal.GetAllIDs(),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCrit: float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			(1 + .02*float64(priest.Talents.DivineProvidence)),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 0.5 * (1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve]),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healFromSP := spellCoeff * spell.HealingPower(target)

			selfHealing := sim.Roll(bp, die) + healFromSP
			spell.CalcAndDealHealing(sim, &priest.Unit, selfHealing, spell.OutcomeHealingCrit)

			targetHealing := sim.Roll(bp, die) + healFromSP
			spell.CalcAndDealHealing(sim, target, targetHealing, spell.OutcomeHealingCrit)
		},
	})
}
