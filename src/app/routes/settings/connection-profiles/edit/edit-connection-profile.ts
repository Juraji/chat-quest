import {Component, computed, effect, inject, Signal} from '@angular/core';
import {Notifications} from '@components/notifications';
import {ConnectionProfiles} from '@api/clients';
import {ActivatedRoute, Router} from '@angular/router';
import {PageHeader} from '@components/page-header';
import {formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {ConnectionProfile, isNew, LlmModel, ProviderType} from '@api/model';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {Dropdown} from '@components/dropdown';

@Component({
  selector: 'app-edit-connection-profile',
  imports: [
    PageHeader,
    FormsModule,
    ReactiveFormsModule,
    Dropdown
  ],
  templateUrl: './edit-connection-profile.html'
})
export class EditConnectionProfile {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications);
  private readonly connectionProfiles = inject(ConnectionProfiles)

  readonly templates: Signal<ConnectionProfile[]> = routeDataSignal(this.activatedRoute, 'templates')
  readonly profile: Signal<ConnectionProfile> = routeDataSignal(this.activatedRoute, 'profile')
  readonly models: Signal<LlmModel[]> = routeDataSignal(this.activatedRoute, 'models')

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
      if (!!inputP) {
        this.formGroup.reset(inputP)
      }
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value

    const update: ConnectionProfile = {
      ...this.profile(),
      ...formValue,
    }

    this.connectionProfiles
      .save(update)
      .subscribe(profile => {
        this.notifications.toast("Connection Profile saved!")
        this.router.navigate(['..', profile.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()}
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
      this.connectionProfiles
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
}
