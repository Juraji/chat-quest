import {Component, computed, effect, inject, Signal} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {formArray, formControl, formGroup, readOnlyControl, routeDataSignal, TypedFormArray} from '@util/ng';
import {Character, Characters} from '@db/characters';
import {ReactiveFormsModule, Validators} from '@angular/forms';
import {Notifications} from '@components/notifications';
import {PageHeader} from '@components/page-header/page-header';
import {CharacterImportExport} from '@util/import-export';
import {Tags} from '@db/tags';
import {firstValueFrom} from 'rxjs';
import {downloadBlob} from '@util/blobs';
import {AvatarControl} from '@components/avatar-control';
import {TagsControl} from '@components/tags-control/tags-control';

@Component({
  selector: 'app-character-edit-page',
  imports: [
    ReactiveFormsModule,
    PageHeader,
    AvatarControl,
    TagsControl,
  ],
  templateUrl: './character-edit-page.html',
})
export class CharacterEditPage {
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characterExport = inject(CharacterImportExport)
  private readonly tags = inject(Tags)

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

    // Chat Behavior
    firstMessage: formControl(''),
    alternateGreetings: formArray([]),

    // Group Chat Behavior
    groupGreetings: formArray([]),
    groupTalkativeness: formControl(0.5),

    // Extended props
    history: formControl(''),
    likelyActions: formArray([]),
    unlikelyActions: formArray([]),
    dialogueExamples: formArray([]),

    // Scenario
    scenario: formControl(''),
  })

  readonly alternateGreetingsFA: TypedFormArray<string> =
    this.formGroup.get('alternateGreetings') as TypedFormArray<string>
  readonly groupGreetingsFA: TypedFormArray<string> =
    this.formGroup.get('groupGreetings') as TypedFormArray<string>
  readonly likelyActionsFA: TypedFormArray<string> =
    this.formGroup.get('likelyActions') as TypedFormArray<string>
  readonly unlikelyActionsFA: TypedFormArray<string> =
    this.formGroup.get('unlikelyActions') as TypedFormArray<string>
  readonly dialogueExamplesFA: TypedFormArray<string> =
    this.formGroup.get('dialogueExamples') as TypedFormArray<string>

  readonly dialogExample: string = `{{user}}: Hi, I am User. What's your name?
{{char}}: *looks at {{user}} startled* "Ow, hello. I didn't notice you there. Nice to meet you, my name is {char}."`

  constructor() {
    effect(() => {
      const character = this.character()
      this.formGroup.reset(character)

      this.setControlsTo(this.alternateGreetingsFA, character.alternateGreetings)
      this.setControlsTo(this.groupGreetingsFA, character.groupGreetings)
      this.setControlsTo(this.likelyActionsFA, character.likelyActions)
      this.setControlsTo(this.unlikelyActionsFA, character.unlikelyActions)
      this.setControlsTo(this.dialogueExamplesFA, character.dialogueExamples)
    });
  }

  onAddControl(arr: TypedFormArray<string>, value: string = '') {
    this.addControlTo(arr, value)
    arr.markAsDirty()
  }

  onRemoveControl(arr: TypedFormArray<string>, idx: number) {
    arr.removeAt(idx)
    arr.markAsDirty()
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

  async onExportCharacter() {
    const character = this.character()
    const allTags = await firstValueFrom(this.tags.getAll())
    const characterTags = character.tagIds
      .map(tagId => allTags.find(t => t.id === tagId))
      .filter(t => !!t)
      .map(t => t.label)

    const blob = await this.characterExport.exportToFile(character, characterTags)
    downloadBlob(blob, `ChatQuest_v1_${character.name}.json`)
  }

  onExportScenario() {
    const sceneDescription = this.formGroup.get('scenario')!.value
    if (sceneDescription.trim() === '') return

    this.router.navigate(['/manage/scenarios/new'], {
      queryParams: {sceneDescription}
    })
  }

  private setControlsTo(arr: TypedFormArray<string>, values: string[]) {
    arr.clear()
    values.forEach(value => this.addControlTo(arr, value))
  }

  private addControlTo(arr: TypedFormArray<string>, value: string = '') {
    arr.push(formControl(value, [Validators.required]))
  }
}
