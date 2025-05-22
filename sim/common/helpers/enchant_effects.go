package helpers

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func AddWeaponDamageEnchant(id int32, amount float64) {
	core.AddWeaponEffect(id, func(agent core.Agent, slot proto.ItemSlot) {
		w := core.Ternary(slot == proto.ItemSlot_ItemSlotMainHand, agent.GetCharacter().AutoAttacks.MH(), agent.GetCharacter().AutoAttacks.OH())
		w.BaseDamageMin += amount
		w.BaseDamageMax += amount
	})
}
