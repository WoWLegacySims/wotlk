import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';
import {
	WarriorShout,
} from '../core/proto/warrior.js';
import { ActionId } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecProtectionWarrior>({
	fieldName: 'startingRage',
	label: 'Starting Rage',
	labelTooltip: 'Initial rage at the start of each iteration.',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecProtectionWarrior, WarriorShout>({
	fieldName: 'shout',
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionId.fromSpellId(47436), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
	],
});

export const ShatteringThrow = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecProtectionWarrior>({
	fieldName: 'useShatteringThrow',
	id: ActionId.fromSpellId(64382),
});

export const Munch = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecProtectionWarrior>({
 	fieldName: 'munch',
	label: 'Munching',
	labelTooltip: 'When two crits occur at the same time (20 ms window), only the latter will count towards deep wounds',
});
