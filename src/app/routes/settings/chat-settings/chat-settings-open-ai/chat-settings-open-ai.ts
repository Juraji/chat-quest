import {Component, effect, inject} from '@angular/core';
import {formControl, formGroup} from '@util/ng';
import {ChatService, OpenAiSettings} from '@ai/chat';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Collapse} from '@components/collapse';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'app-chat-settings-open-ai',
  imports: [
    ReactiveFormsModule,
    Collapse,
  ],
  templateUrl: './chat-settings-open-ai.html',
})
export class ChatSettingsOpenAi {
  private readonly notifications = inject(Notifications);
  private readonly chat = inject(ChatService)

  readonly formGroup = formGroup<OpenAiSettings>({
    baseUri: formControl('https://api.openapi.com/v1', [Validators.required]),
    apiKey: formControl('', [Validators.required]),

    temperature: formControl(1.0, [Validators.required, Validators.min(0)]),
    maxTokens: formControl(256, [Validators.required, Validators.min(0)]),
    topP: formControl(0.9, [Validators.required, Validators.min(0)]),
    stream: formControl(true),
    stop: formControl(''),
  })

  constructor() {
    effect(() => {
      const settings = this.chat.settings()
      if (!!settings) {
        this.formGroup.reset(settings)
      }
    });
  }

  onSubmit() {
    if (this.formGroup.invalid) return
    const settings = this.formGroup.getRawValue()
    this.chat.settings.set(settings)
    this.notifications.toast("Settings updated!")
  }

  onTestConnection() {
    if (this.formGroup.dirty || this.formGroup.invalid) return
    this.chat
      .getModels()
      .subscribe({
        next: data => {
          const models = data.data.map(m => `'${m.id}'`).join(', ')
          this.notifications.toast(`Connection successful. Found models ${models}.`);
        },
        error: error => this.notifications.toast(`Failed to connect: ${error.message}`, 'DANGER')
      })
  }
}
