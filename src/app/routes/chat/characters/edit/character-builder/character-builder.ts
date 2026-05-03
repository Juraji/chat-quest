import {Component, effect, inject, Signal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {CharacterBuilderRequest, Characters} from '@api/characters';
import {Notifications} from '@components/notifications';
import {ActivatedRoute, Router} from '@angular/router';
import {CharacterEditFormService} from '../character-edit-form.service';
import {ReactiveFormsModule} from '@angular/forms';
import {Instruction} from '@api/instructions';
import {World} from '@api/worlds';
import {LlmModelView} from '@api/providers';
import {LlmLabelPipe} from '@components/llm-label.pipe';
import {CharacterBuilderFormService} from '../character-builder-form.service';

@Component({
  selector: 'character-builder',
  imports: [
    ReactiveFormsModule,
    LlmLabelPipe
  ],
  templateUrl: './character-builder.html',
})
export class CharacterBuilder {
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)

  private readonly characterformService = inject(CharacterEditFormService)
  private readonly builderformService = inject(CharacterBuilderFormService)

  readonly builderFormGroup = this.builderformService.builderFormGroup
  readonly characterFormGroup = this.characterformService.characterFG

  get builderInvalid(): boolean {
    return this.builderFormGroup.invalid || this.characterFormGroup.invalid
  }

  readonly instructions: Signal<Instruction[]> = routeDataSignal<Instruction[]>(
    this.activatedRoute, 'instructions', l => l.filter(i => i.type === "CHARACTER_BUILDER"))
  readonly llmModels: Signal<LlmModelView[]> = routeDataSignal<LlmModelView[]>(
    this.activatedRoute, 'llmModels', l => l.filter(m => m.modelType === 'CHAT_MODEL'))
  readonly worlds: Signal<World[]> = routeDataSignal<World[]>(this.activatedRoute, 'worlds')

  constructor() {
    effect(() => {
      const l = this.instructions()
      const c = this.builderFormGroup.get('instructionId')!!
      if (l.length > 0 && c.value == 0) c.reset(l[0].id)
    });
    effect(() => {
      const l = this.llmModels()
      const c = this.builderFormGroup.get('llmModelId')!!
      if (l.length > 0 && c.value == 0) c.reset(l[0].id)
    });

    this.builderformService.runFormCheck()
  }

  protected onBuilderSubmit() {
    if (this.builderInvalid) return

    const props = this.builderFormGroup.value

    const request: CharacterBuilderRequest = {
      character: this.characterformService.characterFG.value,
      description: props.description,
      worldId: props.worldId,
      instructionId: props.instructionId,
      llmModelId: props.llmModelId,
    }

    if (props.ignoreCurrentDescriptions) {
      request.character.appearance = null
      request.character.personality = null
      request.character.history = null
    }

    this.notifications
      .run("Building character...", "INFO", () => this.characters.buildCharacter(request))
      .subscribe(result => {
        this.characterformService.characterFG.setValue(result)
        this.router.navigate(['../descriptions'], {relativeTo: this.activatedRoute})
      })
  }
}
