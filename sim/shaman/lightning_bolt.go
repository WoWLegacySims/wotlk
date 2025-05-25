package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) registerLightningBoltSpell() {
	dbc := shamaninfo.LightningBolt.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	shaman.LightningBolt = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(false, dbc))
	shaman.LightningBoltLO = shaman.RegisterSpell(shaman.newLightningBoltSpellConfig(true, dbc))
}

func (shaman *Shaman) newLightningBoltSpellConfig(isLightningOverload bool, dbc *spellinfo.SpellInfo) core.SpellConfig {
	bp, die := dbc.GetBPDie(0, shaman.Level)
	coef := (dbc.GetCoefficient(0)) * dbc.GetLevelPenalty(shaman.Level)

	spellConfig := shaman.newElectricSpellConfig(
		core.ActionID{SpellID: dbc.SpellID},
		(dbc.BaseCost/100)*core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthShatterGarb, 2), 0.95, 1),
		dbc.CastTime,
		isLightningOverload,
		shamaninfo.LightningBolt.GetAllIDs())

	switch shaman.Ranged().ID {
	case 28066:
		spellConfig.ManaCost.FlatModifier -= 15
	}

	if shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLightningBolt) {
		spellConfig.DamageMultiplier += 0.04
	}

	if shaman.HasSetBonus(ItemSetSkyshatterRegalia, 4) {
		spellConfig.DamageMultiplier += 0.05
	}

	var lbDotSpell *core.Spell
	if !isLightningOverload && shaman.HasSetBonus(ItemSetWorldbreakerGarb, 4) {
		lbDotSpell = shaman.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 64930},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			Flags:            core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "Electrified",
				},
				TickLength:    time.Second * 2,
				NumberOfTicks: 2,

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				spell.Dot(target).ApplyOrReset(sim)
			},
		})
	}

	dmgBonus := shaman.electricSpellBonusDamage(coef)
	spellCoeff := coef + 0.04*float64(shaman.Talents.Shamanism)*dbc.GetLevelPenalty(shaman.Level)

	canLO := !isLightningOverload && shaman.Talents.LightningOverload > 0
	lightningOverloadChance := float64(shaman.Talents.LightningOverload) * 0.11
	spellConfig.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		baseDamage := dmgBonus + sim.Roll(bp, die) + spellCoeff*spell.SpellPower()
		result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

		if !isLightningOverload && lbDotSpell != nil && result.DidCrit() {
			lbDot := lbDotSpell.Dot(target)

			newDamage := result.Damage * 0.08
			outstandingDamage := core.TernaryFloat64(lbDot.IsActive(), lbDot.SnapshotBaseDamage*float64(lbDot.NumberOfTicks-lbDot.TickCount), 0)

			lbDot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(lbDot.NumberOfTicks)
			lbDot.SnapshotAttackerMultiplier = 1
			lbDotSpell.Cast(sim, target)
		}

		if canLO && result.Landed() && sim.RandomFloat("LB Lightning Overload") < lightningOverloadChance {
			shaman.LightningBoltLO.Cast(sim, target)
		}

		spell.DealDamage(sim, result)
	}

	return spellConfig
}
