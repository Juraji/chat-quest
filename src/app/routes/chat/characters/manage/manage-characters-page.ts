import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {PageHeader} from '@components/page-header';
import {CharacterCard} from '@components/cards/character-card';
import {NewItemCard} from '@components/cards/new-item-card';
import {Tag, Tags} from '@api/tags';
import {CharacterWithTags} from '@api/characters';

@Component({
  selector: 'app-manage-characters',
  imports: [
    PageHeader,
    RouterLink,
    CharacterCard,
    NewItemCard
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss'],
})
export class ManageCharactersPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly tags = inject(Tags)

  readonly allTags: Signal<Tag[]> = this.tags.cachedTags
  readonly characters: Signal<CharacterWithTags[]> = routeDataSignal(this.activatedRoute, 'characters');

  readonly selectedTag: WritableSignal<Tag | null> = signal(null)
  readonly filteredCharacters: Signal<CharacterWithTags[]> = computed(() => {
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
