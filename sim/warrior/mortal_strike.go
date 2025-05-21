package warrior

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warriorinfo"
)

func (warrior *Warrior) registerMortalStrikeSpell(cdTimer *core.Timer) {
	dbc := warriorinfo.MortalStrike.GetMaxRank(warrior.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(1, warrior.Level)

	if !warrior.Talents.MortalStrike {
		return
	}

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   30,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second*6 - time.Millisecond*[]time.Duration{0, 300, 700, 1000}[warrior.Talents.ImprovedMortalStrike],
			},
		},

		BonusCrit: core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerBattlegear, 4), 10, 0),
		DamageMultiplier: 1 +
			[]float64{0, 0.03, 0.06, 0.1}[warrior.Talents.ImprovedMortalStrike] +
			core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfMortalStrike), 0.1, 0) +
			core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 4), 0.05, 0),
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bp +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
