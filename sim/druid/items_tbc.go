package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// T4 Balance
var ItemSetMalorneRegalia = core.NewItemSet(core.ItemSet{
	Name: "Malorne Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			manaMetrics := druid.NewManaMetrics(core.ActionID{SpellID: 37295})

			druid.RegisterAura(core.Aura{
				Label:    "Malorne Regalia 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						return
					}
					if !result.Landed() {
						return
					}
					if sim.RandomFloat("malorne 2p") > 0.05 {
						return
					}
					spell.Unit.AddMana(sim, 120, manaMetrics)
				},
			})
		},
		4: func(agent core.Agent) {
			// Currently this is handled in druid.go (reducing CD of innervate)
		},
	},
})

// T5 Balance
var ItemSetNordrassilRegalia = core.NewItemSet(core.ItemSet{
	Name: "Nordrassil Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		4: func(agent core.Agent) {
			// Implemented in starfire.go.
		},
	},
})

// T6 Balance
var ItemSetThunderheartRegalia = core.NewItemSet(core.ItemSet{
	Name: "Thunderheart Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in moonfire.go
		},
		4: func(agent core.Agent) {
			// Implemented in starfire.go
		},
	},
})

// T4 feral
var ItemSetMalorneHarness = core.NewItemSet(core.ItemSet{
	Name: "Malorne Harness",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()

			procChance := 0.04
			rageMetrics := druid.NewRageMetrics(core.ActionID{SpellID: 37306})
			energyMetrics := druid.NewEnergyMetrics(core.ActionID{SpellID: 37311})

			druid.RegisterAura(core.Aura{
				Label:    "Malorne 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMelee) {
						if sim.RandomFloat("Malorne 2pc") < procChance {
							if druid.InForm(Bear) {
								druid.AddRage(sim, 10, rageMetrics)
							} else if druid.InForm(Cat) {
								druid.AddEnergy(sim, 20, energyMetrics)
							}
						}
					}
				},
			})
		},
		4: func(agent core.Agent) {
			druid := agent.(DruidAgent).GetDruid()
			if druid.InForm(Bear) {
				druid.AddStat(stats.Armor, 1400)
			} else if druid.InForm(Cat) {
				druid.AddStat(stats.Strength, 30)
			}
		},
	},
})

// T5 Feral
var ItemSetNordrassilHarness = core.NewItemSet(core.ItemSet{
	Name: "Nordrassil Harness",
	Bonuses: map[int32]core.ApplyEffect{
		4: func(agent core.Agent) {
			// Implemented in lacerate.go.
		},
	},
})

// T6 Feral
var ItemSetThunderheartHarness = core.NewItemSet(core.ItemSet{
	Name: "Thunderheart Harness",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented in mangle.go.
		},
		4: func(agent core.Agent) {
			// Implemented in swipe.go.
		},
	},
})

func init() {

	core.NewItemEffect(30664, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		var procAura *core.Aura
		if druid.InForm(Moonkin) {
			procAura = druid.NewTemporaryStatsAura("Living Root Moonkin Proc", core.ActionID{SpellID: 37343}, stats.Stats{stats.SpellPower: 209}, time.Second*15)
		} else if druid.InForm(Bear) {
			procAura = druid.NewTemporaryStatsAura("Living Root Bear Proc", core.ActionID{SpellID: 37340}, stats.Stats{stats.Armor: 4070}, time.Second*15)
		} else if druid.InForm(Cat) {
			procAura = druid.NewTemporaryStatsAura("Living Root Cat Proc", core.ActionID{SpellID: 37341}, stats.Stats{stats.Strength: 64}, time.Second*15)
		} else {
			return
		}

		druid.RegisterAura(core.Aura{
			Label:    "Living Root of the Wildheart",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
				if druid.InForm(Moonkin) && sim.RandomFloat("Living Root of the Wildheart") < 0.03 {
					procAura.Activate(sim)
				}
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				if sim.RandomFloat("Living Root of the Wildheart") > 0.03 {
					return
				}

				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(32486, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		// Not in the game yet so cant test; this logic assumes that:
		// - does not affect the starfire which procs it
		// - can proc off of any completed cast, not just hits
		actionID := core.ActionID{ItemID: 32486}

		var procAura *core.Aura
		if druid.InForm(Moonkin) {
			procAura = druid.NewTemporaryStatsAura("Ashtongue Talisman Proc", actionID, stats.Stats{stats.SpellPower: 150}, time.Second*8)
		} else if druid.InForm(Bear | Cat) {
			procAura = druid.NewTemporaryStatsAura("Ashtongue Talisman Proc", actionID, stats.Stats{stats.Strength: 140}, time.Second*8)
		} else {
			return
		}

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Ashtongue Talisman",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}
				if druid.Starfire.IsEqual(spell) {
					if sim.RandomFloat("Ashtongue Talisman") < 0.25 {
						procAura.Activate(sim)
					}
				} else if druid.IsMangle(spell) {
					if sim.RandomFloat("Ashtongue Talisman") < 0.4 {
						procAura.Activate(sim)
					}
				}
			},
		}))
	})

	core.NewItemEffect(32257, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := druid.NewTemporaryStatsAura("Idol of the White Stag Proc", core.ActionID{SpellID: 41037}, stats.Stats{stats.AttackPower: 94}, time.Second*20)

		core.MakePermanent(druid.RegisterAura(core.Aura{
			Label: "Idol of the White Stag",
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if druid.IsMangle(spell) {
					procAura.Activate(sim)
				}
			},
		}))
	})

	core.NewItemEffect(33509, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := druid.NewTemporaryStatsAura("Idol of Terror Proc", core.ActionID{SpellID: 43738}, stats.Stats{stats.Agility: 65}, time.Second*10)

		core.MakeProcTriggerAura(&druid.Unit, core.ProcTrigger{
			Name:       "Idol of Terror",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.85,
			ICD:        time.Second * 10,
			CustomCheck: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return druid.IsMangle(spell)
			},
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(33510, func(agent core.Agent) {
		druid := agent.(DruidAgent).GetDruid()

		procAura := druid.NewTemporaryStatsAura("Idol of the Unseen Moon Proc", core.ActionID{SpellID: 43740}, stats.Stats{stats.SpellPower: 140}, time.Second*10)

		druid.RegisterAura(core.Aura{
			Label:    "Idol of the Unseen Moon",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if druid.Moonfire.IsEqual(spell) {
					if sim.RandomFloat("Idol of the Unseen Moon") > 0.5 {
						return
					}
					procAura.Activate(sim)
				}
			},
		})
	})
}
