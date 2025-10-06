package druid

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func (druid *Druid) registerForceOfNatureCD() {
	if !druid.Talents.ForceOfNature {
		return
	}

	forceOfNatureAura := druid.RegisterAura(core.Aura{
		Label:    "Force of Nature",
		ActionID: core.ActionID{SpellID: 65861},
		Duration: time.Second * 30,
	})
	druid.ForceOfNature = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID: core.ActionID{SpellID: 65861},
		Flags:    core.SpellFlagAPL,
		ManaCost: core.ManaCostOptions{
			BaseCost:   0.12,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			druid.Treant1.EnableWithTimeout(sim, druid.Treant1, time.Second*30)
			druid.Treant2.EnableWithTimeout(sim, druid.Treant2, time.Second*30)
			druid.Treant3.EnableWithTimeout(sim, druid.Treant3, time.Second*30)

			forceOfNatureAura.Activate(sim)

			// Animation delay, courtesy of our DK friends
			pa := core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Second*1,
				Priority:     core.ActionPriorityAuto,
				OnAction: func(s *core.Simulation) {
				},
			}
			sim.AddPendingAction(&pa)
		},
	})
}

type TreantPet struct {
	core.Pet
	druidOwner *Druid

	snapshotStat stats.Stats
}

func (druid *Druid) NewTreant() *TreantPet {
	treant := &TreantPet{
		Pet:        core.NewPet("Treant", &druid.Character, treantBaseStats, treantBasePercentageStats, druid.treantStatInheritance(), false, false),
		druidOwner: druid,
	}
	treant.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	treant.AddStatDependency(stats.Agility, stats.MeleeCrit, treant.CritRatingPerCritChance*core.CritPerAgi[proto.Class_ClassRogue][treant.Level])

	treant.PseudoStats.DamageDealtMultiplier = 1 + 0.05*float64(druid.Talents.Brambles)
	treant.EnableAutoAttacks(treant, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  float64(druid.Level) - float64(druid.Level)/4,
			BaseDamageMax:  float64(druid.Level) - float64(druid.Level)/4,
			SwingSpeed:     2,
			CritMultiplier: druid.BalanceCritMultiplier(),
		},
		AutoSwingMelee: true,
	})

	treant.Pet.OnPetEnable = treant.enable
	treant.Pet.OnPetDisable = treant.disable

	druid.AddPet(treant)
	return treant
}

func (treant *TreantPet) GetPet() *core.Pet {
	return &treant.Pet
}

func (treant *TreantPet) enable(sim *core.Simulation) {
	// Snapshot spellpower
	treant.snapshotStat = stats.Stats{stats.Strength: treant.druidOwner.GetStat(stats.SpellPower) * 0.5}
	treant.AddStatsDynamic(sim, treant.snapshotStat)
}

func (treant *TreantPet) disable(sim *core.Simulation) {
	treant.AddStatsDynamic(sim, treant.snapshotStat.Invert())
}

func (treant *TreantPet) Initialize() {
}

func (treant *TreantPet) Reset(_ *core.Simulation) {
}

func (treant *TreantPet) ExecuteCustomRotation(_ *core.Simulation) {
}

func (druid *Druid) treantStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats, ownerPseudoStats stats.PseudoStats) stats.Stats {

		return stats.Stats{ //still need to nail down shadow fiend crit scaling, but removing owner crit scaling after further investigation
			// 3 x sp to ap, lol
			stats.AttackPower: (ownerStats[stats.SpellPower] + ownerPseudoStats.NatureSpellPower) * 1.05,
			stats.MeleeHit:    druid.CalculateHitInheritance(stats.SpellHit, stats.MeleeHit),
			stats.Expertise:   druid.CalculateHitInheritance(stats.SpellHit, stats.Expertise),
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.3,
		}
	}
}

var treantBaseStats = core.PetBaseStats[core.Pet_Unknown][1].Stats

var treantBasePercentageStats = stats.Stats{
	stats.MeleeCrit: 5,
}
