package vanilla

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(867, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeUndead {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 30
		}
	})

	core.NewItemEffect(1465, func(a core.Agent) {
		character := a.GetCharacter()
		if character.CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			character.PseudoStats.MobTypeAttackPower += 18
		}
	})

	core.NewItemEffect(3566, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 30
		}
	})
	core.AddEffectsToTest = true
}
