package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerArcaneBarrageSpell() {
	if !mage.Talents.ArcaneBarrage {
		return
	}

	dbc := mageinfo.ArcaneBarrage.GetMaxRank(mage.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, mage.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(mage.Level)

	mage.ArcaneBarrage = mage.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:   mageinfo.ArcaneBarrage.GetAllIDs(),
		SpellSchool:  core.SpellSchoolFrost,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage | BarrageSpells | core.SpellFlagAPL,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.18,
			Multiplier: 1 *
				(1 - .01*float64(mage.Talents.ArcaneFocus)) *
				core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneBarrage), .8, 1),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 3,
			},
		},

		BonusHit:         float64(mage.Talents.ArcaneFocus),
		DamageMultiplier: 1 + .04*float64(mage.Talents.TormentTheWeak),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + coef*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
