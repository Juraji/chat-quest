import {Component, computed, input, InputSignal, Signal} from '@angular/core';
import {Scenario} from '@api/model';

@Component({
  selector: 'app-scenario-card',
  imports: [],
  templateUrl: './scenario-card.html',
  host: {
    '[class.chat-quest-card]': 'true',
  }
})
export class ScenarioCard {
  readonly scenario: InputSignal<Scenario> = input.required()
  readonly name: Signal<string> = computed(() => this.scenario().name)
  readonly avatarUrl: Signal<Nullable<string>> = computed(() => this.scenario().avatarUrl)
}
