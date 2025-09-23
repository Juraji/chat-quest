import {Component, computed, effect, inject, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ConnectionProfilesOverview} from './components/connection-profiles';
import {InstructionOverview} from './components/instructions';
import {ActivatedRoute} from '@angular/router';
import {CQPreferences, Preferences, PreferencesUpdated} from '@api/preferences';
import {Notifications} from '@components/notifications';
import {SSE} from '@api/sse';
import {formControl, formGroup, routeDataSignal, routeQueryParamSignal} from '@util/ng';
import {LlmModelView} from '@api/providers';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Instruction} from '@api/instructions';
import {LlmLabelPipe} from '@components/llm-label.pipe';

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

  private readonly instructions: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'instructions');
  private readonly llmModelViews: Signal<LlmModelView[]> = routeDataSignal(this.activatedRoute, 'llmModelViews');

  private readonly _prefs: Signal<CQPreferences> = routeDataSignal(this.activatedRoute, 'preferences');
  readonly prefs: WritableSignal<CQPreferences> = linkedSignal(() => this._prefs());

  private readonly validate: Signal<boolean> = routeQueryParamSignal(this.activatedRoute, 'validate', v => !!v);

  readonly chatInstructionTemplates: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'CHAT'));
  readonly memoryInstructionTemplates: Signal<Instruction[]> =
    computed(() => this.instructions().filter(i => i.type === 'MEMORIES'));
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
    memoryIncludeChatNotes: formControl(false)
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
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const update: CQPreferences = {
      ...this.prefs(),
      ...this.formGroup.value
    }

    this.preferences
      .save(update)
      .subscribe(() => this.notifications.toast("Preferences updated"))
  }
}
