import {Component, computed, effect, inject, Signal} from '@angular/core';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {booleanSignal, formControl, formGroup, readOnlyControl, routeDataSignal, toControlValueSignal} from '@util/ng';
import {InstructionTemplate, InstructionType, isNew} from '@api/model';
import {ActivatedRoute, Router} from '@angular/router';
import {InstructionTemplates} from '@api/clients';
import {Notifications} from '@components/notifications';
import {RenderedMessage} from '@components/rendered-message/rendered-message';
import {TokenCount} from '@components/token-count/token-count';

@Component({
  selector: 'app-edit-instruction-template-page',
  imports: [
    FormsModule,
    PageHeader,
    ReactiveFormsModule,
    RenderedMessage,
    TokenCount
  ],
  templateUrl: './edit-instruction-template.html',
})
export class EditInstructionTemplate {
  private readonly templates = inject(InstructionTemplates)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)

  readonly template: Signal<InstructionTemplate> = routeDataSignal(this.activatedRoute, 'template');

  readonly isNew = computed(() => isNew(this.template()))

  readonly formGroup = formGroup<InstructionTemplate>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    type: formControl<InstructionType>('CHAT', [Validators.required]),
    temperature: formControl<Nullable<number>>(null, [Validators.min(0.01)]),
    systemPrompt: formControl('', [Validators.required]),
    instruction: formControl('', [Validators.required]),
  })

  readonly editSystemPrompt = booleanSignal(false)
  readonly systemPromptValue: Signal<string> = toControlValueSignal(this.formGroup, 'systemPrompt')

  readonly editInstruction = booleanSignal(false)
  readonly instructionValue: Signal<string> = toControlValueSignal(this.formGroup, 'instruction')

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
    const update: InstructionTemplate = {
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
