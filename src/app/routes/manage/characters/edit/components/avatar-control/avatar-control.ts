import {Component, computed, forwardRef, inject, Signal, signal, WritableSignal} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {AvatarImageCrop} from './avatar-image-crop';
import {BooleanSignal, booleanSignal} from '@util/ng';

@Component({
  selector: 'app-avatar-control',
  imports: [
    AvatarImageCrop
  ],
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
  private notifications = inject(Notifications)

  private onChange: (value: Nullable<string>) => void = () => null
  private onTouched: () => void = () => null

  readonly selectedFileForCrop: WritableSignal<Blob | null> = signal(null)
  readonly currentValue: WritableSignal<Nullable<string>> = signal(null)
  readonly isDisabled: BooleanSignal = booleanSignal(false)
  readonly isSet: Signal<boolean> = computed(() => this.currentValue() != null)

  constructor() {
  }

  writeValue(obj: Nullable<string>): void {
    this.currentValue.set(obj)
  }

  registerOnChange(fn: (value: Nullable<string>) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn
  }

  setDisabledState?(isDisabled: boolean): void {
    this.isDisabled.set(isDisabled)
  }

  onFileSelected(e: Event) {
    e.preventDefault();

    const input = e.target as HTMLInputElement;
    const file = input.files?.item(0)

    if (!file) return;

    if (!file?.type?.startsWith('image/')) {
      this.notifications.toast("Please select an image file.", "DANGER")
      return
    }

    this.selectedFileForCrop.set(file)
    input.value = ''
    this.onTouched()
  }

  onCropResult(file: Blob) {
    const reader = new FileReader();

    reader.onloadend = () => {
      this.currentValue.set(reader.result as string)
      this.selectedFileForCrop.set(null)
      this.onChange(this.currentValue())
    }

    reader.onerror = () => {
      this.notifications.toast("An unexpected error occurred processing the cropped image.", "DANGER")
      this.selectedFileForCrop.set(null)
    }

    reader.readAsDataURL(file);
  }

  onCropCanceled() {
    this.selectedFileForCrop.set(null)
  }

  onClear(e: Event) {
    e.stopPropagation();
    this.currentValue.set(null)
    this.onChange(this.currentValue())
  }
}
