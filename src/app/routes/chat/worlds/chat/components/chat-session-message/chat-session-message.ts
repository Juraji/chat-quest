import {Component, computed, inject, input, InputSignal, Signal} from '@angular/core';
import {ChatMessage} from '@api/chat-sessions';
import {RenderedMessage} from '@components/rendered-message';
import {Character} from '@api/characters';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'chat-session-message',
  imports: [
    RenderedMessage
  ],
  templateUrl: './chat-session-message.html',
  styleUrl: './chat-session-message.scss'
})
export class ChatSessionMessage {
  private readonly activatedRoute = inject(ActivatedRoute);

  private readonly allCharacters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'allCharacters')

  readonly message: InputSignal<ChatMessage> = input.required()

  readonly content: Signal<string> = computed(() => this.message().content)
  readonly isUser: Signal<boolean> = computed(() => this.message().isUser)
  readonly createdAt: Signal<string> = computed(() => this.message().createdAt!)
  readonly memoryId: Signal<Nullable<number>> = computed(() => this.message().memoryId)
  readonly character: Signal<Nullable<Character>> = computed(() => {
    const cId = this.message().characterId
    const allCharacters = this.allCharacters();
    return !!cId ? allCharacters.find(c => c.id === cId) : null;
  })
}
