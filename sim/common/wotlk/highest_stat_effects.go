package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	newDMCGreatnessEffect := func(itemID int32) {
		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			auraIDs := map[stats.Stat]core.ActionID{
				stats.Strength:  core.ActionID{SpellID: 60229},
				stats.Agility:   core.ActionID{SpellID: 60233},
				stats.Intellect: core.ActionID{SpellID: 60234},
				stats.Spirit:    core.ActionID{SpellID: 60235},
			}

			hsa := helpers.NewHighestStatAura(
				[]stats.Stat{
					stats.Strength,
					stats.Agility,
					stats.Intellect,
					stats.Spirit,
				},
				func(stat stats.Stat) *core.Aura {
					bonus := stats.Stats{}
					bonus[stat] = 300
					actionId := auraIDs[stat]
					return character.NewTemporaryStatsAura("DMC Greatness "+stat.StatName()+" Proc", actionId, bonus, time.Second*15)
				})

			hsa.Init(character)
			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       "DMC Greatness",
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt | core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
				ProcMask:   core.ProcMaskDirect | core.ProcMaskSpellHealing | core.ProcMaskProc,
				Harmful:    true,
				ProcChance: 0.35,
				ICD:        time.Second * 45,
				ActionID:   core.ActionID{ItemID: itemID},
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
			hsa.Get(character).Icd = triggerAura.Icd
		})
	}

	newDeathsChoiceEffect := func(itemID int32, auraIDs map[stats.Stat]core.ActionID, name string, amount float64) {
		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()

			hsa := helpers.NewHighestStatAura(
				[]stats.Stat{
					stats.Strength,
					stats.Agility,
				},
				func(stat stats.Stat) *core.Aura {
					bonus := stats.Stats{}
					bonus[stat] = amount
					actionId := auraIDs[stat]
					return character.NewTemporaryStatsAura(name+" "+stat.StatName()+" Proc", actionId, bonus, time.Second*15)
				})

			hsa.Init(character)
			triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name,
				Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
				ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
				Harmful:    true,
				ProcChance: 0.35,
				ActionID:   core.ActionID{ItemID: itemID},
				ICD:        time.Second * 45,
				Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
					hsa.Get(character).Activate(sim)
				},
			})
			hsa.Get(character).Icd = triggerAura.Icd
		})
	}
	normalAuraIDs := map[stats.Stat]core.ActionID{
		stats.Strength: core.ActionID{SpellID: 67708},
		stats.Agility:  core.ActionID{SpellID: 67703},
	}

	heroicAuraIDs := map[stats.Stat]core.ActionID{
		stats.Strength: core.ActionID{SpellID: 67773},
		stats.Agility:  core.ActionID{SpellID: 67772},
	}

	core.AddEffectsToTest = false
	newDMCGreatnessEffect(42987)
	newDMCGreatnessEffect(44253)
	newDMCGreatnessEffect(44254)
	newDeathsChoiceEffect(47131, heroicAuraIDs, "Deaths Verdict H", 510)
	newDeathsChoiceEffect(47303, normalAuraIDs, "Deaths Choice", 450)
	newDeathsChoiceEffect(47464, heroicAuraIDs, "Deaths Choice H", 510)
	core.AddEffectsToTest = true
	newDMCGreatnessEffect(44255)
	newDeathsChoiceEffect(47115, normalAuraIDs, "Deaths Verdict", 450)

}
