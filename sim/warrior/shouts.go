package warrior

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warriorinfo"
)

const ShoutExpirationThreshold = time.Second * 3

func (warrior *Warrior) makeShoutSpellHelper(actionID core.ActionID, allyAuras core.AuraArray) *core.Spell {
	return warrior.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL | core.SpellFlagHelpful,

		RageCost: core.RageCostOptions{
			Cost: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, aura := range allyAuras {
				if aura != nil {
					aura.Activate(sim)
				}
			}
		},

		RelatedAuras: []core.AuraArray{allyAuras},
	})
}

func (warrior *Warrior) registerShouts() {
	dbc := warriorinfo.BattleShout.GetMaxRank(warrior.Level)
	solarian := core.TernaryFloat64(warrior.HasTrinketEquipped(30466), 110, 0)
	if dbc != nil {
		warrior.BattleShout = warrior.makeShoutSpellHelper(core.ActionID{SpellID: dbc.SpellID}, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
			return core.BattleShoutAura(unit, warrior.Talents.CommandingPresence, warrior.Talents.BoomingVoice, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfBattle), solarian)
		}))
	}
	dbc = warriorinfo.CommandingShout.GetMaxRank(warrior.Level)
	if dbc != nil {
		warrior.CommandingShout = warrior.makeShoutSpellHelper(core.ActionID{SpellID: dbc.SpellID}, warrior.NewAllyAuraArray(func(unit *core.Unit) *core.Aura {
			return core.CommandingShoutAura(unit, warrior.Talents.CommandingPresence, warrior.Talents.BoomingVoice, warrior.HasMinorGlyph(proto.WarriorMinorGlyph_GlyphOfCommand))
		}))
	}
}
