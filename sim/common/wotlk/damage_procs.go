package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false

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
		School: core.SpellSchoolFire,
		MinDmg: 1024,
		MaxDmg: 1536,
	})

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
		School: core.SpellSchoolShadow,
		MinDmg: 1168,
		MaxDmg: 1752,
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
		School: core.SpellSchoolHoly,
		MinDmg: 1024,
		MaxDmg: 1536,
	})

	core.AddEffectsToTest = true

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
		School: core.SpellSchoolArcane,
		MinDmg: 1504,
		MaxDmg: 2256,
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
		School: core.SpellSchoolShadow,
		MinDmg: 788,
		MaxDmg: 1312,
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
		School: core.SpellSchoolShadow,
		MinDmg: 1750,
		MaxDmg: 2250,
	})
}
