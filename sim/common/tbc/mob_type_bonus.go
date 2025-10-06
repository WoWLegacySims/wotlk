package tbc

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(29398, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 39
		}
	})

	core.NewItemEffect(30787, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			a.GetCharacter().PseudoStats.MobTypeSpellPower += 185
		}
	})

	core.NewItemEffect(30788, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 93
		}
	})

	core.NewItemEffect(30789, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 150
		}
	})

	core.NewItemEffect(31745, func(a core.Agent) {
		if a.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeDemon {
			a.GetCharacter().PseudoStats.MobTypeAttackPower += 93
		}
	})
	core.AddEffectsToTest = true
}
