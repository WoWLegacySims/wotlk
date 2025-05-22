package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/shamaninfo"
)

func (shaman *Shaman) newTotemSpellConfig(baseCost float64, spellID int32) core.SpellConfig {
	return core.SpellConfig{
		ActionID: core.ActionID{SpellID: spellID},
		Flags:    SpellFlagTotem | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: baseCost,
			Multiplier: 1 -
				0.05*float64(shaman.Talents.TotemicFocus) -
				0.02*float64(shaman.Talents.MentalQuickness),
		},
	}
}

func (shaman *Shaman) registerWrathOfAirTotemSpell() {
	if shaman.Level < 64 {
		return
	}
	config := shaman.newTotemSpellConfig(0.11, 3738)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WrathOfAirTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerWindfuryTotemSpell() {
	if shaman.Level < 32 {
		return
	}
	config := shaman.newTotemSpellConfig(0.11, 8512)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[AirTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.WindfuryTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerManaSpringTotemSpell() {
	dbc := shamaninfo.ManaSpringTotem.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	config := shaman.newTotemSpellConfig(0.04, dbc.SpellID)
	config.SpellRanks = shamaninfo.ManaSpringTotem.GetAllIDs()
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.ManaSpringTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerHealingStreamTotemSpell() {
	dbc := shamaninfo.HealingStreamTotem.GetMaxRank(shaman.Level)
	dbcEffect := shamaninfo.HealingStream.GetMaxRank(shaman.Level)
	if dbc == nil || dbcEffect == nil {
		return
	}

	bp, _ := dbcEffect.GetBPDie(0, shaman.Level)

	config := shaman.newTotemSpellConfig(0.03, dbc.SpellID)
	config.SpellRanks = shamaninfo.HealingStreamTotem.GetAllIDs()
	hsHeal := shaman.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 52042},
		SpellSchool:      core.SpellSchoolNature,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            core.SpellFlagHelpful | core.SpellFlagNoOnCastComplete,
		DamageMultiplier: 1 + (.02 * float64(shaman.Talents.Purification)) + 0.15*float64(shaman.Talents.RestorativeTotems),
		CritMultiplier:   1,
		ThreatMultiplier: 1 - (float64(shaman.Talents.HealingGrace) * 0.05),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			healing := bp + spell.HealingPower(target)*0.0827
			spell.CalcAndDealHealing(sim, target, healing, spell.OutcomeHealing)
		},
	})
	config.Hot = core.DotConfig{
		Aura: core.Aura{
			Label: "HealingStreamHot",
		},
		NumberOfTicks: 150,
		TickLength:    time.Second * 2,
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			hsHeal.Cast(sim, target)
		},
	}
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		shaman.TotemExpirations[WaterTotem] = sim.CurrentTime + time.Second*300
		for _, agent := range shaman.Party.Players {
			spell.Hot(&agent.GetCharacter().Unit).Activate(sim)
		}
	}
	shaman.HealingStreamTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTotemOfWrathSpell() {
	dbc := shamaninfo.TotemofWrath.GetMaxRank(shaman.Level)
	if dbc == nil || !shaman.Talents.TotemOfWrath {
		return
	}
	config := shaman.newTotemSpellConfig(0.05, dbc.SpellID)
	config.SpellRanks = shamaninfo.TotemofWrath.GetAllIDs()
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*300
		shaman.applyToWDebuff(sim)
	}
	shaman.TotemOfWrath = shaman.RegisterSpell(config)
}

func (shaman *Shaman) applyToWDebuff(sim *core.Simulation) {
	for _, target := range sim.Encounter.TargetUnits {
		auraDef := core.TotemOfWrathDebuff(target)
		auraDef.Activate(sim)
	}
}

func (shaman *Shaman) registerFlametongueTotemSpell() {
	dbc := shamaninfo.FlametongueTotem.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	config := shaman.newTotemSpellConfig(0.11, dbc.SpellID)
	config.SpellRanks = shamaninfo.FlametongueTotem.GetAllIDs()
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[FireTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.FlametongueTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStrengthOfEarthTotemSpell() {
	dbc := shamaninfo.StrengthofEarthTotem.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	config := shaman.newTotemSpellConfig(0.1, dbc.SpellID)
	config.SpellRanks = shamaninfo.StrengthofEarthTotem.GetAllIDs()
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.StrengthOfEarthTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerTremorTotemSpell() {
	if shaman.Level < 18 {
		return
	}
	config := shaman.newTotemSpellConfig(0.02, 8143)
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.TremorTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerStoneskinTotemSpell() {
	dbc := shamaninfo.StoneskinTotem.GetMaxRank(shaman.Level)
	if dbc == nil {
		return
	}
	config := shaman.newTotemSpellConfig(0.1, dbc.SpellID)
	config.SpellRanks = shamaninfo.StoneskinTotem.GetAllIDs()
	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
		shaman.TotemExpirations[EarthTotem] = sim.CurrentTime + time.Second*300
	}
	shaman.StoneskinTotem = shaman.RegisterSpell(config)
}

func (shaman *Shaman) registerCallOfTheElements() {
	airTotem := shaman.getAirTotemSpell(shaman.Totems.Air)
	earthTotem := shaman.getEarthTotemSpell(shaman.Totems.Earth)
	fireTotem := shaman.getFireTotemSpell(shaman.Totems.Fire)
	waterTotem := shaman.getWaterTotemSpell(shaman.Totems.Water)

	totalManaCost := 0.0
	anyTotems := false
	if airTotem != nil {
		totalManaCost += airTotem.DefaultCast.Cost
		anyTotems = true
	}
	if earthTotem != nil {
		totalManaCost += earthTotem.DefaultCast.Cost
		anyTotems = true
	}
	if fireTotem != nil {
		totalManaCost += fireTotem.DefaultCast.Cost
		anyTotems = true
	}
	if waterTotem != nil {
		totalManaCost += waterTotem.DefaultCast.Cost
		anyTotems = true
	}

	shaman.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 66842},
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return anyTotems && shaman.CurrentMana() >= totalManaCost
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Save GCD timer value, so we can safely reset it between each totem cast.
			nextGcdAt := shaman.GCD.ReadyAt()

			if airTotem != nil {
				shaman.GCD.Set(sim.CurrentTime)
				airTotem.Cast(sim, target)
			}
			if earthTotem != nil {
				shaman.GCD.Set(sim.CurrentTime)
				earthTotem.Cast(sim, target)
			}
			if fireTotem != nil {
				shaman.GCD.Set(sim.CurrentTime)
				fireTotem.Cast(sim, target)
			}
			if waterTotem != nil {
				shaman.GCD.Set(sim.CurrentTime)
				waterTotem.Cast(sim, target)
			}

			shaman.GCD.Set(nextGcdAt)
		},
	})
}

func (shaman *Shaman) getAirTotemSpell(totemType proto.AirTotem) *core.Spell {
	switch totemType {
	case proto.AirTotem_WrathOfAirTotem:
		return shaman.WrathOfAirTotem
	case proto.AirTotem_WindfuryTotem:
		return shaman.WindfuryTotem
	}
	return nil
}

func (shaman *Shaman) getEarthTotemSpell(totemType proto.EarthTotem) *core.Spell {
	switch totemType {
	case proto.EarthTotem_StrengthOfEarthTotem:
		return shaman.StrengthOfEarthTotem
	case proto.EarthTotem_TremorTotem:
		return shaman.TremorTotem
	case proto.EarthTotem_StoneskinTotem:
		return shaman.StoneskinTotem
	}
	return nil
}

func (shaman *Shaman) getFireTotemSpell(totemType proto.FireTotem) *core.Spell {
	switch totemType {
	case proto.FireTotem_TotemOfWrath:
		return shaman.TotemOfWrath
	case proto.FireTotem_SearingTotem:
		return shaman.SearingTotem
	case proto.FireTotem_MagmaTotem:
		return shaman.MagmaTotem
	case proto.FireTotem_FlametongueTotem:
		return shaman.FlametongueTotem
	}
	return nil
}

func (shaman *Shaman) getWaterTotemSpell(totemType proto.WaterTotem) *core.Spell {
	switch totemType {
	case proto.WaterTotem_ManaSpringTotem:
		return shaman.ManaSpringTotem
	case proto.WaterTotem_HealingStreamTotem:
		return shaman.HealingStreamTotem
	}
	return nil
}
