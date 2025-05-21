package hunter

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/hunterinfo"
)

func (hunter *Hunter) registerArcaneShotSpell(timer *core.Timer) {
	dbc := hunterinfo.ArcaneShot.GetMaxRank(hunter.Level)
	if dbc == nil {
		return
	}
	dmg, _ := dbc.GetBPDie(0, hunter.Level)

	hasGlyph := hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfArcaneShot)
	var manaMetrics *core.ResourceMetrics
	if hasGlyph {
		manaMetrics = hunter.NewManaMetrics(core.ActionID{ItemID: 42898})
	}

	hunter.ArcaneShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.03*float64(hunter.Talents.Efficiency),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(hunter.Talents.ImprovedArcaneShot),
			},
		},

		BonusCrit: 0 +
			2*float64(hunter.Talents.SurvivalInstincts),
		DamageMultiplierAdditive: 1 +
			.03*float64(hunter.Talents.FerociousInspiration) +
			.05*float64(hunter.Talents.ImprovedArcaneShot),
		DamageMultiplier: 1 *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dmg + 0.15*spell.RangedAttackPower(target)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
			if hasGlyph && result.Landed() && (hunter.SerpentSting.Dot(target).IsActive() || hunter.ScorpidStingAuras.Get(target).IsActive()) {
				hunter.AddMana(sim, 0.2*hunter.ArcaneShot.DefaultCast.Cost, manaMetrics)
			}
			spell.DealDamage(sim, result)
		},
	})
}
