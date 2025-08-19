import {Component, computed, inject, input, InputSignal, Signal} from '@angular/core';
import {Character} from '@api/characters';
import {ChatSessions} from '@api/chat-sessions';
import {CharacterCard} from '@components/cards/character-card';
import {Notifications} from '@components/notifications';
import {Scalable} from '@components/scalable/scalable';
import {RouterLink} from '@angular/router';

@Component({
  selector: 'chat-session-participants-block',
  imports: [
    CharacterCard,
    Scalable,
    RouterLink
  ],
  templateUrl: './chat-session-participants-block.html',
})
export class ChatSessionParticipantsBlock {
  private readonly chatSessions = inject(ChatSessions);
  private readonly notifications = inject(Notifications);

  readonly worldId: InputSignal<number> = input.required()
  readonly chatSessionId: InputSignal<number> = input.required()
  readonly allCharacters: InputSignal<Character[]> = input.required()
  readonly participants: InputSignal<Character[]> = input.required()

  readonly available: Signal<Character[]> = computed(() => {
    const all = this.allCharacters()
    const pIds = this.participants().map(x => x.id)
    return all.filter(c => !pIds.includes(c.id));
  })

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
