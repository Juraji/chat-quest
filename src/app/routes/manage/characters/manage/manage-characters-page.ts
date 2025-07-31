import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {CharacterWithTags, Tag} from '@api/model';
import {PageHeader} from '@components/page-header';
import {CharacterCard} from '@components/character-card/character-card';
import {Tags} from '@api/clients';

@Component({
  selector: 'app-manage-characters',
  imports: [
    PageHeader,
    RouterLink,
    CharacterCard
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
