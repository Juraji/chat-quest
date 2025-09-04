import {Component, computed, effect, inject, linkedSignal, Signal, WritableSignal} from '@angular/core';
import {Notifications} from '@components/notifications';
import {ActivatedRoute, Router} from '@angular/router';
import {PageHeader} from '@components/page-header';
import {formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {DropdownContainer, DropdownMenu, DropdownToggle} from '@components/dropdown';
import {AiProviders, ConnectionProfile, LlmModel, LlmModelType, Providers, ProviderType} from '@api/providers';
import {isNew} from '@api/common';
import {arrayReplace} from '@util/array';

@Component({
  selector: 'app-edit-connection-profile',
  imports: [
    PageHeader,
    FormsModule,
    ReactiveFormsModule,
    DropdownContainer,
    DropdownToggle,
    DropdownMenu
  ],
  templateUrl: './edit-connection-profile.html'
})
export class EditConnectionProfile {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications);
  private readonly providers = inject(Providers)

  readonly aiProviders: Signal<AiProviders> = routeDataSignal(this.activatedRoute, 'providers')
  readonly profile: Signal<ConnectionProfile> = routeDataSignal(this.activatedRoute, 'profile')
  private readonly _models: Signal<LlmModel[]> = routeDataSignal(this.activatedRoute, 'models')
  readonly models: WritableSignal<LlmModel[]> = linkedSignal(() => this._models())

  readonly isNew = computed(() => isNew(this.profile()))

  readonly formGroup = formGroup<ConnectionProfile>({
    id: readOnlyControl(0),
    name: formControl('', [Validators.required]),
    providerType: formControl<ProviderType>('OPEN_AI', [Validators.required]),
    baseUrl: formControl('', [Validators.required]),
    apiKey: formControl('', [Validators.required])
  })

  constructor() {
    effect(() => {
      const inputP = this.profile()
      this.formGroup.reset(inputP)
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value

    const update: ConnectionProfile = {
      ...this.profile(),
      ...formValue,
    }

    this.providers
      .save(update)
      .subscribe(profile => {
        this.notifications.toast("Connection Profile saved!")
        this.router.navigate(['..', profile.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onRevertChanges() {
    this.formGroup.reset(this.profile());
  }

  onDeleteCharacter() {
    const p = this.profile()
    if (isNew(p)) return
    const doDelete = confirm(`Are you sure you want to delete this connection?`)

    if (doDelete) {
      this.providers
        .delete(p!.id)
        .subscribe(() => {
          this.notifications.toast("Connection Profile deleted!")
          this.router.navigate(['..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          })
        })
    }
  }

  onCopyFromTemplate(template: ConnectionProfile) {
    const {id, ...patch} = template

    if (!this.isNew()) {
      const doCopy = confirm(`Are you sure you want overwrite the current settings?`)
      if (!doCopy) return
    }

    this.formGroup.patchValue(patch)
    this.formGroup.markAsDirty()
  }

  onRefreshModels() {
    const doRefresh = confirm(
      `Are you sure you want to refresh the available models?

ChatQuest will request the available set of models from the AI Provider.
New models will be added and non-existent models will be removed.

(Unchanged models will remain as they are.)`
    )

    if (!doRefresh) return
    const profileId = this.profile().id

    this.providers
      .refreshModels(profileId)
      .subscribe(() => {
        this.notifications.toast("Models updated!")
        this.router.navigate([], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onModelTypeChange(model: LlmModel, event: Event) {
    const modelType = (event.target as HTMLSelectElement).value as LlmModelType;
    this.providers
      .saveModel({...model, modelType})
      .subscribe(res => {
        this.models.update(models => arrayReplace(models, res, m => m.id === res.id))
        this.notifications.toast("Model updated!");
      })
  }

  onToggleModelDisabled(model: LlmModel) {
    this.providers
      .saveModel({...model, disabled: !model.disabled})
      .subscribe(res => {
        this.models.update(models => arrayReplace(models, res, m => m.id === res.id))
        this.notifications.toast(`Model ${res.disabled ? 'disabled' : 'enabled'}!`);
      })
  }
}
