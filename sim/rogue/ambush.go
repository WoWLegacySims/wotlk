package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) registerAmbushSpell() {
	dbc := rogueinfo.Ambush.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, rogue.Level)

	rogue.Ambush = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  rogueinfo.Ambush.GetAllIDs(),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | SpellFlagBuilder | SpellFlagColdBlooded | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60 - 4*float64(rogue.Talents.SlaughterFromTheShadows)),
			Refund: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !rogue.PseudoStats.InFrontOfTarget && rogue.HasDagger(core.MainHand) && rogue.IsStealthed()
		},

		BonusCrit: float64(rogue.Talents.TurnTheTables)*2 +
			25*float64(rogue.Talents.ImprovedAmbush),
		// All of these use "Apply Aura: Modifies Damage/Healing Done", and stack additively.
		DamageMultiplier: 2.75 * (1 +
			0.02*float64(rogue.Talents.FindWeakness) +
			0.1*float64(rogue.Talents.Opportunity)),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			baseDamage := bp +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
