package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) RegisterHolyFireSpell() {
	dbc := priestinfo.HolyFire.GetMaxRank(priest.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(priest.Level)
	bpDot, _ := dbc.GetBPDie(1, priest.Level)
	coefDot := dbc.GetCoefficient(1) * dbc.GetLevelPenalty(priest.Level)

	hasGlyph := priest.HasMajorGlyph(proto.PriestMajorGlyph_GlyphOfSmite)

	priest.HolyFire = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.11,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2000 - time.Millisecond*100*time.Duration(priest.Talents.DivineFury),
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second * 10,
			},
		},

		BonusCrit:        float64(priest.Talents.HolySpecialization),
		DamageMultiplier: 1 + 0.05*float64(priest.Talents.SearingLight),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "HolyFire",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					if hasGlyph {
						priest.Smite.DamageMultiplier *= 1.2
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if hasGlyph {
						priest.Smite.DamageMultiplier /= 1.2
					}
				},
			},
			NumberOfTicks: 7,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = bpDot + coefDot*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}
