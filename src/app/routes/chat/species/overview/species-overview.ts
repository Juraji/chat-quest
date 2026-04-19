import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {Species} from '@api/species';
import {NewItemCard} from '@components/cards/new-item-card';
import {PageHeader} from '@components/page-header';
import {SpeciesCard} from '@components/cards/species-card';

@Component({
  selector: 'species-overview',
  imports: [
    NewItemCard,
    PageHeader,
    RouterLink,
    SpeciesCard
  ],
  templateUrl: './species-overview.html',
})
export class SpeciesOverview {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly species: Signal<Species[]> = routeDataSignal(this.activatedRoute, 'species');
}
