import {Component, computed, inject, input, InputSignal} from '@angular/core';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {DatePipe} from '@angular/common';
import {ActivatedRoute, Router, RouterLink} from '@angular/router';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'chat-session-record',
  imports: [
    DatePipe,
    RouterLink
  ],
  templateUrl: './chat-session-record.html'
})
export class ChatSessionRecord {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly notifications = inject(Notifications);
  private readonly router = inject(Router)
  private readonly chatSessions = inject(ChatSessions)

  readonly chatSession: InputSignal<ChatSession> = input.required()
  readonly worldId = computed(() => this.chatSession().worldId)
  readonly sessionCreatedAt = computed(() => this.chatSession().createdAt)
  readonly sessionName = computed(() => this.chatSession().name)
  readonly sessionId = computed(() => this.chatSession().id)

  onDeletedChatSession() {
    const doDelete = confirm(`Are you sure you want to delete the session '${this.sessionName()}'?`);
    if (doDelete) {
      this.chatSessions
        .delete(this.worldId(), this.sessionId())
        .subscribe(() => {
          this.notifications.toast("Chat session deleted!")
          this.router.navigate([], {
            relativeTo: this.activatedRoute,
            replaceUrl: true,
            queryParams: {
              u: Date.now()
            }
          })
        })
    }
  }
}
