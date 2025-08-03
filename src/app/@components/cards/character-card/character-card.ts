import {Component, computed, effect, inject, input, InputSignal, signal, Signal, WritableSignal} from '@angular/core';
import {Character, CharacterWithTags, Tag} from "@api/model"
import {Characters} from '@api/clients';
import {TagsControl} from '@components/tags-control/tags-control';

@Component({
  selector: 'app-character-card',
  imports: [
    TagsControl
  ],
  templateUrl: './character-card.html',
  styleUrl: './character-card.scss',
  host: {
    '[class.chat-quest-card]': 'true',
    '[class.favorite]': 'favorite()'
  }
})
export class CharacterCard {
  readonly characters = inject(Characters)

  readonly character: InputSignal<Character | CharacterWithTags> = input.required()
  protected readonly id: Signal<number> = computed(() => this.character().id)
  protected readonly name: Signal<string> = computed(() => this.character().name)
  protected readonly favorite: Signal<boolean> = computed(() => this.character().favorite)
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => this.character().avatarUrl)
  protected readonly tags: WritableSignal<Tag[]> = signal([])

  constructor() {
    effect(() => {
      const char = this.character()
      if ('tags' in char) {
        this.tags.set(char.tags)
      } else {
        this.characters
          .getTags(char.id)
          .subscribe(tags => this.tags.set(tags))
      }
    });
  }
}
