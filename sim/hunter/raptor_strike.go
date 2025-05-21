package hunter

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/hunterinfo"
)

func (hunter *Hunter) registerRaptorStrikeSpell() {
	dbc := hunterinfo.RaptorStrike.GetMaxRank(hunter.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, hunter.Level)

	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.04,
			Multiplier: 1 - 0.2*float64(hunter.Talents.Resourcefulness),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		BonusCrit:        float64(hunter.Talents.SavageStrikes) * 10,
		DamageMultiplier: 1,
		CritMultiplier:   hunter.critMultiplier(false, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bp +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(_ *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	//if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveAutosOnly || !hunter.RaptorStrike.IsReady(sim) || hunter.CurrentMana() < hunter.RaptorStrike.DefaultCast.Cost {
	//	return nil
	//}
	return mhSwingSpell
}
