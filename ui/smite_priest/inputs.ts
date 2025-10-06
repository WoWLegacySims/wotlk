import * as InputHelpers from '../core/components/input_helpers.js';
import { INNERFIRE } from '../core/constants/auras.js';
import { Player } from '../core/player.js';
import { Spec,UnitReference, UnitReference_Type as UnitType  } from '../core/proto/common.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';
import { EventID } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfPowerInfusion = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'powerInfusionTarget',
	id: ActionId.fromSpellId(10060),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecSmitePriest>) => player.getSpecOptions().powerInfusionTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecSmitePriest>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.powerInfusionTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const InnerFire = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'useInnerFire',
	id: ActionIDMap.fromSpellId(INNERFIRE),
	showWhen: player => player.getLevel() >= 12
});

export const Shadowfiend = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'useShadowfiend',
	id: ActionId.fromSpellId(34433),
});
