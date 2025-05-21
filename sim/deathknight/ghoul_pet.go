package deathknight

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type GhoulPet struct {
	core.Pet

	dkOwner *Deathknight

	GhoulFrenzyAura *core.Aura
	Claw            *core.Spell

	uptimePercent float64
}

func (dk *Deathknight) NewArmyGhoulPet(_ int) *GhoulPet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	var nocsHit float64
	if dk.nervesOfColdSteelActive() {
		nocsHit = float64(dk.Talents.NervesOfColdSteel) * dk.MeleeHitRatingPerHitChance
	}
	if dk.HasDraeneiHitAura {
		nocsHit += 1 * dk.MeleeHitRatingPerHitChance
	}

	minDamage := float64(dk.Level) - float64(dk.Level)/4
	maxDamage := float64(dk.Level) + float64(dk.Level)/4

	armyGhoulPetBaseStats := core.PetBaseStats[core.Pet_Unknown][1].Add(stats.Stats{stats.MeleeHit: -nocsHit, stats.Expertise: -nocsHit * dk.GetPetExpertiseScale()})

	armyGhoulPetBasePercentageStats := stats.Stats{stats.MeleeCrit: 5}

	ghoulPet := &GhoulPet{
		Pet:     core.NewPet("Army of the Dead", &dk.Character, armyGhoulPetBaseStats, armyGhoulPetBasePercentageStats, dk.armyGhoulStatInheritance(), false, true),
		dkOwner: dk,
	}

	ghoulPet.PseudoStats.DamageTakenMultiplier *= 0.1
	ghoulPet.PseudoStats.MeleeHasteRatingPerHastePercent = dk.PseudoStats.MeleeHasteRatingPerHastePercent

	dk.SetupGhoul(ghoulPet)

	ghoulPet.EnableAutoAttacks(ghoulPet, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     minDamage,
			BaseDamageMax:     maxDamage,
			SwingSpeed:        2,
			CritMultiplier:    2,
			AttackPowerPerDPS: 17.5,
		},
		AutoSwingMelee: true,
	})

	ghoulPet.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	ghoulPet.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	ghoulPet.AddStatDependency(stats.Agility, stats.MeleeCrit, ghoulPet.CritRatingPerCritChance*core.CritPerAgi[proto.Class_ClassRogue][ghoulPet.Level])

	// command doesn't apply to army ghoul
	if dk.Race == proto.Race_RaceOrc {
		ghoulPet.PseudoStats.DamageDealtMultiplier /= 1.05
	}

	return ghoulPet
}

func (dk *Deathknight) NewGhoulPet(permanent bool) *GhoulPet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	var nocsHit float64
	if dk.nervesOfColdSteelActive() {
		nocsHit = float64(dk.Talents.NervesOfColdSteel) * dk.MeleeHitRatingPerHitChance
	}
	if dk.HasDraeneiHitAura {
		nocsHit += 1 * dk.MeleeHitRatingPerHitChance
	}

	ghoulPetBaseStats := core.PetBaseStats[core.Pet_RisenGhoul][dk.Level].Stats.Add(stats.Stats{stats.MeleeHit: -nocsHit, stats.Expertise: -nocsHit * dk.GetPetExpertiseScale()})
	ghoulPetBasePercentageStats := stats.Stats{stats.MeleeCrit: 5}

	minDamage := float64(dk.Level) - float64(dk.Level)/4
	maxDamage := float64(dk.Level) + float64(dk.Level)/4

	ghoulPet := &GhoulPet{
		Pet:     core.NewPet("Ghoul", &dk.Character, ghoulPetBaseStats, ghoulPetBasePercentageStats, dk.ghoulStatInheritance(), permanent, !permanent),
		dkOwner: dk,
	}

	// NightOfTheDead
	ghoulPet.PseudoStats.DamageTakenMultiplier *= 1.0 - float64(dk.Talents.NightOfTheDead)*0.45
	ghoulPet.PseudoStats.MeleeHasteRatingPerHastePercent = dk.PseudoStats.MeleeHasteRatingPerHastePercent

	dk.SetupGhoul(ghoulPet)

	ghoulPet.EnableAutoAttacks(ghoulPet, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:     minDamage,
			BaseDamageMax:     maxDamage,
			SwingSpeed:        2,
			CritMultiplier:    2,
			AttackPowerPerDPS: 17.5,
		},
		AutoSwingMelee: true,
	})

	ghoulPet.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	ghoulPet.AddStatDependency(stats.Agility, stats.AttackPower, 1)
	ghoulPet.AddStatDependency(stats.Agility, stats.MeleeCrit, ghoulPet.CritRatingPerCritChance*core.CritPerAgi[proto.Class_ClassRogue][ghoulPet.Level])

	if permanent {
		core.ApplyPetConsumeEffects(&ghoulPet.Character, dk.Consumes)
	}

	return ghoulPet
}

func (dk *Deathknight) SetupGhoul(ghoulPet *GhoulPet) {
	ghoulPet.Pet.OnPetEnable = ghoulPet.enable

	ghoulPet.Unit.EnableFocusBar(2, func(sim *core.Simulation) {
		if ghoulPet.GCD.IsReady(sim) {
			ghoulPet.OnGCDReady(sim)
		}
	})

	dk.AddPet(ghoulPet)
}

func (ghoulPet *GhoulPet) GetPet() *core.Pet {
	return &ghoulPet.Pet
}

func (ghoulPet *GhoulPet) Initialize() {
	ghoulPet.Claw = ghoulPet.registerClaw()
}

func (ghoulPet *GhoulPet) Reset(_ *core.Simulation) {
	if !ghoulPet.IsGuardian() {
		ghoulPet.uptimePercent = min(1, max(0, ghoulPet.dkOwner.Inputs.PetUptime))
	} else {
		ghoulPet.uptimePercent = 1.0
	}
}

func (ghoulPet *GhoulPet) ExecuteCustomRotation(sim *core.Simulation) {
	if ghoulPet.uptimePercent < 1.0 { // Apply uptime for permanent pet ghoul
		if sim.GetRemainingDurationPercent() < 1.0-ghoulPet.uptimePercent { // once fight is % completed, disable pet.
			ghoulPet.Pet.Disable(sim)
			return
		}
	}

	if ghoulPet.CurrentFocus() < ghoulPet.Claw.DefaultCast.Cost {
		return
	}

	ghoulPet.Claw.Cast(sim, ghoulPet.CurrentTarget)
}

func (ghoulPet *GhoulPet) enable(sim *core.Simulation) {
	if ghoulPet.IsGuardian() {
		ghoulPet.PseudoStats.MeleeSpeedMultiplier = 1 // guardians are not affected by raid buffs
		ghoulPet.MultiplyMeleeSpeed(sim, ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier)
		return
	}

	ghoulPet.MultiplyMeleeSpeed(sim, ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier)

	ghoulPet.EnableDynamicMeleeSpeed(func(amount float64) {
		ghoulPet.MultiplyMeleeSpeed(sim, amount)

		if sim.Log != nil {
			sim.Log("Ghoul MeleeSpeedMultiplier: %f, ownerMeleeMultiplier: %f\n", ghoulPet.Character.PseudoStats.MeleeSpeedMultiplier, ghoulPet.dkOwner.PseudoStats.MeleeSpeedMultiplier)
		}
	})
}

func (dk *Deathknight) ghoulStatInheritance() core.PetStatInheritance {
	ravenousDead := 1.0 + 0.2*float64(dk.Talents.RavenousDead)
	glyphBonus := 0.0
	if dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfTheGhoul) {
		glyphBonus = 0.4
	}

	baseStatsScale := glyphBonus + 0.7*ravenousDead

	return func(ownerStats stats.Stats, _ stats.PseudoStats) stats.Stats {
		return stats.Stats{
			stats.Stamina:  ownerStats[stats.Stamina] * baseStatsScale,
			stats.Strength: ownerStats[stats.Strength] * baseStatsScale,

			stats.MeleeHit:  ownerStats[stats.MeleeHit],
			stats.Expertise: ownerStats[stats.MeleeHit] * dk.GetPetExpertiseScale(),

			stats.MeleeHaste: ownerStats[stats.MeleeHaste],
		}
	}
}

func (dk *Deathknight) armyGhoulStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats, _ stats.PseudoStats) stats.Stats {
		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.2,
			stats.AttackPower: ownerStats[stats.AttackPower] * 0.065,

			stats.MeleeHit:  ownerStats[stats.MeleeHit],
			stats.Expertise: ownerStats[stats.MeleeHit] * dk.GetPetExpertiseScale(),

			stats.MeleeHaste: ownerStats[stats.MeleeHaste],
		}
	}
}

func (ghoulPet *GhoulPet) registerClaw() *core.Spell {
	return ghoulPet.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47468},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		FocusCost: core.FocusCostOptions{
			Cost:   40,
			Refund: 0.8,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1.5,
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
