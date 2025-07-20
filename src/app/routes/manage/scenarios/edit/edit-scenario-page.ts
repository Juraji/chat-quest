import {Component, computed, effect, inject, Signal} from '@angular/core';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router} from '@angular/router';
import {Scenario, Scenarios} from '@db/scenarios';
import {Notifications} from '@components/notifications';
import {PageHeader} from '@components/page-header/page-header';

@Component({
  selector: 'app-edit-scenario',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    PageHeader
  ],
  templateUrl: './edit-scenario-page.html'
})
export class EditScenarioPage {
  private readonly scenarios = inject(Scenarios)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)

  readonly scenario: Signal<Scenario> = routeDataSignal(this.activatedRoute, 'scenario');
  readonly isNew = computed(() => !this.scenario().id)
  readonly name = computed(() => this.scenario().name);

  readonly formGroup = formGroup<Scenario>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    sceneDescription: formControl('', [Validators.required])
  });

  constructor() {
    effect(() => {
      const scenario = this.scenario()
      this.formGroup.reset(scenario)
    });
  }

  onSubmit() {
    if (this.formGroup.invalid) return

    const formValue: Scenario = this.formGroup.value

    const update: Scenario = {
      ...this.scenario(),
      ...formValue,
    }

    this.scenarios
      .save(update)
      .subscribe(scenario => {
        this.notifications.toast(`Changes to ${scenario.name} saved.`)
        this.router.navigate(['..', scenario.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        });
      })
  }

  onDeleteScenario() {
    if (this.isNew()) return

    const scenario = this.scenario()
    const doDelete = confirm(`Are you sure you want to delete ${scenario.name}? This action cannot be undone.`)

    if (doDelete) {
      this.scenarios
        .delete(scenario.id)
        .subscribe(() => {
          this.notifications.toast(`Deleted ${scenario.name}.`)
          this.router.navigate(['../..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          });
        })
    }
  }

  onRevertChanges() {
    const character = this.scenario()
    this.formGroup.reset(character)
    this.notifications.toast(`Changes to ${character.name} reverted.`)
  }
}
