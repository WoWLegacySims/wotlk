import * as InputHelpers from '../core/components/input_helpers.js';
import { BATTLESHOUT, COMMANDINGSHOUT } from '../core/constants/auras.js';
import { Spec } from '../core/proto/common.js';
import {
	WarriorShout,
} from '../core/proto/warrior.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useRecklessness',
	id: ActionId.fromSpellId(1719),
});

export const ShatteringThrow = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useShatteringThrow',
	id: ActionId.fromSpellId(64382),
});

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const Munch = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'munch',
	label: 'Munching',
	labelTooltip: 'When two crits occur at the same time (20 ms window), only the latter will count towards deep wounds',
});

export const StanceSnapshot = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'stanceSnapshot',
	label: 'Stance Snapshot',
	labelTooltip: 'Ability that is cast at the same time as stance swap will benefit from the bonus of the stance before the swap.',
});

// Allows for auto gemming whilst ignoring expertise cap
// (Useful for Arms)
export const DisableExpertiseGemming = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'disableExpertiseGemming',
	label: 'Disable expertise gemming',
	labelTooltip: 'Disables auto gemming for expertise',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarrior, WarriorShout>({
	fieldName: 'shout',
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionIDMap.fromSpellId(BATTLESHOUT), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionIDMap.fromSpellId(COMMANDINGSHOUT), value: WarriorShout.WarriorShoutCommanding, showWhen: player => player.getLevel() >= 68 },
	],
});
