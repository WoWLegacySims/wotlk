package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.NewItemEffect(940, func(a core.Agent) {
		character := a.GetCharacter()

		metrics := character.NewManaMetrics(core.ActionID{ItemID: 940})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Insight",
			ActionID: core.ActionID{ItemID: 940},
			Duration: time.Second * 10,
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if spell.CurCast.Cost == 0 {
					return
				}
				amount := min(spell.CurCast.Cost, 500)
				spell.Unit.AddMana(sim, amount, metrics)
				aura.Deactivate(sim)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 940},
			SpellSchool: core.SpellSchoolPhysical,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 15,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(7133, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasRageBar()) {
			return
		}

		ragemetrics := character.NewRageMetrics(core.ActionID{ItemID: 7133})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{SpellID: 70537},
				SpellSchool: core.SpellSchoolPhysical,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 60,
						Timer:    character.NewTimer(),
					},
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					character.AddRage(sim, 30, ragemetrics)
				},
			}),
			Type: core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(7507, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasManaBar()) {
			return
		}

		metrics := character.NewManaMetrics(core.ActionID{ItemID: 7507})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: 7507},
				SpellSchool: core.SpellSchoolPhysical,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 60,
						Timer:    character.NewTimer(),
					},
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					amount := sim.Roll(139, 41)
					character.AddMana(sim, amount, metrics)
				},
			}),
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 180
			},
			Type: core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(7508, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasManaBar()) {
			return
		}

		metrics := character.NewManaMetrics(core.ActionID{ItemID: 7508})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: 7508},
				SpellSchool: core.SpellSchoolPhysical,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 60,
						Timer:    character.NewTimer(),
					},
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					amount := sim.Roll(139, 41)
					character.AddMana(sim, amount, metrics)
				},
			}),
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 180
			},
			Type: core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(7515, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasManaBar()) {
			return
		}

		metrics := character.NewManaMetrics(core.ActionID{ItemID: 7515})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: 7515},
				SpellSchool: core.SpellSchoolPhysical,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 60,
						Timer:    character.NewTimer(),
					},
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					amount := sim.Roll(399, 801)
					character.AddMana(sim, amount, metrics)
				},
			}),
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 800
			},
			Type: core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(9397, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasManaBar()) {
			return
		}

		metrics := character.NewManaMetrics(core.ActionID{ItemID: 9397})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: 9397},
				SpellSchool: core.SpellSchoolPhysical,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 60,
						Timer:    character.NewTimer(),
					},
					DefaultCast: core.Cast{
						GCD: core.GCDDefault,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					amount := sim.Roll(389, 21)
					character.AddMana(sim, amount, metrics)
				},
			}),
			ShouldActivate: func(s *core.Simulation, c *core.Character) bool {
				return (c.MaxMana() - c.CurrentMana()) > 410
			},
			Type: core.CooldownTypeMana,
		})
	})
}
