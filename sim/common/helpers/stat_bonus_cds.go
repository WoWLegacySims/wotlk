package helpers

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type StatCDFactory func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration)

// Wraps factory functions so that only the first item is included in tests.
func testFirstOnly(factory StatCDFactory) StatCDFactory {
	first := true
	return func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
		if first {
			first = false
			factory(itemID, bonus, duration, cooldown)
		} else {
			core.AddEffectsToTest = false
			factory(itemID, bonus, duration, cooldown)
			core.AddEffectsToTest = true
		}
	}
}

var NewHasteActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.MeleeHaste: bonus, stats.SpellHaste: bonus}, duration, cooldown)
})

var NewAttackPowerActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.AttackPower: bonus, stats.RangedAttackPower: bonus}, duration, cooldown)
})

var NewSpellPowerActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.SpellPower: bonus}, duration, cooldown)
})

var NewArmorPenActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatOffensiveTrinketEffect(itemID, stats.Stats{stats.ArmorPenetration: bonus}, duration, cooldown)
})

var NewHealthActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Health: bonus}, duration, cooldown)
})

var NewArmorActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.BonusArmor: bonus}, duration, cooldown)
})

var NewBlockValueActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	//core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.BlockValue: bonus}, duration, cooldown)
	// Hack for Lavanthor's Talisman Shared CD being shorter than its effect
	core.NewSimpleStatItemActiveEffect(itemID, stats.Stats{stats.BlockValue: bonus}, duration, cooldown, func(character *core.Character) core.Cooldown {
		return core.Cooldown{
			Timer:    character.GetDefensiveTrinketCD(),
			Duration: time.Second * 20,
		}
	}, nil)
})

var NewDodgeActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Dodge: bonus}, duration, cooldown)
})

var NewParryActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Parry: bonus}, duration, cooldown)
})

var NewSpiritActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.Spirit: bonus}, duration, cooldown)
})

var NewResistsActive = testFirstOnly(func(itemID int32, bonus float64, duration time.Duration, cooldown time.Duration) {
	core.NewSimpleStatDefensiveTrinketEffect(itemID, stats.Stats{stats.ArcaneResistance: bonus, stats.FireResistance: bonus, stats.FrostResistance: bonus, stats.NatureResistance: bonus, stats.ShadowResistance: bonus}, duration, cooldown)
})
