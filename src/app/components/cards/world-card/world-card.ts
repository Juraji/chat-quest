import {Component, computed, input, InputSignal} from '@angular/core';
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
  readonly name = computed(() => this.world().name)
  readonly avatarUrl = computed(() => this.world().avatarUrl)
}
