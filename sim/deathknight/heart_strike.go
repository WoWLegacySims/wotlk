package deathknight

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/deathknightinfo"
)

func (dk *Deathknight) newHeartStrikeSpell(isMainTarget bool, isDrw bool) *core.Spell {
	dbc := deathknightinfo.HeartStrikeInfos.FindMaxRank(dk.Level)
	if dbc == nil {
		return nil
	}
	damage := dbc.Effects[0].BasePoints + 1

	actionID := core.ActionID{SpellID: dbc.SpellID}
	bonusBaseDamage := dk.sigilOfTheDarkRiderBonus()
	diseaseMulti := dk.dkDiseaseMultiplier(0.1)

	critMultiplier := dk.bonusCritMultiplier(dk.Talents.MightOfMograine)

	conf := core.SpellConfig{
		ActionID:    actionID.WithTag(core.TernaryInt32(isMainTarget, 1, 2)),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		RuneCost: core.RuneCostOptions{
			BloodRuneCost:  1,
			RunicPowerGain: 10,
			Refundable:     true,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		BonusCrit: (dk.subversionCritBonus() + dk.annihilationCritBonus()),
		DamageMultiplier: .5 *
			core.TernaryFloat64(isMainTarget, 1, 0.5) *
			dk.thassariansPlateDamageBonus() *
			dk.scourgelordsBattlegearDamageBonus(ScourgelordBonusSpellHS) *
			dk.bloodyStrikesBonus(BloodyStrikesHS),
		CritMultiplier:   critMultiplier,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := damage + bonusBaseDamage

			if isDrw {
				baseDamage += dk.DrwWeaponDamage(sim, spell)
			} else {
				baseDamage += spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}

			activeDiseases := core.TernaryFloat64(isDrw, dk.drwCountActiveDiseases(target), dk.dkCountActiveDiseases(target))
			baseDamage *= 1 + activeDiseases*diseaseMulti

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if isMainTarget {
				if isDrw {
					if dk.Env.GetNumTargets() > 1 {
						dk.RuneWeapon.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(target))
					}
				} else {
					spell.SpendRefundableCost(sim, result)

					if dk.Env.GetNumTargets() > 1 {
						dk.HeartStrikeOffHit.Cast(sim, dk.Env.NextTargetUnit(target))
					}
				}
			}
		},
	}
	if !isMainTarget || isDrw { // off target doesnt need GCD
		conf.RuneCost = core.RuneCostOptions{}
		conf.Cast = core.CastConfig{}
	}

	if isMainTarget {
		conf.Flags |= core.SpellFlagAPL
	}

	if isDrw {
		return dk.RuneWeapon.RegisterSpell(conf)
	} else {
		return dk.RegisterSpell(conf)
	}
}

func (dk *Deathknight) registerHeartStrikeSpell() {
	if !dk.Talents.HeartStrike {
		return
	}

	dk.HeartStrike = dk.newHeartStrikeSpell(true, false)
	dk.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, false)
}

func (dk *Deathknight) registerDrwHeartStrikeSpell() {
	if !dk.Talents.HeartStrike {
		return
	}

	dk.RuneWeapon.HeartStrike = dk.newHeartStrikeSpell(true, true)
	dk.RuneWeapon.HeartStrikeOffHit = dk.newHeartStrikeSpell(false, true)
}
