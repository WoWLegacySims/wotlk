import { Player } from "../../player";
import {
	Class,
	Conjured,
	Consumes,
	Explosive,
	PetFood,
	Potions,
	Spec,
	Stat,
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
	config: ConsumableInputConfig<T>
}

export interface ConsumeInputFactoryArgs<T extends number> {
	consumesFieldName: keyof Consumes,
	// Additional callback if logic besides syncing consumes is required
	onSet?: (eventactionId: EventID, player: Player<any>, newValue: T) => void
	showWhen?: (player: Player<any>) => boolean
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
						(option.config.faction || player.getFaction()) == player.getFaction()
				} as IconEnumValueConfig<Player<any>, T>;
				if (option.config.value) rtn.value = option.config.value

				return rtn
			})),
			equals: (a: T, b: T) => a == b,
			zeroValue: 0 as T,
			changedEvent: (player: Player<any>) => TypedEvent.onAny([player.consumesChangeEmitter, player.gearChangeEmitter]),
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
export const ConjuredFlameCap = { actionId: ActionId.fromItemId(22788), value: Conjured.ConjuredFlameCap };
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

export const ExplosiveSaroniteBomb    = { actionId: ActionId.fromItemId(41119), value: Explosive.ExplosiveSaroniteBomb };
export const ExplosiveCobaltFragBomb  = { actionId: ActionId.fromItemId(40771), value: Explosive.ExplosiveCobaltFragBomb };

export const EXPLOSIVES_CONFIG = [
	{ config: ExplosiveSaroniteBomb, stats: [] ,level:1},
	{ config: ExplosiveCobaltFragBomb, stats: [],level:1 },
] as ConsumableStatOption<Explosive>[];

export const makeExplosivesInput = makeConsumeInputFactory({consumesFieldName: 'fillerExplosive'});

export const ThermalSapper = makeBooleanConsumeInput({actionId: ActionId.fromItemId(42641), fieldName: 'thermalSapper'});
export const ExplosiveDecoy = makeBooleanConsumeInput({actionId: ActionId.fromItemId(40536), fieldName: 'explosiveDecoy'});

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
  }
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
  }
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
	}
});

///////////////////////////////////////////////////////////////////////////
//                                 FOOD
///////////////////////////////////////////////////////////////////////////

export const makeFoodInput = makeConsumeInputFactory({consumesFieldName: 'food'});

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
  { config: RunicHealingPotion,   stats: [Stat.StatStamina] },
  { config: RunicHealingInjector, stats: [Stat.StatStamina] },
  { config: RunicManaPotion,      stats: [Stat.StatIntellect] },
  { config: RunicManaInjector,    stats: [Stat.StatIntellect] },
  { config: IndestructiblePotion, stats: [Stat.StatArmor] },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength] },
  { config: HeroicPotion,         stats: [Stat.StatStamina] },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
] as ConsumableStatOption<Potions>[];

export const PRE_POTIONS_CONFIG = [
  { config: IndestructiblePotion, stats: [Stat.StatArmor] },
  { config: InsaneStrengthPotion, stats: [Stat.StatStrength] },
  { config: HeroicPotion,         stats: [Stat.StatStamina] },
  { config: PotionOfSpeed,        stats: [Stat.StatMeleeHaste, Stat.StatSpellHaste] },
  { config: PotionOfWildMagic,    stats: [Stat.StatMeleeCrit, Stat.StatSpellCrit, Stat.StatSpellPower] },
] as ConsumableStatOption<Potions>[];

export const makePotionsInput = makeConsumeInputFactory({consumesFieldName: 'defaultPotion'});
export const makePrepopPotionsInput = makeConsumeInputFactory({consumesFieldName: 'prepopPotion'});

