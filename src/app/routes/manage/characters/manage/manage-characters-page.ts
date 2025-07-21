import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, Router, RouterLink} from '@angular/router';
import {CharacterCard} from '@components/character-card/character-card';
import {Character, Characters} from '@db/characters';
import {Tag, Tags} from '@db/tags';
import {toSignal} from '@angular/core/rxjs-interop';
import {PageHeader} from '@components/page-header/page-header';
import {Notifications} from '@components/notifications';
import {CharaCardParser} from '@util/chara-card';

@Component({
  selector: 'app-manage-characters-page',
  imports: [
    CharacterCard,
    RouterLink,
    PageHeader
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss']
})
export class ManageCharactersPage {
  private readonly tags = inject(Tags);
  private readonly notifications = inject(Notifications)
  private readonly parser = inject(CharaCardParser)
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly characters = inject(Characters);

  readonly acceptCharaCardTypes: string = this.parser.supportedContentTypes.join(',')

  readonly availableCharacters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters');

  readonly availableTags: Signal<Tag[]> = toSignal(this.tags.getAll(), {initialValue: []})
  readonly selectedTag: WritableSignal<Tag | null> = signal(null)

  readonly filteredCharacters = computed(() => {
    const characters = this.availableCharacters()
    const selectedTag = this.selectedTag()
    if (!!selectedTag) {
      return characters.filter((char) => char.tagIds.includes(selectedTag.id))
    } else {
      return characters
    }
  })

  onToggleSelectedTag(tag: Tag | null) {
    this.selectedTag.update(current => {
      if (tag == null) return null
      return current?.id == tag.id ? null : tag;
    })
  }

  async onCharaCardFileSelected(e: Event) {
    e.preventDefault();

    const input = e.target as HTMLInputElement;
    const file = input.files?.item(0)

    if (!file) return;

    const fileType = file.type;

    if (!(fileType.startsWith('image/') || fileType === 'application/json')) {
      this.notifications.toast('Please select a valid image or json file.', 'DANGER')
      return
    }

    input.value = ''

    try {
      const {character, tags: charTags} = await this.parser.parse(file)
      const doImport = confirm(`Do you want to import "${character.name}" from "${file.name}"?`)
      if (doImport) {
        // TODO: Save and assign tags

        this.characters
          .save(character)
          .subscribe(character => {
            this.notifications.toast(`Character card "${character.name}" imported!`)
            this.router.navigate([character.id], {relativeTo: this.activatedRoute})
          })
      }
    } catch (error) {
      this.notifications.toast(`Error parsing card: ${error}`, 'DANGER')
    }
  }
}
