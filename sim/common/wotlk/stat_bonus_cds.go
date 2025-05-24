package wotlk

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
)

func init() {
	// Keep these separated by stat, ordered by item ID within each group.

	helpers.NewHasteActive(36972, 256, time.Second*20, time.Minute*2) // Tome of Arcane Phenomena
	helpers.NewHasteActive(37558, 122, time.Second*20, time.Minute*2) // Tidal Boon
	helpers.NewHasteActive(37560, 124, time.Second*20, time.Minute*2) // Vial of Renewal
	helpers.NewHasteActive(37562, 140, time.Second*20, time.Minute*2) // Fury of the Crimson Drake
	helpers.NewHasteActive(38070, 148, time.Second*20, time.Minute*2) // Foresight's Anticipation
	helpers.NewHasteActive(38258, 140, time.Second*20, time.Minute*2) // Sailor's Knotted Charm
	helpers.NewHasteActive(38259, 140, time.Second*20, time.Minute*2) // First Mate's Pocketwatch
	helpers.NewHasteActive(38764, 208, time.Second*20, time.Minute*2) // Rune of Finite Variation
	helpers.NewHasteActive(40531, 491, time.Second*20, time.Minute*2) // Mark of Norgannon
	helpers.NewHasteActive(43836, 212, time.Second*20, time.Minute*2) // Thorny Rose Brooch
	helpers.NewHasteActive(45466, 457, time.Second*20, time.Minute*2) // Scale of Fates
	helpers.NewHasteActive(46088, 375, time.Second*20, time.Minute*2) // Platinum Disks of Swiftness
	helpers.NewHasteActive(48722, 512, time.Second*20, time.Minute*2) // Shard of the Crystal Heart

	helpers.NewAttackPowerActive(35937, 328, time.Second*20, time.Minute*2)  // Braxley's Backyard Moonshine
	helpers.NewAttackPowerActive(36871, 280, time.Second*20, time.Minute*2)  // Fury of the Encroaching Storm
	helpers.NewAttackPowerActive(37166, 670, time.Second*20, time.Minute*2)  // Sphere of Red Dragon's Blood
	helpers.NewAttackPowerActive(37556, 248, time.Second*20, time.Minute*2)  // Talisman of the Tundra
	helpers.NewAttackPowerActive(37557, 304, time.Second*20, time.Minute*2)  // Warsong's Fervor
	helpers.NewAttackPowerActive(38080, 264, time.Second*20, time.Minute*2)  // Automated Weapon Coater
	helpers.NewAttackPowerActive(38081, 280, time.Second*20, time.Minute*2)  // Scarab of Isanoth
	helpers.NewAttackPowerActive(38761, 248, time.Second*20, time.Minute*2)  // Talon of Hatred
	helpers.NewAttackPowerActive(39257, 670, time.Second*20, time.Minute*2)  // Loatheb's Shadow
	helpers.NewAttackPowerActive(44014, 432, time.Second*15, time.Minute*2)  // Fezzik's Pocketwatch
	helpers.NewAttackPowerActive(45263, 905, time.Second*20, time.Minute*2)  // Wrathstone
	helpers.NewAttackPowerActive(46086, 752, time.Second*20, time.Minute*2)  // Platinum Disks of Battle
	helpers.NewAttackPowerActive(47734, 1024, time.Second*20, time.Minute*2) // Mark of Supremacy

	helpers.NewSpellPowerActive(35935, 178, time.Second*20, time.Minute*2) // Infused Coldstone Rune
	helpers.NewSpellPowerActive(36872, 173, time.Second*20, time.Minute*2) // Mender of the Oncoming Dawn
	helpers.NewSpellPowerActive(36874, 183, time.Second*20, time.Minute*2) // Horn of the Herald
	helpers.NewSpellPowerActive(37555, 149, time.Second*20, time.Minute*2) // Warsong's Wrath
	helpers.NewSpellPowerActive(37844, 346, time.Second*20, time.Minute*2) // Winged Talisman
	helpers.NewSpellPowerActive(37873, 346, time.Second*20, time.Minute*2) // Mark of the War Prisoner
	helpers.NewSpellPowerActive(38073, 120, time.Second*15, time.Minute*2) // Will of the Red Dragonflight
	helpers.NewSpellPowerActive(38213, 149, time.Second*20, time.Minute*2) // Harbringer's Wrath
	helpers.NewSpellPowerActive(38527, 183, time.Second*20, time.Minute*2) // Strike of the Seas
	helpers.NewSpellPowerActive(38760, 145, time.Second*20, time.Minute*2) // Mendicant's Charm
	helpers.NewSpellPowerActive(38762, 145, time.Second*20, time.Minute*2) // Insignia of Bloody Fire
	helpers.NewSpellPowerActive(38765, 202, time.Second*20, time.Minute*2) // Rune of Infinite Power
	helpers.NewSpellPowerActive(39811, 183, time.Second*20, time.Minute*2) // Badge of the Infiltrator
	helpers.NewSpellPowerActive(39819, 145, time.Second*20, time.Minute*2) // Bloodbinder's Runestone
	helpers.NewSpellPowerActive(39821, 145, time.Second*20, time.Minute*2) // Spiritist's Focus
	helpers.NewSpellPowerActive(42395, 292, time.Second*20, time.Minute*2) // Figurine - Twilight Serpent
	helpers.NewSpellPowerActive(43837, 281, time.Second*20, time.Minute*2) // Soflty Glowing Orb
	helpers.NewSpellPowerActive(44013, 281, time.Second*20, time.Minute*2) // Cannoneer's Fuselighter
	helpers.NewSpellPowerActive(44015, 281, time.Second*20, time.Minute*2) // Cannoneer's Morale
	helpers.NewSpellPowerActive(45148, 534, time.Second*20, time.Minute*2) // Living Flame
	helpers.NewSpellPowerActive(45292, 431, time.Second*20, time.Minute*2) // Energy Siphon
	helpers.NewSpellPowerActive(46087, 440, time.Second*20, time.Minute*2) // Platinum Disks of Sorcery
	helpers.NewSpellPowerActive(48724, 599, time.Second*20, time.Minute*2) // Talisman of Resurgence
	helpers.NewSpellPowerActive(50357, 716, time.Second*20, time.Minute*2) // Maghia's Misguided Quill

	helpers.NewArmorPenActive(37723, 291, time.Second*20, time.Minute*2) // Incisor Fragment

	helpers.NewHealthActive(37638, 3025, time.Second*15, time.Minute*3) // Offering of Sacrifice
	helpers.NewHealthActive(39292, 3025, time.Second*15, time.Minute*3) // Repelling Charge
	helpers.NewHealthActive(42128, 3385, time.Second*15, time.Minute*3) // Battlemaster's Hostility
	helpers.NewHealthActive(42129, 3385, time.Second*15, time.Minute*3) // Battlemaster's Accuracy
	helpers.NewHealthActive(42130, 3385, time.Second*15, time.Minute*3) // Battlemaster's Avidity
	helpers.NewHealthActive(42131, 3385, time.Second*15, time.Minute*3) // Battlemaster's Conviction
	helpers.NewHealthActive(42132, 3385, time.Second*15, time.Minute*3) // Battlemaster's Bravery
	helpers.NewHealthActive(42133, 4608, time.Second*15, time.Minute*3) // Battlemaster's Fury
	helpers.NewHealthActive(42134, 4608, time.Second*15, time.Minute*3) // Battlemaster's Precision
	helpers.NewHealthActive(42135, 4608, time.Second*15, time.Minute*3) // Battlemaster's Vivacity
	helpers.NewHealthActive(42136, 4608, time.Second*15, time.Minute*3) // Battlemaster's Rage
	helpers.NewHealthActive(42137, 4608, time.Second*15, time.Minute*3) // Battlemaster's Ruination
	helpers.NewHealthActive(47080, 4610, time.Second*15, time.Minute*3) // Satrina's Impeding Scarab
	helpers.NewHealthActive(47088, 5186, time.Second*15, time.Minute*3) // Satrina's Impeding Scarab H
	helpers.NewHealthActive(47290, 4610, time.Second*15, time.Minute*3) // Juggernaut's Vitality
	helpers.NewHealthActive(47451, 5186, time.Second*15, time.Minute*3) // Juggernaut's Vitality H
	helpers.NewHealthActive(50235, 4104, time.Second*15, time.Minute*3) // Ick's Rotting Thumb

	helpers.NewArmorActive(36993, 3570, time.Second*20, time.Minute*2) // Seal of the Pantheon
	helpers.NewArmorActive(45313, 5448, time.Second*20, time.Minute*2) // Furnace Stone

	helpers.NewBlockValueActive(37872, 440, time.Second*40, time.Minute*2) // Lavanthor's Talisman

	helpers.NewDodgeActive(40257, 455, time.Second*20, time.Minute*2) // Defender's Code
	helpers.NewDodgeActive(40683, 335, time.Second*20, time.Minute*2) // Valor Medal of the First War
	helpers.NewDodgeActive(44063, 300, time.Second*10, time.Minute*1) // Figurine - Monarch Crab
	helpers.NewDodgeActive(45158, 457, time.Second*20, time.Minute*2) // Heart of Iron
	helpers.NewDodgeActive(49080, 335, time.Second*20, time.Minute*2) // Brawler's Souvenir
	helpers.NewDodgeActive(47735, 512, time.Second*20, time.Minute*2) // Glyph of Indomitability

	helpers.NewParryActive(40372, 375, time.Second*20, time.Minute*2) // Rune of Repulsion
	helpers.NewParryActive(46021, 402, time.Second*20, time.Minute*2) // Royal Seal of King Llane

	helpers.NewSpiritActive(38763, 184, time.Second*20, time.Minute*2) // Futuresight Rune
	helpers.NewSpiritActive(39388, 336, time.Second*20, time.Minute*2) // Spirit-World Glass

	helpers.NewResistsActive(50361, 239, time.Second*10, time.Minute*1) // Sindragona's Flawless Fang
	helpers.NewResistsActive(50364, 268, time.Second*10, time.Minute*1) // Sindragona's Flawless Fang H
}
