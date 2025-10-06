package warrior

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// ///////////////////////////////////////////////////////////////
// TBC Item set
// ///////////////////////////////////////////////////////////////
var ItemSetWarbringerArmor = core.NewItemSet(core.ItemSet{
	Name: "Warbringer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warrior := agent.(WarriorAgent).GetWarrior()
			// You have a chance each time you parry to gain Blade Turning, absorbing 200 damage for 15 sec.
			bladeTurningShield := 0.0

			procaura := warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 2pc Proc",
				ActionID: core.ActionID{SpellID: 37514},
				Duration: time.Second * 15,
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					damage := result.Damage
					if damage >= bladeTurningShield {
						result.Damage -= bladeTurningShield
						aura.Deactivate(sim)
					} else {
						result.Damage -= bladeTurningShield
						bladeTurningShield -= damage
					}
				},
			})
			warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Outcome.Matches(core.OutcomeParry) {
						procaura.Activate(sim)
						bladeTurningShield = 200
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Your Revenge ability causes your next damaging ability to do 10% more damage.
			warrior := agent.(WarriorAgent).GetWarrior()

			// TODO: This needs to apply only to specific abilities, not any source of damage.
			procAura := warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc Proc",
				ActionID: core.ActionID{SpellID: 37516},
				Duration: core.NeverExpires,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Damage > 0 {
						aura.Deactivate(sim)
					}
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellResult) {
					if spell.IsSpell(warrior.Revenge) {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetDestroyerArmor = core.NewItemSet(core.ItemSet{
	Name: "Destroyer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Each time you use your Shield Block ability, you gain 100 block value against a single attack in the next 6 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.RegisterAura(core.Aura{
				Label:     "Destroyer 2pc Proc",
				ActionID:  core.ActionID{SpellID: 37522},
				Duration:  time.Second * 6,
				MaxStacks: 1,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.BlockValue, 100)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warrior.AddStatDynamic(sim, stats.BlockValue, -100)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Outcome.Matches(core.OutcomeBlock) {
						aura.RemoveStack(sim)
					}
				},
			})

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.IsSpell(warrior.ShieldBlock) {
						procAura.Activate(sim)
						procAura.SetStacks(sim, 1)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// You have a chance each time you are hit to gain 200 haste rating for 10 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.NewTemporaryStatsAura("Destroyer 4pc Proc", core.ActionID{SpellID: 37525}, stats.Stats{stats.MeleeHaste: 200}, time.Second*10)

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && sim.RandomFloat("Destroyer 4pc") < 0.07 {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

var ItemSetOnslaughtArmor = core.NewItemSet(core.ItemSet{
	Name: "Onslaught Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the health bonus from your Commanding Shout ability by 170.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Shield Slam ability by 10%.
			// Handled in shield_slam.go.
		},
	},
})

var ItemSetWarbringerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Warbringer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your whirlwind ability costs 5 less rage.
			// Handled in whirlwind.go.
		},
		4: func(agent core.Agent) {
			// You gain an additional 2 rage each time one of your attacks is parried or dodged.
			warrior := agent.(WarriorAgent).GetWarrior()
			rageMetrics := warrior.NewRageMetrics(core.ActionID{SpellID: 37519})

			warrior.RegisterAura(core.Aura{
				Label:    "Warbringer 4pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Outcome.Matches(core.OutcomeDodge | core.OutcomeParry) {
						warrior.AddRage(sim, 2, rageMetrics)
					}
				},
			})
		},
	},
})

var ItemSetDestroyerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Destroyer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Your Overpower ability now grants you 100 attack power for 5 sec.
			warrior := agent.(WarriorAgent).GetWarrior()
			procAura := warrior.NewTemporaryStatsAura("Destroyer 2pc Proc", core.ActionID{SpellID: 37528}, stats.Stats{stats.AttackPower: 100}, time.Second*5)

			warrior.RegisterAura(core.Aura{
				Label:    "Destroyer 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.IsSpell(warrior.Overpower) {
						procAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			// Your Bloodthirst and Mortal Strike abilities cost 5 less rage.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
})

var ItemSetOnslaughtBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Onslaught Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Reduces the rage cost of your Execute ability by 3.
		},
		4: func(agent core.Agent) {
			// Increases the damage of your Mortal Strike and Bloodthirst abilities by 5%.
			// Handled in bloodthirst.go and mortal_strike.go.
		},
	},
})

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(32485, func(a core.Agent) {
		character := a.(WarriorAgent).GetWarrior()
		metrics := character.NewHealthMetrics(core.ActionID{ItemID: 32485})

		aura := character.NewTemporaryStatsAura("Fire Blood", core.ActionID{SpellID: 40459}, stats.Stats{stats.Strength: 55}, time.Second*12)

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Ashtongue Talisman of Valor",
			ActionID:   core.ActionID{ItemID: 32485},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeSpecial,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.IsSpell(character.MortalStrike) || spell.IsSpell(character.Bloodthirst) || spell.IsSpell(character.ShieldSlam) {
					aura.Activate(sim)
					character.GainHealth(sim, 330, metrics)
				}
			},
		})
	})
	core.AddEffectsToTest = true
}
