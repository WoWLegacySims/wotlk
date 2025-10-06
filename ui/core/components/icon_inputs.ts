import { IndividualSimUI } from '../individual_sim_ui.js';
import { Party } from '../party.js';
import { Player } from '../player';
import {
	Consumes,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	RaidBuffs,
	Spec,
} from '../proto/common.js';
import { ActionId, ActionIDMap } from '../proto_utils/action_id.js';
import { Raid } from '../raid';
import { EventID, TypedEvent } from '../typed_event';
import { IconEnumPicker } from './icon_enum_picker';
import { IconPicker } from './icon_picker';
import * as InputHelpers from './input_helpers';

// Component Functions

export type IconInputConfig<ModObject, T> = (
	InputHelpers.TypedIconPickerConfig<ModObject, T> |
	InputHelpers.TypedIconEnumPickerConfig<ModObject, T>
);

export const buildIconInput = (parent: HTMLElement, simUI: IndividualSimUI<Spec>, inputConfig: IconInputConfig<Player<Spec>, any>) => {
	parent.classList.remove('hide')
	if (inputConfig.type == 'icon') {
		return new IconPicker<Player<Spec>, any>(parent, simUI.player, inputConfig, simUI);
	} else if (inputConfig.type == 'iconEnum') {
		return new IconEnumPicker<Player<Spec>, any>(parent, simUI.player, inputConfig, simUI);
	}
};

export function withLabel<ModObject, T>(config: IconInputConfig<ModObject, T>, label: string): IconInputConfig<ModObject, T> {
	config.label = label;
	return config;
}

interface BooleanInputConfig<T> {
	actionId: ActionId | ActionIDMap
	fieldName: keyof T
	value?: number
	faction?: Faction
	showWhen?: (player: Player<Spec>) => boolean
}

export function makeBooleanRaidBuffInput<SpecType extends Spec>(config: BooleanInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter]),
	}, config.actionId, config.fieldName, config.value);
}
export function makeBooleanPartyBuffInput<SpecType extends Spec>(config: BooleanInputConfig<PartyBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, PartyBuffs, Party>({
		getModObject: (player: Player<SpecType>) => player.getParty()!,
		getValue: (party: Party) => party.getBuffs(),
		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
		changeEmitter: (party: Party) => TypedEvent.onAny([party.buffsChangeEmitter,party.getPlayer(0)!.levelChangeEmitter]),
		showWhen: (party: Party) => !config.showWhen || config.showWhen(party.getPlayer(0)!)
	}, config.actionId, config.fieldName, config.value);
}

export function makeBooleanIndividualBuffInput<SpecType extends Spec>(config: BooleanInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter]),
	}, config.actionId, config.fieldName, config.value);
}

export function makeBooleanConsumeInput<SpecType extends Spec>(config: BooleanInputConfig<Consumes>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Consumes, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		getValue: (player: Player<SpecType>) => player.getConsumes(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: Consumes) => player.setConsumes(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.consumesChangeEmitter])
	}, config.actionId, config.fieldName, config.value);
}
export function makeBooleanDebuffInput<SpecType extends Spec>(config: BooleanInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, boolean> {
	return InputHelpers.makeBooleanIconInput<any, Debuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		getValue: (player: Player<SpecType>) => player.getRaid()!.getDebuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: Debuffs) => player.getRaid()!.setDebuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.debuffsChangeEmitter,player.levelChangeEmitter]),
		showWhen: config.showWhen,
	}, config.actionId, config.fieldName, config.value);
}

interface TristateInputConfig<T> {
	actionId: ActionId | ActionIDMap
	impId: ActionId
	fieldName: keyof T
	faction?: Faction
	showWhen?: (player: Player<Spec>) => boolean
}

export function makeTristateRaidBuffInput<SpecType extends Spec>(config: TristateInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter]),
	}, config.actionId, config.impId, config.fieldName);
}

export function makeTristateIndividualBuffInput<SpecType extends Spec>(config: TristateInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter])
	}, config.actionId, config.impId, config.fieldName);
}

export function makeTristateDebuffInput<SpecType extends Spec>(config: TristateInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeTristateIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<SpecType>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => TypedEvent.onAny([raid.debuffsChangeEmitter,raid.getPlayer(0)!.levelChangeEmitter]),
		showWhen: (raid: Raid) => !config.showWhen || config.showWhen(raid.getPlayer(0)!)
	}, config.actionId, config.impId, config.fieldName);
}

interface QuadStateInputConfig<T> {
	actionId: ActionId | ActionIDMap
	impId: ActionId
	impId2: ActionId
	fieldName: keyof T
	faction?: Faction
	showWhen: (player: Player<Spec>) => boolean
}

export function makeQuadstateDebuffInput<SpecType extends Spec>(config: QuadStateInputConfig<Debuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeQuadstateIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<SpecType>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => TypedEvent.onAny([raid.debuffsChangeEmitter, raid.getPlayer(0)!.levelChangeEmitter]),
		showWhen: (raid: Raid) => config.showWhen(raid.getPlayer(0)!),
	}, config.actionId, config.impId, config.impId2, config.fieldName);
}

interface MultiStateInputConfig<T> {
	actionId: ActionId
	numStates: number
	fieldName: keyof T
	multiplier?: number
	faction?: Faction
	showWhen?: (player: Player<Spec>) => boolean
}

export function makeMultistateRaidBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<RaidBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, RaidBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getRaid()!.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: RaidBuffs) => player.getRaid()!.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.getRaid()!.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter]),
	}, config.actionId, config.numStates, config.fieldName, config.multiplier);
}
export function makeMultistatePartyBuffInput<SpecType extends Spec>(actionId: ActionId, numStates: number, fieldName: keyof PartyBuffs, showWhen?: (party: Party) => boolean): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, PartyBuffs, Party>({
		getModObject: (player: Player<SpecType>) => player.getParty()!,
		getValue: (party: Party) => party.getBuffs(),
		setValue: (eventID: EventID, party: Party, newVal: PartyBuffs) => party.setBuffs(eventID, newVal),
		changeEmitter: (party: Party) => TypedEvent.onAny([party.buffsChangeEmitter, party.getPlayer(0)!.levelChangeEmitter]),
		showWhen: showWhen
	}, actionId, numStates, fieldName);
}
export function makeMultistateIndividualBuffInput<SpecType extends Spec>(config: MultiStateInputConfig<IndividualBuffs>): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		showWhen: (player: Player<SpecType>) =>
			(!config.faction || config.faction == player.getFaction()) &&
			(!config.showWhen || config.showWhen(player)),
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.raceChangeEmitter,player.levelChangeEmitter]),
	}, config.actionId, config.numStates, config.fieldName, config.multiplier);
}

export function makeMultistateMultiplierIndividualBuffInput<SpecType extends Spec>(actionId: ActionId, numStates: number, multiplier: number, fieldName: keyof IndividualBuffs, showWhen?: (player: Player<SpecType>) => boolean): InputHelpers.TypedIconPickerConfig<Player<SpecType>, number> {
	return InputHelpers.makeMultistateIconInput<any, IndividualBuffs, Player<SpecType>>({
		getModObject: (player: Player<SpecType>) => player,
		getValue: (player: Player<SpecType>) => player.getBuffs(),
		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter,player.levelChangeEmitter]),
		showWhen: showWhen
	}, actionId, numStates, fieldName, multiplier);
}

export function makeMultistateMultiplierDebuffInput<SpecType extends Spec>(actionId: ActionId, numStates: number, multiplier: number, fieldName: keyof Debuffs): InputHelpers.TypedIconPickerConfig<Player<any>, number> {
	return InputHelpers.makeMultistateIconInput<any, Debuffs, Raid>({
		getModObject: (player: Player<SpecType>) => player.getRaid()!,
		getValue: (raid: Raid) => raid.getDebuffs(),
		setValue: (eventID: EventID, raid: Raid, newVal: Debuffs) => raid.setDebuffs(eventID, newVal),
		changeEmitter: (raid: Raid) => raid.debuffsChangeEmitter,
	}, actionId, numStates, fieldName, multiplier);
}

// interface EnumInputConfig<ModObject, Message, T> {
// 	fieldName: keyof Message
// 	values: Array<IconEnumValueConfig<ModObject, T>>
// 	direction?: IconEnumPickerDirection
// 	numColumns?: number
// 	faction?: Faction
// }

// export function makeEnumIndividualBuffInput<SpecType extends Spec>(config: EnumInputConfig<Player<SpecType>, IndividualBuffs, number>): InputHelpers.TypedIconEnumPickerConfig<Player<SpecType>, number> {
// 	return InputHelpers.makeEnumIconInput<any, IndividualBuffs, Player<SpecType>, number>({
// 		getModObject: (player: Player<SpecType>) => player,
// 		showWhen: (player: Player<SpecType>) =>
// 			(!config.faction || config.faction == player.getFaction()),
// 		getValue: (player: Player<SpecType>) => player.getBuffs(),
// 		setValue: (eventID: EventID, player: Player<SpecType>, newVal: IndividualBuffs) => player.setBuffs(eventID, newVal),
// 		changeEmitter: (player: Player<SpecType>) => TypedEvent.onAny([player.buffsChangeEmitter, player.raceChangeEmitter]),
// 	}, config.fieldName, config.values, config.numColumns, config.direction || IconEnumPickerDirection.Vertical)
// };
