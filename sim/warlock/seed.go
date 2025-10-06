package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerSeedSpell() {
	dbc := warlockinfo.SeedofCorruption.GetMaxRank(warlock.Level)
	dbcExp := warlockinfo.SeedofCorruptionExplosion.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, warlock.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(warlock.Level)

	bpToExp, _ := dbc.GetBPDie(1, warlock.Level)
	coefToExp := dbc.GetCoefficient(1)

	bpExp, dieExp := dbcExp.GetBPDie(0, warlock.Level)
	coefExp := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(warlock.Level)

	if warlock.HasSetBonus(ItemSetOblivionRaiment, 4) {
		bp += 180
	}

	seedExplosion := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbcExp.SpellID},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagHauntSE | core.SpellFlagNoLogs,

		BonusCrit: 0 +
			float64(warlock.Talents.ImprovedCorruption),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion),
		CritMultiplier:   warlock.DefaultSpellCritMultiplier(),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDmg := (sim.Roll(bpExp, dieExp) + coefExp*spell.SpellPower()) * sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDmg, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	seedDamage := 0.0
	warlock.SeedDamageTracker = make([]float64, len(warlock.Env.AllUnits))
	trySeedPop := func(sim *core.Simulation, target *core.Unit, dmg float64) {
		warlock.SeedDamageTracker[target.UnitIndex] += dmg
		if warlock.SeedDamageTracker[target.UnitIndex] > seedDamage {
			warlock.Seed.Dot(target).Deactivate(sim)
			seedExplosion.Cast(sim, target)
		}
	}

	warlock.Seed = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:   warlockinfo.SeedofCorruption.GetAllIDs(),
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskEmpty,
		Flags:        core.SpellFlagHauntSE | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.34,
			Multiplier: 1 - 0.02*float64(warlock.Talents.Suppression),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2000,
			},
		},

		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery) +
			0.01*float64(warlock.Talents.Contagion) +
			core.TernaryFloat64(warlock.Talents.SiphonLife, 0.05, 0),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Seed",
				OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					trySeedPop(sim, aura.Unit, result.Damage)
				},
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					warlock.SeedDamageTracker[aura.Unit.UnitIndex] = 0
				},
			},

			NumberOfTicks: 6,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = bp + coef*dot.Spell.SpellPower()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
				seedDamage = bpToExp + coefToExp*dot.Spell.SpellPower()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					// seed is mutually exclusive with corruption
					warlock.Corruption.Dot(target).Deactivate(sim)

					if warlock.Options.DetonateSeed {
						seedExplosion.Cast(sim, target)
					} else {
						spell.Dot(target).Apply(sim)
					}
				}
			})
		},
	})
}
