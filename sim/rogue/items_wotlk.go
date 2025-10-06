package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

var Arena = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.Resilience, 100)
			agent.GetCharacter().AddStat(stats.AttackPower, 50)
		},
		4: func(agent core.Agent) {
			agent.GetCharacter().AddStat(stats.AttackPower, 150)
			// 10 maximum energy added in rogue.go
		},
	},
})

var Tier10 = core.NewItemSet(core.ItemSet{
	Name: "Shadowblade's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Tricks of the Trade now grants you 15 energy instead of costing energy.
			// Handled in tricks_of_the_trade.go.
		},
		4: func(agent core.Agent) {
			// Gives your melee finishing moves a 13% chance to add 3 combo points to your target.
			// Handled in the finishing move effect applier
		},
	},
})

var Tier9 = core.NewItemSet(core.ItemSet{
	Name:            "VanCleef's Battlegear",
	AlternativeName: "Garona's Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Rupture ability has a chance each time it deals damage to reduce the cost of your next ability by 40 energy.
			rogue := agent.(RogueAgent).GetRogue()
			energyMetrics := rogue.NewEnergyMetrics(core.ActionID{SpellID: 67209})

			procAura := rogue.RegisterAura(core.Aura{
				Label:    "VanCleef's 2pc Proc",
				ActionID: core.ActionID{SpellID: 67209},
				Duration: time.Second * 15,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					rogue.PseudoStats.CostReduction += 40
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					rogue.PseudoStats.CostReduction -= 40
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMeleeSpecial) {
						return
					}

					// doesn't handle multiple dynamic cost reductions at once, or 0-cost default casts
					if actualGain := spell.DefaultCast.Cost - spell.CurCast.Cost; actualGain > 0 {
						energyMetrics.AddEvent(40, actualGain)
						aura.Deactivate(sim)
					}
				},
			})

			icd := core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 15,
			}
			procAura.Icd = &icd
			procChance := 0.02
			rogue.RegisterAura(core.Aura{
				Label:    "VanCleef's 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if !spell.ActionID.IsSpell(rogue.Rupture) {
						return
					}
					if !icd.IsReady(sim) {
						return
					}
					if sim.RandomFloat("VanCleef's 2pc") > procChance {
						return
					}
					icd.Use(sim)
					procAura.Activate(sim)
				},
			})
		},
		4: func(agent core.Agent) {
			// Increases the critical strike chance of your Hemorrhage, Sinister Strike, Backstab, and Mutilate abilities by 5%.
			// Handled in ability sources
		},
	},
})

var Tier8 = core.NewItemSet(core.ItemSet{
	Name: "Terrorblade Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Deadly Poison causes you to gain 1 energy each time it deals damage
			// Handled in poisons.go
		},
		4: func(agent core.Agent) {
			// Increases the damage done by your Rupture by 20%
			// Handled in rupture.go
		},
	},
})

var Tier7 = core.NewItemSet(core.ItemSet{
	Name: "Bonescythe Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage dealt by your Rupture by 10%
			// Handled in rupture.go
		},
		4: func(agent core.Agent) {
			// Reduce the Energy cost of your Combo Moves by 5%
			// Handled in the builder cast modifier
		},
	},
})
