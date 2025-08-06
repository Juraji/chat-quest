import {Component, computed, effect, inject, input, InputSignal} from '@angular/core';
import {InstructionTemplate, LlmModelView, MemoryPreferences} from '@api/model';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, formGroup} from '@util/ng';
import {Notifications} from '@components/notifications';
import {Memories} from '@api/clients/memories';

@Component({
  selector: 'memory-settings',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './memory-settings.html',
})
export class MemorySettings {
  private readonly memories = inject(Memories)
  private readonly notifications = inject(Notifications);

  readonly preferences: InputSignal<MemoryPreferences> = input.required()
  readonly instructionTemplates: InputSignal<InstructionTemplate[]> = input.required()
  readonly llmModelViews: InputSignal<LlmModelView[]> = input.required()

  readonly memInstructionTemplates = computed(() => {
    return this.instructionTemplates().filter(t => t.type === 'MEMORIES')
  })

  readonly formGroup = formGroup<MemoryPreferences>({
    memoriesModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoriesInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
    embeddingModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoryMinP: formControl(0, [Validators.required, Validators.min(0), Validators.max(1)]),
    memoryTriggerAfter: formControl(0, [Validators.required, Validators.min(1)]),
    memoryWindowSize: formControl(0, [Validators.required, Validators.min(1)]),
  })

  constructor() {
    effect(() => {
      const p = this.preferences()
      this.formGroup.reset(p)
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const update: MemoryPreferences = this.formGroup.getRawValue()

    this.memories
      .savePreferences(update)
      .subscribe(res => {
        this.formGroup.reset(res)
        this.notifications.toast('Memory settings updated!')
      })
  }
}
