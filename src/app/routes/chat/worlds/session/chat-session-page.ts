import {
  Component,
  computed,
  effect,
  ElementRef,
  inject,
  linkedSignal,
  Signal,
  viewChild,
  WritableSignal
} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ChatSessionChatInputBlock, ChatSessionDetailsBlock, ChatSessionParticipantsBlock} from './components';
import {ChatSessionMessage} from './components/chat-session-message/chat-session-message';
import {ChatSessionData} from './chat-session-data';
import {MemoryList} from '@components/memory-list';
import {booleanSignal} from '@util/ng';

@Component({
  selector: 'chat-with-page',
  imports: [
    PageHeader,
    ChatSessionDetailsBlock,
    ChatSessionParticipantsBlock,
    ChatSessionChatInputBlock,
    ChatSessionMessage,
    MemoryList,
  ],
  providers: [
    ChatSessionData
  ],
  templateUrl: './chat-session-page.html',
  styleUrls: ["./chat-session-page.scss"],
  host: {
    '[class.focus-mode]': 'focusMode()'
  }
})
export class ChatSessionPage {
  private readonly sessionData = inject(ChatSessionData)

  readonly worldId: Signal<number> = this.sessionData.worldId
  readonly chatSessionName: Signal<string> = computed(() => this.sessionData.chatSession().name)

  readonly messageLimitPageSize = computed(() => this.sessionData.preferences().memoryTriggerAfter)
  readonly messageLimit: WritableSignal<number> = linkedSignal(() => this.messageLimitPageSize())
  readonly olderMessagesAvailable = computed(() => this.messageLimit() < this.sessionData.messages().length)
  readonly olderMessagesShown = computed(() => this.messageLimit() > this.messageLimitPageSize())
  readonly messages = computed(() => {
    const messages = this.sessionData.messages()
    const limit = this.messageLimit()
    if (limit > messages.length) {
      return messages
    } else {
      return messages.slice(messages.length - limit, messages.length)
    }
  })

  readonly focusMode = booleanSignal(false)
  readonly enableBackdrop = booleanSignal(true)
  readonly avatars: Signal<string[]> = computed(() => {
    if (!this.enableBackdrop()) return [];

    const result: string[] = [];
    const messages = this.sessionData.messages();
    const characters = this.sessionData.characters();
    const avatarUris = new Map<number, string>();

    for (const character of characters) {
      if (character.avatarUrl) {
        avatarUris.set(character.id, character.avatarUrl);
      }
    }

    // Scene Avatar
    const s = this.sessionData.chatSession();
    const sceneId = s.scenarioId;
    if (sceneId) {
      const sceneAvatarUri = this.sessionData
        .scenarios()
        .find(sc => sc.id === sceneId)
        ?.avatarUrl;

      if (sceneAvatarUri) {
        result.push(sceneAvatarUri);
      }
    }

    // Character Avatar - find from end of messages array
    if (messages.length > 0 && avatarUris.size > 0) {
      for (let i = messages.length - 1; i >= 0; i--) {
        const message = messages[i];
        if (message.characterId && avatarUris.has(message.characterId)) {
          result.push(avatarUris.get(message.characterId)!);
          break;
        }
      }
    }

    return result;
  })

  protected readonly chatMessagesContainerRef: Signal<ElementRef<HTMLDivElement> | undefined> =
    viewChild('chatMessagesContainer', {read: ElementRef})

  constructor() {
    effect(() => {
      this.sessionData.messages()
      const element = this.chatMessagesContainerRef()?.nativeElement;
      if (!!element) requestAnimationFrame(() => element.scrollTop = element.scrollHeight)
    });
  }

  onLoadOlderMessages() {
    const incrWith = this.messageLimitPageSize()
    this.messageLimit.update(limit => limit + incrWith)
  }

  onResetMessageLimit() {
    this.messageLimit.set(this.messageLimitPageSize())
  }
}
