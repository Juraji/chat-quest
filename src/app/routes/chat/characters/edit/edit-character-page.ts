import {Component, computed, effect, inject, Signal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {isNew} from '@api/common';
import {ActivatedRoute, Router, RouterLink, RouterLinkActive, RouterOutlet} from '@angular/router';
import {Notifications} from '@components/notifications';
import {ReactiveFormsModule} from '@angular/forms';
import {PageHeader} from '@components/page-header';
import {AvatarControl} from '@components/avatar-control';
import {CharacterFormData} from './character-form-data';
import {defer, forkJoin, mergeMap, tap} from 'rxjs';
import {CharacterEditFormService} from './character-edit-form.service';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {Character, Characters} from '@api/characters';

@Component({
  selector: 'app-edit-character-page',
  imports: [
    PageHeader,
    ReactiveFormsModule,
    AvatarControl,
    RouterOutlet,
    RouterLink,
    RouterLinkActive
  ],
  providers: [
    CharacterEditFormService
  ],
  styleUrls: ['./edit-character-page.scss'],
  templateUrl: './edit-character-page.html',
})
export class EditCharacterPage {
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characters = inject(Characters)
  protected readonly activatedRoute = inject(ActivatedRoute)

  private readonly formService = inject(CharacterEditFormService)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')
  readonly dialogueExamples: Signal<string[]> = routeDataSignal(this.activatedRoute, 'dialogueExamples')
  readonly greetings: Signal<string[]> = routeDataSignal(this.activatedRoute, 'greetings')
  readonly characterFormData: Signal<CharacterFormData> = computed(() => ({
    character: this.character(),
    dialogueExamples: this.dialogueExamples(),
    greetings: this.greetings(),
  }))

  readonly isNew = computed(() => isNew(this.character()))
  readonly name = computed(() => this.character().name)
  readonly favorite = computed(() => this.character().favorite)
  readonly avatar = computed(() => this.character().avatarUrl)

  readonly formGroup = this.formService.formGroup

  readonly subMenuItems = [
    {route: 'descriptions', label: 'Descriptions'},
    {route: 'chat-settings', label: 'Chat Settings'},
    {route: 'memories', label: 'Memories'},
  ]

  constructor() {
    effect(() => {
      const formData = this.characterFormData()
      this.formService.resetFormData(formData)
    });

    this.formService.onSubmitRequested
      .pipe(takeUntilDestroyed())
      .subscribe(() => this.onSubmit())
  }

  onSubmit() {
    if (this.formGroup.invalid) return

    const isNew = this.isNew()
    const character = this.character()

    const characterFG = this.formService.characterFG;
    const dialogueExamplesFA = this.formService.dialogueExamplesFA;
    const greetingsFA = this.formService.greetingsFA;

    this.characters
      .save({...character, ...characterFG.value})
      .pipe(
        tap(res => characterFG.reset(res)),
        mergeMap(c => forkJoin({
          character: [c],
          dialogueExamples: defer(() => !isNew && dialogueExamplesFA.dirty
            ? this.characters
              .saveDialogueExamples(c.id, dialogueExamplesFA.value)
              .pipe(tap(() => dialogueExamplesFA.reset()))
            : [null]
          ),
          greetings: defer(() => !isNew && greetingsFA.dirty
            ? this.characters
              .saveGreetings(c.id, greetingsFA.value)
              .pipe(tap(() => greetingsFA.reset()))
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
        const currentSubRoute = this.router.url
          .replace(/.*\/([a-z-]+)\?.*$/, "$1");
        this.router.navigate(["..", res.character.id, currentSubRoute], {
          queryParams: {u: Date.now()},
          replaceUrl: true,
          relativeTo: this.activatedRoute
        })
      })
  }

  onResetForm() {
    this.formService.resetFormData(this.characterFormData())
  }

  onDuplicateCharacter() {
    if (this.isNew()) return
    const characterId = this.character().id

    this.characters
      .duplicate(characterId)
      .subscribe(res => {
        this.notifications.toast(`Character duplicated successfully!`)
        this.router.navigate(['..', res.id], {
          relativeTo: this.activatedRoute,
        })
      })
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
}
