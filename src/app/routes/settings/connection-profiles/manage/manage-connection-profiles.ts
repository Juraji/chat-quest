import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from "@components/page-header/page-header";
import {ActivatedRoute, RouterLink} from '@angular/router';
import {ConnectionProfile} from '@api/model';
import {routeDataSignal} from '@util/ng';
import {EmptyPipe} from '@components/empty-pipe';

@Component({
  selector: 'app-manage-connection-profiles',
  imports: [
    PageHeader,
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './manage-connection-profiles.html'
})
export class ManageConnectionProfiles {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly profiles: Signal<ConnectionProfile[]> = routeDataSignal(this.activatedRoute, 'profiles')
}
