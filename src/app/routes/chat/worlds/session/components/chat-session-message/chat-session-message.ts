import {Component, computed, inject, input, InputSignal, Signal} from '@angular/core';
import {ChatMessage, ChatSessions} from '@api/chat-sessions';
import {RenderedMessage} from '@components/rendered-message';
import {Character} from '@api/characters';
import {Notifications} from '@components/notifications';
import {ChatSessionData} from '../../chat-session-data';

@Component({
  selector: 'chat-session-message',
  imports: [
    RenderedMessage
  ],
  templateUrl: './chat-session-message.html',
  styleUrl: './chat-session-message.scss'
})
export class ChatSessionMessage {
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions)
  private readonly notifications = inject(Notifications)

  readonly worldId: Signal<number> = this.sessionData.worldId

  readonly message: InputSignal<ChatMessage> = input.required()

  readonly content: Signal<string> = computed(() => this.message().content)
  readonly isUser: Signal<boolean> = computed(() => this.message().isUser)
  readonly isGenerating: Signal<boolean> = computed(() => this.message().isGenerating)
  readonly createdAt: Signal<string> = computed(() => this.message().createdAt!)

  readonly character: Signal<Nullable<Character>> = computed(() => {
    const characterId = this.message().characterId
    const characters = this.sessionData.characters()
    if (!characterId) return null
    else return characters.find(({id}) => id === characterId)
  })
  readonly characterName = computed(() => this.character()?.name || 'You')
  readonly characterAvatar = computed(() => this.character()?.avatarUrl)

  onDeleteMessage() {
    const doDelete = confirm(`Are you sure you want to delete the message?
Note that this and all subsequent messages will be deleted and this action can not be undone.`);

    if (doDelete) {
      const message = this.message()

      this.chatSessions
        .deleteMessage(this.worldId(), message.chatSessionId, message.id)
        .subscribe(() => this.notifications.toast('Message deleted!'))
    }
  }
}
