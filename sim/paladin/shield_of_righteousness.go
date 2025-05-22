package paladin

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/paladininfo"
)

func (paladin *Paladin) registerShieldOfRighteousnessSpell() {
	dbc := paladininfo.ShieldofRighteousness.GetMaxRank(paladin.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, paladin.Level)

	var aegisPlateProcAura *core.Aura
	if paladin.HasSetBonus(ItemSetAegisPlate, 4) {
		aegisPlateProcAura = paladin.NewTemporaryStatsAura("Aegis", core.ActionID{SpellID: 64883}, stats.Stats{stats.BlockValue: 225}, time.Second*6)
	}

	paladin.ShieldOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellRanks:  paladininfo.ShieldofRighteousness.GetAllIDs(),
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.06,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   paladin.MeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if aegisPlateProcAura != nil {
				aegisPlateProcAura.Activate(sim)
			}

			bv := paladin.GetShieldBlockValue(float64(paladin.Level)*29.5, float64(paladin.Level)*34.5)
			baseDamage := bp + bv

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}
