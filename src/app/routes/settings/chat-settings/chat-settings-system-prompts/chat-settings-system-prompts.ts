import {Component, effect, inject, signal, Signal, WritableSignal} from '@angular/core';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, formGroup, readOnlyControl} from '@util/ng';
import {toSignal} from '@angular/core/rxjs-interop';
import {SystemPrompt, SystemPrompts, SystemPromptType} from '@db/system-prompts';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'app-chat-settings-system-prompts',
  imports: [
    ReactiveFormsModule,
    FormsModule
  ],
  templateUrl: './chat-settings-system-prompts.html',
  styleUrls: ['./chat-settings-system-prompts.scss']
})
export class ChatSettingsSystemPrompts {
  private readonly systemPrompts = inject(SystemPrompts)
  private readonly notifications = inject(Notifications)

  readonly availablePrompts: Signal<SystemPrompt[]> = toSignal(this.systemPrompts.getAll(true), {initialValue: []})

  readonly editedPrompt: WritableSignal<SystemPrompt | null> = signal(null);
  readonly formGroup = formGroup<SystemPrompt>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    prompt: formControl('', [Validators.required]),
    type: formControl<SystemPromptType>('CHAT', [Validators.required]),
  })

  constructor() {
    effect(() => {
      const p = this.editedPrompt()
      if (!!p) {
        this.formGroup.reset(p)
      } else {
        this.formGroup.reset()
      }
    });
  }

  onAddPrompt() {
    this.editedPrompt.set({
      id: null as any,
      name: '',
      prompt: '',
      type: 'CHAT'
    })
  }

  onSelectPrompt(id: number) {
    const p = this.availablePrompts().find(p => p.id === id)
    if (!p) return
    this.editedPrompt.set(p)
  }

  onCopyPrompt() {
    this.editedPrompt.update(current => {
      if(!current) return current

      return {
        ...current,
        id: null as any,
        name: `${current.name} (Copy)`
      }
    })
  }

  onSubmit() {
    if (this.formGroup.invalid) return

    const prompt = this.formGroup.value

    this.systemPrompts
      .save(prompt)
      .subscribe(p => {
        this.notifications.toast(`Prompt ${p.name} was saved.`)
        this.editedPrompt.set(p)
      })
  }

  onDeletePrompt() {
    const prompt = this.editedPrompt()
    if (!prompt) return

    if (prompt.id == null) {
      this.editedPrompt.set(null)
    } else {
      const doDelete = confirm(`Are you sure you want to delete the system prompt ${prompt.name}?`)

      if (doDelete) {
        this.systemPrompts
          .delete(prompt.id)
          .subscribe(() => {
            this.notifications.toast(`Prompt ${prompt.name} was deleted.`)
            this.editedPrompt.set(null)
          })
      }
    }
  }
}
