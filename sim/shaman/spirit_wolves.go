package shaman

import (
	"math"
	"strconv"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type SpiritWolf struct {
	core.Pet

	shamanOwner *Shaman
}

type SpiritWolves struct {
	SpiritWolf1 *SpiritWolf
	SpiritWolf2 *SpiritWolf
}

func (SpiritWolves *SpiritWolves) EnableWithTimeout(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.EnableWithTimeout(sim, SpiritWolves.SpiritWolf1, time.Second*45)
	SpiritWolves.SpiritWolf2.EnableWithTimeout(sim, SpiritWolves.SpiritWolf2, time.Second*45)
}

func (SpiritWolves *SpiritWolves) CancelGCDTimer(sim *core.Simulation) {
	SpiritWolves.SpiritWolf1.CancelGCDTimer(sim)
	SpiritWolves.SpiritWolf2.CancelGCDTimer(sim)
}

var spiritWolfBaseStats = core.PetBaseStats[core.Pet_Unknown][1].Stats

var spiritWolfBasePercentageStats = stats.Stats{
	// Add 1.8% because pets aren't affected by that component of crit suppression.
	stats.MeleeCrit: (1.1515 + 1.8),
}

func (shaman *Shaman) NewSpiritWolf(index int) *SpiritWolf {
	spiritWolf := &SpiritWolf{
		Pet:         core.NewPet("Spirit Wolf "+strconv.Itoa(index), &shaman.Character, spiritWolfBaseStats, spiritWolfBasePercentageStats, shaman.makeStatInheritance(), false, false),
		shamanOwner: shaman,
	}

	spiritWolf.EnableAutoAttacks(spiritWolf, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  float64(spiritWolf.Level) * 3,
			BaseDamageMax:  float64(spiritWolf.Level) * 5,
			SwingSpeed:     1.5,
			CritMultiplier: 2,
		},
		AutoSwingMelee: true,
	})

	spiritWolf.AddStatDependency(stats.Strength, stats.AttackPower, 2)
	spiritWolf.AddStatDependency(stats.Agility, stats.MeleeCrit, spiritWolf.CritRatingPerCritChance*core.CritPerAgi[proto.Class_ClassRogue][spiritWolf.Level])
	core.ApplyPetConsumeEffects(&spiritWolf.Character, shaman.Consumes)

	shaman.AddPet(spiritWolf)

	return spiritWolf
}

func (shaman *Shaman) makeStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats, _ stats.PseudoStats) stats.Stats {
		ownerHitChance := ownerStats[stats.MeleeHit] / shaman.MeleeHitRatingPerHitChance
		hitRatingFromOwner := math.Floor(ownerHitChance) * shaman.MeleeHitRatingPerHitChance

		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Armor:       ownerStats[stats.Armor] * 0.35,
			stats.AttackPower: ownerStats[stats.AttackPower] * (core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfFeralSpirit), 0.6, 0.3)),
			stats.SpellPower:  ownerStats[stats.SpellPower] * 0.3,

			stats.MeleeHit:  hitRatingFromOwner,
			stats.Expertise: shaman.CalculateHitInheritance(stats.SpellHit, stats.Expertise),
		}
	}
}

func (spiritWolf *SpiritWolf) Initialize() {
	// Nothing
}

func (spiritWolf *SpiritWolf) ExecuteCustomRotation(_ *core.Simulation) {
}

func (spiritWolf *SpiritWolf) Reset(sim *core.Simulation) {
	spiritWolf.Disable(sim)
	if sim.Log != nil {
		spiritWolf.Log(sim, "Base Stats: %s", spiritWolfBaseStats)
		inheritedStats := spiritWolf.shamanOwner.makeStatInheritance()(spiritWolf.shamanOwner.GetStats(), stats.PseudoStats{})
		spiritWolf.Log(sim, "Inherited Stats: %s", inheritedStats)
		spiritWolf.Log(sim, "Total Stats: %s", spiritWolf.GetStats())
	}
}

func (spiritWolf *SpiritWolf) GetPet() *core.Pet {
	return &spiritWolf.Pet
}
