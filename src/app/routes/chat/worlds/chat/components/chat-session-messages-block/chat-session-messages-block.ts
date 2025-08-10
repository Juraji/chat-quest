import {
  Component,
  effect,
  ElementRef,
  inject,
  input,
  InputSignal,
  linkedSignal,
  Signal,
  WritableSignal
} from '@angular/core';
import {
  ChatMessage,
  ChatMessageCreated,
  ChatMessageDeleted,
  ChatMessageUpdated,
  sessionEntityFilter
} from '@api/chat-sessions';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';
import {SSE} from '@api/sse';
import {arrayAddItem, arrayRemoveItem, arrayUpsertItem,} from '@util/array';
import {ChatSessionMessage} from '../chat-session-message/chat-session-message';

@Component({
  selector: 'chat-session-messages-block',
  imports: [
    ChatSessionMessage
  ],
  templateUrl: './chat-session-messages-block.html',
})
export class ChatSessionMessagesBlock {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly elementRef = inject(ElementRef);
  private readonly sse = inject(SSE)

  readonly chatSessionId: InputSignal<number> = input.required()

  private readonly initialMessages: Signal<ChatMessage[]> = routeDataSignal(this.activatedRoute, 'messages')
  readonly messages: WritableSignal<ChatMessage[]> = linkedSignal(() => this.initialMessages())

  constructor() {
    effect(() => {
      this.messages()
      setTimeout(() => {
        // Scroll to bottom of message list when messages updates.
        // This needs to be in a timeout, so the page updates before we scroll.
        const element = this.elementRef.nativeElement;
        element.scrollTop = element.scrollHeight;
      }, 0)
    });

    this.sse
      .on(ChatMessageCreated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayAddItem(prev, message)))
    this.sse
      .on(ChatMessageUpdated, sessionEntityFilter(this.chatSessionId))
      .subscribe(message => this.messages
        .update(prev => arrayUpsertItem(prev, message, ({id}) => id === message.id)))
    this.sse
      .on(ChatMessageDeleted)
      .subscribe(messageId => this.messages
        .update(prev => arrayRemoveItem(prev, ({id}) => id === messageId)))
  }
}
