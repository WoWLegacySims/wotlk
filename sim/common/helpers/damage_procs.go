package helpers

import "github.com/WoWLegacySims/wotlk/sim/core"

type ProcDamageEffect struct {
	ID      int32
	Trigger core.ProcTrigger

	School     core.SpellSchool
	BasePoints float64
	Die        float64
}

func NewProcDamageEffect(config ProcDamageEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		bp := config.BasePoints
		die := config.Die
		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(bp, die), spell.OutcomeMagicHitAndCrit)
			},
		})

		triggerConfig := config.Trigger
		triggerConfig.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			damageSpell.Cast(sim, character.CurrentTarget)
		}
		core.MakeProcTriggerAura(&character.Unit, triggerConfig)
	})
}

type ProcHealEffect struct {
	ID      int32
	Trigger core.ProcTrigger

	School     core.SpellSchool
	BasePoints float64
	Die        float64
}

func NewProcHealEffect(config ProcHealEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		bp := config.BasePoints
		die := config.Die
		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealHealing(sim, &character.Unit, sim.Roll(bp, die), spell.OutcomeMagicHitAndCrit)
			},
		})

		triggerConfig := config.Trigger
		triggerConfig.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			damageSpell.Cast(sim, character.CurrentTarget)
		}
		core.MakeProcTriggerAura(&character.Unit, triggerConfig)
	})
}
