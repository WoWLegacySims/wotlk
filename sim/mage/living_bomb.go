package mage

import (
	"fmt"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerLivingBombSpell() {
	if !mage.Talents.LivingBomb {
		return
	}
	dbc := mageinfo.LivingBomb.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	expId, _ := dbc.GetBPDie(1, mage.Level)
	dbcExp := mageinfo.LivingBombExplosion.GetByID(int32(expId))
	if dbcExp == nil {
		panic(fmt.Sprintf("No Explosion Spell found for Living Bomb %d", dbc.SpellID))
	}
	bp, _ := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)
	bpExp, _ := dbcExp.GetBPDie(0, mage.Level)
	coefExp := dbcExp.GetCoefficient(0) * dbcExp.GetLevelPenalty(mage.Level)

	livingBombExplosionSpell := mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbcExp.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | HotStreakSpells,

		BonusCrit: 0 +
			2*float64(mage.Talents.WorldInFlames) +
			2*float64(mage.Talents.CriticalMass),
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bpExp + coefExp*spell.SpellPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})

	onTick := func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
		dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
	}
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfLivingBomb) {
		onTick = func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
		}
	}

	mage.LivingBomb = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 55360},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		// WorldInFlames doesn't apply to DoT component.
		BonusCrit: 0 +
			2*float64(mage.Talents.CriticalMass),
		DamageMultiplierAdditive: 1 +
			.02*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "LivingBomb",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					livingBombExplosionSpell.Cast(sim, aura.Unit)
				},
			},

			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = bp + coef*dot.Spell.SpellPower()
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: onTick,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
