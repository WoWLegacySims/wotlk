package tbc

import (
	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	helpers.NewWeaponDamageProc(28164, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Tranquillien Flamberge",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  35,
	})

	core.NewItemEffect(28311, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(28311)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 34107},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(104, 21)
				res := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if res.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, res.Damage, spell.OutcomeHealing)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Revenger",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.9,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})
}
