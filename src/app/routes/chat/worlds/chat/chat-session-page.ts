import {Component, computed, inject, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {World} from '@api/worlds';
import {ChatSession, ChatSessionDeleted, ChatSessionUpdated} from '@api/chat-sessions';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {SSE} from '@api/sse';
import {
  ChatSessionChatInputBlock,
  ChatSessionEditBaseSessionBlock,
  ChatSessionMessagesBlock,
  ChatSessionParticipantsBlock
} from './components';
import {entityIdFilter} from '@api/common';

@Component({
  selector: 'chat-with-page',
  imports: [
    PageHeader,
    ChatSessionEditBaseSessionBlock,
    ChatSessionMessagesBlock,
    ChatSessionParticipantsBlock,
    ChatSessionChatInputBlock,

  ],
  templateUrl: './chat-session-page.html',
  styleUrls: ["./chat-session-page.scss"]
})
export class ChatSessionPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly sse = inject(SSE)

  readonly world: Signal<World> = routeDataSignal(this.activatedRoute, 'world')
  readonly worldId: Signal<number> = computed(() => this.world().id)

  private readonly initialChatSession: Signal<ChatSession> = routeDataSignal(this.activatedRoute, 'chatSession')
  readonly chatSession: WritableSignal<ChatSession> = linkedSignal(() => this.initialChatSession())
  readonly chatSessionId: Signal<number> = computed(() => this.chatSession().id)
  readonly chatSessionName: Signal<string> = computed(() => this.chatSession().name)

  constructor() {
    this.sse
      .on(ChatSessionUpdated, entityIdFilter(this.chatSessionId))
      .subscribe(cs => this.chatSession.set(cs))
    this.sse
      .on<number>(ChatSessionDeleted, entityIdFilter(this.chatSessionId))
      .subscribe(() => this.router.navigate(['/chat']))
  }
}
