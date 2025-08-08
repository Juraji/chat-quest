import {Component, computed, effect, inject, input, InputSignal} from '@angular/core';
import {formControl, formGroup} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Worlds} from '@api/worlds/worlds.service';
import {Notifications} from '@components/notifications';
import {ChatPreferences} from '@api/worlds';
import {Instruction} from '@api/instructions';
import {LlmModelView} from '@api/providers';

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
  readonly instructionTemplates: InputSignal<Instruction[]> = input.required()
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
