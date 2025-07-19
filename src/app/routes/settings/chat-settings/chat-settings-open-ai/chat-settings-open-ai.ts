import {Component, effect, inject, Signal} from '@angular/core';
import {formControl, formGroup} from '@util/ng';
import {CHAT_SETTINGS_NAME, ChatService, OpenAiSettings} from '@ai/chat';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Collapse} from '@components/collapse';
import {Notifications} from '@components/notifications';
import {Settings} from '@db/settings/settings';
import {toSignal} from '@angular/core/rxjs-interop';

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
  private readonly settings = inject(Settings)
  private readonly chat = inject(ChatService)

  private readonly chatSettings: Signal<OpenAiSettings | null> = toSignal(this.settings
    .get<OpenAiSettings>(CHAT_SETTINGS_NAME, true), {initialValue: null});

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
      const settings = this.chatSettings()
      if (!!settings) {
        this.formGroup.reset(settings)
      } else {
        this.formGroup.markAsDirty()
      }
    });
  }

  onSubmit() {
    if (this.formGroup.invalid) return
    const settings = this.formGroup.getRawValue()
    this.settings
      .set(CHAT_SETTINGS_NAME, settings)
      .subscribe(() => this.notifications.toast("Open AI settings updated!"))
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
