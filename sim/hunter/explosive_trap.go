package hunter

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func (hunter *Hunter) registerExplosiveTrapSpell(timer *core.Timer) {
	hasGlyph := hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfExplosiveTrap)
	bonusPeriodicDamageMultiplier := .10 * float64(hunter.Talents.TrapMastery)

	hunter.ExplosiveTrap = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 49067},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.19,
			Multiplier: 1 - 0.2*float64(hunter.Talents.Resourcefulness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		DamageMultiplierAdditive: 1 +
			.02*float64(hunter.Talents.TNT),
		CritMultiplier:   hunter.critMultiplier(false, false, false),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "Explosive Trap",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				baseDamage := 90 + 0.1*dot.Spell.RangedAttackPower(target)
				dot.Spell.DamageMultiplierAdditive += bonusPeriodicDamageMultiplier
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if hasGlyph {
						dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeRangedHitAndCritNoBlock)
					} else {
						dot.Spell.CalcAndDealPeriodicDamage(sim, aoeTarget, baseDamage, dot.Spell.OutcomeRangedHit)
					}
				}
				dot.Spell.DamageMultiplierAdditive -= bonusPeriodicDamageMultiplier
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.CurrentTime < 0 {
				// Traps only last 30s.
				if sim.CurrentTime < -time.Second*30 {
					return
				}

				// If using this on prepull, the trap effect will go off when the fight starts
				// instead of immediately.
				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: 0,
					OnAction: func(sim *core.Simulation) {
						for _, aoeTarget := range sim.Encounter.TargetUnits {
							baseDamage := sim.Roll(523, 671) + 0.1*spell.RangedAttackPower(aoeTarget)
							baseDamage *= sim.Encounter.AOECapMultiplier()
							spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeRangedHitAndCritNoBlock)
						}
						hunter.ExplosiveTrap.AOEDot().Apply(sim)
					},
				})
			} else {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := sim.Roll(523, 671) + 0.1*spell.RangedAttackPower(aoeTarget)
					baseDamage *= sim.Encounter.AOECapMultiplier()
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeRangedHitAndCritNoBlock)
				}
				hunter.ExplosiveTrap.AOEDot().Apply(sim)
			}
		},
	})

	timeToTrapWeave := time.Millisecond * time.Duration(hunter.Options.TimeToTrapWeaveMs)
	halfWeaveTime := timeToTrapWeave / 2
	hunter.TrapWeaveSpell = hunter.RegisterSpell(core.SpellConfig{
		ActionID: hunter.ExplosiveTrap.ActionID.WithTag(1),
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagNoMetrics | core.SpellFlagNoLogs | core.SpellFlagAPL,

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.ExplosiveTrap.CanCast(sim, target)
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if sim.CurrentTime < 0 {
				hunter.mayMoveAt = sim.CurrentTime
			}

			// Assume we started running after the most recent ranged auto, so that time
			// can be subtracted from the run in.
			reachLocationAt := hunter.mayMoveAt + halfWeaveTime
			layTrapAt := max(reachLocationAt, sim.CurrentTime)
			doneAt := layTrapAt + halfWeaveTime

			hunter.AutoAttacks.DelayRangedUntil(sim, doneAt+time.Millisecond*500)

			if layTrapAt == sim.CurrentTime {
				hunter.ExplosiveTrap.Cast(sim, target)
				if doneAt > hunter.GCD.ReadyAt() {
					hunter.GCD.Set(doneAt)
				}
			} else {
				// Make sure the GCD doesn't get used while we're waiting.
				hunter.WaitUntil(sim, doneAt)

				core.StartDelayedAction(sim, core.DelayedActionOptions{
					DoAt: layTrapAt,
					OnAction: func(sim *core.Simulation) {
						hunter.GCD.Reset()
						hunter.ExplosiveTrap.Cast(sim, target)
						if doneAt > hunter.GCD.ReadyAt() {
							hunter.GCD.Set(doneAt)
						}
					},
				})
			}
		},
	})
}
