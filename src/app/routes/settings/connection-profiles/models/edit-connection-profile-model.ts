import {Component, computed, effect, inject, input, InputSignal, Signal} from '@angular/core';
import {LlmModel} from '@api/model';
import {BooleanSignal, booleanSignal, formControl, formGroup, readOnlyControl, toControlValueSignal} from '@util/ng';
import {AbstractControl, FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {ConnectionProfiles} from '@api/clients';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
  selector: 'app-edit-connection-profile-model',
  imports: [
    FormsModule,
    ReactiveFormsModule
  ],
  templateUrl: './edit-connection-profile-model.html',
  styleUrls: ['./edit-connection-profile-model.scss'],
})
export class EditConnectionProfileModel {
  private readonly notifications = inject(Notifications)
  private readonly connectionProfiles = inject(ConnectionProfiles)
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly model: InputSignal<LlmModel> = input.required()
  readonly editMode: BooleanSignal = booleanSignal(false)

  readonly formGroup = formGroup<LlmModel>({
    id: readOnlyControl(),
    profileId: readOnlyControl(),
    modelId: readOnlyControl(),
    temperature: formControl(0, [Validators.required, Validators.min(0.01)]),
    maxTokens: formControl(1, [Validators.required, Validators.min(1)]),
    topP: formControl(0, [Validators.required, Validators.min(0), Validators.max(1)]),
    stream: formControl(false),
    stopSequences: formControl(''),
    disabled: formControl(false),
  })

  readonly modelDisabled: Signal<boolean> = computed(() => this.model().disabled)
  readonly formModelDisabled: Signal<boolean> = toControlValueSignal(this.formGroup, 'disabled');

  constructor() {
    effect(() => {
      const model = this.model();
      const edit = this.editMode();
      if (!edit) {
        this.formGroup.reset(model)
      }
    });

    effect(() => {
      const state = this.formModelDisabled();
      const controls: AbstractControl[] = [
        this.formGroup.get('temperature')!,
        this.formGroup.get('maxTokens')!,
        this.formGroup.get('topP')!,
        this.formGroup.get('stream')!,
        this.formGroup.get('stopSequences')!,
      ]

      if (state) {
        controls.forEach(control => control.disable())
      } else {
        controls.forEach(control => control.enable())
      }
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const model = this.model();
    const value = this.formGroup.value;

    const update: LlmModel = {
      ...model,
      ...value,
    }

    this.connectionProfiles
      .saveModel(update)
      .subscribe(() => {
        this.editMode.set(false)
        this.notifications.toast("Model settings saved!")
        this.router.navigate([], {
          queryParams: {u: Date.now()},
          replaceUrl: true,
          relativeTo: this.activatedRoute
        });
      })
  }
}
