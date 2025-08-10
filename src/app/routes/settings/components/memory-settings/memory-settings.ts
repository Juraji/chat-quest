import {booleanAttribute, Component, computed, effect, inject, input, InputSignal} from '@angular/core';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, formGroup, routeQueryParamSignal} from '@util/ng';
import {Notifications} from '@components/notifications';
import {Memories, MemoryPreferences} from '@api/memories';
import {Instruction} from '@api/instructions';
import {LlmModelView} from '@api/providers';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'memory-settings',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './memory-settings.html',
})
export class MemorySettings {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly memories = inject(Memories)
  private readonly notifications = inject(Notifications);

  readonly validate = routeQueryParamSignal(this.activatedRoute, 'validate', booleanAttribute)

  readonly preferences: InputSignal<MemoryPreferences> = input.required()
  readonly instructionTemplates: InputSignal<Instruction[]> = input.required()
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
    effect(() => {
      if (this.validate()) {
        this.formGroup.markAllAsDirty()
      }
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
