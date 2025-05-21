package warrior

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warriorinfo"
)

func (warrior *Warrior) registerThunderClapSpell() {
	dbc := warriorinfo.ThunderClap.GetMaxRank(warrior.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, warrior.Level)

	warrior.ThunderClapAuras = warrior.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
		return core.ThunderClapAura(target, warrior.Talents.ImprovedThunderClap)
	})

	warrior.ThunderClap = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskRangedSpecial,
		Flags:       core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost: 20 -
				float64(warrior.Talents.FocusedRage) -
				[]float64{0, 1, 2, 4}[warrior.Talents.ImprovedThunderClap] -
				core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfResonatingPower), 5, 0),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.StanceMatches(BattleStance | DefensiveStance)
		},

		// Cruelty doesn't apply to Thunder Clap
		BonusCrit:        (float64(warrior.Talents.Incite)*5 - float64(warrior.Talents.Cruelty)*1),
		DamageMultiplier: []float64{1.0, 1.1, 1.2, 1.3}[warrior.Talents.ImprovedThunderClap],
		CritMultiplier:   warrior.critMultiplier(none),
		ThreatMultiplier: 1.85,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := bp + 0.12*spell.MeleeAttackPower()
			baseDamage *= sim.Encounter.AOECapMultiplier()

			for _, aoeTarget := range sim.Encounter.TargetUnits {
				result := spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeRangedHitAndCrit)
				if result.Landed() {
					warrior.ThunderClapAuras.Get(aoeTarget).Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{warrior.ThunderClapAuras},
	})
}

func (warrior *Warrior) CanThunderClapIgnoreStance(sim *core.Simulation) bool {
	return warrior.CurrentRage() >= warrior.ThunderClap.DefaultCast.Cost && warrior.ThunderClap.IsReady(sim)
}

func (warrior *Warrior) ShouldThunderClap(sim *core.Simulation, target *core.Unit, filler bool, maintainOnly bool, ignoreStance bool) bool {
	if ignoreStance && !warrior.CanThunderClapIgnoreStance(sim) {
		return false
	} else if !ignoreStance && !warrior.ThunderClap.CanCast(sim, target) {
		return false
	}

	if filler {
		return true
	}

	return maintainOnly &&
		warrior.ThunderClapAuras.Get(target).ShouldRefreshExclusiveEffects(sim, time.Second*2)
}
