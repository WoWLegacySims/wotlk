package priest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/priestinfo"
)

func (priest *Priest) getMindSearMiseryCoefficient() float64 {
	return 0.286 * (1 + 0.05*float64(priest.Talents.Misery))
}

func (priest *Priest) getMindSearBaseConfig() core.SpellConfig {
	return core.SpellConfig{
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskProc,
		BonusHit:    float64(priest.Talents.ShadowFocus),
		BonusCrit:   float64(priest.Talents.MindMelt) * 2,
		DamageMultiplier: 1 +
			0.02*float64(priest.Talents.Darkness) +
			0.01*float64(priest.Talents.TwinDisciplines),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),
		CritMultiplier:   priest.DefaultSpellCritMultiplier(),
	}
}

func (priest *Priest) getMindSearTickSpell(numTicks int32, id int32) (*core.Spell, float64, float64) {
	dbc := priestinfo.MindSearDamage.GetByID(id)
	if dbc == nil {
		panic(fmt.Sprintf("No Mind Sear found with SpellID %d", id))
	}
	bp, die := dbc.GetBPDie(0, priest.Level)

	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))
	miseryCoeff := priest.getMindSearMiseryCoefficient() * dbc.GetLevelPenalty(priest.Level)

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: dbc.SpellID}.WithTag(numTicks)
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		damage := sim.Roll(bp, die) + miseryCoeff*spell.SpellPower()
		result := spell.CalcAndDealDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

		if result.Landed() {
			priest.AddShadowWeavingStack(sim)
		}
		if result.DidCrit() && hasGlyphOfShadow {
			priest.ShadowyInsightAura.Activate(sim)
		}
	}
	return priest.GetOrRegisterSpell(config), bp, die
}

func (priest *Priest) newMindSearSpell(numTicksIdx int32) *core.Spell {
	dbc := priestinfo.MindSear.GetMaxRank(priest.Level)
	if dbc == nil {
		return nil
	}

	numTicks := numTicksIdx
	flags := core.SpellFlagChanneled | core.SpellFlagNoMetrics
	if numTicksIdx == 0 {
		numTicks = 5
		flags |= core.SpellFlagAPL
	}

	miseryCoeff := priest.getMindSearMiseryCoefficient()
	mindSearTickSpell, bp, die := priest.getMindSearTickSpell(numTicksIdx, dbc.Effects[0].TriggerSpell)

	config := priest.getMindSearBaseConfig()
	config.ActionID = core.ActionID{SpellID: dbc.SpellID}.WithTag(numTicksIdx)
	config.SpellRanks = priestinfo.MindSear.GetAllIDs()
	config.Flags = flags
	config.ManaCost = core.ManaCostOptions{
		BaseCost:   0.28,
		Multiplier: 1 - 0.05*float64(priest.Talents.FocusedMind),
	}
	config.Cast = core.CastConfig{
		DefaultCast: core.Cast{
			GCD: core.GCDDefault,
		},
	}
	config.Dot = core.DotConfig{
		Aura: core.Aura{
			Label: "MindSear-" + strconv.Itoa(int(numTicksIdx)),
		},
		NumberOfTicks:       numTicks,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				if aoeTarget != target {
					mindSearTickSpell.Cast(sim, aoeTarget)
					mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
				}
			}
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
		result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
		if result.Landed() {
			spell.Dot(target).Apply(sim)
			mindSearTickSpell.SpellMetrics[target.UnitIndex].Casts += 1
		}
	}
	config.ExpectedTickDamage = func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
		baseDamage := sim.Roll(bp, die) + miseryCoeff*spell.SpellPower()
		return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicCrit)
	}
	return priest.GetOrRegisterSpell(config)
}
