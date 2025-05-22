package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) registerBackstabSpell() {
	dbc := rogueinfo.Backstab.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, rogue.Level)
	hasGlyph := rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfBackstab)

	rogue.Backstab = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  rogueinfo.Backstab.GetAllIDs(),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60 - 4*float64(rogue.Talents.SlaughterFromTheShadows)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand)
		},

		BonusCrit: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5, 0) +
			float64(rogue.Talents.TurnTheTables)*2 +
			10*float64(rogue.Talents.PuncturingWounds),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively (up to 142%).
		DamageMultiplier: 1.5 * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.Opportunity) +
			0.03*float64(rogue.Talents.Aggression) +
			0.05*float64(rogue.Talents.BladeTwisting) +
			core.TernaryFloat64(rogue.Talents.SurpriseAttacks, 0.1, 0) +
			core.TernaryFloat64(rogue.HasSetBonus(Tier6, 4), 0.06, 0)) *
			(1 + 0.02*float64(rogue.Talents.SinisterCalling)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := bp +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 1, spell.ComboPointMetrics())
				if dot := rogue.Rupture.Dot(target); hasGlyph && dot.IsActive() && dot.NumberOfTicks < dot.MaxStacks+3 {
					dot.NumberOfTicks += 1
					dot.RecomputeAuraDuration()
					dot.UpdateExpires(dot.ExpiresAt() + dot.TickLength)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
