import {Component, inject, Signal} from '@angular/core';
import {CharacterEditFormService} from '../character-edit-form.service';
import {ReactiveFormsModule} from '@angular/forms';
import {booleanSignal, BooleanSignal, controlValueSignal} from '@util/ng';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {RenderedMessage} from '@components/rendered-message';
import {TokenCount} from '@components/token-count';

@Component({
  selector: 'app-character-edit-descriptions',
  imports: [
    ReactiveFormsModule,
    RenderedMessage,
    TokenCount
  ],
  templateUrl: './character-edit-descriptions.html',
})
export class CharacterEditDescriptions {
  private readonly formService = inject(CharacterEditFormService)

  readonly formGroup = this.formService.characterFG
  readonly onFormSubmit = this.formService.requestSubmitFn()

  readonly editAppearance: BooleanSignal = booleanSignal(false)
  readonly appearanceValue: Signal<string> = controlValueSignal(this.formGroup, 'appearance')

  readonly editPersonality: BooleanSignal = booleanSignal(false)
  readonly personalityValue: Signal<string> = controlValueSignal(this.formGroup, 'personality');

  readonly editHistory: BooleanSignal = booleanSignal(false)
  readonly historyValue: Signal<string> = controlValueSignal(this.formGroup, 'history');

  constructor() {
    this.formService.onFormReset
      .pipe(takeUntilDestroyed())
      .subscribe(() => {
        this.editAppearance.set(false)
        this.editPersonality.set(false)
        this.editHistory.set(false)
      })
  }
}
