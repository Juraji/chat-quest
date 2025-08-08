import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from "@components/page-header";
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {NewItemCard} from '@components/cards/new-item-card';
import {World} from '@api/worlds';
import {WorldCard} from '@components/cards/world-card';

@Component({
  imports: [
    PageHeader,
    NewItemCard,
    RouterLink,
    WorldCard,

  ],
  templateUrl: './worlds-overview-page.html'
})
export class WorldsOverviewPage {
  private readonly activatedRoute = inject(ActivatedRoute);
  readonly worlds: Signal<World[]> = routeDataSignal(this.activatedRoute, 'worlds');
}
