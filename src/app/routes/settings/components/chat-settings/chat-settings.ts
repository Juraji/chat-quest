import {Component, computed, effect, input, InputSignal} from '@angular/core';
import {formControl, formGroup} from '@util/ng';
import {ChatPreferences, InstructionTemplate} from '@api/model';
import {ReactiveFormsModule, Validators} from '@angular/forms';

@Component({
  selector: 'chat-settings',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './chat-settings.html'
})
export class ChatSettings {

  readonly preferences: InputSignal<ChatPreferences> = input.required()
  readonly instructionTemplates: InputSignal<InstructionTemplate[]> = input.required()
  readonly chatInstructionTemplates = computed(() => {
    return this.instructionTemplates().filter(t => t.type === 'CHAT')
  })

  readonly formGroup = formGroup<ChatPreferences>({
    chatModelId: formControl<Nullable<number>>(null, [Validators.required]),
    chatInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
  })

  constructor() {
    effect(() => {
      const p = this.preferences()
      this.formGroup.reset(p)
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return
  }
}
