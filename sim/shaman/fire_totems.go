package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) registerSearingTotemSpell() {
	dbc := shamaninfo.SearingTotem.GetMaxRank(shaman.Level)
	dbcEffect := shamaninfo.Attack.GetMaxRank(shaman.Level)
	if dbc == nil || dbcEffect == nil {
		return
	}
	bp, die := dbcEffect.GetBPDie(0, shaman.Level)
	coef := dbcEffect.GetCoefficient(0) * dbc.GetLevelPenalty(shaman.Level)

	shaman.SearingTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  shamaninfo.SearingTotem.GetAllIDs(),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.07,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusHit:         float64(shaman.Talents.ElementalPrecision),
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SearingTotem",
			},
			// These are the real tick values, but searing totem doesn't start its next
			// cast until the previous missile hits the target. We don't have an option
			// for target distance yet so just pretend the tick rate is lower.
			// https://wotlk.wowhead.com/spell=25530/attack
			//NumberOfTicks:        30,
			//TickLength:           time.Second * 2.2,
			NumberOfTicks: 24,
			TickLength:    time.Second * 60 / 24,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := sim.Roll(bp, die) + coef*dot.Spell.SpellPower()
				dot.Spell.CalcAndDealDamage(sim, target, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.MagmaTotem.AOEDot().Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.Dot(sim.GetTargetUnit(0)).Apply(sim)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*60 + 1
		},
	})
}

func (shaman *Shaman) registerMagmaTotemSpell() {
	dbc := shamaninfo.MagmaTotem.GetMaxRank(shaman.Level)
	dbcEffect := shamaninfo.MagmaTotemEffect.GetMaxRank(shaman.Level)
	if dbc == nil || dbcEffect == nil {
		return
	}
	bp, _ := dbcEffect.GetBPDie(0, shaman.Level)
	coef := dbcEffect.GetCoefficient(0) * dbc.GetLevelPenalty(shaman.Level)

	shaman.MagmaTotem = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  shamaninfo.MagmaTotem.GetAllIDs(),
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.27,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
		},

		BonusHit:         float64(shaman.Talents.ElementalPrecision),
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "MagmaTotem",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := bp + coef*dot.Spell.SpellPower()
				baseDamage *= sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.Spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeMagicHitAndCrit)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			shaman.SearingTotem.Dot(shaman.CurrentTarget).Cancel(sim)
			shaman.FireElemental.Disable(sim)
			spell.AOEDot().Apply(sim)
			// +1 needed because of rounding issues with totem tick time.
			shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*20 + 1
		},
	})
}
