package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = true
	helpers.NewWeaponDamageProc(754, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Shortsword of Vengeance",
			PPM:  2,
		},
		BasePoints:  30,
		SpellSchool: core.SpellSchoolHoly,
	})

	helpers.NewWeaponDotProc(809, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Bloodrazor",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolPhysical,
		Ticks:       10,
		Interval:    time.Second * 3,
		BasePoints:  12,
	})

	helpers.NewWeaponDamageWithDotProc(870, helpers.WeaponDamageWithDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Fiery War Axe",
		},
		BasePoints:  154,
		Die:         43,
		DotBP:       8,
		Interval:    time.Second * 2,
		Ticks:       4,
		SpellSchool: core.SpellSchoolFire,
	})

	helpers.NewWeaponExtraAttackProc(871, helpers.WeaponExtraAttack{
		WeaponProc: helpers.WeaponProc{
			Name: "Flurry Axe",
			PPM:  3,
		},
	})

	core.NewItemEffect(1168, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionid := core.ActionID{ItemID: 1168}

		leech := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionid,
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcAndDealDamage(sim, target, 35, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					amount := result.Damage
					spell.CalcAndDealHealing(sim, spell.Unit, amount, spell.OutcomeHealing)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Skullflame Shield Leech",
			ActionID:   actionid,
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.03,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				leech.Cast(sim, spell.Unit)
			},
		})

		aoe := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionid,
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				amount := sim.Roll(74, 51)
				for _, target := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, target, amount, spell.OutcomeExpectedMagicHitAndCrit)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Skullflame Shield AOE",
			ActionID:   actionid,
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.03,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aoe.Cast(sim, spell.Unit)
			},
		})
	})

	core.NewItemEffect(1204, func(a core.Agent) {
		character := a.GetCharacter()

		dmg := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 17154},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 3, spell.OutcomeExpectedMagicHitAndCrit)
			},
		})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "The Green Tower",
			ActionID: core.ActionID{ItemID: 1204},
			Duration: time.Second * 30,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.NatureResistance, 50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.NatureResistance, -50)
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeLanded) {
					dmg.Cast(sim, spell.Unit)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "The Green Tower",
			ActionID:   core.ActionID{ItemID: 1204},
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			ProcChance: 0.01,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.Activate(sim)
			},
		})
	})

	helpers.NewWeaponDamageProcWithExtraDamage(1318, helpers.WeaponDamageProcWithExtraDamage{
		WeaponDamageProc: helpers.WeaponDamageProc{
			WeaponProc: helpers.WeaponProc{
				Name: "Night Reaver",
				PPM:  1,
			},
			SpellSchool: core.SpellSchoolShadow,
			BasePoints:  59,
			Die:         31,
			SpellID:     13480,
		},
		ExtraBP:          0,
		ExtraDie:         5,
		ExtraSpellSchool: core.SpellSchoolShadow,
	})
	core.AddEffectsToTest = false
	helpers.NewWeaponDamageProc(810, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Hammer of the Northern Wind",
			PPM:  1.5,
		},
		BasePoints:  19,
		Die:         11,
		SpellSchool: core.SpellSchoolFrost,
	})

	helpers.NewWeaponDamageProc(811, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Axe of the Deep Woods",
			PPM:  1.35,
		},
		BasePoints:  89,
		Die:         37,
		SpellSchool: core.SpellSchoolNature,
	})

	helpers.NewWeaponDotProc(899, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Venon Web Fang",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 3,
		BasePoints:  3,
	})

	helpers.NewWeaponDamageProc(937, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Black Duskwood Staff",
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  109,
		Die:         31,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 942,
		Trigger: core.ProcTrigger{
			Name:       "Freezing Band",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ProcChance: 1.0,
		},
		School:     core.SpellSchoolFrost,
		BasePoints: 50,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 1131,
		Trigger: core.ProcTrigger{
			Name:       "Totem of Infliction",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ProcChance: 1.0,
		},
		School:     core.SpellSchoolShadow,
		BasePoints: 74,
		Die:        51,
	})

	helpers.NewWeaponDamageProc(1263, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Brain Hacker",
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  199,
		Die:         101,
	})

	helpers.NewWeaponDotProc(1265, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Scorpid Sting",
			PPM:  1.3,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 5,
		BasePoints:  13,
	})

	helpers.NewWeaponDamageProc(1387, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Ghoulfang",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  35,
	})

	helpers.NewWeaponDamageProc(1481, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Grimclaw",
			PPM:  1.5,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  30,
	})

	helpers.NewWeaponDamageProcWithExtraDamage(1482, helpers.WeaponDamageProcWithExtraDamage{
		WeaponDamageProc: helpers.WeaponDamageProc{
			WeaponProc: helpers.WeaponProc{
				Name: "Shadowfang",
				PPM:  1.6,
			},
			SpellSchool: core.SpellSchoolShadow,
			BasePoints:  30,
			SpellID:     13440,
		},
		ExtraBP:          3,
		ExtraDie:         5,
		ExtraSpellSchool: core.SpellSchoolShadow,
	})

	helpers.NewWeaponDotProc(1726, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Poison-tipped Bone Spear",
			PPM:  1.5,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 6,
		BasePoints:  30,
	})

	helpers.NewWeaponDamageProc(1728, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Teebu's Blazing Longsword",
			PPM:  1.4,
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  150,
	})

	helpers.NewWeaponDamageProc(1982, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Nightblade",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  124,
		Die:         151,
	})

	helpers.NewWeaponDamageProc(1986, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Gutrender",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  89,
		Die:         25,
	})

	helpers.NewWeaponDamageProc(2000, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Archeus",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolArcane,
		BasePoints:  130,
	})

	helpers.NewWeaponDamageProc(2099, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Dwarven Hand Cannon",
			PPM:  3,
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  32,
		Die:         17,
	})

	helpers.NewWeaponDamageProc(2163, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Shadowblade",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  109,
		Die:         31,
	})

	helpers.NewWeaponDamageProc(2164, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Gut Ripper",
			PPM:  1.8,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  94,
		Die:         27,
	})

	helpers.NewWeaponDamageProc(2205, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Duskbringer",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  59,
		Die:         41,
	})

	helpers.NewWeaponDamageProc(2256, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Skeletal Club",
			PPM:  0.75,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  30,
	})

	helpers.NewWeaponDamageProc(2263, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Phytoblade",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolNature,
		BasePoints:  35,
	})

	helpers.NewWeaponDotProc(2291, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Kang the Decapicator",
		},
		Ticks:       10,
		Interval:    time.Second * 3,
		BasePoints:  56,
		SpellSchool: core.SpellSchoolPhysical,
	})

	helpers.NewWeaponDamageWithDotProc(2299, helpers.WeaponDamageWithDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Burning War Axe",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolFire,
		Ticks:       3,
		Interval:    time.Second * 2,
		BasePoints:  85,
		Die:         25,
		DotBP:       6,
	})

	helpers.NewWeaponDamageProc(2824, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Hurricane",
			PPM:  3,
		},
		SpellSchool: core.SpellSchoolFrost,
		BasePoints:  30,
		Die:         15,
	})

	helpers.NewWeaponDamageProc(2825, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Bow of Searing Arrows",
			PPM:  3,
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  17,
		Die:         11,
	})

	helpers.NewWeaponDamageProc(2912, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Claw of the Shadowmancer",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  35,
	})

	helpers.NewWeaponDamageWithDotProc(2915, helpers.WeaponDamageWithDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Taran Icebreaker",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolFire,
		Ticks:       4,
		Interval:    time.Second * 2,
		BasePoints:  179,
		Die:         41,
		DotBP:       9,
	})

	helpers.NewWeaponDamageProc(2942, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Iron Knuckles",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  4,
	})

	helpers.NewWeaponDamageProcWithExtraDamage(3194, helpers.WeaponDamageProcWithExtraDamage{
		WeaponDamageProc: helpers.WeaponDamageProc{
			WeaponProc: helpers.WeaponProc{
				Name: "Black Malice",
				PPM:  2,
			},
			BasePoints:  54,
			SpellSchool: core.SpellSchoolShadow,
			Die:         31,
			SpellID:     18205,
		},
		ExtraBP:          0,
		ExtraDie:         6,
		ExtraSpellSchool: core.SpellSchoolShadow,
	})

	helpers.NewWeaponDotProc(3336, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Flesh Piercer",
		},
		SpellSchool: core.SpellSchoolPhysical,
		Ticks:       5,
		Interval:    time.Second * 6,
		BasePoints:  6,
	})

	helpers.NewWeaponDamageProc(3822, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Runic Darkblade",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  35,
	})

	helpers.NewWeaponDamageProc(3854, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Frost Tiger Blade",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolFrost,
		BasePoints:  19,
		Die:         11,
	})

	helpers.NewWeaponDotProcWithExtraDamage(4446, helpers.WeaponDotProcWithExtraDamage{
		WeaponDotProc: helpers.WeaponDotProc{
			WeaponProc: helpers.WeaponProc{
				Name: "Blackvenom Blade",
				PPM:  2,
			},
			SpellID:     13518,
			SpellSchool: core.SpellSchoolNature,
			Ticks:       5,
			Interval:    time.Second * 3,
			BasePoints:  5,
		},
		ExtraBP:          0,
		ExtraDie:         7,
		ExtraSpellSchool: core.SpellSchoolShadow,
	})

	helpers.NewWeaponDotProc(4449, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Naraxis' Fang",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 3,
		BasePoints:  6,
	})

	helpers.NewWeaponDamageProc(5182, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Shiver Blade",
			PPM:  0.8,
		},
		SpellSchool: core.SpellSchoolFrost,
		BasePoints:  35,
	})

	helpers.NewWeaponDotProc(5426, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Serpent's Kiss",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 3,
		BasePoints:  7,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID: 12631,
		Trigger: core.ProcTrigger{
			Name:     "Fiery Plate Gauntlets",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee,
			Outcome:  core.OutcomeLanded,
		},
		School:     core.SpellSchoolFire,
		BasePoints: 4,
	})

	helpers.NewWeaponDotProc(5616, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Gutwrencher",
			PPM:  1.5,
		},
		SpellSchool: core.SpellSchoolPhysical,
		Ticks:       10,
		Interval:    time.Second * 3,
		BasePoints:  8,
	})

	helpers.NewWeaponDotProc(5752, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Wyvern Tailspike",
			PPM:  1.5,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 3,
		BasePoints:  6,
	})

	helpers.NewWeaponDamageProc(5756, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Sliverblade",
			PPM:  1.7,
		},
		SpellSchool: core.SpellSchoolFrost,
		BasePoints:  45,
	})

	helpers.NewWeaponDamageProc(5815, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Glacial Stone",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolFrost,
		BasePoints:  75,
	})

	helpers.NewWeaponDamageProc(6220, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Meteor Shard",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  35,
	})

	helpers.NewWeaponDamageProc(6469, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Venomstrike",
			PPM:  3,
		},
		SpellSchool: core.SpellSchoolNature,
		BasePoints:  30,
		Die:         15,
	})

	helpers.NewWeaponDotProc(6472, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Stinging Viper",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 3,
		BasePoints:  7,
	})

	helpers.NewWeaponDotProc(6738, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Bleeding Crescent",
			PPM:  0.9,
		},
		SpellSchool: core.SpellSchoolPhysical,
		Ticks:       5,
		Interval:    time.Second * 6,
		BasePoints:  9,
	})

	helpers.NewWeaponDamageProc(6831, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Black Menace",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  30,
	})

	helpers.NewWeaponDotProc(6904, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Bite of Serra'kis",
			PPM:  2.2,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       10,
		Interval:    time.Second * 2,
		BasePoints:  4,
	})

	helpers.NewWeaponDotProc(6909, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Strike of the Hydra",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       10,
		Interval:    time.Second * 3,
		BasePoints:  7,
		Aura: core.Aura{
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.BonusArmor, -50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.AddStatDynamic(sim, stats.BonusArmor, 50)
			},
		},
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID:         7284,
		School:     core.SpellSchoolFire,
		BasePoints: 14,
		Die:        11,
		Trigger: core.ProcTrigger{
			Name:       "Red Whelp Gloves",
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.05,
		},
	})

	helpers.NewWeaponDamageProcWithExtraDamage(7730, helpers.WeaponDamageProcWithExtraDamage{
		WeaponDamageProc: helpers.WeaponDamageProc{
			WeaponProc: helpers.WeaponProc{
				Name: "Cobalt Crusher",
				PPM:  1.875,
			},
			SpellID:     18204,
			SpellSchool: core.SpellSchoolFrost,
			BasePoints:  109,
			Die:         11,
		},
		ExtraBP:          5,
		ExtraSpellSchool: core.SpellSchoolFrost,
	})

	helpers.NewProcDamageEffect(helpers.ProcDamageEffect{
		ID:         7747,
		School:     core.SpellSchoolShadow,
		BasePoints: 104,
		Die:        71,
		Trigger: core.ProcTrigger{
			Name:       "Vile Protector",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.01,
		},
	})

	helpers.NewWeaponDotProc(7753, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Bloodspiller",
			PPM:  0.9,
		},
		SpellSchool: core.SpellSchoolPhysical,
		Ticks:       10,
		Interval:    time.Second * 3,
		BasePoints:  13,
	})

	helpers.NewWeaponDamageWithDotProc(7959, helpers.WeaponDamageWithDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Blight",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolNature,
		Ticks:       5,
		Interval:    time.Second * 2,
		BasePoints:  50,
		DotBP:       10,
	})

	helpers.NewWeaponDamageProc(8006, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "The Ziggler",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolNature,
		BasePoints:  9,
		Die:         11,
	})

	helpers.NewWeaponDamageProc(8190, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Hanzo Sword",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  75,
	})

	helpers.NewWeaponDotProc(8224, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Silithid Ripper",
			PPM:  1.5,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  9,
		Ticks:       5,
		Interval:    time.Second * 6,
	})

	helpers.NewWeaponDotProc(8225, helpers.WeaponDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Tainted Pierce",
			PPM:  2,
		},
		SpellSchool: core.SpellSchoolShadow,
		BasePoints:  15,
		Ticks:       3,
		Interval:    time.Second,
	})
	core.AddEffectsToTest = true
	core.NewItemEffect(9372, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(9372)
		spell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 11658},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			Dot: core.DotConfig{
				NumberOfTicks: 5,
				TickLength:    time.Second * 3,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 25
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
				Aura: core.Aura{
					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.AddStatDynamic(sim, stats.AttackPower, -30)
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.AddStatDynamic(sim, stats.AttackPower, 30)
					},
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(89, 121)
				result := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if result.Landed() {
					spell.Dot(target).Activate(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Sul'thraze the Laster",
			ActionID: core.ActionID{ItemID: 9372},
			ProcMask: procmask,
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			PPM:      2.2,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				spell.Cast(sim, result.Target)
			},
		})
	})
	core.AddEffectsToTest = false
	helpers.NewWeaponDamageWithDotProc(9386, helpers.WeaponDamageWithDotProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Excavator's Brand",
			PPM:  1,
		},
		SpellID:     13438,
		BasePoints:  40,
		SpellSchool: core.SpellSchoolFire,
		Ticks:       3,
		Interval:    time.Second * 2,
		DotBP:       3,
	})

	helpers.NewWeaponDamageProc(9412, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Galgann's Fireblaster",
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  11,
		Die:         7,
	})

	helpers.NewWeaponDamageProc(9419, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Galgann's Firehammer",
			PPM:  2.2,
		},
		SpellSchool: core.SpellSchoolFire,
		BasePoints:  79,
		Die:         33,
	})

	helpers.NewWeaponDamageProc(9425, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Pendulum of Doom",
			PPM:  1,
		},
		SpellSchool: core.SpellSchoolPhysical,
		BasePoints:  249,
		Die:         101,
	})

	helpers.NewWeaponDamageProc(9446, helpers.WeaponDamageProc{
		WeaponProc: helpers.WeaponProc{
			Name: "Electrocutioner Leg",
			PPM:  1.75,
		},
		SpellSchool: core.SpellSchoolNature,
		BasePoints:  9,
		Die:         11,
	})
	core.AddEffectsToTest = true
}
