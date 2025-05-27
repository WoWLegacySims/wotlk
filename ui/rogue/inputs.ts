import * as InputHelpers from '../core/components/input_helpers.js';
import { DEADLYPOISON, INSTANTPOISON, WOUNDPOISON } from '../core/constants/auras.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Rogue_Options_PoisonImbue as Poison,
} from '../core/proto/rogue.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MainHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'mhImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Main Hand Poison' },
		{ actionId: ActionIDMap.fromItemId(DEADLYPOISON), value: Poison.DeadlyPoison, showWhen: player => player.getLevel() >= 30 },
		{ actionId: ActionIDMap.fromItemId(INSTANTPOISON), value: Poison.InstantPoison, showWhen: player => player.getLevel() >= 20 },
		{ actionId: ActionIDMap.fromItemId(WOUNDPOISON), value: Poison.WoundPoison, showWhen: player => player.getLevel() >= 32 },
	],
});

export const OffHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'ohImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Off Hand Poison' },
		{ actionId: ActionIDMap.fromItemId(DEADLYPOISON), value: Poison.DeadlyPoison, showWhen: player => player.getLevel() >= 30 },
		{ actionId: ActionIDMap.fromItemId(INSTANTPOISON), value: Poison.InstantPoison, showWhen: player => player.getLevel() >= 20 },
		{ actionId: ActionIDMap.fromItemId(WOUNDPOISON), value: Poison.WoundPoison, showWhen: player => player.getLevel() >= 32 },
	],
});

export const StartingOverkillDuration = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'startingOverkillDuration',
	label: 'Starting Overkill duration',
	labelTooltip: 'Initial Overkill buff duration at the start of each iteration.',
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0
});

export const VanishBreakTime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'vanishBreakTime',
	label: 'Vanish Break Time',
	labelTooltip: 'Time it takes to start attacking after casting Vanish.',
	extraCssClasses: ['experimental'],
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0
})

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: 'Assume bleed always exists for \'Hunger for Blood\' activation. Otherwise will only calculate based on own garrote/rupture.',
	extraCssClasses: ['within-raid-sim-hide'],
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().hungerForBlood
})

export const HonorOfThievesCritRate = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'honorOfThievesCritRate',
	label: 'Honor of Thieves Crit Rate',
	labelTooltip: 'Number of crits other group members generate within 100 seconds',
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().honorAmongThieves > 0
});

export const ApplyPoisonsManually = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'applyPoisonsManually',
	label: 'Configure poisons manually',
	labelTooltip: 'Prevent automatic poison configuration that is based on equipped weapons.',
});
