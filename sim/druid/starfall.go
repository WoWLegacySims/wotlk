package druid

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

// We register two spells to apply two different dot effects and get two entries in Damage/Detailed results
func (druid *Druid) registerStarfallSpell() {
	if !druid.Talents.Starfall {
		return
	}
	dbc := druidinfo.Starfall.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}
	spellId := dbc.SpellID
	dbc = druidinfo.StarfallAura.GetByID(dbc.Effects[0].TriggerSpell)
	if dbc == nil {
		panic(fmt.Sprintf("No Aura found for Starfall %d", spellId))
	}
	dbc = druidinfo.StarfallDirect.GetByID(int32(dbc.Effects[0].BasePoints) + 1)
	if dbc == nil {
		panic(fmt.Sprintf("No Direct Spell found for Starfall %d", spellId))
	}
	dbcSplash := druidinfo.StarfallSplash.GetByID(dbc.Effects[1].TriggerSpell)
	if dbcSplash == nil {
		panic(fmt.Sprintf("No Splash Spell found for Starfall %d", spellId))
	}
	bpDir, dieDir := dbc.GetBPDie(0, druid.Level)
	splashdmg, _ := dbcSplash.GetBPDie(0, druid.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(druid.Level)
	coefSplash := dbcSplash.GetCoefficient(0) * dbcSplash.GetLevelPenalty(druid.Level)

	numberOfTicks := core.TernaryInt32(druid.Env.GetNumTargets() > 1, 20, 10)
	tickLength := time.Second

	starfallTickSpell := druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: dbc.SpellID},
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSuppressedProc,
		Flags:            SpellFlagNaturesGrace,
		BonusCrit:        2 * float64(druid.Talents.NaturesMajesty),
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bpDir, dieDir) + coef*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})

	druid.Starfall = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellRanks:  druidinfo.Starfall.GetAllIDs(),
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL | SpellFlagOmenTrigger,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.35,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (90 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfStarfall), 30, 0)),
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Starfall",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				starfallTickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				druid.StarfallSplash.Dot(target).Apply(sim)
			}
		},
	})

	starfallSplashTickSpell := druid.RegisterSpell(Any, core.SpellConfig{
		ActionID:         core.ActionID{SpellID: dbcSplash.SpellID},
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskSuppressedProc,
		BonusCrit:        2 * float64(druid.Talents.NaturesMajesty),
		DamageMultiplier: 1 * (1 + core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfFocus), 0.1, 0)),
		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := splashdmg + coefSplash*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	druid.StarfallSplash = druid.RegisterSpell(Any, core.SpellConfig{
		ActionID: core.ActionID{SpellID: dbcSplash.SpellID},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "StarfallSplash",
			},
			NumberOfTicks: numberOfTicks,
			TickLength:    tickLength,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				starfallSplashTickSpell.Cast(sim, target)
			},
		},
	})
}
