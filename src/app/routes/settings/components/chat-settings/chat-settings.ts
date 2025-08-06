import {Component, computed, effect, inject, input, InputSignal} from '@angular/core';
import {formControl, formGroup} from '@util/ng';
import {ChatPreferences, InstructionTemplate, LlmModelView} from '@api/model';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Worlds} from '@api/clients/worlds';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'chat-settings',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './chat-settings.html'
})
export class ChatSettings {
  private readonly worlds = inject(Worlds)
  private readonly notifications = inject(Notifications);

  readonly preferences: InputSignal<ChatPreferences> = input.required()
  readonly instructionTemplates: InputSignal<InstructionTemplate[]> = input.required()
  readonly llmModelViews: InputSignal<LlmModelView[]> = input.required()

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

    const update: ChatPreferences = this.formGroup.getRawValue()

    this.worlds
      .saveChatPreferences(update)
      .subscribe(res => {
        this.formGroup.reset(res)
        this.notifications.toast("Chat settings updated!")
      })
  }
}
