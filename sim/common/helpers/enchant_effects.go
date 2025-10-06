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

func AddScope(id int32, amount float64) {
	core.NewEnchantEffect(id, func(a core.Agent) {
		wp := a.GetCharacter().GetRangedWeapon()
		wp.WeaponDamageMin += amount
		wp.WeaponDamageMax += amount
	})
}

func AddShieldSpike(id int32, itemID int32, name string, bp float64, die float64) {
	core.NewEnchantEffect(id, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{ItemID: itemID}

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID,
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(bp, die)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     name,
			Callback: core.CallbackOnSpellHitTaken,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeLanded,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, spell.Unit)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(id, aura)
	})
}

func AddAbsorption(id int32, spellId int32, name string, amount float64, chance float64) {
	core.NewEnchantEffect(id, func(agent core.Agent) {
		character := agent.GetCharacter()
		shield := character.RegisterAura(core.Aura{
			Label:    name,
			ActionID: core.ActionID{SpellID: spellId},
			Duration: core.NeverExpires,
		})

		character.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if shield.IsActive() && (result.Damage > 0) && spell.SpellSchool.Matches(core.SpellSchoolPhysical) {
				result.Damage = max(0, result.Damage-amount)
				shield.Deactivate(sim)
			}
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:        name + " Trigger",
			Callback:    core.CallbackOnSpellHitTaken,
			SpellSchool: core.SpellSchoolPhysical,
			ProcChance:  chance,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				shield.Activate(sim)
			},
		})
	})
}
