import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {DatePipe} from '@angular/common';
import {formControl, formGroup, readOnlyControl} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {ChatSessionData} from '../../chat-session-data';
import {Scenario} from '@api/scenarios';
import {LlmModelView} from '@api/providers';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {debounceTime, filter, map, pairwise} from 'rxjs';
import {LlmLabelPipe} from '@components/llm-label.pipe';
import {Instruction} from '@api/instructions';
import {RouterLink} from '@angular/router';

const TEXT_FIELD_DEBOUNCE_TIME = 2500 // ms

@Component({
  selector: 'chat-session-details-block',
  imports: [
    DatePipe,
    ReactiveFormsModule,
    LlmLabelPipe,
    RouterLink
  ],
  templateUrl: './chat-session-details-block.html',
})
export class ChatSessionDetailsBlock {
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions)
  private readonly notifications = inject(Notifications)

  private readonly llmModels: Signal<LlmModelView[]> = this.sessionData.llmModels
  private readonly instructions: Signal<Instruction[]> = this.sessionData.instructions

  readonly session = this.sessionData.chatSession
  readonly worldName: Signal<string> = computed(() => this.sessionData.world().name)
  readonly createdAt: Signal<Nullable<string>> = computed(() => this.session().createdAt)
  readonly scenarios: Signal<Scenario[]> = this.sessionData.scenarios
  readonly characters = this.sessionData.characters

  readonly chatInstructions: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'CHAT'));
  readonly chatModels: Signal<LlmModelView[]> =
    computed(() => this.llmModels().filter(i => i.modelType === 'CHAT_MODEL'))

  readonly sessionForm = formGroup<ChatSession>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    createdAt: readOnlyControl(),
    name: formControl('', [Validators.required]),
    scenarioId: formControl(null),
    generateMemories: formControl(false),
    useMemories: formControl(false),
    autoArchiveMessages: formControl(false),
    pauseAutomaticResponses: formControl(false),
    currentTimeOfDay: formControl(null),
    chatNotes: formControl(null),
    personaId: formControl(null),
    chatModelId: formControl(null, [Validators.required]),
    chatInstructionId: formControl(null, [Validators.required]),
  })

  readonly selectedModel: Signal<LlmModelView> = computed(() => {
    const {chatModelId} = this.session()
    return this.llmModels().find(l => l.id === chatModelId)!
  })

  constructor() {
    effect(() => {
      const session = this.session();
      this.sessionForm.reset(session);

      const autoArchiveMessagesField = this.sessionForm.get("autoArchiveMessages")!;
      session.generateMemories ? autoArchiveMessagesField.disable() : autoArchiveMessagesField.enable();
    });

    // Common fields with standard debounce
    const standardFields = [
      'generateMemories',
      'useMemories',
      'autoArchiveMessages',
      'pauseAutomaticResponses',
      'personaId',
      'scenarioId',
      'chatModelId',
      'chatInstructionId',
      'currentTimeOfDay'
    ];

    // Fields with debounce time
    const debouncedFields = [
      'name',
      'chatNotes',
    ];

    standardFields.forEach(field => {
      this.sessionForm.get(field)!.valueChanges
        .pipe(
          takeUntilDestroyed(),
          pairwise(), filter(([p, n]) => p !== undefined && p !== n), map(n => n[1])
        )
        .subscribe(value => this.updateSession({[field]: value}));
    });

    debouncedFields.forEach(field => {
      this.sessionForm.get(field)!.valueChanges
        .pipe(
          takeUntilDestroyed(),
          pairwise(), filter(([p, n]) => p !== undefined && p !== n), map(n => n[1]),
          debounceTime(TEXT_FIELD_DEBOUNCE_TIME)
        )
        .subscribe(value => this.updateSession({[field]: value}));
    });
  }

  private updateSession(d: Partial<ChatSession>) {
    const session = this.session();

    const update: ChatSession = {
      ...session,
      ...d,
    }

    this.chatSessions
      .save(this.sessionData.worldId(), update)
      .subscribe(() => this.notifications.toast("Session details updated!"))
  }

}
