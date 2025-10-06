package warlock

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (warlock *Warlock) registerLifeTapSpell() {
	dbc := warlockinfo.LifeTap.GetMaxRank(warlock.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, warlock.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}
	impLifetap := 1.0 + 0.1*float64(warlock.Talents.ImprovedLifeTap)
	manaMetrics := warlock.NewManaMetrics(actionID)

	var petManaMetrics *core.ResourceMetrics
	if warlock.Talents.ManaFeed && warlock.Pet != nil {
		petManaMetrics = warlock.Pet.NewManaMetrics(actionID)
	}

	warlock.LifeTap = warlock.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellRanks:  warlockinfo.LifeTap.GetAllIDs(),
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			restore := (bp + 0.5*spell.SpellPower()) * impLifetap
			warlock.AddMana(sim, restore, manaMetrics)

			if warlock.Talents.ManaFeed && warlock.Pet != nil {
				warlock.Pet.AddMana(sim, restore, petManaMetrics)
			}
			if warlock.GlyphOfLifeTapAura != nil {
				warlock.GlyphOfLifeTapAura.Activate(sim)
			}
		},
	})
}
