import {Injectable} from '@angular/core';
import {CharacterFormData} from './character-form-data';
import {
  formArray,
  formControl,
  formGroup,
  readOnlyControl,
  TypedFormArray,
  TypedFormControl,
  TypedFormGroup
} from '@util/ng';
import {FormArray, Validators} from '@angular/forms';
import {Observable, Subject} from 'rxjs';
import {Character, Tag} from '@api/characters';

@Injectable()
export class CharacterEditFormService {

  readonly formGroup = formGroup<CharacterFormData>({
    character: formGroup({
      id: readOnlyControl(0),
      createdAt: readOnlyControl<Nullable<string>>(null),
      name: formControl('', [Validators.required]),
      favorite: formControl(false),
      avatarUrl: formControl<Nullable<string>>(null),
      appearance: formControl<Nullable<string>>(null),
      personality: formControl<Nullable<string>>(null),
      history: formControl<Nullable<string>>(null),
      groupTalkativeness: formControl(0)
    }),
    tags: formControl([]),
    dialogueExamples: formArray([]),
    greetings: formArray([]),
    groupGreetings: formArray([]),
  })

  readonly characterFG: TypedFormGroup<Character> =
    this.formGroup.get('character') as TypedFormGroup<Character>
  readonly tagsCtrl: TypedFormControl<Tag[]> =
    this.formGroup.get('tags') as TypedFormControl<Tag[]>
  readonly dialogueExamplesFA: TypedFormArray<string> =
    this.formGroup.get('dialogueExamples') as FormArray
  readonly greetingsFA: TypedFormArray<string> =
    this.formGroup.get('greetings') as FormArray
  readonly groupGreetingsFA: TypedFormArray<string> =
    this.formGroup.get('groupGreetings') as FormArray

  private readonly requestSubmit: Subject<void> = new Subject()
  readonly onSubmitRequested: Observable<void> = this.requestSubmit
  private readonly formReset: Subject<void> = new Subject()
  readonly onFormReset: Observable<void> = this.formReset

  resetFormData(formData: CharacterFormData) {
    this.formReset.next()
    this.formGroup.reset(formData)

    this.setControlsTo(this.dialogueExamplesFA, formData.dialogueExamples)
    this.setControlsTo(this.greetingsFA, formData.greetings)
    this.setControlsTo(this.groupGreetingsFA, formData.groupGreetings)
  }

  requestSubmitFn(): () => void {
    return () => this.requestSubmit.next()
  }

  addControlFn(): (arr: TypedFormArray<string>, value?: string) => void {
    return (arr, value) => {
      this.addControlTo(arr, value ?? '')
      arr.markAsDirty()
    }
  }

  removeControlFn(): (arr: TypedFormArray<string>, idx: number) => void {
    return (arr, idx) => {
      arr.removeAt(idx)
      arr.markAsDirty()
    }
  }

  private setControlsTo(arr: TypedFormArray<string>, values: string[]) {
    arr.clear()
    values.forEach(value => this.addControlTo(arr, value))
    arr.reset()
  }

  private addControlTo(arr: TypedFormArray<string>, value: string = '') {
    arr.push(formControl(value, [Validators.required]))
  }
}
