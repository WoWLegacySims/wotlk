package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// D3
var ItemSetOblivionRaiment = core.NewItemSet(core.ItemSet{
	Name: "Oblivion Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// in pet.go constructor
		},
		4: func(agent core.Agent) {
			// in seed.go
		},
	},
})

// T4
var ItemSetVoidheartRaiment = core.NewItemSet(core.ItemSet{
	Name: "Voidheart Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			shadowBonus := warlock.RegisterAura(core.Aura{
				Label:    "Shadowflame",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 37377},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.AddPseudoSpellpowerDynamic(sim, core.SpellSchoolShadow, 135)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.AddPseudoSpellpowerDynamic(sim, core.SpellSchoolShadow, -135)
				},
			})

			fireBonus := warlock.RegisterAura(core.Aura{
				Label:    "Shadowflame Hellfire",
				Duration: time.Second * 15,
				ActionID: core.ActionID{SpellID: 39437},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.AddPseudoSpellpowerDynamic(sim, core.SpellSchoolFire, 135)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.AddPseudoSpellpowerDynamic(sim, core.SpellSchoolFire, -135)
				},
			})

			core.MakeProcTriggerAura(&warlock.Unit, core.ProcTrigger{
				Name:        "Voidheart Raiment 2PC",
				Callback:    core.CallbackOnCastComplete,
				ProcMask:    core.ProcMaskSpellDamage,
				SpellSchool: core.SpellSchoolFire | core.SpellSchoolShadow,
				ProcChance:  0.05,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.SpellSchool.Matches(core.SpellSchoolFire) {
						fireBonus.Activate(sim)
					}
					if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
						shadowBonus.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// implemented in immolate.go and corruption.go
		},
	},
})

// T5
var ItemSetCorruptorRaiment = core.NewItemSet(core.ItemSet{
	Name: "Corruptor Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// heals pet
		},
		4: func(agent core.Agent) {
			// implemented in immolate.go and corruption.go
		},
	},
})

// T6
var ItemSetMaleficRaiment = core.NewItemSet(core.ItemSet{
	Name: "Malefic Raiment",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// heals... not implemented yet
		},
		4: func(agent core.Agent) {
			// Increases damage done by shadowbolt and incinerate by 6%.
			// Implemented in shadowbolt.go and incinerate.go
		},
	},
})

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(32493, func(agent core.Agent) {
		warlock := agent.(WarlockAgent).GetWarlock()
		procAura := warlock.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{SpellID: 40478}, stats.Stats{stats.SpellPower: 220}, time.Second*5)

		warlock.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell == warlock.Corruption && sim.Proc(0.2, "Ashtongue Talisman of Insight") {
					procAura.Activate(sim)
				}
			},
		})
	})
	core.AddEffectsToTest = true
}
