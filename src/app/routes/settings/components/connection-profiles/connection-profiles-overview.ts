import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {EmptyPipe} from '@components/empty.pipe';
import {ConnectionProfile} from '@api/providers';
import {routeDataSignal} from '@util/ng';

@Component({
  selector: 'connection-profiles-overview',
  imports: [
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './connection-profiles-overview.html'
})
export class ConnectionProfilesOverview {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly profiles: Signal<ConnectionProfile[]> = routeDataSignal(this.activatedRoute, 'profiles')
}
