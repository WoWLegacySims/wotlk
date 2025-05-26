package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(28767, func(a core.Agent) {
		character := a.GetCharacter()

		useSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 37208},
			SpellSchool:      core.SpellSchoolPhysical,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(512, 55)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 28767},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 3,
						Timer:    character.NewTimer(),
					},
					IgnoreHaste: true,
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					useSpell.Cast(sim, target)
				},
			}),
			Type: core.CooldownTypeDPS,
		})
	})
	core.AddEffectsToTest = true
}
