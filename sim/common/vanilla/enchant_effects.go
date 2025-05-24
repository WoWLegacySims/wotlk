package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	helpers.AddWeaponDamageEnchant(241, 2)
	helpers.AddWeaponDamageEnchant(250, 1)
	helpers.AddWeaponDamageEnchant(805, 4)
	helpers.AddWeaponDamageEnchant(943, 3)
	helpers.AddWeaponDamageEnchant(963, 7)
	helpers.AddWeaponDamageEnchant(1896, 9)
	helpers.AddWeaponDamageEnchant(1897, 5)

	helpers.AddScope(30, 1)
	helpers.AddScope(32, 2)
	helpers.AddScope(33, 3)
	helpers.AddScope(663, 5)
	helpers.AddScope(664, 7)

	helpers.AddShieldSpike(43, 6042, "Iron Shield Spike", 7, 5)
	helpers.AddShieldSpike(463, 7967, "Mithril Shield Spike", 15, 5)
	helpers.AddShieldSpike(1704, 12645, "Thorium Shield Spike", 19, 11)

	helpers.AddAbsorption(44, 7426, "Minor Absorption", 10, 0.02)
	helpers.AddAbsorption(63, 13538, "Lesser Absorption", 25, 0.05)

	core.NewEnchantEffect(2443, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.FrostSpellPower += 7
	})

	core.NewEnchantEffect(2614, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.ShadowSpellPower += 20
	})

	core.NewEnchantEffect(2615, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.FrostSpellPower += 20
	})

	core.NewEnchantEffect(2616, func(agent core.Agent) {
		agent.GetCharacter().PseudoStats.FireSpellPower += 20
	})

	//Minor Beastslayer
	core.AddWeaponEffect(249, func(agent core.Agent, slot proto.ItemSlot) {
		w := core.Ternary(slot == proto.ItemSlot_ItemSlotMainHand, agent.GetCharacter().AutoAttacks.MH(), agent.GetCharacter().AutoAttacks.OH())
		if agent.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			w.BaseDamageMin += 2
			w.BaseDamageMax += 2
		}
	})

	//Lesser Beastslayer
	core.AddWeaponEffect(853, func(agent core.Agent, slot proto.ItemSlot) {
		w := core.Ternary(slot == proto.ItemSlot_ItemSlotMainHand, agent.GetCharacter().AutoAttacks.MH(), agent.GetCharacter().AutoAttacks.OH())
		if agent.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeBeast {
			w.BaseDamageMin += 6
			w.BaseDamageMax += 6
		}
	})

	// Lesser Elemental Slayer
	core.AddWeaponEffect(854, func(agent core.Agent, slot proto.ItemSlot) {
		w := core.Ternary(slot == proto.ItemSlot_ItemSlotMainHand, agent.GetCharacter().AutoAttacks.MH(), agent.GetCharacter().AutoAttacks.OH())
		if agent.GetCharacter().CurrentTarget.MobType == proto.MobType_MobTypeElemental {
			w.BaseDamageMin += 6
			w.BaseDamageMax += 6
		}
	})

	core.NewEnchantEffect(36, func(a core.Agent) {
		character := a.GetCharacter()

		procMask := character.GetProcMaskForEnchant(36)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 6297},
			SpellSchool:      core.SpellSchoolFire,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := sim.Roll(8, 5)
					baseDamage *= sim.Encounter.AOECapMultiplier()
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Fiery Blaze",
			Callback:   core.CallbackOnSpellHitDealt,
			ProcMask:   procMask,
			Outcome:    core.OutcomeLanded,
			ProcChance: 0.15,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(36, aura)
	})

	core.NewEnchantEffect(803, func(a core.Agent) {
		character := a.GetCharacter()

		procMask := character.GetProcMaskForEnchant(803)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 13897},
			SpellSchool:      core.SpellSchoolFire,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 40, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Fiery Weapon",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      6.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(803, aura)
	})

	core.NewEnchantEffect(912, func(a core.Agent) {
		character := a.GetCharacter()

		procMask := character.GetProcMaskForEnchant(912)

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 13907},
			SpellSchool:      core.SpellSchoolHoly,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				baseDamage := sim.Roll(74, 51)
				spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Demonslaying",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      6.0,
			CustomCheck: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return result.Target.MobType == proto.MobType_MobTypeDemon
			},
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(912, aura)
	})

	core.NewEnchantEffect(1898, func(a core.Agent) {
		character := a.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1898)
		healthMetrics := character.NewHealthMetrics(core.ActionID{SpellID: 20004})

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 20004},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				res := spell.CalcAndDealDamage(sim, target, 30, spell.OutcomeMagicHitAndCrit)
				if res.Landed() {
					character.GainHealth(sim, res.Damage, healthMetrics)
				}
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Lifestealing",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      6.0,
			CustomCheck: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				return result.Target.MobType == proto.MobType_MobTypeDemon
			},
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(1898, aura)
	})

	core.NewEnchantEffect(1899, func(a core.Agent) {
		character := a.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1899)

		modDamage := func(wp *core.Weapon, amount float64) {
			if wp == nil {
				return
			}
			if !wp.SpellSchool.Matches(core.SpellSchoolPhysical) {
				return
			}
			wp.BaseDamageMin += amount
			wp.BaseDamageMax += amount
		}

		debuff := character.CurrentTarget.GetOrRegisterAura(core.Aura{
			Label:    "Unholy Curse",
			ActionID: core.ActionID{SpellID: 20006},
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				modDamage(aura.Unit.AutoAttacks.MH(), -15)
				modDamage(aura.Unit.AutoAttacks.OH(), -15)
				modDamage(aura.Unit.AutoAttacks.Ranged(), -15)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				modDamage(aura.Unit.AutoAttacks.MH(), 15)
				modDamage(aura.Unit.AutoAttacks.OH(), 15)
				modDamage(aura.Unit.AutoAttacks.Ranged(), 15)
			},
		})

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 20006},
			SpellSchool:      core.SpellSchoolShadow,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(43, 13)
				res := spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				if res.Landed() {
					debuff.Activate(sim)
				}
			},
		})

		aura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Unholy Curse",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procMask,
			Outcome:  core.OutcomeLanded,
			PPM:      1.0,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEnchantEffect(1899, aura)
	})

	// ApplyCrusaderEffect will be applied twice if there is two weapons with this enchant.
	//   However, it will automatically overwrite one of them, so it should be ok.
	//   A single application of the aura will handle both mh and oh procs.
	core.NewEnchantEffect(1900, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForEnchant(1900)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		// -4 str per level over 60
		strBonus := 100.0 - 4.0*float64(character.Level-60)
		mhAura := character.NewTemporaryStatsAura("Crusader Enchant MH", core.ActionID{SpellID: 20007, Tag: 1}, stats.Stats{stats.Strength: strBonus}, time.Second*15)
		ohAura := character.NewTemporaryStatsAura("Crusader Enchant OH", core.ActionID{SpellID: 20007, Tag: 2}, stats.Stats{stats.Strength: strBonus}, time.Second*15)

		aura := character.GetOrRegisterAura(core.Aura{
			Label:    "Crusader Enchant",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Crusader") {
					if spell.IsMH() {
						mhAura.Activate(sim)
					} else {
						ohAura.Activate(sim)
					}
				}
			},
		})

		character.ItemSwap.RegisterOnSwapItemForEffectWithPPMManager(1900, 1.0, &ppmm, aura)
	})

	core.NewEnchantEffect(2621, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 0.98
	})
	core.NewEnchantEffect(2613, func(agent core.Agent) {
		character := agent.GetCharacter()
		character.PseudoStats.ThreatMultiplier *= 1.02
	})

}
