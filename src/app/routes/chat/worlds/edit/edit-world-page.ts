import {Component, computed, effect, inject, Signal} from '@angular/core';
import {PageHeader} from '@components/page-header';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {ActivatedRoute, Router, RouterLink, RouterLinkActive, RouterOutlet} from '@angular/router';
import {controlValueSignal, formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {World, Worlds} from '@api/worlds';
import {isNew} from '@api/common';
import {AvatarControl} from '@components/avatar-control';
import {TokenCount} from '@components/token-count';
import {Notifications} from '@components/notifications';
import {Characters} from '@api/characters';

@Component({
  imports: [
    PageHeader,
    ReactiveFormsModule,
    AvatarControl,
    RouterLinkActive,
    RouterOutlet,
    RouterLink,
    TokenCount
  ],
  templateUrl: './edit-world-page.html'
})
export class EditWorldPage {
  readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly worlds = inject(Worlds)
  private readonly characters = inject(Characters)
  private readonly notifications = inject(Notifications)

  readonly world: Signal<World> = routeDataSignal(this.activatedRoute, 'world');
  readonly allCharacters = this.characters.all

  readonly isNew: Signal<boolean> = computed(() => isNew(this.world()))
  readonly name: Signal<string> = computed(() => this.world().name)

  readonly formGroup = formGroup<World>({
    id: readOnlyControl(),
    name: formControl('', [Validators.required]),
    description: formControl<Nullable<string>>(null),
    avatarUrl: formControl<Nullable<string>>(null),
    personaId: formControl(null)
  })

  readonly descriptionValue: Signal<string> = controlValueSignal(this.formGroup, 'description')

  readonly subMenuItems = [
    {route: 'chat-sessions', label: 'Chat Sessions'},
    {route: 'memories', label: 'Memories'},
  ]


  constructor() {
    effect(() => this.onResetForm());
  }

  onResetForm() {
    this.formGroup.reset(this.world())
  }

  onFormSubmit() {
    if (this.formGroup.invalid) return

    const formValue = this.formGroup.value

    const update: World = {
      ...this.world(),
      ...formValue
    }

    this.worlds
      .save(update)
      .subscribe(world => {
        this.notifications.toast("World saved!")
        this.router.navigate(['..', world.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        })
      })
  }

  onDeleteWorld() {
    const world = this.world();
    if (isNew(world)) return
    const doDelete = confirm(`Are you sure you want to delete this world?\n\nAll chat sessions and memories will also be deleted!`)

    if (doDelete) {
      this.worlds
        .delete(world!.id)
        .subscribe(() => {
          this.notifications.toast("World deleted!")
          this.router.navigate(['..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          })
        })
    }
  }
}
