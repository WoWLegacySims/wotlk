package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 37064,
		Trigger: core.ProcTrigger{
			Name:       "Vestige of Haldor",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
			ActionID:   core.ActionID{ItemID: 37064},
		},
		School:     core.SpellSchoolFire,
		BasePoints: 1023,
		Die:        513,
	})
	core.AddEffectsToTest = false
	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 37264,
		Trigger: core.ProcTrigger{
			Name:       "Pendulum of Telluric Currents",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellOrProc,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
			ActionID:   core.ActionID{ItemID: 37264},
		},
		School:     core.SpellSchoolShadow,
		BasePoints: 1167,
		Die:        585,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 39889,
		Trigger: core.ProcTrigger{
			Name:       "Horn of Agent Fury",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
			ActionID:   core.ActionID{ItemID: 39889},
		},
		School:     core.SpellSchoolHoly,
		BasePoints: 1023,
		Die:        523,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 40371,
		Trigger: core.ProcTrigger{
			Name:       "Bandit's Insignia",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
			ActionID:   core.ActionID{ItemID: 40371},
		},
		School:     core.SpellSchoolArcane,
		BasePoints: 1503,
		Die:        753,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 40373,
		Trigger: core.ProcTrigger{
			Name:       "Extract of Necromantic Power",
			Callback:   core.CallbackOnPeriodicDamageDealt,
			Harmful:    true,
			ProcChance: 0.10,
			ICD:        time.Second * 15,
			ActionID:   core.ActionID{ItemID: 40373},
		},
		School:     core.SpellSchoolShadow,
		BasePoints: 787,
		Die:        525,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 42990,
		Trigger: core.ProcTrigger{
			Name:       "DMC Death",
			Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
			Harmful:    true,
			ProcChance: 0.15,
			ICD:        time.Second * 45,
			ActionID:   core.ActionID{ItemID: 42990},
		},
		School:     core.SpellSchoolShadow,
		BasePoints: 1749,
		Die:        501,
	})

	core.NewItemEffect(49302, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(49302)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 21170},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(99, 81)
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, result.Damage, spell.OutcomeHealing)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Reclaimed Shadowstrike",
			ActionID:   core.ActionID{ItemID: 49302},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procmask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.0558,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(49301, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(49301)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 21179},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				bounceCoeff := 1.0
				curTarget := target
				for hitIndex := int32(0); hitIndex < 3; hitIndex++ {
					baseDamage := sim.Roll(149, 101)
					baseDamage *= bounceCoeff
					spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

					bounceCoeff *= 0.7
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Reinforced Thunderstrike",
			ActionID:   core.ActionID{ItemID: 49301},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procmask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.0558,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(49496, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(49496)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 21170},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(99, 81)
				result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.CalcAndDealHealing(sim, &character.Unit, result.Damage, spell.OutcomeHealing)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Reinforced Shadowstrike",
			ActionID:   core.ActionID{ItemID: 49496},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procmask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.0558,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(49497, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(49497)

		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 21179},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				bounceCoeff := 1.0
				curTarget := target
				for hitIndex := int32(0); hitIndex < 3; hitIndex++ {
					baseDamage := sim.Roll(149, 101)
					baseDamage *= bounceCoeff
					spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMagicHitAndCrit)

					bounceCoeff *= 0.7
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Reinforced Thunderstrike",
			ActionID:   core.ActionID{ItemID: 49497},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procmask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.0558,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})

	helpers.NewWeaponDamageProc(49297, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Empowered Deathbringer",
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  1312,
		Die:         375,
	})

	helpers.NewWeaponDamageProc(49296, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Singed Vis'kag the Bloodletter",
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  2000,
	})

	helpers.NewWeaponDamageProc(49500, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Raging Deathbringer",
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  1457,
		Die:         417,
	})

	helpers.NewWeaponDamageProc(49501, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Tempered Vis'kag the Bloodletter",
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  2222,
	})

	core.AddEffectsToTest = true
}
