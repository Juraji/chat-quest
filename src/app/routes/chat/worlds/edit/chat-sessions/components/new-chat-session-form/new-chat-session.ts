import {Component, effect, inject, signal, Signal, WritableSignal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Scenario} from '@api/scenarios';
import {
  BooleanSignal,
  booleanSignal,
  controlValueSignal,
  formControl,
  formGroup,
  readOnlyControl,
  routeDataSignal
} from '@util/ng';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {ChatSession} from '@api/chat-sessions';
import {Character} from '@api/characters';
import {CharacterCard} from '@components/cards/character-card';

interface NewChatSessionForm extends ChatSession {
  hasCharacters: boolean;
}

@Component({
  selector: 'new-chat-session-form',
  imports: [
    ReactiveFormsModule,
    CharacterCard
  ],
  templateUrl: './new-chat-session.html',
  styleUrls: ['./new-chat-session.scss'],
})
export class NewChatSession {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios');
  readonly characters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters');

  readonly formGroup = formGroup<NewChatSessionForm>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    createdAt: readOnlyControl(),
    name: formControl('', [Validators.required]),
    scenarioId: formControl<Nullable<number>>(null),
    enableMemories: formControl(true),
    hasCharacters: formControl(false, [Validators.requiredTrue]),
  })

  readonly useCustomName: BooleanSignal = booleanSignal(false)
  readonly selectedScenarioId: Signal<Nullable<number>> = controlValueSignal(this.formGroup, 'scenarioId');
  readonly selectedCharacterIds: WritableSignal<number[]> = signal([])

  constructor() {
    effect(() => {
      const nameCtrl = this.formGroup.get('name')
      if (this.useCustomName()) nameCtrl?.enable(); else nameCtrl?.disable()
    });
    effect(() => {
      const hasCharactersSelected = this.selectedCharacterIds().length > 0;
      const ctrl = this.formGroup.get('hasCharacters')!

      ctrl.setValue(hasCharactersSelected)
      if (hasCharactersSelected) ctrl.markAsDirty()
    });
    effect(() => {
      if (!this.useCustomName()) {
        const scenarioId = this.selectedScenarioId();
        const characterIds = this.selectedCharacterIds()

        const nameCtrl = this.formGroup.get('name')!
        const scenario = this.scenarios()
          .find(s => s.id === scenarioId)
          ?.name
        const characters = this.characters()
            .filter(c => characterIds.includes(c.id))
            .map(c => c.name)
            .join(' and ')
          || 'no one'

        let name: string
        if (!!scenario) {
          name = `${scenario} with ${characters}`;
        } else {
          name = characters;
        }

        nameCtrl.setValue(name)
      }
    });
  }

  onToggleCharacter(c: Character) {
    this.selectedCharacterIds.update(ids =>
      ids.includes(c.id)
        ? ids.filter(id => id !== c.id)
        : [...ids, c.id])
  }

  onNewChatSession() {
    if (this.formGroup.invalid) return

    const session = this.formGroup.getRawValue()
    const characterIds = this.selectedCharacterIds()

    this.router.navigate(
      [{outlets: {primary: ['chat', 'worlds', 1, 'chat', 'new']}}],
      {
        queryParams: {
          with: characterIds,
          sessionName: session.name,
          scenarioId: session.scenarioId,
          enableMemories: session.enableMemories
        }
      }
    )
  }
}
