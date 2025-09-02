import {Component, computed, inject, input, InputSignal, Signal} from '@angular/core';
import {TagsControl} from '@components/tags-control';
import {BaseCharacter, CharacterListView, Characters, Tag} from '@api/characters';

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

  readonly characterInp: InputSignal<BaseCharacter | number> = input.required({alias: 'character'})

  private readonly character: Signal<Nullable<CharacterListView>> = computed(() => {
    const inp = this.characterInp()
    const all = this.characters.all()
    if (all.length === 0) {
      return null
    }

    const id = typeof inp === 'number' ? inp : inp.id
    return all.find((char) => char.id === id)
  })

  protected readonly id: Signal<number> = computed(() => this.character()?.id || 0)
  protected readonly name: Signal<string> = computed(() => this.character()?.name || '')
  protected readonly favorite: Signal<boolean> = computed(() => this.character()?.favorite || false)
  protected readonly tags: Signal<Tag[]> = computed(() => this.character()?.tags || [])
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => {
    const u = this.character()?.avatarUrl
    return !!u ? `url(${u})` : null;
  })
}
