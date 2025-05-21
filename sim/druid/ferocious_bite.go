package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

func (druid *Druid) registerFerociousBiteSpell() {
	dbc := druidinfo.FerociousBite.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}

	dmgPerComboPoint := dbc.Effects[0].PointsPerCombo + core.TernaryFloat64(druid.Ranged().ID == 25667, 14, 0)
	bp := dbc.Effects[0].BasePoints
	die := dbc.Effects[0].Die
	dmgPerEnergy := dbc.Effects[0].ChainAmplitude

	druid.FerociousBite = druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost:          35,
			Refund:        0.4 * float64(druid.Talents.PrimalPrecision),
			RefundMetrics: druid.PrimalPrecisionRecoveryMetrics,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return druid.ComboPoints() > 0
		},

		BonusCrit: 0 +
			core.TernaryFloat64(druid.HasSetBonus(ItemSetMalfurionsBattlegear, 4), 5, 0.0) +
			core.TernaryFloat64(druid.AssumeBleedActive, 5*float64(druid.Talents.RendAndTear), 0),
		DamageMultiplier: (1 + 0.03*float64(druid.Talents.FeralAggression)) *
			core.TernaryFloat64(druid.HasSetBonus(ItemSetThunderheartHarness, 4), 1.15, 1.0),
		CritMultiplier:   druid.MeleeCritMultiplier(Cat),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			comboPoints := float64(druid.ComboPoints())
			attackPower := spell.MeleeAttackPower()
			excessEnergy := min(druid.CurrentEnergy(), 30)

			baseDamage := sim.Roll(bp, die) +
				dmgPerComboPoint*comboPoints +
				excessEnergy*(dmgPerEnergy+attackPower/410) +
				attackPower*0.07*comboPoints

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				druid.SpendEnergy(sim, excessEnergy, spell.Cost.(*core.EnergyCost).ResourceMetrics)
				druid.SpendComboPoints(sim, spell.ComboPointMetrics())
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}

func (druid *Druid) CurrentFerociousBiteCost() float64 {
	return druid.FerociousBite.ApplyCostModifiers(druid.FerociousBite.DefaultCast.Cost)
}
