import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {RouterLink} from '@angular/router';
import {PageHeader} from '@components/page-header';
import {CharacterCard} from '@components/cards/character-card';
import {NewItemCard} from '@components/cards/new-item-card';
import {CharacterListView, Characters, Tag, Tags} from '@api/characters';
import {Scalable} from '@components/scalable/scalable';

@Component({
  selector: 'app-manage-characters',
  imports: [
    PageHeader,
    RouterLink,
    CharacterCard,
    NewItemCard,
    Scalable
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss'],
})
export class ManageCharactersPage {
  private readonly tags = inject(Tags)
  private readonly charactersService = inject(Characters)

  readonly allTags: Signal<Tag[]> = this.tags.all
  readonly characters: Signal<CharacterListView[]> = this.charactersService.all

  readonly selectedTag: WritableSignal<Tag | null> = signal(null)
  readonly filteredCharacters: Signal<CharacterListView[]> = computed(() => {
    const selectedTag = this.selectedTag()
    const characters = this.characters()

    if (!!selectedTag) {
      return characters
        .filter(c => c.tags.some(t => t.id === selectedTag.id))
    } else {
      return characters
    }
  })
}
