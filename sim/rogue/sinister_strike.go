package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) registerSinisterStrikeSpell() {
	dbc := rogueinfo.SinisterStrike.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, rogue.Level)

	hasGlyphOfSinisterStrike := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfSinisterStrike)

	rogue.SinisterStrike = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier([]float64{45, 42, 40}[rogue.Talents.ImprovedSinisterStrike]),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		BonusCrit: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5, 0) +
			float64(rogue.Talents.TurnTheTables)*2,
		DamageMultiplier: 1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.03*float64(rogue.Talents.Aggression) +
			0.05*float64(rogue.Talents.BladeTwisting) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(Tier6, 4), 0.06, 0),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := bp +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				points := int32(1)
				if hasGlyphOfSinisterStrike && result.DidCrit() {
					if sim.RandomFloat("Glyph of Sinister Strike") < 0.5 {
						points += 1
					}
				}
				rogue.AddComboPoints(sim, points, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
