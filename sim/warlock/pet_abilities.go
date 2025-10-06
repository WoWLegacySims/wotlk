package warlock

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/warlockinfo"
)

func (wp *WarlockPet) registerCleaveSpell() {
	dbc := warlockinfo.Cleave.GetMaxRank(wp.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, wp.Level)

	numHits := min(2, wp.Env.GetNumTargets())

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * 6,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			constBaseDamage := bp + spell.BonusWeaponDamage()

			curTarget := target
			for hitIndex := int32(0); hitIndex < numHits; hitIndex++ {
				baseDamage := constBaseDamage + spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower())
				spell.CalcAndDealDamage(sim, curTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				curTarget = sim.Environment.NextTargetUnit(curTarget)
			}
		},
	})
}

func (wp *WarlockPet) registerInterceptSpell() {
	wp.secondaryAbility = nil // not implemented
}

func (wp *WarlockPet) registerLashOfPainSpell() {
	dbc := warlockinfo.LashofPain.GetMaxRank(wp.Level)
	if dbc == nil {
		return
	}
	bp, _ := dbc.GetBPDie(0, wp.Level)

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: dbc.ManaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (12 - time.Duration(3*wp.owner.Talents.DemonicPower)),
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// TODO: the hidden 5% damage modifier succ currently gets also applies to this ...
			baseDamage := bp + 0.429*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}

func (wp *WarlockPet) registerShadowBiteSpell() {
	dbc := warlockinfo.ShadowBite.GetMaxRank(wp.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, wp.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}

	var petManaMetrics *core.ResourceMetrics
	maxManaMult := 0.04 * float64(wp.owner.Talents.ImprovedFelhunter)
	impFelhunter := wp.owner.Talents.ImprovedFelhunter > 0
	if impFelhunter {
		petManaMetrics = wp.NewManaMetrics(actionID)
	}

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 0.03,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    wp.NewTimer(),
				Duration: time.Second * (6 - time.Duration(2*wp.owner.Talents.ImprovedFelhunter)),
			},
		},

		DamageMultiplier: 1 + 0.03*float64(wp.owner.Talents.ShadowMastery),
		CritMultiplier:   1.5 + 0.1*float64(wp.owner.Talents.Ruin),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + 0.429*spell.SpellPower()

			w := wp.owner
			spells := []*core.Spell{
				w.UnstableAffliction,
				w.Immolate,
				w.CurseOfAgony,
				w.CurseOfDoom,
				w.Corruption,
				w.Conflagrate,
				w.Seed,
				w.DrainSoul,
				// missing: drain life, shadowflame
			}
			counter := 0
			for _, spell := range spells {
				if spell != nil && spell.Dot(target).IsActive() {
					counter++
				}
			}

			baseDamage *= 1 + 0.15*float64(counter)

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if impFelhunter && result.Landed() {
				wp.AddMana(sim, wp.MaxMana()*maxManaMult, petManaMetrics)
			}
			spell.DealDamage(sim, result)
		},
	})
}

func (wp *WarlockPet) registerFireboltSpell() {
	dbc := warlockinfo.Firebolt.GetMaxRank(wp.Level)
	if dbc == nil {
		return
	}
	bp, die := dbc.GetBPDie(0, wp.Level)

	wp.primaryAbility = wp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: dbc.ManaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: dbc.CastTime - (250 * time.Millisecond * time.Duration(wp.owner.Talents.DemonicPower)),
			},
		},

		DamageMultiplier: (1 + 0.1*float64(wp.owner.Talents.ImprovedImp)) *
			(1 + 0.2*core.TernaryFloat64(wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfImp), 1, 0)),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + 0.714*spell.SpellPower()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
		},
	})
}
