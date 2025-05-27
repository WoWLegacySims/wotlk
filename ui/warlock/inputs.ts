import * as InputHelpers from '../core/components/input_helpers.js';
import { CREATEFIRESTONE, CREATESPELLSTONE, DEMONARMOR, FELARMOR } from '../core/constants/auras.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options_WeaponImbue as WeaponImbue,
} from '../core/proto/warlock.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, Armor>({
	fieldName: 'armor',
	values: [
		{ value: Armor.NoArmor, tooltip: 'No Armor' },
		{ actionId: ActionIDMap.fromSpellId(FELARMOR), value: Armor.FelArmor, showWhen: player => player.getLevel() >= 62 },
		{ actionId: ActionIDMap.fromSpellId(DEMONARMOR), value: Armor.DemonArmor },
	],
});

export const WeaponImbueInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, WeaponImbue>({
	fieldName: 'weaponImbue',
	values: [
		{ value: WeaponImbue.NoWeaponImbue, tooltip: 'No Weapon Stone' },
		{ actionId: ActionIDMap.fromSpellId(CREATEFIRESTONE), value: WeaponImbue.GrandFirestone, showWhen: player => player.getLevel() >= 28 },
		{ actionId: ActionIDMap.fromSpellId(CREATESPELLSTONE), value: WeaponImbue.GrandSpellstone, showWhen: player => player.getLevel() >= 36 },
	],
});

export const PetInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarlock, Summon>({
	fieldName: 'summon',
	values: [
		{ value: Summon.NoSummon, tooltip: 'No Pet' },
		{ actionId: ActionId.fromSpellId(688), value: Summon.Imp },
		{ actionId: ActionId.fromSpellId(712), value: Summon.Succubus },
		{ actionId: ActionId.fromSpellId(691), value: Summon.Felhunter },
		{
			actionId: ActionId.fromSpellId(30146), value: Summon.Felguard,
			showWhen: (player: Player<Spec.SpecWarlock>) => player.getTalents().summonFelguard,
		},
	],
	changeEmitter: (player: Player<Spec.SpecWarlock>) => player.changeEmitter,
});

export const DetonateSeed = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarlock>({
	fieldName: 'detonateSeed',
	label: 'Detonate Seed on Cast',
	labelTooltip: 'Simulates raid doing damage to targets such that seed detonates immediately on cast.',
});
