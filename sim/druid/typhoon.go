package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/druidinfo"
)

func (druid *Druid) registerTyphoonSpell() {
	if !druid.Talents.Typhoon {
		return
	}
	dbc := druidinfo.Typhoon.GetMaxRank(druid.Level)
	if dbc == nil {
		return
	}
	dmg, _ := dbc.GetBPDie(1, druid.Level)
	druid.Typhoon = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagOmenTrigger | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.25,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Second * (20 - core.TernaryDuration(druid.HasMajorGlyph(proto.DruidMajorGlyph_GlyphOfMonsoon), 3, 0)),
			},
		},

		DamageMultiplier: 1 +
			0.15*float64(druid.Talents.GaleWinds),
		ThreatMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				baseDamage := dmg + 0.193*spell.SpellPower()
				baseDamage *= sim.Encounter.AOECapMultiplier()
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMagicHitAndCrit)
				}
			})
		},
	})
}
