package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(28727, func(a core.Agent) {
		character := a.GetCharacter()

		stackAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     "Enlightenment",
				ActionID:  core.ActionID{SpellID: 32095},
				Duration:  core.NeverExpires,
				MaxStacks: 20,
			},
			BonusPerStack: stats.Stats{stats.MP5: 26},
		})

		useAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{SpellID: 29601},
			Name:     "Pendant of the Violet Eye",
			Callback: core.CallbackOnCastComplete,
			Duration: time.Second * 20,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				var i any = spell.Cost
				_, ok := i.(core.ManaCost)
				if ok && spell.CurCast.Cost > 0 {
					stackAura.Activate(sim)
				}

			},
		})
		useAura.OnExpire = func(aura *core.Aura, sim *core.Simulation) {
			stackAura.Deactivate(sim)
		}

		character.AddMajorCooldown(core.MajorCooldown{
			Type: core.CooldownTypeMana,
			Spell: character.GetOrRegisterSpell(core.SpellConfig{
				ActionID: core.ActionID{ItemID: 28727},
				ProcMask: core.ProcMaskEmpty,
				Cast: core.CastConfig{
					CD: core.Cooldown{
						Timer:    character.NewTimer(),
						Duration: time.Minute * 2,
					},
				},
				ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
					useAura.Activate(sim)
				},
			}),
		})

	})

	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:      "The Night Blade",
		ID:        31331,
		AuraID:    38307,
		Duration:  time.Second * 10,
		MaxStacks: 3,
		Bonus:     stats.Stats{stats.ArmorPenetration: 62},
		Callback:  core.CallbackOnSpellHitDealt,
		Weapon:    true,
	})

	core.NewItemEffect(31856, func(agent core.Agent) {
		character := agent.GetCharacter()

		apAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Aura of the Crusader",
			ActionID:  core.ActionID{SpellID: 39439},
			Duration:  time.Second * 10,
			MaxStacks: 20,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.AddStatsDynamic(sim, stats.Stats{stats.AttackPower: 6, stats.RangedAttackPower: 6}.Multiply(float64(newStacks-oldStacks)))
			},
		})

		spAura := character.GetOrRegisterAura(core.Aura{
			Label:     "Aura of the Crusader",
			ActionID:  core.ActionID{SpellID: 39441},
			Duration:  time.Second * 10,
			MaxStacks: 10,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.AddStatsDynamic(sim, stats.Stats{stats.SpellPower: 8}.Multiply(float64(newStacks-oldStacks)))
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 31856},
			Name:     "Darkmoon Card: Crusade",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeLanded,
			Handler: func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
					apAura.Activate(sim)
					apAura.AddStack(sim)
				}
				if spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
					spAura.Activate(sim)
					spAura.AddStack(sim)
				}
			},
		})
	})

	core.NewItemEffect(31857, func(agent core.Agent) {
		character := agent.GetCharacter()

		aura := character.GetOrRegisterAura(core.Aura{
			Label:     "Aura of Wrath",
			ActionID:  core.ActionID{SpellID: 39443},
			Duration:  time.Second * 10,
			MaxStacks: 20,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
				character.AddStatsDynamic(sim, stats.Stats{stats.MeleeCrit: 17, stats.SpellCrit: 17}.Multiply(float64(newStacks-oldStacks)))
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if result.Outcome.Matches(core.OutcomeCrit) {
					aura.Deactivate(sim)
				}
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 31857},
			Name:     "Darkmoon Card: Wrath",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: core.ProcMaskDirect,
			Outcome:  core.OutcomeLanded,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				if !result.Outcome.Matches(core.OutcomeCrit) {
					return
				}
				aura.Activate(sim)
				aura.AddStack(sim)
			},
		})
	})

	helpers.NewStackingStatBonusEffect(helpers.StackingStatBonusEffect{
		Name:       "Blackened Naaru Sliver",
		ID:         34427,
		Duration:   time.Second * 20,
		MaxStacks:  10,
		Bonus:      stats.Stats{stats.AttackPower: 44, stats.RangedAttackPower: 44},
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
	core.AddEffectsToTest = true
}
