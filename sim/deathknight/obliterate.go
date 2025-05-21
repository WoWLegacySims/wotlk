package deathknight

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/deathknightinfo"
)

// TODO: Cleanup obliterate the same way we did for plague strike
func (dk *Deathknight) newObliterateHitSpell(isMH bool) *core.Spell {
	dbc := deathknightinfo.Obliterate.GetMaxRank(dk.Level)
	if dbc == nil {
		return nil
	}
	damage := dbc.Effects[0].BasePoints + 1
	actionID := core.ActionID{SpellID: dbc.SpellID}

	bonusBaseDamage := dk.sigilOfAwarenessBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.125)
	diseaseConsumptionChance := []float64{1.0, 0.67, 0.34, 0.0}[dk.Talents.Annihilation]
	deathConvertChance := float64(dk.Talents.DeathRuneMastery) / 3

	conf := core.SpellConfig{
		ActionID:    actionID.WithTag(core.TernaryInt32(isMH, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    dk.threatOfThassarianProcMask(isMH),
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			FrostRuneCost:  1,
			UnholyRuneCost: 1,
			RunicPowerGain: 15 + 2.5*float64(dk.Talents.ChillOfTheGrave) + dk.scourgeborneBattlegearRunicPowerBonus(),
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCrit: (dk.rimeCritBonus() + dk.subversionCritBonus() + dk.annihilationCritBonus() + dk.scourgeborneBattlegearCritBonus()),
		DamageMultiplier: .8 *
			core.TernaryFloat64(isMH, 1, dk.nervesOfColdSteelBonus()) *
			core.TernaryFloat64(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfObliterate), 1.25, 1.0) *
			dk.scourgelordsBattlegearDamageBonus(ScourgelordBonusSpellOB),
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
				// SpellID 66974
				baseDamage = damage/2 +
					bonusBaseDamage +
					spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}
			baseDamage *= dk.RoRTSBonus(target) *
				(1.0 + dk.dkCountActiveDiseases(target)*diseaseMulti) *
				dk.mercilessCombatBonus(sim)

			result := spell.CalcDamage(sim, target, baseDamage, dk.threatOfThassarianOutcomeApplier(spell))

			if isMH {
				spell.SpendRefundableCostAndConvertFrostOrUnholyRune(sim, result, deathConvertChance)
				dk.threatOfThassarianProc(sim, result, dk.ObliterateOhHit)

				if sim.RandomFloat("Annihilation") < diseaseConsumptionChance {
					dk.FrostFeverSpell.Dot(target).Deactivate(sim)
					dk.BloodPlagueSpell.Dot(target).Deactivate(sim)
				}

				if sim.RandomFloat("Rime") < dk.rimeHbChanceProc() {
					dk.FreezingFogAura.Activate(sim)
				}
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

func (dk *Deathknight) registerObliterateSpell() {
	dk.ObliterateMhHit = dk.newObliterateHitSpell(true)
	dk.ObliterateOhHit = dk.newObliterateHitSpell(false)
	dk.Obliterate = dk.ObliterateMhHit
}
