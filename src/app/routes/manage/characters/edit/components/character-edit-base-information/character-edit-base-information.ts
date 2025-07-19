import {Component, effect, inject, input, InputSignal} from '@angular/core';
import {Character, Characters} from '@db/characters';
import {formControl, formGroup} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {TagsControl} from '@components/tags-control/tags-control';
import {ActivatedRoute, Router} from '@angular/router';
import {AvatarControl} from '@components/avatar-control/avatar-control';

type CharacterEditBaseForm = Pick<Character, 'name' | 'appearance' | 'personality' | 'favorite' | 'tagIds' | 'avatar'>

@Component({
  selector: 'app-character-edit-base-information',
  imports: [
    ReactiveFormsModule,
    TagsControl,
    AvatarControl
  ],
  templateUrl: './character-edit-base-information.html'
})
export class CharacterEditBaseInformation {
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)

  readonly character: InputSignal<Character> = input.required()

  readonly formGroup = formGroup<CharacterEditBaseForm>({
    name: formControl('', [Validators.required]),
    appearance: formControl(''),
    personality: formControl(''),
    favorite: formControl(false),
    tagIds: formControl([]),
    avatar: formControl(null),
  })

  constructor() {
    effect(() => {
      const character = this.character()
      this.formGroup.reset(character)
    });
  }

  onSubmit() {
    if (this.formGroup.invalid) return

    const formValue: CharacterEditBaseForm = this.formGroup.getRawValue()

    const update: Character = {
      ...this.character(),
      ...formValue,
    }

    this.characters
      .save(update)
      .subscribe(char => this.router.navigate(["..", char.id], {
        relativeTo: this.activatedRoute,
        queryParams: {u: Date.now()}
      }))
  }
}
