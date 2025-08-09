import {Component, computed, input, InputSignal, Signal} from '@angular/core';
import {World} from '@api/worlds';

@Component({
  selector: 'world-card',
  imports: [],
  templateUrl: './world-card.html',
  host: {
    '[class.chat-quest-card]': 'true',
  }
})
export class WorldCard {
  readonly world: InputSignal<World> = input.required()
  readonly name: Signal<string> = computed(() => this.world().name)
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => {
    const u = this.world().avatarUrl
    return !!u ? `url(${u})` : null;
  })
}
