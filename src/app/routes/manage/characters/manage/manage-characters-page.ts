import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {CharacterCard} from '@components/character-card/character-card';
import {Character} from '@db/characters';
import {Tag, Tags} from '@db/tags';
import {toSignal} from '@angular/core/rxjs-interop';
import {PageHeader} from '@components/page-header/page-header';
import {CharacterImportButton} from './components/character-import-button/character-import-button';

@Component({
  selector: 'app-manage-characters-page',
  imports: [
    CharacterCard,
    RouterLink,
    PageHeader,
    CharacterImportButton
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss']
})
export class ManageCharactersPage {
  private readonly tags = inject(Tags);
  private readonly activatedRoute = inject(ActivatedRoute);

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
}
