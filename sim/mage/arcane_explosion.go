package mage

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerArcaneExplosionSpell() {
	dbc := mageinfo.ArcaneExplosion.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	mage.ArcaneExplosion = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.22,
			Multiplier: core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneExplosion), .9, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHit:         float64(mage.Talents.ArcaneFocus) * 2,
		BonusCrit:        float64(mage.Talents.SpellImpact) * 2,
		DamageMultiplier: 1,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := coef * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(bp, die) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
