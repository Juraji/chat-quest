import {
  Component,
  computed,
  effect,
  ElementRef,
  inject,
  linkedSignal,
  Signal,
  viewChild,
  WritableSignal
} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {World} from '@api/worlds';
import {
  ChatMessage,
  ChatMessageCreated,
  ChatMessageDeleted,
  ChatMessageUpdated,
  ChatParticipantAdded,
  ChatParticipantRemoved,
  ChatSession,
  ChatSessionDeleted,
  ChatSessions,
  ChatSessionUpdated,
  sessionEntityFilter
} from '@api/chat-sessions';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {SSE} from '@api/sse';
import {ChatSessionChatInputBlock, ChatSessionEditBaseSessionBlock, ChatSessionParticipantsBlock} from './components';
import {entityIdFilter} from '@api/common';
import {arrayAddItem, arrayRemoveItem, arrayUpsertItem} from '@util/array';
import {ChatSessionMessage} from './components/chat-session-message/chat-session-message';
import {Character, CharacterCreated, CharacterDeleted, CharacterUpdated} from '@api/characters';
import {Notifications} from '@components/notifications';
import {map} from 'rxjs';
import {FindCharacterPipe} from '@components/find-character.pipe';

@Component({
  selector: 'chat-with-page',
  imports: [
    PageHeader,
    ChatSessionEditBaseSessionBlock,
    ChatSessionParticipantsBlock,
    ChatSessionChatInputBlock,
    ChatSessionMessage,
    FindCharacterPipe,

  ],
  templateUrl: './chat-session-page.html',
  styleUrls: ["./chat-session-page.scss"]
})
export class ChatSessionPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly chatSessions = inject(ChatSessions)
  private readonly notifications = inject(Notifications);
  private readonly router = inject(Router);
  private readonly sse = inject(SSE)

  protected readonly chatMessagesContainerRef: Signal<ElementRef<HTMLDivElement> | undefined> =
    viewChild('chatMessagesContainer', {read: ElementRef})

  readonly world: Signal<World> = routeDataSignal(this.activatedRoute, 'world')
  readonly worldId: Signal<number> = computed(() => this.world().id)

  private readonly initialChatSession: Signal<ChatSession> = routeDataSignal(this.activatedRoute, 'chatSession')
  readonly chatSession: WritableSignal<ChatSession> = linkedSignal(() => this.initialChatSession())
  readonly chatSessionId: Signal<number> = computed(() => this.chatSession().id)
  readonly chatSessionName: Signal<string> = computed(() => this.chatSession().name)

  private readonly initialAllCharacters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'allCharacters')
  readonly allCharacters: WritableSignal<Character[]> = linkedSignal(() => this.initialAllCharacters())

  private readonly initialParticipants: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'participants')
  readonly participants: WritableSignal<Character[]> = linkedSignal(() => this.initialParticipants())

  private readonly initialMessages: Signal<ChatMessage[]> = routeDataSignal(this.activatedRoute, 'messages')
  readonly messages: WritableSignal<ChatMessage[]> = linkedSignal(() => this.initialMessages())

  constructor() {
    effect(() => {
      this.messages()
      const element = this.chatMessagesContainerRef()?.nativeElement;
      if (!!element) requestAnimationFrame(() => element.scrollTop = element.scrollHeight)
    });

    this.sse
      .on(ChatSessionUpdated, entityIdFilter(this.chatSessionId))
      .subscribe(cs => this.chatSession.set(cs))
    this.sse
      .on(ChatSessionDeleted, entityIdFilter(this.chatSessionId))
      .subscribe(() => this.router.navigate(['/chat']))

    this.sse
      .on(CharacterCreated)
      .subscribe(c => this.allCharacters
        .update(prev => arrayAddItem(prev, c)));
    this.sse
      .on(CharacterUpdated)
      .subscribe(c => this.allCharacters
        .update(prev => arrayUpsertItem(prev, c, ({id}) => id === c.id)));
    this.sse
      .on(CharacterDeleted)
      .subscribe(characterId => this.allCharacters
        .update(prev => arrayRemoveItem(prev, ({id}) => id === characterId)))

    this.sse
      .on(ChatParticipantAdded, sessionEntityFilter(this.chatSessionId))
      .pipe(map(({characterId}) => this.allCharacters().find(({id}) => id === characterId)!))
      .subscribe(c => this.participants
        .update(prev => arrayAddItem(prev, c)))
    this.sse
      .on(ChatParticipantRemoved, sessionEntityFilter(this.chatSessionId))
      .pipe(map(({characterId}) => characterId))
      .subscribe(characterId => this.participants
        .update(prev => arrayRemoveItem(prev, ({id}) => id === characterId)))

    this.sse
      .on(ChatMessageCreated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayAddItem(prev, message)))
    this.sse
      .on(ChatMessageUpdated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayUpsertItem(prev, message, ({id}) => id === message.id)))
    this.sse
      .on(ChatMessageDeleted)
      .subscribe(messageId => this.messages
        .update(prev => arrayRemoveItem(prev, ({id}) => id === messageId)))
  }

  onDeleteMessageRequest(message: ChatMessage) {
    const doDelete = confirm(`Are you sure you want to delete the message?
Note that this and all subsequent messages will be deleted and this action can not be undone.`);

    if (doDelete) {
      this.chatSessions
        .deleteMessage(this.worldId(), this.chatSessionId(), message.id)
        .subscribe(() => this.notifications.toast('Message deleted!'))
    }
  }
}
