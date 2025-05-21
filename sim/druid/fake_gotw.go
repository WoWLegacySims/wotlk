package druid

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

// This is 'fake' because it doesnt actually account for any actual buff updating
// this is only used as a 'clearcast fisher' spell
func (druid *Druid) registerFakeGotw() {
	dbc := druidinfo.GiftoftheWild.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}
	baseCost := core.TernaryFloat64(druid.HasMinorGlyph(proto.DruidMinorGlyph_GlyphOfTheWild), 0.5, 1) * dbc.BaseCost

	druid.GiftOfTheWild = druid.RegisterSpell(Humanoid|Moonkin|Tree, core.SpellConfig{
		ActionID: core.ActionID{SpellID: dbc.SpellID},
		Flags:    SpellFlagOmenTrigger | core.SpellFlagHelpful | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   baseCost,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})
}
