import {Component, computed, effect, inject, input, InputSignal, Signal} from '@angular/core';
import {ChatMessage, ChatSessions} from '@api/chat-sessions';
import {RenderedMessage} from '@components/rendered-message';
import {BaseCharacter} from '@api/characters';
import {Notifications} from '@components/notifications';
import {ChatSessionData} from '../../chat-session-data';
import {booleanSignal, BooleanSignal, formControl, formGroup} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Memories} from '@api/memories';
import {mergeMap} from 'rxjs';
import {ActivatedRoute, Router} from '@angular/router';
import {System} from '@api/system';

type MessageFormGroup = Pick<ChatMessage, 'content'>

@Component({
  selector: 'chat-session-message',
  imports: [
    RenderedMessage,
    ReactiveFormsModule
  ],
  templateUrl: './chat-session-message.html',
  styleUrl: './chat-session-message.scss',
  host: {
    '[class.is-archived]': 'isArchived()',
    '[class.is-generating]': 'isGenerating()'
  }
})
export class ChatSessionMessage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions)
  private readonly memories = inject(Memories)
  private readonly notifications = inject(Notifications)
  private readonly system = inject(System)

  readonly worldId: Signal<number> = this.sessionData.worldId

  readonly message: InputSignal<ChatMessage> = input.required()

  readonly content: Signal<string> = computed(() => this.message().content)
  readonly isUser: Signal<boolean> = computed(() => this.message().isUser)
  readonly isGenerating: Signal<boolean> = computed(() => this.message().isGenerating)
  readonly isArchived: Signal<boolean> = computed(() => this.message().isArchived)
  readonly createdAt: Signal<string> = computed(() => this.message().createdAt!)

  readonly character: Signal<Nullable<BaseCharacter>> = computed(() => {
    const characterId = this.message().characterId
    const characters = this.sessionData.characters()
    if (!characterId) return null
    else return characters.find(({id}) => id === characterId)
  })
  readonly characterName = computed(() => this.character()?.name || 'You')
  readonly characterAvatar = computed(() => this.character()?.avatarUrl)

  readonly editMessage: BooleanSignal = booleanSignal(false)
  readonly editMessageFormGroup = formGroup<MessageFormGroup>({
    content: formControl('', [Validators.required]),
  })

  constructor() {
    effect(() => {
      const content = this.message().content
      this.editMessageFormGroup.reset({content})
    });
  }

  onDeleteMessage() {
    const doDelete = confirm(`Are you sure you want to delete this message?
Note that this and all subsequent messages will be deleted and this action can not be undone.`);

    if (doDelete) {
      const message = this.message()

      this.chatSessions
        .deleteMessage(this.worldId(), message.chatSessionId, message.id)
        .subscribe(() => this.notifications.toast('Message deleted!'))
    }
  }

  onEditMessageSubmit() {
    if (this.editMessageFormGroup.invalid) return

    const worldId = this.worldId();
    const sessionId = this.sessionData.chatSessionId()
    const msg = this.message()
    const value: MessageFormGroup = this.editMessageFormGroup.value

    const update: ChatMessage = {
      ...msg,
      ...value
    }

    this.chatSessions
      .saveMessage(worldId, sessionId, update)
      .subscribe(() => {
        this.editMessage.set(false)
        this.notifications.toast('Message updated!');
      })
  }

  onGenerateMemory() {
    const worldId = this.worldId();
    const {id, content} = this.message();
    const contentPreview = content.substring(0, 50)

    this.notifications.toast(`Requesting memory generation for "<span class="text-info">${contentPreview}</span>"...`);
    this.memories
      .generateMemoriesForMessage(worldId, id)
      .subscribe(() => this.notifications.toast('Memory generation completed!'))
  }

  onRegenerateMessage() {
    const doRegen = confirm(`Are you sure you want to regenerate this message? This will delete all messages after this one.`)
    if (!doRegen) return
    const worldId = this.worldId()
    const {id, chatSessionId, characterId} = this.message();

    this.chatSessions
      .deleteMessage(worldId, chatSessionId, id)
      .pipe(mergeMap(() => this.chatSessions
        .triggerParticipantResponse(worldId, chatSessionId, characterId!)))
      .subscribe()
  }

  onForkChat() {
    const worldId = this.worldId()
    const {id, chatSessionId} = this.message();

    this.chatSessions
      .forkChatSession(worldId, chatSessionId, id)
      .subscribe(session => {
        this.notifications.toast('Chat session forked!');
        this.router.navigate(
          ['..', session.id],
          {relativeTo: this.activatedRoute}
        );
      });
  }

  onCancelGeneration() {
    this.system
      .stopCurrentGeneration()
      .subscribe(() => this.notifications.toast('Generation cancelled!'));
  }
}
