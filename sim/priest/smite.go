package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) RegisterSmiteSpell() {
	dbc := priestinfo.Smite.GetMaxRank(priest.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	coef := 0.1833 * dbc.GetLevelPenalty(priest.Level)

	priest.Smite = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  priestinfo.Smite.GetAllIDs(),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.15,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2500 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
		},

		BonusCrit:        float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
