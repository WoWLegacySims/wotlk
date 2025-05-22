package warlock

import (
	"math"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerDarkPactSpell() {
	if !warlock.Talents.DarkPact {
		return
	}

	dbc := warlockinfo.DarkPact.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	baseRestore, _ := dbc.GetBPDie(0, warlock.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}
	manaMetrics := warlock.NewManaMetrics(actionID)
	petManaMetrics := warlock.Pet.NewManaMetrics(actionID)
	hasGlyph := warlock.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfLifeTap)

	warlock.DarkPact = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellRanks:  warlockinfo.DarkPact.GetAllIDs(),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),
		FlatThreatBonus:  80,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Glyph activates and applies SP before coeff calculations are done.
			if hasGlyph {
				warlock.GlyphOfLifeTapAura.Activate(sim)
			}

			maxDrain := baseRestore + 0.96*warlock.GetStat(stats.SpellPower)
			actualDrain := math.Min(maxDrain, warlock.Pet.CurrentMana())

			warlock.Pet.SpendMana(sim, actualDrain, petManaMetrics)
			warlock.AddMana(sim, actualDrain, manaMetrics)
		},
	})
}
