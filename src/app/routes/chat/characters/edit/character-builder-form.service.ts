import {inject, Injectable} from '@angular/core';
import {CharacterBuilderRequest} from '@api/characters';
import {formControl, formGroup} from '@util/ng';
import {Validators} from '@angular/forms';
import {CharacterEditFormService} from './character-edit-form.service';

interface CharacterBuilderRequestForm extends Omit<CharacterBuilderRequest, 'character'> {
  ignoreCurrentDescriptions: boolean
}

@Injectable()
export class CharacterBuilderFormService {
  private readonly characterFormService = inject(CharacterEditFormService)

  readonly builderFormGroup = formGroup<CharacterBuilderRequestForm>({
    description: formControl('', [Validators.required]),
    worldId: formControl(null),
    instructionId: formControl(0, [Validators.required]),
    llmModelId: formControl(0, [Validators.required]),
    ignoreCurrentDescriptions: formControl(true),
  })

  runFormCheck() {
    this.characterFormService.formGroup.markAllAsDirty({emitEvent: true})
    this.builderFormGroup.markAllAsDirty({emitEvent: true})
  }
}
