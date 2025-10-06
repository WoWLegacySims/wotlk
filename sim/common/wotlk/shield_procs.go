package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(40865, func(a core.Agent) {
		character := a.GetCharacter()
		shieldStrength := 0.0
		metrics := character.NewHealthMetrics(core.ActionID{SpellID: 55019})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: 40865},
			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Sonic Shield",
					ActionID: core.ActionID{SpellID: 55019},
					Duration: time.Second * 12,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						shieldStrength = 100
					},
					OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if result.Damage > 0 {
							absorb := min(result.Damage, shieldStrength)
							shieldStrength -= absorb
							character.GainHealth(sim, absorb, metrics)

							if shieldStrength == 0 {
								aura.Deactivate(sim)
							}
						}
					},
					OnPeriodicDamageTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						aura.OnSpellHitTaken(aura, sim, spell, result)
					},
				},
				SelfOnly: true,
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.SelfShield().Activate(sim)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Jang'thraze the Protector",
			ActionID:   core.ActionID{ItemID: 40865},
			ProcMask:   core.ProcMaskMelee,
			Callback:   core.CallbackOnSpellHitTaken,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.5,
			ICD:        time.Minute,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})

	core.AddEffectsToTest = true
}
