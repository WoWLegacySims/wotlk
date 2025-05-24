package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.NewItemEffect(833, func(a core.Agent) {
		character := a.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 17712},
			SpellSchool:      core.SpellSchoolPhysical,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ProcMask:         core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 30,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				amount := sim.Roll(299, 401)
				spell.CalcAndDealHealing(sim, spell.Unit, amount, spell.OutcomeHealingCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	core.NewItemEffect(4262, func(a core.Agent) {
		character := a.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 9163},
			SpellSchool:      core.SpellSchoolHoly,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ProcMask:         core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 5,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				amount := sim.Roll(449, 101)
				spell.CalcAndDealHealing(sim, spell.Unit, amount, spell.OutcomeHealingCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	core.NewItemEffect(4264, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasRageBar() || character.HasEnergyBar()) {
			return
		}

		var ragemetrics *core.ResourceMetrics
		if character.HasRageBar() {
			ragemetrics = character.NewRageMetrics(core.ActionID{ItemID: 4264})
		}
		var energymetrics *core.ResourceMetrics
		if character.HasEnergyBar() {
			energymetrics = character.NewEnergyMetrics(core.ActionID{ItemID: 4264})
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 9174},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 20,
					Timer:    character.NewTimer(),
				},
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				if character.HasRageBar() {
					character.AddRage(sim, 30, ragemetrics)
				}
				if character.HasEnergyBar() {
					character.AddEnergy(sim, 30, energymetrics)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(6972, func(a core.Agent) {
		character := a.GetCharacter()
		if !(character.HasRageBar()) {
			return
		}

		var ragemetrics *core.ResourceMetrics
		if character.HasRageBar() {
			ragemetrics = character.NewRageMetrics(core.ActionID{ItemID: 6972})
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 70537},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 60,
					Timer:    character.NewTimer(),
				},
				DefaultCast: core.Cast{
					GCD: core.GCDDefault,
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				character.AddRage(sim, 30, ragemetrics)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(5079, func(a core.Agent) {
		character := a.GetCharacter()

		aura := character.CurrentTarget.GetOrRegisterAura(core.Aura{
			Label:    "Cold Eye",
			ActionID: core.ActionID{SpellID: 1139},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyAttackSpeed(sim, 1/1.05)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyAttackSpeed(sim, 1.05)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 1139},
			SpellSchool: core.SpellSchoolFrost,
			ProcMask:    core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 5,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	helpers.NewProcStatDebuffEffect(helpers.ProcStatDebuffEffect{
		Name:     "Dazzling Longsword",
		ID:       869,
		AuraID:   13752,
		Debuff:   stats.Stats{stats.BonusArmor: -100},
		Duration: time.Second * 30,
		Callback: core.CallbackOnSpellHitDealt,
		PPM:      1.6,
		Outcome:  core.OutcomeLanded,
		Weapon:   true,
	})

	helpers.NewProcStatDebuffEffect(helpers.ProcStatDebuffEffect{
		Name:     "Howling Blade",
		ID:       6331,
		AuraID:   13490,
		Debuff:   stats.Stats{stats.AttackPower: -30, stats.RangedAttackPower: -30},
		Duration: time.Second * 30,
		Callback: core.CallbackOnSpellHitDealt,
		PPM:      1.6,
		Outcome:  core.OutcomeLanded,
		Weapon:   true,
	})

	core.NewItemEffect(1447, func(a core.Agent) {
		character := a.GetCharacter()

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Ring of Saviors",
			ActionID: core.ActionID{ItemID: 1447},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.BonusArmor, 300)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.BonusArmor, -300)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 1447},
			SpellSchool: core.SpellSchoolPhysical,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 30,
					Timer:    character.NewTimer(),
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				aura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeSurvival,
		})
	})

	helpers.NewProcStatDebuffEffect(helpers.ProcStatDebuffEffect{
		Name:     "Sword of Decay",
		ID:       1727,
		AuraID:   13528,
		Debuff:   stats.Stats{stats.AttackPower: -20},
		Duration: time.Second * 30,
		Callback: core.CallbackOnSpellHitDealt,
		PPM:      1.5,
		Outcome:  core.OutcomeLanded,
		Weapon:   true,
	})

	helpers.NewProcStatDebuffEffect(helpers.ProcStatDebuffEffect{
		Name:     "Phantom Blade",
		ID:       7691,
		AuraID:   9806,
		Debuff:   stats.Stats{stats.BonusArmor: -100},
		Duration: time.Second * 20,
		Callback: core.CallbackOnSpellHitDealt,
		Outcome:  core.OutcomeLanded,
		PPM:      1.5,
		Weapon:   true,
	})

	core.NewItemEffect(11815, func(agent core.Agent) {
		character := agent.GetCharacter()
		if !character.AutoAttacks.AutoSwingMelee {
			return
		}

		config := *character.AutoAttacks.MHConfig()
		config.ActionID = core.ActionID{ItemID: 11815}
		handOfJusticeSpell := character.GetOrRegisterSpell(config)

		procChance := 0.02
		if character.Level > 60 {
			procChance *= 1 - (float64(character.Level-60) / 30)
		}

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Hand of Justice",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskMelee,
			Outcome:    core.OutcomeLanded,
			ICD:        time.Second * 2,
			ProcChance: procChance,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				spell.Unit.AutoAttacks.MaybeReplaceMHSwing(sim, handOfJusticeSpell).Cast(sim, result.Target)
			},
		})
	})

	helpers.NewProcHealEffect(helpers.ProcHealEffect{
		ID:         7939,
		School:     core.SpellSchoolPhysical,
		BasePoints: 59,
		Die:        41,
		Trigger: core.ProcTrigger{
			Name:       "Truesilver Breastplate",
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskDirect,
			Harmful:    true,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.03,
		},
	})

	core.NewItemEffect(19019, func(agent core.Agent) { //Thunderfury
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(19019)

		procActionID := core.ActionID{SpellID: 21992}

		singleTargetSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(1),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 0.5,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)
			},
		})

		makeDebuffAura := func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				Label:    "Thunderfury",
				ActionID: procActionID,
				Duration: time.Second * 12,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, -25)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					target.AddStatDynamic(sim, stats.NatureResistance, 25)
				},
			})
		}

		numHits := min(5, character.Env.GetNumTargets())
		debuffAuras := make([]*core.Aura, len(character.Env.Encounter.TargetUnits))
		for i, target := range character.Env.Encounter.TargetUnits {
			debuffAuras[i] = makeDebuffAura(target)
		}

		bounceSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    procActionID.WithTag(2),
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskEmpty,

			ThreatMultiplier: 1,
			FlatThreatBonus:  63,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				curTarget := target
				for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
					result := spell.CalcDamage(sim, curTarget, 0, spell.OutcomeMagicHit)
					if result.Landed() {
						debuffAuras[target.Index].Activate(sim)
					}
					spell.DealDamage(sim, result)
					curTarget = sim.Environment.NextTargetUnit(curTarget)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Thunderfury",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      6,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				singleTargetSpell.Cast(sim, result.Target)
				bounceSpell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(21625, func(agent core.Agent) { // Scarab Brooch
		character := agent.GetCharacter()

		shieldSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 26470},
			SpellSchool: core.SpellSchoolNature,
			ProcMask:    core.ProcMaskSpellHealing,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Shield: core.ShieldConfig{
				Aura: core.Aura{
					Label:    "Scarab Brooch Shield",
					Duration: time.Second * 30,
				},
			},
		})

		activeAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Persistent Shield",
			ActionID: core.ActionID{SpellID: 26467},
			Callback: core.CallbackOnHealDealt,
			Duration: time.Second * 30,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				shieldSpell.Shield(result.Target).Apply(sim, result.Damage*0.15)
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 21625},
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				activeAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})
}
