package deathknight

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/deathknightinfo"
)

func (dk *Deathknight) newFrostStrikeHitSpell(isMH bool) *core.Spell {
	dbc := deathknightinfo.FrostStrike.GetMaxRank(dk.Level)
	if dbc == nil {
		return nil
	}
	damage := dbc.Effects[0].BasePoints + 1

	baseActionID := core.ActionID{SpellID: dbc.SpellID}

	bonusBaseDamage := dk.sigilOfTheVengefulHeartFrostStrike()

	actionID := baseActionID.WithTag(1)
	if !isMH {
		actionID = baseActionID.WithTag(2)
	}

	conf := core.SpellConfig{
		ActionID:    actionID,
		SpellRanks:  deathknightinfo.FrostStrike.GetAllIDs(),
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfFrostStrike), 32, 40),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCrit: (dk.annihilationCritBonus() + dk.darkrunedBattlegearCritBonus()),
		DamageMultiplier: .55 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			dk.bloodOfTheNorthCoeff(),
		CritMultiplier:   dk.bonusCritMultiplier(dk.Talents.GuileOfGorefiend),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = damage +
					bonusBaseDamage +
					spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				// SpellID 66962
				baseDamage = damage/2 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.glacielRotBonus(target) *
				dk.RoRTSBonus(target) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCost(sim, result)
				dk.threatOfThassarianProc(sim, result, dk.FrostStrikeOhHit)
			}

			spell.DealDamage(sim, result)
		},
	}

	if !isMH {
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	} else {
		conf.Flags |= core.SpellFlagAPL
	}

	return dk.RegisterSpell(conf)
}

func (dk *Deathknight) registerFrostStrikeSpell() {
	if !dk.Talents.FrostStrike {
		return
	}

	dk.FrostStrikeMhHit = dk.newFrostStrikeHitSpell(true)
	dk.FrostStrikeOhHit = dk.newFrostStrikeHitSpell(false)
	dk.FrostStrike = dk.FrostStrikeMhHit
}
