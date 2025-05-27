import * as AURAS from "../../constants/auras";
import { Expansion, Stat, TristateEffect } from "../../proto/common";
import { ActionId, ActionIDMap } from "../../proto_utils/action_id";
import {
  makeBooleanDebuffInput,
  makeBooleanIndividualBuffInput,
  makeBooleanPartyBuffInput,
  makeBooleanRaidBuffInput,
  makeMultistateIndividualBuffInput,
  makeMultistateMultiplierIndividualBuffInput,
  makeMultistatePartyBuffInput,
  makeMultistateRaidBuffInput,
  makeQuadstateDebuffInput,
  makeTristateDebuffInput,
  makeTristateIndividualBuffInput,
  makeTristateRaidBuffInput,
  withLabel,
} from "../icon_inputs";
import { IconPicker } from "../icon_picker";
import * as InputHelpers from '../input_helpers';
import { MultiIconPicker } from "../multi_icon_picker";
import { IconPickerStatOption, PickerStatOptions } from "./stat_options";

///////////////////////////////////////////////////////////////////////////
//                                 RAID BUFFS
///////////////////////////////////////////////////////////////////////////

export const AllStatsBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.MARKOFTHEWILD), impId: ActionId.fromSpellId(17051), fieldName: 'giftOfTheWild'}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(49634), fieldName: 'drumsOfTheWild', showWhen: player => player.getLevel() >= 80}),
], 'Stats');

export const AllStatsPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25898), fieldName: 'blessingOfKings', showWhen: player => player.getLevel() >= 20}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromItemId(49633), fieldName: 'drumsOfForgottenKings', showWhen: player => player.getLevel() >= 80}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary', showWhen: player => player.getLevel() >= 30}),
], 'Stats %');

export const ArmorBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.DEVOTIONAURA), impId: ActionId.fromSpellId(20140), fieldName: 'devotionAura'}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.STONESKINTOTEM), impId: ActionId.fromSpellId(16293), fieldName: 'stoneskinTotem', showWhen: player => player.getLevel() >= 4}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.ARMOR), fieldName: 'scrollOfProtection'}),
], 'Armor');

export const AttackPowerBuff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.BLESSINGOFMIGHT), impId: ActionId.fromSpellId(20045), fieldName: 'blessingOfMight', showWhen: player => player.getLevel() >= 4}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.BATTLESHOUT), impId: ActionId.fromSpellId(12861), fieldName: 'battleShout'}),
], 'AP');

export const AttackPowerPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(53138), fieldName: 'abominationsMight', showWhen: player => player.getLevel() >= 58}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(30809), fieldName: 'unleashedRage', showWhen: player => player.getLevel() >= 37}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(19506), fieldName: 'trueshotAura', showWhen: player => player.getLevel() >= 40}),
], 'Atk Pwr %');

export const Bloodlust = withLabel(
  makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(2825), fieldName: 'bloodlust', showWhen: player => player.getLevel() >= 70}),
  'Lust',
);

export const DamagePercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(31869), fieldName: 'sanctifiedRetribution', showWhen: player => player.getLevel() >= 30}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(31583), fieldName: 'arcaneEmpowerment', showWhen: player => player.getLevel() >= 42}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(34460), fieldName: 'ferociousInspiration', showWhen: player => player.getLevel() >= 42}),
], 'Dmg %');

export const DamageReductionPercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(57472), fieldName: 'renewedHope', showWhen: player => player.getLevel() >= 45}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(25899), fieldName: 'blessingOfSanctuary', showWhen: player => player.getLevel() >= 30}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(50720), fieldName: 'vigilance', showWhen: player => player.getLevel() >= 40}),
], 'Mit %');

export const DefensiveCooldownBuff = InputHelpers.makeMultiIconInput([
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(6940), numStates: 11, fieldName: 'handOfSacrifices', showWhen: player => player.getLevel() >= 46}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(53530), numStates: 11, fieldName: 'divineGuardians', showWhen: player => player.getLevel() >= 26}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(33206), numStates: 11, fieldName: 'painSuppressions', showWhen: player => player.getLevel() >= 50}),
	makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(47788), numStates: 11, fieldName: 'guardianSpirits', showWhen: player => player.getLevel() >= 60}),
], 'Defensive CDs');

export const HastePercentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(53648), fieldName: 'swiftRetribution', showWhen: player => player.getLevel() >= 52}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(48396), fieldName: 'moonkinAura', value: TristateEffect.TristateEffectImproved, showWhen: player => player.getLevel() >= 43}),
], 'Haste %');

export const HealthBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.COMMANDINGSHOUT), impId: ActionId.fromSpellId(12861), fieldName: 'commandingShout', showWhen: player => player.getLevel() >= 68}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.BLOODPACT), impId: ActionId.fromSpellId(18696), fieldName: 'bloodPact', showWhen: player => player.getLevel() >= 4}),
], 'Health');

export const IntellectBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.ARCANEINTELLECT), fieldName: 'arcaneBrilliance'}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FELINTELLIGENCE), impId: ActionId.fromSpellId(54038), fieldName: 'felIntelligence', showWhen: player => player.getLevel() >= 32}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.INTELLECT), fieldName: 'scrollOfIntellect'}),
], 'Int');

export const MeleeCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(17007), impId: ActionId.fromSpellId(34300), fieldName: 'leaderOfThePack', showWhen: player => player.getLevel() >= 40}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(29801), fieldName: 'rampage', showWhen: player => player.getLevel() >= 50}),
], 'Melee Crit');

export const MeleeHasteBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(55610), fieldName: 'icyTalons', showWhen: player => player.getLevel() >= 58}),
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(65990), impId: ActionId.fromSpellId(29193), fieldName: 'windfuryTotem', showWhen: player => player.getLevel() >= 30}),
], 'Melee Haste');

export const MP5Buff = InputHelpers.makeMultiIconInput([
	makeTristateIndividualBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.BLESSINGOFWISDOM), impId: ActionId.fromSpellId(20245), fieldName: 'blessingOfWisdom', showWhen: player => player.getLevel() >= 14}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.MANASPRINGTOTEM), impId: ActionId.fromSpellId(16206), fieldName: 'manaSpringTotem', showWhen: player => player.getLevel() >= 26}),
], 'MP5');

export const ReplenishmentBuff = InputHelpers.makeMultiIconInput([
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(48160), fieldName: 'vampiricTouch', showWhen: player => player.getLevel() >= 50}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(31878), fieldName: 'judgementsOfTheWise', showWhen: player => player.getLevel() >= 40}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(53292), fieldName: 'huntingParty', showWhen: player => player.getLevel() >= 55}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(54118), fieldName: 'improvedSoulLeech', showWhen: player => player.getLevel() >= 45}),
	makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(44561), fieldName: 'enduringWinter', showWhen: player => player.getLevel() >= 51}),
], 'Replen');

export const ResistanceBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.SHADOWPROTECTION), fieldName: 'shadowProtection', showWhen: player => player.getLevel() >= 30}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.SHADOWRESISTANCEAURA), fieldName: 'shadowResistanceAura', showWhen: player => player.getLevel() >= 24}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.NATURERESISTANCETOTEM), fieldName: 'natureResistanceTotem', showWhen: player => player.getLevel() >= 30}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.ASPECTOFTHEWILD), fieldName: 'aspectOfTheWild', showWhen: player => player.getLevel() >= 46}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FROSTRESISTANCEAURA), fieldName: 'frostResistanceAura', showWhen: player => player.getLevel() >= 32}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FROSTRESISTANCETOTEM), fieldName: 'frostResistanceTotem', showWhen: player => player.getLevel() >= 24}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FIRERESISTANCETOTEM), fieldName: 'fireResistanceTotem', showWhen: player => player.getLevel() >= 24}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FIRERESISTANCEAURA), fieldName: 'fireResistanceAura', showWhen: player => player.getLevel() >= 24}),
], 'Resistances');

export const RevitalizeBuff = InputHelpers.makeMultiIconInput([
  makeMultistateMultiplierIndividualBuffInput(ActionId.fromSpellId(26982), 101, 10, 'revitalizeRejuvination', player => player.getLevel() >= 42),
  makeMultistateMultiplierIndividualBuffInput(ActionId.fromSpellId(53251), 101, 10, 'revitalizeWildGrowth', player => player.getLevel() >= 60),
], 'Revit', ActionId.fromSpellId(48545))

export const SpellCritBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionId.fromSpellId(24907), impId: ActionId.fromSpellId(48396), fieldName: 'moonkinAura', showWhen: player => player.getLevel() >= 40}),
	makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(51470), fieldName: 'elementalOath', showWhen: player => player.getLevel() >= 46}),
], 'Spell Crit');

export const SpellHasteBuff = withLabel(
  makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(3738), fieldName: 'wrathOfAirTotem', showWhen: player => player.getLevel() >= 64}),
  'Spell Haste',
);

export const SpellPowerBuff = InputHelpers.makeMultiIconInput([
	makeMultistateRaidBuffInput({actionId: ActionId.fromSpellId(47240), numStates: 2000, fieldName: 'demonicPactSp', multiplier: 20, showWhen: player => player.getLevel() >= 55}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.TOTEMOFWRATH), fieldName: 'totemOfWrath', showWhen: player => player.getLevel() >= 50}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FLAMETONGUETOTEM), fieldName: 'flametongueTotem', showWhen: player => player.getLevel() >= 28}),
], 'Spell Power');

export const SpiritBuff = InputHelpers.makeMultiIconInput([
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.DIVINESPIRIT), fieldName: 'divineSpirit', showWhen: player => player.getLevel() >= 30}),
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.FELINTELLIGENCE), impId: ActionId.fromSpellId(54038), fieldName: 'felIntelligence', showWhen: player => player.getLevel() >= 32}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.SPIRIT), fieldName: 'scrollOfSpirit'}),
], 'Spirit');

export const StaminaBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.POWERWORDFORTITUDE), impId: ActionId.fromSpellId(14767), fieldName: 'powerWordFortitude'}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.STAMINA), fieldName: 'scrollOfStamina'}),
], 'Stamina');

export const StrengthAndAgilityBuff = InputHelpers.makeMultiIconInput([
	makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.STRENGTHOFEARTHTOTEM), impId: ActionId.fromSpellId(52456), fieldName: 'strengthOfEarthTotem', showWhen: player => player.getLevel() >= 10}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.HORNOFWINTER), fieldName: 'hornOfWinter', showWhen: player => player.getLevel() >= 65}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.AGILITY), fieldName: 'scrollOfAgility'}),
	makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.STRENGTH), fieldName: 'scrollOfStrength'}),
], 'Str/Agi');

// Misc Buffs
export const StrengthOfWrynn = makeBooleanRaidBuffInput({actionId: ActionId.fromSpellId(73828), fieldName: 'strengthOfWrynn', showWhen: player => player.getLevel() >= 80});
export const RetributionAura = makeBooleanRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.RETRIBUTIONAURA), fieldName: 'retributionAura', showWhen: player => player.getLevel() >= 16});
export const BraidedEterniumChain = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31025), fieldName: 'braidedEterniumChain', showWhen: player => player.getLevel() >= 70});
export const ChainOfTheTwilightOwl = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31035), fieldName: 'chainOfTheTwilightOwl', showWhen: player => player.getLevel() >= 70});
export const HeroicPresence = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(6562), fieldName: 'heroicPresence'});
export const EyeOfTheNight = makeBooleanPartyBuffInput({actionId: ActionId.fromSpellId(31033), fieldName: 'eyeOfTheNight', showWhen: player => player.getLevel() >= 70});
export const Thorns = makeTristateRaidBuffInput({actionId: ActionIDMap.fromSpellId(AURAS.THORNS), impId: ActionId.fromSpellId(16840), fieldName: 'thorns', showWhen: player => player.getLevel() >= 6});
export const ManaTideTotem = makeMultistatePartyBuffInput(ActionId.fromSpellId(16190), 5, 'manaTideTotems', party => party.getPlayer(0)!.getLevel() >= 40);
export const Innervate = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(29166), numStates: 11, fieldName: 'innervates', showWhen: player => player.getLevel() >= 40});
export const PowerInfusion = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(10060), numStates: 11, fieldName: 'powerInfusions', showWhen: player => player.getLevel() >= 40});
export const FocusMagic = makeBooleanIndividualBuffInput({actionId: ActionId.fromSpellId(54648), fieldName: 'focusMagic', showWhen: player => player.getLevel() >= 20});
export const TricksOfTheTrade = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(57933), numStates: 20, fieldName: 'tricksOfTheTrades', showWhen: player => player.getLevel() >= 75});
export const UnholyFrenzy = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(49016), numStates: 11, fieldName: 'unholyFrenzy', showWhen: player => player.getLevel() >= 58});

///////////////////////////////////////////////////////////////////////////
//                                 DEBUFFS
///////////////////////////////////////////////////////////////////////////

export const MajorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(7386), fieldName: 'sunderArmor', showWhen: player => player.getLevel() >= 10}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(8647), fieldName: 'exposeArmor', showWhen: player => player.getLevel() >= 14}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(55754), fieldName: 'acidSpit', showWhen: player => player.getLevel() >= 60}),
], 'Major ArP');

export const MinorArmorDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(770), impId: ActionId.fromSpellId(33602), fieldName: 'faerieFire', showWhen: player => player.getLevel() >= 18}),
	makeTristateDebuffInput({actionId:ActionIDMap.fromSpellId(AURAS.CURSEOFWEAKNESS), impId: ActionId.fromSpellId(18180), fieldName: 'curseOfWeakness', showWhen: player => player.getLevel() >= 4}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(56631), fieldName: 'sting', showWhen: player => player.getLevel() >= 60 && player.getExpansion() >= Expansion.ExpansionTbc}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53598), fieldName: 'sporeCloud', showWhen: player => player.getLevel() >= 60 && player.getExpansion() >= Expansion.ExpansionTbc}),
], 'Minor ArP');

export const AttackPowerDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(26016), fieldName: 'vindication', showWhen: player => player.getLevel() >= 21}),
	makeTristateDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.DEMORALIZINGSHOUT), impId: ActionId.fromSpellId(12879), fieldName: 'demoralizingShout', showWhen: player => player.getLevel() >= 14}),
	makeTristateDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.DEMORALIZINGROAR), impId: ActionId.fromSpellId(16862), fieldName: 'demoralizingRoar', showWhen: player => player.getLevel() >= 10}),
	makeTristateDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.CURSEOFWEAKNESS), impId: ActionId.fromSpellId(18180), fieldName: 'curseOfWeakness', showWhen: player => player.getLevel() >= 4}),
	makeBooleanDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.DEMORALIZINGSCREECH), fieldName: 'demoralizingScreech', showWhen: player => player.getLevel() >= 10}),
], 'Atk Pwr');

export const BleedDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48564), fieldName: 'mangle', showWhen: player => player.getLevel() >= 50}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(46855), fieldName: 'trauma', showWhen: player => player.getLevel() >= 36}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(57393), fieldName: 'stampede', showWhen: player => player.getLevel() >= 68 && player.getExpansion() == Expansion.ExpansionWotlk}),
], 'Bleed');

export const CritDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.TOTEMOFWRATH), fieldName: 'totemOfWrath', showWhen: player => player.getLevel() >= 50}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(20337), fieldName: 'heartOfTheCrusader', showWhen: player => player.getLevel() >= 17}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(58410), fieldName: 'masterPoisoner', showWhen: player => player.getLevel() >= 52}),
], 'Crit');

export const MeleeAttackSpeedDebuff = InputHelpers.makeMultiIconInput([
	makeTristateDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.THUNDERCLAP), impId: ActionId.fromSpellId(12666), fieldName: 'thunderClap', showWhen: player => player.getLevel() >= 6}),
	makeTristateDebuffInput({actionId: ActionId.fromSpellId(55095), impId: ActionId.fromSpellId(51456), fieldName: 'frostFever', showWhen: player => player.getLevel() >= 58}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53696), fieldName: 'judgementsOfTheJust', showWhen: player => player.getLevel() >= 56}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48485), fieldName: 'infectedWounds', showWhen: player => player.getLevel() >= 47}),
], 'Atk Speed');

export const MeleeHitDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(65855), fieldName: 'insectSwarm', showWhen: player => player.getLevel() >= 30}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(3043), fieldName: 'scorpidSting', showWhen: player => player.getLevel() >= 22}),
], 'Miss');

export const PhysicalDamageDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(29859), fieldName: 'bloodFrenzy', showWhen: player => player.getLevel() >= 51}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(58413), fieldName: 'savageCombat', showWhen: player => player.getLevel() >= 51}),
], 'Phys Vuln');

export const SpellCritDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(17803), fieldName: 'shadowMastery', showWhen: player => player.getLevel() >= 10}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(12873), fieldName: 'improvedScorch', showWhen: player => player.getLevel() >= 25}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(28593), fieldName: 'wintersChill', showWhen: player => player.getLevel() >= 35}),
], 'Spell Crit');

export const SpellHitDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(33198), fieldName: 'misery', showWhen: player => player.getLevel() >= 47}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(33602), fieldName: 'faerieFire', value: TristateEffect.TristateEffectImproved, showWhen: player => player.getLevel() >= 42}),
], 'Spell Hit');

export const SpellDamageDebuff = InputHelpers.makeMultiIconInput([
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(51161), fieldName: 'ebonPlaguebringer', showWhen: player => player.getLevel() >= 58}),
	makeBooleanDebuffInput({actionId: ActionId.fromSpellId(48511), fieldName: 'earthAndMoon', showWhen: player => player.getLevel() >= 57}),
	makeBooleanDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.CURSEOFTHEELEMENTS), fieldName: 'curseOfElements', showWhen: player => player.getLevel() >= 32}),
], 'Spell Dmg');

export const HuntersMark = withLabel(
  makeQuadstateDebuffInput({actionId: ActionIDMap.fromSpellId(AURAS.HUNTERSMARK), impId: ActionId.fromSpellId(19423), impId2: ActionId.fromItemId(42907), fieldName: 'huntersMark', showWhen: player => player.getLevel() >= 6}),
  'Mark',
);
export const JudgementOfWisdom = withLabel(makeBooleanDebuffInput({actionId: ActionId.fromSpellId(53408), fieldName: 'judgementOfWisdom', showWhen: player => player.getLevel() >= 12}), 'JoW');
export const JudgementOfLight = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(20271), fieldName: 'judgementOfLight', showWhen: player => player.getLevel() >= 4});
export const ShatteringThrow = makeMultistateIndividualBuffInput({actionId: ActionId.fromSpellId(64382), numStates: 20, fieldName: 'shatteringThrows', showWhen: player => player.getLevel() >= 80});
export const GiftOfArthas = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(11374), fieldName: 'giftOfArthas', showWhen: player => player.getLevel() >= 38});
export const CrystalYield = makeBooleanDebuffInput({actionId: ActionId.fromSpellId(15235), fieldName: 'crystalYield', showWhen: player => player.getLevel() >= 47});

///////////////////////////////////////////////////////////////////////////
//                                 CONFIGS
///////////////////////////////////////////////////////////////////////////

export const RAID_BUFFS_CONFIG = [
	// Standard buffs
	{
		config: AllStatsBuff,
		picker: MultiIconPicker,
		stats: [],
	},
  {
		config: AllStatsPercentBuff,
		picker: MultiIconPicker,
		stats: [],
	},
  {
    config: HealthBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatHealth],
  },
	{
		config: ArmorBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor],
	},
	{
		config: StaminaBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStamina],
	},
	{
		config: StrengthAndAgilityBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatStrength, Stat.StatAgility],
	},
	{
		config: IntellectBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatIntellect],
	},
	{
		config: SpiritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpirit],
	},
	{
		config: AttackPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: AttackPowerPercentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
	},
	{
		config: MeleeCritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMeleeCrit],
	},
  {
    config: MeleeHasteBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatMeleeHaste],
  },
	{
		config: SpellPowerBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower],
	},
	{
		config: SpellCritBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellCrit],
	},
  {
    config: HastePercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],
  },
  {
    config: DamagePercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: DamageReductionPercentBuff,
    picker: MultiIconPicker,
    stats: [Stat.StatArmor],
  },
	{
		config: ResistanceBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatNatureResistance, Stat.StatShadowResistance, Stat.StatFrostResistance]
		},
	{
		config: DefensiveCooldownBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
	{
		config: MP5Buff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5]
	},
	{
		config: ReplenishmentBuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMP5]
	},
  {
    config: Bloodlust,
    picker: IconPicker,
    stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],
  },
  {
    config: SpellHasteBuff,
    picker: IconPicker,
    stats: [Stat.StatSpellHaste],
  },
  {
    config: RevitalizeBuff,
    picker: MultiIconPicker,
    stats: [],
  },
] as PickerStatOptions[]

export const RAID_BUFFS_MISC_CONFIG = [
  {
    config: StrengthOfWrynn,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: HeroicPresence,
    picker: IconPicker,
    stats: [Stat.StatMeleeHit, Stat.StatSpellHit],
  },
  {
    config: BraidedEterniumChain,
    picker: IconPicker,
    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit],
  },
  {
    config: ChainOfTheTwilightOwl,
    picker: IconPicker,
    stats: [Stat.StatSpellCrit, Stat.StatMeleeCrit],
  },
  {
    config: FocusMagic,
    picker: IconPicker,
    stats: [Stat.StatSpellCrit],
  },
  {
    config: EyeOfTheNight,
    picker: IconPicker,
    stats: [Stat.StatSpellPower],
  },
  {
    config: Thorns,
    picker: IconPicker,
    stats: [Stat.StatArmor],
  },
  {
    config: RetributionAura,
    picker: IconPicker,
    stats: [Stat.StatArmor],
  },
  {
    config: ManaTideTotem,
    picker: IconPicker,
    stats: [Stat.StatMP5],
  },
  {
    config: Innervate,
    picker: IconPicker,
    stats: [Stat.StatMP5],
  },
  {
    config: PowerInfusion,
    picker: IconPicker,
    stats: [Stat.StatMP5, Stat.StatSpellPower],
  },
  {
    config: TricksOfTheTrade,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower, Stat.StatSpellPower],
  },
  {
    config: UnholyFrenzy,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
  },
] as IconPickerStatOption[]

export const DEBUFFS_CONFIG = [
	{
		config: MajorArmorDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmorPenetration]
	},
  {
		config: MinorArmorDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmorPenetration]
	},
  {
		config: PhysicalDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
  {
		config: BleedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower]
	},
  {
		config: SpellDamageDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellPower]
	},
  {
		config: SpellHitDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellHit]
	},
  {
		config: SpellCritDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatSpellCrit]
	},
  {
		config: CritDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit]
	},
  {
		config: AttackPowerDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
  {
		config: MeleeAttackSpeedDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatArmor]
	},
  {
		config: MeleeHitDebuff,
		picker: MultiIconPicker,
		stats: [Stat.StatDodge]
	},
  {
		config: JudgementOfWisdom,
		picker: IconPicker,
		stats: [Stat.StatMP5, Stat.StatIntellect]
	},
	{
		config: HuntersMark,
		picker: IconPicker,
		stats: [Stat.StatRangedAttackPower]
	},
] as PickerStatOptions[];

export const DEBUFFS_MISC_CONFIG = [
  {
    config: JudgementOfLight,
    picker: IconPicker,
    stats: [Stat.StatStamina],
  },
  {
    config: ShatteringThrow,
    picker: IconPicker,
    stats: [Stat.StatArmorPenetration],
  },
  {
    config: GiftOfArthas,
    picker: IconPicker,
    stats: [Stat.StatAttackPower, Stat.StatRangedAttackPower],
  },
  {
    config: CrystalYield,
    picker: IconPicker,
    stats: [Stat.StatArmorPenetration],
  },
] as IconPickerStatOption[];
