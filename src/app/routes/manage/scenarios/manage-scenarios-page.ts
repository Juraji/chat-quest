import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {Scenario} from '@db/scenarios';
import {EmptyPipe} from '@components/empty-pipe';

@Component({
  selector: 'app-manage-scenarios-page',
  imports: [
    EmptyPipe,
    RouterLink
  ],
  templateUrl: './manage-scenarios-page.html'
})
export class ManageScenariosPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly scenarios: Signal<Scenario[]> = routeDataSignal(this.activatedRoute, 'scenarios');
}
