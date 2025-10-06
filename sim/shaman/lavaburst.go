package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	dbc := shamaninfo.LavaBurst.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, shaman.Level)
	coef := (dbc.GetCoefficient(0) + 0.05*float64(shaman.Talents.Shamanism) + core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLava), 0.1, 0)) * dbc.GetLevelPenalty(shaman.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}
	dmgBonus := core.TernaryFloat64(shaman.Ranged().ID == VentureCoLightningRod, 121, 0) +
		core.TernaryFloat64(shaman.Ranged().ID == ThunderfallTotem, 215, 0)

	var lvbDotSpell *core.Spell
	if shaman.HasSetBonus(ItemSetThrallsRegalia, 4) {
		lvbDotSpell = shaman.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 71824},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					Label: "LavaBursted",
				},
				TickLength:    time.Second * 2,
				NumberOfTicks: 3,

				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
				spell.Dot(target).Apply(sim)
			},
		})
	}

	shaman.LavaBurst = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellRanks:  shamaninfo.LavaBurst.GetAllIDs(),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagFocusable | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Second*2 - time.Millisecond*100*time.Duration(shaman.Talents.LightningMastery),
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		BonusHit:         float64(shaman.Talents.ElementalPrecision),
		DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion) + 0.02*float64(shaman.Talents.CallOfFlame),
		CritMultiplier:   shaman.ElementalCritMultiplier([]float64{0, 0.06, 0.12, 0.24}[shaman.Talents.LavaFlows] + core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthShatterGarb, 4), 0.1, 0)),
		ThreatMultiplier: shaman.spellThreatMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dmgBonus + sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if lvbDotSpell != nil && result.Landed() {
				dot := lvbDotSpell.Dot(target)
				dot.SnapshotBaseDamage = result.Damage * 0.1 / float64(dot.NumberOfTicks)
				dot.SnapshotAttackerMultiplier = 1
				dot.Spell.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})
}
