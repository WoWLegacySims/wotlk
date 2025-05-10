package sim

import (
	_ "github.com/WoWLegacySims/wotlk/sim/common"
	dpsDeathKnight "github.com/WoWLegacySims/wotlk/sim/deathknight/dps"
	tankDeathKnight "github.com/WoWLegacySims/wotlk/sim/deathknight/tank"
	"github.com/WoWLegacySims/wotlk/sim/druid/balance"
	"github.com/WoWLegacySims/wotlk/sim/druid/feral"
	restoDruid "github.com/WoWLegacySims/wotlk/sim/druid/restoration"
	feralTank "github.com/WoWLegacySims/wotlk/sim/druid/tank"
	_ "github.com/WoWLegacySims/wotlk/sim/encounters"
	"github.com/WoWLegacySims/wotlk/sim/hunter"
	"github.com/WoWLegacySims/wotlk/sim/mage"
	holyPaladin "github.com/WoWLegacySims/wotlk/sim/paladin/holy"
	protectionPaladin "github.com/WoWLegacySims/wotlk/sim/paladin/protection"
	"github.com/WoWLegacySims/wotlk/sim/paladin/retribution"
	healingPriest "github.com/WoWLegacySims/wotlk/sim/priest/healing"
	"github.com/WoWLegacySims/wotlk/sim/priest/shadow"
	"github.com/WoWLegacySims/wotlk/sim/priest/smite"
	"github.com/WoWLegacySims/wotlk/sim/rogue"
	"github.com/WoWLegacySims/wotlk/sim/shaman/elemental"
	"github.com/WoWLegacySims/wotlk/sim/shaman/enhancement"
	restoShaman "github.com/WoWLegacySims/wotlk/sim/shaman/restoration"
	"github.com/WoWLegacySims/wotlk/sim/warlock"
	dpsWarrior "github.com/WoWLegacySims/wotlk/sim/warrior/dps"
	protectionWarrior "github.com/WoWLegacySims/wotlk/sim/warrior/protection"
)

var registered = false

func RegisterAll() {
	if registered {
		return
	}
	registered = true

	balance.RegisterBalanceDruid()
	feral.RegisterFeralDruid()
	feralTank.RegisterFeralTankDruid()
	restoDruid.RegisterRestorationDruid()
	elemental.RegisterElementalShaman()
	enhancement.RegisterEnhancementShaman()
	restoShaman.RegisterRestorationShaman()
	hunter.RegisterHunter()
	mage.RegisterMage()
	healingPriest.RegisterHealingPriest()
	shadow.RegisterShadowPriest()
	smite.RegisterSmitePriest()
	rogue.RegisterRogue()
	dpsWarrior.RegisterDpsWarrior()
	protectionWarrior.RegisterProtectionWarrior()
	holyPaladin.RegisterHolyPaladin()
	protectionPaladin.RegisterProtectionPaladin()
	retribution.RegisterRetributionPaladin()
	warlock.RegisterWarlock()
	dpsDeathKnight.RegisterDpsDeathknight()
	tankDeathKnight.RegisterTankDeathknight()
}
