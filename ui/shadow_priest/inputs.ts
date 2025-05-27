import * as InputHelpers from '../core/components/input_helpers.js';
import { INNERFIRE } from '../core/constants/auras.js';
import { Spec } from '../core/proto/common.js';
import {
	ShadowPriest_Options_Armor as Armor,
} from '../core/proto/priest.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ArmorInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecShadowPriest, Armor>({
	fieldName: 'armor',
	values: [
		{ value: Armor.NoArmor, tooltip: 'No Inner Fire' },
		{ actionId: ActionIDMap.fromSpellId(INNERFIRE), value: Armor.InnerFire, showWhen: player => player.getLevel() >= 12 },
	],
});
