package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false

	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "Thunder Capacitor",
		ID:        38072,
		MaxStacks: 4,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2500,
			ActionID: core.ActionID{ItemID: 38072},
		},
		School: core.SpellSchoolNature,
		MinDmg: 1181,
		MaxDmg: 1371,
	})
	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "Reign of the Unliving",
		ID:        47182,
		MaxStacks: 3,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
			ActionID: core.ActionID{ItemID: 47182},
		},
		School: core.SpellSchoolFire,
		MinDmg: 1741,
		MaxDmg: 2023,
	})
	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "Reign of the Unliving H",
		ID:        47188,
		MaxStacks: 3,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
			ActionID: core.ActionID{ItemID: 47188},
		},
		School: core.SpellSchoolFire,
		MinDmg: 1959,
		MaxDmg: 2275,
	})

	core.AddEffectsToTest = true

	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "Reign of the Dead",
		ID:        47316,
		MaxStacks: 3,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
			ActionID: core.ActionID{ItemID: 47316},
		},
		School: core.SpellSchoolFire,
		MinDmg: 1741,
		MaxDmg: 2023,
	})
	helpers.NewCapacitorDamageEffect(helpers.CapacitorDamageEffect{
		Name:      "Reign of the Dead H",
		ID:        47477,
		MaxStacks: 3,
		Trigger: core.ProcTrigger{
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskSpellOrProc | core.ProcMaskWeaponProc | core.ProcMaskSuppressedProc,
			Outcome:  core.OutcomeCrit,
			ICD:      time.Millisecond * 2000,
			ActionID: core.ActionID{ItemID: 47477},
		},
		School: core.SpellSchoolFire,
		MinDmg: 1959,
		MaxDmg: 2275,
	})

	// see various posts around https://web.archive.org/web/20100530203708/http://elitistjerks.com/f78/t39136-combat_mutilate_spreadsheets_updated_3_3_a/p96/#post1518212
	NewItemEffectWithHeroic(func(isHeroic bool) {
		name := "Tiny Abomination in a Jar"
		itemID := int32(50351)
		maxStacks := int32(8)
		if isHeroic {
			name += " H"
			itemID = 50706
			maxStacks = 7
		}

		core.NewItemEffect(itemID, func(agent core.Agent) {
			character := agent.GetCharacter()
			if !character.AutoAttacks.AutoSwingMelee {
				return
			}

			var mhSpell *core.Spell
			var ohSpell *core.Spell
			initSpells := func() {
				mhSpell = character.GetOrRegisterSpell(core.SpellConfig{
					ActionID:         core.ActionID{SpellID: 71433}, // "Manifest Anger"
					SpellSchool:      core.SpellSchoolPhysical,
					ProcMask:         core.ProcMaskMeleeMHSpecial,
					Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
					DamageMultiplier: character.AutoAttacks.MHConfig().DamageMultiplier * 0.5,
					CritMultiplier:   character.AutoAttacks.MHConfig().CritMultiplier,
					ThreatMultiplier: character.AutoAttacks.MHConfig().ThreatMultiplier,
					ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
						baseDamage := character.MHWeaponDamage(sim, spell.MeleeAttackPower())
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
					},
				})

				if character.AutoAttacks.IsDualWielding {
					ohSpell = character.GetOrRegisterSpell(core.SpellConfig{
						ActionID:         core.ActionID{SpellID: 71434}, // "Manifest Anger"
						SpellSchool:      core.SpellSchoolPhysical,
						ProcMask:         core.ProcMaskMeleeOHSpecial,
						Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagNoOnCastComplete,
						DamageMultiplier: character.AutoAttacks.OHConfig().DamageMultiplier * 0.5,
						CritMultiplier:   character.AutoAttacks.OHConfig().CritMultiplier,
						ThreatMultiplier: character.AutoAttacks.OHConfig().ThreatMultiplier,
						ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							baseDamage := character.OHWeaponDamage(sim, spell.MeleeAttackPower())
							spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
						},
					})
				}
			}

			firstProc := core.MainHand

			capacitorAura := helpers.MakeCapacitorAura(&character.Unit, helpers.CapacitorAura{
				Aura: core.Aura{
					Label:     name,
					ActionID:  core.ActionID{SpellID: 71432}, // "Motes of Anger", the aura is either 71406 or 71545 (H) ("Anger Capacitor")
					Duration:  core.NeverExpires,
					MaxStacks: maxStacks,
					OnInit: func(aura *core.Aura, sim *core.Simulation) {
						initSpells()
					},
				},
				Handler: func(sim *core.Simulation) {
					if firstProc == core.MainHand {
						mhSpell.Cast(sim, character.CurrentTarget)
					} else {
						ohSpell.Cast(sim, character.CurrentTarget)
					}
				},
			})

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				Name:       name + " Trigger",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrProc,
				Outcome:    core.OutcomeLanded,
				ActionID:   core.ActionID{ItemID: itemID},
				ProcChance: 0.5,
				Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if spell == mhSpell || spell == ohSpell { // can't proc itself
						return
					}
					if !capacitorAura.IsActive() {
						if spell.ProcMask.Matches(core.ProcMaskMeleeMH | core.ProcMaskProc) {
							firstProc = core.MainHand
						} else {
							firstProc = core.OffHand
						}
					}
					capacitorAura.Activate(sim)
					capacitorAura.AddStack(sim)
				},
			})
		})
	})
}
