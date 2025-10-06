package hunter

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/hunterinfo"
)

func (hunter *Hunter) registerAimedShotSpell(timer *core.Timer) {
	if !hunter.Talents.AimedShot {
		return
	}
	dbc := hunterinfo.AimedShot.GetMaxRank(hunter.Level)
	if dbc == nil {
		return
	}
	dmg, _ := dbc.GetBPDie(0, hunter.Level)

	hunter.AimedShot = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  hunterinfo.AimedShot.GetAllIDs(),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
			Multiplier: 1 -
				0.03*float64(hunter.Talents.Efficiency) -
				0.05*float64(hunter.Talents.MasterMarksman),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*10 - core.TernaryDuration(hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfAimedShot), time.Second*2, 0),
			},
		},

		BonusCrit: 0 +
			4*float64(hunter.Talents.ImprovedBarrage) +
			core.TernaryFloat64(hunter.Talents.TrueshotAura && hunter.HasMajorGlyph(proto.HunterMajorGlyph_GlyphOfTrueshotAura), 10, 0),
		DamageMultiplierAdditive: 1 +
			.04*float64(hunter.Talents.Barrage),
		DamageMultiplier: 1 *
			hunter.markedForDeathMultiplier(),
		CritMultiplier:   hunter.critMultiplier(true, true, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.2*spell.RangedAttackPower(target) +
				hunter.AutoAttacks.Ranged().BaseDamage(sim) +
				hunter.AmmoDamageBonus +
				spell.BonusWeaponDamage() +
				dmg
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}
