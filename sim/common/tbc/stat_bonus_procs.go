package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Quagmirran's Eye",
		ID:         27683,
		AuraID:     33297,
		Bonus:      stats.Stats{stats.MeleeHaste: 320, stats.SpellHaste: 320},
		Duration:   time.Second * 6,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Blackout Truncheon",
		ID:       27901,
		AuraID:   33489,
		Bonus:    stats.Stats{stats.MeleeHaste: 132, stats.SpellHaste: 132},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.9,
		ICD:      time.Second * 45,
		Weapon:   true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Hourglass of the Unraveller",
		ID:         28034,
		AuraID:     60066,
		Bonus:      stats.Stats{stats.AttackPower: 300, stats.RangedAttackPower: 300},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeCrit,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Scarab of the Infinite Cycle",
		ID:         28190,
		AuraID:     60061,
		Bonus:      stats.Stats{stats.MeleeHaste: 320, stats.SpellHaste: 320},
		Duration:   time.Second * 6,
		Callback:   core.CallbackOnHealDealt,
		Outcome:    core.OutcomeLanded,
		ProcMask:   core.ProcMaskDirect,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Greatsword of Forlorn Visions",
		ID:       28367,
		AuraID:   34199,
		Bonus:    stats.Stats{stats.BonusArmor: 2750},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Shiffar's Nexus-Horn",
		ID:         28418,
		AuraID:     34321,
		Bonus:      stats.Stats{stats.SpellPower: 225},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.2,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Lionheart Champion",
		ID:         28429,
		AuraID:     34513,
		Bonus:      stats.Stats{stats.Strength: 100},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.07,
		Weapon:     true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Lionheart Executioner",
		ID:         28430,
		AuraID:     34513,
		Bonus:      stats.Stats{stats.Strength: 100},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.07,
		Weapon:     true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Drakefirst Hammer",
		ID:       28437,
		AuraID:   21165,
		Bonus:    stats.Stats{stats.MeleeHaste: 212, stats.SpellHaste: 212},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
		Weapon:   true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Dragonmaw",
		ID:       28438,
		AuraID:   21165,
		Bonus:    stats.Stats{stats.MeleeHaste: 212, stats.SpellHaste: 212},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
		Weapon:   true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Dragonstrike",
		ID:       28439,
		AuraID:   21165,
		Bonus:    stats.Stats{stats.MeleeHaste: 212, stats.SpellHaste: 212},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
		Weapon:   true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Masquerade Gown",
		ID:         28578,
		AuraID:     34584,
		Bonus:      stats.Stats{stats.Spirit: 145},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellHealing | core.ProcMaskSpellDamage,
		ProcChance: 0.1,
		ICD:        time.Second * 30,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Robe of the Elder Scribes",
		ID:         28602,
		AuraID:     34597,
		Bonus:      stats.Stats{stats.SpellPower: 130},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		Duration:   time.Second * 10,
		ProcChance: 0.2,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Eye of Magtheridon",
		ID:       28789,
		AuraID:   34747,
		Bonus:    stats.Stats{stats.SpellPower: 170},
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskSpellDamage,
		Outcome:  core.OutcomeMiss,
		Duration: time.Second * 10,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Dragonspine Trophy",
		ID:       28830,
		AuraID:   34774,
		Bonus:    stats.Stats{stats.MeleeHaste: 325, stats.SpellHaste: 325},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
		ICD:      time.Second * 20,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:        "Band of the Eternal Defender",
		ID:          29297,
		AuraID:      35077,
		Bonus:       stats.Stats{stats.Armor: 800},
		Duration:    time.Second * 10,
		Callback:    core.CallbackOnSpellHitTaken,
		SpellSchool: core.SpellSchoolPhysical,
		Outcome:     core.OutcomeLanded,
		ProcChance:  0.03,
		ICD:         time.Second * 60,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Band of the Eternal Champion",
		ID:       29301,
		AuraID:   35080,
		Bonus:    stats.Stats{stats.AttackPower: 160, stats.RangedAttackPower: 160},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Band of the Eternal Sage",
		ID:         29305,
		AuraID:     35083,
		Bonus:      stats.Stats{stats.SpellPower: 95},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 60,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "The Bladefist",
		ID:       29348,
		AuraID:   35131,
		Bonus:    stats.Stats{stats.MeleeHaste: 180, stats.SpellHaste: 180},
		Duration: time.Second * 10,
		Outcome:  core.OutcomeLanded,
		PPM:      3,
		Weapon:   true,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Heartrazor",
		ID:       29962,
		AuraID:   36041,
		Bonus:    stats.Stats{stats.AttackPower: 270, stats.RangedAttackPower: 270},
		Duration: time.Second * 10,
		Outcome:  core.OutcomeLanded,
		PPM:      2.2,
		Weapon:   true,
	})

	core.NewItemEffect(30090, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(30090)

		aura := character.NewTemporaryStatsAura("World Breaker", core.ActionID{SpellID: 36111}, stats.Stats{stats.MeleeCrit: 900}, time.Second*4)
		aura.OnSpellHitDealt = func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}
			aura.Deactivate(sim)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "World Breaker",
			ActionID: core.ActionID{ItemID: 30090},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procmask,
			Outcome:  core.OutcomeLanded,
			PPM:      1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.Activate(sim)
			},
		})
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Sextant of Unstable Currents",
		ID:         30626,
		AuraID:     38348,
		Bonus:      stats.Stats{stats.SpellPower: 190},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpell,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.2,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Tsunami Talisman",
		ID:         30627,
		AuraID:     42083,
		Bonus:      stats.Stats{stats.AttackPower: 340, stats.RangedAttackPower: 340},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:        "Bullwark Of Azzinoth",
		ID:          32375,
		AuraID:      40407,
		Bonus:       stats.Stats{stats.Armor: 2000},
		Duration:    time.Second * 10,
		Callback:    core.CallbackOnSpellHitTaken,
		SpellSchool: core.SpellSchoolPhysical,
		Outcome:     core.OutcomeLanded,
		ProcChance:  0.02,
		ICD:         time.Second * 60,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Madness of the Betrayer",
		ID:       32505,
		AuraID:   42083,
		Bonus:    stats.Stats{stats.ArmorPenetration: 42},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      3,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:       "Shard of Contempt",
		ID:         34472,
		Bonus:      stats.Stats{stats.AttackPower: 230, stats.RangedAttackPower: 230},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})

	helpers.NewProcStatBonusEffect(helpers.ProcStatBonusEffect{
		Name:     "Commendation of Kael'Thas",
		ID:       34473,
		AuraID:   45057,
		Bonus:    stats.Stats{stats.Dodge: 152},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitTaken,
		ProcMask: core.ProcMaskMelee,
		Outcome:  core.OutcomeLanded,
		ICD:      time.Second * 30,
		CustomCheck: func(aura *core.Aura, _ *core.Simulation, _ *core.Spell, _ *core.SpellResult) bool {
			return aura.Unit.CurrentHealthPercent() <= 0.35
		},
	})
	core.AddEffectsToTest = true
}
