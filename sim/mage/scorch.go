package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerScorchSpell() {
	dbc := mageinfo.Scorch.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	hasImpScorch := mage.Talents.ImprovedScorch > 0
	procChance := float64(mage.Talents.ImprovedScorch) / 3.0

	if hasImpScorch {
		mage.ScorchAuras = mage.NewEnemyAuraArray(core.ImprovedScorchAura)
		mage.CritDebuffCategories = mage.GetEnemyExclusiveCategories(core.SpellCritEffectCategory)
	}

	mage.Scorch = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		BonusCrit: 0 +
			float64(mage.Talents.Incineration+mage.Talents.CriticalMass)*2 +
			float64(mage.Talents.ImprovedScorch)*1,
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower) +
			core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfScorch), 0.2, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if hasImpScorch && result.Landed() && sim.Proc(procChance, "Improved Scorch") {
				mage.ScorchAuras.Get(target).Activate(sim)
			}
			spell.DealDamage(sim, result)
		},
	})

	if hasImpScorch {
		mage.Scorch.RelatedAuras = append(mage.Scorch.RelatedAuras, mage.ScorchAuras)
	}
}
