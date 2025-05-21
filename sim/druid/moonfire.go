package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

func (druid *Druid) registerMoonfireSpell() {
	dbc := druidinfo.Moonfire.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}
	basecost := dbc.BaseCost / 100

	bp, die := dbc.GetBPDie(1, druid.Level)

	coef := dbc.GetCoefficient(1) * dbc.GetLevelPenalty(druid.Level)

	dotDmg := dbc.Effects[0].BasePoints + 1
	dotCoef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(druid.Level)

	numTicks := dbc.Duration/int32(dbc.Effects[0].AuraPeriod) +
		core.TernaryInt32(druid.Talents.NaturesSplendor, 1, 0) +
		core.TernaryInt32(druid.HasSetBonus(ItemSetThunderheartRegalia, 2), 1, 0)

	starfireBonusCrit := float64(druid.Talents.ImprovedInsectSwarm)
	dotCanCrit := druid.HasSetBonus(ItemSetMalfurionsRegalia, 2)

	baseDamageMultiplier := 1 +
		0.05*float64(druid.Talents.ImprovedMoonfire) +
		[]float64{0.0, 0.03, 0.06, 0.1}[druid.Talents.Moonfury]

	malusInitialDamageMultiplier := core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.9, 0)

	bonusPeriodicDamageMultiplier := 0 +
		0.01*float64(druid.Talents.Genesis) +
		core.TernaryFloat64(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMoonfire), 0.75, 0)

	druid.Moonfire = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagNaturesGrace | SpellFlagOmenTrigger | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   basecost,
			Multiplier: 1 - 0.03*float64(druid.Talents.Moonglow),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusCrit:        float64(druid.Talents.ImprovedMoonfire) * 5,
		DamageMultiplier: baseDamageMultiplier - malusInitialDamageMultiplier,

		CritMultiplier:   druid.BalanceCritMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Moonfire",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.BonusCrit += starfireBonusCrit
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.Starfire.BonusCrit -= starfireBonusCrit
				},
			},
			NumberOfTicks: numTicks,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Spell.DamageMultiplier = baseDamageMultiplier + bonusPeriodicDamageMultiplier
				dot.SnapshotBaseDamage = dotDmg + dotCoef*dot.Spell.SpellPower()
				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(attackTable)
				dot.Spell.DamageMultiplier = baseDamageMultiplier - malusInitialDamageMultiplier
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				if dotCanCrit {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeSnapshotCrit)
				} else {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				druid.ExtendingMoonfireStacks = 3
				dot := spell.Dot(target)
				dot.NumberOfTicks = numTicks
				dot.Apply(sim)
			}
			spell.DealDamage(sim, result)
		},
	})
}
