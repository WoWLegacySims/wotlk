package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

func (druid *Druid) registerTigersFurySpell() {
	dbc := druidinfo.TigersFury.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}

	actionID := core.ActionID{SpellID: dbc.SpellID}
	energyMetrics := druid.NewEnergyMetrics(actionID)
	instantEnergy := 20.0 * float64(druid.Talents.KingOfTheJungle)

	dmgBonus, _ := dbc.GetBPDie(0, druid.Level)
	cdReduction := core.TernaryDuration(druid.HasSetBonus(ItemSetDreamwalkerBattlegear, 4), time.Second*3, 0)

	druid.TigersFuryAura = druid.RegisterAura(core.Aura{
		Label:     "Tiger's Fury Aura",
		AuraRanks: druidinfo.TigersFury.GetAllIDs(),
		ActionID:  actionID,
		Duration:  6 * time.Second,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage += dmgBonus
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			druid.PseudoStats.BonusDamage -= dmgBonus
		},
	})

	spell := druid.RegisterSpell(Cat, core.SpellConfig{
		ActionID:   actionID,
		SpellRanks: druidinfo.TigersFury.GetAllIDs(),
		Flags:      core.SpellFlagAPL,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second*30 - cdReduction,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !druid.BerserkAura.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.AddEnergy(sim, instantEnergy, energyMetrics)

			druid.TigersFuryAura.Activate(sim)
		},
	})

	druid.TigersFury = spell
}
