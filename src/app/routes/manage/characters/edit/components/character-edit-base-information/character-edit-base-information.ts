import {Component, input, InputSignal} from '@angular/core';
import {Character} from '@db/characters';
import {TypedFormGroup} from '@util/ng';
import {ReactiveFormsModule} from '@angular/forms';
import {TagsControl} from '@components/tags-control/tags-control';
import {AvatarControl} from '@components/avatar-control';


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
  readonly parentForm: InputSignal<TypedFormGroup<Character>> = input.required()
}
