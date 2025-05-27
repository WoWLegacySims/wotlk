import * as InputHelpers from '../core/components/input_helpers.js';
import { INNERFIRE } from '../core/constants/auras.js';
import { Player } from '../core/player.js';
import { Spec,UnitReference, UnitReference_Type as UnitType  } from '../core/proto/common.js';
import { ActionId, ActionIDMap } from '../core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfPowerInfusion = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'powerInfusionTarget',
	id: ActionId.fromSpellId(10060),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecHealingPriest>) => player.getSpecOptions().powerInfusionTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecHealingPriest>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.powerInfusionTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const InnerFire = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'useInnerFire',
	id: ActionIDMap.fromSpellId(INNERFIRE),
	showWhen: player => player.getLevel() >= 12
});

export const Shadowfiend = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHealingPriest>({
	fieldName: 'useShadowfiend',
	id: ActionId.fromSpellId(34433),
});

export const RapturesPerMinute = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHealingPriest>({
	fieldName: 'rapturesPerMinute',
	label: 'Raptures / Min',
	labelTooltip: 'Number of times to proc Rapture each minute (due to a PWS being fully absorbed).',
	showWhen: (player: Player<Spec.SpecHealingPriest>) => player.getTalents().rapture > 0,
	changeEmitter: (player: Player<Spec.SpecHealingPriest>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});
