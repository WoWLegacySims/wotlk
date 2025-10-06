package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerArcaneBlastSpell() {
	dbc := mageinfo.ArcaneBlast.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, mage.Level)

	abAuraMultiplierPerStack := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneBlast), .18, .15)
	mage.ArcaneBlastAura = mage.GetOrRegisterAura(core.Aura{
		Label:     "Arcane Blast",
		ActionID:  core.ActionID{SpellID: 36032},
		Duration:  time.Second * 8,
		MaxStacks: 4,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			oldMultiplier := 1 + float64(oldStacks)*abAuraMultiplierPerStack
			newMultiplier := 1 + float64(newStacks)*abAuraMultiplierPerStack
			mage.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= newMultiplier / oldMultiplier
			mage.ArcaneBlast.CostMultiplier += 1.75 * float64(newStacks-oldStacks)
		},
	})

	actionID := core.ActionID{SpellID: dbc.SpellID}
	spellCoeff := (dbc.GetCoefficient(0) + .03*float64(mage.Talents.ArcaneEmpowerment)) * dbc.GetLevelPenalty(mage.Level)

	mage.ArcaneBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellRanks:  mageinfo.ArcaneBlast.GetAllIDs(),
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | BarrageSpells | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.07,
			Multiplier: (1 - .01*float64(mage.Talents.ArcaneFocus)) * core.TernaryFloat64(mage.HasSetBonus(ItemSetTirisfalRegalia, 2), 1.05, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: 2500 * time.Millisecond,
			},
		},

		BonusHit: float64(mage.Talents.ArcaneFocus),
		BonusCrit: 0 +
			float64(mage.Talents.Incineration)*2 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5, 0),
		DamageMultiplier: 1 *
			(1 + .04*float64(mage.Talents.TormentTheWeak)) *
			(core.TernaryFloat64(mage.HasSetBonus(ItemSetTirisfalRegalia, 2), 1.05, 1)),
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + spellCoeff*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			mage.ArcaneBlastAura.Activate(sim)
			mage.ArcaneBlastAura.AddStack(sim)
		},
	})
}
