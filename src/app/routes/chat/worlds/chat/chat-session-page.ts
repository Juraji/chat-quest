import {Component, computed, inject, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {World} from '@api/worlds';
import {ChatMessage, ChatParticipant, ChatSession} from '@api/chat-sessions';
import {Character} from '@api/characters';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';
import {SSE} from '@api/sse';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {filter, map} from 'rxjs';

@Component({
  selector: 'chat-with-page',
  imports: [
    PageHeader,
  ],
  templateUrl: './chat-session-page.html'
})
export class ChatSessionPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly sse = inject(SSE)

  readonly world: Signal<World> = routeDataSignal(this.activatedRoute, 'world')
  readonly chatSession: Signal<ChatSession> = routeDataSignal(this.activatedRoute, 'chatSession')
  readonly chatSessionId: Signal<number> = computed(() => this.chatSession().id)
  readonly chatSessionName: Signal<string> = computed(() => this.chatSession().name)

  private readonly allCharacters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'allCharacters')
  private readonly originalParticipants: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'participants')
  readonly participantsMap: WritableSignal<Record<number, Character>> = linkedSignal(() =>
    this.originalParticipants()
      .reduce((acc, next) => {
        acc[next.id] = next
        return acc
      }, {} as Record<number, Character>))
  readonly participants = computed(() => Object.values(this.participantsMap()))

  private readonly originalMessages: Signal<ChatMessage[]> = routeDataSignal(this.activatedRoute, 'messages')
  private readonly messagesMap: WritableSignal<Record<number, ChatMessage>> = linkedSignal(() =>
    this.originalMessages()
      .reduce((acc, next) => {
        acc[next.id] = next
        return acc
      }, {} as Record<number, ChatMessage>))
  readonly messages: Signal<ChatMessage[]> = linkedSignal(() => Object
    .values(this.messagesMap())
    .sort((a, b) => a.createdAt!.localeCompare(b.createdAt!)))

  constructor() {
    const sessionFilter: <T extends { chatSessionId: number }>(input: T) => boolean =
      input => input.chatSessionId === this.chatSessionId()

    this.sse.on<ChatParticipant>('ChatParticipantAdded')
      .pipe(
        takeUntilDestroyed(), filter(sessionFilter),
        map(({characterId}) => this.allCharacters()
          .find(({id}) => id === characterId)!),
      )
      .subscribe(c => this.participantsMap
        .update(map => ({...map, [c.id]: c})))
    this.sse.on<ChatParticipant>('ChatParticipantRemoved')
      .pipe(
        takeUntilDestroyed(), filter(sessionFilter),
        map(({characterId}) => characterId)
      )
      .subscribe(characterId => this.participantsMap
        .update(map => {
          const update = {...map}
          delete update[characterId]
          return update
        }))

    this.sse.on<ChatMessage>('ChatMessageCreatedSignal')
      .pipe(takeUntilDestroyed(), filter(sessionFilter))
      .subscribe(message => this.messagesMap
        .update(map => ({...map, [message.id]: message})))
    this.sse.on<ChatMessage>('ChatMessageUpdatedSignal')
      .pipe(takeUntilDestroyed(), filter(sessionFilter))
      .subscribe(message => this.messagesMap
        .update(map => ({...map, [message.id]: message})))
    this.sse.on<number>('ChatMessageDeletedSignal')
      .pipe(takeUntilDestroyed())
      .subscribe(characterId => this.messagesMap
        .update(map => {
          const update = {...map}
          delete update[characterId]
          return update
        }))
  }
}
