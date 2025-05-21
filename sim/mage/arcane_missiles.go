package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

func (mage *Mage) registerArcaneMissilesSpell() {
	dbc := mageinfo.ArcaneMissiles.GetMaxRank(mage.Level)
	dbcDmg := mageinfo.ArcaneMissilesDamage.GetMaxRank(mage.Level)
	if dbc == nil || dbcDmg == nil {
		return
	}

	ticks := dbc.Duration / dbc.Effects[1].AuraPeriod
	bp, _ := dbcDmg.GetBPDie(0, mage.Level)

	spellCoeff := (dbcDmg.GetCoefficient(0) + 0.03*float64(mage.Talents.ArcaneEmpowerment)) * dbcDmg.GetLevelPenalty(mage.Level)
	hasT8_4pc := mage.HasSetBonus(ItemSetKirinTorGarb, 4)

	mage.ArcaneMissilesTickSpell = mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbcDmg.SpellID},
		SpellSchool: core.SpellSchoolArcane,
		// unlike Mind Flay, this CAN proc JoW. It can also proc trinkets without the "can proc from proc" flag
		// such as illustration of the dragon soul
		// however, it cannot proc Nibelung so we add the ProcMaskNotInSpellbook flag
		ProcMask:         core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:            SpellFlagMage | core.SpellFlagNoLogs,
		MissileSpeed:     20,
		BonusHit:         float64(mage.Talents.ArcaneFocus),
		BonusCrit:        core.TernaryFloat64(mage.HasSetBonus(ItemSetKhadgarsRegalia, 4), 5, 0),
		DamageMultiplier: 1 + .04*float64(mage.Talents.TormentTheWeak),
		DamageMultiplierAdditive: 1 +
			core.TernaryFloat64(mage.HasSetBonus(ItemSetTempestRegalia, 4), .05, 0),
		CritMultiplier:   mage.SpellCritMultiplier(1, mage.bonusCritDamage+core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneMissiles), .25, 0)),
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := bp + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, damage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolArcane,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagChanneled | core.SpellFlagAPL,
		BonusHit:    float64(mage.Talents.ArcaneFocus),
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.31,
			Multiplier: 1 - .01*float64(mage.Talents.ArcaneFocus),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissiles",
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.MissileBarrageAura.IsActive() {
						if !hasT8_4pc || sim.RandomFloat("MageT84PC") > T84PcProcChance {
							mage.MissileBarrageAura.Deactivate(sim)
						}
					}

					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.
					dot := mage.ArcaneMissiles.Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
					mage.ArcaneBlastAura.Deactivate(sim)
				},
			},
			NumberOfTicks:       ticks,
			TickLength:          time.Second,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.ArcaneMissilesTickSpell.Cast(sim, target)
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})
}
