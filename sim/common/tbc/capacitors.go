package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "The Lightning Capacitor",
		ID:        28785,
		MaxStacks: 3,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2500,
			ActionID: core.ActionID{ItemID: 38072},
		},
		School: core.SpellSchoolNature,
		MinDmg: 694,
		MaxDmg: 806,
	})
}
