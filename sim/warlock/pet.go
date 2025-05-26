package warlock

import (
	"math"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type WarlockPet struct {
	core.Pet

	owner *Warlock

	primaryAbility   *core.Spell
	secondaryAbility *core.Spell

	DemonicEmpowermentAura *core.Aura
}

func (warlock *Warlock) NewWarlockPet() *WarlockPet {
	var cfg struct {
		Name            string
		PowerModifier   float64
		Stats           stats.Stats
		PercentageStats stats.Stats
		AutoAttacks     core.AutoAttackOptions
	}
	baseStats := core.PetBaseStats[int32(warlock.Options.Summon)][warlock.Level]
	cfg.Stats = baseStats.Stats
	cfg.Name = warlock.Options.Summon.String()
	if warlock.Options.Summon == proto.Warlock_Options_Imp {
		cfg.PowerModifier = 4.95
		cfg.PercentageStats = stats.Stats{
			stats.MeleeCrit: 3.454,
			stats.SpellCrit: 0.9075,
		}
	} else {
		cfg.AutoAttacks = core.AutoAttackOptions{
			MainHand: core.Weapon{
				BaseDamageMin:  float64(baseStats.Min_dmg),
				BaseDamageMax:  float64(baseStats.Max_dmg),
				SwingSpeed:     2,
				CritMultiplier: 2,
			},
			AutoSwingMelee: true,
		}
		cfg.PowerModifier = 11.5
		cfg.PercentageStats = stats.Stats{
			stats.MeleeCrit: 3.2685,
			stats.SpellCrit: 3.3355,
		}
	}

	wp := &WarlockPet{
		Pet:   core.NewPet(cfg.Name, &warlock.Character, cfg.Stats, cfg.PercentageStats, warlock.makeStatInheritance(), true, false),
		owner: warlock,
	}

	wp.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	wp.AddStat(stats.AttackPower, -20)

	if warlock.HasSetBonus(ItemSetOblivionRaiment, 2) {
		wp.AddStat(stats.MP5, 45)
	}

	if warlock.Options.Summon == proto.Warlock_Options_Imp {
		// imps are mages
		wp.Class = proto.Class_ClassMage
	} else {
		wp.Class = proto.Class_ClassPaladin
	}
	wp.AddStatDependency(stats.Agility, stats.MeleeCrit, wp.CritRatingPerCritChance*core.CritPerAgi[wp.Class][wp.Level])
	wp.EnableManaBarWithModifier(cfg.PowerModifier)

	wp.AddStats(stats.Stats{
		stats.MeleeCrit: float64(warlock.Talents.DemonicTactics) * 2 * wp.CritRatingPerCritChance,
		stats.SpellCrit: float64(warlock.Talents.DemonicTactics) * 2 * wp.CritRatingPerCritChance,
	})

	wp.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.04*float64(warlock.Talents.UnholyPower)

	if warlock.Options.Summon != proto.Warlock_Options_Imp { // imps generally don't meele
		wp.EnableAutoAttacks(wp, cfg.AutoAttacks)
	}

	if warlock.Options.Summon == proto.Warlock_Options_Felguard {
		if wp.owner.HasMajorGlyph(proto.WarlockMajorGlyph_GlyphOfFelguard) {
			wp.MultiplyStat(stats.AttackPower, 1.2)
		}

		statDeps := []*stats.StatDependency{nil}
		for i := 1; i <= 10; i++ {
			statDeps = append(statDeps, wp.NewDynamicMultiplyStat(stats.AttackPower,
				1+float64(i)*(0.05+0.01*float64(warlock.Talents.DemonicBrutality))))
		}

		DemonicFrenzyAura := wp.RegisterAura(core.Aura{
			Label:     "Demonic Frenzy",
			ActionID:  core.ActionID{SpellID: 32851},
			Duration:  time.Second * 10,
			MaxStacks: 10,
			OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
				if oldStacks != 0 {
					aura.Unit.DisableDynamicStatDep(sim, statDeps[oldStacks])
				}
				if newStacks != 0 {
					aura.Unit.EnableDynamicStatDep(sim, statDeps[newStacks])
				}
			},
		})
		wp.RegisterAura(core.Aura{
			Label:    "Demonic Frenzy Hidden Aura",
			Duration: core.NeverExpires,
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				aura.Activate(sim)
			},
			OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !result.Landed() || !spell.ProcMask.Matches(core.ProcMaskMelee) {
					return
				}
				DemonicFrenzyAura.Activate(sim)
				DemonicFrenzyAura.AddStack(sim)
			},
		})
	}

	if warlock.Talents.FelVitality > 0 {
		bonus := 1.0 + 0.05*float64(warlock.Talents.FelVitality)
		wp.MultiplyStat(stats.Intellect, bonus)
		wp.MultiplyStat(stats.Stamina, bonus)
	}

	if warlock.Talents.MasterDemonologist > 0 {
		val := 1.0 + 0.01*float64(warlock.Talents.MasterDemonologist)
		md := core.Aura{
			Label:    "Master Demonologist",
			ActionID: core.ActionID{SpellID: 35706}, // many different spells associated with this talent
			Duration: core.NeverExpires,
			OnGain: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.Warlock_Options_Imp:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= val
				case proto.Warlock_Options_Succubus:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= val
				case proto.Warlock_Options_Felguard:
					aura.Unit.PseudoStats.DamageDealtMultiplier *= val
				}
			},
			OnExpire: func(aura *core.Aura, _ *core.Simulation) {
				switch warlock.Options.Summon {
				case proto.Warlock_Options_Imp:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= val
				case proto.Warlock_Options_Succubus:
					aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= val
				case proto.Warlock_Options_Felguard:
					aura.Unit.PseudoStats.DamageDealtMultiplier /= val
				}
			},
		}

		mdLockAura := warlock.RegisterAura(md)
		mdPetAura := wp.RegisterAura(md)

		masterDemonologist := float64(warlock.Talents.MasterDemonologist)
		masterDemonologistFireCrit := core.TernaryFloat64(warlock.Options.Summon == proto.Warlock_Options_Imp, masterDemonologist, 0)
		masterDemonologistShadowCrit := core.TernaryFloat64(warlock.Options.Summon == proto.Warlock_Options_Succubus, masterDemonologist, 0)

		wp.OnPetEnable = func(sim *core.Simulation) {
			mdLockAura.Activate(sim)
			mdPetAura.Activate(sim)

			spellbook := make([]*core.Spell, 0)
			spellbook = append(spellbook, warlock.Spellbook...)
			spellbook = append(spellbook, wp.Spellbook...)

			for _, spell := range spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					spell.BonusCrit += masterDemonologistFireCrit
				}

				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.BonusCrit += masterDemonologistShadowCrit
				}
			}
		}

		wp.OnPetDisable = func(sim *core.Simulation) {
			mdLockAura.Deactivate(sim)
			mdPetAura.Deactivate(sim)

			spellbook := make([]*core.Spell, 0)
			spellbook = append(spellbook, warlock.Spellbook...)
			spellbook = append(spellbook, wp.Spellbook...)

			for _, spell := range spellbook {
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					spell.BonusCrit -= masterDemonologistFireCrit
				}

				if spell.SpellSchool.Matches(core.SpellSchoolShadow) {
					spell.BonusCrit -= masterDemonologistShadowCrit
				}
			}
		}
	}

	core.ApplyPetConsumeEffects(&wp.Character, warlock.Consumes)

	warlock.AddPet(wp)

	return wp
}

func (wp *WarlockPet) GetPet() *core.Pet {
	return &wp.Pet
}

func (wp *WarlockPet) Initialize() {
	switch wp.owner.Options.Summon {
	case proto.Warlock_Options_Felguard:
		wp.registerCleaveSpell()
		wp.registerInterceptSpell()
	case proto.Warlock_Options_Succubus:
		wp.registerLashOfPainSpell()
	case proto.Warlock_Options_Felhunter:
		wp.registerShadowBiteSpell()
	case proto.Warlock_Options_Imp:
		wp.registerFireboltSpell()
	}
}

func (wp *WarlockPet) Reset(_ *core.Simulation) {
}

func (wp *WarlockPet) ExecuteCustomRotation(sim *core.Simulation) {
	if !wp.primaryAbility.IsReady(sim) {
		wp.WaitUntil(sim, wp.primaryAbility.CD.ReadyAt())
		return
	}

	wp.primaryAbility.Cast(sim, wp.CurrentTarget)
}

func (warlock *Warlock) makeStatInheritance() core.PetStatInheritance {
	improvedDemonicTactics := float64(warlock.Talents.ImprovedDemonicTactics)

	return func(ownerStats stats.Stats, ownerPseudoStats stats.PseudoStats) stats.Stats {
		ownerSP := ownerStats[stats.SpellPower] + math.Max(ownerPseudoStats.FireSpellPower, ownerPseudoStats.ShadowSpellPower)

		return stats.Stats{
			stats.Stamina:          ownerStats[stats.Stamina] * 0.75,
			stats.Intellect:        ownerStats[stats.Intellect] * 0.3,
			stats.Armor:            ownerStats[stats.Armor] * 0.35,
			stats.AttackPower:      ownerSP * 0.57,
			stats.SpellPower:       ownerSP * 0.15,
			stats.SpellPenetration: ownerStats[stats.SpellPenetration],
			stats.SpellCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeCrit:        improvedDemonicTactics * 0.1 * ownerStats[stats.SpellCrit],
			stats.MeleeHit:         warlock.CalculateHitInheritance(stats.SpellHit, stats.MeleeHit),
			stats.SpellHit:         ownerStats[stats.SpellHit],
			stats.Expertise:        warlock.CalculateHitInheritance(stats.SpellHit, stats.Expertise),
			// Resists, 40%
		}
	}
}
