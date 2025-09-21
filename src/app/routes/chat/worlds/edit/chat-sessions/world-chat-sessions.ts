import {Component, inject, Signal} from '@angular/core';
import {ReactiveFormsModule} from '@angular/forms';
import {paramAsId, routeDataSignal, routeParamSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {ChatSession, ChatSessions, chatSessionSortingTransformer} from '@api/chat-sessions';
import {ChatSessionRecord} from './components/chat-session-record/chat-session-record';
import {NEW_ID} from '@api/common';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'world-chat-sessions',
  imports: [
    ReactiveFormsModule,
    ChatSessionRecord
  ],
  templateUrl: './world-chat-sessions.html'
})
export class WorldChatSessions {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly notifications = inject(Notifications);
  private readonly router = inject(Router);
  private readonly chatSessionsService = inject(ChatSessions)

  readonly worldId: Signal<number> = routeParamSignal(this.activatedRoute, 'worldId', paramAsId);
  readonly chatSessions: Signal<ChatSession[]> =
    routeDataSignal(this.activatedRoute, 'chatSessions', chatSessionSortingTransformer);

  onNewChatSession() {
    const worldId = this.worldId();
    const chatSession: ChatSession = {
      id: NEW_ID,
      worldId,
      createdAt: null,
      name: 'New Session',
      scenarioId: null,
      generateMemories: true,
      useMemories: true,
      autoArchiveMessages: true,
      pauseAutomaticResponses: false,
      currentTimeOfDay: null,
      chatNotes: null
    }

    this.chatSessionsService
      .save(worldId, chatSession)
      .subscribe(cs => {
        this.notifications.toast("New session created!")
        this.router.navigate(
          ['..', 'session', cs.id],
          {relativeTo: this.activatedRoute}
        );
      })
  }
}
