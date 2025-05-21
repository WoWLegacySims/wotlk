package shaman

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

// Variables that control the Fire Elemental.
const (
	// 7.5 CPM
	maxFireBlastCasts = 15
	maxFireNovaCasts  = 15
)

type FireElemental struct {
	core.Pet

	FireBlast *core.Spell
	FireNova  *core.Spell

	FireShieldAura *core.Aura

	shamanOwner *Shaman
}

func (shaman *Shaman) NewFireElemental(bonusSpellPower float64) *FireElemental {
	basestats := core.PetBaseStats[core.Pet_GreaterFireElemental][shaman.Level].Stats

	fireElemental := &FireElemental{
		Pet:         core.NewPet("Greater Fire Elemental", &shaman.Character, basestats, fireElementalPetBasePercentageStats, shaman.fireElementalStatInheritance(), false, true),
		shamanOwner: shaman,
	}
	fireElemental.Class = proto.Class_ClassPaladin
	fireElemental.EnableManaBar(15)
	fireElemental.EnableAutoAttacks(fireElemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:  float64(fireElemental.Level) * 2.5,
			BaseDamageMax:  float64(fireElemental.Level) * 4.5,
			SwingSpeed:     2,
			CritMultiplier: 2, // Pretty sure this is right.
			SpellSchool:    core.SpellSchoolFire,
		},
		AutoSwingMelee: true,
	})

	if bonusSpellPower > 0 {
		fireElemental.AddStat(stats.SpellPower, float64(bonusSpellPower))
		fireElemental.AddStat(stats.AttackPower, float64(bonusSpellPower)*3)
	}

	if shaman.hasHeroicPresence || shaman.Race == proto.Race_RaceDraenei {
		fireElemental.AddStats(stats.Stats{
			stats.MeleeHit:  -shaman.MeleeHitRatingPerHitChance,
			stats.SpellHit:  -shaman.SpellHitRatingPerHitChance,
			stats.Expertise: -shaman.SpellHitRatingPerHitChance * 17.0 / 26.0 * shaman.ExpertisePerQuarterPercentReduction / shaman.SpellHitRatingPerHitChance,
		})
	}

	fireElemental.OnPetEnable = fireElemental.enable
	fireElemental.OnPetDisable = fireElemental.disable

	shaman.AddPet(fireElemental)

	return fireElemental
}

func (fireElemental *FireElemental) enable(sim *core.Simulation) {
	fireElemental.FireShieldAura.Activate(sim)
}

func (fireElemental *FireElemental) disable(sim *core.Simulation) {
	fireElemental.FireShieldAura.Deactivate(sim)
}

func (fireElemental *FireElemental) GetPet() *core.Pet {
	return &fireElemental.Pet
}

func (fireElemental *FireElemental) Initialize() {

	fireElemental.registerFireBlast()
	fireElemental.registerFireNova()
	fireElemental.registerFireShieldAura()
}

func (fireElemental *FireElemental) Reset(_ *core.Simulation) {

}

func (fireElemental *FireElemental) ExecuteCustomRotation(sim *core.Simulation) {
	/*
		TODO this is a little dirty, can probably clean this up, the rotation might go through some more overhauls,
		the random AI is hard to emulate.
	*/
	target := fireElemental.CurrentTarget
	fireBlastCasts := fireElemental.FireBlast.SpellMetrics[0].Casts
	fireNovaCasts := fireElemental.FireNova.SpellMetrics[0].Casts

	if fireBlastCasts == maxFireBlastCasts && fireNovaCasts == maxFireNovaCasts {
		return
	}

	if fireElemental.FireNova.DefaultCast.Cost > fireElemental.CurrentMana() {
		return
	}

	random := sim.RandomFloat("Fire Elemental Pet Spell")

	//Melee the other 30%
	if random >= .65 {
		if !fireElemental.TryCast(sim, target, fireElemental.FireNova, maxFireNovaCasts) {
			fireElemental.TryCast(sim, target, fireElemental.FireBlast, maxFireBlastCasts)
		}
	} else if random >= .35 {
		if !fireElemental.TryCast(sim, target, fireElemental.FireBlast, maxFireBlastCasts) {
			fireElemental.TryCast(sim, target, fireElemental.FireNova, maxFireNovaCasts)
		}
	}

	if !fireElemental.GCD.IsReady(sim) {
		return
	}

	fireElemental.WaitUntil(sim, sim.CurrentTime+time.Second)
}

func (fireElemental *FireElemental) TryCast(sim *core.Simulation, target *core.Unit, spell *core.Spell, maxCastCount int32) bool {
	if maxCastCount == spell.SpellMetrics[0].Casts {
		return false
	}

	if !spell.Cast(sim, target) {
		return false
	}
	// all spell casts reset the elemental's swing timer
	fireElemental.AutoAttacks.StopMeleeUntil(sim, sim.CurrentTime+spell.CurCast.CastTime, false)
	return true
}

var fireElementalPetBasePercentageStats = stats.Stats{
	// TODO : Log digging and my own samples this seems to be around the 5% mark.
	stats.MeleeCrit: (5 + 1.8),
	stats.SpellCrit: 2.61,
}

func (shaman *Shaman) fireElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats, pseudoStats stats.PseudoStats) stats.Stats {
		spellPower := ownerStats[stats.SpellPower] + pseudoStats.FireSpellPower
		return stats.Stats{
			stats.Stamina:     ownerStats[stats.Stamina] * 0.3,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.30,
			stats.SpellPower:  spellPower,
			stats.AttackPower: spellPower * 3,

			stats.MeleeHit:  shaman.CalculateHitInheritance(stats.SpellHit, stats.MeleeHit),
			stats.SpellHit:  ownerStats[stats.SpellHit],
			stats.Expertise: shaman.CalculateHitInheritance(stats.SpellHit, stats.Expertise),
		}
	}
}
