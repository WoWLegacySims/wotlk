package paladin

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// T4 Ret
var ItemSetJusticarBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Justicar Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// sim/debuffs.go handles this (and paladin/judgement.go)
		},
		4: func(agent core.Agent) {
			// soc/sor/sov.go
		},
	},
})

// T5 Ret
var ItemSetCrystalforgeBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Crystalforge Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// judgement.go
		},
		4: func(agent core.Agent) {
			// TODO: if we implement healing, this heals party.
		},
	},
})

// Tier 6 ret
var ItemSetLightbringerBattlegear = core.NewItemSet(core.ItemSet{
	Name: "Lightbringer Battlegear",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			paladin := agent.(PaladinAgent).GetPaladin()
			manaMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 38428})

			paladin.RegisterAura(core.Aura{
				Label:    "Lightbringer Battlegear 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !spell.ProcMask.Matches(core.ProcMaskMelee) {
						return
					}
					if sim.RandomFloat("lightbringer 2pc") > 0.2 {
						return
					}
					paladin.AddMana(sim, 50, manaMetrics)
				},
			})
		},
		4: func(agent core.Agent) {
			// Implemented in hammer_of_wrath.go
		},
	},
})

func (paladin *Paladin) getItemSetLightbringerBattlegearBonus4() float64 {
	return core.TernaryFloat64(paladin.HasSetBonus(ItemSetLightbringerBattlegear, 4), .1, 0)
}

// T4 Prot
var ItemSetJusticarArmor = core.NewItemSet(core.ItemSet{
	Name: "Justicar Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage dealt by your Seal of Righteousness, Seal of
			// Vengeance, and Seal of Corruption by 10%.
			// Implemented in seals.go.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by Holy Shield by 15.
			// Implemented in holy_shield.go.
		},
	},
})

// T5 Prot
var ItemSetCrystalforgeArmor = core.NewItemSet(core.ItemSet{
	Name: "Crystalforge Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the damage from your Retribution Aura by 15.
			// TODO
		},
		4: func(agent core.Agent) {
			// Each time you use your Holy Shield ability, you gain 100 Block Value
			// against a single attack in the next 6 seconds.
			paladin := agent.(PaladinAgent).GetPaladin()

			procAura := paladin.RegisterAura(core.Aura{
				Label:    "Crystalforge 4pc Proc",
				ActionID: core.ActionID{SpellID: 37191},
				Duration: time.Second * 6,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.BlockValue, 100)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					paladin.AddStatDynamic(sim, stats.BlockValue, -100)
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellResult) {
					if spellEffect.Outcome.Matches(core.OutcomeBlock) {
						aura.Deactivate(sim)
					}
				},
			})

			paladin.RegisterAura(core.Aura{
				Label:    "Crystalforge 2pc",
				Duration: core.NeverExpires,
				OnReset: func(aura *core.Aura, sim *core.Simulation) {
					aura.Activate(sim)
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.IsSpell(paladin.HolyShield) {
						procAura.Activate(sim)
					}
				},
			})
		},
	},
})

// T6 Prot
var ItemSetLightbringerArmor = core.NewItemSet(core.ItemSet{
	Name: "Lightbringer Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			// Increases the mana gained from your Spiritual Attunement ability by 10%.
		},
		4: func(agent core.Agent) {
			// Increases the damage dealt by Consecration by 10%.
		},
	},
})

func init() {

	core.NewItemEffect(27484, func(agent core.Agent) {
		paladin := agent.(PaladinAgent).GetPaladin()
		procAura := paladin.NewTemporaryStatsAura("Libram of Avengement Proc", core.ActionID{SpellID: 48835}, stats.Stats{stats.MeleeCrit: 53, stats.SpellCrit: 53}, time.Second*5)

		paladin.RegisterAura(core.Aura{
			Label:    "Libram of Avengement",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.Flags.Matches(SpellFlagSecondaryJudgement) {
					procAura.Activate(sim)
				}
			},
		})
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Tome of Fiery Redemption",
		ID:         30447,
		AuraID:     37198,
		Bonus:      stats.Stats{stats.SpellPower: 290},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpell,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Libram of Divine Judgement",
		ID:         33503,
		AuraID:     43747,
		Bonus:      stats.Stats{stats.AttackPower: 200},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.4,
		CustomCheck: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
			return spell.IsSpellAction(20467)
		},
	})
}
