package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) registerCircleOfHealingSpell() {
	if !priest.Talents.CircleOfHealing {
		return
	}
	dbc := priestinfo.CircleofHealing.GetMaxRank(priest.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(priest.Level)

	numTargets := 5 + core.TernaryInt32(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfCircleOfHealing), 1, 0)
	targets := priest.Env.Raid.GetFirstNPlayersOrPets(numTargets)

	priest.CircleOfHealing = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  priestinfo.CircleofHealing.GetAllIDs(),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellHealing,
		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.21,
			Multiplier: 1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCrit: float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 *
			(1 + .02*float64(priest.Talents.SpiritualHealing)) *
			(1 + .01*float64(priest.Talents.BlessedResilience)) *
			(1 + .02*float64(priest.Talents.FocusedPower)) *
			(1 + .02*float64(priest.Talents.DivineProvidence)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetCrimsonAcolytesRaiment, 4), 1.1, 1),
		CritMultiplier:   priest.DefaultHealingCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healFromSP := coef * spell.HealingPower(target)
			for _, aoeTarget := range targets {
				baseHealing := sim.Roll(bp, die) + healFromSP
				spell.CalcAndDealHealing(sim, aoeTarget, baseHealing, spell.OutcomeHealingCrit)
			}
		},
	})
}
