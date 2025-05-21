package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func init() {

	core.NewItemEffect(12590, func(agent core.Agent) {
		character := agent.GetCharacter()

		procMask := character.GetProcMaskForItem(12590)
		ppmm := character.AutoAttacks.NewPPMManager(1.0, procMask)

		effectAura := character.NewTemporaryStatsAura("Felstriker Proc", core.ActionID{SpellID: 16551}, stats.Stats{stats.MeleeCrit: 100 * core.CritRatingPerCritChance}, time.Second*3)

		character.GetOrRegisterAura(core.Aura{
			Label:    "Felstriker",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() {
					return
				}

				if ppmm.Proc(sim, spell.ProcMask, "Felstriker") {
					effectAura.Activate(sim)
				}
			},
		})
	})
}
