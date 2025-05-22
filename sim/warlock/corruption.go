package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerCorruptionSpell() {
	dbc := warlockinfo.Corruption.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, warlock.Level)
	coef := (0.2 + 0.02*float64(warlock.Talents.EmpoweredCorruption) + 0.01*float64(warlock.Talents.EverlastingAffliction)) * dbc.GetLevelPenalty(warlock.Level)
	ticks := dbc.Duration / dbc.Effects[0].AuraPeriod
	canCrit := warlock.Talents.Pandemic

	warlock.Corruption = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  warlockinfo.Corruption.GetAllIDs(),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagHauntSE | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   dbc.BaseCost / 100,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCrit: 0 +
			3*float64(warlock.Talents.Malediction) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDarkCovensRegalia, 2), 5, 0),
		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion) +
			0.02*float64(warlock.Talents.ImprovedCorruption) +
			core.TernaryFloat64(warlock.Talents.SiphonLife, 0.05, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetGuldansRegalia, 4), 0.1, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, 1),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Corruption",
			},
			NumberOfTicks:       ticks,
			TickLength:          time.Second * 3,
			AffectedByCastSpeed: warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfQuickDecay),

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = bp + coef*dot.Spell.SpellPower()
				if !isRollover {
					attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
					dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
					dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if canCrit {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, dot.OutcomeExpectedMagicSnapshotCrit)
			} else {
				baseDmg := bp + coef*spell.SpellPower()
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicCrit)
			}
		},
	})
}
