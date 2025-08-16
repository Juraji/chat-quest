import {Component, computed, effect, inject, input, InputSignal, signal, Signal, WritableSignal} from '@angular/core';
import {TagsControl} from '@components/tags-control';
import {BaseCharacter, CharacterListView, Characters} from '@api/characters';
import {Tag} from '@api/tags';

@Component({
  selector: 'app-character-card',
  imports: [
    TagsControl
  ],
  templateUrl: './character-card.html',
  styleUrl: './character-card.scss',
  host: {
    '[class.item-card]': 'true',
    '[class.favorite]': 'favorite()'
  }
})
export class CharacterCard {
  readonly characters = inject(Characters)

  readonly character: InputSignal<BaseCharacter> = input.required()
  protected readonly id: Signal<number> = computed(() => this.character().id)
  protected readonly name: Signal<string> = computed(() => this.character().name)
  protected readonly favorite: Signal<boolean> = computed(() => this.character().favorite)
  protected readonly tags: WritableSignal<Tag[]> = signal([])
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => {
    const u = this.character().avatarUrl
    return !!u ? `url(${u})` : null;
  })

  constructor() {
    effect(() => {
      const char = this.character()
      if ('tags' in char) {
        this.tags.set((char as CharacterListView).tags)
      } else {
        this.characters
          .getTags(char.id)
          .subscribe(tags => this.tags.set(tags))
      }
    });
  }
}
