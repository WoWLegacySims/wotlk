package helpers

import "github.com/WoWLegacySims/wotlk/sim/core"

type CapacitorHandler func(*core.Simulation)

type CapacitorAura struct {
	Aura    core.Aura
	Handler CapacitorHandler
}

// Creates an aura which activates a handler function upon reaching a certain number of stacks.
func MakeCapacitorAura(unit *core.Unit, config CapacitorAura) *core.Aura {
	handler := config.Handler
	config.Aura.OnStacksChange = func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
		if newStacks == aura.MaxStacks {
			handler(sim)
			aura.SetStacks(sim, 0)
		}
	}
	return unit.RegisterAura(config.Aura)
}

type CapacitorDamageEffect struct {
	Name      string
	ID        int32
	MaxStacks int32
	Trigger   core.ProcTrigger

	School     core.SpellSchool
	BasePoints float64
	Die        float64
}

func NewCapacitorDamageEffect(config CapacitorDamageEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		bp := config.BasePoints
		die := config.Die
		damageSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: config.ID},
			SpellSchool: config.School,
			ProcMask:    core.ProcMaskEmpty,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(bp, die), spell.OutcomeMagicHitAndCrit)
			},
		})

		capacitorAura := MakeCapacitorAura(&character.Unit, CapacitorAura{
			Aura: core.Aura{
				Label:     config.Name,
				ActionID:  core.ActionID{ItemID: config.ID},
				Duration:  core.NeverExpires,
				MaxStacks: config.MaxStacks,
			},
			Handler: func(sim *core.Simulation) {
				damageSpell.Cast(sim, character.CurrentTarget)
			},
		})

		config.Trigger.Name = config.Name + " Trigger"
		config.Trigger.Handler = func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			capacitorAura.Activate(sim)
			capacitorAura.AddStack(sim)
		}
		core.MakeProcTriggerAura(&character.Unit, config.Trigger)
	})
}
