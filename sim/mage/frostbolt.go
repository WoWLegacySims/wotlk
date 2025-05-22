package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerFrostboltSpell() {
	dbc := mageinfo.Frostbolt.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(1, mage.Level)
	spellCoeff := (dbc.GetCoefficient(1) + 0.05*float64(mage.Talents.EmpoweredFrostbolt)) * dbc.GetLevelPenalty(mage.Level)
	basecost := dbc.BaseCost / 100

	replProcChance := float64(mage.Talents.EnduringWinter) / 3
	var replSrc core.ReplenishmentSource
	if replProcChance > 0 {
		replSrc = mage.Env.Raid.NewReplenishmentSource(core.ActionID{SpellID: 44561})
	}

	mage.Frostbolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:   mageinfo.Frostbolt.GetAllIDs(),
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | core.SpellFlagAPL,
		MissileSpeed: 28,

		ManaCost: core.ManaCostOptions{
			BaseCost: basecost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Second*3 - time.Millisecond*100*time.Duration(mage.Talents.ImprovedFrostbolt+mage.Talents.EmpoweredFrostbolt),
			},
		},

		BonusCrit: 0 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5, 0),
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.ChilledToTheBone) +
			core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfFrostbolt), .05, 0) +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+float64(mage.Talents.IceShards)/3),
		ThreatMultiplier: 1 - (0.1/3)*float64(mage.Talents.FrostChanneling),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
				if replProcChance == 1 || sim.RandomFloat("Enduring Winter") < replProcChance {
					mage.Env.Raid.ProcReplenishment(sim, replSrc)
				}
			})
		},
	})
}
