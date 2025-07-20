import {Component, computed, effect, input, InputSignal} from '@angular/core';
import {formControl, TypedFormArray, TypedFormGroup} from '@util/ng';
import {Character} from '@db/characters';
import {FormArray, ReactiveFormsModule, Validators} from '@angular/forms';

@Component({
  selector: 'app-character-edit-chat-defaults',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './character-edit-chat-defaults.html'
})
export class CharacterEditChatDefaults {
  readonly parentForm: InputSignal<TypedFormGroup<Character>> = input.required()
  readonly character: InputSignal<Character> = input.required()

  readonly alternateGreetingsFA =
    computed(() => this.parentForm().get('alternateGreetings') as FormArray)
  readonly groupGreetingsFA =
    computed(() => this.parentForm().get('groupGreetings') as FormArray)

  constructor() {
    effect(() => {
      const {
        alternateGreetings,
        groupGreetings,
      } = this.character()

      const ag = this.alternateGreetingsFA()
      ag.clear()
      alternateGreetings.forEach(value => this.addControl(ag, value))

      const gg = this.groupGreetingsFA()
      gg.clear()
      groupGreetings.forEach(value => this.addControl(gg, value))
    });
  }

  private addControl(arr: TypedFormArray<string>, value: string = '') {
    arr.push(formControl(value, [Validators.required]))
  }

  onAddAlternateGreetingAction() {
    const fa = this.alternateGreetingsFA()
    this.addControl(fa)
    fa.markAsDirty()
  }

  onRemoveAlternateGreetingAction(idx: number) {
    const fa = this.alternateGreetingsFA()
    fa.removeAt(idx)
    fa.markAsDirty()
  }

  onAddGroupGreetingAction() {
    const fa = this.groupGreetingsFA()
    this.addControl(fa)
    fa.markAsDirty()
  }

  onRemoveGroupGreetingAction(idx: number) {
    const fa = this.groupGreetingsFA()
    fa.removeAt(idx)
    fa.markAsDirty()
  }
}
