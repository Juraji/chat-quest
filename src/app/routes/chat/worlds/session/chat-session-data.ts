import {computed, inject, Injectable, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {SSE} from '@api/sse';
import {routeDataSignal} from '@util/ng';
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
  ChatSessionUpdated,
  sessionEntityFilter
} from '@api/chat-sessions';
import {Character, CharacterCreated, CharacterDeleted, CharacterUpdated} from '@api/characters';
import {Scenario, ScenarioCreated, ScenarioDeleted, ScenarioUpdated} from '@api/scenarios';
import {entityIdFilter} from '@api/common';
import {arrayAddItem, arrayRemoveItem, arrayUpsertItem} from '@util/array';
import {map} from 'rxjs';

@Injectable()
export class ChatSessionData {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly sse = inject(SSE)

  private readonly _world: Signal<World> = routeDataSignal(this.activatedRoute, 'world')
  readonly world: WritableSignal<World> = linkedSignal(() => this._world())
  readonly worldId: Signal<number> = computed(() => this.world().id)

  private readonly _chatSession: Signal<ChatSession> = routeDataSignal(this.activatedRoute, 'chatSession')
  readonly chatSession: WritableSignal<ChatSession> = linkedSignal(() => this._chatSession())
  readonly chatSessionId: Signal<number> = computed(() => this.chatSession().id)

  private readonly _participants: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'participants')
  readonly participants: WritableSignal<Character[]> = linkedSignal(() => this._participants())

  private readonly _messages: Signal<ChatMessage[]> = routeDataSignal(this.activatedRoute, 'messages')
  readonly messages: WritableSignal<ChatMessage[]> = linkedSignal(() => this._messages())

  private readonly _characters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters')
  readonly characters: WritableSignal<Character[]> = linkedSignal(() => this._characters())

  private readonly _scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios')
  readonly scenarios: WritableSignal<Scenario[]> = linkedSignal(() => this._scenarios())

  constructor() {
    this.sse
      .on(ChatSessionUpdated, entityIdFilter(this.chatSessionId))
      .subscribe(cs => this.chatSession.set(cs))
    this.sse
      .on(ChatSessionDeleted, entityIdFilter(this.chatSessionId))
      .subscribe(() => this.router.navigate(['/chat/worlds', this.worldId()]))

    this.sse
      .on(CharacterCreated)
      .subscribe(c => this.characters
        .update(prev => arrayAddItem(prev, c)));
    this.sse
      .on(CharacterUpdated)
      .subscribe(c => this.characters
        .update(prev => arrayUpsertItem(prev, c, ({id}) => id === c.id)));
    this.sse
      .on(CharacterDeleted)
      .subscribe(characterId => this.characters
        .update(prev => arrayRemoveItem(prev, ({id}) => id === characterId)))

    this.sse
      .on(ScenarioCreated)
      .subscribe(s => this.scenarios
        .update(prev => arrayAddItem(prev, s)))
    this.sse
      .on(ScenarioUpdated)
      .subscribe(s => this.scenarios
        .update(prev => arrayUpsertItem(prev, s, ({id}) => id === s.id)));
    this.sse
      .on(ScenarioDeleted)
      .subscribe(scenarioId => this.scenarios
        .update(prev => arrayRemoveItem(prev, ({id}) => id === scenarioId)))

    this.sse
      .on(ChatParticipantAdded, sessionEntityFilter(this.chatSessionId))
      .pipe(map(({characterId}) => this.characters().find(({id}) => id === characterId)!))
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
}
