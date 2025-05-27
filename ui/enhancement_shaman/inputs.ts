import * as InputHelpers from '../core/components/input_helpers.js';
import { FLAMETONGUEWEAPON, FLAMETONGUEWEAPONDR, FROSTBRANDWEAPON, LIGHTNINGSHIELD, ROCKBITERWEAPON, WATERSHIELD, WINDFURYWEAPON } from '../core/constants/auras.js';
import { ItemSlot, Spec } from '../core/proto/common.js';
import {
	ShamanImbue,
	ShamanShield,
	ShamanSyncType,
} from '../core/proto/shaman.js';
import { ActionIDMap } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionIDMap.fromSpellId(WATERSHIELD), value: ShamanShield.WaterShield, showWhen: player => player.getLevel() >= 20},
		{ actionId: ActionIDMap.fromSpellId(LIGHTNINGSHIELD), value: ShamanShield.LightningShield, showWhen: player => player.getLevel() >= 8},
	],
});

export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueMh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Main Hand Enchant' },
		{ actionId: ActionIDMap.fromSpellId(WINDFURYWEAPON), value: ShamanImbue.WindfuryWeapon, showWhen: player => player.getLevel() >= 30},
		{ actionId: ActionIDMap.fromSpellId(FLAMETONGUEWEAPON), value: ShamanImbue.FlametongueWeapon, text: 'Max', showWhen: player => player.getLevel() >= 10},
		{ actionId: ActionIDMap.fromSpellId(FLAMETONGUEWEAPONDR), value: ShamanImbue.FlametongueWeaponDownrank, text: 'Down', showWhen: player => player.getLevel() >= 40 },
		{ actionId: ActionIDMap.fromSpellId(FROSTBRANDWEAPON), value: ShamanImbue.FrostbrandWeapon, showWhen: player => player.getLevel() >= 20 },
		{ actionId: ActionIDMap.fromSpellId(ROCKBITERWEAPON), value: ShamanImbue.RockbiterWeapon, showWhen: player => player.getLevel() < 30 },
	],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueOh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Off Hand Enchant' },
		{ actionId: ActionIDMap.fromSpellId(WINDFURYWEAPON), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionIDMap.fromSpellId(FLAMETONGUEWEAPON), value: ShamanImbue.FlametongueWeapon, text: 'Max' },
		{ actionId: ActionIDMap.fromSpellId(FLAMETONGUEWEAPONDR), value: ShamanImbue.FlametongueWeaponDownrank, text: 'Down' },
		{ actionId: ActionIDMap.fromSpellId(FROSTBRANDWEAPON), value: ShamanImbue.FrostbrandWeapon },
	],
	showWhen: player => player.getEquippedItem(ItemSlot.ItemSlotOffHand) != null,
	changeEmitter: player => TypedEvent.onAny([player.gearChangeEmitter,player.levelChangeEmitter,player.specOptionsChangeEmitter]),
});

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman, ShamanSyncType>({
	fieldName: 'syncType',
	label: 'Sync/Stagger Setting',
	labelTooltip:
		`Choose your sync or stagger option Perfect
		<ul>
			<li><div>Auto: Will auto pick sync options based on your weapons attack speeds</div></li>
			<li><div>None: No Sync or Staggering, used for mismatched weapon speeds</div></li>
			<li><div>Perfect Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>
			<li><div>Delayed Offhand: Adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window</div></li>
		</ul>`,
	values: [
		{ name: "Automatic", value: ShamanSyncType.Auto },
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});
