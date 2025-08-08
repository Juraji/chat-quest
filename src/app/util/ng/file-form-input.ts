import {Directive, ElementRef, forwardRef, HostListener} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';

@Directive({
  selector: 'input[appFileFormInput][type=file][formControlName]',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => FileFormInput),
      multi: true,
    },
  ],
})
export class FileFormInput implements ControlValueAccessor {
  constructor(private el: ElementRef<HTMLInputElement>) {
  }

  private onChange: (file: File | null) => void = () => null
  private onTouched: () => void = () => null

  writeValue(file: File | null): void {
    // Clear file input when setting value to null
    if (file === null && this.el.nativeElement.files) {
      this.el.nativeElement.value = '';
    }
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean) {
    this.el.nativeElement.disabled = isDisabled;
  }

  @HostListener('change')
  onFileChange() {
    this.onTouched();
    const files = this.el.nativeElement.files;
    if (files && files[0]) {
      this.onChange(files[0]);
    } else {
      this.onChange(null);
    }
  }
}
