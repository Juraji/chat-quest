import {Component, computed, input, InputSignal, Signal} from '@angular/core';
import {Species} from '@api/species';

@Component({
  selector: 'species-card',
  imports: [],
  templateUrl: './species-card.html',
  host: {
    '[class.item-card]': 'true',
  }
})
export class SpeciesCard {
  readonly species: InputSignal<Species> = input.required()
  readonly name: Signal<string> = computed(() => this.species().name)
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => {
    const u = this.species().avatarUrl
    return !!u ? `url(${u})` : null;
  })
}
