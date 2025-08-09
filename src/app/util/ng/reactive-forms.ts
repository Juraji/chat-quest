import {AbstractControl, AsyncValidatorFn, FormArray, FormControl, FormGroup, ValidatorFn} from '@angular/forms';
import {Observable} from 'rxjs';
import {Signal} from '@angular/core';
import {toSignal} from '@angular/core/rxjs-interop';

export type FormControls<T> = {
  [P in keyof T]: AbstractControl<T[P]>
}

// NG Reactive Forms typing is broken af, this fixes it, mostly.
export type TypedFormGroup<T> =
  FormGroup<FormControls<T>> & { value: T, valueChanges: Observable<T> }
export type TypedFormControl<T> =
  FormControl<T> & { value: T, valueChanges: Observable<T> }
export type TypedFormArray<T> =
  FormArray<AbstractControl<T>>

export function formGroup<T extends object>(
  controls: FormControls<T>,
  validators?: ValidatorFn | ValidatorFn[]
): TypedFormGroup<T> {
  return new FormGroup(controls, {validators, updateOn: 'blur'}) as TypedFormGroup<T>
}

export function formArray<T>(
  controls: AbstractControl<T>[],
  validators?: ValidatorFn | ValidatorFn[]
): TypedFormArray<T> {
  return new FormArray<AbstractControl<T>>(controls, {validators, updateOn: 'blur'}) as TypedFormArray<T>
}

export function formControl<T>(
  value: T,
  validators?: ValidatorFn | ValidatorFn[],
  asyncValidators?: AsyncValidatorFn | AsyncValidatorFn[],
): TypedFormControl<T> {
  return new FormControl(
    {value, disabled: false},
    {nonNullable: true, updateOn: 'change', validators, asyncValidators}
  ) as TypedFormControl<T>
}

export function readOnlyControl<T>(value: T = "" as T): AbstractControl<T, T> {
  return new FormControl<T>({value, disabled: true}, {nonNullable: true})
}

export function controlValueSignal<T>(control: AbstractControl<T>): Signal<T>
export function controlValueSignal<T>(control: AbstractControl, path: string | (string | number)[]): Signal<T>
export function controlValueSignal<T>(control: AbstractControl, path?: string | (string | number)[]): Signal<T> {
  if(!!path) {
    const ctrl = control.get(path);
    if (!ctrl) throw new Error(`Could not find control for path: ${path}`);
    return toSignal(ctrl.valueChanges, {initialValue: ctrl.value})
  } else {
    return toSignal(control.valueChanges, {initialValue: control.value})
  }
}
