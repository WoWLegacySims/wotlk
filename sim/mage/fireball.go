package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

var SpellRanks = core.SpellRanks{Ranks: []*core.SpellRank{
	{
		ActionID:     core.ActionID{SpellID: 133},
		MinLevel:     1,
		MinEffect:    14,
		MaxEffect:    22,
		LevelScaling: 0.6,
		OverTime:     2,
		Ticks:        2,
		BaseCost:     0.08,
		CastTime:     1500 * time.Millisecond,
		Coefficient:  0.123,
	},
	{
		ActionID:     core.ActionID{SpellID: 143},
		MinLevel:     6,
		MinEffect:    31,
		MaxEffect:    45,
		LevelScaling: 0.8,
		OverTime:     3,
		Ticks:        3,
		BaseCost:     0.11,
		CastTime:     2000 * time.Millisecond,
		Coefficient:  0.271,
	},
	{
		ActionID:     core.ActionID{SpellID: 145},
		MinLevel:     12,
		MinEffect:    53,
		MaxEffect:    73,
		LevelScaling: 1,
		OverTime:     6,
		Ticks:        3,
		BaseCost:     0.14,
		CastTime:     2500 * time.Millisecond,
		Coefficient:  0.5,
	},
	{
		ActionID:     core.ActionID{SpellID: 3140},
		MinLevel:     18,
		MinEffect:    84,
		MaxEffect:    116,
		LevelScaling: 1.3,
		OverTime:     12,
		Ticks:        4,
		BaseCost:     0.16,
		CastTime:     3000 * time.Millisecond,
		Coefficient:  0.793,
	},
	{
		ActionID:     core.ActionID{SpellID: 8400},
		MinLevel:     24,
		MinEffect:    139,
		MaxEffect:    187,
		LevelScaling: 1.8,
		OverTime:     20,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 8401},
		MinLevel:     30,
		MinEffect:    199,
		MaxEffect:    265,
		LevelScaling: 2.1,
		OverTime:     28,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 8402},
		MinLevel:     36,
		MinEffect:    255,
		MaxEffect:    335,
		LevelScaling: 2.4,
		OverTime:     32,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 10148},
		MinLevel:     42,
		MinEffect:    318,
		MaxEffect:    414,
		LevelScaling: 2.7,
		OverTime:     40,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 10149},
		MinLevel:     48,
		MinEffect:    392,
		MaxEffect:    506,
		LevelScaling: 3,
		OverTime:     52,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 10150},
		MinLevel:     54,
		MinEffect:    475,
		MaxEffect:    609,
		LevelScaling: 3.4,
		OverTime:     60,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 10151},
		MinLevel:     60,
		MinEffect:    561,
		MaxEffect:    715,
		LevelScaling: 3.7,
		OverTime:     72,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 25306},
		MinLevel:     60,
		MinEffect:    596,
		MaxEffect:    760,
		LevelScaling: 3.8,
		OverTime:     76,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 27070},
		MinLevel:     66,
		MinEffect:    633,
		MaxEffect:    805,
		LevelScaling: 4,
		OverTime:     84,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 38692},
		MinLevel:     70,
		MinEffect:    717,
		MaxEffect:    913,
		LevelScaling: 4.2,
		OverTime:     92,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 42832},
		MinLevel:     74,
		MinEffect:    783,
		MaxEffect:    997,
		LevelScaling: 4.6,
		OverTime:     100,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
	{
		ActionID:     core.ActionID{SpellID: 42833},
		MinLevel:     78,
		MinEffect:    888,
		MaxEffect:    1132,
		LevelScaling: 5.2,
		OverTime:     116,
		Ticks:        4,
		BaseCost:     0.19,
		CastTime:     3500 * time.Millisecond,
		Coefficient:  1,
	},
}}

func (mage *Mage) registerFireballSpell() {
	rank := SpellRanks.FindMaxRank(mage.Level)

	if rank == nil {
		return
	}

	spellCoeff := rank.Coefficient + 0.05*float64(mage.Talents.EmpoweredFire)
	levelDamage := float64(min(mage.Level, rank.MinLevel+4)-rank.MinLevel) * rank.LevelScaling

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFireball)

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:     rank.ActionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | HotStreakSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: rank.BaseCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				CastTime: rank.CastTime -
					time.Millisecond*100*time.Duration(mage.Talents.ImprovedFireball) -
					core.TernaryDuration(hasGlyph, time.Millisecond*150, 0),
			},
		},

		BonusCrit: 0 +
			2*float64(mage.Talents.CriticalMass) +
			float64(mage.Talents.ImprovedScorch) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5, 0),
		DamageMultiplier: 1 *
			(1 + .04*float64(mage.Talents.TormentTheWeak)),
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.SpellImpact) +
			.02*float64(mage.Talents.FirePower) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Fireball",
			},
			NumberOfTicks: rank.Ticks,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = rank.OverTime / 4.0
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(rank.MinEffect, rank.MaxEffect) + levelDamage + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() && !hasGlyph {
					spell.Dot(target).Apply(sim)
				}
				spell.DealDamage(sim, result)
			})
		},
	})
}
