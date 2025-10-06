package rogue

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) newMutilateHitSpell(isMH bool, damage float64, id int32) *core.Spell {
	actionID := core.ActionID{SpellID: id}
	procMask := core.ProcMaskMeleeMHSpecial
	if !isMH {
		actionID = core.ActionID{SpellID: id}
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | SpellFlagBuilder | SpellFlagColdBlooded,

		BonusCrit: core.TernaryFloat64(rogue.HasSetBonus(Tier9, 4), 5, 0) +
			float64(rogue.Talents.TurnTheTables)*2 +
			5*float64(rogue.Talents.PuncturingWounds),

		DamageMultiplierAdditive: 1 +
			0.1*float64(rogue.Talents.Opportunity) +
			0.02*float64(rogue.Talents.FindWeakness) +
			core.TernaryFloat64(rogue.HasSetBonus(Tier6, 4), 0.06, 0),
		DamageMultiplier: 1 *
			core.TernaryFloat64(isMH, 1, rogue.dwsMultiplier()),
		CritMultiplier:   rogue.MeleeCritMultiplier(true),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = damage + spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			} else {
				baseDamage = damage + spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower())
			}
			// TODO: Add support for all poison effects
			if rogue.DeadlyPoison.Dot(target).IsActive() || rogue.woundPoisonDebuffAuras.Get(target).IsActive() {
				baseDamage *= 1.2
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (rogue *Rogue) registerMutilateSpell() {
	if !rogue.Talents.Mutilate {
		return
	}

	dbc := rogueinfo.Mutilate.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	dbcDamage := rogueinfo.MutilateMH.GetByID(dbc.Effects[1].TriggerSpell)
	if dbcDamage == nil {
		panic(fmt.Sprintf("No Trigger spell found for Mutilate %d", dbc.SpellID))
	}
	bp, _ := dbcDamage.GetBPDie(0, rogue.Level)

	rogue.MutilateMH = rogue.newMutilateHitSpell(true, bp, dbc.Effects[1].TriggerSpell)
	rogue.MutilateOH = rogue.newMutilateHitSpell(false, bp, dbc.Effects[2].TriggerSpell)

	rogue.Mutilate = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  rogueinfo.Mutilate.GetAllIDs(),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:   rogue.costModifier(60 - core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfMutilate), 5, 0)),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit) // Miss/Dodge/Parry/Hit
			if result.Landed() {
				rogue.AddComboPoints(sim, 2, spell.ComboPointMetrics())
				rogue.MutilateOH.Cast(sim, target)
				rogue.MutilateMH.Cast(sim, target)
			} else {
				spell.IssueRefund(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
