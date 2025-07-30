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
import {AvatarControl} from '@components/avatar-control';
import {CharacterFormData} from './character-form-data';
import {defer, forkJoin, mergeMap, tap} from 'rxjs';
import {TagsControl} from '@components/tags-control/tags-control';

@Component({
  selector: 'app-edit-character-page',
  imports: [
    PageHeader,
    ReactiveFormsModule,
    AvatarControl,
    TagsControl
  ],
  templateUrl: './edit-character-page.html',
})
export class EditCharacterPage {
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characters = inject(Characters)
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly characterFormData: Signal<CharacterFormData> = routeDataSignal(this.activatedRoute, 'characterFormData')

  readonly character = computed(() => this.characterFormData()?.character)
  readonly isNew = computed(() => isNew(this.character()))
  readonly name = computed(() => this.character().name)
  readonly favorite = computed(() => this.character().favorite)

  readonly formGroup = formGroup<CharacterFormData>({
    character: formGroup({
      id: readOnlyControl(0),
      createdAt: readOnlyControl<Nullable<string>>(null),
      name: formControl('', [Validators.required]),
      favorite: formControl(false),
      avatarUrl: formControl<Nullable<string>>(null)
    }),
    characterDetails: formGroup({
      characterId: readOnlyControl(0),
      appearance: formControl<Nullable<string>>(null),
      personality: formControl<Nullable<string>>(null),
      history: formControl<Nullable<string>>(null),
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
    effect(() => this.onResetForm());
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

    const isNew = this.isNew()
    const {
      character,
      characterDetails,
    } = this.characterFormData()

    this.characters
      .save({...character, ...this.characterFG.value})
      .pipe(
        tap(res => this.characterFG.reset(res)),
        mergeMap(c => forkJoin({
          character: [c],
          characterDetails: defer(() => !isNew && this.characterDetailsFG.dirty
            ? this.characters
              .saveDetails({...characterDetails, ...this.characterDetailsFG.value, characterId: c.id})
              .pipe(tap(res => this.characterDetailsFG.reset(res)))
            : [null]
          ),
          tags: defer(() => !isNew && this.tagsCtrl.dirty
            ? this.characters
              .saveTags(c.id, this.tagsCtrl.value.map(t => t.id))
              .pipe(tap(() => this.tagsCtrl.reset()))
            : [null]
          ),
          dialogueExamples: defer(() => !isNew && this.dialogueExamplesFA.dirty
            ? this.characters
              .saveDialogueExamples(c.id, this.dialogueExamplesFA.value)
              .pipe(tap(() => this.dialogueExamplesFA.reset()))
            : [null]
          ),
          greetings: defer(() => !isNew && this.greetingsFA.dirty
            ? this.characters
              .saveGreetings(c.id, this.greetingsFA.value)
              .pipe(tap(() => this.greetingsFA.reset()))
            : [null]
          ),
          groupGreetings: defer(() => !isNew && this.groupGreetingsFA.dirty
            ? this.characters
              .saveGroupGreetings(c.id, this.groupGreetingsFA.value)
              .pipe(tap(() => this.groupGreetingsFA.reset()))
            : [null]
          ),
        })),
        tap({
          error: () => this.notifications
            .toast("Character save (partially failed).", "DANGER")
        })
      )
      .subscribe(res => {
        this.notifications.toast(`${res.character.name} saved!`)
        this.router.navigate(["..", res.character.id], {
          queryParams: {u: Date.now()},
          replaceUrl: true,
          relativeTo: this.activatedRoute
        })
      })
  }

  onResetForm() {
    const formData = this.characterFormData()
    this.formGroup.reset(formData)

    this.setControlsTo(this.dialogueExamplesFA, formData.dialogueExamples)
    this.setControlsTo(this.greetingsFA, formData.greetings)
    this.setControlsTo(this.groupGreetingsFA, formData.groupGreetings)
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
