package vanilla

import (
	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 12631,
		Trigger: core.ProcTrigger{
			Name:       "Fiery Plate Gauntlets",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ProcChance: 1.0,
		},
		School: core.SpellSchoolFire,
		MinDmg: 4,
		MaxDmg: 4,
	})
	core.AddEffectsToTest = true
}
