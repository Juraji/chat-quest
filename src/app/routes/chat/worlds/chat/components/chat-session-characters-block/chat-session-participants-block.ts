import {Component, computed, inject, input, InputSignal, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {Character, characterSortingTransformer} from '@api/characters';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';
import {SSE} from '@api/sse';
import {ChatParticipantAdded, ChatParticipantRemoved, ChatSessions, sessionEntityFilter} from '@api/chat-sessions';
import {map} from 'rxjs';
import {arrayAddItem, arrayRemoveItem} from '@util/array';
import {CharacterCard} from '@components/cards/character-card';
import {Collapse} from '@components/collapse';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'chat-session-participants-block',
  imports: [
    CharacterCard,
    Collapse
  ],
  templateUrl: './chat-session-participants-block.html',
  styleUrls: ['./chat-session-participants-block.scss']
})
export class ChatSessionParticipantsBlock {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly chatSessions = inject(ChatSessions);
  private readonly notifications = inject(Notifications);
  private readonly sse = inject(SSE)

  readonly worldId: InputSignal<number> = input.required()
  readonly chatSessionId: InputSignal<number> = input.required()

  private readonly allCharacters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'allCharacters', characterSortingTransformer)
  private readonly initialParticipants: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'participants')

  readonly participants: WritableSignal<Character[]> = linkedSignal(() => this.initialParticipants())
  readonly available: Signal<Character[]> = computed(() => {
    const all = this.allCharacters()
    const pIds = this.participants().map(x => x.id)
    return all.filter(c => !pIds.includes(c.id));
  })

  constructor() {
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
  }

  onAddParticipant(char: Character) {
    this.chatSessions
      .addParticipant(this.worldId(), this.chatSessionId(), char.id)
      .subscribe(() => this.notifications.toast(`${char.name} has been added to the session`))
  }

  onRemoveParticipant(char: Character) {
    this.chatSessions
      .removeParticipant(this.worldId(), this.chatSessionId(), char.id)
      .subscribe(() => this.notifications.toast(`${char.name} has been removed from the session`))
  }
}
