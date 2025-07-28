import {Component, computed, effect, inject, Signal} from '@angular/core';
import {
  formArray,
  formControl,
  formGroup,
  readOnlyControl,
  routeDataSignal,
  TypedFormArray,
  TypedFormControl,
  TypedFormGroup
} from '@util/ng';
import {Character, CharacterDetails, isNew, Tag} from '@api/model';
import {ActivatedRoute, Router} from '@angular/router';
import {Notifications} from '@components/notifications';
import {Characters} from '@api/clients';
import {FormArray, ReactiveFormsModule, Validators} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {AvatarControl} from './components/avatar-control';

interface CharacterMasterForm {
  character: Character
  characterDetails: CharacterDetails
  tags: Tag[]
  dialogueExamples: string[]
  greetings: string[]
  groupGreetings: string[]
}

@Component({
  selector: 'app-edit-character-page',
  imports: [
    PageHeader,
    ReactiveFormsModule,
    AvatarControl
  ],
  templateUrl: './edit-character-page.html',
})
export class EditCharacterPage {
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')

  readonly isNew = computed(() => isNew(this.character()))
  readonly name = computed(() => this.character().name)
  readonly favorite = computed(() => this.character().favorite)

  readonly characterDetails: Signal<CharacterDetails> = routeDataSignal(this.activatedRoute, 'characterDetails')
  readonly tags: Signal<Tag[]> = routeDataSignal(this.activatedRoute, 'tags')
  readonly dialogueExamples: Signal<string[]> = routeDataSignal(this.activatedRoute, 'dialogueExamples')
  readonly greetings: Signal<string[]> = routeDataSignal(this.activatedRoute, 'greetings')
  readonly groupGreetings: Signal<string[]> = routeDataSignal(this.activatedRoute, 'groupGreetings')

  readonly formGroup = formGroup<CharacterMasterForm>({
    character: formGroup({
      id: readOnlyControl(0),
      createdAt: readOnlyControl<Nullable<number>>(null),
      name: formControl('', [Validators.required]),
      favorite: formControl(false),
      avatarUrl: formControl<Nullable<string>>(null)
    }),
    characterDetails: formGroup({
      characterId: readOnlyControl(0),
      appearance: formControl<Nullable<string>>(null),
      personality: formControl<Nullable<string>>(null),
      history: formControl<Nullable<string>>(null),
      scenario: formControl<Nullable<string>>(null),
      groupTalkativeness: formControl(0)
    }),
    tags: formControl([]),
    dialogueExamples: formArray([]),
    greetings: formArray([]),
    groupGreetings: formArray([]),
  })

  readonly characterFG: TypedFormGroup<Character> =
    this.formGroup.get('character') as TypedFormGroup<Character>
  readonly characterDetailsFG: TypedFormGroup<CharacterDetails> =
    this.formGroup.get('characterDetails') as TypedFormGroup<CharacterDetails>
  readonly tagsCtrl: TypedFormControl<Tag[]> =
    this.formGroup.get('tags') as TypedFormControl<Tag[]>
  readonly dialogueExamplesFA: TypedFormArray<string> =
    this.formGroup.get('dialogueExamples') as FormArray
  readonly greetingsFA: TypedFormArray<string> =
    this.formGroup.get('greetings') as FormArray
  readonly groupGreetingsFA: TypedFormArray<string> =
    this.formGroup.get('groupGreetings') as FormArray

  constructor() {
    effect(() => {
      const character = this.character()
      this.characterFG.reset(character)
    });
    effect(() => {
      const details = this.characterDetails()
      this.characterDetailsFG.reset(details)
    });
    effect(() => {
      const tags = this.tags()
      this.tagsCtrl.reset(tags)
    });
    effect(() => {
      const examples = this.dialogueExamples()
      this.setControlsTo(this.dialogueExamplesFA, examples)
    })
    effect(() => {
      const greetings = this.greetings()
      this.setControlsTo(this.greetingsFA, greetings)
    })
    effect(() => {
      const greetings = this.groupGreetings()
      this.setControlsTo(this.groupGreetingsFA, greetings)
    })
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
    throw new Error('Not implemented.')
  }

  onRevertChanges() {
    const character = this.character()
    this.characterFG.reset(character)

    const details = this.characterDetails()
    this.characterDetailsFG.reset(details)

    const tags = this.tags()
    this.tagsCtrl.reset(tags)

    const examples = this.dialogueExamples()
    this.setControlsTo(this.dialogueExamplesFA, examples)

    const greetings = this.greetings()
    this.setControlsTo(this.greetingsFA, greetings)

    const groupGreetings = this.groupGreetings()
    this.setControlsTo(this.groupGreetingsFA, groupGreetings)
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
          this.router.navigate(['..'], {
            relativeTo: this.activatedRoute,
            replaceUrl: true
          });
        })
    }
  }

  private setControlsTo(arr: TypedFormArray<string>, values: string[]) {
    arr.clear()
    values.forEach(value => this.addControlTo(arr, value))
    arr.reset()
  }

  private addControlTo(arr: TypedFormArray<string>, value: string = '') {
    arr.push(formControl(value, [Validators.required]))
  }
}
