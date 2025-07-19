import {Component, computed, effect, input, InputSignal, signal, Signal, WritableSignal} from '@angular/core';
import {Character} from '@db/characters';
import {RouterLink} from '@angular/router';
import {TagsControl} from '@components/tags-control/tags-control';

@Component({
  selector: 'app-character-card',
  imports: [
    RouterLink,
    TagsControl
  ],
  templateUrl: './character-card.html',
  styleUrl: './character-card.scss',
  host: {
    '[class.favorite]': 'favorite()'
  }
})
export class CharacterCard {
  readonly character: InputSignal<Character> = input.required()
  readonly id: Signal<number> = computed(() => this.character().id)
  readonly name: Signal<string> = computed(() => this.character().name)
  readonly favorite: Signal<boolean> = computed(() => this.character().favorite)
  readonly tagIds: Signal<number[]> = computed(() => this.character().tagIds)
  readonly avatar: Signal<Blob | null> = computed(() => this.character().avatar)

  readonly avatarImageUrl: WritableSignal<string> = signal('')

  constructor() {
    effect(() => {
      const blob = this.avatar()
      this.avatarImageUrl.update(current => {
        if (!!current) URL.revokeObjectURL(current)
        if (!!blob) return URL.createObjectURL(blob)
        else return ''
      })
    });
  }
}
