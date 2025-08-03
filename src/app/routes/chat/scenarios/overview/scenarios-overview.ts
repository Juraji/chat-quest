import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from "@components/page-header";
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {NewItemCard} from '@components/cards/new-item-card/new-item-card';
import {Scenario} from '@api/model';
import {ScenarioCard} from '@components/cards/scenario-card/scenario-card';

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
