package tbc

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {
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
			CustomCheck: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) bool {
				var i any = spell.Cost
				_, ok := i.(core.ManaCost)
				return ok && spell.CurCast.Cost > 0
			},
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				stackAura.Activate(sim)
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
}
