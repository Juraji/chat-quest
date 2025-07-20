import {Component, computed, inject, signal, Signal, WritableSignal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {CharacterCard} from '@components/character-card/character-card';
import {Character} from '@db/characters';
import {Tag, Tags} from '@db/tags';
import {toSignal} from '@angular/core/rxjs-interop';
import {PageHeader} from '@components/page-header/page-header';

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
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly tags = inject(Tags);

  readonly characters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters');

  readonly availableTags: Signal<Tag[]> = toSignal(this.tags.getAll(), {initialValue: []})
  readonly selectedTag: WritableSignal<Tag | null> = signal(null)

  readonly filteredCharacters = computed(() => {
    const characters = this.characters()
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
