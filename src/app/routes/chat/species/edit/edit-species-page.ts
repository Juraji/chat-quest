import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Notifications} from '@components/notifications';
import {Species, SpeciesS} from '@api/species';
import {
  booleanSignal,
  BooleanSignal,
  controlValueSignal,
  formControl,
  formGroup,
  readOnlyControl,
  routeDataSignal
} from '@util/ng';
import {isNew} from '@api/common';
import {FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {AvatarControl} from '@components/avatar-control';
import {RenderedMessage} from '@components/rendered-message';
import {TokenCount} from '@components/token-count';

@Component({
  selector: 'edit-species-page',
  imports: [
    FormsModule,
    ReactiveFormsModule,
    PageHeader,
    AvatarControl,
    RenderedMessage,
    TokenCount
  ],
  templateUrl: './edit-species-page.html',
  styleUrl: './edit-species-page.scss'
})
export class EditSpeciesPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly notifications = inject(Notifications);
  private readonly speciesS = inject(SpeciesS)
  private readonly router = inject(Router);

  readonly species: Signal<Species> = routeDataSignal(this.activatedRoute, 'species')
  readonly isNew: Signal<boolean> = computed(() => isNew(this.species()))
  readonly name: Signal<string> = computed(() => this.species().name)
  readonly avatar: Signal<Nullable<string>> = computed(() => this.species().avatarUrl)

  readonly formGroup = formGroup<Species>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    description: formControl('', [Validators.required]),
    avatarUrl: formControl<Nullable<string>>(null),
  })

  readonly editDescription: BooleanSignal = booleanSignal(false)
  readonly descriptionValue: Signal<string> = controlValueSignal(this.formGroup, 'description')

  constructor() {
    effect(() => {
      const input = this.species()
      this.formGroup.reset(input)
      this.editDescription.set(isNew(input))
    });
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value
    const update: Species = {
      ...this.species(),
      ...formValue,
    }

    this.speciesS
      .save(update)
      .subscribe(s => {
        this.notifications.toast("Species saved!")
        this.router.navigate(['..', s.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onResetForm() {
    this.formGroup.reset(this.species());
  }

  onDeleteSpecies() {
    const s = this.species()
    if (isNew(s)) return
    const doDelete = confirm(`Are you sure you want to delete this species?`)

    if (doDelete) {
      this.speciesS
        .delete(s!.id)
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
