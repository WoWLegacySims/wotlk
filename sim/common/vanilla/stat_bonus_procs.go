package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Destiny",
		ID:       647,
		AuraID:   17152,
		Bonus:    stats.Stats{stats.Strength: 200},
		Duration: time.Second * 10,
		PPM:      1.3,
		Weapon:   true,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Guardian Talisman",
		ID:         1490,
		AuraID:     4070,
		Bonus:      stats.Stats{stats.BonusArmor: 350},
		Duration:   time.Second * 15,
		ProcChance: 0.02,
		Outcome:    core.OutcomeLanded,
		Callback:   core.CallbackOnSpellHitTaken,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Wall of the Dead",
		ID:         1979,
		AuraID:     19409,
		Bonus:      stats.Stats{stats.BonusArmor: 150},
		Duration:   time.Second * 20,
		ProcChance: 0.03,
		Outcome:    core.OutcomeLanded,
		Callback:   core.CallbackOnSpellHitTaken,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Hand of Edward the Odd",
		ID:       2243,
		AuraID:   18803,
		Bonus:    stats.Stats{stats.MeleeHaste: 320, stats.SpellHaste: 320},
		Duration: time.Second * 4,
		PPM:      2,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
		Weapon:   true,
	})

	core.NewItemEffect(6622, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(6622)

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Zeal",
			ActionID: core.ActionID{SpellID: 8191},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.BonusArmor, 150)
				character.PseudoStats.BonusDamage += 10
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.BonusArmor, -150)
				character.PseudoStats.BonusDamage -= 10
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 6622},
			Name:     "Sword of Zeal",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	core.NewItemEffect(6660, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(6660)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 8348},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			Hot: core.DotConfig{
				SelfOnly:      true,
				NumberOfTicks: 6,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 13
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotHealing(sim, &character.Unit, dot.Spell.OutcomeHealing)
				},
				Aura: core.Aura{Label: "Julie's Dagger"},
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 6660},
			Name:     "Julie's Dagger",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      2.6,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, &character.Unit)
			},
		})
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Blade of the Basilisk",
		ID:       8223,
		AuraID:   10351,
		Bonus:    stats.Stats{stats.Defense: 50},
		Duration: time.Second * 5,
		PPM:      3.4,
		Weapon:   true,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
	})

	core.NewItemEffect(9418, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(9418)

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Strength of Stone",
			ActionID: core.ActionID{SpellID: 12731},
			Duration: time.Second * 8,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.BonusDamage += 10
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.BonusDamage -= 10
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 9418},
			Name:     "Stoneslayer",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.3,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "The Jackhammer",
		ID:       9423,
		AuraID:   13533,
		Bonus:    stats.Stats{stats.MeleeHaste: 300, stats.SpellHaste: 300},
		Duration: time.Second * 10,
		PPM:      2,
		Weapon:   true,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Felstriker",
		ID:       12590,
		AuraID:   16551,
		Bonus:    stats.Stats{stats.MeleeCrit: 10000000},
		Duration: time.Second * 3,
		PPM:      1,
		Weapon:   true,
		Outcome:  core.OutcomeLanded,
		Callback: core.CallbackOnSpellHitDealt,
	})
	core.AddEffectsToTest = true
}
