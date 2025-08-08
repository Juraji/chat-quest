import {Component, computed, inject, Signal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {ActivatedRoute} from '@angular/router';
import {formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {World} from '@api/worlds';
import {isNew} from '@api/common';

@Component({
  imports: [
    PageHeader,
    ReactiveFormsModule
  ],
  templateUrl: './edit-world-page.html'
})
export class EditWorldPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly world: Signal<World> = routeDataSignal(this.activatedRoute, 'world');

  readonly isNew: Signal<boolean> = computed(() => isNew(this.world()))
  readonly name: Signal<string> = computed(() => this.world().name)

  readonly formGroup = formGroup<World>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    description: formControl(''),
  })

  onResetForm() {

  }

  onFormSubmit() {

  }

  onDeleteScenario() {

  }
}
