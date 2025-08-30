import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
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

@Component({
  selector: 'chat-session-details-block',
  imports: [
    DatePipe,
    ReactiveFormsModule
  ],
  templateUrl: './chat-session-details-block.html',
})
export class ChatSessionDetailsBlock {
  private readonly sessionData = inject(ChatSessionData)
  private readonly chatSessions = inject(ChatSessions)
  private readonly preferences = inject(Preferences)
  private readonly notifications = inject(Notifications)

  readonly worldName: Signal<string> = computed(() => this.sessionData.world().name)
  readonly createdAt: Signal<Nullable<string>> = computed(() => this.sessionData.chatSession().createdAt)
  readonly scenarios: Signal<Scenario[]> = this.sessionData.scenarios
  readonly llModels: Signal<LlmModelView[]> = this.sessionData.llmModels

  readonly nameControl: TypedFormControl<string> = formControl('', [Validators.required])
  readonly enableMemoriesControl: TypedFormControl<boolean> = formControl(false)
  readonly pauseAutomaticResponsesControl: TypedFormControl<boolean> = formControl(false)
  readonly scenarioControl: TypedFormControl<Nullable<number>> = formControl(null)
  readonly chatModelControl: TypedFormControl<Nullable<number>> = formControl(null, [Validators.required])

  constructor() {
    effect(() => {
      const session = this.sessionData.chatSession();
      this.nameControl.reset(session.name, {emitEvent: false});
      this.enableMemoriesControl.reset(session.enableMemories, {emitEvent: false})
      this.scenarioControl.reset(session.scenarioId, {emitEvent: false})
    });
    effect(() => {
      const {chatModelId} = this.sessionData.preferences()
      this.chatModelControl.reset(chatModelId, {emitEvent: false})
    });

    this.nameControl.valueChanges
      .pipe(takeUntilDestroyed(), debounceTime(1000))
      .subscribe(name => this.updateSession({name}))
    this.enableMemoriesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(enableMemories => this.updateSession({enableMemories}))
    this.pauseAutomaticResponsesControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(pauseAutomaticResponses => this.updateSession({pauseAutomaticResponses}))
    this.scenarioControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(scenarioId => this.updateSession({scenarioId}))
    this.chatModelControl.valueChanges
      .pipe(takeUntilDestroyed())
      .subscribe(chatModelId => this.updateChatPrefs({chatModelId}))
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
