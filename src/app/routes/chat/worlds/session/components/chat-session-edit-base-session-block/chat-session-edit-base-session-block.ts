import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {DatePipe} from '@angular/common';
import {formControl, formGroup, readOnlyControl} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {ChatSessionData} from '../../chat-session-data';
import {Scenario} from '@api/scenarios';

@Component({
  selector: 'chat-session-edit-base-session-block',
  imports: [
    DatePipe,
    ReactiveFormsModule
  ],
  templateUrl: './chat-session-edit-base-session-block.html',
})
export class ChatSessionEditBaseSessionBlock {
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions)
  private readonly notifications = inject(Notifications)

  readonly worldName: Signal<string> = computed(() => this.sessionData.chatSession().name)
  readonly createdAt: Signal<Nullable<string>> = computed(() => this.sessionData.chatSession().createdAt)
  readonly scenarios: Signal<Scenario[]> = this.sessionData.scenarios

  readonly formGroup = formGroup<ChatSession>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    createdAt: readOnlyControl(),
    name: formControl('', [Validators.required]),
    scenarioId: formControl(0),
    enableMemories: formControl(true),
  })

  constructor() {
    effect(() => {
      const session = this.sessionData.chatSession();
      this.formGroup.reset(session)
    });
  }

  onUpdateSession() {
    if (this.formGroup.invalid) return

    const update: ChatSession = {
      ...this.sessionData.chatSession(),
      ...this.formGroup.value
    }

    this.chatSessions
      .save(this.sessionData.worldId(), update)
      .subscribe(() => this.notifications.toast("Session details updated!"))
  }
}
