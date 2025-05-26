package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	core.NewItemEffect(24114, func(agent core.Agent) { // Braided Eternium Chain
		agent.GetCharacter().PseudoStats.BonusDamage += 5
	})

	core.NewItemEffect(27529, func(a core.Agent) {
		character := a.GetCharacter()
		metrics := character.NewHealthMetrics(core.ActionID{SpellID: 33089})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Vigilance of the Colossus",
			ActionID: core.ActionID{SpellID: 33089},
			Duration: time.Second * 20,
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeBlock) {
					character.GainHealth(sim, 120, metrics)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 27529},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 2,
						Timer:    character.NewTimer(),
					},
				},
				ProcMask: core.ProcMaskEmpty,
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					aura.Activate(sim)
				},
			}),
		})
	})

	core.NewItemEffect(27770, func(a core.Agent) {
		character := a.GetCharacter()
		shieldStrength := 0.0
		metrics := character.NewHealthMetrics(core.ActionID{SpellID: 39228})

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Vigilance of the Colossus",
			ActionID: core.ActionID{SpellID: 39228},
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				shieldStrength = 1150
			},
			OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Damage > 0 {
					absorb := min(result.Damage, shieldStrength, 68)
					shieldStrength -= absorb
					character.GainHealth(sim, absorb, metrics)

					if shieldStrength == 0 {
						aura.Deactivate(sim)
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 27770},
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Duration: time.Minute * 2,
						Timer:    character.NewTimer(),
					},
				},
				ProcMask: core.ProcMaskEmpty,
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					aura.Activate(sim)
				},
			}),
		})
	})

	core.NewItemEffect(27896, func(a core.Agent) {
		character := a.GetCharacter()
		metrics := character.NewManaMetrics(core.ActionID{SpellID: 15603})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Alembic of Infernal Power",
			ActionID:   core.ActionID{ItemID: 27896},
			Callback:   core.CallbackOnSpellHitTaken,
			ProcMask:   core.ProcMaskDirect,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.02,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				character.AddMana(sim, 260, metrics)
			},
		})
	})

	core.NewItemEffect(27920, func(a core.Agent) {
		character := a.GetCharacter()
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 33510})
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33510})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Mark of Conquest",
			ActionID: core.ActionID{ItemID: 27920},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			Outcome:  core.OutcomeLanded,
			PPM:      2,
			ICD:      time.Second * 25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMelee) {
					amount := sim.Roll(89, 31)
					character.GainHealth(sim, amount, healthMetrics)
				} else {
					amount := sim.Roll(127, 45)
					character.AddMana(sim, amount, manaMetrics)
				}

			},
		})
	})

	core.NewItemEffect(27921, func(a core.Agent) {
		character := a.GetCharacter()
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 33510})
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33510})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Mark of Conquest",
			ActionID: core.ActionID{ItemID: 27921},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskMelee | core.ProcMaskRanged,
			Outcome:  core.OutcomeLanded,
			PPM:      2,
			ICD:      time.Second * 25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMelee) {
					amount := sim.Roll(89, 31)
					character.GainHealth(sim, amount, healthMetrics)
				} else {
					amount := sim.Roll(127, 45)
					character.AddMana(sim, amount, manaMetrics)
				}

			},
		})
	})

	core.NewItemEffect(27922, func(a core.Agent) {
		character := a.GetCharacter()
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33511})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mark of Conquest",
			ActionID:   core.ActionID{ItemID: 27922},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 17,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(127, 45)
				character.AddMana(sim, amount, manaMetrics)
			},
		})
	})

	core.NewItemEffect(27924, func(a core.Agent) {
		character := a.GetCharacter()
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33511})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mark of Conquest",
			ActionID:   core.ActionID{ItemID: 27924},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 17,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(127, 45)
				character.AddMana(sim, amount, manaMetrics)
			},
		})
	})

	core.NewItemEffect(27926, func(a core.Agent) {
		character := a.GetCharacter()
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33522})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mark of Conquest",
			ActionID:   core.ActionID{ItemID: 27926},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(127, 45)
				character.AddMana(sim, amount, manaMetrics)
			},
		})
	})

	core.NewItemEffect(27927, func(a core.Agent) {
		character := a.GetCharacter()
		manaMetrics := character.NewManaMetrics(core.ActionID{SpellID: 33522})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Mark of Conquest",
			ActionID:   core.ActionID{ItemID: 27927},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   core.ProcMaskSpellDamage,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			ICD:        time.Second * 25,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(127, 45)
				character.AddMana(sim, amount, manaMetrics)
			},
		})
	})

	core.NewItemEffect(28370, func(a core.Agent) {
		character := a.GetCharacter()

		regen := 0.15 - float64(max(character.Level-70, 0))*0.005

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Bangle of Endless Blessings",
			ActionID: core.ActionID{SpellID: 38334},
			Duration: time.Second * 15,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SpiritRegenRateCasting += regen
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SpiritRegenRateCasting -= regen
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Bangle of Endless Blessings",
			ActionID:   core.ActionID{SpellID: 38334},
			Callback:   core.CallbackOnCastComplete,
			ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
			ProcChance: 0.1,
			ICD:        time.Second * 50,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})

		activeAura := character.NewTemporaryStatsAura("Endless Blessings", core.ActionID{SpellID: 34210}, stats.Stats{stats.Spirit: 130}, time.Second*20)

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: character.RegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{SpellID: 34210},
				ProcMask: core.ProcMaskEmpty,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					activeAura.Activate(sim)
				},
			}),
			Type: core.CooldownTypeMana,
		})
	})

	core.NewItemEffect(29182, func(a core.Agent) {
		character := a.GetCharacter()
		procMask := character.GetProcMaskForItem(29182)

		aura := character.CurrentTarget.GetOrRegisterAura(core.Aura{
			Label:    "Temporal Rift",
			ActionID: core.ActionID{SpellID: 35353},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyAttackSpeed(sim, 1/1.1)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.MultiplyAttackSpeed(sim, 1.01)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Riftmaker",
			ActionID: core.ActionID{ItemID: 29182},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				aura.Activate(sim)
			},
		})

	})

	core.NewItemEffect(29996, func(agent core.Agent) { // Rod of the Sun King
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(29996)
		pppm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		actionID := core.ActionID{ItemID: 29996}

		var resourceMetricsRage *core.ResourceMetrics
		var resourceMetricsEnergy *core.ResourceMetrics
		if character.HasRageBar() {
			resourceMetricsRage = character.NewRageMetrics(actionID)
		}
		if character.HasEnergyBar() {
			resourceMetricsEnergy = character.NewEnergyMetrics(actionID)
		}

		character.GetOrRegisterAura(core.Aura{
			Label:    "Rod of the Sun King",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if pppm.Proc(sim, spell.ProcMask, "Rod of the Sun King") {
					switch spell.Unit.GetCurrentPowerBar() {
					case core.RageBar:
						spell.Unit.AddRage(sim, 5, resourceMetricsRage)
					case core.EnergyBar:
						spell.Unit.AddEnergy(sim, 10, resourceMetricsEnergy)
					}
				}
			},
		})
	})

	core.NewItemEffect(30620, func(a core.Agent) {
		character := a.GetCharacter()

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeSurvival,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID:    core.ActionID{ItemID: 30620},
				SpellSchool: core.SpellSchoolNature,
				ProcMask:    core.ProcMaskSpellHealing,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				Hot: core.DotConfig{
					Aura: core.Aura{
						Label:    "Spyglass of the Hidden Fleet",
						ActionID: core.ActionID{SpellID: 38325},
					},
					SelfOnly:      true,
					NumberOfTicks: 4,
					TickLength:    time.Second * 3,
					OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
						dot.SnapshotBaseDamage = 325
					},
					OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
						dot.CalcAndDealPeriodicSnapshotHealing(sim, target, dot.OutcomeTick)
					},
				},
			}),
		})
	})

	core.NewItemEffect(30892, func(agent core.Agent) { //Beast-tamer's Shoulders
		for _, pet := range agent.GetCharacter().Pets {
			if pet.IsGuardian() {
				continue // not sure if this applies to guardians.
			}
			pet.PseudoStats.DamageDealtMultiplier *= 1.03
			pet.AddStat(stats.MeleeCrit, pet.CritRatingPerCritChance*2)
		}
	})

	core.NewItemEffect(31193, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(31193)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 24585},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(48, 54) + spell.SpellPower()
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Blade of Unquenched Thirst Trigger",
			ActionID:   core.ActionID{ItemID: 31193},
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.02,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	core.NewItemEffect(31322, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(31322)
		metrics := character.NewManaMetrics(core.ActionID{SpellID: 38284})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "The Hammer of Destiniy",
			ActionID: core.ActionID{ItemID: 31322},
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procmask,
			Outcome:  core.OutcomeLanded,
			PPM:      2.5,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				amount := sim.Roll(169, 61)
				character.AddMana(sim, amount, metrics)
			},
		})
	})

	core.NewItemEffect(31328, func(agent core.Agent) { //Beast-tamer's Shoulders
		for _, pet := range agent.GetCharacter().Pets {
			if pet.IsGuardian() {
				continue // not sure if this applies to guardians.
			}
			pet.AddStats(stats.Stats{stats.AttackPower: 70, stats.Armor: 490, stats.Stamina: 52})
		}
	})

	helpers.NewWeaponExtraAttackProc(31332, helpers.WeaponExtraAttack{
		WeaponProc: helpers.WeaponProc{
			Name: "Blinkstrike",
			PPM:  1.54,
		},
	})

	core.NewItemEffect(32262, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(32262)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		procSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 40291},
			SpellSchool: core.SpellSchoolShadow,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 20, spell.OutcomeMagicHitAndCrit)
			},
		})

		procAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Siphon Essence",
			ActionID: core.ActionID{SpellID: 40291},
			Duration: time.Second * 6,
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}

				procSpell.Cast(sim, result.Target)
			},
		})

		character.GetOrRegisterAura(core.Aura{
			Label:    "Syphon of the Nathrezim",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Syphon Of The Nathrezim") {
					procAura.Activate(sim)
				}
			},
		})
	})

}
