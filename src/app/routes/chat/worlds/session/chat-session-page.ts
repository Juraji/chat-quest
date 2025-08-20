import {Component, computed, effect, ElementRef, inject, Signal, viewChild} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ChatMessage} from '@api/chat-sessions';
import {ChatSessionChatInputBlock, ChatSessionEditBaseSessionBlock, ChatSessionParticipantsBlock} from './components';
import {ChatSessionMessage} from './components/chat-session-message/chat-session-message';
import {ChatSessionData} from './chat-session-data';

@Component({
  selector: 'chat-with-page',
  imports: [
    PageHeader,
    ChatSessionEditBaseSessionBlock,
    ChatSessionParticipantsBlock,
    ChatSessionChatInputBlock,
    ChatSessionMessage,
  ],
  providers: [
    ChatSessionData
  ],
  templateUrl: './chat-session-page.html',
  styleUrls: ["./chat-session-page.scss"]
})
export class ChatSessionPage {
  private readonly sessionData = inject(ChatSessionData)

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
