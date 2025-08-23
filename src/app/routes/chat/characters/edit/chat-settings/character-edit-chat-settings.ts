import {Component, computed, inject} from '@angular/core';
import {CharacterEditFormService} from '../character-edit-form.service';
import {ReactiveFormsModule} from '@angular/forms';
import {booleanSignal, controlValueSignal} from '@util/ng';
import {EmptyPipe} from '@components/empty.pipe';
import {RenderedMessage} from '@components/rendered-message';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {TokenCount} from '@components/token-count';

@Component({
  selector: 'app-character-edit-chat-settings',
  imports: [
    ReactiveFormsModule,
    EmptyPipe,
    RenderedMessage,
    TokenCount,
  ],
  templateUrl: './character-edit-chat-settings.html',
  styleUrls: ['./character-edit-chat-settings.scss']
})
export class CharacterEditChatSettings {
  private readonly formService = inject(CharacterEditFormService)

  readonly formGroup = this.formService.formGroup

  readonly editDialogueExamples = booleanSignal(false)
  readonly dialogueExamplesFA = this.formService.dialogueExamplesFA
  readonly dialogueExamples = controlValueSignal(this.dialogueExamplesFA)

  readonly editGreetings = booleanSignal(false)
  readonly greetingsFA = this.formService.greetingsFA
  readonly greetings = controlValueSignal(this.greetingsFA)

  readonly editGroupGreetings = booleanSignal(false)
  readonly groupGreetingsFA = this.formService.groupGreetingsFA
  readonly groupGreetings = controlValueSignal(this.groupGreetingsFA)
  readonly groupTalkativeness = controlValueSignal<number>(this.formGroup, ['character', 'groupTalkativeness'])
  readonly groupTalkativenessPercent = computed(() => (this.groupTalkativeness() * 100) + '%')

  readonly onFormSubmit = this.formService.requestSubmitFn()
  readonly onAddControl = this.formService.addControlFn()
  readonly onRemoveControl = this.formService.removeControlFn()

  constructor() {
    this.formService.onFormReset
      .pipe(takeUntilDestroyed())
      .subscribe(() => {
        this.editDialogueExamples.set(false)
        this.editGreetings.set(false)
        this.editGroupGreetings.set(false)
      })

  }
}
