package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// black book is only ever used pre fight, after which we switch to a real trinket. For this reason we implement it as a
// cooldown and only allow it being cast before combat starts during prepull actions.
func (warlock *Warlock) registerBlackBook() {
	if warlock.Options.Summon == proto.Warlock_Options_NoSummon {
		return
	}

	effectAura := warlock.Pet.NewTemporaryStatsAura("Blessing of the Black Book", core.ActionID{SpellID: 23720},
		stats.Stats{stats.SpellPower: 200, stats.AttackPower: 325, stats.Armor: 1600}, 30*time.Second)

	spell := warlock.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 23720},
		SpellSchool: core.SpellSchoolShadow,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 5 * time.Minute,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return sim.CurrentTime < 0
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			effectAura.Activate(sim)
		},
	})

	warlock.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}
