package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	core.NewItemEffect(24114, func(agent core.Agent) { // Braided Eternium Chain
		agent.GetCharacter().PseudoStats.BonusDamage += 5
	})

	core.NewItemEffect(29996, func(agent core.Agent) { // Rod of the Sun King
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(29996)
		pppm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		actionID := core.ActionID{ItemID: 29996}

		var resourceMetricsRage *core.ResourceMetrics
		var resourceMetricsEnergy *core.ResourceMetrics
		if character.HasRageBar() {
			resourceMetricsRage = character.NewRageMetrics(actionID)
		}
		if character.HasEnergyBar() {
			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rod of the Sun King",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if pppm.Proc(sim, spell.ProcMask, "Rod of the Sun King") {
					switch spell.Unit.GetCurrentPowerBar() {
					case core.RageBar:
						spell.Unit.AddRage(sim, 5, resourceMetricsRage)
					case core.EnergyBar:
						spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
					}
				}
			},
		})
	})

	core.NewItemEffect(30892, func(agent core.Agent) { //Beast-tamer's Shoulders
		for _, pet := range agent.GetCharacter().Pets {
			if pet.IsGuardian() {
				continue // not sure if this applies to guardians.
			}
			pet.PseudoStats.DamageDealtMultiplier *= 1.03
			pet.AddStat(stats.MeleeCrit, pet.CritRatingPerCritChance*2)
		}
	})

	core.NewItemEffect(31193, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(31193)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 24585},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(48, 54) + spell.SpellPower()
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Blade of Unquenched Thirst Trigger",
			ActionID:   core.ActionID{ItemID: 31193},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.02,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(32262, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(32262)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 40291},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeMagicHitAndCrit)
			},
		})

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Siphon Essence",
			ActionID: core.ActionID{SpellID: 40291},
			Duration: time.Second * 6,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, result.Target)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Syphon of the Nathrezim",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Syphon Of The Nathrezim") {
					procAura.Activate(sim)
				}
			},
		})
	})

}
