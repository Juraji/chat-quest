import {Component, computed, input, InputSignal, Signal} from '@angular/core';
import {Scenario} from '@api/scenarios';

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
  protected readonly avatarUrl: Signal<Nullable<string>> = computed(() => {
    const u = this.scenario().avatarUrl
    return !!u ? `url(${u})` : null;
  })
}
