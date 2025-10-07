import {Component, computed, effect, inject, signal, Signal, WritableSignal} from '@angular/core';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {booleanSignal, controlValueSignal, formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {Notifications} from '@components/notifications';
import {RenderedMessage} from '@components/rendered-message/rendered-message';
import {Instruction, Instructions, InstructionType} from '@api/instructions';
import {isNew, NEW_ID} from '@api/common';

@Component({
  selector: 'app-edit-instruction-template-page',
  imports: [
    FormsModule,
    PageHeader,
    ReactiveFormsModule,
    RenderedMessage,
  ],
  templateUrl: './edit-instruction.html',
})
export class EditInstruction {
  private readonly instructions = inject(Instructions)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)


  readonly instruction: Signal<Instruction> = routeDataSignal(this.activatedRoute, 'template');

  readonly isNew = computed(() => isNew(this.instruction()))

  readonly formGroup = formGroup<Instruction>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    type: formControl<InstructionType>('CHAT', [Validators.required]),
    temperature: formControl<number>(1.1, [Validators.required, Validators.min(0)]),
    maxTokens: formControl<number>(300, [Validators.required, Validators.min(1)]),
    topP: formControl<number>(0.95, [Validators.required, Validators.min(0.01), Validators.max(1.0)]),
    presencePenalty: formControl<number>(1.1, [Validators.required, Validators.min(0)]),
    frequencyPenalty: formControl<number>(1.1, [Validators.required, Validators.min(0)]),
    stream: formControl<boolean>(true),
    stopSequences: formControl<Nullable<string>>(null),
    includeReasoning: formControl(false),

    reasoningPrefix: formControl('', [Validators.required, Validators.maxLength(50)]),
    reasoningSuffix: formControl('', [Validators.required, Validators.maxLength(50)]),
    characterIdPrefix: formControl('', [Validators.required, Validators.maxLength(50)]),
    characterIdSuffix: formControl('', [Validators.required, Validators.maxLength(50)]),

    systemPrompt: formControl('', [Validators.required]),
    worldSetup: formControl('', [Validators.required]),
    instruction: formControl('', [Validators.required]),
  })

  readonly selectedTab: WritableSignal<number> = signal(0)

  readonly editSystemPrompt = booleanSignal(false)
  readonly systemPromptValue: Signal<string> = controlValueSignal(this.formGroup, 'systemPrompt')

  readonly editWorldSetup = booleanSignal(false)
  readonly worldSetupValue: Signal<string> = controlValueSignal(this.formGroup, 'worldSetup')

  readonly editInstruction = booleanSignal(false)
  readonly instructionValue: Signal<string> = controlValueSignal(this.formGroup, 'instruction')

  constructor() {
    effect(() => {
      const input = this.instruction()
      this.formGroup.reset(input)

      const n = isNew(input)
      this.editSystemPrompt.set(n)
      this.editWorldSetup.set(n)
      this.editInstruction.set(n)
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value
    const update: Instruction = {
      ...this.instruction(),
      ...formValue
    }

    this.instructions
      .save(update)
      .subscribe(res => {
        this.notifications.toast("Instruction Template saved!")
        this.router.navigate(['..', res.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onRevertChanges() {
    this.formGroup.reset(this.instruction());
  }

  onDuplicateInstruction() {
    const instruction = this.instruction()

    const newInstruction = {
      ...instruction,
      id: NEW_ID,
      name: instruction.name + ' (copy)',
    }

    this.instructions
      .save(newInstruction)
      .subscribe(res => {
        this.notifications.toast(`Instruction copied as "${res.name}"!`)
        this.router.navigate(['..', res.id], {
          relativeTo: this.activatedRoute
        })
      })
  }

  onDeleteTemplate() {
    const t = this.instruction();
    if (isNew(t)) return
    const doDelete = confirm(`Are you sure you want to delete this template?`)

    if (doDelete) {
      this.instructions
        .delete(t!.id)
        .subscribe(() => {
          this.notifications.toast("Instruction Template deleted!")
          this.router.navigate(['../..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          })
        })
    }
  }
}
