package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func (mage *Mage) registerFireballSpell() {
	dbc := core.FindMaxRank(FireballInfos, mage.Level)

	if dbc == nil {
		return
	}

	spellCoeff := dbc.Effects[0].Coefficient + 0.05*float64(mage.Talents.EmpoweredFire)
	levelDamage := float64(min(mage.Level, dbc.MaxLevel)-dbc.MinLevel) * dbc.Effects[0].LevelScaling

	hasGlyph := mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFireball)

	mage.Fireball = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: dbc.SpellID},
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | HotStreakSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: dbc.BaseCost / 100,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				CastTime: dbc.CastTime -
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
			NumberOfTicks: dbc.Duration / dbc.Period,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = dbc.Effects[1].BasePoints + dbc.Effects[1].Die
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(dbc.Effects[0].BasePoints+1, dbc.Effects[0].BasePoints+dbc.Effects[0].Die) + levelDamage + spellCoeff*spell.SpellPower()
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
