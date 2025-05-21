package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) registerThunderstormSpell() {
	if !shaman.Talents.Thunderstorm {
		return
	}
	dbc := shamaninfo.Thunderstorm.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, shaman.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(shaman.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}
	manaMetrics := shaman.NewManaMetrics(actionID)

	manaRestore := 0.08
	if shaman.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfThunderstorm) {
		manaRestore = 0.1
	}

	shaman.Thunderstorm = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		Flags:       core.SpellFlagAPL,
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		BonusHit:         float64(shaman.Talents.ElementalPrecision),
		BonusCrit:        core.TernaryFloat64(shaman.Talents.CallOfThunder, 5, 0),
		DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion),
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.spellThreatMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics)

			if shaman.thunderstormInRange {
				dmgFromSP := coef * spell.SpellPower()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					baseDamage := sim.Roll(bp, die) + dmgFromSP
					baseDamage *= sim.Encounter.AOECapMultiplier()
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			}
		},
	})
}
