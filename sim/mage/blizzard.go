package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerBlizzardSpell() {
	dbc := mageinfo.Blizzard.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	var improvedBlizzardProcApplication *core.Spell
	if mage.Talents.ImprovedBlizzard > 0 {
		auras := mage.NewEnemyAuraArray(func(unit *core.Unit, _ int32) *core.Aura {
			return unit.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 12488},
				Label:    "Improved Blizzard",
				Duration: time.Millisecond * 1500,
			})
		})
		improvedBlizzardProcApplication = mage.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 12488},
			ProcMask: core.ProcMaskProc,
			Flags:    SpellFlagMage | core.SpellFlagNoLogs,
			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				auras.Get(target).Activate(sim)
			},
		})
	}

	blizzardTickSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: dbc.Effects[0].TriggerSpell},
		SpellSchool:      core.SpellSchoolFrost,
		ProcMask:         core.ProcMaskSpellDamage,
		Flags:            SpellFlagMage,
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3),
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := bp + coef*spell.SpellPower()
			damage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, damage, spell.OutcomeMagicHitAndCrit)

				if improvedBlizzardProcApplication != nil {
					improvedBlizzardProcApplication.Cast(sim, aoeTarget)
				}
			}
		},
	})

	mage.Blizzard = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost: 0.74,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Blizzard",
			},
			NumberOfTicks:       8,
			TickLength:          time.Second * 1,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				blizzardTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	})
}
