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
  styleUrls: ["./chat-session-page.scss"]
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
