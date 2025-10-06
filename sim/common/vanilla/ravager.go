package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
)

func init() {
	core.NewItemEffect(7717, func(a core.Agent) {
		character := a.GetCharacter()
		procmask := character.GetProcMaskForItem(7717)

		actionID := core.ActionID{SpellID: 9632}

		var spellOH *core.Spell

		if character.AutoAttacks.IsDualWielding {
			spellOH = character.RegisterSpell(core.SpellConfig{
				ActionID:    actionID.WithTag(2),
				SpellSchool: core.SpellSchoolPhysical,
				ProcMask:    core.ProcMaskMeleeOHSpecial,
				Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

				DamageMultiplier: 1,
				CritMultiplier:   character.DefaultMeleeCritMultiplier(),
				ThreatMultiplier: 1,
			})
		}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID:    actionID.WithTag(1),
			SpellSchool: core.SpellSchoolPhysical,
			ProcMask:    core.ProcMaskMeleeMHSpecial,
			Flags:       core.SpellFlagChanneled | core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

			DamageMultiplier: 1,
			CritMultiplier:   character.DefaultMeleeCritMultiplier(),
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				IsAOE: true,
				Aura: core.Aura{
					Label: "Ravager Bladestorm",
				},
				NumberOfTicks: 3,
				TickLength:    time.Second * 3,
				OnTick: func(sim *core.Simulation, _ *core.Unit, dot *core.Dot) {
					spell := dot.Spell
					for _, target := range sim.Encounter.TargetUnits {
						baseDamage := spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + 5
						spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

					}

					if spellOH != nil {
						for _, target := range sim.Encounter.TargetUnits {
							baseDamage := (spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) + 5) * 0.5
							spellOH.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)
						}
					}
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				spell.AOEDot().Apply(sim)
			},
		})

		core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID: core.ActionID{ItemID: 7717},
			Name:     "Ravager",
			Callback: core.CallbackOnSpellHitDealt,
			ProcMask: procmask,
			Outcome:  core.OutcomeLanded,
			PPM:      0.7,
			Handler: func(sim *core.Simulation, _ *core.Spell, result *core.SpellResult) {
				if character.ChanneledDot != nil {
					return
				}
				character.AutoAttacks.StopMeleeUntil(sim, time.Second*9, false)
				spell.Cast(sim, &character.Unit)
			},
		})
	})

}
