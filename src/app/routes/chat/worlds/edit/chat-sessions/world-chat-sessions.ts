import {Component, inject, Signal} from '@angular/core';
import {ReactiveFormsModule} from '@angular/forms';
import {NewChatSession} from './components/new-chat-session-form/new-chat-session';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';
import {ChatSession, chatSessionSortingTransformer} from '@api/chat-sessions';
import {ChatSessionRecord} from './components/chat-session-record/chat-session-record';

@Component({
  selector: 'world-chat-sessions',
  imports: [
    ReactiveFormsModule,
    NewChatSession,
    ChatSessionRecord
  ],
  templateUrl: './world-chat-sessions.html'
})
export class WorldChatSessions {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly chatSessions: Signal<ChatSession[]> =
    routeDataSignal(this.activatedRoute, 'chatSessions', chatSessionSortingTransformer);
}
