package helpers

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type StackingStatBonusEffect struct {
	Name       string
	ID         int32
	AuraID     int32
	Bonus      stats.Stats
	Duration   time.Duration
	MaxStacks  int32
	Callback   core.AuraCallback
	ProcMask   core.ProcMask
	SpellFlags core.SpellFlag
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
}

func NewStackingStatBonusEffect(config StackingStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ID}
		}
		procAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     config.Name + " Proc",
				ActionID:  auraID,
				Duration:  config.Duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ID},
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			SpellFlags: config.SpellFlags,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			},
		})
	})
}

type StackingStatBonusCD struct {
	Name        string
	ID          int32
	AuraID      int32
	Bonus       stats.Stats
	Duration    time.Duration
	MaxStacks   int32
	CD          time.Duration
	Callback    core.AuraCallback
	ProcMask    core.ProcMask
	SpellFlags  core.SpellFlag
	Outcome     core.HitOutcome
	Harmful     bool
	ProcChance  float64
	IsDefensive bool
}

func NewStackingStatBonusCD(config StackingStatBonusCD) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		auraID := core.ActionID{SpellID: config.AuraID}
		if auraID.IsEmptyAction() {
			auraID = core.ActionID{ItemID: config.ID}
		}
		buffAura := core.MakeStackingAura(character, core.StackingStatAura{
			Aura: core.Aura{
				Label:     config.Name + " Aura",
				ActionID:  auraID,
				Duration:  config.Duration,
				MaxStacks: config.MaxStacks,
			},
			BonusPerStack: config.Bonus,
		})

		core.ApplyProcTriggerCallback(&character.Unit, buffAura, core.ProcTrigger{
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			SpellFlags: config.SpellFlags,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				buffAura.AddStack(sim)
			},
		})

		var sharedTimer *core.Timer
		if config.IsDefensive {
			sharedTimer = character.GetDefensiveTrinketCD()
		} else {
			sharedTimer = character.GetOffensiveTrinketCD()
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: config.ID},
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: config.CD,
				},
				SharedCD: core.Cooldown{
					Timer:    sharedTimer,
					Duration: config.Duration,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})
}
