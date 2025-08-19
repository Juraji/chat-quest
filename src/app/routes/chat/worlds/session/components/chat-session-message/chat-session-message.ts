import {Component, computed, input, InputSignal, output, Signal} from '@angular/core';
import {ChatMessage} from '@api/chat-sessions';
import {RenderedMessage} from '@components/rendered-message';
import {Character} from '@api/characters';

@Component({
  selector: 'chat-session-message',
  imports: [
    RenderedMessage
  ],
  templateUrl: './chat-session-message.html',
  styleUrl: './chat-session-message.scss'
})
export class ChatSessionMessage {
  readonly character: InputSignal<Nullable<Character>> = input()
  readonly message: InputSignal<ChatMessage> = input.required()

  readonly manageMemoryRequest = output()
  readonly deleteMessageRequest = output()

  readonly content: Signal<string> = computed(() => this.message().content)
  readonly isUser: Signal<boolean> = computed(() => this.message().isUser)
  readonly isGenerating: Signal<boolean> = computed(() => this.message().isGenerating)
  readonly createdAt: Signal<string> = computed(() => this.message().createdAt!)
  readonly memoryId: Signal<Nullable<number>> = computed(() => this.message().memoryId)

  readonly characterName = computed(() => this.character()?.name || 'You')
  readonly characterAvatar = computed(() => this.character()?.avatarUrl)
}
