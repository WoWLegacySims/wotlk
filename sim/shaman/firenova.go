package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) registerFireNovaSpell() {
	dbc := shamaninfo.FireNovaDamage.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, shaman.Level)
	coef := dbc.GetCoefficient(0) * dbc.GetLevelPenalty(shaman.Level)
	fireNovaGlyphCDReduction := core.TernaryInt32(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFireNova), 3, 0)
	impFireNovaCDReduction := shaman.Talents.ImprovedFireNova * 2
	fireNovaCooldown := 10 - fireNovaGlyphCDReduction - impFireNovaCDReduction

	shaman.FireNova = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagFocusable | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.22,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * time.Duration(fireNovaCooldown),
			},
		},

		BonusHit:         float64(shaman.Talents.ElementalPrecision),
		DamageMultiplier: 1 + float64(shaman.Talents.CallOfFlame)*0.05 + float64(shaman.Talents.ImprovedFireNova)*0.1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: shaman.spellThreatMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dmgFromSP := coef * spell.SpellPower()
			for _, aoeTarget := range sim.Encounter.TargetUnits {
				baseDamage := sim.Roll(bp, die) + dmgFromSP
				baseDamage *= sim.Encounter.AOECapMultiplier()
				spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
			}
		},
	})
}

func (shaman *Shaman) IsFireNovaCastable(sim *core.Simulation) bool {
	return shaman.FireNova.IsReady(sim) && shaman.Totems.Fire > 0
}
