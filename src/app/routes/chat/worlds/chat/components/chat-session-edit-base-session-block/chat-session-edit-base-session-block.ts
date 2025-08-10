import {Component, effect, inject, input, InputSignal} from '@angular/core';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {World} from '@api/worlds';
import {DatePipe} from '@angular/common';
import {formControl, formGroup, readOnlyControl} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'chat-session-edit-base-session-block',
  imports: [
    DatePipe,
    ReactiveFormsModule
  ],
  templateUrl: './chat-session-edit-base-session-block.html',
})
export class ChatSessionEditBaseSessionBlock {
  private readonly chatSessions = inject(ChatSessions)
  private readonly notifications = inject(Notifications)

  readonly world: InputSignal<World> = input.required()
  readonly chatSession: InputSignal<ChatSession> = input.required()

  readonly formGroup = formGroup<ChatSession>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    createdAt: readOnlyControl(),
    name: formControl('', [Validators.required]),
    scenarioId: readOnlyControl(),
    enableMemories: formControl(true),
  })

  constructor() {
    effect(() => {
      const session = this.chatSession();
      this.formGroup.reset(session)
    });
  }

  onUpdateSession() {
    if (this.formGroup.invalid) return

    const update: ChatSession = {
      ...this.chatSession(),
      ...this.formGroup.value
    }

    this.chatSessions
      .save(this.world().id, update)
      .subscribe(() => this.notifications.toast("Session details updated!"))
  }
}
