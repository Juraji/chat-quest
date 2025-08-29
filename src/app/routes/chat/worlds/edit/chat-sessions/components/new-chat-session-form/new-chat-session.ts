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
import {Character, characterSortingTransformer} from '@api/characters';
import {CharacterCard} from '@components/cards/character-card';
import {ChatSession} from '@api/chat-sessions';
import {arrayAdd, arrayRemove} from '@util/array';
import {Scalable} from '@components/scalable/scalable';

@Component({
  selector: 'new-chat-session-form',
  imports: [
    ReactiveFormsModule,
    CharacterCard,
    Scalable
  ],
  templateUrl: './new-chat-session.html',
})
export class NewChatSession {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios');
  readonly characters: Signal<Character[]> =
    routeDataSignal<Character[]>(this.activatedRoute, 'characters', characterSortingTransformer);

  readonly formGroup = formGroup<ChatSession>({
    id: readOnlyControl(),
    worldId: readOnlyControl(),
    createdAt: readOnlyControl(),
    name: formControl('', [Validators.required]),
    scenarioId: formControl<Nullable<number>>(null),
    enableMemories: formControl(true),
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
      if (hasCharactersSelected) this.formGroup.markAllAsDirty()
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
          || 'No one'

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
        ? arrayRemove(ids, id => id === c.id)
        : arrayAdd(ids, c.id))
  }

  onNewChatSession() {
    if (this.formGroup.invalid) return

    const session = this.formGroup.getRawValue()
    const characterIds = this.selectedCharacterIds()

    this.router.navigate(
      [{outlets: {primary: ['chat', 'worlds', 1, 'session', 'new']}}],
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
