import {Component, computed, inject, Signal} from '@angular/core';
import {Character} from '@api/characters';
import {ChatParticipant, ChatSessions} from '@api/chat-sessions';
import {CharacterCard} from '@components/cards/character-card';
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

  private readonly worldId: Signal<number> = this.sessionData.worldId
  private readonly chatSessionId: Signal<number> = this.sessionData.chatSessionId
  readonly participants: Signal<ChatParticipant[]> = this.sessionData.participants

  readonly available: Signal<Character[]> = computed(() => {
    const all = this.sessionData.characters()
    const pIds = this.sessionData.participants().map(x => x.characterId)
    return all.filter(c => !pIds.includes(c.id));
  })

  onAddParticipant(char: Character) {
    this.chatSessions
      .addParticipant(this.worldId(), this.chatSessionId(), char.id, false)
      .subscribe()
  }

  onRemoveParticipant(p: ChatParticipant) {
    this.chatSessions
      .removeParticipant(this.worldId(), this.chatSessionId(), p.characterId)
      .subscribe()
  }

  onTriggerResponse(p: ChatParticipant) {
    this.chatSessions
      .triggerParticipantResponse(this.worldId(), this.chatSessionId(), p.characterId)
      .subscribe()
  }

  onToggleMute(p: ChatParticipant) {
    this.chatSessions
      .addParticipant(this.worldId(), this.chatSessionId(), p.characterId, !p.muted)
      .subscribe()
  }
}
