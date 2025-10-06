package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
)

func (warlock *Warlock) registerConflagrateSpell() {
	if !warlock.Talents.Conflagrate {
		return
	}

	hasGlyphOfConflag := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfConflagrate)
	warlock.Conflagrate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 17962},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.16,
			Multiplier: 1 - []float64{0, .04, .07, .10}[warlock.Talents.Cataclysm],
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: time.Second * 10,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warlock.Immolate.Dot(target).IsActive() || warlock.ShadowflameDot.Dot(target).IsActive()
		},

		BonusCrit: 0 +
			core.TernaryFloat64(warlock.Talents.Devastation, 5, 0) +
			5*float64(warlock.Talents.FireAndBrimstone),
		DamageMultiplierAdditive: 1 +
			warlock.GrandFirestoneBonus() +
			0.03*float64(warlock.Talents.Emberstorm) +
			0.03*float64(warlock.Talents.Aftermath) +
			0.1*float64(warlock.Talents.ImprovedImmolate) +
			core.TernaryFloat64(warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImmolate), 0.1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetDeathbringerGarb, 2), 0.1, 0) +
			core.TernaryFloat64(warlock.HasSetBonus(ItemSetGuldansRegalia, 4), 0.1, 0),
		CritMultiplier:   warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Conflagrate",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDamage, _ := warlock.getConflagBaseDot(target)

				dot.SnapshotBaseDamage = 0.4 / 3 * baseDamage
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)

				// DoT does not benefit from firestone and also not from spellstone
				dot.Spell.DamageMultiplierAdditive -= warlock.GrandFirestoneBonus()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplierAdditive += warlock.GrandFirestoneBonus()
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage, dot := warlock.getConflagBaseDot(target)
			baseDamage *= 0.6
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if !result.Landed() {
				return
			}

			spell.Dot(target).Apply(sim)

			if !hasGlyphOfConflag {
				dot.Deactivate(sim)
			}
		},
	})
}

func (warlock *Warlock) getConflagBaseDot(target *core.Unit) (float64, *core.Dot) {
	dot := warlock.Immolate.Dot(target)
	ticks := 5.0
	if !dot.IsActive() {
		dot = warlock.ShadowflameDot.Dot(target)
		ticks = 4.0
	}
	return dot.SnapshotBaseDamage * ticks, dot
}
