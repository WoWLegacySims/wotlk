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
