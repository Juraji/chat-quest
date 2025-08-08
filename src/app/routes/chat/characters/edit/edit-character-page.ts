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
import {TagsControl} from '@components/tags-control/tags-control';
import {CharacterEditFormService} from './character-edit-form.service';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {Character, CharacterDetails, Characters} from '@api/characters';
import {Tag} from '@api/tags';

@Component({
  selector: 'app-edit-character-page',
  imports: [
    PageHeader,
    ReactiveFormsModule,
    AvatarControl,
    TagsControl,
    RouterOutlet,
    RouterLink,
    RouterLinkActive
  ],
  providers: [
    CharacterEditFormService
  ],
  templateUrl: './edit-character-page.html',
})
export class EditCharacterPage {
  private readonly router = inject(Router)
  private readonly notifications = inject(Notifications)
  private readonly characters = inject(Characters)
  protected readonly activatedRoute = inject(ActivatedRoute)

  private readonly formService = inject(CharacterEditFormService)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')
  readonly characterDetails: Signal<CharacterDetails> = routeDataSignal(this.activatedRoute, 'characterDetails')
  readonly tags: Signal<Tag[]> = routeDataSignal(this.activatedRoute, 'tags')
  readonly dialogueExamples: Signal<string[]> = routeDataSignal(this.activatedRoute, 'dialogueExamples')
  readonly greetings: Signal<string[]> = routeDataSignal(this.activatedRoute, 'greetings')
  readonly groupGreetings: Signal<string[]> = routeDataSignal(this.activatedRoute, 'groupGreetings')
  readonly characterFormData: Signal<CharacterFormData> = computed(() => ({
    character: this.character(),
    characterDetails: this.characterDetails(),
    tags: this.tags(),
    dialogueExamples: this.dialogueExamples(),
    greetings: this.greetings(),
    groupGreetings: this.groupGreetings(),
  }))

  readonly isNew = computed(() => isNew(this.character()))
  readonly name = computed(() => this.character().name)
  readonly favorite = computed(() => this.character().favorite)

  readonly formGroup = this.formService.formGroup

  readonly subMenuItems = [
    {route: 'chat-settings', label: 'Chat Settings'},
    {route: 'descriptions', label: 'Descriptions'},
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
    const characterDetails = this.characterDetails()

    const characterFG = this.formService.characterFG;
    const characterDetailsFG = this.formService.characterDetailsFG;
    const tagsCtrl = this.formService.tagsCtrl;
    const dialogueExamplesFA = this.formService.dialogueExamplesFA;
    const greetingsFA = this.formService.greetingsFA;
    const groupGreetingsFA = this.formService.groupGreetingsFA;

    this.characters
      .save({...character, ...characterFG.value})
      .pipe(
        tap(res => characterFG.reset(res)),
        mergeMap(c => forkJoin({
          character: [c],
          characterDetails: defer(() => !isNew && characterDetailsFG.dirty
            ? this.characters
              .saveDetails({...characterDetails, ...characterDetailsFG.value, characterId: c.id})
              .pipe(tap(res => characterDetailsFG.reset(res)))
            : [null]
          ),
          tags: defer(() => !isNew && tagsCtrl.dirty
            ? this.characters
              .saveTags(c.id, tagsCtrl.value.map(t => t.id))
              .pipe(tap(() => tagsCtrl.reset()))
            : [null]
          ),
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
          groupGreetings: defer(() => !isNew && groupGreetingsFA.dirty
            ? this.characters
              .saveGroupGreetings(c.id, groupGreetingsFA.value)
              .pipe(tap(() => groupGreetingsFA.reset()))
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
