import {Component, inject, Signal} from '@angular/core';
import {ChatMessage, ChatSessions} from '@api/chat-sessions';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {controlValueSignal, formControl, formGroup} from '@util/ng';
import {TokenCount} from '@components/token-count';
import {NEW_ID} from '@api/common';
import {ChatSessionData} from '../../chat-session-data';

interface ChatInputForm {
  message: string;
}

@Component({
  selector: 'chat-session-chat-input-block',
  imports: [
    ReactiveFormsModule,
    TokenCount
  ],
  templateUrl: './chat-session-chat-input-block.html',
})
export class ChatSessionChatInputBlock {
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions);

  readonly worldId: Signal<number> = this.sessionData.worldId
  readonly chatSessionId: Signal<number> = this.sessionData.chatSessionId

  readonly formGroup = formGroup<ChatInputForm>({
    message: formControl('', [Validators.required])
  })

  readonly messageValue: Signal<string> = controlValueSignal(this.formGroup, 'message')

  onSendMessage() {
    if (this.formGroup.invalid) return

    const worldId = this.worldId();
    const chatSessionId = this.chatSessionId();
    const content = this.formGroup.value.message

    const chatMessage: ChatMessage = {
      id: NEW_ID,
      chatSessionId,
      createdAt: null,
      isUser: true,
      isSystem: false,
      isGenerating: false,
      characterId: null,
      content,
    }

    this.chatSessions
      .saveMessage(worldId, chatSessionId, chatMessage)
      .subscribe(() => this.formGroup.reset({message: ''}))
  }

  onInputKeyDown(e: KeyboardEvent) {
    if (e.ctrlKey && (e.code === 'Enter' || e.code === 'NumpadEnter')) {
      e.preventDefault();
      this.onSendMessage()
    }
  }
}
