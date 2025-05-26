package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

var Tier6 = core.NewItemSet(core.ItemSet{
	Name: "Slayer's Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the haste from your Slice and Dice ability by 5%.
			// Handled in slice_and_dice.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by your Backstab, Sinister Strike, Mutilate, and Hemorrhage abilities by 6%.
			// Handled in the corresponding ability files.
		},
	},
})

var Tier5 = core.NewItemSet(core.ItemSet{
	Name: "Deathmantle",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Eviscerate and Envenom abilities cause 40 extra damage per combo point.
			// Handled in eviscerate.go.
		},
		4: func(agent core.Agent) {
			// Your attacks have a chance to make your next finishing move cost no energy.
			rogue := agent.(RogueAgent).GetRogue()

			rogue.DeathmantleProcAura = rogue.RegisterAura(core.Aura{
				Label:    "Deathmantle 4pc Proc",
				ActionID: core.ActionID{SpellID: 37171},
				Duration: time.Second * 15,
			})

			ppmm := rogue.AutoAttacks.NewPPMManager(1.0, core.ProcMaskMelee)

			rogue.RegisterAura(core.Aura{
				Label:    "Deathmantle 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}

					// https://tbc.wowhead.com/spell=37170/free-finisher-chance, proc mask = 20.
					if !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}

					if !ppmm.Proc(sim, spell.ProcMask, "Deathmantle 4pc") {
						return
					}

					rogue.DeathmantleProcAura.Activate(sim)
				},
			})
		},
	},
})

func (rogue *Rogue) applyDeathmantle(sim *core.Simulation, spell *core.Spell, _ *core.Cast) {
	if rogue.DeathmantleProcAura != nil && rogue.DeathmantleProcAura.IsActive() {
		spell.CostMultiplier = 0
		rogue.DeathmantleProcAura.Deactivate(sim)
	}
}

var Tier4 = core.NewItemSet(core.ItemSet{
	Name: "Netherblade",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the duration of your Slice and Dice ability by 3 sec.
			// Handled in slice_and_dice.go.
		},
		4: func(agent core.Agent) {
			// Your finishing moves have a 15% chance to grant you an extra combo point.
			// Handled in talents.go.
		},
	},
})

var D3 = core.NewItemSet(core.ItemSet{
	Name: "Assassination Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
		},
		4: func(agent core.Agent) {
			// Your Eviscerate and Envenom abilities cost 10 less energy.
			// Handled in eviscerate.go.
		},
	},
})

func init() {
	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Warp-Spring-Coil",
		ID:         30450,
		AuraID:     37174,
		Bonus:      stats.Stats{stats.ArmorPenetration: 142},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeMHSpecial,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 30,
	})

	core.NewItemEffect(32492, func(agent core.Agent) {
		rogue := agent.(RogueAgent).GetRogue()
		procAura := rogue.NewTemporaryStatsAura("Ashtongue Talisman Proc", core.ActionID{ItemID: 32492}, stats.Stats{stats.MeleeCrit: 145}, time.Second*10)

		var numPoints int32

		rogue.RegisterAura(core.Aura{
			Label:    "Ashtongue Talisman",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				numPoints = 0
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				// Need to store the points because they get spent before OnSpellHit is called.
				numPoints = rogue.ComboPoints()

				if spell == rogue.SliceAndDice {
					// SND won't call OnSpellHit, so we have to add the effect now.
					if p := 0.2 * float64(numPoints); sim.Proc(p, "AshtongueTalismanOfLethality") {
						procAura.Activate(sim)
					}
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.Flags.Matches(SpellFlagFinisher) {
					return
				}

				if p := 0.2 * float64(numPoints); sim.Proc(p, "AshtongueTalismanOfLethality") {
					procAura.Activate(sim)
				}
			},
		})
	})

}
