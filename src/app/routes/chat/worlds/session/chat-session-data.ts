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
  ChatParticipant,
  ChatParticipantAdded,
  ChatParticipantRemoved,
  ChatSession,
  ChatSessionDeleted,
  ChatSessionUpdated,
  sessionEntityFilter
} from '@api/chat-sessions';
import {Characters} from '@api/characters';
import {Scenario, ScenarioCreated, ScenarioDeleted, ScenarioUpdated} from '@api/scenarios';
import {entityIdFilter} from '@api/common';
import {arrayAdd, arrayRemove, arrayReplace} from '@util/array';
import {LlmModelView} from '@api/providers';
import {CQPreferences, PreferencesUpdated} from '@api/preferences';
import {MemoryCreated} from '@api/memories';
import {Notifications} from '@components/notifications';

@Injectable()
export class ChatSessionData {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly notifications = inject(Notifications);
  private readonly router = inject(Router);
  private readonly sse = inject(SSE)
  private readonly charactersService = inject(Characters)

  private readonly _world: Signal<World> = routeDataSignal(this.activatedRoute, 'world')
  readonly world: WritableSignal<World> = linkedSignal(() => this._world())
  readonly worldId: Signal<number> = computed(() => this.world().id)

  private readonly _chatSession: Signal<ChatSession> = routeDataSignal(this.activatedRoute, 'chatSession')
  readonly chatSession: WritableSignal<ChatSession> = linkedSignal(() => this._chatSession())
  readonly chatSessionId: Signal<number> = computed(() => this.chatSession().id)

  private readonly _participants: Signal<ChatParticipant[]> = routeDataSignal(this.activatedRoute, 'participants')
  readonly participants: WritableSignal<ChatParticipant[]> = linkedSignal(() => this._participants())

  private readonly _messages: Signal<ChatMessage[]> = routeDataSignal(this.activatedRoute, 'messages')
  readonly messages: WritableSignal<ChatMessage[]> = linkedSignal(() => this._messages())

  readonly characters = this.charactersService.all

  private readonly _scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios')
  readonly scenarios: WritableSignal<Scenario[]> = linkedSignal(() => this._scenarios())

  private readonly _preferences: Signal<CQPreferences> = routeDataSignal(this.activatedRoute, 'preferences')
  readonly preferences: WritableSignal<CQPreferences> = linkedSignal(() => this._preferences())

  readonly llmModels: Signal<LlmModelView[]> = routeDataSignal(this.activatedRoute, 'llmModels')

  constructor() {
    this.sse
      .on(ChatSessionUpdated, entityIdFilter(this.chatSessionId))
      .subscribe(cs => this.chatSession.set(cs))
    this.sse
      .on(ChatSessionDeleted, entityIdFilter(this.chatSessionId))
      .subscribe(() => this.router.navigate(['/chat/worlds', this.worldId()]))

    this.sse
      .on(ScenarioCreated)
      .subscribe(s => this.scenarios
        .update(prev => arrayAdd(prev, s)))
    this.sse
      .on(ScenarioUpdated)
      .subscribe(s => this.scenarios
        .update(prev => arrayReplace(prev, s, ({id}) => id === s.id)));
    this.sse
      .on(ScenarioDeleted)
      .subscribe(scenarioId => this.scenarios
        .update(prev => arrayRemove(prev, ({id}) => id === scenarioId)))

    this.sse
      .on(ChatParticipantAdded, sessionEntityFilter(this.chatSessionId))
      .subscribe(p => this.participants
        .update(prev => arrayReplace(prev, p, (op) => op.characterId === p.characterId)));
    this.sse
      .on(ChatParticipantRemoved, sessionEntityFilter(this.chatSessionId))
      .subscribe(p => this.participants
        .update(prev => arrayRemove(prev, ({characterId}) => characterId === p.characterId)))

    this.sse
      .on(ChatMessageCreated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayAdd(prev, message)))
    this.sse
      .on(ChatMessageUpdated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayReplace(prev, message, ({id}) => id === message.id)))
    this.sse
      .on(ChatMessageDeleted)
      .subscribe(messageId => this.messages
        .update(prev => arrayRemove(prev, ({id}) => id === messageId)))

    this.sse
      .on(PreferencesUpdated)
      .subscribe(prefs => this.preferences.set(prefs))

    this.sse
      .on(MemoryCreated)
      .subscribe(m => {
        let message: string
        if (!!m.characterId) {
          const char = this.characters().find(c => c.id === m.characterId)!
          message = `<span>A new memory was created for ${char.name}:</span><br/><p class="text-body-emphasis">${m.content}</p>`
        } else {
          message = `<span>A new general memory was created:</span><br/><p class="text-body-emphasis">${m.content}</p>`
        }

        this.notifications.toast(message)
      })
  }
}
