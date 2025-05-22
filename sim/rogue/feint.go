package rogue

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/rogueinfo"
)

func (rogue *Rogue) registerFeintSpell() {
	dbc := rogueinfo.Feint.GetMaxRank(rogue.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, rogue.Level)

	rogue.Feint = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  rogueinfo.Feint.GetAllIDs(),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		EnergyCost: core.EnergyCostOptions{
			Cost: core.TernaryFloat64(rogue.HasMajorGlyph(proto.RogueMajorGlyph_GlyphOfFeint), 0, 20),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    rogue.NewTimer(),
				Duration: time.Second * 10,
			},
			IgnoreHaste: true,
		},
		FlatThreatBonus:  bp,
		DamageMultiplier: 0,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			rogue.BreakStealth(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
		},
	})
}
