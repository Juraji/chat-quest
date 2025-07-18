import {Component, computed, input, InputSignal, Signal} from '@angular/core';
import {Character} from '@db/characters';
import {RouterLink} from '@angular/router';

@Component({
  selector: 'app-character-card',
  imports: [
    RouterLink
  ],
  templateUrl: './character-card.html',
  styleUrl: './character-card.scss',
  host: {
    '[class.card]': 'true',
    '[class.favorite]': 'favorite()'
  }
})
export class CharacterCard {
  readonly character: InputSignal<Character> = input.required()
  readonly id: Signal<number> = computed(() => this.character().id)
  readonly name: Signal<string> = computed(() => this.character().name)
  readonly favorite: Signal<boolean> = computed(() => this.character().favorite)
  readonly avatar: Signal<Blob | null> = computed(() => this.character().avatar)
}
