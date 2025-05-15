package deathknight

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/deathknightinfo"
)

func (dk *Deathknight) registerBloodBoilSpell() {
	dbc := deathknightinfo.BloodBoil.FindMaxRank(dk.Level)
	if dbc == nil {
		return
	}
	minDamage := dbc.Effects[0].BasePoints
	maxDamage := minDamage + dbc.Effects[0].Die
	minDamage += 1
	actionID := core.ActionID{SpellID: dbc.SpellID}

	// TODO: Handle blood boil correctly -
	//  There is no refund and you only get RP on at least one of the effects hitting.
	dk.BloodBoil = dk.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: dk.bloodyStrikesBonus(BloodyStrikesBB),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1.0,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := (sim.Roll(minDamage, maxDamage) + 0.06*dk.getImpurityBonus(spell)) * dk.RoRTSBonus(aoeTarget) * core.TernaryFloat64(dk.DiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

				if aoeTarget == target {
					spell.SpendRefundableCost(sim, result)
				}
			}
		},
	})
}

func (dk *Deathknight) registerDrwBloodBoilSpell() {
	dbc := deathknightinfo.BloodBoil.FindMaxRank(dk.Level)
	if dbc == nil {
		return
	}
	minDamage := dbc.Effects[0].BasePoints
	maxDamage := minDamage + dbc.Effects[0].Die
	minDamage += 1
	actionID := core.ActionID{SpellID: dbc.SpellID}

	dk.RuneWeapon.BloodBoil = dk.RuneWeapon.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: dk.bloodyStrikesBonus(BloodyStrikesBB),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.MightOfMograine),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := (sim.Roll(minDamage, maxDamage) + 0.06*dk.RuneWeapon.getImpurityBonus(spell)) * core.TernaryFloat64(dk.DrwDiseasesAreActive(aoeTarget), 1.5, 1.0)
				baseDamage *= sim.Encounter.AOECapMultiplier()

				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}
