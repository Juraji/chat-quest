import {Component, computed, inject, Signal} from '@angular/core';
import {Character} from '@api/characters';
import {ChatSessions} from '@api/chat-sessions';
import {CharacterCard} from '@components/cards/character-card';
import {Notifications} from '@components/notifications';
import {Scalable} from '@components/scalable/scalable';
import {RouterLink} from '@angular/router';
import {ChatSessionData} from '../../chat-session-data';

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
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions);
  private readonly notifications = inject(Notifications);

  private readonly worldId: Signal<number> = this.sessionData.worldId
  private readonly chatSessionId: Signal<number> = this.sessionData.chatSessionId
  readonly participants: Signal<Character[]> = this.sessionData.participants

  readonly available: Signal<Character[]> = computed(() => {
    const all = this.sessionData.characters()
    const pIds = this.sessionData.participants().map(x => x.id)
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

  onTriggerResponse(char: Character) {
    this.chatSessions
      .triggerParticipantResponse(this.worldId(), this.chatSessionId(), char.id)
      .subscribe()
  }
}
