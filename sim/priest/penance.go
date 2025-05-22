package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) registerPenanceHealSpell() {
	priest.PenanceHeal = priest.makePenanceSpell(true)
}

func (priest *Priest) RegisterPenanceSpell() {
	priest.Penance = priest.makePenanceSpell(false)
}

func (priest *Priest) makePenanceSpell(isHeal bool) *core.Spell {
	if !priest.Talents.Penance {
		return nil
	}

	dbc := priestinfo.Penance.GetMaxRank(priest.Level)
	dbcDmg := priestinfo.PenanceDamage.GetMaxRank(priest.Level)
	dbcHeal := priestinfo.PenanceHeal.GetMaxRank(priest.Level)
	if dbc == nil || dbcDmg == nil || dbcHeal == nil {
		return nil
	}
	bpDmg, _ := dbcDmg.GetBPDie(0, priest.Level)
	coefDmg := dbcDmg.GetCoefficient(0) * dbcDmg.GetLevelPenalty(priest.Level)

	bpHeal, dieHeal := dbcHeal.GetBPDie(0, priest.Level)
	coefHeal := dbcHeal.GetCoefficient(0) * dbcHeal.GetLevelPenalty(priest.Level)

	var procMask core.ProcMask
	flags := core.SpellFlagChanneled | core.SpellFlagAPL
	if isHeal {
		flags |= core.SpellFlagHelpful
		procMask = core.ProcMaskSpellHealing
	} else {
		procMask = core.ProcMaskSpellDamage
	}

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  priestinfo.Penance.GetAllIDs(),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    procMask,
		Flags:       flags,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.16,
			Multiplier: 1 *
				(1 - 0.05*float64(priest.Talents.ImprovedHealing)) *
				(1 - []float64{0, .04, .07, .10}[priest.Talents.MentalAgility]),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Duration(float64(time.Second*12-core.TernaryDuration(priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfPenance), time.Second*2, 0)) * (1 - .1*float64(priest.Talents.Aspiration))),
			},
		},

		DamageMultiplier: 1 +
			core.TernaryFloat64(isHeal,
				1*
					(1+.02*float64(priest.Talents.SpiritualHealing))*
					(1+.01*float64(priest.Talents.BlessedResilience))*
					(1+.02*float64(priest.Talents.FocusedPower)),
				.05*float64(priest.Talents.SearingLight)) +
			.01*float64(priest.Talents.TwinDisciplines),
		CritMultiplier:   core.TernaryFloat64(isHeal, priest.DefaultHealingCritMultiplier(), priest.DefaultSpellCritMultiplier()),
		ThreatMultiplier: 0,

		Dot: core.Ternary(!isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := bpDmg + coefDmg*dot.Spell.SpellPower()
				dot.Spell.CalcAndDealPeriodicDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			},
		}, core.DotConfig{}),
		Hot: core.Ternary(isHeal, core.DotConfig{
			Aura: core.Aura{
				Label: "Penance",
			},
			NumberOfTicks:       2,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseHealing := sim.Roll(bpHeal, dieHeal) + coefHeal*dot.Spell.HealingPower(target)
				dot.Spell.CalcAndDealPeriodicHealing(sim, target, baseHealing, dot.Spell.OutcomeHealingCrit)
			},
		}, core.DotConfig{}),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if isHeal {
				spell.SpellMetrics[target.UnitIndex].Hits--
				hot := spell.Hot(target)
				hot.Apply(sim)
				// Do immediate tick
				hot.TickOnce(sim)
			} else {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.SpellMetrics[target.UnitIndex].Hits--
					dot := spell.Dot(target)
					dot.Apply(sim)
					// Do immediate tick
					dot.TickOnce(sim)
				}
				spell.DealOutcome(sim, result)
			}
		},
	})
}
