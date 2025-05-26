package tbc

import (
	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
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
			ActionID: core.ActionID{ItemID: 28311},
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.9,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	helpers.NewWeaponDamageProc(28573, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Despair",
			PPM:  0.65,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  600,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID:         28579,
		School:     core.SpellSchoolNature,
		BasePoints: 221,
		Die:        111,
		Trigger: core.ProcTrigger{
			Name:     "Romulo'S Poison Vial",
			ActionID: core.ActionID{SpellID: 34586},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMeleeOrRanged,
			Outcome:  core.OutcomeLanded,
			PPM:      1.5,
		},
	})

	core.NewItemEffect(28774, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(28774)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 34696},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(284, 31)
				res := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if res.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, res.Damage, spell.OutcomeHealing)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Glaive of the Pit",
			ActionID: core.ActionID{ItemID: 28774},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.6,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(31858, func(agent core.Agent) {
		character := agent.GetCharacter()

		bp := 94.0
		die := 21.0
		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 31858},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(bp, die), spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Darkmoon Card: Vengeance",
			ActionID:   core.ActionID{ItemID: 31858},
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskDirect,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				damageSpell.Cast(sim, spell.Unit)
			},
		})
	})

	core.AddEffectsToTest = true
}
