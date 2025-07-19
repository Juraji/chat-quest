import {Component, effect, forwardRef, signal, WritableSignal} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';

type Avatar = Blob | null

@Component({
  selector: 'app-avatar-control',
  imports: [],
  templateUrl: './avatar-control.html',
  styleUrl: './avatar-control.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => AvatarControl),
      multi: true
    }
  ],
  host: {
    '[class.disabled]': 'isDisabled()'
  }
})
export class AvatarControl implements ControlValueAccessor {
  private onChange: (value: Avatar) => void = () => null
  private onTouched: () => void = () => null

  readonly currentValue: WritableSignal<Avatar> = signal(null)
  readonly isDisabled: WritableSignal<boolean> = signal(false)

  readonly imageUrl = signal('')

  constructor() {
    effect(() => {
      const blob = this.currentValue()
      this.imageUrl.update(current => {
        if (!!current) URL.revokeObjectURL(current)
        if (!!blob) return URL.createObjectURL(blob)
        else return ''
      })
    });
  }

  writeValue(obj: Avatar): void {
    this.currentValue.set(obj)
  }

  registerOnChange(fn: (value: Avatar) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn
  }

  setDisabledState?(isDisabled: boolean): void {
    this.isDisabled.set(isDisabled)
  }

  onSetFile(e: Event) {
    e.preventDefault();

    const input = e.target as HTMLInputElement;
    const file = input.files?.item(0)
    this.currentValue.set(!!file ? file : null)

    this.onTouched()
    this.onChange(this.currentValue())
    input.value = ''
  }
}
