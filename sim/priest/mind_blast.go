package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) registerMindBlastSpell() {
	dbc := priestinfo.MindBlast.GetMaxRank(priest.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, priest.Level)
	coef := dbc.GetCoefficient(0)
	spellCoeff := coef * (1 + 0.05*float64(priest.Talents.Misery)) * dbc.GetLevelPenalty(priest.Level)
	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	var replSrc core.ReplenishmentSource
	if priest.Talents.VampiricTouch {
		replSrc = priest.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 48160})
	}

	// From Improved Mind Blast
	mindTraumaSpell := priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48301},
		ProcMask:    core.ProcMaskProc,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNoMetrics,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			MindTraumaAura(target).Activate(sim)
		},
	})

	priest.MindBlast = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  priestinfo.MindBlast.GetAllIDs(),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.17,
			Multiplier: 1 *
				(1 - 0.05*float64(priest.Talents.FocusedMind)) *
				core.TernaryFloat64(priest.HasSetBonus(ItemSetValorous, 2), 0.9, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: time.Second*8 - time.Millisecond*500*time.Duration(priest.Talents.ImprovedMindBlast),
			},
		},

		BonusHit:  0 + float64(priest.Talents.ShadowFocus),
		BonusCrit: float64(priest.Talents.MindMelt) * 2,
		DamageMultiplier: 1 *
			(1 + 0.02*float64(priest.Talents.Darkness)) *
			core.TernaryFloat64(priest.HasSetBonus(ItemSetAbsolution, 4), 1.1, 1),
		CritMultiplier:   priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + spellCoeff*spell.SpellPower()
			baseDamage *= priest.MindBlastModifier

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if result.DidCrit() && hasGlyphOfShadow {
				priest.ShadowyInsightAura.Activate(sim)
			}
			if result.DidCrit() && priest.ImprovedSpiritTap != nil {
				priest.ImprovedSpiritTap.Activate(sim)
			}
			spell.DealDamage(sim, result)

			if priest.Talents.VampiricTouch && priest.VampiricTouch.CurDot().IsActive() {
				priest.Env.Raid.ProcReplenishment(sim, replSrc)
			}

			if priest.Talents.Shadowform && priest.Talents.ImprovedMindBlast > 0 {
				if sim.RandomFloat("Improved Mind Blast") < 0.2*float64(priest.Talents.ImprovedMindBlast) {
					mindTraumaSpell.Cast(sim, target)
				}
			}
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := bp + die/2 + spellCoeff*spell.SpellPower()
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}

func MindTraumaAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:    "Mind Trauma",
		ActionID: core.ActionID{SpellID: 48301},
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.HealingTakenMultiplier /= 0.8
		},
	})
}
