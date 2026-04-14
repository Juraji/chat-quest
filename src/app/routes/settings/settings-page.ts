import {Component, computed, effect, inject, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ConnectionProfilesOverview} from './components/connection-profiles';
import {InstructionOverview} from './components/instructions';
import {ActivatedRoute} from '@angular/router';
import {CQPreferences, Preferences, PreferencesUpdated} from '@api/preferences';
import {Notifications} from '@components/notifications';
import {SSE} from '@api/sse';
import {booleanSignal, BooleanSignal, formControl, formGroup, routeDataSignal, routeQueryParamSignal} from '@util/ng';
import {LlmModelView} from '@api/providers';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Instruction} from '@api/instructions';
import {LlmLabelPipe} from '@components/llm-label.pipe';
import {ChatQuestUIConfig} from '@config/config';

@Component({
  selector: 'app-settings-page',
  imports: [
    PageHeader,
    ConnectionProfilesOverview,
    InstructionOverview,
    ReactiveFormsModule,
    LlmLabelPipe,
  ],
  templateUrl: './settings-page.html'
})
export class SettingsPage {
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly preferences = inject(Preferences)
  private readonly notifications = inject(Notifications)
  private readonly sse = inject(SSE)
  private readonly config = inject(ChatQuestUIConfig)

  private readonly instructions: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'instructions');
  private readonly llmModelViews: Signal<LlmModelView[]> = routeDataSignal(this.activatedRoute, 'llmModelViews');

  private readonly _prefs: Signal<CQPreferences> = routeDataSignal(this.activatedRoute, 'preferences');
  readonly prefs: WritableSignal<CQPreferences> = linkedSignal(() => this._prefs());

  private readonly validate: Signal<boolean> = routeQueryParamSignal(this.activatedRoute, 'validate', v => !!v);

  readonly uiSettingsAdvancedOpened: BooleanSignal = booleanSignal(false)

  readonly chatInstructionTemplates: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'CHAT'));
  readonly memoryInstructionTemplates: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'MEMORIES'));
  readonly titleGenerationInstructionTemplates: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'TITLE_GENERATION'));
  readonly chatModels: Signal<LlmModelView[]> =
    computed(() => this.llmModelViews().filter(i => i.modelType === 'CHAT_MODEL'))
  readonly embeddingModels: Signal<LlmModelView[]> =
    computed(() => this.llmModelViews().filter(i => i.modelType === 'EMBEDDING_MODEL'))

  readonly formGroup = formGroup<CQPreferences>({
    chatModelId: formControl<Nullable<number>>(null, [Validators.required]),
    chatInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
    maxMessagesInContext: formControl(0, [Validators.required, Validators.min(1)]),
    embeddingModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoriesModelId: formControl<Nullable<number>>(null, [Validators.required]),
    memoriesInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
    memoryMinP: formControl(0, [Validators.required, Validators.min(0.01), Validators.max(1.0)]),
    memoryTriggerAfter: formControl(0, [Validators.required, Validators.min(1)]),
    memoryWindowSize: formControl(0, [Validators.required, Validators.min(1)]),
    memoryIncludeChatSize: formControl(0, [Validators.required, Validators.min(1)]),
    memoryIncludeChatNotes: formControl(false),
    titleGenerationModelId: formControl<Nullable<number>>(null, [Validators.required]),
    titleGenerationInstructionId: formControl<Nullable<number>>(null, [Validators.required]),
    titleGenerationMessageWindow: formControl<number>(0, [Validators.required, Validators.min(1)]),
  })

  readonly uiSettingsFormGroup = formGroup<ChatQuestUIConfig>({
    apiBaseUrl: formControl('', [Validators.required]),
    sseMaxReconnectAttempts: formControl(0, [Validators.required, Validators.min(1)]),
    sseMinReconnectDelayMillis: formControl(0, [Validators.required, Validators.min(100)]),
    maxMessagesInChatView: formControl(0, [Validators.required, Validators.min(1)]),
  })

  constructor() {
    effect(() => this.onRevertChanges());
    effect(() => {
      if (this.validate()) {
        this.formGroup.markAsDirty()
        this.notifications.toast("There is an issue with your settings. Please check.", "WARNING")
      }
    });

    this.sse
      .on(PreferencesUpdated)
      .subscribe(prefs => this.prefs.set(prefs));
  }

  onRevertChanges() {
    const prefs = this.prefs()
    this.formGroup.reset(prefs)
    this.uiSettingsFormGroup.reset(this.config)
  }

  onFormSubmit() {
    if (this.formGroup.invalid || this.uiSettingsFormGroup.invalid) return

    if (this.uiSettingsFormGroup.dirty) {
      Object.assign(this.config, this.uiSettingsFormGroup.value)
    }

    const update: CQPreferences = {
      ...this.prefs(),
      ...this.formGroup.value
    }

    this.preferences
      .save(update)
      .subscribe(() => this.notifications.toast("Preferences updated"))
  }
}
