import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ChatSession, ChatSessions, TimeOfDay} from '@api/chat-sessions';
import {DatePipe} from '@angular/common';
import {formControl, TypedFormControl} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {ChatSessionData} from '../../chat-session-data';
import {Scenario} from '@api/scenarios';
import {LlmModelView} from '@api/providers';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {debounceTime} from 'rxjs';
import {CQPreferences, Preferences} from '@api/preferences';
import {World, Worlds} from '@api/worlds';
import {LlmLabelPipe} from '@components/llm-label.pipe';
import {Instruction} from '@api/instructions';
import {RouterLink} from '@angular/router';

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
  private readonly worlds = inject(Worlds)
  private readonly chatSessions = inject(ChatSessions)
  private readonly preferences = inject(Preferences)
  private readonly notifications = inject(Notifications)

  private readonly llmModels: Signal<LlmModelView[]> = this.sessionData.llmModels
  private readonly instructions: Signal<Instruction[]> = this.sessionData.instructions
  readonly worldName: Signal<string> = computed(() => this.sessionData.world().name)
  readonly createdAt: Signal<Nullable<string>> = computed(() => this.sessionData.chatSession().createdAt)
  readonly scenarios: Signal<Scenario[]> = this.sessionData.scenarios
  readonly characters = this.sessionData.characters

  readonly chatInstructions: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'CHAT'));
  readonly chatModels: Signal<LlmModelView[]> =
    computed(() => this.llmModels().filter(i => i.modelType === 'CHAT_MODEL'))

  readonly nameControl: TypedFormControl<string> = formControl('', [Validators.required])
  readonly generateMemoriesControl: TypedFormControl<boolean> = formControl(false)
  readonly useMemoriesControl: TypedFormControl<boolean> = formControl(false)
  readonly autoArchiveMessagesControl: TypedFormControl<boolean> = formControl(false)
  readonly pauseAutomaticResponsesControl: TypedFormControl<boolean> = formControl(false)
  readonly personaControl: TypedFormControl<Nullable<number>> = formControl(null)
  readonly scenarioControl: TypedFormControl<Nullable<number>> = formControl(null)
  readonly chatModelControl: TypedFormControl<Nullable<number>> = formControl(null, [Validators.required])
  readonly chatInstructionControl: TypedFormControl<Nullable<number>> = formControl(null, [Validators.required])
  readonly timeOfDayControl: TypedFormControl<Nullable<TimeOfDay>> = formControl(null, [Validators.required])
  readonly chatNotesControl: TypedFormControl<Nullable<string>> = formControl(null, [Validators.required])

  readonly selectedModel: Signal<LlmModelView> = computed(() => {
    const {chatModelId} = this.sessionData.preferences()
    return this.llmModels().find(l => l.id === chatModelId)!
  })

  constructor() {
    effect(() => {
      const world = this.sessionData.world()
      this.personaControl.reset(world.personaId, {emitEvent: false})
    });
    effect(() => {
      const session = this.sessionData.chatSession();
      this.nameControl.reset(session.name, {emitEvent: false});
      this.generateMemoriesControl.reset(session.generateMemories, {emitEvent: false})
      this.useMemoriesControl.reset(session.useMemories, {emitEvent: false})
      this.autoArchiveMessagesControl.reset(session.autoArchiveMessages, {emitEvent: false})
      this.pauseAutomaticResponsesControl.reset(session.pauseAutomaticResponses, {emitEvent: false})
      this.scenarioControl.reset(session.scenarioId, {emitEvent: false})
      this.timeOfDayControl.reset(session.currentTimeOfDay, {emitEvent: false})
      this.chatNotesControl.reset(session.chatNotes, {emitEvent: false})

      if (session.generateMemories) {
        this.autoArchiveMessagesControl.disable({emitEvent: false})
      } else {
        this.autoArchiveMessagesControl.enable({emitEvent: false})
      }
    });
    effect(() => {
      const prefs = this.sessionData.preferences()
      this.chatModelControl.reset(prefs.chatModelId, {emitEvent: false})
      this.chatInstructionControl.reset(prefs.chatInstructionId, {emitEvent: false})
    });

    this.nameControl.valueChanges
      .pipe(takeUntilDestroyed(), debounceTime(1000))
      .subscribe(name => this.updateSession({name}))
    this.generateMemoriesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(generateMemories => this.updateSession({generateMemories}))
    this.useMemoriesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(useMemories => this.updateSession({useMemories}))
    this.autoArchiveMessagesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(autoArchiveMessages => this.updateSession({autoArchiveMessages}))
    this.pauseAutomaticResponsesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(pauseAutomaticResponses => this.updateSession({pauseAutomaticResponses}))
    this.personaControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(personaId => this.updateWorld({personaId}))
    this.scenarioControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(scenarioId => this.updateSession({scenarioId}))
    this.chatModelControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(chatModelId => this.updateChatPrefs({chatModelId}))
    this.chatInstructionControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(chatInstructionId => this.updateChatPrefs({chatInstructionId}))
    this.timeOfDayControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(currentTimeOfDay => this.updateSession({currentTimeOfDay}))
    this.chatNotesControl.valueChanges
      .pipe(takeUntilDestroyed(), debounceTime(1000))
      .subscribe(chatNotes => this.updateSession({chatNotes}))
  }

  private updateWorld(w: Partial<World>) {
    const world = this.sessionData.world()

    const update: World = {
      ...world,
      ...w
    }

    this.worlds
      .save(update)
      .subscribe(() => this.notifications.toast("World data updated!"))
  }

  private updateSession(d: Partial<ChatSession>) {
    const session = this.sessionData.chatSession();

    const update: ChatSession = {
      ...session,
      ...d,
    }

    this.chatSessions
      .save(this.sessionData.worldId(), update)
      .subscribe(() => this.notifications.toast("Session details updated!"))
  }

  private updateChatPrefs(p: Partial<CQPreferences>) {
    const prefs = this.sessionData.preferences()

    const update: CQPreferences = {
      ...prefs,
      ...p
    }

    this.preferences
      .save(update)
      .subscribe(() => this.notifications.toast("Preferences updated!"))
  }
}
