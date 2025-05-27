import { Player } from "../../player";
import {
	Class,
	Conjured,
	Consumes,
	Explosive_Big,
	Explosive_Medium,
	Explosive_Small,
	ItemSlot,
	PetFood,
	Potions,
	Profession,
	Spec,
	Stat,
	WeaponImbue,
} from "../../proto/common";
import {BattleElixir,Flask,GuardianElixir} from '../../proto/consumes_gen.js'
import { ActionId } from "../../proto_utils/action_id";
import { EventID, TypedEvent } from "../../typed_event";
import { IconEnumValueConfig } from "../icon_enum_picker";
import { makeBooleanConsumeInput } from "../icon_inputs";
import * as InputHelpers from '../input_helpers';
import { ActionInputConfig, ItemStatOption } from "./stat_options";

export interface ConsumableInputConfig<T> extends ActionInputConfig<T> {
	value: T,
}

export interface ConsumableStatOption<T> extends ItemStatOption<T> {
	config: ConsumableInputConfig<T>,
	level?: number,
	condition?: (player:Player<any>) => boolean
}

export interface ImbueConsumableStatOption extends ConsumableStatOption<WeaponImbue> {
	type?: string,
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void
	showWhen?: (player: Player<any>) => boolean,
	filter?: (option: ConsumableStatOption<T>, player: Player<any>) => boolean
}

function makeConsumeInputFactory<T extends number>(args: ConsumeInputFactoryArgs<T>): (options: ConsumableStatOption<T>[], tooltip?: string) => InputHelpers.TypedIconEnumPickerConfig<Player<any>, T> {
	return (options: ConsumableStatOption<T>[], tooltip?: string) => {
		return {
			type: 'iconEnum',
			tooltip: tooltip,
			numColumns: options.length > 5 ? 2 : 1,
			values: [
				{ value: 0 } as unknown as IconEnumValueConfig<Player<any>, T>,
			].concat(options.map(option => {
				const rtn = {
					actionId: option.config.actionId,
					showWhen: (player: Player<any>) =>
						(!option.config.showWhen || option.config.showWhen(player)) &&
						(option.config.faction || player.getFaction()) == player.getFaction() &&
						(!args.filter || args.filter(option, player))
				} as IconEnumValueConfig<Player<any>, T>;
				if (option.config.value) rtn.value = option.config.value

				return rtn
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.gearChangeEmitter, player.levelChangeEmitter, player.professionChangeEmitter, player.sim.expansionChangeEmitter]),
			showWhen: (player: Player<any>) => !args.showWhen || args.showWhen(player),
			getValue: (player: Player<any>) => player.getConsumes()[args.consumesFieldName] as T,
			setValue: (eventID: EventID, player: Player<any>, newValue: number) => {
				const newConsumes = player.getConsumes();
				if (newConsumes[args.consumesFieldName] === newValue){
					return;
				}

				(newConsumes[args.consumesFieldName] as number) = newValue;
				TypedEvent.freezeAllAndDo(() => {
					player.setConsumes(eventID, newConsumes);
					if (args.onSet) {
						args.onSet(eventID, player, newValue as T);
					}
				});
			},
		};
	};
}


///////////////////////////////////////////////////////////////////////////
//                                 CONJURED
///////////////////////////////////////////////////////////////////////////

export const ConjuredDarkRune = { actionId: ActionId.fromItemId(12662), value: Conjured.ConjuredDarkRune };
export const ConjuredFlameCap = { actionId: ActionId.fromItemId(22788), value: Conjured.ConjuredFlameCap};
export const ConjuredHealthstone = { actionId: ActionId.fromItemId(22105), value: Conjured.ConjuredHealthstone };
export const ConjuredRogueThistleTea = {
  actionId: ActionId.fromItemId(7676),
  value: Conjured.ConjuredRogueThistleTea,
  showWhen: (player: Player<Spec>) => player.getClass() == Class.ClassRogue
};

export const CONJURED_CONFIG = [
  { config: ConjuredRogueThistleTea, stats: [] },
	{ config: ConjuredHealthstone, stats: [Stat.StatStamina] },
  { config: ConjuredDarkRune, stats: [Stat.StatIntellect] },
  { config: ConjuredFlameCap, stats: [] },
] as ConsumableStatOption<Conjured>[]

export const makeConjuredInput = makeConsumeInputFactory({consumesFieldName: 'defaultConjured'});

///////////////////////////////////////////////////////////////////////////
//                                 EXPLOSIVES
///////////////////////////////////////////////////////////////////////////

export const EXPLOSIVES_CONFIG = [
	{ config: { actionId: ActionId.fromItemId(41119), value: Explosive_Small.ExplosiveSaroniteBomb, showWhen: player => player.hasProfession(Profession.Engineering) }, stats: []},
	{ config: { actionId: ActionId.fromItemId(40771), value: Explosive_Small.ExplosiveCobaltFragBomb, showWhen: player => player.hasProfession(Profession.Engineering) }, stats: []},
	{ config: { actionId: ActionId.fromItemId(23826), value: Explosive_Small.TheBiggerOne, showWhen: player => player.hasProfession(Profession.Engineering) }, stats: []},
	{ config: { actionId: ActionId.fromItemId(18588), value: Explosive_Small.EzThroDynamiteII}, stats: []},
	{ config: { actionId: ActionId.fromItemId(6714), value: Explosive_Small.EzThroDynamite}, stats: []},
	{ config: { actionId: ActionId.fromItemId(18641), value: Explosive_Small.DenseDynamite, showWhen: player => player.hasProfession(Profession.Engineering) }, stats: []},
	{ config: { actionId: ActionId.fromItemId(4378), value: Explosive_Small.HeavyDynamite, showWhen: player => player.hasProfession(Profession.Engineering) }, stats: []},
] as ConsumableStatOption<Explosive_Small>[];

export const BIG_EXPLOSIVES_CONFIG = [
	{ config: {actionId: ActionId.fromItemId(42641), value: Explosive_Big.ThermalSapper, showWhen: player => player.hasProfession(Profession.Engineering)}, stats: []},
	{ config: {actionId: ActionId.fromItemId(23827), value: Explosive_Big.SuperSapperCharge, showWhen: player => player.hasProfession(Profession.Engineering)}, stats: []},
	{ config: {actionId: ActionId.fromItemId(10646), value: Explosive_Big.GoblinSapperCharge, showWhen: player => player.hasProfession(Profession.Engineering)}, stats: []},
] as ConsumableStatOption<Explosive_Big>[];

export const DECOY_EXPLOSIVE_CONFIG = [
	{ config: {actionId: ActionId.fromItemId(40536), value: Explosive_Medium.ExplosiveDecoy, showWhen: player => player.hasProfession(Profession.Engineering)}, stats: []},
] as ConsumableStatOption<Explosive_Medium>[];

export const makeExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'explosiveSmall'
});

export const makeBigExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'explosiveBig'
});

export const makeDecoyExplosivesInput = makeConsumeInputFactory({
	consumesFieldName: 'explosiveMedium'
});

///////////////////////////////////////////////////////////////////////////
//                                 FLASKS + ELIXIRS
///////////////////////////////////////////////////////////////////////////

// Flasks

export const makeFlasksInput = makeConsumeInputFactory({
    consumesFieldName: 'flask',
    onSet: (eventID: EventID, player: Player<any>, newValue: Flask) => {
        if (newValue) {
        	const newConsumes = player.getConsumes();
      		newConsumes.battleElixir = BattleElixir.BattleElixirUnknown;
      		newConsumes.guardianElixir = GuardianElixir.GuardianElixirUnknown;
      		player.setConsumes(eventID, newConsumes);
    	}
  	},
	filter: (option,player) => {
		return !option.level || player.getLevel() >= option.level
	},
});

// Battle Elixirs

export const makeBattleElixirsInput = makeConsumeInputFactory({
    consumesFieldName: 'battleElixir',
    onSet: (eventID: EventID, player: Player<any>, newValue: BattleElixir) => {
	    if (newValue) {
            const newConsumes = player.getConsumes();
            newConsumes.flask = Flask.FlaskUnknown;
    	    player.setConsumes(eventID, newConsumes);
		}
  	},
  	filter: (option,player) => {
		return !option.level || player.getLevel() >= option.level
	},
});

// Guardian Elixirs

export const makeGuardianElixirsInput = makeConsumeInputFactory({
	consumesFieldName: 'guardianElixir',
	onSet: (eventID: EventID, player: Player<any>, newValue: GuardianElixir) => {
		if (newValue) {
			const newConsumes = player.getConsumes();
			newConsumes.flask = Flask.FlaskUnknown;
			player.setConsumes(eventID, newConsumes);
		}
	},
	filter: (option,player) => {
		return !option.level || player.getLevel() >= option.level
	},
});

export const IMBUE_CONFIG = [
{ config : {actionId: ActionId.fromItemId(20749), value: WeaponImbue.BrilliantWizardOil}, stats: [Stat.StatSpellCrit, Stat.StatSpellPower],level: 45},
{ config : {actionId: ActionId.fromItemId(20748), value: WeaponImbue.BrilliantManaOil}, stats: [Stat.StatMP5, Stat.StatSpellPower],level: 45},
{ config : {actionId: ActionId.fromItemId(22522), value: WeaponImbue.SuperiorWizardOil}, stats: [Stat.StatSpellPower],level: 58},
{ config : {actionId: ActionId.fromItemId(22521), value: WeaponImbue.SuperiorManaOil}, stats: [Stat.StatMP5],level: 58},
{ config : {actionId: ActionId.fromItemId(20750), value: WeaponImbue.WizardOil}, stats: [Stat.StatSpellPower],level: 40},
{ config : {actionId: ActionId.fromItemId(20747), value: WeaponImbue.LesserManaOil}, stats: [Stat.StatMP5],level: 30},
{ config : {actionId: ActionId.fromItemId(20746), value: WeaponImbue.LesserWizardOil}, stats: [Stat.StatSpellPower],level: 30},
{ config : {actionId: ActionId.fromItemId(20745), value: WeaponImbue.MinorManaOil}, stats: [Stat.StatMP5],level: 5},
{ config : {actionId: ActionId.fromItemId(20744), value: WeaponImbue.MinorWizardOil}, stats: [Stat.StatSpellPower],level: 5},
{ config : {actionId: ActionId.fromItemId(23529), value: WeaponImbue.AdamantiteSharpeningStone}, stats: [Stat.StatMeleeCrit, Stat.StatAttackPower],level: 60, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(28421), value: WeaponImbue.AdamantiteWeightStone}, stats: [Stat.StatMeleeCrit, Stat.StatAttackPower],level: 60, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(23528), value: WeaponImbue.FelSharpeningStone}, stats: [Stat.StatAttackPower],level: 50, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(28420), value: WeaponImbue.FelWeightstone}, stats: [Stat.StatAttackPower],level: 50, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(18262), value: WeaponImbue.ElementalSharpeningStone}, stats: [Stat.StatMeleeCrit],level: 50},
{ config : {actionId: ActionId.fromItemId(12404), value: WeaponImbue.DenseSharpeningStone}, stats: [Stat.StatAttackPower],level: 35, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(12643), value: WeaponImbue.DenseWeightstone}, stats: [Stat.StatAttackPower],level: 35, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(7964), value: WeaponImbue.SolidSharpeningStone}, stats: [Stat.StatAttackPower],level: 25, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(7965), value: WeaponImbue.SolidWeightStone}, stats: [Stat.StatAttackPower],level: 25, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(2871), value: WeaponImbue.HeavySharpeningStone}, stats: [Stat.StatAttackPower],level: 15, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(3241), value: WeaponImbue.HeavyWeightStone,}, stats: [Stat.StatAttackPower],level: 15, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(2863), value: WeaponImbue.CoarseSharpeningStone}, stats: [Stat.StatAttackPower],level: 5, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(3240), value: WeaponImbue.CoarseWeightStone}, stats: [Stat.StatAttackPower],level: 5, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(2862), value: WeaponImbue.RoughSharpeningStone }, stats: [Stat.StatAttackPower],level: 1, type: "sharp"},
{ config : {actionId: ActionId.fromItemId(3239), value: WeaponImbue.RoughWeightStone}, stats: [Stat.StatAttackPower],level: 1, type: "blunt"},
{ config : {actionId: ActionId.fromItemId(34539), value: WeaponImbue.RighteousWeaponCoating}, stats: [Stat.StatAttackPower],level: 70},
{ config : {actionId: ActionId.fromItemId(34538), value: WeaponImbue.BlessedWeaponCoating}, stats: [Stat.StatMP5],level: 70},
{ config : {actionId: ActionId.fromItemId(30696), value: WeaponImbue.ConsecratedWeapon}, stats: [Stat.StatAttackPower], condition: player => { return player.isClass(Class.ClassPaladin) && player.getEquippedItem(ItemSlot.ItemSlotMainHand)!.item.ilvl < 138 && player.hasTrinketEquipped(30696)
}},
] as ImbueConsumableStatOption[];

export const makeMHImbueInput = makeConsumeInputFactory({
	consumesFieldName: 'mhImbue',
	showWhen: player => {
		const mh = player.getGear().getEquippedItem(ItemSlot.ItemSlotMainHand);
		return mh != null && mh?.item.ilvl <= 165;
	},
	filter: (option,player) => {
		const opt = option as ImbueConsumableStatOption
		return (!opt.level || player.getLevel() >= opt.level) &&
		(!opt.type || (opt.type === "sharp" && player.getGear().hasSharpMHWeapon()) || (opt.type === "blunt" && player.getGear().hasBluntMHWeapon())) && (!opt.condition || opt.condition(player))
	},
});
export const makeOHImbueInput = makeConsumeInputFactory({
	consumesFieldName: 'ohImbue',
	showWhen: player => {
		const oh = player.getGear().getEquippedItem(ItemSlot.ItemSlotOffHand);
		return oh != null && oh?.item.ilvl <= 165;
	},
	filter: (option,player) => {
		const opt = option as ImbueConsumableStatOption
		return (!opt.level || player.getLevel() >= opt.level) && (!opt.type || (opt.type === "sharp" && player.getGear().hasSharpOHWeapon()) || (opt.type === "blunt" && player.getGear().hasBluntOHWeapon())) && (!opt.condition || opt.condition(player))
	},
});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const makeFoodInput = makeConsumeInputFactory({
	consumesFieldName: 'food',
	filter: (option,player) => {
		return !option.level || player.getLevel() >= option.level
	},});

///////////////////////////////////////////////////////////////////////////
//                                 PET
///////////////////////////////////////////////////////////////////////////

export const SpicedMammothTreats = makeBooleanConsumeInput({actionId: ActionId.fromItemId(43005), fieldName: 'petFood', value: PetFood.PetFoodSpicedMammothTreats});
export const PetScrollOfAgilityV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27498), fieldName: 'petScrollOfAgility', value: 5});
export const PetScrollOfStrengthV = makeBooleanConsumeInput({actionId: ActionId.fromItemId(27503), fieldName: 'petScrollOfStrength', value: 5});

///////////////////////////////////////////////////////////////////////////
//                                 POTIONS
///////////////////////////////////////////////////////////////////////////

export const RunicHealingPotion   = { actionId: ActionId.fromItemId(33447), value: Potions.RunicHealingPotion };
export const RunicHealingInjector = { actionId: ActionId.fromItemId(41166), value: Potions.RunicHealingInjector };
export const RunicManaPotion      = { actionId: ActionId.fromItemId(33448), value: Potions.RunicManaPotion };
export const RunicManaInjector    = { actionId: ActionId.fromItemId(42545), value: Potions.RunicManaInjector };
export const IndestructiblePotion = { actionId: ActionId.fromItemId(40093), value: Potions.IndestructiblePotion };
export const PotionOfSpeed        = { actionId: ActionId.fromItemId(40211), value: Potions.PotionOfSpeed };
export const PotionOfWildMagic    = { actionId: ActionId.fromItemId(40212), value: Potions.PotionOfWildMagic };

export const DestructionPotion    = { actionId: ActionId.fromItemId(22839), value: Potions.DestructionPotion };
export const HastePotion          = { actionId: ActionId.fromItemId(22838), value: Potions.HastePotion };
export const MightyRagePotion     = { actionId: ActionId.fromItemId(13442), value: Potions.MightyRagePotion };
export const SuperManaPotion      = { actionId: ActionId.fromItemId(22832), value: Potions.SuperManaPotion };
export const FelManaPotion        = { actionId: ActionId.fromItemId(31677), value: Potions.FelManaPotion };
export const InsaneStrengthPotion = { actionId: ActionId.fromItemId(22828), value: Potions.InsaneStrengthPotion };
export const IronshieldPotion     = { actionId: ActionId.fromItemId(22849), value: Potions.IronshieldPotion };
export const HeroicPotion         = { actionId: ActionId.fromItemId(22837), value: Potions.HeroicPotion };

export const POTIONS_CONFIG = [
  { config: RunicHealingPotion,   stats: [Stat.StatStamina],level:70 },
  { config: RunicHealingInjector, stats: [Stat.StatStamina],level:70 },
  { config: RunicManaPotion,      stats: [Stat.StatIntellect],level:70 },
  { config: RunicManaInjector,    stats: [Stat.StatIntellect],level:70 },
  { config: IndestructiblePotion, stats: [Stat.StatArmor],level:70 },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength],level:70 },
  { config: HeroicPotion,         stats: [Stat.StatStamina],level:70 },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],level:70 },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower],level:70 },

  { config: DestructionPotion,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower],level:60 },
  { config: HastePotion,          stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],level:60 },
  { config: IronshieldPotion,     stats: [Stat.StatArmor],level:60 },
  { config: SuperManaPotion,      stats: [Stat.StatIntellect],level:60 },
  { config: FelManaPotion,        stats: [Stat.StatRangedAttackPower],level:60 },

] as ConsumableStatOption<Potions>[];

export const PRE_POTIONS_CONFIG = [
  { config: IndestructiblePotion, stats: [Stat.StatArmor],level:60 },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength],level:60 },
  { config: HeroicPotion,         stats: [Stat.StatStamina],level:60 },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],level:60 },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower],level:60 },

  { config: DestructionPotion,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower],level:60 },
  { config: HastePotion,          stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste],level:60 },
  { config: IronshieldPotion,     stats: [Stat.StatArmor],level:60 },
] as ConsumableStatOption<Potions>[];

export const makePotionsInput = makeConsumeInputFactory({consumesFieldName: 'defaultPotion'});
export const makePrepopPotionsInput = makeConsumeInputFactory({consumesFieldName: 'prepopPotion'});

