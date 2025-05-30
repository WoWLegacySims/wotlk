package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.AddEffectsToTest = false
	core.NewItemEffect(744, func(a core.Agent) {
		character := a.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 12257},
			SpellSchool:      core.SpellSchoolFire,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, target := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, target, 50, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{ItemID: 744},
			SpellSchool: core.SpellSchoolFire,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 30,
					Timer:    character.NewTimer(),
				},
				IgnoreHaste: true,
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				core.StartPeriodicAction(sim, core.PeriodicActionOptions{
					NumTicks: 5,
					Period:   time.Second,
					OnAction: func(sim *core.Simulation) {
						procSpell.Cast(sim, target)
					},
				})
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(8348, func(a core.Agent) {
		character := a.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{ItemID: 8348},
			SpellSchool:      core.SpellSchoolFire,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			ProcMask:         core.ProcMaskEmpty,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 5,
					Timer:    character.NewTimer(),
				},
				IgnoreHaste: true,
			},
			Dot: core.DotConfig{
				NumberOfTicks: 4,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.SnapshotBaseDamage = 10
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(440, 119)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(43660, func(a core.Agent) {
		character := a.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 59197},
			SpellSchool:      core.SpellSchoolFire,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 2,
					Timer:    character.NewTimer(),
				},
				IgnoreHaste: true,
			},
			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, target := range sim.Encounter.TargetUnits {
					dmg := sim.Roll(277, 45)
					spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})

	core.NewItemEffect(43663, func(a core.Agent) {
		character := a.GetCharacter()

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:         core.ActionID{SpellID: 59199},
			SpellSchool:      core.SpellSchoolNature,
			ProcMask:         core.ProcMaskEmpty,
			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultSpellCritMultiplier(),
			ThreatMultiplier: 1,
			Cast: core.CastConfig{
				CD: core.Cooldown{
					Duration: time.Minute * 2,
					Timer:    character.NewTimer(),
				},
				IgnoreHaste: true,
			},
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				dmg := sim.Roll(11, 789)
				spell.CalcAndDealDamage(sim, target, dmg, spell.OutcomeMagicHitAndCrit)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: spell,
			Type:  core.CooldownTypeDPS,
		})
	})
	core.AddEffectsToTest = true
}
