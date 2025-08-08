import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from "@components/page-header";
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {NewItemCard} from '@components/cards/new-item-card';
import {ScenarioCard} from '@components/cards/scenario-card';
import {Scenario} from '@api/scenarios';

@Component({
  selector: 'app-scenarios-overview',
  imports: [
    PageHeader,
    RouterLink,
    NewItemCard,
    ScenarioCard
  ],
  templateUrl: './scenarios-overview.html'
})
export class ScenariosOverview {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios');
}
