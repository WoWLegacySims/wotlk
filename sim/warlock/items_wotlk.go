package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// T7
var ItemSetPlagueheartGarb = core.NewItemSet(core.ItemSet{
	Name: "Plagueheart Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			const bonusCrit = 10
			warlock.DemonicSoulAura = warlock.RegisterAura(core.Aura{
				Label:    "Demonic Soul",
				ActionID: core.ActionID{SpellID: 61595},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCrit += bonusCrit
					warlock.Incinerate.BonusCrit += bonusCrit
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.ShadowBolt.BonusCrit -= bonusCrit
					warlock.Incinerate.BonusCrit -= bonusCrit
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warlock.ShadowBolt || spell == warlock.Incinerate {
						warlock.DemonicSoulAura.Deactivate(sim)
					}
				},
			})

			warlock.RegisterAura(core.Aura{
				Label: "2pT7 Hidden Aura",
				// ActionID: core.ActionID{SpellID: 60170},
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if (spell == warlock.Corruption || spell == warlock.Immolate) && sim.Proc(0.15, "2pT7") {
						warlock.DemonicSoulAura.Activate(sim)
					}
				},
			})
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			warlock.SpiritsoftheDamnedAura = warlock.RegisterAura(core.Aura{
				Label:    "Spirits of the Damned",
				ActionID: core.ActionID{SpellID: 61082},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, 300)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddStatDynamic(sim, stats.Spirit, -300)
				},
			})

			warlock.RegisterAura(core.Aura{
				Label:    "4pT7 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell == warlock.LifeTap {
						warlock.SpiritsoftheDamnedAura.Activate(sim)
					}
				},
			})
		},
	},
})

// T8
var ItemSetDeathbringerGarb = core.NewItemSet(core.ItemSet{
	Name: "Deathbringer Garb",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented
		},
		4: func(agent core.Agent) {
			// Implemented
		},
	},
})

// T9
var ItemSetGuldansRegalia = core.NewItemSet(core.ItemSet{
	Name:            "Gul'dan's Regalia",
	AlternativeName: "Kel'Thuzad's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			if warlock.Pet != nil {
				warlock.Pet.AddStats(stats.Stats{
					stats.MeleeCrit: 10 * warlock.Pet.CritRatingPerCritChance,
					stats.SpellCrit: 10 * warlock.Pet.CritRatingPerCritChance,
				})
			}
		},
		4: func(agent core.Agent) {
			// Implemented
		},
	},
})

// T10
var ItemSetDarkCovensRegalia = core.NewItemSet(core.ItemSet{
	Name: "Dark Coven's Regalia",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Implemented
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()

			deviousMindsAura := warlock.RegisterAura(core.Aura{
				Label:    "Devious Minds",
				ActionID: core.ActionID{SpellID: 70840},
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
				},
			})

			var petDeviousMindsAura *core.Aura
			if warlock.Pet != nil {
				petDeviousMindsAura = warlock.Pet.RegisterAura(core.Aura{
					Label:    "Devious Minds",
					ActionID: core.ActionID{SpellID: 70840},
					Duration: time.Second * 10,
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
					},
				})
			}

			warlock.RegisterAura(core.Aura{
				Label:    "4pT10 Hidden Aura",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == warlock.UnstableAffliction || spell == warlock.Immolate {
						if sim.Proc(0.15, "4pT10") {
							deviousMindsAura.Activate(sim)
							if petDeviousMindsAura != nil {
								petDeviousMindsAura.Activate(sim)
							}
						}
					}
				},
			})
		},
	},
})

var ItemSetGladiatorsFelshroud = core.NewItemSet(core.ItemSet{
	Name: "Gladiator's Felshroud",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 29)
		},
		4: func(agent core.Agent) {
			warlock := agent.(WarlockAgent).GetWarlock()
			warlock.AddStat(stats.SpellPower, 88)
		},
	},
})
