import {Component, computed, input, InputSignal, Signal} from '@angular/core';
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

  readonly dialogExample: string = `{user}: Hi, I am User. What's your name?
{char}: *looks at {user} startled* "Ow, hello. I didn't notice you there. Nice to meet you, my name is {char}."`

  readonly likelyActionsFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('likelyActions') as FormArray)
  readonly unlikelyActionsFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('unlikelyActions') as FormArray)
  readonly dialogueExamplesFA: Signal<TypedFormArray<string>> =
    computed(() => this.parentForm().get('dialogueExamples') as FormArray)

  onAddLikelyAction() {
    this.likelyActionsFA().push(formControl('', [Validators.required]))
  }

  onRemoveLikelyAction(idx: number) {
    this.likelyActionsFA().removeAt(idx)
  }

  onAddUnlikelyAction() {
    this.unlikelyActionsFA().push(formControl('', [Validators.required]))
  }

  onRemoveUnlikelyAction(idx: number) {
    this.unlikelyActionsFA().removeAt(idx)
  }

  onAddDialogExample() {
    this.dialogueExamplesFA().push(formControl('', [Validators.required]))
  }

  onRemoveDialogExample(idx: number) {
    this.dialogueExamplesFA().removeAt(idx)
  }
}
