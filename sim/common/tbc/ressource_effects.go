package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(31322, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(31322)
		metrics := character.NewManaMetrics(core.ActionID{SpellID: 38284})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "The Hammer of Destiniy",
			ActionID: core.ActionID{ItemID: 31322},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procmask,
			Outcome:  core.OutcomeLanded,
			PPM:      2.5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(169, 61)
				character.AddMana(sim, amount, metrics)
			},
		})
	})

	core.NewItemEffect(34430, func(a core.Agent) {
		character := a.GetCharacter()
		metrics := character.NewManaMetrics(core.ActionID{ItemID: 34430})

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeMana,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 34430},
				ProcMask: core.ProcMaskEmpty,
				Flags:    core.SpellFlagChanneled,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 5,
					},
				},
				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "Evocation",
					},
					SelfOnly:            true,
					NumberOfTicks:       8,
					TickLength:          time.Second * 1,
					AffectedByCastSpeed: true,
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						character.AddMana(sim, 250, metrics)
					},
				},
			}),
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 2250
			},
		})
	})

	core.NewItemEffect(35703, func(a core.Agent) {
		character := a.GetCharacter()
		metrics := character.NewManaMetrics(core.ActionID{ItemID: 35703})

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeMana,
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 900
			},
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 35703},
				ProcMask: core.ProcMaskEmpty,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 3,
					},
				},
				Hot: core.DotConfig{
					Aura: core.Aura{
						Label: "Evocation",
					},
					SelfOnly:      true,
					NumberOfTicks: 12,
					TickLength:    time.Second,
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						character.AddMana(sim, 75, metrics)
					},
				},
			}),
		})
	})
	core.AddEffectsToTest = true
}
