package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// T4
var ItemSetAldorRegalia = core.NewItemSet(core.ItemSet{
	Name: "Aldor Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Interruption avoidance.
		},
		4: func(agent core.Agent) {
			// Reduces the cooldown on PoM/Blast Wave/Ice Block.
		},
	},
})

// T5
var ItemSetTirisfalRegalia = core.NewItemSet(core.ItemSet{
	Name: "Tirisfal Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage and mana cost of Arcane Blast by 5%.
			// Implemented in arcane_blast.go.
		},
		4: func(agent core.Agent) {
			mage := agent.(MageAgent).GetMage()
			// Your spell critical strikes grant you up to 70 spell damage for 6 sec.
			procAura := mage.NewTemporaryStatsAura("Tirisfal 4pc Proc", core.ActionID{SpellID: 37443}, stats.Stats{stats.SpellPower: 70}, time.Second*6)
			mage.RegisterAura(core.Aura{
				Label:    "Tirisfal 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if result.Outcome.Matches(core.OutcomeCrit) {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

// T6 Sunwell
var ItemSetTempestRegalia = core.NewItemSet(core.ItemSet{
	Name: "Tempest Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the duration of your Evocation ability by 2 sec.
			// Implemented in evocation.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Fireball, Frostbolt, and Arcane Missles abilities by 5%.
			// Implemented in the files for those spells.
		},
	},
})

func init() {
	core.NewItemEffect(32488, func(a core.Agent) {
		character := a.GetCharacter()
		aura := character.NewTemporaryStatsAura("Insight of the Ashtongue", core.ActionID{SpellID: 40483}, stats.Stats{stats.SpellHaste: 145}, time.Second*5)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Ashtongue Talisman of Insight",
			ActionID:   core.ActionID{ItemID: 32488},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeCrit,
			ProcChance: 0.5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.Activate(sim)
			},
		})
	})
}
