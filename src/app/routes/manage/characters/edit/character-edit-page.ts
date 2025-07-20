import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {formArray, formControl, formGroup, readOnlyControl, routeDataSignal} from '@util/ng';
import {Character, Characters} from '@db/characters';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {
  CharacterEditBaseInformation
} from './components/character-edit-base-information/character-edit-base-information';
import {
  CharacterEditExtendedDetails
} from './components/character-edit-extended-details/character-edit-extended-details';
import {Notifications} from '@components/notifications';
import {PageHeader} from '@components/page-header/page-header';
import {CharacterEditChatDefaults} from './components/character-edit-chat-defaults/character-edit-chat-defaults';

@Component({
  selector: 'app-character-edit-page',
  imports: [
    ReactiveFormsModule,
    CharacterEditBaseInformation,
    CharacterEditExtendedDetails,
    PageHeader,
    CharacterEditChatDefaults,
  ],
  templateUrl: './character-edit-page.html',
})
export class CharacterEditPage {
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')
  readonly isNew: Signal<boolean> = computed(() => !this.character().id)
  readonly name: Signal<string> = computed(() => this.character().name)
  readonly isFavorite: Signal<boolean> = computed(() => this.character().favorite)

  readonly formGroup = formGroup<Character>({
    id: readOnlyControl(),
    // Base props
    name: formControl('', [Validators.required]),
    appearance: formControl(''),
    personality: formControl(''),
    avatar: formControl(null),
    tagIds: formControl([]),
    favorite: formControl(false),

    // Extended props
    history: formControl(''),
    likelyActions: formArray([]),
    unlikelyActions: formArray([]),
    dialogueExamples: formArray([]),

    // Chat Defaults
    scenario: formControl(''),
    firstMessage: formControl(''),
    alternateGreetings: formArray([]),
    groupGreetings: formArray([]),
    groupTalkativeness: formControl(1.0),
  })

  constructor() {
    effect(() => {
      const character = this.character()
      this.formGroup.reset(character)
    });
  }

  onSubmit() {
    if (this.formGroup.invalid) return

    const formValue: Character = this.formGroup.value

    const update: Character = {
      ...this.character(),
      ...formValue,
    }

    this.characters
      .save(update)
      .subscribe(char => {
        this.notifications.toast(`Changes to ${char.name} saved.`)
        this.router.navigate(['..', char.id], {
          relativeTo: this.activatedRoute,
          queryParams: {u: Date.now()},
          replaceUrl: true
        });
      })
  }

  onRevertChanges() {
    const character = this.character()
    this.formGroup.reset(character)
    this.notifications.toast(`Changes to ${character.name} reverted.`)
  }

  onDeleteCharacter() {
    if (this.isNew()) return

    const character = this.character()
    const doDelete = confirm(`Are you sure you want to delete ${character.name}? This action cannot be undone.`)

    if (doDelete) {
      this.characters
        .delete(character.id)
        .subscribe(() => {
          this.notifications.toast(`Deleted ${character.name}.`)
          this.router.navigate(['../..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          });
        })
    }
  }
}
