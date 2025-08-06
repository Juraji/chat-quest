import {Component, computed, effect, input, InputSignal} from '@angular/core';
import {InstructionTemplate, MemoryPreferences} from '@api/model';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, formGroup} from '@util/ng';

@Component({
  selector: 'memory-settings',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './memory-settings.html',
})
export class MemorySettings {

  readonly preferences: InputSignal<MemoryPreferences> = input.required()
  readonly instructionTemplates: InputSignal<InstructionTemplate[]> = input.required()
  readonly memInstructionTemplates = computed(() => {
    return this.instructionTemplates().filter(t => t.type === 'MEMORIES')
  })

  readonly formGroup = formGroup<MemoryPreferences>({
    memoriesModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoriesInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
    embeddingModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoryMinP: formControl(0, [Validators.required, Validators.min(0), Validators.min(3)]),
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

  }
}
