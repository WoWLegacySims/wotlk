package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(9380, func(a core.Agent) {
		character := a.GetCharacter()
		shieldStrength := 0.0
		procmask := character.GetProcMaskForItem(9380)
		metrics := character.NewHealthMetrics(core.ActionID{SpellID: 10618})

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: 9380},
			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Jang'thraze the Protector",
					ActionID: core.ActionID{SpellID: 10618},
					Duration: time.Second * 20,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						shieldStrength = sim.Roll(54, 31)
					},
					OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if spell.SpellSchool.Matches(core.SpellSchoolPhysical) && result.Damage > 0 {
							absorb := min(result.Damage, shieldStrength)
							shieldStrength -= absorb
							character.GainHealth(sim, absorb, metrics)

							if shieldStrength == 0 {
								aura.Deactivate(sim)
							}
						}
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
			ActionID:   core.ActionID{ItemID: 9380},
			ProcMask:   procmask,
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.04,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})
	core.AddEffectsToTest = true
}
