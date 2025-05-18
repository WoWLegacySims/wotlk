import { EventID, TypedEvent } from '../typed_event.js';
import { Input, InputConfig } from './input.js';

/**
 * Data for creating a number picker.
 */
export interface NumberPickerConfig<ModObject> extends InputConfig<ModObject, number> {
	float?: boolean,
	positive?: boolean,
	max?: () => number,
	min?: () => number,
}

// UI element for picking an arbitrary number field.
export class NumberPicker<ModObject> extends Input<ModObject, number> {
	private readonly inputElem: HTMLInputElement;
	private float: boolean;
	private positive: boolean;
	private max?: () => number;
	private min?: () => number;

	constructor(parent: HTMLElement | null, modObject: ModObject, config: NumberPickerConfig<ModObject>) {
		super(parent, 'number-picker-root', modObject, config);
		this.float = config.float || false;
		this.positive = config.positive || false;
		this.min = config.min;
		this.max = config.max;

		this.inputElem = document.createElement('input');
		this.inputElem.type = 'text';
		this.inputElem.classList.add('form-control', 'number-picker-input');

		this.inputElem.onchange = e => {
			let value: number
			if (this.float) {
				value = parseFloat(this.inputElem.value)
			} else {
				value = parseInt(this.inputElem.value)
			}

			if(this.positive) {
				value = Math.abs(value)
			}

			if(this.min) {
				value = Math.max(value,this.min())
			}

			if(this.max) {
				value = Math.min(value,this.max())
			}

			if (this.float) {
				this.inputElem.value = value.toFixed(2)
			} else {
				this.inputElem.value = value.toString()
			}
		}

		this.rootElem.appendChild(this.inputElem);
		this.init();

		this.inputElem.addEventListener('change', event => {
			this.inputChanged(TypedEvent.nextEventID());
		});

		this.inputElem.addEventListener('input', event => {
			this.updateSize();
		});
		this.updateSize();
	}

	getInputElem(): HTMLElement {
		return this.inputElem;
	}

	getInputValue(): number {
		if (this.float) {
			return parseFloat(this.inputElem.value || '') || 0;
		} else {
			return parseInt(this.inputElem.value || '') || 0;
		}
	}

	setInputValue(newValue: number) {
		if (this.float)
			this.inputElem.value = newValue.toFixed(2);
		else
			this.inputElem.value = String(newValue);
	}

	private updateSize() {
		const newSize = Math.max(3, this.inputElem.value.length);
		if (this.inputElem.size != newSize)
			this.inputElem.size = newSize;
	}
}
