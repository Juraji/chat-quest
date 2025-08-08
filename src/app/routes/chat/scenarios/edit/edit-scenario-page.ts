import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Notifications} from '@components/notifications';
import {
  BooleanSignal,
  booleanSignal,
  formControl,
  formGroup,
  readOnlyControl,
  routeDataSignal,
  toControlValueSignal
} from '@util/ng';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {AvatarControl} from '@components/avatar-control';
import {RenderedMessage} from '@components/rendered-message';
import {TokenCount} from '@components/token-count';
import {Scenario, Scenarios} from '@api/scenarios';
import {isNew} from '@api/common';

@Component({
  selector: 'app-edit-scenario',
  imports: [
    FormsModule,
    PageHeader,
    ReactiveFormsModule,
    AvatarControl,
    RenderedMessage,
    TokenCount
  ],
  templateUrl: './edit-scenario-page.html'
})
export class EditScenarioPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly notifications = inject(Notifications);
  private readonly scenarios = inject(Scenarios)
  private readonly router = inject(Router);

  readonly scenario: Signal<Scenario> = routeDataSignal(this.activatedRoute, 'scenario');
  readonly isNew: Signal<boolean> = computed(() => isNew(this.scenario()))
  readonly name: Signal<string> = computed(() => this.scenario().name)

  readonly formGroup = formGroup<Scenario>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    description: formControl('', [Validators.required]),
    avatarUrl: formControl<Nullable<string>>(null),
    linkedCharacterId: formControl<Nullable<number>>(null)
  })

  readonly editDescription: BooleanSignal = booleanSignal(false)
  readonly descriptionValue: Signal<string> = toControlValueSignal(this.formGroup, 'description')

  constructor() {
    effect(() => {
      const input = this.scenario()
      this.formGroup.reset(input)
      this.editDescription.set(isNew(input))
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value

    const update: Scenario = {
      ...this.scenario(),
      ...formValue
    }

    this.scenarios
      .save(update)
      .subscribe(scenario => {
        this.notifications.toast("Scenario saved!")
        this.router.navigate(['..', scenario.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onResetForm() {
    this.formGroup.reset(this.scenario());
  }

  onDeleteScenario() {
    const scene = this.scenario();
    if (isNew(scene)) return
    const doDelete = confirm(`Are you sure you want to delete this scenario?`)

    if (doDelete) {
      this.scenarios
        .delete(scene!.id)
        .subscribe(() => {
          this.notifications.toast("Scenario deleted!")
          this.router.navigate(['..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          })
        })
    }
  }
}
