package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	core.NewItemEffect(2802, func(a core.Agent) {
		character := a.GetCharacter()

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Blazing Emblem",
			ActionID: core.ActionID{ItemID: 2802},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FireResistance, 50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.FireResistance, -50)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					amount := min(25, result.Damage)
					spell.CalcAndDealHealing(sim, &character.Unit, amount, spell.OutcomeHealing)
				}
			},
			OnPeriodicDamageTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.OnSpellHitTaken(aura, sim, spell, result)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 2802},
			SpellSchool: core.SpellSchoolFire,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 10,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	core.NewItemEffect(8367, func(a core.Agent) {
		character := a.GetCharacter()
		shieldStrength := 0.0

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 8367},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Hour,
						Timer:    character.NewTimer(),
					},
				},
				Shield: core.ShieldConfig{
					Aura: core.Aura{
						Label:    "Elemental Protection",
						ActionID: core.ActionID{SpellID: 10618},
						Duration: time.Minute * 2,
						OnGain: func(aura *core.Aura, sim *core.Simulation) {
							shieldStrength = 600
						},
						OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
							if spell.SpellSchool.Matches(core.SpellSchoolMagic) && result.Damage > 0 {
								absorb := min(result.Damage, shieldStrength)
								shieldStrength -= absorb
								spell.CalcAndDealHealing(sim, &character.Unit, absorb, spell.OutcomeHealing)
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
			}),
		})
	})
}
