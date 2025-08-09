import {Component, computed, effect, inject, Signal} from '@angular/core';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {booleanSignal, controlValueSignal, formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {Notifications} from '@components/notifications';
import {RenderedMessage} from '@components/rendered-message/rendered-message';
import {TokenCount} from '@components/token-count/token-count';
import {Instruction, Instructions, InstructionType} from '@api/instructions';
import {isNew} from '@api/common';

@Component({
  selector: 'app-edit-instruction-template-page',
  imports: [
    FormsModule,
    PageHeader,
    ReactiveFormsModule,
    RenderedMessage,
    TokenCount
  ],
  templateUrl: './edit-instruction.html',
})
export class EditInstruction {
  private readonly templates = inject(Instructions)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)

  readonly template: Signal<Instruction> = routeDataSignal(this.activatedRoute, 'template');

  readonly isNew = computed(() => isNew(this.template()))

  readonly formGroup = formGroup<Instruction>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    type: formControl<InstructionType>('CHAT', [Validators.required]),
    temperature: formControl<Nullable<number>>(null, [Validators.min(0.01)]),
    systemPrompt: formControl('', [Validators.required]),
    instruction: formControl('', [Validators.required]),
  })

  readonly editSystemPrompt = booleanSignal(false)
  readonly systemPromptValue: Signal<string> = controlValueSignal(this.formGroup, 'systemPrompt')

  readonly editInstruction = booleanSignal(false)
  readonly instructionValue: Signal<string> = controlValueSignal(this.formGroup, 'instruction')

  constructor() {
    effect(() => {
      const input = this.template()
      this.formGroup.reset(input)
      this.editSystemPrompt.set(isNew(input))
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value
    const update: Instruction = {
      ...this.template(),
      ...formValue
    }

    this.templates
      .save(update)
      .subscribe(template => {
        this.notifications.toast("Instruction Template saved!")
        this.router.navigate(['..', template.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onRevertChanges() {
    this.formGroup.reset(this.template());
  }

  onDeleteTemplate() {
    const t = this.template();
    if (isNew(t)) return
    const doDelete = confirm(`Are you sure you want to delete this template?`)

    if (doDelete) {
      this.templates
        .delete(t!.id)
        .subscribe(() => {
          this.notifications.toast("Instruction Template deleted!")
          this.router.navigate(['..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          })
        })
    }
  }
}
