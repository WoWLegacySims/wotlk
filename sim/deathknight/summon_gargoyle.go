package deathknight

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerSummonGargoyleCD() {
	if !dk.Talents.SummonGargoyle {
		return
	}

	dk.SummonGargoyleAura = dk.RegisterAura(core.Aura{
		Label:    "Summon Gargoyle",
		ActionID: core.ActionID{SpellID: 49206},
		Duration: time.Second * 30,
	})

	dk.SummonGargoyle = dk.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 49206},
		Flags:    core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 60,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			dk.Gargoyle.EnableWithTimeout(sim, dk.Gargoyle, time.Second*30)
			dk.Gargoyle.CancelGCDTimer(sim)

			// Add a dummy aura to show in metrics
			dk.SummonGargoyleAura.Activate(sim)

			// Start casting after a 2.5s delay to simulate the summon animation
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + dk.GargoyleSummonDelay,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
					dk.OnGargoyleStartFirstCast()
					dk.Gargoyle.GargoyleStrike.Cast(sim, dk.CurrentTarget)
				},
			}
			sim.AddPendingAction(&pa)
		},
	})

	dk.AddMajorCooldown(core.MajorCooldown{
		Spell: dk.SummonGargoyle,
		Type:  core.CooldownTypeDPS,
	})
	if dk.Inputs.IsDps {
		// We use this for defining the min cast time of gargoyle,
		// but we don't cast it with the MCD system in the dps sim
		dk.GetMajorCooldown(dk.SummonGargoyle.ActionID).Disable()
	}
}

type GargoylePet struct {
	core.Pet

	dkOwner *Deathknight

	GargoyleStrike *core.Spell
}

func (dk *Deathknight) NewGargoyle() *GargoylePet {
	// Remove any hit that would be given by NocS as it does not translate to pets
	var nocsHit float64
	if dk.nervesOfColdSteelActive() {
		nocsHit = float64(dk.Talents.NervesOfColdSteel) * dk.MeleeHitRatingPerHitChance
	}
	if dk.HasDraeneiHitAura {
		nocsHit += 1 * dk.MeleeHitRatingPerHitChance
	}

	gargoyleStats := core.PetBaseStats[core.Pet_Unknown][1].Stats.Add(stats.Stats{stats.Mana: 28 + 10*float64(dk.Level), stats.Health: 28 + 30*float64(dk.Level), stats.SpellHit: -nocsHit * dk.GetPetSpellHitScale()})

	gargoyle := &GargoylePet{
		Pet: core.NewPet("Gargoyle", &dk.Character, gargoyleStats, stats.Stats{stats.SpellCrit: 5}, func(ownerStats stats.Stats) stats.Stats {
			return stats.Stats{
				stats.AttackPower: ownerStats[stats.AttackPower],
				stats.SpellHit:    ownerStats[stats.MeleeHit] * dk.GetPetSpellHitScale(),
				stats.SpellHaste:  ownerStats[stats.MeleeHaste] * PetSpellHasteScale,
			}
		}, false, true),
		dkOwner: dk,
	}

	// NightOfTheDead
	gargoyle.PseudoStats.DamageTakenMultiplier *= 1.0 - float64(dk.Talents.NightOfTheDead)*0.45

	gargoyle.OnPetEnable = func(sim *core.Simulation) {
		gargoyle.PseudoStats.CastSpeedMultiplier = 1 // guardians are not affected by raid buffs
		gargoyle.MultiplyCastSpeed(dk.PseudoStats.MeleeSpeedMultiplier)
	}

	dk.AddPet(gargoyle)

	return gargoyle
}

func (garg *GargoylePet) GetPet() *core.Pet {
	return &garg.Pet
}

func (garg *GargoylePet) Initialize() {
	garg.registerGargoyleStrikeSpell()
}

func (garg *GargoylePet) Reset(_ *core.Simulation) {
}

func (garg *GargoylePet) ExecuteCustomRotation(sim *core.Simulation) {
	garg.GargoyleStrike.Cast(sim, garg.CurrentTarget)
}

func (garg *GargoylePet) registerGargoyleStrikeSpell() {
	attackPowerModifier := (1.0 + 0.04*float64(garg.dkOwner.Talents.Impurity)) / 3.0

	garg.GargoyleStrike = garg.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 51963},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				CastTime: time.Millisecond * 2000,
				GCD:      core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1.5,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 2.05*sim.Roll(51, 69) + attackPowerModifier*spell.MeleeAttackPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.DealDamage(sim, result)
		},
	})
}
