package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.NewItemEffect(38212, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Death Knight's Anguish Proc",
				ActionID:  core.ActionID{SpellID: 54697},
				Duration:  time.Second * 20,
				MaxStacks: 10,
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
						aura.AddStack(sim)
					}
				},
			},
			BonusPerStack: stats.Stats{stats.MeleeCrit: 15, stats.SpellCrit: 15},
		})

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Death Knight's Anguish",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMeleeOrRanged,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.1,
			ActionID:   core.ActionID{SpellID: 54696},
			ICD:        time.Second * 45,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
		procAura.Icd = triggerAura.Icd
	})

	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Majestic Dragon Figurine",
		ID:        40430,
		AuraID:    60525,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.Spirit: 18},
		Callback:  core.CallbackOnCastComplete,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Fury of the Five Fights",
		ID:        40431,
		AuraID:    60314,
		Duration:  time.Second * 10,
		MaxStacks: 20,
		Bonus:     stats.Stats{stats.AttackPower: 16, stats.RangedAttackPower: 16},
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeOrRanged,
		Harmful:   true,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Illustration of the Dragon Soul",
		ID:        40432,
		AuraID:    60486,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 20},
		Callback:  core.CallbackOnCastComplete,
		ProcMask:  core.ProcMaskSpellHealing | core.ProcMaskSpellDamage,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:       "DMC Berserker",
		ID:         42989,
		AuraID:     60196,
		Duration:   time.Second * 12,
		MaxStacks:  3,
		Bonus:      stats.Stats{stats.MeleeCrit: 35, stats.SpellCrit: 35},
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnSpellHitTaken,
		Harmful:    true,
		ProcChance: 0.5,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Eye of the Broodmother",
		ID:        45308,
		AuraID:    65006,
		Duration:  time.Second * 10,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.SpellPower: 26},
		Callback:  core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicHealDealt | core.CallbackOnPeriodicDamageDealt,
	})

	core.AddEffectsToTest = false

	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Solance of the Defeated",
		ID:        47041,
		AuraID:    67696,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 16},
		Callback:  core.CallbackOnCastComplete,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Solace of the Defeated H",
		ID:        47059,
		AuraID:    67750,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 18},
		Callback:  core.CallbackOnCastComplete,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Solace of the Fallen",
		ID:        47271,
		AuraID:    67696,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 16},
		Callback:  core.CallbackOnCastComplete,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Solace of the Fallen H",
		ID:        47432,
		AuraID:    67750,
		Duration:  time.Second * 10,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MP5: 18},
		Callback:  core.CallbackOnCastComplete,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Muradin's Spyglass",
		ID:        50340,
		AuraID:    71570,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 18},
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
		Harmful:   true,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Muradin's Spyglass H",
		ID:        50345,
		AuraID:    71572,
		Duration:  time.Second * 10,
		MaxStacks: 10,
		Bonus:     stats.Stats{stats.SpellPower: 20},
		Callback:  core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
		Harmful:   true,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:       "Unidentifiable Organ",
		ID:         50341,
		AuraID:     71575,
		Duration:   time.Second * 10,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.Stamina: 24},
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.6,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:       "Unidentifiable Organ H",
		ID:         50344,
		AuraID:     71577,
		Duration:   time.Second * 10,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.Stamina: 27},
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.6,
	})
	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "Herkuml War Token",
		ID:        50355,
		AuraID:    71396,
		Duration:  time.Second * 10,
		MaxStacks: 20,
		Bonus:     stats.Stats{stats.AttackPower: 17, stats.RangedAttackPower: 17},
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMeleeOrRanged,
		Harmful:   true,
	})

	// Stacking CD effects

	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:        "Meteorite Crystal",
		ID:          46051,
		AuraID:      65000,
		Duration:    time.Second * 20,
		MaxStacks:   20,
		Bonus:       stats.Stats{stats.MP5: 85},
		CD:          time.Minute * 2,
		Callback:    core.CallbackOnCastComplete,
		IsDefensive: true,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Victor's Call",
		ID:        47725,
		AuraID:    67737,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 215, stats.RangedAttackPower: 215},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Victor's Call H",
		ID:        47948,
		AuraID:    67746,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 250, stats.RangedAttackPower: 250},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Talisman of Volatile Power",
		ID:        47726,
		AuraID:    67735,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 57, stats.SpellHaste: 57},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnCastComplete,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Talisman of Volatile Power H",
		ID:        47946,
		AuraID:    67743,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 64, stats.SpellHaste: 64},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnCastComplete,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:        "Fervor of the Frostborn",
		ID:          47727,
		AuraID:      67727,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1265},
		CD:          time.Minute * 2,
		Callback:    core.CallbackOnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:        "Ferver of the Frostborn H",
		ID:          47949,
		AuraID:      67741,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1422},
		CD:          time.Minute * 2,
		Callback:    core.CallbackOnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:       "Binding Light",
		ID:         47728,
		AuraID:     67723,
		Duration:   time.Second * 20,
		MaxStacks:  8,
		Bonus:      stats.Stats{stats.SpellPower: 66},
		CD:         time.Minute * 2,
		Callback:   core.CallbackOnCastComplete,
		SpellFlags: core.SpellFlagHelpful,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:       "Binding Light H",
		ID:         47947,
		AuraID:     67739,
		Duration:   time.Second * 20,
		MaxStacks:  8,
		Bonus:      stats.Stats{stats.SpellPower: 74},
		CD:         time.Minute * 2,
		Callback:   core.CallbackOnCastComplete,
		SpellFlags: core.SpellFlagHelpful,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Fetish of Volatile Power",
		ID:        47879,
		AuraID:    67735,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 57, stats.SpellHaste: 57},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnCastComplete,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Fetish of Volatile Power H",
		ID:        48018,
		AuraID:    67743,
		Duration:  time.Second * 20,
		MaxStacks: 8,
		Bonus:     stats.Stats{stats.MeleeHaste: 64, stats.SpellHaste: 64},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnCastComplete,
		ProcMask:  core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:       "Binding Stone",
		ID:         47880,
		AuraID:     67723,
		Duration:   time.Second * 20,
		MaxStacks:  8,
		Bonus:      stats.Stats{stats.SpellPower: 66},
		CD:         time.Minute * 2,
		Callback:   core.CallbackOnCastComplete,
		SpellFlags: core.SpellFlagHelpful,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:       "Binding Stone H",
		ID:         48019,
		AuraID:     67739,
		Duration:   time.Second * 20,
		MaxStacks:  8,
		Bonus:      stats.Stats{stats.SpellPower: 74},
		CD:         time.Minute * 2,
		Callback:   core.CallbackOnCastComplete,
		SpellFlags: core.SpellFlagHelpful,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Vengeance of the Forsaken",
		ID:        47881,
		AuraID:    67737,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 215, stats.RangedAttackPower: 215},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:      "Vengeance of the Forsaken H",
		ID:        48020,
		AuraID:    67746,
		Duration:  time.Second * 20,
		MaxStacks: 5,
		Bonus:     stats.Stats{stats.AttackPower: 250, stats.RangedAttackPower: 250},
		CD:        time.Minute * 2,
		Callback:  core.CallbackOnSpellHitDealt,
		ProcMask:  core.ProcMaskMelee,
		Outcome:   core.OutcomeLanded,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:        "Eitrigg's Oath",
		ID:          47882,
		AuraID:      67727,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1265},
		CD:          time.Minute * 2,
		Callback:    core.CallbackOnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})
	helpers.NewStackingStatBonusCD(helpers.StackingStatBonusCD{
		Name:        "Eitrigg's Oath H",
		ID:          48021,
		AuraID:      67741,
		Duration:    time.Second * 20,
		MaxStacks:   5,
		Bonus:       stats.Stats{stats.Armor: 1422},
		CD:          time.Minute * 2,
		Callback:    core.CallbackOnSpellHitTaken,
		Outcome:     core.OutcomeLanded,
		IsDefensive: true,
	})

	core.AddEffectsToTest = true
}
