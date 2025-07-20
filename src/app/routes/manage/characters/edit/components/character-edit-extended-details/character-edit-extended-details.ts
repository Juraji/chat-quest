import {Component, computed, effect, input, InputSignal, Signal} from '@angular/core';
import {Character} from '@db/characters';
import {FormArray, ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, TypedFormArray, TypedFormGroup} from '@util/ng';

@Component({
  selector: 'app-character-edit-extended-details',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './character-edit-extended-details.html'
})
export class CharacterEditExtendedDetails {
  readonly parentForm: InputSignal<TypedFormGroup<Character>> = input.required()
  readonly character: InputSignal<Character> = input.required()

  readonly dialogExample: string = `{user}: Hi, I am User. What's your name?
{char}: *looks at {user} startled* "Ow, hello. I didn't notice you there. Nice to meet you, my name is {char}."`

  readonly likelyActionsFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('likelyActions') as FormArray)
  readonly unlikelyActionsFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('unlikelyActions') as FormArray)
  readonly dialogueExamplesFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('dialogueExamples') as FormArray)

  constructor() {
    effect(() => {
      const {
        likelyActions,
        unlikelyActions,
        dialogueExamples,
      } = this.character()

      const la = this.likelyActionsFA()
      la.clear()
      likelyActions.forEach(action => this.addControl(la, action));

      const ua = this.unlikelyActionsFA()
      ua.clear()
      unlikelyActions.forEach(action => this.addControl(ua, action));

      const de = this.dialogueExamplesFA()
      de.clear()
      dialogueExamples.forEach(action => this.addControl(de, action));
    });
  }

  onAddLikelyAction() {
    const fa = this.likelyActionsFA()
    this.addControl(fa)
    fa.markAsDirty()
  }

  onRemoveLikelyAction(idx: number) {
    const fa = this.likelyActionsFA()
    fa.removeAt(idx)
    fa.markAsDirty()
  }

  onAddUnlikelyAction() {
    const fa = this.unlikelyActionsFA()
    this.addControl(fa)
    fa.markAsDirty()
  }

  onRemoveUnlikelyAction(idx: number) {
    const fa = this.unlikelyActionsFA()
    fa.removeAt(idx)
    fa.markAsDirty()
  }

  onAddDialogExample() {
    const fa = this.dialogueExamplesFA()
    this.addControl(fa)
    fa.markAsDirty()
  }

  onRemoveDialogExample(idx: number) {
    const fa = this.dialogueExamplesFA()
    fa.removeAt(idx)
    fa.markAsDirty()
  }

  private addControl(arr: TypedFormArray<string>, value: string = '') {
    arr.push(formControl(value, [Validators.required]))
  }
}
