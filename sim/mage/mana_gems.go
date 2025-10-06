package mage

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/mageinfo"
)

type ManaGem struct {
	ID      int32
	Name    string
	Charges int
	Spell   int32
	Level   int32
}

var Managems = [6]*ManaGem{
	{ID: 33312, Name: "Mana Sapphire", Charges: 3, Spell: 42987, Level: 77},
	{ID: 22044, Name: "Mana Emerald", Charges: 3, Spell: 27103, Level: 68},
	{ID: 8008, Name: "Mana Ruby", Charges: 1, Spell: 10058, Level: 58},
	{ID: 8007, Name: "Mana Citrine", Charges: 1, Spell: 10057, Level: 48},
	{ID: 5513, Name: "Mana Jade", Charges: 1, Spell: 10052, Level: 38},
	{ID: 5514, Name: "Mana Agate", Charges: 1, Spell: 5405, Level: 28},
}

func (mage *Mage) registerManaGemsCD() {
	var gem *ManaGem
	var pos int
	for i, g := range Managems {
		if mage.Level >= g.Level {
			gem = g
			pos = i
			break
		}
	}
	if gem == nil {
		return
	}
	dbc := mageinfo.ReplenishMana.GetByID(gem.Spell)
	bp, die := dbc.GetBPDie(0, mage.Level)
	useDownCharges := 0
	totalcharges := gem.Charges

	var bpDown, dieDown = 0.0, 0.0
	if pos < 5 {
		downGem := Managems[pos+1]
		dbc := mageinfo.ReplenishMana.GetByID(downGem.Spell)
		bpDown, dieDown = dbc.GetBPDie(0, mage.Level)
		useDownCharges = downGem.Charges
		totalcharges += useDownCharges
	}

	actionID := core.ActionID{ItemID: gem.ID}
	manaMetrics := mage.NewManaMetrics(actionID)
	hasT7_2pc := mage.HasSetBonus(ItemSetFrostfireGarb, 2)
	var gemAura *core.Aura
	if hasT7_2pc {
		gemAura = mage.NewTemporaryStatsAura("Improved Mana Gems T7", core.ActionID{SpellID: 61062}, stats.Stats{stats.SpellPower: 225}, 15*time.Second)
	}

	var serpentCoilAura *core.Aura
	if mage.HasTrinketEquipped(30720) {
		serpentCoilAura = mage.NewTemporaryStatsAura("Serpent-Coil Braid", core.ActionID{ItemID: 30720}, stats.Stats{stats.SpellPower: 225}, 15*time.Second)
	}

	manaMultiplier := core.TernaryFloat64(mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfManaGem), 1.4, 1) *
		(1 +
			core.TernaryFloat64(serpentCoilAura != nil, 0.25, 0) +
			core.TernaryFloat64(hasT7_2pc, 0.25, 0))

	bpDown = bpDown * manaMultiplier
	dieDown = dieDown * manaMultiplier
	bp = bp * manaMultiplier
	die = die * manaMultiplier

	var remainingManaGems int
	mage.RegisterResetEffect(func(sim *core.Simulation) {
		remainingManaGems = totalcharges
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Minute * 2,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return remainingManaGems != 0
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			var manaGain float64
			if remainingManaGems > useDownCharges {
				// Mana Sapphire: Restores 3330 to 3500 mana. (2 Min Cooldown)
				manaGain = sim.Roll(bp, die)
			} else {
				// Mana Emerald: Restores 2340 to 2460 mana. (2 Min Cooldown)
				manaGain = sim.Roll(bpDown, dieDown)
			}

			if gemAura != nil {
				gemAura.Activate(sim)
			}
			if serpentCoilAura != nil {
				serpentCoilAura.Activate(sim)
			}

			mage.AddMana(sim, manaGain, manaMetrics)

			remainingManaGems--
			if remainingManaGems == 0 {
				// Disable this cooldown since we're out of emeralds.
				mage.GetMajorCooldown(actionID).Disable()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell:    spell,
		Priority: core.CooldownPriorityDefault,
		Type:     core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Only pop if we have less than the max mana provided by the gem minus 1mp5 tick.
			totalRegen := character.ManaRegenPerSecondWhileCasting() * 5
			maxManaGain := bp + die
			if remainingManaGems <= useDownCharges {
				maxManaGain = bpDown + dieDown
			}
			if character.MaxMana()-(character.CurrentMana()+totalRegen) < maxManaGain {
				return false
			}

			return true
		},
	})
}
