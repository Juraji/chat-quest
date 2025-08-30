import {Component, computed, effect, ElementRef, inject, Signal, viewChild} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ChatMessage} from '@api/chat-sessions';
import {ChatSessionChatInputBlock, ChatSessionDetailsBlock, ChatSessionParticipantsBlock} from './components';
import {ChatSessionMessage} from './components/chat-session-message/chat-session-message';
import {ChatSessionData} from './chat-session-data';
import {MemoryList} from '@components/memory-list';

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
  readonly messages: Signal<ChatMessage[]> = this.sessionData.messages

  protected readonly chatMessagesContainerRef: Signal<ElementRef<HTMLDivElement> | undefined> =
    viewChild('chatMessagesContainer', {read: ElementRef})

  constructor() {
    effect(() => {
      this.sessionData.messages()
      const element = this.chatMessagesContainerRef()?.nativeElement;
      if (!!element) requestAnimationFrame(() => element.scrollTop = element.scrollHeight)
    });
  }
}
