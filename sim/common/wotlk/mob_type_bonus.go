package wotlk

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func init() {
	core.AddEffectsToTest = false

	core.NewItemEffect(37018, func(a core.Agent) {
		character := a.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 40
		}
	})

	core.AddEffectsToTest = true
}
